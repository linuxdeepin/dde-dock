package swapsched

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"dbus/org/freedesktop/login1"

	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/cgroup"
	oldDBusLib "pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/log"
)

const (
	loginServiceName = "org.freedesktop.login1"
	loginObjPath     = "/org/freedesktop/login1"

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
	sw := newHelper()
	sw.init()
	d.sessionWatcher = sw

	service := loader.GetService()
	err = service.Export(sw)
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

	methods *struct {
		Prepare func() `in:"sessionID"`
	}
}

func (sw *Helper) GetDBusExportInfo() dbusutil.ExportInfo {
	return dbusutil.ExportInfo{
		Path:      dbusPath,
		Interface: dbusInterface,
	}
}

func newHelper() *Helper {
	loginManager, err := login1.NewManager(loginServiceName, loginObjPath)
	if err != nil {
		panic(err)
	}
	return &Helper{
		loginManager: loginManager,
	}
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

func (sw *Helper) getSessionUid(sessionID0 string) (uint32, error) {
	sessions, err := sw.loginManager.ListSessions()
	if err != nil {
		return 0, err
	}

	for _, session := range sessions {
		// session fields: sessionID, uid, username, seat, sessionObjPath
		if len(session) < 2 {
			return 0, errors.New("len(session) < 3")
		}

		sessionID, ok := session[0].(string)
		if !ok {
			return 0, errors.New("type of session[0] is not string")
		}

		uid, ok := session[1].(uint32)
		if !ok {
			return 0, errors.New("type of session[1] is not uint32")
		}

		if sessionID == sessionID0 {
			return uid, nil
		}
	}

	return 0, errors.New("not found session")
}

func (sw *Helper) init() {
	sw.loginManager.ConnectSessionRemoved(func(sessionID string, sessionObjPath oldDBusLib.ObjectPath) {
		logger.Debug("session removed", sessionID, sessionObjPath)
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
