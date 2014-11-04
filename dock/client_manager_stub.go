package dock

import (
	"pkg.linuxdeepin.com/lib/dbus"
)

func (m *ClientManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.daemon.Dock",
		ObjectPath: "/dde/dock/ClientManager",
		Interface:  "dde.dock.ClientManager",
	}
}
