package grub2

func (v *Grub2) setPropDefaultEntry(value string) (changed bool) {
	if v.DefaultEntry != value {
		v.DefaultEntry = value
		v.emitPropChangedDefaultEntry(value)
		return true
	}
	return false
}

func (v *Grub2) emitPropChangedDefaultEntry(value string) error {
	return v.service.EmitPropertyChanged(v, "DefaultEntry", value)
}

func (v *Grub2) setPropEnableTheme(value bool) (changed bool) {
	if v.EnableTheme != value {
		v.EnableTheme = value
		v.emitPropChangedEnableTheme(value)
		return true
	}
	return false
}

func (v *Grub2) emitPropChangedEnableTheme(value bool) error {
	return v.service.EmitPropertyChanged(v, "EnableTheme", value)
}

func (v *Grub2) setPropResolution(value string) (changed bool) {
	if v.Resolution != value {
		v.Resolution = value
		v.emitPropChangedResolution(value)
		return true
	}
	return false
}

func (v *Grub2) emitPropChangedResolution(value string) error {
	return v.service.EmitPropertyChanged(v, "Resolution", value)
}

func (v *Grub2) setPropTimeout(value uint32) (changed bool) {
	if v.Timeout != value {
		v.Timeout = value
		v.emitPropChangedTimeout(value)
		return true
	}
	return false
}

func (v *Grub2) emitPropChangedTimeout(value uint32) error {
	return v.service.EmitPropertyChanged(v, "Timeout", value)
}

func (v *Grub2) setPropUpdating(value bool) (changed bool) {
	if v.Updating != value {
		v.Updating = value
		v.emitPropChangedUpdating(value)
		return true
	}
	return false
}

func (v *Grub2) emitPropChangedUpdating(value bool) error {
	return v.service.EmitPropertyChanged(v, "Updating", value)
}

func (v *Theme) setPropUpdating(value bool) (changed bool) {
	if v.Updating != value {
		v.Updating = value
		v.emitPropChangedUpdating(value)
		return true
	}
	return false
}

func (v *Theme) emitPropChangedUpdating(value bool) error {
	return v.service.EmitPropertyChanged(v, "Updating", value)
}

func (v *Theme) setPropItemColor(value string) (changed bool) {
	if v.ItemColor != value {
		v.ItemColor = value
		v.emitPropChangedItemColor(value)
		return true
	}
	return false
}

func (v *Theme) emitPropChangedItemColor(value string) error {
	return v.service.EmitPropertyChanged(v, "ItemColor", value)
}

func (v *Theme) setPropSelectedItemColor(value string) (changed bool) {
	if v.SelectedItemColor != value {
		v.SelectedItemColor = value
		v.emitPropChangedSelectedItemColor(value)
		return true
	}
	return false
}

func (v *Theme) emitPropChangedSelectedItemColor(value string) error {
	return v.service.EmitPropertyChanged(v, "SelectedItemColor", value)
}
