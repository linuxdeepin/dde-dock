package main

import (
	"dlib/dbus"
	"dlib/gio-2.0"
	"fmt"
	"strconv"
	"strings"
)

type KeyOwnerRet struct {
	Success bool
	ID      int32
}

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_KEY_BINDING_NAME,
		_KEY_BINDING_PATH,
		_KEY_BINDING_IFC,
	}
}

func (m *Manager) SystemList() []int32 {
	sysIDList := []int32{}

	for k, _ := range currentSystemBindings {
		sysIDList = append(sysIDList, int32(k))
	}

	return sysIDList
}

func (m *Manager) CustomList() []int32 {
	return GetCustomIdList()
}

func (m *Manager) HasOwnID(id int32) bool {
	for i, _ := range currentSystemBindings {
		if id == i {
			return true
		}
	}

	customIDList := GetCustomIdList()
	for _, v := range customIDList {
		if id == v {
			return true
		}
	}

	return false
}

func (m *Manager) HasOwnShortcut(accel string) *KeyOwnerRet {
	tmp := GenericKeyInfo(accel)
	if tmp == nil {
		fmt.Println("shortcut error...\n")
		return &KeyOwnerRet{Success: true, ID: -1}
	}

	accelList := GetKeyAccelList()
	for k, v := range accelList {
		if CompareKeyInfo(tmp, v) {
			return &KeyOwnerRet{Success: true, ID: k}
		}
	}

	return &KeyOwnerRet{Success: false, ID: -1}
}

func (m *Manager) GetBindingName(id int32) string {
	if !m.HasOwnID(id) {
		return ""
	}

	if id >= 0 && id < _CUSTOM_KEY_BASE {
		return currentSystemBindings[id]
	} else {
		gs := NewCustomGSettings(id)
		return gs.GetString(_CUSTOM_KEY_NAME)
	}
	return ""
}

func (m *Manager) GetBindingExec(id int32) string {
	if !m.HasOwnID(id) {
		return ""
	}

	if id >= _CUSTOM_KEY_BASE {
		return GetCustomValue(id, _CUSTOM_KEY_ACTION)
	} else if id >= 0 && id < 300 {
		values := GSDGetValue(id)
		strArray := strings.Split(values, ";")
		if len(strArray) == 2 {
			return strArray[0]
		}
	}

	return ""
}

func (m *Manager) GetBindingAccel(id int32) string {
	if !m.HasOwnID(id) {
		return ""
	}

	shortcut := ""
	if id > _CUSTOM_KEY_BASE {
		shortcut = GetCustomValue(id, _CUSTOM_KEY_SHORTCUT)
	} else if id >= 0 && id < 300 {
		values := GSDGetValue(id)
		strArray := strings.Split(values, ";")
		if len(strArray) == 2 {
			shortcut = strArray[1]
		}
	} else if id >= 300 && id < 600 {
		return MediaGetValue(id)
	} else if id >= 600 && id < 800 {
		shortcut = WMGetValue(id)
	} else if id >= 800 && id < 900 {
		shortcut = CompizShiftValue(id)
	} else if id >= 900 && id < 1000 {
		shortcut = CompizPutValue(id)
	}

	return FormatShortcut(shortcut)
}

func (m *Manager) ChangeKeyBinding(id int32, accel string) *KeyOwnerRet {
	if !m.HasOwnID(id) {
		return &KeyOwnerRet{Success: false, ID: -1}
	}

	ret := m.HasOwnShortcut(accel)
	if ret.Success {
		return &KeyOwnerRet{Success: false, ID: ret.ID}
	}

	if id >= _CUSTOM_KEY_BASE {
		ModifyCustomKey(m, id, _CUSTOM_KEY_SHORTCUT, accel)
	}
	return &KeyOwnerRet{Success: true, ID: -1}
}

func (m *Manager) AddCustomBinding(name, shortcut, action string) int32 {
	id := GetMaxIdFromCustom() + 1
	gs := NewCustomGSettings(id)
	SetGSettings(gs, id, name, shortcut, action)

	customList := GetCustomIdList()
	customList = append(customList, id)
	SetCustomList(customList)

	gs.Connect("changed::shortcut", func(s *gio.Settings, key string) {
		tmpPairs := GetCustomPairs()
		BindingKeysPairs(m.customAccelMap, false)
		BindingKeysPairs(tmpPairs, true)
		m.customAccelMap = tmpPairs
	})

	return id
}

func (m *Manager) DeleteCustomBinding(id int32) {
	if id < _CUSTOM_KEY_BASE {
		return
	}

	if !m.HasOwnID(id) {
		return
	}

	gs := NewCustomGSettings(id)
	ResetCustomGSettings(gs)

	tmpList := []int32{}
	customList := GetCustomIdList()
	for _, v := range customList {
		if v == id {
			continue
		}
		tmpList = append(tmpList, v)
	}
	SetCustomList(tmpList)
}

func ModifyCustomKey(m *Manager, id int32, key, value string) {
	gs := NewCustomGSettings(id)
	ModifyGSetingsKey(gs, key, value)
}

func ResetCustomGSettings(gs *gio.Settings) {
	gs.Reset(_CUSTOM_KEY_ID)
	gs.Reset(_CUSTOM_KEY_NAME)
	gs.Reset(_CUSTOM_KEY_SHORTCUT)
	gs.Reset(_CUSTOM_KEY_ACTION)
}

func SetGSettings(gs *gio.Settings, id int32, name, shortcut, action string) {
	gs.SetInt(_CUSTOM_KEY_ID, int(id))
	gs.SetString(_CUSTOM_KEY_NAME, name)
	gs.SetString(_CUSTOM_KEY_SHORTCUT, shortcut)
	gs.SetString(_CUSTOM_KEY_ACTION, action)

	gio.SettingsSync()
}

func ModifyGSetingsKey(gs *gio.Settings, key, value string) {
	gs.SetString(key, value)
	gio.SettingsSync()
}

func SetCustomList(customList []int32) {
	strList := []string{}
	for _, v := range customList {
		str := strconv.FormatInt(int64(v), 10)
		strList = append(strList, str)
	}

	customGSettings.SetStrv(_CUSTOM_KEY_LIST, strList)
	gio.SettingsSync()
}

func GetCustomValue(id int32, key string) string {
	customList := GetCustomIdList()

	for _, v := range customList {
		if id == v {
			gs := NewCustomGSettings(id)
			return gs.GetString(key)
		}
	}

	return ""
}

/*
 * Listen custom 'key-list' changed
 */
func ListenKeyList(m *Manager) {
	customGSettings.Connect("changed::key-list", func(s *gio.Settings, name string) {
		tmpPairs := GetCustomPairs()
		BindingKeysPairs(m.customAccelMap, false)
		BindingKeysPairs(tmpPairs, true)
		m.customAccelMap = tmpPairs
		m.CustomBindList = GetCustomIdList()
		dbus.NotifyChange(m, "CustomBindList")
	})
}

func ListenCustomKey(m *Manager) {
	customList := GetCustomIdList()

	for _, k := range customList {
		print("id: ", k, "\n")
		gs := NewCustomGSettings(k)

		gs.Connect("changed::shortcut", func(s *gio.Settings, key string) {
			tmpPairs := GetCustomPairs()
			BindingKeysPairs(m.customAccelMap, false)
			BindingKeysPairs(tmpPairs, true)
			m.customAccelMap = tmpPairs
		})
	}
}
