/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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
	. "pkg.deepin.io/lib/gettext"
)

func isWirelessSecurityNeedShowWepKey(data connectionData) bool {
	return isSettingRequireSecret(getSettingWirelessSecurityWepKeyFlags(data))
}

func isWirelessSecurityNeedShowPsk(data connectionData) bool {
	return isSettingRequireSecret(getSettingWirelessSecurityPskFlags(data))
}

// Get available keys
func getSettingWirelessSecurityAvailableKeys(data connectionData) (keys []string) {
	vkKeyMgmt := getSettingVkWirelessSecurityKeyMgmt(data)
	switch vkKeyMgmt {
	default:
		logger.Error("invalid value", vkKeyMgmt)
	case "none":
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_WIRELESS_SECURITY_SETTING_NAME, nm.NM_SETTING_WIRELESS_SECURITY_KEY_MGMT)
	case "wep":
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_WIRELESS_SECURITY_SETTING_NAME, nm.NM_SETTING_WIRELESS_SECURITY_KEY_MGMT)
		if isWirelessSecurityNeedShowWepKey(data) {
			keys = appendAvailableKeys(data, keys, nm.NM_SETTING_WIRELESS_SECURITY_SETTING_NAME, nm.NM_SETTING_WIRELESS_SECURITY_WEP_KEY0)
		}
	case "wpa-psk":
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_WIRELESS_SECURITY_SETTING_NAME, nm.NM_SETTING_WIRELESS_SECURITY_KEY_MGMT)
		if isWirelessSecurityNeedShowPsk(data) {
			keys = appendAvailableKeys(data, keys, nm.NM_SETTING_WIRELESS_SECURITY_SETTING_NAME, nm.NM_SETTING_WIRELESS_SECURITY_PSK)
		}
	case "wpa-eap":
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_WIRELESS_SECURITY_SETTING_NAME, nm.NM_SETTING_WIRELESS_SECURITY_KEY_MGMT)
	}
	return
}

// Get available values
func getSettingWirelessSecurityAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case nm.NM_SETTING_WIRELESS_SECURITY_KEY_MGMT:
		values = []kvalue{
			kvalue{"none", Tr("WEP")},
			kvalue{"ieee8021x", Tr("Dynamic WEP")},
			kvalue{"wpa-none", Tr("Ad-Hoc WPA-PSK")},
			kvalue{"wpa-psk", Tr("WPA-PSK Infrastructure")},
			kvalue{"wpa-eap", Tr("WPA Enterprise")},
		}
	case nm.NM_SETTING_WIRELESS_SECURITY_GROUP:
		values = []kvalue{
			kvalue{"wep40", Tr("WEP40")},
			kvalue{"wep104", Tr("WEP104")},
			kvalue{"tkip", Tr("TKIP")},
			kvalue{"ccmp", Tr("CCMP")},
		}
	case nm.NM_SETTING_WIRELESS_SECURITY_AUTH_ALG:
		values = []kvalue{
			kvalue{"open", Tr("Open")},
			kvalue{"shared", Tr("Shared")},
			kvalue{"leap", Tr("LEAP")},
		}
	}
	return
}

// Check whether the values are correct
func checkSettingWirelessSecurityValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)

	// check key-mgmt
	ensureSettingWirelessSecurityKeyMgmtNoEmpty(data, errs)
	switch getSettingWirelessSecurityKeyMgmt(data) {
	default:
		rememberError(errs, nm.NM_SETTING_WIRELESS_SECURITY_SETTING_NAME, nm.NM_SETTING_WIRELESS_SECURITY_KEY_MGMT, nmKeyErrorInvalidValue)
		return
	case "none": // wep
		ensureSettingWirelessSecurityWepKeyTypeNoEmpty(data, errs)
		if isWirelessSecurityNeedShowWepKey(data) {
			ensureSettingWirelessSecurityWepKey0NoEmpty(data, errs)
		}
	case "ieee8021x": // dynamic wep
	case "wpa-none": // wpa-psk ad-hoc
	case "wpa-psk": // wpa-psk infrastructure
		if isWirelessSecurityNeedShowPsk(data) {
			ensureSettingWirelessSecurityPskNoEmpty(data, errs)
		}
	case "wpa-eap": // wpa enterprise
		ensureSectionSetting8021xExists(data, errs, nm.NM_SETTING_VK_WIRELESS_SECURITY_KEY_MGMT)
	}

	// check wep-key-type
	checkSettingWirelessSecurityWepKeyType(data, errs)

	// check wep-key0
	checkSettingWirelessSecurityWepKey0(data, errs)

	// check psk
	checkSettingWirelessSecurityPsk(data, errs)

	return
}

func checkSettingWirelessSecurityWepKeyType(data connectionData, errs sectionErrors) {
	if !isSettingWirelessSecurityWepKeyTypeExists(data) {
		return
	}
	wepKeyType := getSettingWirelessSecurityWepKeyType(data)
	if wepKeyType != nm.NM_WEP_KEY_TYPE_KEY && wepKeyType != nm.NM_WEP_KEY_TYPE_PASSPHRASE {
		rememberError(errs, nm.NM_SETTING_WIRELESS_SECURITY_SETTING_NAME, nm.NM_SETTING_WIRELESS_SECURITY_WEP_KEY_TYPE, nmKeyErrorInvalidValue)
	}
}
func checkSettingWirelessSecurityWepKey0(data connectionData, errs sectionErrors) {
	if !isSettingWirelessSecurityWepKey0Exists(data) {
		return
	}
	wepKey0 := getSettingWirelessSecurityWepKey0(data)
	wepKeyType := getSettingWirelessSecurityWepKeyType(data)
	if wepKeyType == nm.NM_WEP_KEY_TYPE_KEY {
		if !isPasswordValid(passTypeWifiWepKey, wepKey0) {
			rememberError(errs, nm.NM_SETTING_WIRELESS_SECURITY_SETTING_NAME, nm.NM_SETTING_WIRELESS_SECURITY_WEP_KEY0, nmKeyErrorInvalidValue)
		}
	} else if wepKeyType == nm.NM_WEP_KEY_TYPE_PASSPHRASE {
		if !isPasswordValid(passTypeWifiWepPassphrase, wepKey0) {
			rememberError(errs, nm.NM_SETTING_WIRELESS_SECURITY_SETTING_NAME, nm.NM_SETTING_WIRELESS_SECURITY_WEP_KEY0, nmKeyErrorInvalidValue)
		}
	}
}
func checkSettingWirelessSecurityPsk(data connectionData, errs sectionErrors) {
	if !isSettingWirelessSecurityPskExists(data) {
		return
	}
	psk := getSettingWirelessSecurityPsk(data)
	if !isPasswordValid(passTypeWifiWpaPsk, psk) {
		rememberError(errs, nm.NM_SETTING_WIRELESS_SECURITY_SETTING_NAME, nm.NM_SETTING_WIRELESS_SECURITY_PSK, nmKeyErrorInvalidValue)
	}
}

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
