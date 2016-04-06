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
	"fmt"
	"pkg.deepin.io/lib/dbus"
	"sync"
)

type batteryGroup struct {
	manager            *Manager
	batterys           []*batteryDevice
	changeLock         sync.Mutex
	destroyLock        sync.Mutex
	propChangedHandler func()
	getInfoMethod      func(*batteryDevice) (*batteryInfo, error)
	InfoMap            map[string]*batteryInfo
}

func NewBatteryGroup(m *Manager) (*batteryGroup, error) {
	batGroup := &batteryGroup{
		manager:            m,
		getInfoMethod:      getBatteryInfoFromUPower,
		propChangedHandler: m.updateBatteryGroupInfo,
	}
	batGroup.setDeviceList()

	upower := m.upower
	upower.ConnectDeviceAdded(func(path dbus.ObjectPath) {
		if batGroup != nil {
			batGroup.AddBatteryDevice(path)
		}
		batGroup.propChangedHandler()
	})

	upower.ConnectDeviceRemoved(func(path dbus.ObjectPath) {
		if batGroup != nil {
			batGroup.RemoveBatteryDevice(path)
		}
		batGroup.propChangedHandler()
	})
	return batGroup, nil
}

func (batGroup *batteryGroup) NewBatteryDevice(dev *libupower.Device) *batteryDevice {
	return newBatteryDevice(dev, batGroup.propChangedHandler)
}

func (batGroup *batteryGroup) AddBatteryDevice(path dbus.ObjectPath) {
	batGroup.changeLock.Lock()
	defer batGroup.changeLock.Unlock()
	if batGroup.isBatteryPathExist(path) {
		return
	}

	dev, err := libupower.NewDevice(upowerDBusDest, dbus.ObjectPath(path))
	if err != nil {
		logger.Warning("New Device Failed:", err)
		return
	}

	if dev.Type.Get() != DeviceTypeBattery {
		return
	}

	batGroup.batterys = append(batGroup.batterys, batGroup.NewBatteryDevice(dev))
}

func (batGroup *batteryGroup) RemoveBatteryDevice(path dbus.ObjectPath) {
	batGroup.changeLock.Lock()
	defer batGroup.changeLock.Unlock()

	var tmpList []*batteryDevice
	for _, battery := range batGroup.batterys {
		if battery.Path == path {
			battery.Destroy()
			continue
		}

		tmpList = append(tmpList, battery)
	}

	batGroup.batterys = tmpList
}

func (batGroup *batteryGroup) UpdateInfo() error {
	if len(batGroup.batterys) == 0 {
		return fmt.Errorf("No battery device")
	}

	if batGroup.getInfoMethod == nil {
		return fmt.Errorf("batteryGroup.getInfoMethod is nil")
	}

	batGroup.InfoMap = make(map[string]*batteryInfo, len(batGroup.batterys))
	for _, battery := range batGroup.batterys {
		key := battery.NativePath.Get()
		batInfo, err := batGroup.getInfoMethod(battery)
		if err != nil {
			return fmt.Errorf("getInfo failed: %v", err)
		}
		logger.Debugf("%v %#v", key, batInfo)
		batGroup.InfoMap[key] = batInfo
	}
	return nil
}

func (batGroup *batteryGroup) Destroy() {
	batGroup.destroyLock.Lock()
	defer batGroup.destroyLock.Unlock()

	for _, battery := range batGroup.batterys {
		battery.Destroy()
	}
	batGroup.batterys = nil
	batGroup.propChangedHandler = nil
}

func (batGroup *batteryGroup) setDeviceList() {
	batGroup.changeLock.Lock()
	defer batGroup.changeLock.Unlock()

	devs, err := batGroup.manager.upower.EnumerateDevices()
	if err != nil {
		logger.Error("Can't EnumerateDevices", err)
		return
	}

	var tmpList []*batteryDevice
	for _, path := range devs {
		logger.Debug("Check device whether is battery:", path)
		dev, err := libupower.NewDevice(upowerDBusDest, path)
		if err != nil {
			logger.Warning("New Device Object Failed:", err)
			continue
		}

		if dev.Type.Get() != DeviceTypeBattery {
			logger.Debugf("Not battery device: %v, type %v", path, dev.Type.Get())
			continue
		}

		logger.Debug("Add battery:", path)
		tmpList = append(tmpList, batGroup.NewBatteryDevice(dev))
	}

	batGroup.batterys = tmpList
	logger.Debug("Battery device list:", batGroup.batterys)
}

func (batGroup *batteryGroup) isBatteryPathExist(path dbus.ObjectPath) bool {
	for _, battery := range batGroup.batterys {
		if battery.Path == path {
			return true
		}
	}
	return false
}

type batteryInfo struct {
	IsPresent        bool
	State            uint32
	Percentage       float64
	Energy           float64
	EnergyFull       float64
	EnergyFullDesign float64
	EnergyEmpty      float64
	TimeToFull       int64
	TimeToEmpty      int64
}

func getBatteryInfoFromUPower(b *batteryDevice) (*batteryInfo, error) {
	return &batteryInfo{
		IsPresent:        b.IsPresent.Get(),
		State:            b.State.Get(),
		Percentage:       b.Percentage.Get(),
		Energy:           b.Energy.Get(),
		EnergyFull:       b.EnergyFull.Get(),
		EnergyFullDesign: b.EnergyFullDesign.Get(),
		EnergyEmpty:      b.EnergyEmpty.Get(),
		TimeToFull:       b.TimeToFull.Get(),
		TimeToEmpty:      b.TimeToEmpty.Get(),
	}, nil
}
