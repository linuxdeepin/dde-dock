package main

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

// TODO Check whether the values are correct
func checkSettingWirelessSecurityValues(data _ConnectionData) (errs map[string]string) {
	errs = make(map[string]string)
	return
}

// Get and set key's value generally
func generalGetSettingWirelessSecurityKeyJSON(data _ConnectionData, key string) (value string) {
	switch key {
	default:
		LOGGER.Error("generalGetSettingWirelessSecurityKey: invalide key", key)
	case NM_SETTING_WIRELESS_SECURITY_KEY_MGMT:
		value = getSettingWirelessSecurityKeyMgmtJSON(data)
	case NM_SETTING_WIRELESS_SECURITY_WEP_TX_KEYIDX:
		value = getSettingWirelessSecurityWepTxKeyidxJSON(data)
	case NM_SETTING_WIRELESS_SECURITY_AUTH_ALG:
		value = getSettingWirelessSecurityAuthAlgJSON(data)
	case NM_SETTING_WIRELESS_SECURITY_PROTO:
		value = getSettingWirelessSecurityProtoJSON(data)
	case NM_SETTING_WIRELESS_SECURITY_PAIRWISE:
		value = getSettingWirelessSecurityPairwiseJSON(data)
	case NM_SETTING_WIRELESS_SECURITY_GROUP:
		value = getSettingWirelessSecurityGroupJSON(data)
	case NM_SETTING_WIRELESS_SECURITY_LEAP_USERNAME:
		value = getSettingWirelessSecurityLeapUsernameJSON(data)
	case NM_SETTING_WIRELESS_SECURITY_WEP_KEY0:
		value = getSettingWirelessSecurityWepKey0JSON(data)
	case NM_SETTING_WIRELESS_SECURITY_WEP_KEY1:
		value = getSettingWirelessSecurityWepKey1JSON(data)
	case NM_SETTING_WIRELESS_SECURITY_WEP_KEY2:
		value = getSettingWirelessSecurityWepKey2JSON(data)
	case NM_SETTING_WIRELESS_SECURITY_WEP_KEY3:
		value = getSettingWirelessSecurityWepKey3JSON(data)
	case NM_SETTING_WIRELESS_SECURITY_WEP_KEY_FLAGS:
		value = getSettingWirelessSecurityWepKeyFlagsJSON(data)
	case NM_SETTING_WIRELESS_SECURITY_WEP_KEY_TYPE:
		value = getSettingWirelessSecurityWepKeyTypeJSON(data)
	case NM_SETTING_WIRELESS_SECURITY_PSK:
		value = getSettingWirelessSecurityPskJSON(data)
	case NM_SETTING_WIRELESS_SECURITY_PSK_FLAGS:
		value = getSettingWirelessSecurityPskFlagsJSON(data)
	case NM_SETTING_WIRELESS_SECURITY_LEAP_PASSWORD:
		value = getSettingWirelessSecurityLeapPasswordJSON(data)
	case NM_SETTING_WIRELESS_SECURITY_LEAP_PASSWORD_FLAGS:
		value = getSettingWirelessSecurityLeapPasswordFlagsJSON(data)
	}
	return
}

// TODO use logic setter
func generalSetSettingWirelessSecurityKeyJSON(data _ConnectionData, key, value string) {
	switch key {
	default:
		LOGGER.Error("generalSetSettingWirelessSecurityKey: invalide key", key)
	case NM_SETTING_WIRELESS_SECURITY_KEY_MGMT:
		setSettingWirelessSecurityKeyMgmtJSON(data, value)
	case NM_SETTING_WIRELESS_SECURITY_WEP_TX_KEYIDX:
		setSettingWirelessSecurityWepTxKeyidxJSON(data, value)
	case NM_SETTING_WIRELESS_SECURITY_AUTH_ALG:
		setSettingWirelessSecurityAuthAlgJSON(data, value)
	case NM_SETTING_WIRELESS_SECURITY_PROTO:
		setSettingWirelessSecurityProtoJSON(data, value)
	case NM_SETTING_WIRELESS_SECURITY_PAIRWISE:
		setSettingWirelessSecurityPairwiseJSON(data, value)
	case NM_SETTING_WIRELESS_SECURITY_GROUP:
		setSettingWirelessSecurityGroupJSON(data, value)
	case NM_SETTING_WIRELESS_SECURITY_LEAP_USERNAME:
		setSettingWirelessSecurityLeapUsernameJSON(data, value)
	case NM_SETTING_WIRELESS_SECURITY_WEP_KEY0:
		setSettingWirelessSecurityWepKey0JSON(data, value)
	case NM_SETTING_WIRELESS_SECURITY_WEP_KEY1:
		setSettingWirelessSecurityWepKey1JSON(data, value)
	case NM_SETTING_WIRELESS_SECURITY_WEP_KEY2:
		setSettingWirelessSecurityWepKey2JSON(data, value)
	case NM_SETTING_WIRELESS_SECURITY_WEP_KEY3:
		setSettingWirelessSecurityWepKey3JSON(data, value)
	case NM_SETTING_WIRELESS_SECURITY_WEP_KEY_FLAGS:
		setSettingWirelessSecurityWepKeyFlagsJSON(data, value)
	case NM_SETTING_WIRELESS_SECURITY_WEP_KEY_TYPE:
		setSettingWirelessSecurityWepKeyTypeJSON(data, value)
	case NM_SETTING_WIRELESS_SECURITY_PSK:
		setSettingWirelessSecurityPskJSON(data, value)
	case NM_SETTING_WIRELESS_SECURITY_PSK_FLAGS:
		setSettingWirelessSecurityPskFlagsJSON(data, value)
	case NM_SETTING_WIRELESS_SECURITY_LEAP_PASSWORD:
		setSettingWirelessSecurityLeapPasswordJSON(data, value)
	case NM_SETTING_WIRELESS_SECURITY_LEAP_PASSWORD_FLAGS:
		setSettingWirelessSecurityLeapPasswordFlagsJSON(data, value)
	}
	return
}
