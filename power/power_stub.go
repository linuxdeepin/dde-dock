package main

import "dlib/dbus"

func (p *Power) setPropLidIsPresent(v bool) {
	if p.LidIsPresent != v {
		p.LidIsPresent = v
		dbus.NotifyChange(p, "LidIsPresent")
	}
}

func (p *Power) setPropIdleDelay(v int32) {
	if p.IdleDelay != v {
		p.IdleDelay = v
		dbus.NotifyChange(p, "IdleDelay")
	}
}
func (p *Power) setPropSuspendDelay(v int32) {
	if p.SuspendDelay != v {
		p.SuspendDelay = v
		dbus.NotifyChange(p, "SuspendDelay")
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
	case "CurrentPlan":
		if v, ok := oldv.(int32); ok {
			p.setPlan(v)
		}
	}
}
