package main

import "dlib/logger"

import "dlib"
import "dlib/dbus"
import "os"
import "dde-daemon/accounts"
import "dde-daemon/systeminfo"

type Manager struct{}

var Logger = logger.NewLogger("com.deepin.daemon")
var _manager *Manager

func GetManager() *Manager {
	if _manager == nil {
		_manager = &Manager{}
	}

	return _manager
}

func (obj *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon",
		"/com/deepin/daemon",
		"com.deepin.daemon",
	}
}

func main() {
	if !dlib.UniqueOnSystem("com.deepin.daemon") {
		Logger.Warning("There already has an dde-daemon running.")
		return
	}

	defer Logger.EndTracing()
	Logger.SetRestartCommand("/usr/lib/deepin-daemon/dde-system-daemon")

	if err := dbus.InstallOnSystem(GetManager()); err != nil {
		Logger.Error("Install DBus Failed:", err)
		panic(err)
	}
	dbus.DealWithUnhandledMessage()

	go accounts.Start()
	go systeminfo.Start()

	if err := dbus.Wait(); err != nil {
		Logger.Error("dde-daemon lost dbus")
		os.Exit(-1)
	} else {
		os.Exit(0)
	}
}
