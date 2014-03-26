package main

import (
	"dlib/dbus"
)

func (s *Setting) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"dde.dock.Daemon",
		"/dde/dock/DockSetting",
		"dde.dock.DockSetting",
	}
}
