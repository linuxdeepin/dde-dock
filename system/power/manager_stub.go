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
	"gir/gudev-1.0"
	"pkg.deepin.io/dde/api/powersupply/battery"
	"pkg.deepin.io/lib/dbus"
)

const (
	dbusDest = "com.deepin.system.Power"
	dbusPath = "/com/deepin/system/Power"
	dbusIFC  = dbusDest
)

func (*Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}

func (m *Manager) GetBatteries() []*Battery {
	ret := make([]*Battery, 0, len(m.batteries))
	for _, bat := range m.batteries {
		ret = append(ret, bat)
	}
	return ret
}

func (m *Manager) RefreshBatteries() {
	logger.Debug("RefreshBatteries")
	for _, bat := range m.batteries {
		bat.Refresh()
	}
}

func (m *Manager) setPropOnBattery(val bool) {
	logger.Debug("set OnBattery", val)
	if m.OnBattery != val {
		m.OnBattery = val
		dbus.NotifyChange(m, "OnBattery")
	}
}

func (m *Manager) setPropHasBattery(val bool) {
	if m.HasBattery != val {
		m.HasBattery = val
		dbus.NotifyChange(m, "HasBattery")
	}
}

func (m *Manager) setPropBatteryStatus(val battery.Status) {
	if m.BatteryStatus != val {
		m.BatteryStatus = val
		dbus.NotifyChange(m, "BatteryStatus")
	}
}

func (m *Manager) setPropBatteryPercentage(val float64) {
	if m.BatteryPercentage != val {
		m.BatteryPercentage = val
		dbus.NotifyChange(m, "BatteryPercentage")
	}
}

func (m *Manager) setPropBatteryTimeToEmpty(val uint64) {
	if m.BatteryTimeToEmpty != val {
		m.BatteryTimeToEmpty = val
		dbus.NotifyChange(m, "BatteryTimeToEmpty")
	}
}

func (m *Manager) setPropBatteryTimeToFull(val uint64) {
	if m.BatteryTimeToFull != val {
		m.BatteryTimeToFull = val
		dbus.NotifyChange(m, "BatteryTimeToFull")
	}
}

func (m *Manager) Debug(cmd string) {
	logger.Debug("Debug", cmd)
	switch cmd {
	case "init-batteries":
		devices := m.gudevClient.QueryBySubsystem("power_supply")
		for _, dev := range devices {
			m.addBattery(dev)
		}
		logger.Debug("initBatteries done")
		for _, dev := range devices {
			dev.Unref()
		}

	case "remove-all-batteries":
		var devices []*gudev.Device
		for _, bat := range m.batteries {
			devices = append(devices, bat.newDevice())
		}

		for _, dev := range devices {
			m.removeBattery(dev)
			dev.Unref()
		}

	case "destroy":
		m.destroy()

	default:
		logger.Warning("Command not support")
	}
}
