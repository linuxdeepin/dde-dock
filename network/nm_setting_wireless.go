package main

import (
	"fmt"
)

const NM_SETTING_WIRELESS_SETTING_NAME = "802-11-wireless"

const (
	// SSID of the WiFi network. Must be specified.
	NM_SETTING_WIRELESS_SSID = "ssid"

	// WiFi network mode; one of 'infrastructure', 'adhoc' or 'ap'. If
	// blank, infrastructure is assumed.
	NM_SETTING_WIRELESS_MODE = "mode"

	// 802.11 frequency band of the network. One of 'a' for 5GHz
	// 802.11a or 'bg' for 2.4GHz 802.11. This will lock associations
	// to the WiFi network to the specific band, i.e. if 'a' is
	// specified, the device will not associate with the same network
	// in the 2.4GHz band even if the network's settings are
	// compatible. This setting depends on specific driver capability
	// and may not work with all drivers.
	NM_SETTING_WIRELESS_BAND = "band"

	// Wireless channel to use for the WiFi connection. The device
	// will only join (or create for Ad-Hoc networks) a WiFi network
	// on the specified channel. Because channel numbers overlap
	// between bands, this property also requires the 'band' property
	// to be set.
	NM_SETTING_WIRELESS_CHANNEL = "channel"

	// If specified, directs the device to only associate with the
	// given access point. This capability is highly driver dependent
	// and not supported by all devices. Note: this property does not
	// control the BSSID used when creating an Ad-Hoc network and is
	// unlikely to in the future.
	NM_SETTING_WIRELESS_BSSID = "bssid"

	// If non-zero, directs the device to only use the specified
	// bitrate for communication with the access point. Units are in
	// Kb/s, ie 5500 = 5.5 Mbit/s. This property is highly driver
	// dependent and not all devices support setting a static bitrate.
	NM_SETTING_WIRELESS_RATE = "rate"

	// If non-zero, directs the device to use the specified transmit
	// power. Units are dBm. This property is highly driver dependent
	// and not all devices support setting a static transmit power.
	NM_SETTING_WIRELESS_TX_POWER = "tx-power"

	// If specified, this connection will only apply to the WiFi
	// device whose permanent MAC address matches. This property does
	// not change the MAC address of the device (i.e. MAC spoofing).
	NM_SETTING_WIRELESS_MAC_ADDRESS = "mac-address"

	// If specified, request that the WiFi device use this MAC address
	// instead of its permanent MAC address. This is known as MAC
	// cloning or spoofing.
	NM_SETTING_WIRELESS_CLONED_MAC_ADDRESS = "cloned-mac-address"

	// A list of permanent MAC addresses of Wi-Fi devices to which
	// this connection should never apply. Each MAC address should be
	// given in the standard hex-digits-and-colons notation (eg
	// '00:11:22:33:44:55').
	NM_SETTING_WIRELESS_MAC_ADDRESS_BLACKLIST = "mac-address-blacklist"

	// If non-zero, only transmit packets of the specified size or
	// smaller, breaking larger packets up into multiple Ethernet
	// frames.
	NM_SETTING_WIRELESS_MTU = "mtu"

	// A list of BSSIDs (each BSSID formatted as a MAC address like
	// 00:11:22:33:44:55') that have been detected as part of the WiFI
	// network. NetworkManager internally tracks previously seen
	// BSSIDs. The property is only meant for reading and reflects the
	// BBSID list of NetworkManager. The changes you make to this
	// property will not be preserved.
	NM_SETTING_WIRELESS_SEEN_BSSIDS = "seen-bssids"

	// If the wireless connection has any security restrictions, like
	// 802.1x, WEP, or WPA, set this property to
	// '802-11-wireless-security' and ensure the connection contains a
	// valid 802-11-wireless-security setting.
	NM_SETTING_WIRELESS_SEC = "security"

	// If TRUE, indicates this network is a non-broadcasting network
	// that hides its SSID. In this case various workarounds may take
	// place, such as probe-scanning the SSID for more reliable
	// network discovery. However, these workarounds expose inherent
	// insecurities with hidden SSID networks, and thus hidden SSID
	// networks should be used with caution.
	NM_SETTING_WIRELESS_HIDDEN = "hidden"
)

const (
	// Indicates Ad-Hoc mode where no access point is expected to be
	// present.
	NM_SETTING_WIRELESS_MODE_ADHOC = "adhoc"

	// Indicates AP/master mode where the wireless device is started
	// as an access point/hotspot.
	//
	// Since: 0.9.8
	NM_SETTING_WIRELESS_MODE_AP = "ap"

	// Indicates infrastructure mode where an access point is expected
	// to be present for this connection.
	NM_SETTING_WIRELESS_MODE_INFRA = "infrastructure"
)

func newWirelessConnectionData(id, uuid string, ssid []byte, keyFlag int) (data _ConnectionData) {
	data = make(_ConnectionData)

	addSettingField(data, fieldConnection)
	setSettingConnectionId(data, id)
	setSettingConnectionUuid(data, uuid)
	setSettingConnectionType(data, typeWireless)

	addSettingField(data, fieldWireless)
	setSettingWirelessSsid(data, ssid)

	if keyFlag != ApKeyNone {
		addSettingField(data, fieldWirelessSecurity)
		setSettingWirelessSec(data, fieldWirelessSecurity)
		switch keyFlag {
		case ApKeyWep:
			setSettingWirelessSecurityKeyMgmt(data, "none")
		case ApKeyPsk:
			setSettingWirelessSecurityKeyMgmt(data, "wpa-psk")
			setSettingWirelessSecurityAuthAlg(data, "open")
		case ApKeyEap:
			setSettingWirelessSecurityKeyMgmt(data, "wpa-eap")
			setSettingWirelessSecurityAuthAlg(data, "open")
		}
	}

	addSettingField(data, fieldIPv4)
	setSettingIp4ConfigMethod(data, NM_SETTING_IP4_CONFIG_METHOD_AUTO)

	addSettingField(data, fieldIPv6)
	setSettingIp6ConfigMethod(data, NM_SETTING_IP6_CONFIG_METHOD_AUTO)

	return

	// TODO remove

	// data[fieldConnection] = make(map[string]dbus.Variant)
	// data[fieldIPv4] = make(map[string]dbus.Variant)
	// data[fieldIPv6] = make(map[string]dbus.Variant)
	// data[fieldWireless] = make(map[string]dbus.Variant)

	// data[fieldConnection]["id"] = dbus.MakeVariant(id)
	// uuid := newUUID()
	// data[fieldConnection]["uuid"] = dbus.MakeVariant(uuid)
	// data[fieldConnection]["type"] = dbus.MakeVariant(fieldWireless)

	// data[fieldWireless]["ssid"] = dbus.MakeVariant([]uint8(ssid))

	// if keyFlag != ApKeyNone {
	// 	data[fieldWirelessSecurity] = make(map[string]dbus.Variant)
	// 	data[fieldWireless]["security"] = dbus.MakeVariant(fieldWirelessSecurity)
	// 	switch keyFlag {
	// 	case ApKeyWep:
	// 		data[fieldWirelessSecurity]["key-mgmt"] = dbus.MakeVariant("none")
	// 	case ApKeyPsk:
	// 		data[fieldWirelessSecurity]["key-mgmt"] = dbus.MakeVariant("wpa-psk")
	// 		data[fieldWirelessSecurity]["auth-alg"] = dbus.MakeVariant("open")
	// 	case ApKeyEap:
	// 		data[fieldWirelessSecurity]["key-mgmt"] = dbus.MakeVariant("wpa-eap")
	// 		data[fieldWirelessSecurity]["auth-alg"] = dbus.MakeVariant("open")
	// 	}
	// }

	// data[fieldIPv4]["method"] = dbus.MakeVariant("auto")
	// data[fieldIPv6]["method"] = dbus.MakeVariant("auto")
}

// TODO Get available keys
func getSettingWirelessAvailableKeys(data _ConnectionData) (keys []string) {
	keys = []string{
		NM_SETTING_WIRELESS_MODE,
	}
	return
}

// Get available values
func getSettingWirelessAvailableValues(key string) (values []string, customizable bool) {
	customizable = true
	switch key {
	case NM_SETTING_WIRELESS_MODE:
		values = []string{
			NM_SETTING_WIRELESS_MODE_ADHOC,
			NM_SETTING_WIRELESS_MODE_AP,
			NM_SETTING_WIRELESS_MODE_INFRA,
		}
		customizable = false
	case NM_SETTING_WIRELESS_SEC:
		values = []string{
			fieldWirelessSecurity,
			field8021x,
		}
		customizable = true
	}
	return
}

// Check whether the values are correct
func checkSettingWirelessValues(data _ConnectionData) (errs map[string]string) {
	errs = make(map[string]string)

	// check ssid
	if !isSettingWirelessSsidExists(data) {
		rememberError(errs, NM_SETTING_WIRELESS_SSID, NM_KEY_ERROR_MISSING_VALUE)
		return
	}
	ssid := getSettingWirelessSsid(data)
	if len(ssid) == 0 {
		rememberError(errs, NM_SETTING_WIRELESS_SSID, NM_KEY_ERROR_EMPTY_VALUE)
		return
	}

	// check security
	if isSettingWirelessSecExists(data) {
		securityField := getSettingWirelessSec(data)
		if !isSettingFieldExists(data, securityField) {
			rememberError(errs, NM_SETTING_WIRELESS_SEC, fmt.Sprintf(NM_KEY_ERROR_MISSING_SECTION, securityField))
		}
	}

	return
}

// Set JSON value generally
func generalSetSettingWirelessKeyJSON(data _ConnectionData, key, valueJSON string) {
	switch key {
	default:
		LOGGER.Error("generalSetSettingWirelessKey: invalide key", key)
	case NM_SETTING_WIRELESS_SSID:
		setSettingWirelessSsidJSON(data, valueJSON)
	case NM_SETTING_WIRELESS_MODE:
		setSettingWirelessModeJSON(data, valueJSON)
	case NM_SETTING_WIRELESS_BAND:
		setSettingWirelessBandJSON(data, valueJSON)
	case NM_SETTING_WIRELESS_CHANNEL:
		setSettingWirelessChannelJSON(data, valueJSON)
	case NM_SETTING_WIRELESS_BSSID:
		setSettingWirelessBssidJSON(data, valueJSON)
	case NM_SETTING_WIRELESS_RATE:
		setSettingWirelessRateJSON(data, valueJSON)
	case NM_SETTING_WIRELESS_TX_POWER:
		setSettingWirelessTxPowerJSON(data, valueJSON)
	case NM_SETTING_WIRELESS_MAC_ADDRESS:
		setSettingWirelessMacAddressJSON(data, valueJSON)
	case NM_SETTING_WIRELESS_CLONED_MAC_ADDRESS:
		setSettingWirelessClonedMacAddressJSON(data, valueJSON)
	case NM_SETTING_WIRELESS_MAC_ADDRESS_BLACKLIST:
		setSettingWirelessMacAddressBlacklistJSON(data, valueJSON)
	case NM_SETTING_WIRELESS_MTU:
		setSettingWirelessMtuJSON(data, valueJSON)
	case NM_SETTING_WIRELESS_SEEN_BSSIDS:
		setSettingWirelessSeenBssidsJSON(data, valueJSON)
	case NM_SETTING_WIRELESS_SEC:
		setSettingWirelessSecJSON(data, valueJSON)
	case NM_SETTING_WIRELESS_HIDDEN:
		setSettingWirelessHiddenJSON(data, valueJSON)
	}
	return
}
