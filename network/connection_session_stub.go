package network

import (
	"dlib/dbus"
	"fmt"
)

func (s *ConnectionSession) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Network",
		string(s.sessionPath),
		"com.deepin.daemon.ConnectionSession",
	}
}

func (s *ConnectionSession) updateProps() {
	s.updatePropType()
	s.updatePropAvailableVirtualSections()
	s.updatePropAvailableSections()
	s.updatePropAvailableKeys()
	s.updatePropErrors()

	// update Data property at end, for that this was used by font-end
	// to update widget value that with proeprty "alwaysUpdate", which
	// should only update value when visible, so it depends on
	// "AvailableSections" and "AvailableKeys"
	s.updatePropData()
}

func (s *ConnectionSession) updatePropData() {
	dbus.NotifyChange(s, "Data")
}

func (s *ConnectionSession) updatePropType() {
	s.Type = getCustomConnectionType(s.Data)
	dbus.NotifyChange(s, "Type")
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
	s.Errors = make(sessionErrors)
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

	// append errors when setting key
	for section, sectionErrors := range s.settingKeyErrors {
		for k, v := range sectionErrors {
			s.Errors[section][k] = v
		}
	}

	// check if vpn missing plugin
	switch getCustomConnectionType(s.Data) {
	case connectionVpnL2tp, connectionVpnOpenconnect, connectionVpnPptp, connectionVpnVpnc, connectionVpnOpenvpn:
		if isKeyAvailable(s.Data, vsectionVpn, NM_SETTING_VK_VPN_MISSING_PLUGIN) {
			if _, ok := s.Errors[vsectionVpn]; ok {
				s.Errors[vsectionVpn][NM_SETTING_VK_VPN_MISSING_PLUGIN] =
					fmt.Sprintf(NM_KEY_ERROR_MISSING_DEPENDS_PACKAGE, getSettingVkVpnMissingPlugin(s.Data))
			} else {
				logger.Errorf("missing section, errors[%s]", vsectionVpn)
			}
		}
	}

	dbus.NotifyChange(s, "Errors")
}
