package mock

import (
	"fmt"
	ifc "pkg.deepin.io/dde/daemon/launcher/interfaces"
	s "pkg.deepin.io/dde/daemon/launcher/setting"
)

type SettingCore struct {
	values   map[string]int32
	handlers map[string]func(ifc.SettingCore, string)
}

func (m *SettingCore) GetEnum(k string) int32 {
	return m.values[k]
}

func (m *SettingCore) SetEnum(key string, v int32) bool {
	m.values[key] = v

	detailSignal := fmt.Sprintf("changed::%s", key)
	if fn, ok := m.handlers[detailSignal]; ok {
		fn(m, key)
	}
	return true
}

func (m *SettingCore) Connect(signalName string, fn interface{}) {
	f := fn.(func(ifc.SettingCore, string))
	m.handlers[signalName] = f
}

func (m *SettingCore) Unref() {
}

func NewSettingCore() *SettingCore {
	s := &SettingCore{
		values: map[string]int32{
			s.CategoryDisplayModeKey: int32(s.CategoryDisplayModeIcon),
			s.SortMethodkey:          int32(s.SortMethodByName),
		},
		handlers: map[string]func(ifc.SettingCore, string){},
	}

	return s
}
