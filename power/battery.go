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
	"io/ioutil"
	"pkg.deepin.io/lib/dbus"
	"strconv"
	"strings"
	"sync"
	"time"
)

type batteryDevice libupower.Device

type batteryGroup struct {
	BatList            []*batteryDevice
	propChangedHandler func()
	changeLock         sync.Mutex
	destroyLock        sync.Mutex
	quitListener       chan struct{}
}

var batteryStateMap = map[string]uint32{
	"Unknown":          BatteryStateUnknown,
	"Charging":         BatteryStateCharging,
	"Discharging":      BatteryStateDischarging,
	"Empty":            BatteryStateEmpty,
	"FullCharged":      BatteryStateFullyCharged,
	"PendingCharge":    BatteryStatePendingCharge,
	"PendingDischarge": BatteryStatePendingDischarge,
}

func NewBatteryGroup(handler func()) *batteryGroup {
	batGroup := batteryGroup{}

	batGroup.propChangedHandler = handler
	batGroup.setDeviceList()

	if isSWPlatform() {
		batGroup.quitListener = make(chan struct{})
		go batGroup.stateListener()
	}

	return &batGroup
}

func (batGroup *batteryGroup) AddBatteryDevice(path dbus.ObjectPath) {
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

func (batGroup *batteryGroup) RemoveBatteryDevice(path dbus.ObjectPath) {
	batGroup.changeLock.Lock()
	defer batGroup.changeLock.Unlock()

	var tmpList []*batteryDevice
	for _, battery := range batGroup.BatList {
		if battery.Path == path {
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
	if isSWPlatform() {
		batGroup.quitListener = make(chan struct{})
		go batGroup.stateListener()
	}
	batGroup.setDeviceList()
}

func (batGroup *batteryGroup) GetBatteryInfo() (bool, uint32, float64, error) {
	if len(batGroup.BatList) == 0 {
		return false,
			BatteryStateUnknown,
			0,
			fmt.Errorf("No battery device")
	}

	battery := batGroup.BatList[0]
	if isSWPlatform() {
		file := "/sys/class/power_supply/" + battery.NativePath.Get() + "/uevent"
		logger.Debug("[Battery Info] file:", file)
		return getBatteryStateFromFile(file)
	}

	return battery.IsPresent.Get(),
		battery.State.Get(),
		battery.Percentage.Get(),
		nil
}

func (batGroup *batteryGroup) Destroy() {
	batGroup.destroyLock.Lock()
	defer batGroup.destroyLock.Unlock()

	if isSWPlatform() {
		if batGroup.quitListener != nil {
			close(batGroup.quitListener)
			batGroup.quitListener = nil
		}
	}

	for _, battery := range batGroup.BatList {
		battery.Destroy()
	}

	batGroup.BatList = nil
	batGroup.propChangedHandler = nil
}

func (batGroup *batteryGroup) setDeviceList() {
	batGroup.changeLock.Lock()
	defer batGroup.changeLock.Unlock()

	devs, err := upower.EnumerateDevices()
	if err != nil {
		logger.Error("Can't EnumerateDevices", err)
		return
	}

	var tmpList []*batteryDevice
	for _, path := range devs {
		logger.Debug("Check device whether is battery:", path)
		dev, err := libupower.NewDevice(UPOWER_BUS_NAME, path)
		if err != nil {
			logger.Warning("New Device Object Failed:", err)
			continue
		}

		if dev.Type.Get() != DeviceTypeBattery {
			logger.Debug("Not battery device:", path)
			continue
		}

		logger.Debug("Add battery:", path)
		tmpList = append(tmpList,
			newBatteryDevice(dev, batGroup.propChangedHandler))
	}

	batGroup.BatList = tmpList
	logger.Debug("Battery device list:", batGroup.BatList)
}

func (batGroup *batteryGroup) isBatteryPathExist(path dbus.ObjectPath) bool {
	for _, battery := range batGroup.BatList {
		if battery.Path == path {
			return true
		}
	}

	return false
}

func (batGroup *batteryGroup) stateListener() {
	if batGroup.propChangedHandler == nil {
		return
	}

	for {
		select {
		case <-batGroup.quitListener:
			return
		case <-time.After(time.Second * 5):
			if batGroup.propChangedHandler != nil {
				batGroup.propChangedHandler()
			}
		}
	}
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

func getBatteryStateFromFile(file string) (bool, uint32, float64, error) {
	data, err := readBatteryInfoFile(file)
	if err != nil {
		return false, BatteryStateUnknown, 0, err
	}

	isPresent, ok := data["POWER_SUPPLY_PRESENT"]
	if !ok {
		return false, BatteryStateUnknown, 0,
			fmt.Errorf("Invalid battery state file: %v", file)
	}

	status, ok := data["POWER_SUPPLY_STATUS"]
	if !ok {
		return false, BatteryStateUnknown, 0,
			fmt.Errorf("Invalid battery state file: %v", file)
	}
	state, ok := batteryStateMap[status]
	if !ok {
		return false, BatteryStateUnknown, 0,
			fmt.Errorf("Invalid battery state file: %v", file)
	}

	capacity, ok := data["POWER_SUPPLY_CAPACITY"]
	if !ok {
		return false, BatteryStateUnknown, 0,
			fmt.Errorf("Invalid battery state file: %v", file)
	}
	percentage, err := strconv.ParseFloat(capacity, 10)
	if err != nil {
		return false, BatteryStateUnknown, 0, err
	}

	logger.Debug("[Battery State] from file:", isPresent, state, percentage)
	return (isPresent == "1"), state, percentage, nil
}

func readBatteryInfoFile(file string) (map[string]string, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	var data = make(map[string]string)
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		array := strings.Split(line, "=")
		if len(array) != 2 {
			continue
		}
		data[array[0]] = strings.TrimSpace(array[1])
	}

	return data, nil
}
