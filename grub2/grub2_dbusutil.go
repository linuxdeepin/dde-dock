package grub2

func (v *Grub2) setPropDefaultEntry(value string) (changed bool) {
	v.PropsMaster.Lock()
	if v.DefaultEntry != value {
		v.DefaultEntry = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "DefaultEntry", value)
	}
	return
}

func (v *Grub2) getPropDefaultEntry() string {
	v.PropsMaster.RLock()
	value := v.DefaultEntry
	v.PropsMaster.RUnlock()
	return value
}

func (v *Grub2) setPropEnableTheme(value bool) (changed bool) {
	v.PropsMaster.Lock()
	if v.EnableTheme != value {
		v.EnableTheme = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "EnableTheme", value)
	}
	return
}

func (v *Grub2) getPropEnableTheme() bool {
	v.PropsMaster.RLock()
	value := v.EnableTheme
	v.PropsMaster.RUnlock()
	return value
}

func (v *Grub2) setPropResolution(value string) (changed bool) {
	v.PropsMaster.Lock()
	if v.Resolution != value {
		v.Resolution = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "Resolution", value)
	}
	return
}

func (v *Grub2) getPropResolution() string {
	v.PropsMaster.RLock()
	value := v.Resolution
	v.PropsMaster.RUnlock()
	return value
}

func (v *Grub2) setPropTimeout(value uint32) (changed bool) {
	v.PropsMaster.Lock()
	if v.Timeout != value {
		v.Timeout = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "Timeout", value)
	}
	return
}

func (v *Grub2) getPropTimeout() uint32 {
	v.PropsMaster.RLock()
	value := v.Timeout
	v.PropsMaster.RUnlock()
	return value
}

func (v *Grub2) setPropUpdating(value bool) (changed bool) {
	v.PropsMaster.Lock()
	if v.Updating != value {
		v.Updating = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "Updating", value)
	}
	return
}

func (v *Grub2) getPropUpdating() bool {
	v.PropsMaster.RLock()
	value := v.Updating
	v.PropsMaster.RUnlock()
	return value
}

func (v *Theme) setPropUpdating(value bool) (changed bool) {
	v.PropsMaster.Lock()
	if v.Updating != value {
		v.Updating = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "Updating", value)
	}
	return
}

func (v *Theme) getPropUpdating() bool {
	v.PropsMaster.RLock()
	value := v.Updating
	v.PropsMaster.RUnlock()
	return value
}

func (v *Theme) setPropItemColor(value string) (changed bool) {
	v.PropsMaster.Lock()
	if v.ItemColor != value {
		v.ItemColor = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "ItemColor", value)
	}
	return
}

func (v *Theme) getPropItemColor() string {
	v.PropsMaster.RLock()
	value := v.ItemColor
	v.PropsMaster.RUnlock()
	return value
}

func (v *Theme) setPropSelectedItemColor(value string) (changed bool) {
	v.PropsMaster.Lock()
	if v.SelectedItemColor != value {
		v.SelectedItemColor = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "SelectedItemColor", value)
	}
	return
}

func (v *Theme) getPropSelectedItemColor() string {
	v.PropsMaster.RLock()
	value := v.SelectedItemColor
	v.PropsMaster.RUnlock()
	return value
}
