package lastore

import (
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/linuxdeepin/go-dbus-factory/com.deepin.lastore"
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.system.power"
	ofdbus "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.dbus"
	"github.com/linuxdeepin/go-dbus-factory/org.freedesktop.notifications"

	"pkg.deepin.io/dde/daemon/common/dsync"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/proxy"
	"pkg.deepin.io/lib/gettext"
)

type Lastore struct {
	service        *dbusutil.Service
	sysSigLoop     *dbusutil.SignalLoop
	sessionSigLoop *dbusutil.SignalLoop
	jobStatus      map[dbus.ObjectPath]CacheJobInfo
	lang           string
	inhibitFd      dbus.UnixFD

	power         *power.Power
	core          *lastore.Lastore
	sysDBusDaemon *ofdbus.DBus
	notifications *notifications.Notifications

	syncConfig *dsync.Config

	notifiedBattery     bool
	notifyIdHidMap      map[uint32]dbusutil.SignalHandlerId
	lastoreRule         dbusutil.MatchRule
	jobsPropsChangedHId dbusutil.SignalHandlerId

	// prop:
	PropsMu            sync.RWMutex
	SourceCheckEnabled bool

	methods *struct {
		SetSourceCheckEnabled func() `in:"val"`
		IsDiskSpaceSufficient func() `out:"result"`
	}
}

type CacheJobInfo struct {
	Id       string
	Status   Status
	Name     string
	Progress float64
	Type     string
}

func newLastore(service *dbusutil.Service) (*Lastore, error) {
	l := &Lastore{
		service:   service,
		jobStatus: make(map[dbus.ObjectPath]CacheJobInfo),
		inhibitFd: -1,
		lang:      QueryLang(),
	}

	logger.Debugf("CurrentLang: %q", l.lang)
	systemBus, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	sessionBus, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}

	l.sysSigLoop = dbusutil.NewSignalLoop(systemBus, 100)
	l.sysSigLoop.Start()

	l.sessionSigLoop = dbusutil.NewSignalLoop(sessionBus, 10)
	l.sessionSigLoop.Start()

	l.initCore(systemBus)
	l.initNotify(sessionBus)
	l.initSysDBusDaemon(systemBus)
	l.initPower(systemBus)

	l.syncConfig = dsync.NewConfig("updater", &syncConfig{l: l},
		l.sessionSigLoop, dbusPath, logger)
	return l, nil
}

func (l *Lastore) initPower(systemBus *dbus.Conn) {
	l.power = power.NewPower(systemBus)
	l.power.InitSignalExt(l.sysSigLoop, true)
	err := l.power.HasBattery().ConnectChanged(func(hasValue bool, hasBattery bool) {
		if !hasBattery {
			l.notifiedBattery = false
		}
	})
	if err != nil {
		logger.Warning(err)
	}
}

func (l *Lastore) initNotify(sessionBus *dbus.Conn) {
	l.notifications = notifications.NewNotifications(sessionBus)
	l.notifications.InitSignalExt(l.sessionSigLoop, true)

	l.notifyIdHidMap = make(map[uint32]dbusutil.SignalHandlerId)
	_, err := l.notifications.ConnectNotificationClosed(func(id uint32, reason uint32) {
		logger.Debug("notification closed id", id)
		hid, ok := l.notifyIdHidMap[id]
		if ok {
			logger.Debugf("remove id: %d, hid: %d", id, hid)
			delete(l.notifyIdHidMap, id)

			// delay call removeHandler
			time.AfterFunc(100*time.Millisecond, func() {
				l.notifications.RemoveHandler(hid)
			})
		}
	})
	if err != nil {
		logger.Warning(err)
	}
}

func (l *Lastore) initSysDBusDaemon(systemBus *dbus.Conn) {
	l.sysDBusDaemon = ofdbus.NewDBus(systemBus)
	l.sysDBusDaemon.InitSignalExt(l.sysSigLoop, true)
	_, err := l.sysDBusDaemon.ConnectNameOwnerChanged(
		func(name string, oldOwner string, newOwner string) {
			if name == l.core.ServiceName_() {
				if newOwner == "" {
					l.offline()
				} else {
					l.online()
				}
			}
		})
	if err != nil {
		logger.Warning(err)
	}
}

func (l *Lastore) initCore(systemBus *dbus.Conn) {
	l.core = lastore.NewLastore(systemBus)
	l.lastoreRule = dbusutil.NewMatchRuleBuilder().
		Sender(l.core.ServiceName_()).
		Type("signal").
		Interface("org.freedesktop.DBus.Properties").
		Member("PropertiesChanged").
		Build()
	err := l.lastoreRule.AddTo(systemBus)
	if err != nil {
		logger.Warning(err)
	}

	l.core.InitSignalExt(l.sysSigLoop, false)
	err = l.core.JobList().ConnectChanged(func(hasValue bool, value []dbus.ObjectPath) {
		if !hasValue {
			return
		}
		l.updateJobList(value)
	})
	if err != nil {
		logger.Warning(err)
	}

	l.jobsPropsChangedHId = l.sysSigLoop.AddHandler(&dbusutil.SignalRule{
		Name: "org.freedesktop.DBus.Properties.PropertiesChanged",
	}, func(sig *dbus.Signal) {
		if len(sig.Body) != 3 {
			return
		}
		props, _ := sig.Body[1].(map[string]dbus.Variant)
		ifc, _ := sig.Body[0].(string)
		if ifc == "com.deepin.lastore.Job" {
			l.updateCacheJobInfo(sig.Path, props)
		}
	})

	jobList, err := l.core.JobList().Get(0)
	if err != nil {
		logger.Warning(err)
	}

	l.updateJobList(jobList)
}

func (l *Lastore) destroy() {
	l.sessionSigLoop.Stop()
	l.sysSigLoop.Stop()
	l.syncConfig.Destroy()

	systemBus := l.sysSigLoop.Conn()
	err := l.lastoreRule.RemoveFrom(systemBus)
	if err != nil {
		logger.Warning(err)
	}

	l.sysSigLoop.RemoveHandler(l.jobsPropsChangedHId)
	l.power.RemoveHandler(proxy.RemoveAllHandlers)
	l.core.RemoveHandler(proxy.RemoveAllHandlers)
	l.notifications.RemoveHandler(proxy.RemoveAllHandlers)
}

func (l *Lastore) GetInterfaceName() string {
	return "com.deepin.LastoreSessionHelper"
}

// updateJobList clean invalid cached Job status
// The list is the newest JobList.
func (l *Lastore) updateJobList(list []dbus.ObjectPath) {
	var invalids []dbus.ObjectPath
	for jobPath := range l.jobStatus {
		safe := false
		for _, p := range list {
			if p == jobPath {
				safe = true
				break
			}
		}
		if !safe {
			invalids = append(invalids, jobPath)
		}
	}
	for _, jobPath := range invalids {
		delete(l.jobStatus, jobPath)
	}
	logger.Debugf("UpdateJobList: %v - %v", list, invalids)
}

func (l *Lastore) offline() {
	logger.Info("Lastore.Daemon Offline")
	l.jobStatus = make(map[dbus.ObjectPath]CacheJobInfo)
}

func (l *Lastore) online() {
	logger.Info("Lastore.Daemon Online")
}

func (l *Lastore) createJobFailedActions(jobId string) []NotifyAction {
	ac := []NotifyAction{
		{
			Id:   "retry",
			Name: gettext.Tr("Retry"),
			Callback: func() {
				err := l.core.StartJob(dbus.FlagNoAutoStart, jobId)
				logger.Infof("StartJob %q : %v", jobId, err)
			},
		},
		{
			Id:   "cancel",
			Name: gettext.Tr("Cancel"),
			Callback: func() {
				err := l.core.CleanJob(dbus.FlagNoAutoStart, jobId)
				logger.Infof("CleanJob %q : %v", jobId, err)
			},
		},
	}
	return ac
}

func (l *Lastore) createUpdateActions() []NotifyAction {
	ac := []NotifyAction{
		{
			Id:   "update",
			Name: gettext.Tr("Update Now"),
			Callback: func() {
				go func() {
					err := exec.Command("dde-control-center","-m", "update", "-p", "Checking").Run()
					if err != nil {
						logger.Warningf("createUpdateActions: %v",err)
					}
				}()
			},
		},
	}


	return ac
}

func (l *Lastore) notifyJob(path dbus.ObjectPath) {
	l.checkBattery()

	info := l.jobStatus[path]
	status := info.Status
	logger.Debugf("notifyJob: %q %q --> %v", path, status, info)
	switch guestJobTypeFromPath(path) {
	case InstallJobType:
		switch status {
		case FailedStatus:
			l.notifyInstall(info.Name, false, l.createJobFailedActions(info.Id))
		case SucceedStatus:
			if info.Progress == 1 {
				l.notifyInstall(info.Name, true, nil)
			}
		}
	case RemoveJobType:
		switch status {
		case FailedStatus:
			l.notifyRemove(info.Name, false, l.createJobFailedActions(info.Id))
		case SucceedStatus:
			l.notifyRemove(info.Name, true, nil)
		}

	case CleanJobType:
		if status == SucceedStatus &&
			strings.Contains(info.Name, "+notify") {
			l.notifyAutoClean()
		}
	case UpdateSourceJobType:
		val, _ := l.core.UpdatablePackages().Get(0)
		if status == EndStatus && len(val) > 0 {
			l.notifyUpdateSource(l.createUpdateActions())
		}
	}
}

func (*Lastore) IsDiskSpaceSufficient() (bool, *dbus.Error) {
	avail, err := queryVFSAvailable("/")
	if err != nil {
		return false, dbusutil.ToError(err)
	}
	return avail > 1024*1024*10 /* 10 MB */, nil
}

func (l *Lastore) updateCacheJobInfo(path dbus.ObjectPath, props map[string]dbus.Variant) {
	info := l.jobStatus[path]
	oldStatus := info.Status

	systemBus := l.sysSigLoop.Conn()
	job, err := lastore.NewJob(systemBus, path)
	if err != nil {
		logger.Warning(err)
		return
	}

	if info.Id == "" {
		if v, ok := props["Id"]; ok {
			info.Id, _ = v.Value().(string)
		}
		if info.Id == "" {
			id, _ := job.Id().Get(dbus.FlagNoAutoStart)
			info.Id = id
		}
	}

	if info.Name == "" {
		if v, ok := props["Name"]; ok {
			info.Name, _ = v.Value().(string)
		}
		if info.Name == "" {
			name, _ := job.Name().Get(dbus.FlagNoAutoStart)

			if name == "" {
				pkgs, _ := job.Packages().Get(dbus.FlagNoAutoStart)
				if len(pkgs) == 0 {
					name = "unknown"
				} else {
					name = PackageName(pkgs[0], l.lang)
				}
			}

			info.Name = name
		}
	}

	if v, ok := props["Progress"]; ok {
		info.Progress, _ = v.Value().(float64)
	}

	if v, ok := props["Status"]; ok {
		status := v.Value().(string)
		info.Status = Status(status)
	}

	if info.Type == "" {
		if v, ok := props["Type"]; ok {
			info.Type, _ = v.Value().(string)
		}
		if info.Type == "" {
			info.Type, _ = job.Type().Get(dbus.FlagNoAutoStart)
		}
	} else {
		if v, ok := props["Type"]; ok {
			info.Type, _ = v.Value().(string)
		}
	}
	l.jobStatus[path] = info
	logger.Debugf("updateCacheJobInfo: %#v", info)

	if oldStatus != info.Status {
		l.notifyJob(path)
	}
}

// guestJobTypeFromPath guest the JobType from object path
// We can't get the JobType when the DBusObject destroyed.
func guestJobTypeFromPath(path dbus.ObjectPath) string {
	_path := string(path)
	for _, jobType := range []string{
		// job types:
		InstallJobType, DownloadJobType, RemoveJobType,
		UpdateSourceJobType, DistUpgradeJobType, CleanJobType,
	} {
		if strings.Contains(_path, jobType) {
			return jobType
		}
	}
	return ""
}

var MinBatteryPercent = 30.0

func (l *Lastore) checkBattery() {
	if l.notifiedBattery {
		return
	}
	hasBattery, _ := l.power.HasBattery().Get(0)
	onBattery, _ := l.power.OnBattery().Get(0)
	percent, _ := l.power.BatteryPercentage().Get(0)
	if hasBattery && onBattery && percent <= MinBatteryPercent {
		l.notifiedBattery = true
		l.notifyLowPower()
	}
}
