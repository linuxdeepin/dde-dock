package dock

import (
	"pkg.linuxdeepin.com/lib/dbus"
)

func (m *DockedAppManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Dock",
		"/dde/dock/DockedAppManager",
		"dde.dock.DockedAppManager",
	}
}
