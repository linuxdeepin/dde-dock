package main

import (
	"dlib/dbus"
)

func (e *EntryProxyer) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"dde.dock.EntryManager",
		entryPathPrefix + e.entryId,
		"dde.dock.EntryProxyer",
	}
}
