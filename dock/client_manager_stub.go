package dock

import (
	"pkg.linuxdeepin.com/lib/dbus"
)

func (m *ClientManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Dock",
		"/dde/dock/ClientManager",
		"dde.dock.ClientManager",
	}
}
