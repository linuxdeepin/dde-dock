package sessionwatcher

import (
	"pkg.deepin.io/lib/dbus"
)

const (
	dbusDest = "com.deepin.daemon.SessionWatcher"
	dbusPath = "/com/deepin/daemon/SessionWatcher"
	dbusIFC  = dbusDest
)

func (*Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}
