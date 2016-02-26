/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package network

import (
	"fmt"
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
		keys = appendAvailableKeys(data, keys, sectionWirelessSecurity, NM_SETTING_WIRELESS_SECURITY_KEY_MGMT)
	case "wep":
		keys = appendAvailableKeys(data, keys, sectionWirelessSecurity, NM_SETTING_WIRELESS_SECURITY_KEY_MGMT)
		if isWirelessSecurityNeedShowWepKey(data) {
			keys = appendAvailableKeys(data, keys, sectionWirelessSecurity, NM_SETTING_WIRELESS_SECURITY_WEP_KEY0)
		}
	case "wpa-psk":
		keys = appendAvailableKeys(data, keys, sectionWirelessSecurity, NM_SETTING_WIRELESS_SECURITY_KEY_MGMT)
		if isWirelessSecurityNeedShowPsk(data) {
			keys = appendAvailableKeys(data, keys, sectionWirelessSecurity, NM_SETTING_WIRELESS_SECURITY_PSK)
		}
	case "wpa-eap":
		keys = appendAvailableKeys(data, keys, sectionWirelessSecurity, NM_SETTING_WIRELESS_SECURITY_KEY_MGMT)
	}
	return
}

// Get available values
func getSettingWirelessSecurityAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case NM_SETTING_WIRELESS_SECURITY_KEY_MGMT:
		values = []kvalue{
			kvalue{"none", Tr("WEP")},
			kvalue{"ieee8021x", Tr("Dynamic WEP")},
			kvalue{"wpa-none", Tr("Ad-Hoc WPA-PSK")},
			kvalue{"wpa-psk", Tr("WPA-PSK Infrastructure")},
			kvalue{"wpa-eap", Tr("WPA Enterprise")},
		}
	case NM_SETTING_WIRELESS_SECURITY_GROUP:
		values = []kvalue{
			kvalue{"wep40", Tr("WEP40")},
			kvalue{"wep104", Tr("WEP104")},
			kvalue{"tkip", Tr("TKIP")},
			kvalue{"ccmp", Tr("CCMP")},
		}
	case NM_SETTING_WIRELESS_SECURITY_AUTH_ALG:
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
		rememberError(errs, sectionWirelessSecurity, NM_SETTING_WIRELESS_SECURITY_KEY_MGMT, NM_KEY_ERROR_INVALID_VALUE)
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
		ensureSectionSetting8021xExists(data, errs, NM_SETTING_VK_WIRELESS_SECURITY_KEY_MGMT)
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
	if wepKeyType != 1 && wepKeyType != 2 {
		rememberError(errs, sectionWirelessSecurity, NM_SETTING_WIRELESS_SECURITY_WEP_KEY_TYPE, NM_KEY_ERROR_INVALID_VALUE)
	}
}
func checkSettingWirelessSecurityWepKey0(data connectionData, errs sectionErrors) {
	if !isSettingWirelessSecurityWepKey0Exists(data) {
		return
	}
	wepKey0 := getSettingWirelessSecurityWepKey0(data)
	wepKeyType := getSettingWirelessSecurityWepKeyType(data)
	if wepKeyType == 1 {
		// If set to 1 and the keys are hexadecimal, they must be
		// either 10 or 26 characters in length. If set to 1 and the
		// keys are ASCII keys, they must be either 5 or 13 characters
		// in length.
		if len(wepKey0) != 10 && len(wepKey0) != 26 && len(wepKey0) != 5 && len(wepKey0) != 13 {
			rememberError(errs, sectionWirelessSecurity, NM_SETTING_WIRELESS_SECURITY_WEP_KEY0, NM_KEY_ERROR_INVALID_VALUE)
		}
	} else if wepKeyType == 2 {
		// If set to 2, the passphrase is hashed using the de-facto
		// MD5 method to derive the actual WEP key.
		if len(wepKey0) == 0 {
			rememberError(errs, sectionWirelessSecurity, NM_SETTING_WIRELESS_SECURITY_WEP_KEY0, NM_KEY_ERROR_INVALID_VALUE)
		}
	}
}
func checkSettingWirelessSecurityPsk(data connectionData, errs sectionErrors) {
	if !isSettingWirelessSecurityPskExists(data) {
		return
	}
	psk := getSettingWirelessSecurityPsk(data)
	// If the key is 64-characters long, it must contain only
	// hexadecimal characters and is interpreted as a hexadecimal WPA
	// key. Otherwise, the key must be between 8 and 63 ASCII
	// characters (as specified in the 802.11i standard) and is
	// interpreted as a WPA passphrase
	if len(psk) < 8 || len(psk) > 64 {
		// TODO
		rememberError(errs, sectionWirelessSecurity, NM_SETTING_WIRELESS_SECURITY_PSK, NM_KEY_ERROR_INVALID_VALUE)
	}
}

// Virtual key getter and setter
func getSettingVkWirelessSecurityKeyMgmt(data connectionData) (value string) {
	if !isSettingSectionExists(data, sectionWirelessSecurity) {
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
		err = fmt.Errorf(NM_KEY_ERROR_INVALID_VALUE)
	case "none":
		// removeSettingWirelessSec(data) // TODO:
		removeSettingSection(data, sectionWirelessSecurity)
		removeSettingSection(data, section8021x)
	case "wep":
		// setSettingWirelessSec(data, sectionWirelessSecurity) // TODO:
		addSettingSection(data, sectionWirelessSecurity)
		removeSettingSection(data, section8021x)

		removeSettingKeyBut(data, sectionWirelessSecurity,
			NM_SETTING_WIRELESS_SECURITY_KEY_MGMT,
			NM_SETTING_WIRELESS_SECURITY_AUTH_ALG,
			NM_SETTING_WIRELESS_SECURITY_WEP_KEY0,
			NM_SETTING_WIRELESS_SECURITY_WEP_KEY_FLAGS,
			NM_SETTING_WIRELESS_SECURITY_WEP_KEY_TYPE,
		)
		setSettingWirelessSecurityKeyMgmt(data, "none")
		setSettingWirelessSecurityAuthAlg(data, "open")
		setSettingWirelessSecurityWepKeyFlags(data, NM_SETTING_SECRET_FLAG_NONE)
		setSettingWirelessSecurityWepKeyType(data, 1)
	case "wpa-psk":
		// setSettingWirelessSec(data, sectionWirelessSecurity) // TODO:
		addSettingSection(data, sectionWirelessSecurity)
		removeSettingSection(data, section8021x)

		removeSettingKeyBut(data, sectionWirelessSecurity,
			NM_SETTING_WIRELESS_SECURITY_KEY_MGMT,
			NM_SETTING_WIRELESS_SECURITY_PSK,
			NM_SETTING_WIRELESS_SECURITY_PSK_FLAGS,
		)
		setSettingWirelessSecurityKeyMgmt(data, "wpa-psk")
		setSettingWirelessSecurityPskFlags(data, NM_SETTING_SECRET_FLAG_NONE)
	case "wpa-eap":
		// setSettingWirelessSec(data, sectionWirelessSecurity) // TODO:
		addSettingSection(data, sectionWirelessSecurity)
		addSettingSection(data, section8021x)

		removeSettingKeyBut(data, sectionWirelessSecurity,
			NM_SETTING_WIRELESS_SECURITY_KEY_MGMT,
		)
		setSettingWirelessSecurityKeyMgmt(data, "wpa-eap")
		err = logicSetSetting8021xEap(data, []string{"tls"})
	}
	return
}
