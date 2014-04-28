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

// TODO
func (s *ConnectionSession) updatePropAllowSave(v bool) {
	s.AllowSave = v
	dbus.NotifyChange(s, "AllowSave")
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
		s.Errors[page] = make(FieldKeyErrors)
		fields := s.pageToFields(page)
		for _, field := range fields {
			if isSettingFieldExists(s.data, field) { // TODO
				errs := generalCheckSettingValues(s.data, field)
				for k, v := range errs {
					s.Errors[page][k] = v
				}
			}
		}
	}
	dbus.NotifyChange(s, "Errors")
}
