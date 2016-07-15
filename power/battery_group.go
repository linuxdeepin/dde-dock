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
	"pkg.deepin.io/lib/dbus"
	"sync"
)

type batteryGroup struct {
	manager    *Manager
	batteryMap map[string]batteryDevice

	changeLock  sync.Mutex
	destroyLock sync.Mutex

	percentage         float64
	timeToEmpty        int64
	state              batteryStateType
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

	manager.setPropBatteryIsPresent(batteryDisplay, true)
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
			batGroup.updateDisplayPercentage()

		case "Energy":
			batGroup.updateDisplayPercentage()
			batGroup.updateDisplayTimeToEmpty()
			manager.checkBatteryPowerLevel(batGroup)

		case "EnergyRate":
			batGroup.updateDisplayPercentage()
			batGroup.updateDisplayTimeToEmpty()
			manager.checkBatteryPowerLevel(batGroup)

		case "State":
			state := newVal.(batteryStateType)
			manager.setPropBatteryState(path, state)
			batGroup.updateDisplayState()
		}
	})

	device.SetInfo(batInfo)
}

func (batGroup *batteryGroup) Remove(path string) {
	if batDevice, ok := batGroup.batteryMap[path]; ok {
		logger.Debug("Remove ", path)
		batDevice.Destroy()
		delete(batGroup.batteryMap, path)
		// update manager properties
		manager := batGroup.manager
		delete(manager.BatteryPercentage, path)
		delete(manager.BatteryIsPresent, path)
		delete(manager.BatteryState, path)
		if batGroup.batteryDevicesCount() == 0 {
			// remove battery Display
			delete(manager.BatteryPercentage, batteryDisplay)
			delete(manager.BatteryIsPresent, batteryDisplay)
			delete(manager.BatteryState, batteryDisplay)
		}
		dbus.NotifyChange(manager, "BatteryPercentage")
		dbus.NotifyChange(manager, "BatteryIsPresent")
		dbus.NotifyChange(manager, "BatteryState")
		batGroup.updateDisplayState()
		batGroup.updateDisplayPercentage()
		batGroup.updateDisplayTimeToEmpty()
		manager.checkBatteryPowerLevel(batGroup)
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
		energyRate += batInfo.EnergyRate
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

func (batGroup *batteryGroup) getState() batteryStateType {
	/* If one battery is charging, then the composite is charging
	* If all batteries are discharging, then the composite is discharging
	* If all batteries are fully charged, then they're all fully charged
	* Everything else is unknown */
	var stateTotal batteryStateType = BatteryStateUnknown
	for path, dev := range batGroup.batteryMap {
		batInfo := dev.GetInfo()
		state := batInfo.State
		logger.Debugf("path %q state %s", path, state)

		if state == BatteryStateCharging {

			stateTotal = BatteryStateCharging

		} else if stateTotal != BatteryStateCharging &&
			state == BatteryStateDischarging {

			stateTotal = BatteryStateDischarging

		} else if stateTotal == BatteryStateUnknown &&
			state == BatteryStateFullyCharged {

			stateTotal = BatteryStateFullyCharged

		}
	}
	return stateTotal
}

func (batGroup *batteryGroup) updateDisplayState() {
	batGroup.state = batGroup.getState()
	logger.Debug("display state", batGroup.state)
	if batGroup.batteryDevicesCount() > 0 {
		batGroup.manager.setPropBatteryState(batteryDisplay, batGroup.state)
	}
}

func (batGroup *batteryGroup) updateDisplayPercentage() {
	batGroup.percentage = batGroup.getPercentage()
	logger.Debug("display percentage:", batGroup.percentage)
	if batGroup.batteryDevicesCount() > 0 {
		batGroup.manager.setPropBatteryPercentage(batteryDisplay, batGroup.percentage)
	}
}

func (batGroup *batteryGroup) updateDisplayTimeToEmpty() {
	batGroup.timeToEmpty = batGroup.getTimeToEmpty()
	logger.Debug("display timeToEmpty:", batGroup.timeToEmpty)
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
	for _, dev := range batGroup.batteryMap {
		batInfo := dev.GetInfo()
		if !batInfo.Inited {
			logger.Debug("Battery info not inited")
			return batteryPowerLevelUnknown
		}
	}

	manager := batGroup.manager
	if usePercentageForPolicy {
		return manager.getBatteryPowerLevelByPercentage(batGroup.percentage)
	} else {
		// use time to empty for policy
		powerLevel := manager.getBatteryPowerLevelByTimeToEmpty(batGroup.timeToEmpty)

		if powerLevel == batteryPowerLevelAbnormal {
			logger.Debug("Try use percentage for policy")
			return manager.getBatteryPowerLevelByPercentage(batGroup.percentage)
		}
		return powerLevel
	}
}
