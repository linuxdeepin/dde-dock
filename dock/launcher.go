package dock

import (
	"pkg.deepin.io/dde/daemon/appinfo"
	"pkg.deepin.io/lib/dbus"
)

var dbusConn, _ = dbus.SessionBus()

const (
	launcherDest    = "com.deepin.dde.daemon.Launcher"
	launcherObjPath = "/com/deepin/dde/daemon/Launcher"
	fullMethodName  = "com.deepin.dde.daemon.Launcher.MarkLaunched"
)

func markAsLaunched(appId string) {
	if dbusConn == nil {
		return
	}
	obj := dbusConn.Object(launcherDest, dbus.ObjectPath(launcherObjPath))
	obj.Call(fullMethodName, 0, appId)
}

func recordFrequency(appId string) {
	f, err := appinfo.GetFrequencyRecordFile()
	if err == nil {
		appinfo.SetFrequency(appId, appinfo.GetFrequency(appId, f)+1, f) // FIXME: DesktopID???
		f.Free()
	}
}
