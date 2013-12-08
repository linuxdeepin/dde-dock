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
	success bool
	id      int32
}

const (
	_KEY_BINDING_NAME = "com.deepin.daemon.KeyBinding"
	_KEY_BINDING_PATH = "/com/deepin/daemon/KeyBinding"
	_KEY_BINDING_IFC  = "com.deepin.daemon.KeyBinding"

	_KEY_BINDING_ID       = "com.deepin.daemon.key-binding"
	_KEY_BINDING_ADD_ID   = "com.deepin.daemon.key-binding.custom"
	_KEY_BINDING_ADD_PATH = "/com/deepin/daemon/key-binding/profiles/"

	_WM_BINDING_ID     = "org.gnome.desktop.wm.keybindings"
	_PRESET_BINDING_ID = "org.gnome.settings-daemon.plugins.key-bindings"
	_MEDIA_BINDING_ID  = "org.gnome.settings-daemon.plugins.media-keys"

	_KEY_COUNT_BASE = 1000
	_KEY_COUNT      = "count"
	_KEY_ID         = "id"
	_KEY_NAME       = "name"
	_KEY_SHORTCUT   = "shortcut"
	_KEY_ACTION     = "action"
)

var (
	busConn          *dbus.Conn
	bindingGSettings = gio.NewSettings(_KEY_BINDING_ID)
	presetGSettings  = gio.NewSettings(_PRESET_BINDING_ID)
	mediaGSettings   = gio.NewSettings(_MEDIA_BINDING_ID)
	wmGSettings      = gio.NewSettings(_WM_BINDING_ID)
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

	count := bindingGSettings.GetInt(_KEY_COUNT)
	for i := 0; i < count; i++ {
		customIDList = append(customIDList,
			int32(_KEY_COUNT_BASE+i))
	}

	return customIDList
}

func (binding *KeyBinding) HasOwnerID(id int32) KeyOwnerRet {
	accel := ""
	if id >= _KEY_COUNT_BASE {
		gs := NewCustomGSettings(id)
		accel += gs.GetString(_KEY_SHORTCUT)
	} else if id >= 0 && id < 300 {
		values := PresetGetValue(id)
		strArray := strings.Split(values, ";")
		if len(strArray) == 2 {
			accel += strArray[1]
		}
	} else if id >= 300 && id < 600 {
		values := MediaGetValue(id)
		accel += values
	} else if id >= 600 && id < _KEY_COUNT_BASE {
		values := WMGetValue(id)
		accel += values
	}

	fmt.Println(accel)
	accelList := GetKeyAccelList()
	for k, v := range accelList {
		if accel == v {
			fmt.Println("v:", v)
			ret := KeyOwnerRet{success: true, id: k}
			fmt.Println(ret)
			return ret
		}
	}

	return KeyOwnerRet{success: false, id: -1}
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

func (binding *KeyBinding) ChangeKeyBinding(id int32, accel string) KeyOwnerRet {
	/*
	 *ret := binding.HasOwnerID(id)
	 *if ret.success {
	 *        return ret
	 *}
	 */

	if id >= _KEY_COUNT_BASE {
		ModifyCustomKey(binding, id, _KEY_SHORTCUT, accel)
	}
	return KeyOwnerRet{success: true, id: -1}
}

func (binding *KeyBinding) AddCustomBinding(name, shortcut, action string) int32 {
	count := binding.KeyBindingCount
	id := _KEY_COUNT_BASE + count
	gs := NewCustomGSettings(id)
	SetGSettings(gs, id, name, shortcut, action)

	count++
	bindingGSettings.SetInt(_KEY_COUNT, int(count))

	return id
}

func (binding *KeyBinding) DeleteCustomBinding(id int32) {
	if id < _KEY_COUNT_BASE {
		return
	}
	UpdateBindingList(id)

	cnt := binding.KeyBindingCount
	if cnt > 0 {
		bindingGSettings.SetInt(_KEY_COUNT, int(cnt-1))
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
	cnt := bindingGSettings.GetInt(_KEY_COUNT) + _KEY_COUNT_BASE
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

	count := bindingGSettings.GetInt(_KEY_COUNT)
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

		return presetGSettings.GetString(keyName)
	}

	return ""
}

func MediaGetValue(id int32) string {
	if id >= 300 && id < 600 {
		keyName := currentSystemBindings[id]

		return mediaGSettings.GetString(keyName)
	}

	return ""
}

func WMGetValue(id int32) string {
	if id >= 600 && id < 1000 {
		keyName := currentSystemBindings[id]

		values := wmGSettings.GetStrv(keyName)
		strRet := ""

		for _, v := range values {
			strRet += v
		}
		return strRet
	}

	return ""
}

func NewKeyBinding() *KeyBinding {
	var err error
	busConn, err = dbus.SessionBus()
	if err != nil {
		panic("Get Session Bus Connect Failed")
	}

	binding := KeyBinding{}
	binding.KeyBindingCount = int32(bindingGSettings.GetInt(_KEY_COUNT))

	bindingGSettings.Connect("changed::count", func(s *gio.Settings, name string) {
		binding.KeyBindingCount = int32(s.GetInt(name))
		dbus.NotifyChange(&binding, "KeyBindingCount")
	})

	return &binding
}

func main() {
	binding := NewKeyBinding()
	dbus.InstallOnAny(busConn, binding)
	dlib.StartLoop()
}
