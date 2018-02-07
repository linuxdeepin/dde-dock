/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package network

import (
	"fmt"
	"pkg.deepin.io/dde/daemon/network/nm"
	"pkg.deepin.io/lib/dbus"
)

func (s *ConnectionSession) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.daemon.Network",
		ObjectPath: string(s.sessionPath),
		Interface:  "com.deepin.daemon.ConnectionSession",
	}
}

func (s *ConnectionSession) setProps() {
	correctConnectionData(s.data)
	s.setPropType()
	s.setPropAllowDelete()
	s.setPropAvailableVirtualSections()
	s.setPropAvailableSections()
	s.setPropAvailableKeys()
	s.setPropErrors()

	// notify connection data changed at end, for that this was used
	// by font-end to update widget value that with proeprty
	// "alwaysUpdate", which should only update value when visible, so
	// it depends on "AvailableSections" and "AvailableKeys"
	dbus.Emit(s, "ConnectionDataChanged")
}

func (s *ConnectionSession) setPropType() {
	s.Type = getCustomConnectionType(s.data)
	dbus.NotifyChange(s, "Type")
}

func (s *ConnectionSession) setPropAllowDelete() {
	//any connection was allowed to deleted
	//if !s.connectionExists || (isNmObjectPathValid(s.devPath) &&
	//	nmGeneralGetDeviceUniqueUuid(s.devPath) == s.Uuid) {
	//	s.AllowDelete = false
	//} else {
	//	s.AllowDelete = true
	//}
	s.AllowDelete = true
	dbus.NotifyChange(s, "AllowDelete")
}

func (s *ConnectionSession) setPropAvailableVirtualSections() {
	s.AvailableVirtualSections = getAvailableVsections(s.data)
	dbus.NotifyChange(s, "AvailableVirtualSections")
}

func (s *ConnectionSession) setPropAvailableSections() {
	s.AvailableSections = getAvailableSections(s.data)
	dbus.NotifyChange(s, "AvailableSections")
}

func (s *ConnectionSession) setPropAvailableKeys() {
	s.AvailableKeys = make(map[string][]string) // clear structure
	for _, section := range getAvailableSections(s.data) {
		s.AvailableKeys[section] = generalGetSettingAvailableKeys(s.data, section)
	}
	dbus.NotifyChange(s, "AvailableKeys")
}

func (s *ConnectionSession) setPropErrors() {
	s.Errors = make(sessionErrors)
	for _, section := range getAvailableSections(s.data) {
		s.Errors[section] = make(sectionErrors)
		if isSettingExists(s.data, section) {
			// check error only section exists
			errs := generalCheckSettingValues(s.data, section)
			for k, v := range errs {
				s.Errors[section][k] = v
			}
		}
	}

	// clear setting key errors that not available
	for errSection, _ := range s.settingKeyErrors {
		if !isStringInArray(errSection, s.AvailableSections) {
			delete(s.settingKeyErrors, errSection)
		} else {
			for errKey, _ := range s.settingKeyErrors[errSection] {
				if !isStringInArray(errKey, s.AvailableKeys[errSection]) {
					delete(s.settingKeyErrors[errSection], errKey)
				}
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
	switch getCustomConnectionType(s.data) {
	case connectionVpnL2tp, connectionVpnOpenconnect, connectionVpnPptp, connectionVpnVpnc, connectionVpnOpenvpn:
		if isKeyAvailable(s.data, nm.NM_SETTING_VS_VPN, nm.NM_SETTING_VK_VPN_MISSING_PLUGIN) {
			if _, ok := s.Errors[nm.NM_SETTING_VS_VPN]; ok {
				s.Errors[nm.NM_SETTING_VS_VPN][nm.NM_SETTING_VK_VPN_MISSING_PLUGIN] =
					fmt.Sprintf(nmKeyErrorMissingDependsPackage, getSettingVkVpnMissingPlugin(s.data))
			} else {
				logger.Errorf("missing section, errors[%s]", nm.NM_SETTING_VS_VPN)
			}
		}
	}

	dbus.NotifyChange(s, "Errors")
}
