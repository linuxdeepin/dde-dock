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
	"path/filepath"
)

type upowerBatteryDevice struct {
	udevice *libupower.Device
	info    *batteryInfo
}

func newUpowerBatteryDevice(dev *libupower.Device) *upowerBatteryDevice {
	return &upowerBatteryDevice{
		udevice: dev,
	}
}

func (dev *upowerBatteryDevice) Destroy() {
	if dev.udevice != nil {
		libupower.DestroyDevice(dev.udevice)
	}
}

func (dev *upowerBatteryDevice) GetPath() string {
	return filepath.Base(string(dev.udevice.Path))
}

func (dev *upowerBatteryDevice) GetInfo() *batteryInfo {
	return dev.info
}

func (dev *upowerBatteryDevice) SetInfo(bi *batteryInfo) {
	dev.info = bi
	udevice := dev.udevice

	// init info
	bi.setIsPresent(udevice.IsPresent.Get())
	bi.setState(batteryStateType(udevice.State.Get()))
	bi.setEnergyEmpty(udevice.EnergyEmpty.Get())
	bi.setEnergyFullDesign(udevice.EnergyFullDesign.Get())
	bi.setEnergyFull(udevice.EnergyFull.Get())
	bi.setEnergy(udevice.Energy.Get())
	bi.setTimeToEmpty(udevice.TimeToEmpty.Get())
	bi.setTimeToFull(udevice.TimeToFull.Get())
	bi.setEnergyRate(udevice.EnergyRate.Get())
	bi.setPercentage(udevice.Percentage.Get())
	bi.Inited = true

	// Connect changed
	udevice.IsPresent.ConnectChanged(func() {
		bi.setIsPresent(udevice.IsPresent.Get())
	})
	udevice.State.ConnectChanged(func() {
		bi.setState(batteryStateType(udevice.State.Get()))
	})
	udevice.Energy.ConnectChanged(func() {
		bi.setEnergy(udevice.Energy.Get())
	})
	udevice.EnergyFull.ConnectChanged(func() {
		bi.setEnergyFull(udevice.EnergyFull.Get())
	})
	udevice.EnergyFullDesign.ConnectChanged(func() {
		bi.setEnergyFullDesign(udevice.EnergyFullDesign.Get())
	})
	udevice.EnergyEmpty.ConnectChanged(func() {
		bi.setEnergyEmpty(udevice.EnergyEmpty.Get())
	})
	udevice.TimeToEmpty.ConnectChanged(func() {
		bi.setTimeToEmpty(udevice.TimeToEmpty.Get())
	})
	udevice.TimeToFull.ConnectChanged(func() {
		bi.setTimeToFull(udevice.TimeToFull.Get())
	})
	udevice.EnergyRate.ConnectChanged(func() {
		bi.setEnergyRate(udevice.EnergyRate.Get())
	})
	udevice.Percentage.ConnectChanged(func() {
		bi.setPercentage(udevice.Percentage.Get())
	})
}
