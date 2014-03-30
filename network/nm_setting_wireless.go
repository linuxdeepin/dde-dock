package main

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

func newWirelessConnectionData(id, uuid, ssid string, keyFlag int) (data _ConnectionData) {
	data = make(_ConnectionData)

	addConnectionDataField(data, fieldConnection)
	setSettingConnectionId(data, id)
	setSettingConnectionUuid(data, uuid)
	setSettingConnectionType(data, typeWireless)

	addConnectionDataField(data, fieldWireless)
	setSettingWirelessSsid(data, ssid)

	if keyFlag != ApKeyNone {
		addConnectionDataField(data, fieldWirelessSecurity)
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

	addConnectionDataField(data, fieldIPv4)
	setSettingIp4ConfigMethod(data, NM_SETTING_IP4_CONFIG_METHOD_AUTO)

	addConnectionDataField(data, fieldIPv6)
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

// Get key type
func getSettingWirelessKeyType(key string) (t ktype) {
	switch key {
	default:
		t = ktypeUnknown
	case NM_SETTING_WIRELESS_SSID:
		t = ktypeArrayByte
	case NM_SETTING_WIRELESS_MODE:
		t = ktypeString
	case NM_SETTING_WIRELESS_BAND:
		t = ktypeString
	case NM_SETTING_WIRELESS_CHANNEL:
		t = ktypeUint32
	case NM_SETTING_WIRELESS_BSSID:
		t = ktypeArrayByte
	case NM_SETTING_WIRELESS_RATE:
		t = ktypeUint32
	case NM_SETTING_WIRELESS_TX_POWER:
		t = ktypeUint32
	case NM_SETTING_WIRELESS_MAC_ADDRESS:
		t = ktypeArrayByte
	case NM_SETTING_WIRELESS_CLONED_MAC_ADDRESS:
		t = ktypeArrayByte
	case NM_SETTING_WIRELESS_MAC_ADDRESS_BLACKLIST:
		t = ktypeArrayString
	case NM_SETTING_WIRELESS_MTU:
		t = ktypeUint32
	case NM_SETTING_WIRELESS_SEEN_BSSIDS:
		t = ktypeArrayString
	case NM_SETTING_WIRELESS_SEC:
		t = ktypeString
	case NM_SETTING_WIRELESS_HIDDEN:
		t = ktypeBoolean
	}
	return
}

// Get and set key's value generally
func generalGetSettingWirelessKey(data _ConnectionData, key string) (value string) {
	switch key {
	default:
		LOGGER.Error("generalGetSettingWirelessKey: invalide key", key)
	case NM_SETTING_WIRELESS_SSID:
		value = getSettingWirelessSsid(data)
	case NM_SETTING_WIRELESS_MODE:
		value = getSettingWirelessMode(data)
	case NM_SETTING_WIRELESS_BAND:
		value = getSettingWirelessBand(data)
	case NM_SETTING_WIRELESS_CHANNEL:
		value = getSettingWirelessChannel(data)
	case NM_SETTING_WIRELESS_BSSID:
		value = getSettingWirelessBssid(data)
	case NM_SETTING_WIRELESS_RATE:
		value = getSettingWirelessRate(data)
	case NM_SETTING_WIRELESS_TX_POWER:
		value = getSettingWirelessTxPower(data)
	case NM_SETTING_WIRELESS_MAC_ADDRESS:
		value = getSettingWirelessMacAddress(data)
	case NM_SETTING_WIRELESS_CLONED_MAC_ADDRESS:
		value = getSettingWirelessClonedMacAddress(data)
	case NM_SETTING_WIRELESS_MAC_ADDRESS_BLACKLIST:
		value = getSettingWirelessMacAddressBlacklist(data)
	case NM_SETTING_WIRELESS_MTU:
		value = getSettingWirelessMtu(data)
	case NM_SETTING_WIRELESS_SEEN_BSSIDS:
		value = getSettingWirelessSeenBssids(data)
	case NM_SETTING_WIRELESS_SEC:
		value = getSettingWirelessSec(data)
	case NM_SETTING_WIRELESS_HIDDEN:
		value = getSettingWirelessHidden(data)
	}
	return
}

// TODO use logic setter
func generalSetSettingWirelessKey(data _ConnectionData, key, value string) {
	switch key {
	default:
		LOGGER.Error("generalSetSettingWirelessKey: invalide key", key)
	case NM_SETTING_WIRELESS_SSID:
		setSettingWirelessSsid(data, value)
	case NM_SETTING_WIRELESS_MODE:
		setSettingWirelessMode(data, value)
	case NM_SETTING_WIRELESS_BAND:
		setSettingWirelessBand(data, value)
	case NM_SETTING_WIRELESS_CHANNEL:
		setSettingWirelessChannel(data, value)
	case NM_SETTING_WIRELESS_BSSID:
		setSettingWirelessBssid(data, value)
	case NM_SETTING_WIRELESS_RATE:
		setSettingWirelessRate(data, value)
	case NM_SETTING_WIRELESS_TX_POWER:
		setSettingWirelessTxPower(data, value)
	case NM_SETTING_WIRELESS_MAC_ADDRESS:
		setSettingWirelessMacAddress(data, value)
	case NM_SETTING_WIRELESS_CLONED_MAC_ADDRESS:
		setSettingWirelessClonedMacAddress(data, value)
	case NM_SETTING_WIRELESS_MAC_ADDRESS_BLACKLIST:
		setSettingWirelessMacAddressBlacklist(data, value)
	case NM_SETTING_WIRELESS_MTU:
		setSettingWirelessMtu(data, value)
	case NM_SETTING_WIRELESS_SEEN_BSSIDS:
		setSettingWirelessSeenBssids(data, value)
	case NM_SETTING_WIRELESS_SEC:
		setSettingWirelessSec(data, value)
	case NM_SETTING_WIRELESS_HIDDEN:
		setSettingWirelessHidden(data, value)
	}
	return
}

// Getter
func getSettingWirelessSsid(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_SSID, getSettingWirelessKeyType(NM_SETTING_WIRELESS_SSID))
	return
}
func getSettingWirelessMode(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_MODE, getSettingWirelessKeyType(NM_SETTING_WIRELESS_MODE))
	return
}
func getSettingWirelessBand(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_BAND, getSettingWirelessKeyType(NM_SETTING_WIRELESS_BAND))
	return
}
func getSettingWirelessChannel(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_CHANNEL, getSettingWirelessKeyType(NM_SETTING_WIRELESS_CHANNEL))
	return
}
func getSettingWirelessBssid(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_BSSID, getSettingWirelessKeyType(NM_SETTING_WIRELESS_BSSID))
	return
}
func getSettingWirelessRate(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_RATE, getSettingWirelessKeyType(NM_SETTING_WIRELESS_RATE))
	return
}
func getSettingWirelessTxPower(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_TX_POWER, getSettingWirelessKeyType(NM_SETTING_WIRELESS_TX_POWER))
	return
}
func getSettingWirelessMacAddress(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_MAC_ADDRESS, getSettingWirelessKeyType(NM_SETTING_WIRELESS_MAC_ADDRESS))
	return
}
func getSettingWirelessClonedMacAddress(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_CLONED_MAC_ADDRESS, getSettingWirelessKeyType(NM_SETTING_WIRELESS_CLONED_MAC_ADDRESS))
	return
}
func getSettingWirelessMacAddressBlacklist(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_MAC_ADDRESS_BLACKLIST, getSettingWirelessKeyType(NM_SETTING_WIRELESS_MAC_ADDRESS_BLACKLIST))
	return
}
func getSettingWirelessMtu(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_MTU, getSettingWirelessKeyType(NM_SETTING_WIRELESS_MTU))
	return
}
func getSettingWirelessSeenBssids(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_SEEN_BSSIDS, getSettingWirelessKeyType(NM_SETTING_WIRELESS_SEEN_BSSIDS))
	return
}
func getSettingWirelessSec(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_SEC, getSettingWirelessKeyType(NM_SETTING_WIRELESS_SEC))
	return
}
func getSettingWirelessHidden(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_HIDDEN, getSettingWirelessKeyType(NM_SETTING_WIRELESS_HIDDEN))
	return
}

// Setter
func setSettingWirelessSsid(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_SSID, value, getSettingWirelessKeyType(NM_SETTING_WIRELESS_SSID))
}
func setSettingWirelessMode(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_MODE, value, getSettingWirelessKeyType(NM_SETTING_WIRELESS_MODE))
}
func setSettingWirelessBand(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_BAND, value, getSettingWirelessKeyType(NM_SETTING_WIRELESS_BAND))
}
func setSettingWirelessChannel(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_CHANNEL, value, getSettingWirelessKeyType(NM_SETTING_WIRELESS_CHANNEL))
}
func setSettingWirelessBssid(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_BSSID, value, getSettingWirelessKeyType(NM_SETTING_WIRELESS_BSSID))
}
func setSettingWirelessRate(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_RATE, value, getSettingWirelessKeyType(NM_SETTING_WIRELESS_RATE))
}
func setSettingWirelessTxPower(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_TX_POWER, value, getSettingWirelessKeyType(NM_SETTING_WIRELESS_TX_POWER))
}
func setSettingWirelessMacAddress(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_MAC_ADDRESS, value, getSettingWirelessKeyType(NM_SETTING_WIRELESS_MAC_ADDRESS))
}
func setSettingWirelessClonedMacAddress(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_CLONED_MAC_ADDRESS, value, getSettingWirelessKeyType(NM_SETTING_WIRELESS_CLONED_MAC_ADDRESS))
}
func setSettingWirelessMacAddressBlacklist(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_MAC_ADDRESS_BLACKLIST, value, getSettingWirelessKeyType(NM_SETTING_WIRELESS_MAC_ADDRESS_BLACKLIST))
}
func setSettingWirelessMtu(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_MTU, value, getSettingWirelessKeyType(NM_SETTING_WIRELESS_MTU))
}
func setSettingWirelessSeenBssids(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_SEEN_BSSIDS, value, getSettingWirelessKeyType(NM_SETTING_WIRELESS_SEEN_BSSIDS))
}
func setSettingWirelessSec(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_SEC, value, getSettingWirelessKeyType(NM_SETTING_WIRELESS_SEC))
}
func setSettingWirelessHidden(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_HIDDEN, value, getSettingWirelessKeyType(NM_SETTING_WIRELESS_HIDDEN))
}

// Remover
func removeSettingWirelessSsid(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_SSID)
}
func removeSettingWirelessMode(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_MODE)
}
func removeSettingWirelessBand(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_BAND)
}
func removeSettingWirelessChannel(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_CHANNEL)
}
func removeSettingWirelessBssid(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_BSSID)
}
func removeSettingWirelessRate(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_RATE)
}
func removeSettingWirelessTxPower(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_TX_POWER)
}
func removeSettingWirelessMacAddress(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_MAC_ADDRESS)
}
func removeSettingWirelessClonedMacAddress(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_CLONED_MAC_ADDRESS)
}
func removeSettingWirelessMacAddressBlacklist(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_MAC_ADDRESS_BLACKLIST)
}
func removeSettingWirelessMtu(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_MTU)
}
func removeSettingWirelessSeenBssids(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_SEEN_BSSIDS)
}
func removeSettingWirelessSec(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_SEC)
}
func removeSettingWirelessHidden(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_HIDDEN)
}
