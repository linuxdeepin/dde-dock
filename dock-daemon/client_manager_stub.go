package main

import (
	"dlib/dbus"
)

func (m *ClientManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Dock",
		"/dde/dock/ClientManager",
		"dde.dock.ClientManager",
	}
}
