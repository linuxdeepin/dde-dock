package main

import (
	"dlib/dbus"
)

func (m *DockedAppManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"dde.dock.EntryManager",
		"/dde/dock/DockedAppManager",
		"dde.dock.DockedAppManager",
	}
}
