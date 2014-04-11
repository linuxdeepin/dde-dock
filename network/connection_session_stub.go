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
		session.Errors[page] = make(map[string]string)
		fields := session.pageToFields(page)
		for _, field := range fields {
			if isSettingFieldExists(session.data, field) { // TODO
				errs := generalCheckSettingValues(session.data, field)
				for k, v := range errs {
					session.Errors[page][k] = v
				}
			}
		}
	}
	dbus.NotifyChange(session, "Errors")
}
