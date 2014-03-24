package main

import (
	"dlib/dbus"
)

func (editor *ConnectionEditor) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Network",
		"/com/deepin/daemon/ConnectionEditor",
		"com.deepin.daemon.ConnectionEditor",
	}
}

func (editor *ConnectionEditor) updatePropCurrentUUID(v string) {
	editor.CurrentUUID = v
	dbus.NotifyChange(editor, "CurrentUUID")
}

func (editor *ConnectionEditor) updatePropHasChanged(v bool) {
	editor.HasChanged = v
	dbus.NotifyChange(editor, "HasChanged")
}

func (editor *ConnectionEditor) updatePropCurrentFields() {
	// get current fields through current page
	editor.CurrentFields = editor.listFields(editor.currentPage)
	dbus.NotifyChange(editor, "CurrentFields")
}

func (editor *ConnectionEditor) updatePropCurrentErrors(v string) {
	// TODO
	dbus.NotifyChange(editor, "CurrentErrors")
}
