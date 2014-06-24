package dock

import (
	"pkg.linuxdeepin.com/lib/dbus"
)

func (s *Setting) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Dock",
		"/dde/dock/DockSetting",
		"dde.dock.DockSetting",
	}
}
