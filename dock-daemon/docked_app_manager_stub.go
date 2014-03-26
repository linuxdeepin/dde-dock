package main

import (
	"dlib/dbus"
)

func (m *DockedAppManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"dde.dock.Daemon",
		"/dde/dock/DockedAppManager",
		"dde.dock.DockedAppManager",
	}
}
