package main

import "dlib/dbus"

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

func (p *Power) OnPropertiesChanged(key string, oldv interface{}) {
	switch key {
	case "BatterySuspendDelay":
		p.coreSettings.SetInt("battery-suspend-delay", int(p.BatterySuspendDelay))
	case "BatteryIdleDelay":
		p.coreSettings.SetInt("battery-idle-delay", int(p.BatteryIdleDelay))
	case "LinePowerSuspendDelay":
		p.coreSettings.SetInt("ac-suspend-delay", int(p.LinePowerSuspendDelay))
	case "LinePowerIdleDelay":
		p.coreSettings.SetInt("ac-idle-delay", int(p.LinePowerSuspendDelay))
	}
}
