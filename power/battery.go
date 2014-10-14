/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package power

import (
	libupower "dbus/org/freedesktop/upower"
	"fmt"
	"pkg.linuxdeepin.com/lib/dbus"
	"sync"
)

type batteryDevice libupower.Device

type batteryGroup struct {
	BatList            []*batteryDevice
	propChangedHandler func()
	changeLock         sync.Mutex
	destroyLock        sync.Mutex
}

func NewBatteryGroup(handler func()) *batteryGroup {
	batGroup := batteryGroup{}

	batGroup.propChangedHandler = handler
	batGroup.setDeviceLsit()

	return &batGroup
}

func (batGroup *batteryGroup) AddBatteryDevice(path string) {
	batGroup.changeLock.Lock()
	defer batGroup.changeLock.Unlock()
	if batGroup.isBatteryPathExist(path) {
		return
	}

	dev, err := libupower.NewDevice(UPOWER_BUS_NAME, dbus.ObjectPath(path))
	if err != nil {
		logger.Warning("New Device Failed:", err)
		return
	}

	if dev.Type.Get() != DeviceTypeBattery {
		return
	}

	batGroup.BatList = append(batGroup.BatList,
		newBatteryDevice(dev, batGroup.propChangedHandler))
}

func (batGroup *batteryGroup) RemoveBatteryDevice(path string) {
	batGroup.changeLock.Lock()
	defer batGroup.changeLock.Unlock()

	var tmpList []*batteryDevice
	for _, battery := range batGroup.BatList {
		if string(battery.Path) == path {
			battery.Destroy()
			continue
		}

		tmpList = append(tmpList, battery)
	}

	batGroup.BatList = tmpList
}

func (batGroup *batteryGroup) SetPropChnagedHandler(handler func()) {
	if handler == nil {
		return
	}

	batGroup.Destroy()
	batGroup.propChangedHandler = handler
	batGroup.setDeviceLsit()
}

func (batGroup *batteryGroup) GetBatteryInfo() (bool, uint32, float64, error) {
	if len(batGroup.BatList) == 0 {
		return false,
			BatteryStateUnknown,
			0,
			fmt.Errorf("No battery device")
	}

	battery := batGroup.BatList[0]
	return battery.IsPresent.Get(),
		battery.State.Get(),
		battery.Percentage.Get(),
		nil
}

func (batGroup *batteryGroup) Destroy() {
	batGroup.destroyLock.Lock()
	defer batGroup.destroyLock.Unlock()

	for _, battery := range batGroup.BatList {
		battery.Destroy()
	}

	batGroup.BatList = nil
	batGroup.propChangedHandler = nil
}

func (batGroup *batteryGroup) setDeviceLsit() {
	batGroup.changeLock.Lock()
	defer batGroup.changeLock.Unlock()

	devs, err := upower.EnumerateDevices()
	if err != nil {
		logger.Error("Can't EnumerateDevices", err)
		return
	}

	var tmpList []*batteryDevice
	for _, path := range devs {
		dev, err := libupower.NewDevice(UPOWER_BUS_NAME, path)
		if err != nil {
			logger.Warning("New Device Object Failed:", err)
			continue
		}

		if dev.Type.Get() != DeviceTypeBattery {
			logger.Debug("Not battery device:", path)
			continue
		}

		tmpList = append(tmpList,
			newBatteryDevice(dev, batGroup.propChangedHandler))
	}

	batGroup.BatList = tmpList
}

func (batGroup *batteryGroup) isBatteryPathExist(path string) bool {
	for _, battery := range batGroup.BatList {
		if string(battery.Path) == path {
			return true
		}
	}

	return false
}

func newBatteryDevice(dev *libupower.Device, handler func()) *batteryDevice {
	battery := (*batteryDevice)(dev)
	if handler != nil {
		battery.listenProperties(handler)
	}

	return battery
}

func (battery *batteryDevice) listenProperties(handler func()) {
	battery.Percentage.ConnectChanged(func() {
		if battery == nil {
			return
		}

		handler()
	})
}

func (battery *batteryDevice) Destroy() {
	libupower.DestroyDevice((*libupower.Device)(battery))
}
