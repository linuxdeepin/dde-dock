package main

import (
	"fmt"
)

const (
	//sync with com.deepin.daemon.power.schema
	PowerPlanCustom          = 0
	PowerPlanPowerSaver      = 1
	PowerPlanBalanced        = 2
	PowerPlanHighPerformance = 3
)

func (p *Power) setBatteryIdleDelay(delay int32) {
	p.setPropBatteryIdleDelay(delay)

	if p.BatteryPlan.Get() == PowerPlanCustom && int32(p.coreSettings.GetInt("battery-idle-delay")) != delay {
		p.coreSettings.SetInt("battery-idle-delay", int(delay))
	}
	p.updateIdletimer()
}

func (p *Power) setBatterySuspendDelay(delay int32) {
	p.setPropBatterySuspendDelay(delay)

	if p.BatteryPlan.Get() == PowerPlanCustom && int32(p.coreSettings.GetInt("battery-suspend-delay")) != delay {
		p.coreSettings.SetInt("battery-suspend-delay", int(delay))
	}
	p.updateIdletimer()
}

func (p *Power) setBatteryPlan(plan int32) {
	switch plan {
	case PowerPlanHighPerformance:
		p.setBatteryIdleDelay(0)
		p.setBatterySuspendDelay(0)
	case PowerPlanBalanced:
		p.setBatteryIdleDelay(600)
		p.setBatterySuspendDelay(0)
	case PowerPlanPowerSaver:
		p.setBatteryIdleDelay(300)
		p.setBatterySuspendDelay(600)
	case PowerPlanCustom:
		p.setBatteryIdleDelay(int32(p.coreSettings.GetInt("battery-idle-delay")))
		p.setBatterySuspendDelay(int32(p.coreSettings.GetInt("battery-suspend-delay")))
	}
}

func (p *Power) setLinePowerIdleDelay(delay int32) {
	p.setPropLinePowerIdleDelay(delay)

	if p.LinePowerPlan.Get() == PowerPlanCustom && int32(p.coreSettings.GetInt("ac-idle-delay")) != delay {
		p.coreSettings.SetInt("ac-idle-delay", int(delay))
	}
	p.updateIdletimer()
}

func (p *Power) setLinePowerSuspendDelay(delay int32) {
	p.setPropLinePowerSuspendDelay(delay)

	if p.LinePowerPlan.Get() == PowerPlanCustom && int32(p.coreSettings.GetInt("ac-suspend-delay")) != delay {
		p.coreSettings.SetInt("ac-suspend-delay", int(delay))
	}
	p.updateIdletimer()
}

func (p *Power) setLinePowerPlan(plan int32) {
	switch plan {
	case PowerPlanHighPerformance:
		p.setLinePowerIdleDelay(0)
		p.setLinePowerSuspendDelay(0)
	case PowerPlanBalanced:
		p.setLinePowerIdleDelay(600)
		p.setLinePowerSuspendDelay(0)
	case PowerPlanPowerSaver:
		p.setLinePowerIdleDelay(300)
		p.setLinePowerSuspendDelay(600)
	case PowerPlanCustom:
		p.setLinePowerIdleDelay(int32(p.coreSettings.GetInt("ac-idle-delay")))
		p.setLinePowerSuspendDelay(int32(p.coreSettings.GetInt("ac-suspend-delay")))
	}
}

func (p *Power) updateIdletimer() {
	var min, delta int32
	var idle, suspend int32
	if p.OnBattery {
		idle = p.BatteryIdleDelay
		suspend = p.BatterySuspendDelay
	} else {
		idle = p.LinePowerIdleDelay
		suspend = p.LinePowerSuspendDelay
	}

	if idle < suspend {
		min = idle
		delta = suspend - idle
	} else {
		min = suspend
		delta = idle - suspend
	}
	if err := p.screensaver.SetTimeout(uint32(min)/10, uint32(delta)/10, true); err != nil {
		LOGGER.Error("Failed set ScreenSaver's timeout:", err)
	} else {
		LOGGER.Info("Set ScreenTimeout to ", uint32(min), uint32(delta))
	}
}

func (p *Power) initPlan() {
	p.screensaver.ConnectIdleOn(p.handleIdleOn)
	p.screensaver.ConnectIdleOff(p.handleIdleOff)
	p.screensaver.ConnectCycleActive(p.handleCycleActive)
	p.updateIdletimer()
}

func (p *Power) handleIdleOn() {
	if p.OnBattery {
		if p.BatteryIdleDelay < p.BatterySuspendDelay {
			fmt.Println("ON battery Idel>>>>>>>>>")
		} else {
			fmt.Println("ON battery Suspend>>>>>>>>>")
			doSuspend()
		}
	} else {
		if p.LinePowerIdleDelay < p.LinePowerSuspendDelay {
			fmt.Println("ON LinePower Idel>>>>>>>>>")
		} else {
			fmt.Println("ON LinePower Suspend>>>>>>>>>")
			doSuspend()
		}
	}
}
func (*Power) handleIdleOff() {
	fmt.Println("OFF>>>>>>>>>")
}
func (p *Power) handleCycleActive() {
	if p.OnBattery {
		fmt.Println("ON battery Suspend>>>>>>>>>")
	} else {
		fmt.Println("ON LinePower Suspend>>>>>>>>>")
	}
	doSuspend()
}
