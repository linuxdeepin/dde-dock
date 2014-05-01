package main

import (
	"dlib"
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

func newWirelessConnection(id string, ssid []byte, secType apSecType) (uuid string) {
	logger.Debugf("new wireless connection, id=%s, ssid=%s, secType=%d", id, ssid, secType)
	uuid = newUUID()
	data := newWirelessConnectionData(id, uuid, ssid, secType)
	nmAddConnection(data)
	return
}

func newWirelessConnectionData(id, uuid string, ssid []byte, secType apSecType) (data connectionData) {
	data = make(connectionData)

	addSettingField(data, fieldConnection)
	setSettingConnectionId(data, id)
	setSettingConnectionUuid(data, uuid)
	setSettingConnectionType(data, typeWireless)

	addSettingField(data, fieldWireless)
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

	initSettingFieldIpv4(data)
	initSettingFieldIpv6(data)

	return
}

// Get available keys
func getSettingWirelessAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, fieldWireless, NM_SETTING_WIRELESS_SSID)
	keys = appendAvailableKeys(data, keys, fieldWireless, NM_SETTING_WIRELESS_MODE)
	switch getSettingWirelessMode(data) {
	case NM_SETTING_WIRELESS_MODE_INFRA:
	case NM_SETTING_WIRELESS_MODE_ADHOC:
		keys = appendAvailableKeys(data, keys, fieldWireless, NM_SETTING_WIRELESS_BAND)
		// TODO
		keys = appendAvailableKeys(data, keys, fieldWireless, NM_SETTING_WIRELESS_CHANNEL)
	case NM_SETTING_WIRELESS_MODE_AP:
		// TODO
	}
	keys = appendAvailableKeys(data, keys, fieldWireless, NM_SETTING_WIRELESS_MAC_ADDRESS)
	keys = appendAvailableKeys(data, keys, fieldWireless, NM_SETTING_WIRELESS_CLONED_MAC_ADDRESS)
	keys = appendAvailableKeys(data, keys, fieldWireless, NM_SETTING_WIRELESS_MTU)
	return
}

// Get available values
func getSettingWirelessAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case NM_SETTING_WIRELESS_MODE:
		values = []kvalue{
			kvalue{NM_SETTING_WIRELESS_MODE_INFRA, dlib.Tr("Infrastructure")},
			kvalue{NM_SETTING_WIRELESS_MODE_ADHOC, dlib.Tr("Ad-Hoc")},
			kvalue{NM_SETTING_WIRELESS_MODE_AP, dlib.Tr("AP-Hotspot")},
		}
	case NM_SETTING_WIRELESS_BAND:
		values = []kvalue{
			kvalue{"", dlib.Tr("Automatic")},
			kvalue{"a", dlib.Tr("A (5 GHz)")},
			kvalue{"bg", dlib.Tr("BG (2.4 GHz)")},
		}
	case NM_SETTING_WIRELESS_CHANNEL:
		// TODO
	case NM_SETTING_WIRELESS_MAC_ADDRESS:
		// get wireless devices mac address
		devPaths, err := nmGetDevices()
		if err == nil {
			for _, p := range devPaths {
				if dev, err := nmNewDevice(p); err == nil && dev.DeviceType.Get() == NM_DEVICE_TYPE_WIFI {
					hwAddr, err := nmGetWirelessDeviceHwAddr(p)
					if err == nil {
						values = append(values, kvalue{hwAddr, hwAddr + " (" + dev.Interface.Get() + ")"})
					}
				}
			}
		}
	}
	return
}

// Check whether the values are correct
func checkSettingWirelessValues(data connectionData) (errs fieldErrors) {
	errs = make(map[string]string)

	// check ssid
	ensureSettingWirelessSsidNoEmpty(data, errs)

	// check security
	if isSettingWirelessSecExists(data) {
		securityField := getSettingWirelessSec(data)
		if !isSettingFieldExists(data, securityField) {
			rememberError(errs, fieldWireless, NM_SETTING_WIRELESS_SEC, fmt.Sprintf(NM_KEY_ERROR_MISSING_SECTION, securityField))
		}
	}

	// machine address will be checked when setting key
	return
}
