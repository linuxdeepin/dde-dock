package power

type batteryInfo struct {
	IsPresent        bool
	State            batteryStateType
	Percentage       float64
	Energy           float64
	EnergyFull       float64
	EnergyFullDesign float64
	EnergyEmpty      float64
	EnergyRate       float64
	TimeToFull       int64
	TimeToEmpty      int64
	propertyChange   func(name string, oldVal interface{}, newVal interface{})
}

func newBatteryInfo() *batteryInfo {
	bi := &batteryInfo{}
	return bi
}

func (bi *batteryInfo) OnPropertyChange(handler func(string, interface{}, interface{})) {
	bi.propertyChange = handler
}

func (bi *batteryInfo) setIsPresent(isPresent bool) {
	if bi.IsPresent != isPresent {
		oldVal := bi.IsPresent
		bi.IsPresent = isPresent
		bi.propertyChange("IsPresent", oldVal, isPresent)
	}
}

func (bi *batteryInfo) setState(state batteryStateType) {
	if bi.State != state {
		oldVal := bi.State
		bi.State = state
		bi.propertyChange("State", oldVal, state)
	}
}

func (bi *batteryInfo) setPercentage(p float64) {
	if bi.Percentage != p {
		oldVal := bi.Percentage
		bi.Percentage = p
		bi.propertyChange("Percentage", oldVal, p)
	}
}

func (bi *batteryInfo) setEnergy(energy float64) {
	if bi.Energy != energy {
		oldVal := bi.Energy
		bi.Energy = energy
		bi.propertyChange("Energy", oldVal, energy)
	}
}

func (bi *batteryInfo) setEnergyFull(energy float64) {
	if bi.EnergyFull != energy {
		oldVal := bi.EnergyFull
		bi.EnergyFull = energy
		bi.propertyChange("EnergyFull", oldVal, energy)
	}
}

func (bi *batteryInfo) setEnergyFullDesign(energy float64) {
	if bi.EnergyFullDesign != energy {
		oldVal := bi.EnergyFullDesign
		bi.EnergyFullDesign = energy
		bi.propertyChange("EnergyFullDesign", oldVal, energy)
	}
}

func (bi *batteryInfo) setEnergyEmpty(energy float64) {
	if bi.EnergyEmpty != energy {
		oldVal := bi.EnergyEmpty
		bi.EnergyEmpty = energy
		bi.propertyChange("EnergyEmpty", oldVal, energy)
	}
}

func (bi *batteryInfo) setEnergyRate(rate float64) {
	if bi.EnergyRate != rate {
		oldVal := bi.EnergyRate
		bi.EnergyRate = rate
		bi.propertyChange("EnergyRate", oldVal, rate)
	}
}

func (bi *batteryInfo) setTimeToFull(time int64) {
	if bi.TimeToFull != time {
		oldVal := bi.TimeToFull
		bi.TimeToFull = time
		bi.propertyChange("TimeToFull", oldVal, time)
	}
}

func (bi *batteryInfo) setTimeToEmpty(time int64) {
	if bi.TimeToEmpty != time {
		oldVal := bi.TimeToEmpty
		bi.TimeToEmpty = time
		bi.propertyChange("TimeToEmpty", oldVal, time)
	}
}
