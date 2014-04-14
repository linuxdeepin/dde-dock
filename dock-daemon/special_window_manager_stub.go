package main

import (
	"dlib/dbus"
)

func (m *SpecialWindowManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Dock",
		"/dde/dock/SpecialWindowManager",
		"dde.dock.SpecialWindowManager",
	}
}
