package main

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

func (s *ConnectionSession) updatePropAvailablePages() {
	s.AvailablePages = s.listPages()
	dbus.NotifyChange(s, "AvailablePages")
}

func (s *ConnectionSession) updatePropAvailableKeys() {
	s.AvailableKeys = make(map[string][]string) // clear structure
	for _, page := range s.listPages() {
		s.AvailableKeys[page] = s.listKeys(page)
	}
	dbus.NotifyChange(s, "AvailableKeys")
}

func (s *ConnectionSession) updatePropErrors() {
	for _, page := range s.listPages() {
		s.Errors[page] = make(sectionErrors)
		sections := s.pageToSections(page)
		for _, section := range sections {
			// check error only section exists
			if isSettingSectionExists(s.data, section) {
				errs := generalCheckSettingValues(s.data, section)
				for k, v := range errs {
					s.Errors[page][k] = v
				}
			}
		}
	}
	// append errors when setting keys
	for page, pageErrors := range s.settingKeyErrors {
		for k, v := range pageErrors {
			s.Errors[page][k] = v
		}
	}
	dbus.NotifyChange(s, "Errors")
}
