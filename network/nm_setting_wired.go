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
	"pkg.deepin.io/dde/daemon/network/nm"
	"pkg.deepin.io/lib/dbus"
)

func newWiredConnectionForDevice(id, uuid string, devPath dbus.ObjectPath, active bool) (cpath dbus.ObjectPath, err error) {
	logger.Infof("new wired connection, id=%s, uuid=%s, devPath=%s", id, uuid, devPath)
	data := newWiredConnectionData(id, uuid)
	hwAddr, _ := nmGeneralGetDeviceHwAddr(devPath, true)
	setSettingWiredMacAddress(data, convertMacAddressToArrayByte(hwAddr))
	setSettingConnectionAutoconnect(data, true)
	if active {
		cpath, _, err = nmAddAndActivateConnection(data, devPath, false)
	} else {
		cpath, err = nmAddConnection(data)
	}
	return
}

func newWiredConnectionData(id, uuid string) (data connectionData) {
	data = make(connectionData)

	addSetting(data, nm.NM_SETTING_CONNECTION_SETTING_NAME)
	setSettingConnectionId(data, id)
	setSettingConnectionUuid(data, uuid)
	setSettingConnectionType(data, nm.NM_SETTING_WIRED_SETTING_NAME)

	initSettingSectionWired(data)

	initSettingSectionIpv4(data)
	initSettingSectionIpv6(data)
	return
}

func initSettingSectionWired(data connectionData) {
	addSetting(data, nm.NM_SETTING_WIRED_SETTING_NAME)
	setSettingWiredDuplex(data, "full")
}

// Get available keys
func getSettingWiredAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_WIRED_SETTING_NAME, nm.NM_SETTING_WIRED_MAC_ADDRESS)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_WIRED_SETTING_NAME, nm.NM_SETTING_WIRED_CLONED_MAC_ADDRESS)
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_WIRED_SETTING_NAME, nm.NM_SETTING_WIRED_MTU)
	if isSettingWiredMtuExists(data) {
		keys = append(keys, nm.NM_SETTING_WIRED_MTU)
	}
	return
}

// Get available values
func getSettingWiredAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case nm.NM_SETTING_WIRED_MAC_ADDRESS:
		// get all wired devices mac address
		for iface, hwAddr := range nmGeneralGetAllDeviceHwAddr(nm.NM_DEVICE_TYPE_ETHERNET) {
			values = append(values, kvalue{hwAddr, hwAddr + " (" + iface + ")"})
		}
	}
	return
}

// Check whether the values are correct
func checkSettingWiredValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)
	// hardware address will be checked when setting key
	return
}
