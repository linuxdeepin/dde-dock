package grub2

import (
	"pkg.deepin.io/lib/dbusutil"
)

func (v *Grub2) setPropDefaultEntry(service *dbusutil.Service, value string) (changed bool) {
	v.PropsMu.Lock()
	if v.DefaultEntry != value {
		v.DefaultEntry = value
		changed = true
	}
	v.PropsMu.Unlock()
	if service != nil && changed {
		service.EmitPropertyChanged(v, "DefaultEntry", value)
	}
	return
}

func (v *Grub2) getPropDefaultEntry() string {
	v.PropsMu.RLock()
	value := v.DefaultEntry
	v.PropsMu.RUnlock()
	return value
}

func (v *Grub2) setPropEnableTheme(service *dbusutil.Service, value bool) (changed bool) {
	v.PropsMu.Lock()
	if v.EnableTheme != value {
		v.EnableTheme = value
		changed = true
	}
	v.PropsMu.Unlock()
	if service != nil && changed {
		service.EmitPropertyChanged(v, "EnableTheme", value)
	}
	return
}

func (v *Grub2) getPropEnableTheme() bool {
	v.PropsMu.RLock()
	value := v.EnableTheme
	v.PropsMu.RUnlock()
	return value
}

func (v *Grub2) setPropResolution(service *dbusutil.Service, value string) (changed bool) {
	v.PropsMu.Lock()
	if v.Resolution != value {
		v.Resolution = value
		changed = true
	}
	v.PropsMu.Unlock()
	if service != nil && changed {
		service.EmitPropertyChanged(v, "Resolution", value)
	}
	return
}

func (v *Grub2) getPropResolution() string {
	v.PropsMu.RLock()
	value := v.Resolution
	v.PropsMu.RUnlock()
	return value
}

func (v *Grub2) setPropTimeout(service *dbusutil.Service, value uint32) (changed bool) {
	v.PropsMu.Lock()
	if v.Timeout != value {
		v.Timeout = value
		changed = true
	}
	v.PropsMu.Unlock()
	if service != nil && changed {
		service.EmitPropertyChanged(v, "Timeout", value)
	}
	return
}

func (v *Grub2) getPropTimeout() uint32 {
	v.PropsMu.RLock()
	value := v.Timeout
	v.PropsMu.RUnlock()
	return value
}

func (v *Grub2) setPropUpdating(service *dbusutil.Service, value bool) (changed bool) {
	v.PropsMu.Lock()
	if v.Updating != value {
		v.Updating = value
		changed = true
	}
	v.PropsMu.Unlock()
	if service != nil && changed {
		service.EmitPropertyChanged(v, "Updating", value)
	}
	return
}

func (v *Grub2) getPropUpdating() bool {
	v.PropsMu.RLock()
	value := v.Updating
	v.PropsMu.RUnlock()
	return value
}

func (v *Theme) setPropUpdating(service *dbusutil.Service, value bool) (changed bool) {
	v.PropsMu.Lock()
	if v.Updating != value {
		v.Updating = value
		changed = true
	}
	v.PropsMu.Unlock()
	if service != nil && changed {
		service.EmitPropertyChanged(v, "Updating", value)
	}
	return
}

func (v *Theme) getPropUpdating() bool {
	v.PropsMu.RLock()
	value := v.Updating
	v.PropsMu.RUnlock()
	return value
}

func (v *Theme) setPropItemColor(service *dbusutil.Service, value string) (changed bool) {
	v.PropsMu.Lock()
	if v.ItemColor != value {
		v.ItemColor = value
		changed = true
	}
	v.PropsMu.Unlock()
	if service != nil && changed {
		service.EmitPropertyChanged(v, "ItemColor", value)
	}
	return
}

func (v *Theme) getPropItemColor() string {
	v.PropsMu.RLock()
	value := v.ItemColor
	v.PropsMu.RUnlock()
	return value
}

func (v *Theme) setPropSelectedItemColor(service *dbusutil.Service, value string) (changed bool) {
	v.PropsMu.Lock()
	if v.SelectedItemColor != value {
		v.SelectedItemColor = value
		changed = true
	}
	v.PropsMu.Unlock()
	if service != nil && changed {
		service.EmitPropertyChanged(v, "SelectedItemColor", value)
	}
	return
}

func (v *Theme) getPropSelectedItemColor() string {
	v.PropsMu.RLock()
	value := v.SelectedItemColor
	v.PropsMu.RUnlock()
	return value
}
