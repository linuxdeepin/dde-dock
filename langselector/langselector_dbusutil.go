package langselector

func (v *LangSelector) setPropCurrentLocale(value string) (changed bool) {
	v.PropsMaster.Lock()
	if v.CurrentLocale != value {
		v.CurrentLocale = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "CurrentLocale", value)
	}
	return
}

func (v *LangSelector) getPropCurrentLocale() string {
	v.PropsMaster.RLock()
	value := v.CurrentLocale
	v.PropsMaster.RUnlock()
	return value
}

func (v *LangSelector) setPropLocaleState(value int32) (changed bool) {
	v.PropsMaster.Lock()
	if v.LocaleState != value {
		v.LocaleState = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "LocaleState", value)
	}
	return
}

func (v *LangSelector) getPropLocaleState() int32 {
	v.PropsMaster.RLock()
	value := v.LocaleState
	v.PropsMaster.RUnlock()
	return value
}
