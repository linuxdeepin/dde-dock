/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package keybinding

import (
	"dbus/com/deepin/daemon/network"
	"encoding/json"
	"pkg.deepin.io/lib/dbus"
)

var (
	// enable by 'go build -ldflags "-X pkg.deepin.io/dde/daemon/keybinding.ManageWireless=enabled"'
	ManageWireless = "disabled"
)

func toggleWireless() error {
	net, err := Network.NewNetworkManager("com.deepin.daemon.Network",
		"/com/deepin/daemon/Network")
	if err != nil {
		return err
	}
	defer Network.DestroyNetworkManager(net)

	list := getWirelessDevice(net.Devices.Get())
	enabled := false
	for _, dev := range list {
		ok, _ := net.IsDeviceEnabled(dbus.ObjectPath(dev))
		if ok {
			enabled = true
			break
		}
	}

	for _, dev := range list {
		net.EnableDevice(dbus.ObjectPath(dev), !enabled)
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
