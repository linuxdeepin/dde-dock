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
		switch page {
		default:
			LOGGER.Error("updatePropErrors: invalid page name", page)
		case pageGeneral:
			session.Errors[page] = checkSettingConnectionValues(session.data)
		case pageEthernet:
			session.Errors[page] = checkSettingWiredValues(session.data)
		case pageWifi:
			session.Errors[page] = checkSettingWirelessValues(session.data)
		case pageIPv4:
			session.Errors[page] = checkSettingIp4ConfigValues(session.data)
		case pageIPv6:
			session.Errors[page] = checkSettingIp6ConfigValues(session.data)
		case pageSecurity: // TODO
			switch session.connectionType {
			case typeWired:
			case typeWireless:
				// switch method {
				// session.Errors[page] = checkSetting8021xValues(session.data)
				// session.Errors[page] = checkSettingWirelessSecurityValues(session.data)
				// }
			}
		}
	}
	dbus.NotifyChange(session, "Errors")
}
