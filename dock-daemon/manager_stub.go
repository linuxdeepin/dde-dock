package main

import (
	"dlib/dbus"
)

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"dde.dock.EntryManager",
		"/dde/dock/EntryManager",
		"dde.dock.EntryManager",
	}
}
