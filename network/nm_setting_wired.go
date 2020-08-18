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
	"net"

	dbus "github.com/godbus/dbus"
	"pkg.deepin.io/dde/daemon/network/nm"
)

func newWiredConnectionForDevice(id, uuid string, devPath dbus.ObjectPath, active bool) (cpath dbus.ObjectPath, err error) {
	logger.Infof("new wired connection, id=%s, uuid=%s, devPath=%s", id, uuid, devPath)
	data := newWiredConnectionData(id, uuid, devPath)

	setSettingConnectionAutoconnect(data, true)
	cpath, err = nmAddConnection(data)
	if err != nil {
		return "/", err
	}
	if active {
		_, err = nmActivateConnection(cpath, devPath)
		if err != nil {
			logger.Warningf("failed to activate connection cpath: %v, devPath: %v, err: %v",
				cpath, devPath, err)
		}
	}
	return cpath, nil
}

func newWiredConnectionData(id, uuid string, devPath dbus.ObjectPath) (data connectionData) {
	data = make(connectionData)

	addSetting(data, nm.NM_SETTING_CONNECTION_SETTING_NAME)
	setSettingConnectionId(data, id)
	setSettingConnectionUuid(data, uuid)
	setSettingConnectionType(data, nm.NM_SETTING_WIRED_SETTING_NAME)

	initSettingSectionWired(data, devPath)

	initSettingSectionIpv4(data)
	initSettingSectionIpv6(data)
	return
}

func initSettingSectionWired(data connectionData, devPath dbus.ObjectPath) {
	addSetting(data, nm.NM_SETTING_WIRED_SETTING_NAME)
	setSettingWiredDuplex(data, "full")

	hwAddr, err := nmGeneralGetDeviceHwAddr(devPath, true)
	if err != nil {
		return
	}

	macAddr, err := net.ParseMAC(hwAddr)
	if err != nil {
		return
	}

	setSettingWiredMacAddress(data, macAddr)
}
