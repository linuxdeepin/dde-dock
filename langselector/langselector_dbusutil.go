package langselector

import (
	"pkg.deepin.io/lib/dbusutil"
)

func (v *LangSelector) setPropCurrentLocale(service *dbusutil.Service, value string) (changed bool) {
	v.PropsMu.Lock()
	if v.CurrentLocale != value {
		v.CurrentLocale = value
		changed = true
	}
	v.PropsMu.Unlock()
	if service != nil && changed {
		service.EmitPropertyChanged(v, "CurrentLocale", value)
	}
	return
}

func (v *LangSelector) getPropCurrentLocale() string {
	v.PropsMu.RLock()
	value := v.CurrentLocale
	v.PropsMu.RUnlock()
	return value
}

func (v *LangSelector) setPropLocaleState(service *dbusutil.Service, value int32) (changed bool) {
	v.PropsMu.Lock()
	if v.LocaleState != value {
		v.LocaleState = value
		changed = true
	}
	v.PropsMu.Unlock()
	if service != nil && changed {
		service.EmitPropertyChanged(v, "LocaleState", value)
	}
	return
}

func (v *LangSelector) getPropLocaleState() int32 {
	v.PropsMu.RLock()
	value := v.LocaleState
	v.PropsMu.RUnlock()
	return value
}
