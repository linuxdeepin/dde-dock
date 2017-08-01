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
	"os"
	"pkg.deepin.io/dde/daemon/network/nm"
	"pkg.deepin.io/lib/dbus"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/utils"
)

// initialize available values
var availableValuesWirelessChannelA []kvalue
var availableValuesWirelessChannelBg []kvalue

func initAvailableValuesWirelessChannel() {
	availableValuesWirelessChannelA = []kvalue{
		kvalue{"", Tr("Default")},
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

	availableValuesWirelessChannelBg = []kvalue{
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
}

func newWirelessHotspotConnectionForDevice(id, uuid string, devPath dbus.ObjectPath, active bool) (cpath dbus.ObjectPath, err error) {
	logger.Infof("new wireless hotspot connection, id=%s, uuid=%s, devPath=%s", id, uuid, devPath)
	data := newWirelessHotspotConnectionData(id, uuid)
	setSettingConnectionInterfaceName(data, nmGetDeviceInterface(devPath))
	setSettingWirelessSsid(data, []byte(os.Getenv("USER")+"-hotspot"))
	setSettingWirelessSecurityKeyMgmt(data, "none")
	hwAddr, _ := nmGeneralGetDeviceHwAddr(devPath)
	setSettingWirelessMacAddress(data, convertMacAddressToArrayByte(hwAddr))
	if active {
		cpath, _, err = nmAddAndActivateConnection(data, devPath, true)
	} else {
		cpath, err = nmAddConnection(data)
	}
	return
}

func newWirelessConnection(id string, ssid []byte, secType apSecType) (uuid string) {
	logger.Debugf("new wireless connection, id=%s, ssid=%s, secType=%d", id, ssid, secType)
	uuid = utils.GenUuid()
	data := newWirelessConnectionData(id, uuid, ssid, secType)
	nmAddConnection(data)
	return
}

// new connection data
func newWirelessConnectionData(id, uuid string, ssid []byte, secType apSecType) (data connectionData) {
	data = make(connectionData)

	addSetting(data, nm.NM_SETTING_CONNECTION_SETTING_NAME)
	setSettingConnectionId(data, id)
	setSettingConnectionUuid(data, uuid)
	setSettingConnectionType(data, nm.NM_SETTING_WIRELESS_SETTING_NAME)

	addSetting(data, nm.NM_SETTING_WIRELESS_SETTING_NAME)
	if ssid != nil {
		setSettingWirelessSsid(data, ssid)
	}
	setSettingWirelessMode(data, nm.NM_SETTING_WIRELESS_MODE_INFRA)

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
	logicSetSettingWirelessMode(data, nm.NM_SETTING_WIRELESS_MODE_ADHOC)
	setSettingConnectionAutoconnect(data, false)
	return
}

func newWirelessHotspotConnectionData(id, uuid string) (data connectionData) {
	data = newWirelessConnectionData(id, uuid, nil, apSecNone)
	logicSetSettingWirelessMode(data, nm.NM_SETTING_WIRELESS_MODE_AP)
	setSettingConnectionAutoconnect(data, false)
	return
}

// Get available keys
func getSettingWirelessAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_WIRELESS_SETTING_NAME, nm.NM_SETTING_WIRELESS_SSID)
	switch getSettingWirelessMode(data) {
	case nm.NM_SETTING_WIRELESS_MODE_INFRA:
	case nm.NM_SETTING_WIRELESS_MODE_ADHOC:
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_WIRELESS_SETTING_NAME, nm.NM_SETTING_WIRELESS_BAND)
		if isSettingWirelessBandExists(data) {
			keys = appendAvailableKeys(data, keys, nm.NM_SETTING_WIRELESS_SETTING_NAME, nm.NM_SETTING_WIRELESS_CHANNEL)
		}
	case nm.NM_SETTING_WIRELESS_MODE_AP:
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_WIRELESS_SETTING_NAME, nm.NM_SETTING_WIRELESS_BAND)
		if isSettingWirelessBandExists(data) {
			keys = appendAvailableKeys(data, keys, nm.NM_SETTING_WIRELESS_SETTING_NAME, nm.NM_SETTING_WIRELESS_CHANNEL)
		}
	}
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_WIRELESS_SETTING_NAME, nm.NM_SETTING_WIRELESS_MAC_ADDRESS)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_WIRELESS_SETTING_NAME, nm.NM_SETTING_WIRELESS_MTU)

	// hide some wireless options for better user experience
	// keys = appendAvailableKeys(data, keys, nm.NM_SETTING_WIRELESS_SETTING_NAME, nm.NM_SETTING_WIRELESS_MODE)
	// keys = appendAvailableKeys(data, keys, nm.NM_SETTING_WIRELESS_SETTING_NAME, nm.NM_SETTING_WIRELESS_CLONED_MAC_ADDRESS)
	return
}

// Get available values
func getSettingWirelessAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case nm.NM_SETTING_WIRELESS_MODE:
		values = []kvalue{
			kvalue{nm.NM_SETTING_WIRELESS_MODE_INFRA, Tr("Infrastructure")},
			kvalue{nm.NM_SETTING_WIRELESS_MODE_ADHOC, Tr("Ad-Hoc")},
			kvalue{nm.NM_SETTING_WIRELESS_MODE_AP, Tr("AP-Hotspot")},
		}
	case nm.NM_SETTING_WIRELESS_BAND:
		values = []kvalue{
			kvalue{"", Tr("Automatic")},
			kvalue{"a", Tr("A (5 GHz)")},
			kvalue{"bg", Tr("BG (2.4 GHz)")},
		}
	case nm.NM_SETTING_WIRELESS_CHANNEL:
		if isSettingWirelessBandExists(data) {
			switch getSettingWirelessBand(data) {
			case "a":
				values = availableValuesWirelessChannelA
			case "bg":
				values = availableValuesWirelessChannelBg
			}
		}
	case nm.NM_SETTING_WIRELESS_MAC_ADDRESS:
		// get wireless devices mac address
		for iface, hwAddr := range nmGeneralGetAllDeviceHwAddr(nm.NM_DEVICE_TYPE_WIFI) {
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

	// machine address will be checked when setting key
	return
}

// Logic setter
func logicSetSettingWirelessMode(data connectionData, value string) (err error) {
	// for ad-hoc or ap-hotspot mode, wpa-eap security is invalid, and
	// set ip4 method to "shared"
	if value != nm.NM_SETTING_WIRELESS_MODE_INFRA {
		if getSettingVkWirelessSecurityKeyMgmt(data) == "wpa-eap" {
			logicSetSettingVkWirelessSecurityKeyMgmt(data, "wpa-psk")
		}
		setSettingIP4ConfigMethod(data, nm.NM_SETTING_IP4_CONFIG_METHOD_SHARED)
	}
	setSettingWirelessMode(data, value)
	return
}
func logicSetSettingWirelessBand(data connectionData, value string) (err error) {
	removeSettingWirelessChannel(data)
	setSettingWirelessBand(data, value)
	return
}
