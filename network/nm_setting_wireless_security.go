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
	"errors"
	"fmt"

	"pkg.deepin.io/dde/daemon/network/nm"
)

// Virtual key getter and setter
func getSettingVkWirelessSecurityKeyMgmt(data connectionData) (value string) {
	if !isSettingExists(data, nm.NM_SETTING_WIRELESS_SECURITY_SETTING_NAME) {
		value = "none"
		return
	}
	keyMgmt := getSettingWirelessSecurityKeyMgmt(data)
	switch keyMgmt {
	case "none":
		value = "wep"
	case "wpa-psk":
		value = "wpa-psk"
	case "wpa-eap":
		value = "wpa-eap"
	}
	return
}

func getApSecTypeFromConnData(data connectionData) (apSecType, error) {
	if !isSettingExists(data, nm.NM_SETTING_WIRELESS_SECURITY_SETTING_NAME) {
		return apSecNone, nil
	}
	keyMgmt := getSettingWirelessSecurityKeyMgmt(data)
	switch keyMgmt {
	case "none":
		if "open" == getSettingWirelessSecurityAuthAlg(data) {
			return apSecWep, nil
		}
	case "wpa-psk":
		return apSecPsk, nil

	case "wpa-eap":
		return apSecEap, nil
	}

	return apSecNone, errors.New("unknown apSecType")
}

func logicSetSettingVkWirelessSecurityKeyMgmt(data connectionData, value string) (err error) {
	switch value {
	default:
		logger.Error("invalid value", value)
		err = fmt.Errorf(nmKeyErrorInvalidValue)
	case "none":
		removeSetting(data, nm.NM_SETTING_WIRELESS_SECURITY_SETTING_NAME)
		removeSetting(data, nm.NM_SETTING_802_1X_SETTING_NAME)
	case "wep":
		addSetting(data, nm.NM_SETTING_WIRELESS_SECURITY_SETTING_NAME)
		removeSetting(data, nm.NM_SETTING_802_1X_SETTING_NAME)

		removeSettingKeyBut(data, nm.NM_SETTING_WIRELESS_SECURITY_SETTING_NAME,
			nm.NM_SETTING_WIRELESS_SECURITY_KEY_MGMT,
			nm.NM_SETTING_WIRELESS_SECURITY_AUTH_ALG,
			nm.NM_SETTING_WIRELESS_SECURITY_WEP_KEY0,
			nm.NM_SETTING_WIRELESS_SECURITY_WEP_KEY_FLAGS,
			nm.NM_SETTING_WIRELESS_SECURITY_WEP_KEY_TYPE,
		)
		setSettingWirelessSecurityKeyMgmt(data, "none")
		setSettingWirelessSecurityAuthAlg(data, "open")
		setSettingWirelessSecurityWepKeyFlags(data, nm.NM_SETTING_SECRET_FLAG_NONE)
		setSettingWirelessSecurityWepKeyType(data, nm.NM_WEP_KEY_TYPE_KEY)
	case "wpa-psk":
		addSetting(data, nm.NM_SETTING_WIRELESS_SECURITY_SETTING_NAME)
		removeSetting(data, nm.NM_SETTING_802_1X_SETTING_NAME)

		removeSettingKeyBut(data, nm.NM_SETTING_WIRELESS_SECURITY_SETTING_NAME,
			nm.NM_SETTING_WIRELESS_SECURITY_KEY_MGMT,
			nm.NM_SETTING_WIRELESS_SECURITY_PSK,
			nm.NM_SETTING_WIRELESS_SECURITY_PSK_FLAGS,
		)
		setSettingWirelessSecurityKeyMgmt(data, "wpa-psk")
		setSettingWirelessSecurityPskFlags(data, nm.NM_SETTING_SECRET_FLAG_NONE)
	case "wpa-eap":
		addSetting(data, nm.NM_SETTING_WIRELESS_SECURITY_SETTING_NAME)
		addSetting(data, nm.NM_SETTING_802_1X_SETTING_NAME)

		removeSettingKeyBut(data, nm.NM_SETTING_WIRELESS_SECURITY_SETTING_NAME,
			nm.NM_SETTING_WIRELESS_SECURITY_KEY_MGMT,
		)
		setSettingWirelessSecurityKeyMgmt(data, "wpa-eap")
		err = logicSetSetting8021xEap(data, []string{"tls"})
	}
	return
}
