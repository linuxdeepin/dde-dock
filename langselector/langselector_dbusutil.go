package langselector

func (v *LangSelector) setPropCurrentLocale(value string) (changed bool) {
	if v.CurrentLocale != value {
		v.CurrentLocale = value
		v.emitPropChangedCurrentLocale(value)
		return true
	}
	return false
}

func (v *LangSelector) emitPropChangedCurrentLocale(value string) error {
	return v.service.EmitPropertyChanged(v, "CurrentLocale", value)
}

func (v *LangSelector) setPropLocaleState(value int32) (changed bool) {
	if v.LocaleState != value {
		v.LocaleState = value
		v.emitPropChangedLocaleState(value)
		return true
	}
	return false
}

func (v *LangSelector) emitPropChangedLocaleState(value int32) error {
	return v.service.EmitPropertyChanged(v, "LocaleState", value)
}
