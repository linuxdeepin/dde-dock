package power

import (
	"pkg.deepin.io/dde/api/powersupply/battery"
)

func (v *Manager) setPropOnBattery(value bool) (changed bool) {
	v.PropsMaster.Lock()
	if v.OnBattery != value {
		v.OnBattery = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "OnBattery", value)
	}
	return
}

func (v *Manager) getPropOnBattery() bool {
	v.PropsMaster.RLock()
	value := v.OnBattery
	v.PropsMaster.RUnlock()
	return value
}

func (v *Manager) setPropHasLidSwitch(value bool) (changed bool) {
	v.PropsMaster.Lock()
	if v.HasLidSwitch != value {
		v.HasLidSwitch = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "HasLidSwitch", value)
	}
	return
}

func (v *Manager) getPropHasLidSwitch() bool {
	v.PropsMaster.RLock()
	value := v.HasLidSwitch
	v.PropsMaster.RUnlock()
	return value
}

func (v *Manager) setPropHasBattery(value bool) (changed bool) {
	v.PropsMaster.Lock()
	if v.HasBattery != value {
		v.HasBattery = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "HasBattery", value)
	}
	return
}

func (v *Manager) getPropHasBattery() bool {
	v.PropsMaster.RLock()
	value := v.HasBattery
	v.PropsMaster.RUnlock()
	return value
}

func (v *Manager) setPropBatteryPercentage(value float64) (changed bool) {
	v.PropsMaster.Lock()
	if v.BatteryPercentage != value {
		v.BatteryPercentage = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "BatteryPercentage", value)
	}
	return
}

func (v *Manager) getPropBatteryPercentage() float64 {
	v.PropsMaster.RLock()
	value := v.BatteryPercentage
	v.PropsMaster.RUnlock()
	return value
}

func (v *Manager) setPropBatteryStatus(value battery.Status) (changed bool) {
	v.PropsMaster.Lock()
	if v.BatteryStatus != value {
		v.BatteryStatus = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "BatteryStatus", value)
	}
	return
}

func (v *Manager) getPropBatteryStatus() battery.Status {
	v.PropsMaster.RLock()
	value := v.BatteryStatus
	v.PropsMaster.RUnlock()
	return value
}

func (v *Manager) setPropBatteryTimeToEmpty(value uint64) (changed bool) {
	v.PropsMaster.Lock()
	if v.BatteryTimeToEmpty != value {
		v.BatteryTimeToEmpty = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "BatteryTimeToEmpty", value)
	}
	return
}

func (v *Manager) getPropBatteryTimeToEmpty() uint64 {
	v.PropsMaster.RLock()
	value := v.BatteryTimeToEmpty
	v.PropsMaster.RUnlock()
	return value
}

func (v *Manager) setPropBatteryTimeToFull(value uint64) (changed bool) {
	v.PropsMaster.Lock()
	if v.BatteryTimeToFull != value {
		v.BatteryTimeToFull = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "BatteryTimeToFull", value)
	}
	return
}

func (v *Manager) getPropBatteryTimeToFull() uint64 {
	v.PropsMaster.RLock()
	value := v.BatteryTimeToFull
	v.PropsMaster.RUnlock()
	return value
}

func (v *Battery) setPropSysfsPath(value string) (changed bool) {
	v.PropsMaster.Lock()
	if v.SysfsPath != value {
		v.SysfsPath = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "SysfsPath", value)
	}
	return
}

func (v *Battery) getPropSysfsPath() string {
	v.PropsMaster.RLock()
	value := v.SysfsPath
	v.PropsMaster.RUnlock()
	return value
}

func (v *Battery) setPropIsPresent(value bool) (changed bool) {
	v.PropsMaster.Lock()
	if v.IsPresent != value {
		v.IsPresent = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "IsPresent", value)
	}
	return
}

func (v *Battery) getPropIsPresent() bool {
	v.PropsMaster.RLock()
	value := v.IsPresent
	v.PropsMaster.RUnlock()
	return value
}

func (v *Battery) setPropManufacturer(value string) (changed bool) {
	v.PropsMaster.Lock()
	if v.Manufacturer != value {
		v.Manufacturer = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "Manufacturer", value)
	}
	return
}

func (v *Battery) getPropManufacturer() string {
	v.PropsMaster.RLock()
	value := v.Manufacturer
	v.PropsMaster.RUnlock()
	return value
}

func (v *Battery) setPropModelName(value string) (changed bool) {
	v.PropsMaster.Lock()
	if v.ModelName != value {
		v.ModelName = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "ModelName", value)
	}
	return
}

func (v *Battery) getPropModelName() string {
	v.PropsMaster.RLock()
	value := v.ModelName
	v.PropsMaster.RUnlock()
	return value
}

func (v *Battery) setPropSerialNumber(value string) (changed bool) {
	v.PropsMaster.Lock()
	if v.SerialNumber != value {
		v.SerialNumber = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "SerialNumber", value)
	}
	return
}

func (v *Battery) getPropSerialNumber() string {
	v.PropsMaster.RLock()
	value := v.SerialNumber
	v.PropsMaster.RUnlock()
	return value
}

func (v *Battery) setPropName(value string) (changed bool) {
	v.PropsMaster.Lock()
	if v.Name != value {
		v.Name = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "Name", value)
	}
	return
}

func (v *Battery) getPropName() string {
	v.PropsMaster.RLock()
	value := v.Name
	v.PropsMaster.RUnlock()
	return value
}

func (v *Battery) setPropTechnology(value string) (changed bool) {
	v.PropsMaster.Lock()
	if v.Technology != value {
		v.Technology = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "Technology", value)
	}
	return
}

func (v *Battery) getPropTechnology() string {
	v.PropsMaster.RLock()
	value := v.Technology
	v.PropsMaster.RUnlock()
	return value
}

func (v *Battery) setPropEnergy(value float64) (changed bool) {
	v.PropsMaster.Lock()
	if v.Energy != value {
		v.Energy = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "Energy", value)
	}
	return
}

func (v *Battery) getPropEnergy() float64 {
	v.PropsMaster.RLock()
	value := v.Energy
	v.PropsMaster.RUnlock()
	return value
}

func (v *Battery) setPropEnergyFull(value float64) (changed bool) {
	v.PropsMaster.Lock()
	if v.EnergyFull != value {
		v.EnergyFull = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "EnergyFull", value)
	}
	return
}

func (v *Battery) getPropEnergyFull() float64 {
	v.PropsMaster.RLock()
	value := v.EnergyFull
	v.PropsMaster.RUnlock()
	return value
}

func (v *Battery) setPropEnergyFullDesign(value float64) (changed bool) {
	v.PropsMaster.Lock()
	if v.EnergyFullDesign != value {
		v.EnergyFullDesign = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "EnergyFullDesign", value)
	}
	return
}

func (v *Battery) getPropEnergyFullDesign() float64 {
	v.PropsMaster.RLock()
	value := v.EnergyFullDesign
	v.PropsMaster.RUnlock()
	return value
}

func (v *Battery) setPropEnergyRate(value float64) (changed bool) {
	v.PropsMaster.Lock()
	if v.EnergyRate != value {
		v.EnergyRate = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "EnergyRate", value)
	}
	return
}

func (v *Battery) getPropEnergyRate() float64 {
	v.PropsMaster.RLock()
	value := v.EnergyRate
	v.PropsMaster.RUnlock()
	return value
}

func (v *Battery) setPropVoltage(value float64) (changed bool) {
	v.PropsMaster.Lock()
	if v.Voltage != value {
		v.Voltage = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "Voltage", value)
	}
	return
}

func (v *Battery) getPropVoltage() float64 {
	v.PropsMaster.RLock()
	value := v.Voltage
	v.PropsMaster.RUnlock()
	return value
}

func (v *Battery) setPropPercentage(value float64) (changed bool) {
	v.PropsMaster.Lock()
	if v.Percentage != value {
		v.Percentage = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "Percentage", value)
	}
	return
}

func (v *Battery) getPropPercentage() float64 {
	v.PropsMaster.RLock()
	value := v.Percentage
	v.PropsMaster.RUnlock()
	return value
}

func (v *Battery) setPropCapacity(value float64) (changed bool) {
	v.PropsMaster.Lock()
	if v.Capacity != value {
		v.Capacity = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "Capacity", value)
	}
	return
}

func (v *Battery) getPropCapacity() float64 {
	v.PropsMaster.RLock()
	value := v.Capacity
	v.PropsMaster.RUnlock()
	return value
}

func (v *Battery) setPropStatus(value battery.Status) (changed bool) {
	v.PropsMaster.Lock()
	if v.Status != value {
		v.Status = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "Status", value)
	}
	return
}

func (v *Battery) getPropStatus() battery.Status {
	v.PropsMaster.RLock()
	value := v.Status
	v.PropsMaster.RUnlock()
	return value
}

func (v *Battery) setPropTimeToEmpty(value uint64) (changed bool) {
	v.PropsMaster.Lock()
	if v.TimeToEmpty != value {
		v.TimeToEmpty = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "TimeToEmpty", value)
	}
	return
}

func (v *Battery) getPropTimeToEmpty() uint64 {
	v.PropsMaster.RLock()
	value := v.TimeToEmpty
	v.PropsMaster.RUnlock()
	return value
}

func (v *Battery) setPropTimeToFull(value uint64) (changed bool) {
	v.PropsMaster.Lock()
	if v.TimeToFull != value {
		v.TimeToFull = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "TimeToFull", value)
	}
	return
}

func (v *Battery) getPropTimeToFull() uint64 {
	v.PropsMaster.RLock()
	value := v.TimeToFull
	v.PropsMaster.RUnlock()
	return value
}

func (v *Battery) setPropUpdateTime(value int64) (changed bool) {
	v.PropsMaster.Lock()
	if v.UpdateTime != value {
		v.UpdateTime = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "UpdateTime", value)
	}
	return
}

func (v *Battery) getPropUpdateTime() int64 {
	v.PropsMaster.RLock()
	value := v.UpdateTime
	v.PropsMaster.RUnlock()
	return value
}
