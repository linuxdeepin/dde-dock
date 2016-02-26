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
	"dbus/com/deepin/daemon/display"
	"fmt"
	"time"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/dpms"
)

const (
	// 高性能模式下空闲检测超时
	HighPerformanceIdleTime = 15 * 60
	// 高性能模式下挂起超时
	HighPerformanceSuspendTime = 0
	// 均衡模式下空闲检测超时
	BlancedIdleTime = 10 * 60
	// 均衡模式下挂起超时
	BlancedSuspendTime = 0
	// 节能模式下空闲检测超时
	PowerSaverIdleTime = 5 * 60
	// 节能模式下挂起超时
	PowerSaverSuspendTime = 15 * 60
)

const (
	//sync with com.deepin.daemon.power.schema
	// 电源计划：自定义
	PowerPlanCustom = 0
	// 电源计划：节能模式
	PowerPlanPowerSaver = 1
	// 电源计划：均衡模式
	PowerPlanBalanced = 2
	// 电源计划：高性能模式
	PowerPlanHighPerformance = 3
)

func (p *Power) setBatteryIdleDelay(delay int32) {
	p.setPropBatteryIdleDelay(delay)

	if p.BatteryPlan.Get() == PowerPlanCustom && int32(p.coreSettings.GetInt("battery-idle-delay")) != delay {
		p.coreSettings.SetInt("battery-idle-delay", delay)
	}
	p.updateIdletimer()
	p.updatePlanInfo()
}

func (p *Power) setBatterySuspendDelay(delay int32) {
	p.setPropBatterySuspendDelay(delay)

	if p.BatteryPlan.Get() == PowerPlanCustom && int32(p.coreSettings.GetInt("battery-suspend-delay")) != delay {
		p.coreSettings.SetInt("battery-suspend-delay", delay)
	}
	p.updateIdletimer()
	p.updatePlanInfo()
}

func (p *Power) setBatteryPlan(plan int32) {
	switch plan {
	case PowerPlanHighPerformance:
		p.setBatteryIdleDelay(HighPerformanceIdleTime)
		p.setBatterySuspendDelay(HighPerformanceSuspendTime)
	case PowerPlanBalanced:
		p.setBatteryIdleDelay(BlancedIdleTime)
		p.setBatterySuspendDelay(BlancedSuspendTime)
	case PowerPlanPowerSaver:
		p.setBatteryIdleDelay(PowerSaverIdleTime)
		p.setBatterySuspendDelay(PowerSaverSuspendTime)
	case PowerPlanCustom:
		p.setBatteryIdleDelay(p.coreSettings.GetInt("battery-idle-delay"))
		p.setBatterySuspendDelay(p.coreSettings.GetInt("battery-suspend-delay"))
	}
}

func (p *Power) setLinePowerIdleDelay(delay int32) {
	p.setPropLinePowerIdleDelay(delay)

	if p.LinePowerPlan.Get() == PowerPlanCustom && int32(p.coreSettings.GetInt("ac-idle-delay")) != delay {
		p.coreSettings.SetInt("ac-idle-delay", delay)
	}
	p.updateIdletimer()
	p.updatePlanInfo()
}

func (p *Power) setLinePowerSuspendDelay(delay int32) {
	p.setPropLinePowerSuspendDelay(delay)

	if p.LinePowerPlan.Get() == PowerPlanCustom && int32(p.coreSettings.GetInt("ac-suspend-delay")) != delay {
		p.coreSettings.SetInt("ac-suspend-delay", delay)
	}
	p.updateIdletimer()
	p.updatePlanInfo()
}

func (p *Power) setLinePowerPlan(plan int32) {
	switch plan {
	case PowerPlanHighPerformance:
		p.setLinePowerIdleDelay(HighPerformanceIdleTime)
		p.setLinePowerSuspendDelay(HighPerformanceSuspendTime)
	case PowerPlanBalanced:
		p.setLinePowerIdleDelay(BlancedIdleTime)
		p.setLinePowerSuspendDelay(BlancedSuspendTime)
	case PowerPlanPowerSaver:
		p.setLinePowerIdleDelay(PowerSaverIdleTime)
		p.setLinePowerSuspendDelay(PowerSaverSuspendTime)
	case PowerPlanCustom:
		p.setLinePowerIdleDelay(int32(p.coreSettings.GetInt("ac-idle-delay")))
		p.setLinePowerSuspendDelay(int32(p.coreSettings.GetInt("ac-suspend-delay")))
	}
}

var suspendDelta int32 = 0

func (p *Power) updateIdletimer() {
	var min int32
	var idle, suspend int32
	if p.OnBattery {
		idle = p.BatteryIdleDelay
		suspend = p.BatterySuspendDelay
	} else {
		idle = p.LinePowerIdleDelay
		suspend = p.LinePowerSuspendDelay
	}
	if idle == 0 {
		idle = 0xfffffff
	}
	if suspend == 0 {
		suspend = 0xfffffff
	}

	if idle < suspend {
		min = idle
		suspendDelta = suspend - idle
	} else {
		min = suspend
		suspendDelta = idle - suspend
	}
	if suspendDelta > 0xfffff {
		suspendDelta = 0
	}
	if min > 0xffffff {
		min = 0
	}
	if err := p.screensaver.SetTimeout(uint32(min), 0, false); err != nil {
		logger.Error("Failed set ScreenSaver's timeout:", err)
	}
}

func (p *Power) updatePlanInfo() {
	acDelay := p.coreSettings.GetInt("ac-idle-delay")
	acSuspend := p.coreSettings.GetInt("ac-suspend-delay")
	batteryDelay := p.coreSettings.GetInt("battery-idle-delay")
	batterySuspend := p.coreSettings.GetInt("battery-suspend-delay")

	info := fmt.Sprintf(`{
		"PowerLine":{"Custom":[%d,%d], "PowerSaver":[%d,%d], "Balanced":[%d,%d],"HighPerformance":[%d,%d]},
		"Battery":{"Custom":[%d,%d], "PowerSaver":[%d,%d], "Balanced":[%d,%d],"HighPerformance":[%d,%d]}
	}`, acDelay, acSuspend, PowerSaverIdleTime, PowerSaverSuspendTime,
		BlancedIdleTime, BlancedSuspendTime, HighPerformanceIdleTime, HighPerformanceSuspendTime,
		batteryDelay, batterySuspend, PowerSaverIdleTime, PowerSaverSuspendTime,
		BlancedIdleTime, BlancedSuspendTime, HighPerformanceIdleTime, HighPerformanceSuspendTime,
	)
	p.setPropPlanInfo(info)
}

var dpmsOn func()
var dpmsOff func()

func (p *Power) initPlan() {
	p.screensaver.ConnectIdleOn(p.handleIdleOn)
	p.screensaver.ConnectIdleOff(p.handleIdleOff)
	p.updateIdletimer()
	con, _ := xgb.NewConn()
	dpms.Init(con)
	dpmsOn = func() { dpms.ForceLevel(con, dpms.DPMSModeOn) }
	dpmsOff = func() { dpms.ForceLevel(con, dpms.DPMSModeOff) }
}

var stopAnimation []chan bool

func doIdleAction() {
	dp, _ := display.NewDisplay("com.deepin.daemon.Display", "/com/deepin/daemon/Display")

	stoper := make(chan bool)
	stopAnimation = append(stopAnimation, stoper)
	for _, p := range dp.Monitors.Get() {
		m, err := display.NewMonitor("com.deepin.daemon.Display", p)
		if err != nil {
			logger.Warningf("create monitor %v failed:%v", p, err)
			continue
		}

		go func() {
			outputs := m.Outputs.Get()
			for v := 0.8; v > 0.1; v -= 0.05 {
				<-time.After(time.Millisecond * time.Duration(float64(400)*(v)))

				select {
				case <-stoper:
					for _, name := range outputs {
						dp.ResetBrightness(name)
					}
					dpmsOn()
					return

				default:
					for _, name := range outputs {
						dp.ChangeBrightness(name, v)
					}
				}
			}
		}()
	}

	dpmsOff()
	if suspendDelta != 0 {
		for {
			select {
			case <-time.After(time.Second * time.Duration(suspendDelta)):
				doSuspend()
				return
			case <-stoper:
				return
			}
		}
	}
}

func (p *Power) handleIdleOn() {
	if p.OnBattery {
		if p.BatteryIdleDelay == 0 && p.BatterySuspendDelay == 0 {
			return
		}

		if (p.BatteryIdleDelay != 0 && p.BatteryIdleDelay < p.BatterySuspendDelay) ||
			p.BatterySuspendDelay == 0 {
			doIdleAction()
		} else {
			doSuspend()
		}
	} else {
		if p.LinePowerIdleDelay == 0 && p.LinePowerSuspendDelay == 0 {
			return
		}

		if (p.LinePowerIdleDelay != 0 && p.LinePowerIdleDelay < p.LinePowerSuspendDelay) ||
			p.LinePowerSuspendDelay == 0 {
			doIdleAction()
		} else {
			doSuspend()
		}
	}
}

func (*Power) handleIdleOff() {
	for _, c := range stopAnimation {
		if c == nil {
			continue
		}
		close(c)
	}
	stopAnimation = nil

	dpmsOn()
	dp, _ := display.NewDisplay("com.deepin.daemon.Display", "/com/deepin/daemon/Display")
	defer display.DestroyDisplay(dp)
	for _, p := range dp.Monitors.Get() {
		m, _ := display.NewMonitor("com.deepin.daemon.Display", p)
		defer display.DestroyMonitor(m)
		for _, name := range m.Outputs.Get() {
			dp.ResetBrightness(name)
		}
	}
}
