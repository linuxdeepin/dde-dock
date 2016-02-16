/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package power

import "pkg.deepin.io/lib/dbus"

const (
	dbusDest = "com.deepin.daemon.Power"
	dbusPath = "/com/deepin/daemon/Power"
	dbusIFC  = "com.deepin.daemon.Power"
)

func (*Power) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}

func (p *Power) setPropLidIsPresent(v bool) {
	if p.LidIsPresent != v {
		p.LidIsPresent = v
		dbus.NotifyChange(p, "LidIsPresent")
	}
}

func (p *Power) setPropBatteryIdleDelay(v int32) {
	if p.BatteryIdleDelay != v {
		p.BatteryIdleDelay = v
		dbus.NotifyChange(p, "BatteryIdleDelay")
	}
}
func (p *Power) setPropBatterySuspendDelay(v int32) {
	if p.BatterySuspendDelay != v {
		p.BatterySuspendDelay = v
		dbus.NotifyChange(p, "BatterySuspendDelay")
	}
}

func (p *Power) setPropLinePowerIdleDelay(v int32) {
	if p.LinePowerIdleDelay != v {
		p.LinePowerIdleDelay = v
		dbus.NotifyChange(p, "LinePowerIdleDelay")
	}
}
func (p *Power) setPropLinePowerSuspendDelay(v int32) {
	if p.LinePowerSuspendDelay != v {
		p.LinePowerSuspendDelay = v
		dbus.NotifyChange(p, "LinePowerSuspendDelay")
	}
}

func (p *Power) setPropOnBattery(v bool) {
	if p.OnBattery != v {
		p.OnBattery = v
		dbus.NotifyChange(p, "OnBattery")
	}
}

func (p *Power) setPropBatteryIsPresent(v bool) {
	if p.BatteryIsPresent != v {
		p.BatteryIsPresent = v
		dbus.NotifyChange(p, "BatteryIsPresent")
	}
}

func (p *Power) setPropBatteryPercentage(v float64) {
	if p.BatteryPercentage != v {
		p.BatteryPercentage = v
		dbus.NotifyChange(p, "BatteryPercentage")
	}
}

func (p *Power) setPropBatteryState(v uint32) {
	if p.BatteryState != v {
		p.BatteryState = v
		dbus.NotifyChange(p, "BatteryState")
	}
}

func (p *Power) setPropPlanInfo(v string) {
	if p.PlanInfo != v {
		p.PlanInfo = v
		dbus.NotifyChange(p, "PlanInfo")
	}
}

func (p *Power) OnPropertiesChanged(key string, oldv interface{}) {
	switch key {
	case "BatterySuspendDelay":
		v, ok := oldv.(int32)
		if ok && p.BatterySuspendDelay == v {
			return
		}
		p.setBatterySuspendDelay(p.BatterySuspendDelay)
	case "BatteryIdleDelay":
		v, ok := oldv.(int32)
		if ok && p.BatteryIdleDelay == v {
			return
		}
		p.setBatteryIdleDelay(p.BatteryIdleDelay)
	case "LinePowerSuspendDelay":
		v, ok := oldv.(int32)
		logger.Info("[Power] changed:", key, p.LinePowerSuspendDelay, v)
		if ok && p.LinePowerSuspendDelay == v {
			return
		}
		p.setLinePowerSuspendDelay(p.LinePowerSuspendDelay)
	case "LinePowerIdleDelay":
		v, ok := oldv.(int32)
		logger.Info("[Power] changed:", key, p.LinePowerIdleDelay, v)
		if ok && p.LinePowerIdleDelay == v {
			return
		}
		p.setLinePowerIdleDelay(p.LinePowerIdleDelay)
	}
}
