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
	"errors"
	"math"
	"sync"
)

type batteryGroup struct {
	manager    *Manager
	batteryMap map[string]batteryDevice

	changeLock  sync.Mutex
	destroyLock sync.Mutex

	propChangedHandler func()
}

func NewBatteryGroup(m *Manager) (*batteryGroup, error) {
	batGroup := &batteryGroup{
		manager: m,
	}

	batGroup.batteryMap = make(map[string]batteryDevice, 0)

	logger.Debug("powerSupplyDataBackend:", m.powerSupplyDataBackend)
	switch m.powerSupplyDataBackend {
	case powerSupplyDataBackendUPower:
		batGroup.initUPowerBackend()

	case powerSupplyDataBackendPoll:
		batGroup.initPollBackend()
	default:
		return nil, errors.New("unknown data backend")
	}

	return batGroup, nil
}

func (batGroup *batteryGroup) isBatteryDeviceExist(path string) bool {
	_, ok := batGroup.batteryMap[path]
	return ok
}

func (batGroup *batteryGroup) Add(device batteryDevice) {
	manager := batGroup.manager
	path := device.GetPath()

	if batGroup.isBatteryDeviceExist(path) {
		logger.Warning("Add failed, battery device existed")
		return
	}

	batGroup.batteryMap[path] = device
	batInfo := newBatteryInfo()
	batInfo.OnPropertyChange(func(property string, oldVal interface{}, newVal interface{}) {
		logger.Debugf("%q propertyChange %q: %v => %v", path, property, oldVal, newVal)

		switch property {
		case "IsPresent":
			isPresent := newVal.(bool)
			manager.setPropBatteryIsPresent(path, isPresent)

		case "Percentage":
			percentage := newVal.(float64)
			manager.setPropBatteryPercentage(path, percentage)

		case "Energy":
			manager.checkBatteryPowerLevel(batGroup)

		case "EnergyRate":
			manager.checkBatteryPowerLevel(batGroup)

		case "State":
			state := newVal.(batteryStateType)
			manager.setPropBatteryState(path, state)
		}
	})

	device.SetInfo(batInfo)
}

func (batGroup *batteryGroup) Remove(path string) {
	if batDevice, ok := batGroup.batteryMap[path]; ok {
		batDevice.Destroy()
		delete(batGroup.batteryMap, path)
	}
}

func (batGroup *batteryGroup) Destroy() {
	batGroup.destroyLock.Lock()
	defer batGroup.destroyLock.Unlock()

	for _, dev := range batGroup.batteryMap {
		dev.Destroy()
	}
	batGroup.batteryMap = nil
}

func (batGroup *batteryGroup) batteryDevicesCount() int {
	return len(batGroup.batteryMap)
}

func (batGroup *batteryGroup) getTimeToEmpty() int64 {
	var energySum, energyRate float64
	for path, dev := range batGroup.batteryMap {
		batInfo := dev.GetInfo()
		logger.Debugf("path %q, batInfo: %#v", path, batInfo)
		energySum += batInfo.Energy
		// 假设只有正在使用的电池 engeryRate 大于 0
		if batInfo.EnergyRate > 0 {
			energyRate = batInfo.EnergyRate
		}
	}
	hours := energySum / energyRate
	logger.Debugf("getTimeToEmpty: %v/%v = %.1f h", energySum, energyRate, hours)
	return int64(hours * 3600)
}

func (batGroup *batteryGroup) getPercentage() float64 {
	var energySum, energyFullSum float64
	for path, dev := range batGroup.batteryMap {
		batInfo := dev.GetInfo()
		logger.Debugf("path %q, batInfo: %#v", path, batInfo)
		energySum += batInfo.Energy
		energyFullSum += batInfo.EnergyFull
	}
	logger.Debugf("%v/%v", energySum, energyFullSum)
	return math.Floor((energySum / energyFullSum) * 100.0)
}

func (m *Manager) getBatteryPowerLevelByPercentage(percentage float64) uint32 {
	switch {
	case percentage < batteryPercentageAbnormal:
		return batteryPowerLevelAbnormal

	case percentage <= m.batteryPercentageExhausted:
		return batteryPowerLevelExhausted

	case percentage <= m.batteryPercentageVeryLow:
		return batteryPowerLevelVeryLow

	case percentage <= m.batteryPercentageLow:
		return batteryPowerLevelLow

	default:
		return batteryPowerLevelSufficient
	}
}

func (m *Manager) getBatteryPowerLevelByTimeToEmpty(time int64) uint32 {
	switch {
	case time < timeToEmptyAbnormal:
		return batteryPowerLevelAbnormal

	case time <= m.timeToEmptyExhausted:
		return batteryPowerLevelExhausted

	case time <= m.timeToEmptyVeryLow:
		return batteryPowerLevelVeryLow

	case time <= m.timeToEmptyLow:
		return batteryPowerLevelLow

	default:
		return batteryPowerLevelSufficient
	}
}

func (batGroup *batteryGroup) getPowerLevel(usePercentageForPolicy bool) uint32 {
	manager := batGroup.manager
	if usePercentageForPolicy {
		percentage := batGroup.getPercentage()
		logger.Debugf("sum percentage: %.2f%%", percentage)
		return manager.getBatteryPowerLevelByPercentage(percentage)
	} else {
		// use time to empty for policy
		timeToEmpty := batGroup.getTimeToEmpty()
		logger.Debug("sum timeToEmpty(secs):", timeToEmpty)
		powerLevel := manager.getBatteryPowerLevelByTimeToEmpty(timeToEmpty)

		if powerLevel == batteryPowerLevelAbnormal {
			logger.Debug("Try use percentage for policy")
			percentage := batGroup.getPercentage()
			logger.Debugf("sum percentage: %.2f%%", percentage)
			return manager.getBatteryPowerLevelByPercentage(percentage)
		}
		return powerLevel
	}
}
