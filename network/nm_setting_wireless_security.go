/**
 * Copyright (c) 2014 Deepin, Inc.
 *               2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package network

import (
	. "pkg.deepin.io/lib/gettext"
	"fmt"
)

// https://developer.gnome.org/libnm-util/0.9/NMSettingWirelessSecurity.html
// https://developer.gnome.org/NetworkManager/unstable/ref-settings.html

const NM_SETTING_WIRELESS_SECURITY_SETTING_NAME = "802-11-wireless-security"

const (
	// Key management used for the connection. One of 'none' (WEP),
	// 'ieee8021x' (Dynamic WEP), 'wpa-none' (Ad-Hoc WPA-PSK),
	// 'wpa-psk' (infrastructure WPA-PSK), or 'wpa-eap'
	// (WPA-Enterprise). This property must be set for any WiFi
	// connection that uses security.
	// Default value: NULL
	NM_SETTING_WIRELESS_SECURITY_KEY_MGMT = "key-mgmt"

	// When static WEP is used (ie, key-mgmt = 'none') and a
	// non-default WEP key index is used by the AP, put that WEP key
	// index here. Valid values are 0 (default key) through 3. Note
	// that some consumer access points (like the Linksys WRT54G)
	// number the keys 1 - 4.
	// Allowed values: <= 3
	// Default value: 0
	NM_SETTING_WIRELESS_SECURITY_WEP_TX_KEYIDX = "wep-tx-keyidx"

	// When WEP is used (ie, key-mgmt = 'none' or 'ieee8021x')
	// indicate the 802.11 authentication algorithm required by the AP
	// here. One of 'open' for Open System, 'shared' for Shared Key,
	// or 'leap' for Cisco LEAP. When using Cisco LEAP (ie, key-mgmt =
	// 'ieee8021x' and auth-alg = 'leap') the 'leap-username' and
	// 'leap-password' properties must be specified.
	// Default value: NULL
	NM_SETTING_WIRELESS_SECURITY_AUTH_ALG = "auth-alg"

	// List of strings specifying the allowed WPA protocol versions to
	// use. Each element may be one 'wpa' (allow WPA) or 'rsn' (allow
	// WPA2/RSN). If not specified, both WPA and RSN connections are
	// allowed.
	NM_SETTING_WIRELESS_SECURITY_PROTO = "proto"

	// A list of pairwise encryption algorithms which prevents
	// connections to Wi-Fi networks that do not utilize one of the
	// algorithms in the list. For maximum compatibility leave this
	// property empty. Each list element may be one of 'wep40',
	// 'wep104', 'tkip', or 'ccmp'.
	NM_SETTING_WIRELESS_SECURITY_PAIRWISE = "pairwise"

	// A list of group/broadcast encryption algorithms which prevents
	// connections to Wi-Fi networks that do not utilize one of the
	// algorithms in the list. For maximum compatibility leave this
	// property empty. Each list element may be one of 'wep40',
	// 'wep104', 'tkip', or 'ccmp'.
	NM_SETTING_WIRELESS_SECURITY_GROUP = "group"

	// The login username for legacy LEAP connections (ie, key-mgmt =
	// 'ieee8021x' and auth-alg = 'leap').
	// Default value: NULL
	NM_SETTING_WIRELESS_SECURITY_LEAP_USERNAME = "leap-username"

	// Flags indicating how to handle "leap-password".
	// Allowed values: <= 7
	// Default value: 0
	NM_SETTING_WIRELESS_SECURITY_LEAP_PASSWORD_FLAGS = "leap-password-flags"

	// The login password for legacy LEAP connections (ie, key-mgmt =
	// 'ieee8021x' and auth-alg = 'leap').
	// Default value: NULL
	NM_SETTING_WIRELESS_SECURITY_LEAP_PASSWORD = "leap-password"

	// Index 0 WEP key. This is the WEP key used in most networks. See
	// the 'wep-key-type' property for a description of how this key
	// is interpreted.
	// Default value: NULL
	NM_SETTING_WIRELESS_SECURITY_WEP_KEY0 = "wep-key0"

	// Index 1 WEP key. This WEP index is not used by most
	// networks. See the 'wep-key-type' property for a description of
	// how this key is interpreted.
	// Default value: NULL
	NM_SETTING_WIRELESS_SECURITY_WEP_KEY1 = "wep-key1"

	// Index 2 WEP key. This WEP index is not used by most
	// networks. See the 'wep-key-type' property for a description of
	// how this key is interpreted.
	// Default value: NULL
	NM_SETTING_WIRELESS_SECURITY_WEP_KEY2 = "wep-key2"

	// Index 3 WEP key. This WEP index is not used by most
	// networks. See the 'wep-key-type' property for a description of
	// how this key is interpreted.
	// Default value: NULL
	NM_SETTING_WIRELESS_SECURITY_WEP_KEY3 = "wep-key3"

	// Flags indicating how to handle NMSettingWirelessSecurity WEP keys.
	// Allowed values: <= 7
	// Default value: 0
	NM_SETTING_WIRELESS_SECURITY_WEP_KEY_FLAGS = "wep-key-flags"

	// Controls the interpretation of WEP keys. Allowed values are
	// those given by NMWepKeyType. If set to 1 and the keys are
	// hexadecimal, they must be either 10 or 26 characters in
	// length. If set to 1 and the keys are ASCII keys, they must be
	// either 5 or 13 characters in length. If set to 2, the
	// passphrase is hashed using the de-facto MD5 method to derive
	// the actual WEP key.
	// Allowed values: <= 2
	// Default value: 0
	NM_SETTING_WIRELESS_SECURITY_WEP_KEY_TYPE = "wep-key-type"

	// Pre-Shared-Key for WPA networks. If the key is 64-characters
	// long, it must contain only hexadecimal characters and is
	// interpreted as a hexadecimal WPA key. Otherwise, the key must
	// be between 8 and 63 ASCII characters (as specified in the
	// 802.11i standard) and is interpreted as a WPA passphrase, and
	// is hashed to derive the actual WPA-PSK used when connecting to
	// the WiFi network.
	// Default value: NULL
	NM_SETTING_WIRELESS_SECURITY_PSK = "psk"

	// Flags indicating how to handle "psk"
	// Allowed values: <= 7
	// Default value: 0
	NM_SETTING_WIRELESS_SECURITY_PSK_FLAGS = "psk-flags"
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
		removeSettingWirelessSec(data)
		removeSettingSection(data, sectionWirelessSecurity)
		removeSettingSection(data, section8021x)
	case "wep":
		setSettingWirelessSec(data, sectionWirelessSecurity)
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
		setSettingWirelessSec(data, sectionWirelessSecurity)
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
		setSettingWirelessSec(data, sectionWirelessSecurity)
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
