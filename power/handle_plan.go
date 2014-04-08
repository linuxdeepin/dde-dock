package main

const (
	//sync with com.deepin.daemon.power.schema
	PowerPlanCustom          = 0
	PowerPlanPowerSaver      = 1
	PowerPlanBalanced        = 2
	PowerPlanHighPerformance = 3
)

func (p *Power) setIdleDelay(delay int32) {
	p.setPropIdleDelay(delay)

	if p.CurrentPlan == PowerPlanCustom && int32(p.coreSettings.GetInt("idle-delay")) != delay {
		p.coreSettings.SetInt("idle-delay", int(delay))
	}
}

func (p *Power) setSuspendDelay(delay int32) {
	p.setPropSuspendDelay(delay)

	if p.CurrentPlan == PowerPlanCustom && int32(p.coreSettings.GetInt("suspend-delay")) != delay {
		p.coreSettings.SetInt("suspend-delay", int(delay))
	}
}

func (p *Power) setPlan(plan int32) {
	switch plan {
	case PowerPlanHighPerformance:
		p.setIdleDelay(0)
		p.setSuspendDelay(0)
	case PowerPlanBalanced:
		p.setIdleDelay(600)
		p.setSuspendDelay(0)
	case PowerPlanPowerSaver:
		p.setIdleDelay(300)
		p.setSuspendDelay(600)
	case PowerPlanCustom:
		p.setIdleDelay(int32(p.coreSettings.GetInt("idle-delay")))
		p.setSuspendDelay(int32(p.coreSettings.GetInt("suspend-delay")))
	}
}
