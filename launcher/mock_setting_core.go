package launcher

import (
	"fmt"
	. "pkg.linuxdeepin.com/dde-daemon/launcher/interfaces"
	. "pkg.linuxdeepin.com/dde-daemon/launcher/setting"
)

type MockSettingCore struct {
	values   map[string]int
	handlers map[string]func(SettingCoreInterface, string)
}

func (m *MockSettingCore) GetEnum(k string) int {
	return m.values[k]
}

func (m *MockSettingCore) SetEnum(key string, v int) bool {
	m.values[key] = v

	detailSignal := fmt.Sprintf("changed::%s", key)
	if fn, ok := m.handlers[detailSignal]; ok {
		fn(m, key)
	}
	return true
}

func (m *MockSettingCore) Connect(signalName string, fn interface{}) {
	f := fn.(func(SettingCoreInterface, string))
	m.handlers[signalName] = f
}

func (m *MockSettingCore) Unref() {
}

func NewMockSettingCore() *MockSettingCore {
	s := &MockSettingCore{
		values: map[string]int{
			CategoryDisplayModeKey: int(CategoryDisplayModeIcon),
			SortMethodkey:          int(SortMethodByName),
		},
		handlers: map[string]func(SettingCoreInterface, string){},
	}

	return s
}
