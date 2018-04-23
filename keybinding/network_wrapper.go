/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
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

package keybinding

import (
	"encoding/json"

	"github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.network"
	"pkg.deepin.io/lib/dbus1"
)

var (
	// enable by 'go build -ldflags "-X pkg.deepin.io/dde/daemon/keybinding.ManageWireless=enabled"'
	ManageWireless = "disabled"
)

func toggleWireless(sessionConn *dbus.Conn) error {
	net := network.NewNetwork(sessionConn)
	devices, err := net.Devices().Get(0)
	if err != nil {
		return err
	}
	list := getWirelessDevice(devices)
	enabled := false
	for _, dev := range list {
		ok, _ := net.IsDeviceEnabled(0, dbus.ObjectPath(dev))
		if ok {
			enabled = true
			break
		}
	}

	for _, dev := range list {
		net.EnableDevice(0, dbus.ObjectPath(dev), !enabled)
	}
	return nil
}

type deviceInfo struct {
	Path string `json:"Path"`
}

type wirelessDevice struct {
	Devices []deviceInfo `json:"wireless"`
}

func getWirelessDevice(value string) []string {
	var wireless wirelessDevice
	err := json.Unmarshal([]byte(value), &wireless)
	if err != nil {
		return nil
	}
	var list []string
	for _, dev := range wireless.Devices {
		list = append(list, dev.Path)
	}
	return list
}
