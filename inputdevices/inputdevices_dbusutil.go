package inputdevices

func (v *Mouse) setPropDeviceList(value string) (changed bool) {
	if v.DeviceList != value {
		v.DeviceList = value
		v.emitPropChangedDeviceList(value)
		return true
	}
	return false
}

func (v *Mouse) emitPropChangedDeviceList(value string) error {
	return v.service.EmitPropertyChanged(v, "DeviceList", value)
}

func (v *Mouse) setPropExist(value bool) (changed bool) {
	if v.Exist != value {
		v.Exist = value
		v.emitPropChangedExist(value)
		return true
	}
	return false
}

func (v *Mouse) emitPropChangedExist(value bool) error {
	return v.service.EmitPropertyChanged(v, "Exist", value)
}

func (v *Touchpad) setPropExist(value bool) (changed bool) {
	if v.Exist != value {
		v.Exist = value
		v.emitPropChangedExist(value)
		return true
	}
	return false
}

func (v *Touchpad) emitPropChangedExist(value bool) error {
	return v.service.EmitPropertyChanged(v, "Exist", value)
}

func (v *Touchpad) setPropDeviceList(value string) (changed bool) {
	if v.DeviceList != value {
		v.DeviceList = value
		v.emitPropChangedDeviceList(value)
		return true
	}
	return false
}

func (v *Touchpad) emitPropChangedDeviceList(value string) error {
	return v.service.EmitPropertyChanged(v, "DeviceList", value)
}

func (v *TrackPoint) setPropDeviceList(value string) (changed bool) {
	if v.DeviceList != value {
		v.DeviceList = value
		v.emitPropChangedDeviceList(value)
		return true
	}
	return false
}

func (v *TrackPoint) emitPropChangedDeviceList(value string) error {
	return v.service.EmitPropertyChanged(v, "DeviceList", value)
}

func (v *TrackPoint) setPropExist(value bool) (changed bool) {
	if v.Exist != value {
		v.Exist = value
		v.emitPropChangedExist(value)
		return true
	}
	return false
}

func (v *TrackPoint) emitPropChangedExist(value bool) error {
	return v.service.EmitPropertyChanged(v, "Exist", value)
}

func (v *Wacom) setPropDeviceList(value string) (changed bool) {
	if v.DeviceList != value {
		v.DeviceList = value
		v.emitPropChangedDeviceList(value)
		return true
	}
	return false
}

func (v *Wacom) emitPropChangedDeviceList(value string) error {
	return v.service.EmitPropertyChanged(v, "DeviceList", value)
}

func (v *Wacom) setPropExist(value bool) (changed bool) {
	if v.Exist != value {
		v.Exist = value
		v.emitPropChangedExist(value)
		return true
	}
	return false
}

func (v *Wacom) emitPropChangedExist(value bool) error {
	return v.service.EmitPropertyChanged(v, "Exist", value)
}

func (v *Wacom) setPropMapOutput(value string) (changed bool) {
	if v.MapOutput != value {
		v.MapOutput = value
		v.emitPropChangedMapOutput(value)
		return true
	}
	return false
}

func (v *Wacom) emitPropChangedMapOutput(value string) error {
	return v.service.EmitPropertyChanged(v, "MapOutput", value)
}
