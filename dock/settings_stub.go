package dock

import (
	"dlib/dbus"
)

func (s *Setting) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Dock",
		"/dde/dock/DockSetting",
		"dde.dock.DockSetting",
	}
}
