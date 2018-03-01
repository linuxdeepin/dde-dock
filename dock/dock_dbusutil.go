package dock

import (
	"github.com/BurntSushi/xgb/xproto"
)

func (v *AppEntry) setPropId(value string) (changed bool) {
	if v.Id != value {
		v.Id = value
		v.emitPropChangedId(value)
		return true
	}
	return false
}

func (v *AppEntry) emitPropChangedId(value string) error {
	return v.service.EmitPropertyChanged(v, "Id", value)
}

func (v *AppEntry) setPropIsActive(value bool) (changed bool) {
	if v.IsActive != value {
		v.IsActive = value
		v.emitPropChangedIsActive(value)
		return true
	}
	return false
}

func (v *AppEntry) emitPropChangedIsActive(value bool) error {
	return v.service.EmitPropertyChanged(v, "IsActive", value)
}

func (v *AppEntry) setPropName(value string) (changed bool) {
	if v.Name != value {
		v.Name = value
		v.emitPropChangedName(value)
		return true
	}
	return false
}

func (v *AppEntry) emitPropChangedName(value string) error {
	return v.service.EmitPropertyChanged(v, "Name", value)
}

func (v *AppEntry) setPropIcon(value string) (changed bool) {
	if v.Icon != value {
		v.Icon = value
		v.emitPropChangedIcon(value)
		return true
	}
	return false
}

func (v *AppEntry) emitPropChangedIcon(value string) error {
	return v.service.EmitPropertyChanged(v, "Icon", value)
}

func (v *AppEntry) setPropMenu(value string) (changed bool) {
	if v.Menu != value {
		v.Menu = value
		v.emitPropChangedMenu(value)
		return true
	}
	return false
}

func (v *AppEntry) emitPropChangedMenu(value string) error {
	return v.service.EmitPropertyChanged(v, "Menu", value)
}

func (v *AppEntry) setPropDesktopFile(value string) (changed bool) {
	if v.DesktopFile != value {
		v.DesktopFile = value
		v.emitPropChangedDesktopFile(value)
		return true
	}
	return false
}

func (v *AppEntry) emitPropChangedDesktopFile(value string) error {
	return v.service.EmitPropertyChanged(v, "DesktopFile", value)
}

func (v *AppEntry) setPropWindowInfos(value windowInfosType) (changed bool) {
	if !v.WindowInfos.Equal(value) {
		v.WindowInfos = value
		v.emitPropChangedWindowInfos(value)
		return true
	}
	return false
}

func (v *AppEntry) emitPropChangedWindowInfos(value windowInfosType) error {
	return v.service.EmitPropertyChanged(v, "WindowInfos", value)
}

func (v *AppEntry) setPropCurrentWindow(value xproto.Window) (changed bool) {
	if v.CurrentWindow != value {
		v.CurrentWindow = value
		v.emitPropChangedCurrentWindow(value)
		return true
	}
	return false
}

func (v *AppEntry) emitPropChangedCurrentWindow(value xproto.Window) error {
	return v.service.EmitPropertyChanged(v, "CurrentWindow", value)
}

func (v *AppEntry) setPropIsDocked(value bool) (changed bool) {
	if v.IsDocked != value {
		v.IsDocked = value
		v.emitPropChangedIsDocked(value)
		return true
	}
	return false
}

func (v *AppEntry) emitPropChangedIsDocked(value bool) error {
	return v.service.EmitPropertyChanged(v, "IsDocked", value)
}
