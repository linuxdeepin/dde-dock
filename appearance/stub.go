package appearance

import (
	"pkg.deepin.io/lib/dbus"
)

const (
	dbusDest = "com.deepin.daemon.Appearance"
	dbusPath = "/com/deepin/daemon/Appearance"
	dbusIFC  = "com.deepin.daemon.Appearance"
)

func (*Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}
