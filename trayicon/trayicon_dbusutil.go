package trayicon

func (v *TrayManager) setPropTrayIcons(value []uint32) {
	v.TrayIcons = value
	v.emitPropChangedTrayIcons(value)
}

func (v *TrayManager) emitPropChangedTrayIcons(value []uint32) error {
	return v.service.EmitPropertyChanged(v, "TrayIcons", value)
}
