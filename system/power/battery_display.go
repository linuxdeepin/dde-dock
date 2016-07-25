/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package power

import (
	"pkg.deepin.io/dde/api/powersupply/battery"
	"pkg.deepin.io/lib/dbus"
	"time"
)

func (m *Manager) refreshBatteryDisplay() {
	logger.Debug("refreshBatteryDisplay")
	defer dbus.Emit(m, "BatteryDisplayUpdate", time.Now().Unix())
	batteryCount := len(m.batteries)
	if batteryCount == 0 {
		m.resetBatteryDisplay()
		return
	}

	var energyNowTotal, energyFullTotal, powerNowTotal uint64
	for _, bat := range m.batteries {
		energyNowTotal += bat.EnergyNow
		energyFullTotal += bat.EnergyFull
		powerNowTotal += bat.PowerNow
	}

	percentage := rightPercentage(
		float64(energyNowTotal) / float64(energyFullTotal) * 100.0)

	status := m.getBatteryDisplayStatus()

	var timeToEmpty, timeToFull uint64

	timeToEmpty = uint64(float64(energyNowTotal) / float64(powerNowTotal) * 3600)
	timeToFull = uint64(float64(energyFullTotal-energyNowTotal) / float64(powerNowTotal) * 3600)

	// report
	m.setPropHasBattery(true)
	m.setPropBatteryPercentage(percentage)
	m.setPropBatteryStatus(status)
	m.setPropBatteryTimeToEmpty(timeToEmpty)
	m.setPropBatteryTimeToFull(timeToFull)

	logger.Debugf("energyNowTotal: %vµAh", energyNowTotal)
	logger.Debugf("energyFullTotal: %vµAh", energyFullTotal)
	logger.Debugf("powerNowTotal: %vµA", powerNowTotal)
	logger.Debugf("percentage: %.1f%%", percentage)
	logger.Debug("status:", status, uint32(status))
	logger.Debugf("timeToEmpty %v (%vs), timeToFull %v (%vs)",
		time.Duration(timeToEmpty)*time.Second,
		timeToEmpty,
		time.Duration(timeToFull)*time.Second,
		timeToFull)
}

func _getBatteryDisplayStatus(batteries []*Battery) battery.Status {
	var statusSlice []battery.Status
	for _, bat := range batteries {
		statusSlice = append(statusSlice, bat.Status)
	}
	return battery.GetDisplayStatus(statusSlice)
}

func (m *Manager) getBatteryDisplayStatus() battery.Status {
	return _getBatteryDisplayStatus(m.GetBatteries())
}

func (m *Manager) resetBatteryDisplay() {
	m.setPropHasBattery(false)
	m.setPropBatteryPercentage(0)
	m.setPropBatteryTimeToFull(0)
	m.setPropBatteryTimeToEmpty(0)
	m.setPropBatteryStatus(battery.StatusUnknown)
}

func rightPercentage(val float64) float64 {
	if val < 0.0 {
		val = 0.0
	} else if val > 100.0 {
		val = 100.0
	}
	return val
}
