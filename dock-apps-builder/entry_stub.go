package main

import "crypto/md5"
import "encoding/hex"

import "dlib/dbus"

func (e *AppEntry) GetDBusInfo() dbus.DBusInfo {
	hasher := md5.New()
	hasher.Write([]byte(e.Id))
	// DBusObjectPath can't be start with digital number
	id := "d" + hex.EncodeToString(hasher.Sum(nil))
	return dbus.DBusInfo{
		"dde.dock.entry." + id,
		"/dde/dock/entry/v1/" + id,
		"dde.dock.Entry",
	}
}
