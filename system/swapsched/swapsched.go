package swapsched

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"dbus/org/freedesktop/login1"

	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

const (
	loginDest         = "org.freedesktop.login1"
	loginObjPath      = "/org/freedesktop/login1"
	cGroupControllers = "memory,freezer"
	cGroupRoot        = "/sys/fs/cgroup"
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
	logger.Debug("swap sched helper start")
	sw := newHelper()
	sw.init()
	d.sessionWatcher = sw

	err := dbus.InstallOnSystem(sw)
	if err != nil {
		logger.Warning(err)
		return err
	}
	dbus.DealWithUnhandledMessage()
	return nil
}

func (d *Daemon) Stop() error {
	// TODO:
	return nil
}

type Helper struct {
	loginManager *login1.Manager
}

func (sw *Helper) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.daemon.SwapSchedHelper",
		ObjectPath: "/com/deepin/daemon/SwapSchedHelper",
		Interface:  "com.deepin.daemon.SwapSchedHelper",
	}
}

func newHelper() *Helper {
	loginManager, err := login1.NewManager(loginDest, loginObjPath)
	if err != nil {
		panic(err)
	}
	return &Helper{
		loginManager: loginManager,
	}
}

func (sw *Helper) Prepare(sessionID string) error {
	username, err := sw.getSessionUsername(sessionID)
	if err != nil {
		return err
	}

	err = createDDECGroups(username, sessionID)
	if err != nil {
		logger.Warning("failed to create cgroup:", err)
		return err
	}

	return nil
}

func (sw *Helper) getSessionUsername(sessionID0 string) (string, error) {
	sessions, err := sw.loginManager.ListSessions()
	if err != nil {
		return "", err
	}

	for _, session := range sessions {
		// session fields: sessionID, uid, username, seat, sessionObjPath
		if len(session) < 3 {
			return "", errors.New("len(session) < 3")
		}

		sessionID, ok := session[0].(string)
		if !ok {
			return "", errors.New("type of session[0] is not string")
		}

		username, ok := session[2].(string)
		if !ok {
			return "", errors.New("type of session[2] is not string")
		}

		if sessionID == sessionID0 {
			return username, nil
		}
	}

	return "", errors.New("not found session")
}

func (sw *Helper) init() {
	sw.loginManager.ConnectSessionRemoved(func(sessionID string, sessionObjPath dbus.ObjectPath) {
		logger.Debug("session removed", sessionID, sessionObjPath)
		go func() {
			time.Sleep(time.Second * 10)
			_, err := os.Stat(filepath.Join(cGroupRoot, "memory", sessionID+"@dde"))
			if err == nil {
				// path exist
				err = deleteDDECGroups(sessionID)
				if err != nil {
					logger.Warning("failed to delete cgroup:", err)
				}
			}
		}()
	})
}

func createDDECGroups(username, sessionID string) error {
	user := username + ":" + username
	dir := sessionID + "@dde/"

	err := createCGroup(user, dir+"uiapps")
	if err != nil {
		return err
	}

	err = createCGroup(user, dir+"DE")
	if err != nil {
		return err
	}
	return nil
}

func createCGroup(user, path string) error {
	cmdline := fmt.Sprintf("cgcreate -t %s -a %s -g %s:%s", user, user, cGroupControllers, path)
	logger.Debug("exec cmd:", cmdline)
	cmd := exec.Command("cgcreate", "-t", user, "-a", user, "-g", cGroupControllers+":"+path)
	out, err := cmd.CombinedOutput()
	if len(out) > 0 {
		logger.Debugf("cgcreate output: %s", out)
	}
	return err
}

func deleteDDECGroups(sessionID string) error {
	path := sessionID + "@dde"
	cmdline := fmt.Sprintf("cgdelete -r -g %s:%s", cGroupControllers, path)
	logger.Debug("exec cmd:", cmdline)
	cmd := exec.Command("cgdelete", "-r", "-g", cGroupControllers+":"+path)
	out, err := cmd.CombinedOutput()
	if len(out) > 0 {
		logger.Debugf("cgdelete output: %s", out)
	}
	return err
}
