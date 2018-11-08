package gesture

import (
	"pkg.deepin.io/lib/dbus1"
)

func (m *Manager) SetLongPressDuration(duration uint32) *dbus.Error {
	if m.tsSetting.GetInt(tsSchemaKeyLongPress) == int32(duration) {
		return nil
	}
	err := m.sysDaemon.SetLongPressDuration(0, duration)
	if err != nil {
		return dbus.NewError("SetLongPressDuration",
			[]interface{}{err.Error()})
	}
	m.tsSetting.SetInt(tsSchemaKeyLongPress, int32(duration))
	return nil
}

func (m *Manager) GetLongPressDuration() (uint32, *dbus.Error) {
	return uint32(m.tsSetting.GetInt(tsSchemaKeyLongPress)), nil
}
