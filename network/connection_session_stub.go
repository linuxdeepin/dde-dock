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
func (session *ConnectionSession) updatePropAllowSave(v bool) {
	session.AllowSave = v
	dbus.NotifyChange(session, "AllowSave")
}

func (session *ConnectionSession) updatePropAvailableKeys() {
	session.AvailableKeys = make(map[string][]string) // clear structure
	for _, page := range session.ListPages() {
		session.AvailableKeys[page] = session.listKeys(page)
	}
	dbus.NotifyChange(session, "AvailableKeys")
}

func (session *ConnectionSession) updatePropErrors() {
	for _, page := range session.ListPages() {
		field := session.pageToField(page)
		if isSettingFieldExists(session.data, field) { // TODO
			session.Errors[page] = generalCheckSettingValues(session.data, field)
		}
	}
	dbus.NotifyChange(session, "Errors")
}
