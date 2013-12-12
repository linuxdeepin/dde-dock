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
	KeyBindingCount int32
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

	_KEY_COUNT_BASE = 10000
	_KEY_COUNT      = "count"
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
	customIDList := []int32{}

	count := _bindingGSettings.GetInt(_KEY_COUNT)
	for i := 0; i < count; i++ {
		customIDList = append(customIDList,
			int32(_KEY_COUNT_BASE+i))
	}

	return customIDList
}

func (binding *KeyBinding) HasOwnID(id int32) bool {
	/*sysIDList := binding.GetSystemList ()*/
	/*customIDList := binding.GetCustomList ()*/
	for i, _ := range currentSystemBindings {
		if id == i {
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
	} else if id >= 600 && id < 1000 {
		return WMGetValue(id)
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
	count := binding.KeyBindingCount
	id := _KEY_COUNT_BASE + count
	gs := NewCustomGSettings(id)
	SetGSettings(gs, id, name, shortcut, action)

	count++
	_bindingGSettings.SetInt(_KEY_COUNT, int(count))

	return id
}

func (binding *KeyBinding) DeleteCustomBinding(id int32) {
	if id < _KEY_COUNT_BASE {
		return
	}
	UpdateBindingList(id)

	cnt := binding.KeyBindingCount
	if cnt > 0 {
		_bindingGSettings.SetInt(_KEY_COUNT, int(cnt-1))
	}
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

func UpdateBindingList(id int32) {
	cnt := _bindingGSettings.GetInt(_KEY_COUNT) + _KEY_COUNT_BASE
	fmt.Println("id:", id)
	fmt.Println("cnt:", cnt)

	i := id
	for ; i < int32(cnt-1); i++ {
		gsSrc := NewCustomGSettings(i)
		gsDest := NewCustomGSettings(i + 1)
		ReplaceGSettings(gsSrc, gsDest)
	}

	fmt.Println("i:", i)
	gs := NewCustomGSettings(i)
	ResetGSettings(gs)
}

func ReplaceGSettings(src, dest *gio.Settings) {
	SetGSettings(src,
		int32(src.GetInt(_KEY_ID)),
		dest.GetString(_KEY_NAME),
		dest.GetString(_KEY_SHORTCUT),
		dest.GetString(_KEY_ACTION),
	)
}

func ResetGSettings(gs *gio.Settings) {
	gs.Reset(_KEY_ID)
	gs.Reset(_KEY_NAME)
	gs.Reset(_KEY_SHORTCUT)
	gs.Reset(_KEY_ACTION)

	gio.SettingsSync()
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
		} else if k >= 600 && k < 1000 {
			values := WMGetValue(k)
			accelList[k] = values
		}
	}

	count := _bindingGSettings.GetInt(_KEY_COUNT)
	for i := 0; i < count; i++ {
		id := int32(_KEY_COUNT_BASE + i)
		gs := NewCustomGSettings(id)
		values := gs.GetString(_KEY_SHORTCUT)
		accelList[id] = values
	}

	return accelList
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
	if id >= 600 && id < 1000 {
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

func NewKeyBinding() *KeyBinding {
	binding := KeyBinding{}
	binding.KeyBindingCount = int32(_bindingGSettings.GetInt(_KEY_COUNT))

	_bindingGSettings.Connect("changed::count", func(s *gio.Settings, name string) {
		binding.KeyBindingCount = int32(s.GetInt(name))
		dbus.NotifyChange(&binding, "KeyBindingCount")
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
