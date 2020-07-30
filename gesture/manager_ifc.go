package gesture

import (
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
)

func (m *Manager) SetLongPressDuration(duration uint32) *dbus.Error {
	if m.tsSetting.GetInt(tsSchemaKeyLongPress) == int32(duration) {
		return nil
	}
	err := m.sysDaemon.SetLongPressDuration(0, duration)
	if err != nil {
		return dbusutil.ToError(err)
	}
	m.tsSetting.SetInt(tsSchemaKeyLongPress, int32(duration))
	return nil
}

func (m *Manager) GetLongPressDuration() (uint32, *dbus.Error) {
	return uint32(m.tsSetting.GetInt(tsSchemaKeyLongPress)), nil
}

func (m *Manager) SetShortPressDuration(duration uint32) *dbus.Error {
	if m.tsSetting.GetInt(tsSchemaKeyShortPress) == int32(duration) {
		return nil
	}
	err := m.gesture.SetShortPressDuration(0, duration)
	if err != nil {
		return dbusutil.ToError(err)
	}
	m.tsSetting.SetInt(tsSchemaKeyShortPress, int32(duration))
	return nil
}

func (m *Manager) GetShortPressDuration() (uint32, *dbus.Error) {
	return uint32(m.tsSetting.GetInt(tsSchemaKeyShortPress)), nil

}
