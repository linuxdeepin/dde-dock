package main

import (
	"dlib/dbus"
	"dlib/dbus/property"
	"dlib/gio-2.0"
	"strconv"
)

type KeyBinding struct {
	KeyBindingCount int32 `access:"read"`
}

const (
	_KEY_BINDING_NAME = "com.deepin.daemon.KeyBinding"
	_KEY_BINDING_PATH = "/com/deepin/daemon/KeyBinding"
	_KEY_BINDING_IFC  = "com.deepin.daemon.KeyBinding"

	_KEY_BINDING_ID       = "com.deepin.daemon.key-binding"
	_KEY_BINDING_ADD_ID   = "com.deepin.daemon.key-binding.key"
	_KEY_BINDING_ADD_PATH = "/com/deepin/daemon/key-binding/profiles/"

	_KEY_COUNT    = "count"
	_KEY_NAME     = "name"
	_KEY_SHORTCUT = "shortcut"
	_KEY_ACTION   = "action"
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

func NewKeyBinding() *KeyBinding {
	var err error
	busConn, err = dbus.SessionBus()
	if err != nil {
		panic("Get Session Bus Connect Failed")
	}

	binding := KeyBinding{}
	binding.KeyBindingCount = bindingGSettings.GetInt(_KEY_COUNT)

	bindingGSettings.Connect("changed::count", func(s *gio.Settings, name string) {
		binding.KeyBindingCount = s.GetInt(name)
		dbus.NotifyChange(binding, "KeyBindingCount")
	})

	return &binding
}

func (binding *KeyBinding) AddCustomBinding(name, shortcut, action string) string {
	count := binding.KeyBindingCount
	id := "custom" + strconv.FormatInt(int64(count), 10)
	gs := NewCustomGSettings(id)
	AddGSettings(gs, name, shortcut, action)

	count++
	bindingGSettings.SetInt(_KEY_COUNT, count)

	return id
}

func (binding *KeyBinding) ModifyCustomKey(id, key, value string) bool {
	gs := NewCustomGSettings(id)

	ModifyGSetingsKey(gs, key, value)

	return true
}

func (binding *KeyBinding) DeleteCustomBinding(id string) {
	gs := NewCustomGSettings(id)

	ResetGSettings(gs)
}

func NewCustomGSettings(id string) *gio.Settings {
	customId := id + "/"
	gs := gio.NewSettingsWithPath(_KEY_BINDING_ADD_ID, _KEY_BINDING_ADD_PATH+customId)

	return gs
}

func AddGSettings(gs *gio.Settings, name, shortcut, action string) {
	gs.SetString(_KEY_NAME, name)
	gs.SetString(_KEY_SHORTCUT, shortcut)
	gs.SetString(_KEY_ACTION, action)

	gio.SettingsSync()
}

func ModifyGSetingsKey(gs *gio.Settings, key, value string) {
	gs.SetString(key, value)

	gio.SettingsSync()
}

func UpdateBindingList (id int32) {
	cnt := bindingGSettings.GetInt(_KEY_COUNT)

	for i := id; i < (cnt - 1); i++ {
		customSrc := "custom" + strconv.FormatInt(int64(cnt), 10)
		customDest := "custom" + strconv.FormatInt(int64(cnt + 1), 10)
		gsSrc := NewCustomGSettings (customSrc)
		gsDest := NewCustomGSettings (customDest)
	}
}

func ReplaceGSettings (src, dest *gio.Settings) {
}

func ResetGSettings(gs *gio.Settings) {
	gs.Reset(_KEY_NAME)
	gs.Reset(_KEY_SHORTCUT)
	gs.Reset(_KEY_ACTION)

	gio.SettingsSync()
}

func main() {
	binding := KeyBinding{}
	dbus.InstallOnSession(&binding)
	select {}
}
