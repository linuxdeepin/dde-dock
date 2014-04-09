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
		switch field {
		default:
			LOGGER.Error("updatePropErrors: invalid field name", field)
		case field8021x:
			session.Errors[page] = checkSetting8021xValues(session.data)
		case fieldConnection:
			session.Errors[page] = checkSettingConnectionValues(session.data)
		case fieldIPv4:
			session.Errors[page] = checkSettingIp4ConfigValues(session.data)
		case fieldIPv6:
			session.Errors[page] = checkSettingIp6ConfigValues(session.data)
		case fieldWired:
			session.Errors[page] = checkSettingWiredValues(session.data)
		case fieldWireless:
			session.Errors[page] = checkSettingWirelessValues(session.data)
		case fieldWirelessSecurity:
			session.Errors[page] = checkSettingWirelessSecurityValues(session.data)
		}
	}
	dbus.NotifyChange(session, "Errors")
}
