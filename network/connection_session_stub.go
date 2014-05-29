package network

import (
	"dlib/dbus"
)

func (s *ConnectionSession) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Network",
		string(s.sessionPath),
		"com.deepin.daemon.ConnectionSession",
	}
}

func (s *ConnectionSession) updatePropConnectionType() {
	dbus.NotifyChange(s, "ConnectionType")
}

func (s *ConnectionSession) updatePropAvailableSections() {
	s.AvailableSections = getAvailableVsections(s.data)
	dbus.NotifyChange(s, "AvailableSections")
}

func (s *ConnectionSession) updatePropAvailableKeys() {
	s.AvailableKeys = make(map[string][]string) // clear structure
	for _, vsection := range getAvailableVsections(s.data) {
		s.AvailableKeys[vsection] = getAvailableKeysOfVsection(s.data, vsection)
	}
	dbus.NotifyChange(s, "AvailableKeys")
}

func (s *ConnectionSession) updatePropErrors() {
	for _, vsection := range getAvailableVsections(s.data) {
		s.Errors[vsection] = make(sectionErrors)
		sections := getRelatedSectionsOfVsection(s.data, vsection)
		for _, section := range sections {
			// check error only section exists
			if isSettingSectionExists(s.data, section) {
				errs := generalCheckSettingValues(s.data, section)
				for k, v := range errs {
					s.Errors[vsection][k] = v
				}
			}
		}
	}
	// append errors when setting keys
	for vsection, vsectionErrors := range s.settingKeyErrors {
		for k, v := range vsectionErrors {
			s.Errors[vsection][k] = v
		}
	}
	dbus.NotifyChange(s, "Errors")
}
