package swapsched

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	login1 "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.login1"

	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/cgroup"
	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/log"
)

const (
	dbusServiceName = "com.deepin.daemon.SwapSchedHelper"
	dbusPath        = "/com/deepin/daemon/SwapSchedHelper"
	dbusInterface   = dbusServiceName
)

var logger = log.NewLogger("daemon/system/swapsched")

func init() {
	loader.Register(newDaemon(logger))
}

type Daemon struct {
	*loader.ModuleBase
	sessionWatcher *Helper
}

func newDaemon(logger *log.Logger) *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("swapsched", daemon, logger)
	return daemon
}

func (d *Daemon) GetDependencies() []string {
	return []string{}
}

func (d *Daemon) Start() error {
	err := cgroup.Init()
	if err != nil {
		return err
	}

	logger.Debug("swap sched helper start")
	sw, err := newHelper()
	if err != nil {
		return err
	}

	sw.init()
	d.sessionWatcher = sw

	service := loader.GetService()
	err = service.Export(dbusPath, sw)
	if err != nil {
		return err
	}

	err = service.RequestName(dbusServiceName)
	if err != nil {
		return err
	}

	return nil
}

func (d *Daemon) Stop() error {
	// TODO:
	return nil
}

type Helper struct {
	loginManager *login1.Manager
	sysSigLoop   *dbusutil.SignalLoop

	// nolint
	methods *struct {
		Prepare func() `in:"sessionID"`
	}
}

func (*Helper) GetInterfaceName() string {
	return dbusInterface
}

func newHelper() (*Helper, error) {
	systemBus, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	sysSigLoop := dbusutil.NewSignalLoop(systemBus, 10)
	sysSigLoop.Start()
	loginManager := login1.NewManager(systemBus)
	return &Helper{
		loginManager: loginManager,
		sysSigLoop:   sysSigLoop,
	}, nil
}

func (sw *Helper) Prepare(sessionID string) *dbus.Error {
	uid, err := sw.getSessionUid(sessionID)
	if err != nil {
		return dbusutil.ToError(err)
	}

	err = createDDECGroups(uid, sessionID)
	if err != nil {
		logger.Warning("failed to create cgroup:", err)
		return dbusutil.ToError(err)
	}

	return nil
}

func (sw *Helper) getSessionUid(sessionID string) (uint32, error) {
	sessions, err := sw.loginManager.ListSessions(0)
	if err != nil {
		return 0, err
	}

	for _, session := range sessions {
		if session.SessionId == sessionID {
			return session.UID, nil
		}
	}

	return 0, errors.New("not found session")
}

func (sw *Helper) init() {
	sw.loginManager.InitSignalExt(sw.sysSigLoop, true)
	_, err := sw.loginManager.ConnectSessionRemoved(
		func(sessionID string, sessionPath dbus.ObjectPath) {
			logger.Debug("session removed", sessionID, sessionPath)
			memMountPoint := cgroup.GetSubSysMountPoint(cgroup.Memory)
			_, err := os.Stat(filepath.Join(memMountPoint, sessionID+"@dde"))
			if err == nil {
				// path exit
				go func() {
					time.Sleep(10 * time.Second)
					err := deleteDDECGroups(sessionID)
					if err != nil {
						logger.Warning("failed to delete DDE cgroups:", err)
					}
				}()
			}
		})

	if err != nil {
		logger.Warning(err)
	}
}

func createDDECGroups(uid uint32, sessionID string) error {
	dir := sessionID + "@dde/"
	err := createCGroup(uid, dir+"uiapps")
	if err != nil {
		return err
	}

	err = createCGroup(uid, dir+"DE")
	if err != nil {
		return err
	}
	return nil
}

func createCGroup(uid uint32, name string) error {
	cg := newCgroup(name)
	uid0 := int(uid)
	cg.SetUidGid(uid0, uid0, uid0, uid0)
	logger.Debugf("create cgroup %s, uid: %d", name, uid)
	return cg.Create(false)
}

func deleteDDECGroups(sessionID string) error {
	logger.Debugf("delete cgroup for session %s", sessionID)
	cg := newCgroup(sessionID + "@dde")
	return cg.Delete(cgroup.DeleteFlagRecursive)
}

func newCgroup(name string) *cgroup.Cgroup {
	cg := cgroup.NewCgroup(name)
	cg.AddController(cgroup.Memory)
	cg.AddController(cgroup.Freezer)
	cg.AddController(cgroup.Blkio)
	return cg
}
