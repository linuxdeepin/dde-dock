package main

import (
	"dlib"
	"dlib/dbus"
	"dlib/gio-2.0"
	"fmt"
	"strconv"
)

type KeyBinding struct {
	KeyBindingCount int32
}

const (
	_KEY_BINDING_NAME = "com.deepin.daemon.KeyBinding"
	_KEY_BINDING_PATH = "/com/deepin/daemon/KeyBinding"
	_KEY_BINDING_IFC  = "com.deepin.daemon.KeyBinding"

	_KEY_BINDING_ID       = "com.deepin.daemon.key-binding"
	_KEY_BINDING_ADD_ID   = "com.deepin.daemon.key-binding.custom"
	_KEY_BINDING_ADD_PATH = "/com/deepin/daemon/key-binding/profiles/"

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
)

func (binding *KeyBinding) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_KEY_BINDING_NAME,
		_KEY_BINDING_PATH,
		_KEY_BINDING_IFC,
	}
}

/*
func (binding *KeyBinding) GetSystemList() []int32 {
	return nil
}

func (binding *KeyBinding) GetCustomList() []int32 {
	return nil
}

func (binding *KeyBinding) HasOwnerID(id int32) bool {
	return true
}

func (binding *KeyBinding) GetBindingName(id int32) string {
	return ""
}

func (binding *KeyBinding) GetBindingExec(id int32) string {
	return ""
}

func (binding *KeyBinding) GetBindingAccel(id int32) string {
	return ""
}

func (binding *KeyBinding) AddKeyBinding(name, exec string) int32 {
	return 0
}

func (binding *KeyBinding) ChangeKeyBinding(id int32, accel string) (bool, int32) {
	return true, 0
}

func (binding *KeyBinding) DeleteKeyBinding(id int32) {
}
*/

func (binding *KeyBinding) AddCustomBinding(name, shortcut, action string) int32 {
	count := binding.KeyBindingCount
	id := _KEY_COUNT_BASE + count
	gs := NewCustomGSettings(id)
	SetGSettings(gs, id, name, shortcut, action)

	count++
	bindingGSettings.SetInt(_KEY_COUNT, int(count))

	return id
}

func (binding *KeyBinding) ModifyCustomKey(id int32, key, value string) bool {
	gs := NewCustomGSettings(id)

	ModifyGSetingsKey(gs, key, value)

	return true
}

func (binding *KeyBinding) DeleteCustomBinding(id int32) {
	UpdateBindingList(id)

	cnt := binding.KeyBindingCount
	if cnt > 0 {
		bindingGSettings.SetInt(_KEY_COUNT, int(cnt-1))
	}
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
