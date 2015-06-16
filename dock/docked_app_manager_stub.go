package dock

import (
	"pkg.linuxdeepin.com/lib/dbus"
)

func (m *DockedAppManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.daemon.Dock",
		ObjectPath: "/dde/dock/DockedAppManager",
		Interface:  "dde.dock.DockedAppManager",
	}
}

func (m *DockedAppManager) destroy() {
	if m.core != nil {
		m.core.Unref()
	}
	dbus.UnInstallObject(m)
}
