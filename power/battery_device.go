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
)

type batteryDevice libupower.Device

func newBatteryDevice(dev *libupower.Device, handler func()) *batteryDevice {
	battery := (*batteryDevice)(dev)
	if handler != nil {
		battery.listenProperties(handler)
	}
	return battery
}

func (battery *batteryDevice) listenProperties(handler func()) {
	battery.Percentage.ConnectChanged(handler)
	battery.State.ConnectChanged(handler)
	battery.IsPresent.ConnectChanged(handler)
	battery.TimeToEmpty.ConnectChanged(handler)
	battery.TimeToFull.ConnectChanged(handler)
}

func (battery *batteryDevice) Destroy() {
	libupower.DestroyDevice((*libupower.Device)(battery))
}
