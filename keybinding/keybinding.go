package main

import (
	"dlib"
	"dlib/dbus"
	"dlib/gio-2.0"
	"fmt"
	"strconv"
	"strings"
)

type KeyBinding struct {
	KeyBindingList []int32
}

type KeyOwnerRet struct {
	Success bool
	ID      int32
}

const (
	_KEY_BINDING_NAME = "com.deepin.daemon.KeyBinding"
	_KEY_BINDING_PATH = "/com/deepin/daemon/KeyBinding"
	_KEY_BINDING_IFC  = "com.deepin.daemon.KeyBinding"

	_KEY_BINDING_ID       = "com.deepin.dde.key-binding"
	_KEY_BINDING_ADD_ID   = "com.deepin.dde.key-binding.custom"
	_KEY_BINDING_ADD_PATH = "/com/deepin/dde/key-binding/profiles/"

	_WM_BINDING_ID     = "org.gnome.desktop.wm.keybindings"
	_PRESET_BINDING_ID = "org.gnome.settings-daemon.plugins.key-bindings"
	_MEDIA_BINDING_ID  = "org.gnome.settings-daemon.plugins.media-keys"
	_COMPIZ_SHIFT_ID   = "org.compiz.shift"
	_COMPIZ_SHIFT_PATH = "/org/compiz/profiles/shift/"
	_COMPIZ_PUT_ID     = "org.compiz.put"
	_COMPIZ_PUT_PATH   = "/org/compiz/profiles/put/"

	_KEY_COUNT_BASE = 10000
	_KEY_COUNT      = "count"
	_KEY_LIST       = "key-list"
	_KEY_ID         = "id"
	_KEY_NAME       = "name"
	_KEY_SHORTCUT   = "shortcut"
	_KEY_ACTION     = "action"
)

var (
	_bindingGSettings = gio.NewSettings(_KEY_BINDING_ID)
	_presetGSettings  = gio.NewSettings(_PRESET_BINDING_ID)
	_mediaGSettings   = gio.NewSettings(_MEDIA_BINDING_ID)
	_wmGSettings      = gio.NewSettings(_WM_BINDING_ID)
	_shiftGSettings   = gio.NewSettingsWithPath(_COMPIZ_SHIFT_ID,
		_COMPIZ_SHIFT_PATH)
	_putGSettings = gio.NewSettingsWithPath(_COMPIZ_PUT_ID,
		_COMPIZ_PUT_PATH)
)

func (binding *KeyBinding) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_KEY_BINDING_NAME,
		_KEY_BINDING_PATH,
		_KEY_BINDING_IFC,
	}
}

func (binding *KeyBinding) GetSystemList() []int32 {
	sysIDList := []int32{}

	for k, _ := range currentSystemBindings {
		sysIDList = append(sysIDList, int32(k))
	}

	return sysIDList
}

func (binding *KeyBinding) GetCustomList() []int32 {
	return GetCustomIdList()
}

func (binding *KeyBinding) HasOwnID(id int32) bool {
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

func (binding *KeyBinding) HasOwnShortcut(accel string) *KeyOwnerRet {
	accelList := GetKeyAccelList()
	for k, v := range accelList {
		if accel == v {
			return &KeyOwnerRet{Success: true, ID: k}
		}
	}

	return &KeyOwnerRet{Success: false, ID: -1}
}

func (binding *KeyBinding) GetBindingName(id int32) string {
	if id >= 0 && id < _KEY_COUNT_BASE {
		return currentSystemBindings[id]
	} else {
		gs := NewCustomGSettings(id)
		return gs.GetString(_KEY_NAME)
	}
	return ""
}

func (binding *KeyBinding) GetBindingExec(id int32) string {
	if id >= _KEY_COUNT_BASE {
		gs := NewCustomGSettings(id)
		return gs.GetString(_KEY_ACTION)
	} else if id >= 0 && id < 300 {
		values := PresetGetValue(id)
		strArray := strings.Split(values, ";")
		if len(strArray) == 2 {
			return strArray[0]
		}
	}

	return ""
}

func (binding *KeyBinding) GetBindingAccel(id int32) string {
	if id > _KEY_COUNT_BASE {
		gs := NewCustomGSettings(id)
		return gs.GetString(_KEY_SHORTCUT)
	} else if id >= 0 && id < 300 {
		values := PresetGetValue(id)
		strArray := strings.Split(values, ";")
		if len(strArray) == 2 {
			return strArray[1]
		}
	} else if id >= 300 && id < 600 {
		return MediaGetValue(id)
	} else if id >= 600 && id < 800 {
		return WMGetValue(id)
	} else if id >= 800 && id < 900 {
		return CompizShiftValue(id)
	} else if id >= 900 && id < 1000 {
		return CompizPutValue(id)
	}

	return ""
}

func (binding *KeyBinding) ChangeKeyBinding(id int32, accel string) *KeyOwnerRet {
	ret := binding.HasOwnShortcut(accel)
	if ret.Success {
		return &KeyOwnerRet{Success: false, ID: ret.ID}
	}

	if id >= _KEY_COUNT_BASE {
		ModifyCustomKey(binding, id, _KEY_SHORTCUT, accel)
	}
	return &KeyOwnerRet{Success: true, ID: -1}
}

func (binding *KeyBinding) AddCustomBinding(name, shortcut, action string) int32 {
	id := GetMaxIdFromCustom() + 1
	gs := NewCustomGSettings(id)
	SetGSettings(gs, id, name, shortcut, action)

	customList := GetCustomIdList()
	customList = append(customList, id)
	SetCustomList(customList)

	return id
}

func (binding *KeyBinding) DeleteCustomBinding(id int32) {
	if id < _KEY_COUNT_BASE {
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

func ModifyCustomKey(binding *KeyBinding, id int32, key, value string) {
	gs := NewCustomGSettings(id)
	ModifyGSetingsKey(gs, key, value)
}

func NewCustomGSettings(id int32) *gio.Settings {
	customId := strconv.FormatInt(int64(id), 10) + "/"
	gs := gio.NewSettingsWithPath(_KEY_BINDING_ADD_ID, _KEY_BINDING_ADD_PATH+customId)

	return gs
}

func ResetCustomGSettings(gs *gio.Settings) {
	gs.Reset(_KEY_NAME)
	gs.Reset(_KEY_SHORTCUT)
	gs.Reset(_KEY_ACTION)
}

func SetGSettings(gs *gio.Settings, id int32, name, shortcut, action string) {
	gs.SetInt(_KEY_ID, int(id))
	gs.SetString(_KEY_NAME, name)
	gs.SetString(_KEY_SHORTCUT, shortcut)
	gs.SetString(_KEY_ACTION, action)

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

	_bindingGSettings.SetStrv(_KEY_LIST, strList)
	gio.SettingsSync()
}

func GetMaxIdFromCustom() int32 {
	customList := GetCustomIdList()
	max := int32(0)

	for _, v := range customList {
		if max < v {
			max = v
		}
	}

	return max
}

func GetKeyAccelList() map[int32]string {
	accelList := make(map[int32]string)

	for k, _ := range currentSystemBindings {
		if k >= 0 && k < 300 {
			values := PresetGetValue(k)
			strArray := strings.Split(values, ";")
			if len(strArray) == 2 {
				accelList[k] = strArray[1]
			}
		} else if k >= 300 && k < 600 {
			values := MediaGetValue(k)
			accelList[k] = values
		} else if k >= 600 && k < 800 {
			values := WMGetValue(k)
			accelList[k] = values
		} else if k >= 800 && k < 900 {
			values := CompizShiftValue(k)
			accelList[k] = values
		} else if k >= 900 && k < 1000 {
			values := CompizPutValue(k)
			accelList[k] = values
		}
	}

	customList := GetCustomIdList()
	for _, v := range customList {
		gs := NewCustomGSettings(v)
		values := gs.GetString(_KEY_SHORTCUT)
		accelList[v] = values
	}

	return accelList
}

func GetCustomIdList() []int32 {
	customIDList := []int32{}
	strList := _bindingGSettings.GetStrv(_KEY_LIST)

	for _, v := range strList {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			fmt.Println("get custom list failed:", err)
			continue
		}
		customIDList = append(customIDList, int32(id))
	}

	return customIDList
}

func PresetGetValue(id int32) string {
	if id >= 0 && id < 300 {
		keyName := currentSystemBindings[id]

		return _presetGSettings.GetString(keyName)
	}

	return ""
}

func MediaGetValue(id int32) string {
	if id >= 300 && id < 600 {
		keyName := currentSystemBindings[id]

		return _mediaGSettings.GetString(keyName)
	}

	return ""
}

func WMGetValue(id int32) string {
	if id >= 600 && id < 800 {
		keyName := currentSystemBindings[id]

		values := _wmGSettings.GetStrv(keyName)
		strRet := ""

		for _, v := range values {
			strRet += v
		}
		return strRet
	}

	return ""
}

func CompizShiftValue(id int32) string {
	if id >= 800 && id < 900 {
		keyName := currentSystemBindings[id]
		values := _shiftGSettings.GetString(keyName)

		return values
	}

	return ""
}

func CompizPutValue(id int32) string {
	if id >= 900 && id < 1000 {
		keyName := currentSystemBindings[id]
		values := _putGSettings.GetString(keyName)

		return values
	}

	return ""
}

func NewKeyBinding() *KeyBinding {
	binding := KeyBinding{}
	binding.KeyBindingList = GetCustomIdList()

	_bindingGSettings.Connect("changed::key-list", func(s *gio.Settings, name string) {
		binding.KeyBindingList = GetCustomIdList()
		dbus.NotifyChange(&binding, "KeyBindingList")
	})

	return &binding
}

func main() {
	binding := NewKeyBinding()
	err := dbus.InstallOnSession(binding)
	if err != nil {
		panic("Get Session Bus Connect Failed")
	}
	dlib.StartLoop()
}
