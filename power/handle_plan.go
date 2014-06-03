package power

import (
	"dbus/com/deepin/daemon/display"
	"fmt"
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/dpms"
	"time"
)

const (
	HighPerformanceIdleTime    = 15 * 60
	HighPerformanceSuspendTime = 0
	BlancedIdleTime            = 10 * 60
	BlancedSuspendTime         = 0
	PowerSaverIdleTime         = 5 * 60
	PowerSaverSuspendTime      = 15 * 60
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
	p.updatePlanInfo()
}

func (p *Power) setBatterySuspendDelay(delay int32) {
	p.setPropBatterySuspendDelay(delay)

	if p.BatteryPlan.Get() == PowerPlanCustom && int32(p.coreSettings.GetInt("battery-suspend-delay")) != delay {
		p.coreSettings.SetInt("battery-suspend-delay", int(delay))
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
	p.updatePlanInfo()
}

func (p *Power) setLinePowerSuspendDelay(delay int32) {
	p.setPropLinePowerSuspendDelay(delay)

	if p.LinePowerPlan.Get() == PowerPlanCustom && int32(p.coreSettings.GetInt("ac-suspend-delay")) != delay {
		p.coreSettings.SetInt("ac-suspend-delay", int(delay))
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
		p.setLinePowerSuspendDelay(PowerSaverIdleTime)
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
		LOGGER.Error("Failed set ScreenSaver's timeout:", err)
	} else {
		LOGGER.Info("Set ScreenTimeout to ", uint32(min), uint32(suspendDelta))
	}
}

func (p *Power) updatePlanInfo() {
	info := fmt.Sprintf(`{
		"PowerLine":{"Custom":[%d,%d], "PowerSaver":[%d,%d], "Balanced":[%d,%d],"HighPerformance":[%d,%d]},
		"Battery":{"Custom":[%d,%d], "PowerSaver":[%d,%d], "Balanced":[%d,%d],"HighPerformance":[%d,%d]}
	}`, p.LinePowerIdleDelay, p.LinePowerSuspendDelay, PowerSaverIdleTime, PowerSaverSuspendTime,
		BlancedIdleTime, BlancedSuspendTime, HighPerformanceIdleTime, HighPerformanceSuspendTime,
		p.BatteryIdleDelay, p.BatterySuspendDelay, PowerSaverIdleTime, PowerSaverSuspendTime,
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
	defer display.DestroyDisplay(dp)

	stoper := make(chan bool)
	stopAnimation = append(stopAnimation, stoper)
	for _, p := range dp.Monitors.Get() {
		go func() {
			m, _ := display.NewMonitor("com.deepin.daemon.Display", p)
			defer display.DestroyMonitor(m)

			for v := 0.8; v > 0.1; v -= 0.05 {
				<-time.After(time.Millisecond * time.Duration(float64(400)*(v)))

				select {
				case <-stoper:
					for _, name := range m.Outputs.Get() {
						dp.ResetBrightness(name)
					}
					dpmsOn()
					return

				default:
					for _, name := range m.Outputs.Get() {
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
		if p.BatteryIdleDelay < p.BatterySuspendDelay || p.BatterySuspendDelay == 0 {
			doIdleAction()
		} else {
			doSuspend()
		}
	} else {
		if p.LinePowerIdleDelay < p.LinePowerSuspendDelay || p.LinePowerSuspendDelay == 0 {
			doIdleAction()
		} else {
			doSuspend()
		}
	}
}

func (*Power) handleIdleOff() {
	for _, c := range stopAnimation {
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
