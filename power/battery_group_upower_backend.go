/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package power

import (
	libupower "dbus/org/freedesktop/upower"
	"errors"
	"pkg.deepin.io/lib/dbus"
)

func newBatteryDeviceWithUPowerDBusObjPath(path dbus.ObjectPath) (batteryDevice, error) {
	dev, err := libupower.NewDevice(upowerDBusDest, path)
	if err != nil {
		logger.Warning("New battery device failed:", err)
		return nil, err
	}
	if dev.Type.Get() != DeviceTypeBattery {
		return nil, errors.New("device not battery")
	}
	return newUpowerBatteryDevice(dev), nil
}

func (batGroup *batteryGroup) AddUPowerBatteryDevice(path dbus.ObjectPath) {
	batGroup.changeLock.Lock()
	defer batGroup.changeLock.Unlock()

	batteryDevice, err := newBatteryDeviceWithUPowerDBusObjPath(path)
	if err == nil {
		batGroup.Add(batteryDevice)
	}
}

func (batGroup *batteryGroup) RemoveUPowerBatteryDevice(path dbus.ObjectPath) {
	batGroup.changeLock.Lock()
	defer batGroup.changeLock.Unlock()
	pathStr := "upower://" + string(path)
	batGroup.Remove(pathStr)
}

func (batGroup *batteryGroup) initUpowerBatteryDevices() {
	batGroup.changeLock.Lock()
	defer batGroup.changeLock.Unlock()

	devs, err := batGroup.manager.upower.EnumerateDevices()
	if err != nil {
		logger.Error("Can't EnumerateDevices", err)
		return
	}
	for _, devObjPath := range devs {
		batteryDevice, err := newBatteryDeviceWithUPowerDBusObjPath(devObjPath)
		if err == nil {
			batGroup.Add(batteryDevice)
		}
	}
}

func (batGroup *batteryGroup) initUPowerBackend() {
	batGroup.initUpowerBatteryDevices()
	upower := batGroup.manager.upower
	upower.ConnectDeviceAdded(batGroup.AddUPowerBatteryDevice)
	upower.ConnectDeviceRemoved(batGroup.RemoveUPowerBatteryDevice)
}
