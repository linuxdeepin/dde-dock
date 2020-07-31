/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package network

import (
	"os"

	"pkg.deepin.io/dde/daemon/network/nm"
	dbus "pkg.deepin.io/lib/dbus1"
)

func newWirelessHotspotConnectionForDevice(id, uuid string, devPath dbus.ObjectPath, active bool) (cpath dbus.ObjectPath, err error) {
	logger.Infof("new wireless hotspot connection, id=%s, uuid=%s, devPath=%s", id, uuid, devPath)
	data := newWirelessHotspotConnectionData(id, uuid)
	setSettingConnectionInterfaceName(data, nmGetDeviceInterface(devPath))
	setSettingWirelessSsid(data, []byte(os.Getenv("USER")))
	setSettingWirelessSecurityKeyMgmt(data, "none")
	hwAddr, _ := nmGeneralGetDeviceHwAddr(devPath, true)
	setSettingWirelessMacAddress(data, convertMacAddressToArrayByte(hwAddr))
	if active {
		cpath, _, err = nmAddAndActivateConnection(data, devPath, true)
	} else {
		cpath, err = nmAddConnection(data)
	}
	return
}

// new connection data
func newWirelessConnectionData(id, uuid string, ssid []byte, secType apSecType) (data connectionData) {
	logger.Debug("newWirelessConnectionData: secType:", secType)
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
	var err error 
	switch secType {
	case apSecNone:
		err = logicSetSettingVkWirelessSecurityKeyMgmt(data, "none")
	case apSecWep:
		err = logicSetSettingVkWirelessSecurityKeyMgmt(data, "wep")
	case apSecPsk:
		err = logicSetSettingVkWirelessSecurityKeyMgmt(data, "wpa-psk")
	case apSecEap:
		err = logicSetSettingVkWirelessSecurityKeyMgmt(data, "wpa-eap")
	}
	if err != nil {
		logger.Debug("failed to set VKWirelessSecutiryKeyMgmt")
		return 
	}

	initSettingSectionIpv4(data)
	initSettingSectionIpv6(data)

	return
}

func newWirelessHotspotConnectionData(id, uuid string) (data connectionData) {
	data = newWirelessConnectionData(id, uuid, nil, apSecNone)
	err := logicSetSettingWirelessMode(data, nm.NM_SETTING_WIRELESS_MODE_AP)
	if err != nil {
		logger.Debug("failed to set WirelessMode")
		return 
	}
	setSettingConnectionAutoconnect(data, false)
	return
}

// Logic setter
func logicSetSettingWirelessMode(data connectionData, value string) (err error) {
	// for ad-hoc or ap-hotspot mode, wpa-eap security is invalid, and
	// set ip4 method to "shared"
	if value != nm.NM_SETTING_WIRELESS_MODE_INFRA {
		if getSettingVkWirelessSecurityKeyMgmt(data) == "wpa-eap" {
			err = logicSetSettingVkWirelessSecurityKeyMgmt(data, "wpa-psk")
			if err != nil {
				logger.Debug("failed to set VkWirelessKeyMgmt")
				return err 
			}
		}
		setSettingIP4ConfigMethod(data, nm.NM_SETTING_IP4_CONFIG_METHOD_SHARED)
	}
	setSettingWirelessMode(data, value)
	return
}
