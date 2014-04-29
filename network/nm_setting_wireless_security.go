package main

import (
	"dlib"
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

// Get available keys
func getSettingWirelessSecurityAvailableKeys(data connectionData) (keys []string) {
	vkKeyMgmt := getSettingVkWirelessSecurityKeyMgmt(data)
	switch vkKeyMgmt {
	default:
		logger.Error("invalid value", vkKeyMgmt)
	case "none":
		keys = appendAvailableKeys(keys, fieldWirelessSecurity, NM_SETTING_WIRELESS_SECURITY_KEY_MGMT)
	case "wep":
		keys = appendAvailableKeys(keys, fieldWirelessSecurity, NM_SETTING_WIRELESS_SECURITY_KEY_MGMT)
		keys = appendAvailableKeys(keys, fieldWirelessSecurity, NM_SETTING_WIRELESS_SECURITY_WEP_KEY0)
	case "wpa-psk":
		keys = appendAvailableKeys(keys, fieldWirelessSecurity, NM_SETTING_WIRELESS_SECURITY_KEY_MGMT)
		keys = appendAvailableKeys(keys, fieldWirelessSecurity, NM_SETTING_WIRELESS_SECURITY_PSK)
	case "wpa-eap":
		keys = appendAvailableKeys(keys, fieldWirelessSecurity, NM_SETTING_WIRELESS_SECURITY_KEY_MGMT)
	}
	return
}

// Get available values
func getSettingWirelessSecurityAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case NM_SETTING_WIRELESS_SECURITY_KEY_MGMT:
		values = []kvalue{
			kvalue{"none", dlib.Tr("wep")},
			kvalue{"ieee8021x", dlib.Tr("dynamic wep")},
			kvalue{"wpa-none", dlib.Tr("wpa-psk ad-hoc")},
			kvalue{"wpa-psk", dlib.Tr("wpa-psk infrastructure")},
			kvalue{"wpa-eap", dlib.Tr("wpa enterprise")},
		}
	case NM_SETTING_WIRELESS_SECURITY_GROUP:
		values = []kvalue{
			kvalue{"wep40", dlib.Tr("wep40")},
			kvalue{"wep104", dlib.Tr("wep104")},
			kvalue{"tkip", dlib.Tr("tkip")},
			kvalue{"ccmp", dlib.Tr("ccmp")},
		}
	case NM_SETTING_WIRELESS_SECURITY_AUTH_ALG:
		values = []kvalue{
			kvalue{"open", dlib.Tr("open")},
			kvalue{"shared", dlib.Tr("shared")},
			kvalue{"leap", dlib.Tr("leap")},
		}
	}
	return
}

// Check whether the values are correct
func checkSettingWirelessSecurityValues(data connectionData) (errs fieldErrors) {
	errs = make(map[string]string)

	// check key-mgmt
	ensureSettingWirelessSecurityKeyMgmtNoEmpty(data, errs)
	switch getSettingWirelessSecurityKeyMgmt(data) {
	default:
		rememberError(errs, fieldWirelessSecurity, NM_SETTING_WIRELESS_SECURITY_KEY_MGMT, NM_KEY_ERROR_INVALID_VALUE)
		return
	case "none": // wep
		ensureSettingWirelessSecurityWepKeyTypeNoEmpty(data, errs)
		ensureSettingWirelessSecurityWepKey0NoEmpty(data, errs)
	case "ieee8021x": // dynamic wep
	case "wpa-none": // wpa-psk ad-hoc
	case "wpa-psk": // wpa-psk infrastructure
		ensureSettingWirelessSecurityPskNoEmpty(data, errs)
	case "wpa-eap": // wpa enterprise
		ensureFieldSetting8021xExists(data, errs, NM_SETTING_VK_WIRELESS_SECURITY_KEY_MGMT)
	}

	// check wep-key-type
	checkSettingWirelessSecurityWepKeyType(data, errs)

	// check wep-key0
	checkSettingWirelessSecurityWepKey0(data, errs)

	// check psk
	checkSettingWirelessSecurityPsk(data, errs)

	return
}

func checkSettingWirelessSecurityWepKeyType(data connectionData, errs fieldErrors) {
	if !isSettingWirelessSecurityWepKeyTypeExists(data) {
		return
	}
	wepKeyType := getSettingWirelessSecurityWepKeyType(data)
	if wepKeyType != 1 && wepKeyType != 2 {
		rememberError(errs, fieldWirelessSecurity, NM_SETTING_WIRELESS_SECURITY_WEP_KEY_TYPE, NM_KEY_ERROR_INVALID_VALUE)
	}
}
func checkSettingWirelessSecurityWepKey0(data connectionData, errs fieldErrors) {
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
			rememberError(errs, fieldWirelessSecurity, NM_SETTING_WIRELESS_SECURITY_WEP_KEY0, NM_KEY_ERROR_INVALID_VALUE)
		}
	} else if wepKeyType == 2 {
		// If set to 2, the passphrase is hashed using the de-facto
		// MD5 method to derive the actual WEP key.
		if len(wepKey0) == 0 {
			rememberError(errs, fieldWirelessSecurity, NM_SETTING_WIRELESS_SECURITY_WEP_KEY0, NM_KEY_ERROR_INVALID_VALUE)
		}
	}
}
func checkSettingWirelessSecurityPsk(data connectionData, errs fieldErrors) {
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
		rememberError(errs, fieldWirelessSecurity, NM_SETTING_WIRELESS_SECURITY_PSK, NM_KEY_ERROR_INVALID_VALUE)
	}
}

// Virtual key getter and setter
func getSettingVkWirelessSecurityKeyMgmt(data connectionData) (value string) {
	if !isSettingFieldExists(data, fieldWirelessSecurity) {
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
		removeSettingField(data, fieldWirelessSecurity)
		removeSettingField(data, field8021x)
	case "wep":
		addSettingField(data, fieldWirelessSecurity)
		removeSettingField(data, field8021x)
		setSettingWirelessSecurityKeyMgmt(data, "none")
		setSettingWirelessSecurityAuthAlg(data, "open")
		setSettingWirelessSecurityWepKeyType(data, 1)
	case "wpa-psk":
		addSettingField(data, fieldWirelessSecurity)
		removeSettingField(data, field8021x)
		setSettingWirelessSecurityKeyMgmt(data, "wpa-psk")
	case "wpa-eap":
		addSettingField(data, fieldWirelessSecurity)
		addSettingField(data, field8021x)
		setSettingWirelessSecurityKeyMgmt(data, "wpa-eap")
		err = logicSetSetting8021xEap(data, []string{"tls"})
	}
	return
}
