package dock

import (
	"pkg.deepin.io/lib/dbus"
)

func (s *Setting) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.daemon.Dock",
		ObjectPath: "/dde/dock/DockSetting",
		Interface:  "dde.dock.DockSetting",
	}
}
