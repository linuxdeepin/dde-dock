package launcher

import (
	"fmt"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
	. "pkg.deepin.io/dde/daemon/launcher/setting"
)

type MockSettingCore struct {
	values   map[string]int32
	handlers map[string]func(SettingCore, string)
}

func (m *MockSettingCore) GetEnum(k string) int32 {
	return m.values[k]
}

func (m *MockSettingCore) SetEnum(key string, v int32) bool {
	m.values[key] = v

	detailSignal := fmt.Sprintf("changed::%s", key)
	if fn, ok := m.handlers[detailSignal]; ok {
		fn(m, key)
	}
	return true
}

func (m *MockSettingCore) Connect(signalName string, fn interface{}) {
	f := fn.(func(SettingCore, string))
	m.handlers[signalName] = f
}

func (m *MockSettingCore) Unref() {
}

func NewMockSettingCore() *MockSettingCore {
	s := &MockSettingCore{
		values: map[string]int32{
			CategoryDisplayModeKey: int32(CategoryDisplayModeIcon),
			SortMethodkey:          int32(SortMethodByName),
		},
		handlers: map[string]func(SettingCore, string){},
	}

	return s
}
