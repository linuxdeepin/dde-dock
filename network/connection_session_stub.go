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

func (s *ConnectionSession) updatePropAvailableVirtualSections() {
	s.AvailableVirtualSections = getAvailableVsections(s.Data)
	dbus.NotifyChange(s, "AvailableVirtualSections")
}

func (s *ConnectionSession) updatePropAvailableSections() {
	s.AvailableSections = getAvailableSections(s.Data)
	dbus.NotifyChange(s, "AvailableSections")
}

func (s *ConnectionSession) updatePropAvailableKeys() {
	s.AvailableKeys = make(map[string][]string) // clear structure
	for _, section := range getAvailableSections(s.Data) {
		s.AvailableKeys[section] = generalGetSettingAvailableKeys(s.Data, section)
	}
	dbus.NotifyChange(s, "AvailableKeys")
}

func (s *ConnectionSession) updatePropErrors() {
	for _, section := range getAvailableSections(s.Data) {
		s.Errors[section] = make(sectionErrors)
		if isSettingSectionExists(s.Data, section) {
			// check error only section exists
			errs := generalCheckSettingValues(s.Data, section)
			for k, v := range errs {
				s.Errors[section][k] = v
			}
		}
	}
	// append errors when setting keys
	for section, sectionErrors := range s.settingKeyErrors {
		for k, v := range sectionErrors {
			s.Errors[section][k] = v
		}
	}

	// TODO remove
	// for _, vsection := range getAvailableVsections(s.Data) {
	// 	s.Errors[vsection] = make(sectionErrors)
	// 	sections := getRelatedSectionsOfVsection(s.Data, vsection)
	// 	for _, section := range sections {
	// 		// check error only section exists
	// 		if isSettingSectionExists(s.Data, section) {
	// 			errs := generalCheckSettingValues(s.Data, section)
	// 			for k, v := range errs {
	// 				s.Errors[vsection][k] = v
	// 			}
	// 		}
	// 	}
	// }
	// // append errors when setting keys
	// for vsection, vsectionErrors := range s.settingKeyErrors {
	// 	for k, v := range vsectionErrors {
	// 		s.Errors[vsection][k] = v
	// 	}
	// }
	dbus.NotifyChange(s, "Errors")
}
