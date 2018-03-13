package timedate

func (v *Manager) setPropCanNTP(value bool) (changed bool) {
	if v.CanNTP != value {
		v.CanNTP = value
		v.emitPropChangedCanNTP(value)
		return true
	}
	return false
}

func (v *Manager) emitPropChangedCanNTP(value bool) error {
	return v.service.EmitPropertyChanged(v, "CanNTP", value)
}

func (v *Manager) setPropNTP(value bool) (changed bool) {
	if v.NTP != value {
		v.NTP = value
		v.emitPropChangedNTP(value)
		return true
	}
	return false
}

func (v *Manager) emitPropChangedNTP(value bool) error {
	return v.service.EmitPropertyChanged(v, "NTP", value)
}

func (v *Manager) setPropLocalRTC(value bool) (changed bool) {
	if v.LocalRTC != value {
		v.LocalRTC = value
		v.emitPropChangedLocalRTC(value)
		return true
	}
	return false
}

func (v *Manager) emitPropChangedLocalRTC(value bool) error {
	return v.service.EmitPropertyChanged(v, "LocalRTC", value)
}

func (v *Manager) setPropTimezone(value string) (changed bool) {
	if v.Timezone != value {
		v.Timezone = value
		v.emitPropChangedTimezone(value)
		return true
	}
	return false
}

func (v *Manager) emitPropChangedTimezone(value string) error {
	return v.service.EmitPropertyChanged(v, "Timezone", value)
}
