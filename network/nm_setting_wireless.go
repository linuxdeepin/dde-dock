package network

import (
	. "dlib/gettext"
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

// initialize available values
var availableValuesWirelessChannelA = []kvalue{
	kvalue{"", "Default"},
	kvalue{"7", "7 (5035 MHz)"},
	kvalue{"8", "8 (5040 MHz)"},
	kvalue{"9", "9 (5045 MHz)"},
	kvalue{"11", "11 (5055 MHz)"},
	kvalue{"12", "12 (5060 MHz)"},
	kvalue{"16", "16 (5080 MHz)"},
	kvalue{"34", "34 (5170 MHz)"},
	kvalue{"36", "36 (5180 MHz)"},
	kvalue{"38", "38 (5190 MHz)"},
	kvalue{"40", "40 (5200 MHz)"},
	kvalue{"42", "42 (5210 MHz)"},
	kvalue{"44", "44 (5220 MHz)"},
	kvalue{"46", "46 (5230 MHz)"},
	kvalue{"48", "48 (5240 MHz)"},
	kvalue{"50", "50 (5250 MHz)"},
	kvalue{"52", "52 (5260 MHz)"},
	kvalue{"56", "56 (5280 MHz)"},
	kvalue{"58", "58 (5290 MHz)"},
	kvalue{"60", "60 (5300 MHz)"},
	kvalue{"64", "64 (5320 MHz)"},
	kvalue{"100", "100 (5500 MHz)"},
	kvalue{"104", "104 (5520 MHz)"},
	kvalue{"108", "108 (5540 MHz)"},
	kvalue{"112", "112 (5560 MHz)"},
	kvalue{"116", "116 (5580 MHz)"},
	kvalue{"120", "120 (5600 MHz)"},
	kvalue{"124", "124 (5620 MHz)"},
	kvalue{"128", "128 (5640 MHz)"},
	kvalue{"132", "132 (5660 MHz)"},
	kvalue{"136", "136 (5680 MHz)"},
	kvalue{"140", "140 (5700 MHz)"},
	kvalue{"149", "149 (5745 MHz)"},
	kvalue{"152", "152 (5760 MHz)"},
	kvalue{"153", "153 (5765 MHz)"},
	kvalue{"157", "157 (5785 MHz)"},
	kvalue{"160", "160 (5800 MHz)"},
	kvalue{"161", "161 (5805 MHz)"},
	kvalue{"165", "165 (5825 MHz)"},
	kvalue{"183", "183 (4915 MHz)"},
	kvalue{"184", "184 (4920 MHz)"},
	kvalue{"185", "185 (4925 MHz)"},
	kvalue{"187", "187 (4935 MHz)"},
	kvalue{"188", "188 (4945 MHz)"},
	kvalue{"192", "192 (4960 MHz)"},
	kvalue{"196", "196 (4980 MHz)"},
}

var availableValuesWirelessChannelBg = []kvalue{
	kvalue{"", Tr("Default")},
	kvalue{"1", "1 (2412 MHz)"},
	kvalue{"2", "2 (2417 MHz)"},
	kvalue{"3", "3 (2422 MHz)"},
	kvalue{"4", "4 (2427 MHz)"},
	kvalue{"5", "5 (2432 MHz)"},
	kvalue{"6", "6 (2437 MHz)"},
	kvalue{"7", "7 (2442 MHz)"},
	kvalue{"8", "8 (2447 MHz)"},
	kvalue{"9", "9 (2452 MHz)"},
	kvalue{"10", "10 (2457 MHz)"},
	kvalue{"11", "11 (2462 MHz)"},
	kvalue{"12", "12 (2467 MHz)"},
	kvalue{"13", "13 (2472 MHz)"},
	kvalue{"14", "14 (2484 MHz)"},
}

// new connection data
func newWirelessConnection(id string, ssid []byte, secType apSecType) (uuid string) {
	logger.Debugf("new wireless connection, id=%s, ssid=%s, secType=%d", id, ssid, secType)
	uuid = genUuid()
	data := newWirelessConnectionData(id, uuid, ssid, secType)
	nmAddConnection(data)
	return
}

func newWirelessConnectionData(id, uuid string, ssid []byte, secType apSecType) (data connectionData) {
	data = make(connectionData)

	addSettingSection(data, sectionConnection)
	setSettingConnectionId(data, id)
	setSettingConnectionUuid(data, uuid)
	setSettingConnectionType(data, NM_SETTING_WIRELESS_SETTING_NAME)

	addSettingSection(data, sectionWireless)
	setSettingWirelessSsid(data, ssid)
	setSettingWirelessMode(data, NM_SETTING_WIRELESS_MODE_INFRA)

	switch secType {
	case apSecNone:
		logicSetSettingVkWirelessSecurityKeyMgmt(data, "none")
	case apSecWep:
		logicSetSettingVkWirelessSecurityKeyMgmt(data, "wep")
	case apSecPsk:
		logicSetSettingVkWirelessSecurityKeyMgmt(data, "wpa-psk")
	case apSecEap:
		logicSetSettingVkWirelessSecurityKeyMgmt(data, "wpa-eap")
	}

	initSettingSectionIpv4(data)
	initSettingSectionIpv6(data)

	return
}

func newWirelessAdhocConnectionData(id, uuid string) (data connectionData) {
	data = newWirelessConnectionData(id, uuid, nil, apSecNone)
	logicSetSettingWirelessMode(data, NM_SETTING_WIRELESS_MODE_ADHOC)
	setSettingConnectionAutoconnect(data, false)
	return
}

func newWirelessHotspotConnectionData(id, uuid string) (data connectionData) {
	data = newWirelessConnectionData(id, uuid, nil, apSecNone)
	logicSetSettingWirelessMode(data, NM_SETTING_WIRELESS_MODE_AP)
	setSettingConnectionAutoconnect(data, false)
	return
}

// Get available keys
func getSettingWirelessAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, sectionWireless, NM_SETTING_WIRELESS_SSID)
	// hide wireless mode option for better user experience
	// keys = appendAvailableKeys(data, keys, sectionWireless, NM_SETTING_WIRELESS_MODE)
	switch getSettingWirelessMode(data) {
	case NM_SETTING_WIRELESS_MODE_INFRA:
	case NM_SETTING_WIRELESS_MODE_ADHOC:
		keys = appendAvailableKeys(data, keys, sectionWireless, NM_SETTING_WIRELESS_BAND)
		if isSettingWirelessBandExists(data) {
			keys = appendAvailableKeys(data, keys, sectionWireless, NM_SETTING_WIRELESS_CHANNEL)
		}
	case NM_SETTING_WIRELESS_MODE_AP:
		keys = appendAvailableKeys(data, keys, sectionWireless, NM_SETTING_WIRELESS_BAND)
		if isSettingWirelessBandExists(data) {
			keys = appendAvailableKeys(data, keys, sectionWireless, NM_SETTING_WIRELESS_CHANNEL)
		}
	}
	keys = appendAvailableKeys(data, keys, sectionWireless, NM_SETTING_WIRELESS_MAC_ADDRESS)
	keys = appendAvailableKeys(data, keys, sectionWireless, NM_SETTING_WIRELESS_CLONED_MAC_ADDRESS)
	keys = appendAvailableKeys(data, keys, sectionWireless, NM_SETTING_WIRELESS_MTU)
	return
}

// Get available values
func getSettingWirelessAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case NM_SETTING_WIRELESS_MODE:
		values = []kvalue{
			kvalue{NM_SETTING_WIRELESS_MODE_INFRA, Tr("Infrastructure")},
			kvalue{NM_SETTING_WIRELESS_MODE_ADHOC, Tr("Ad-Hoc")},
			kvalue{NM_SETTING_WIRELESS_MODE_AP, Tr("AP-Hotspot")},
		}
	case NM_SETTING_WIRELESS_BAND:
		values = []kvalue{
			kvalue{"", Tr("Automatic")},
			kvalue{"a", Tr("A (5 GHz)")},
			kvalue{"bg", Tr("BG (2.4 GHz)")},
		}
	case NM_SETTING_WIRELESS_CHANNEL:
		if isSettingWirelessBandExists(data) {
			switch getSettingWirelessBand(data) {
			case "a":
				values = availableValuesWirelessChannelA
			case "bg":
				values = availableValuesWirelessChannelBg
			}
		}
	case NM_SETTING_WIRELESS_MAC_ADDRESS:
		// get wireless devices mac address
		for iface, hwAddr := range nmGeneralGetAllDeviceHwAddr(NM_DEVICE_TYPE_WIFI) {
			values = append(values, kvalue{hwAddr, hwAddr + " (" + iface + ")"})
		}
	}
	return
}

// Check whether the values are correct
func checkSettingWirelessValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)

	// check ssid
	ensureSettingWirelessSsidNoEmpty(data, errs)

	// check security
	if isSettingWirelessSecExists(data) {
		securitySection := getSettingWirelessSec(data)
		if !isSettingSectionExists(data, securitySection) {
			rememberError(errs, sectionWireless, NM_SETTING_WIRELESS_SEC, fmt.Sprintf(NM_KEY_ERROR_MISSING_SECTION, securitySection))
		}
	}

	// machine address will be checked when setting key
	return
}

// Logic setter
func logicSetSettingWirelessMode(data connectionData, value string) (err error) {
	// for ad-hoc or ap-hotspot mode, wpa-eap security is invalid
	if value != NM_SETTING_WIRELESS_MODE_INFRA {
		if getSettingVkWirelessSecurityKeyMgmt(data) == "wpa-eap" {
			logicSetSettingVkWirelessSecurityKeyMgmt(data, "wpa-psk")
		}
	}
	setSettingWirelessMode(data, value)
	return
}
func logicSetSettingWirelessBand(data connectionData, value string) (err error) {
	removeSettingWirelessChannel(data)
	setSettingWirelessBand(data, value)
	return
}
