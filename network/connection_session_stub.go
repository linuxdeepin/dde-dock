package main

import (
	"dlib/dbus"
)

func (session *ConnectionSession) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Network",
		string(session.objPath),
		"com.deepin.daemon.ConnectionSession",
	}
}

// TODO
// func (session *ConnectionSession) updatePropCurrentUUID(v string) {
// 	session.CurrentUUID = v
// 	dbus.NotifyChange(session, "CurrentUUID")
// }

func (session *ConnectionSession) updatePropHasChanged(v bool) {
	session.HasChanged = v
	dbus.NotifyChange(session, "HasChanged")
}

func (session *ConnectionSession) updatePropCurrentFields() {
	// get fields through current page, show or hide some fields when
	// target fileds toggled

	// TODO processing logic

	session.CurrentFields = session.listFields(session.currentPage)
	dbus.NotifyChange(session, "CurrentFields")
}

func (session *ConnectionSession) updatePropCurrentErrors(v string) {
	// TODO
	dbus.NotifyChange(session, "CurrentErrors")
}
