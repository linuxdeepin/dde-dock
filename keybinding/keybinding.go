package main

import (
	/*"dlib"*/
	"dlib/dbus"
)

type KeyBinding struct{}

func (binding *KeyBinding) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.KeyBinding",
		"/com/deepin/daemon/KeyBinding",
		"com.deepin.daemon.KeyBinding",
	}
}

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

func main() {
	binding := KeyBinding{}
	dbus.InstallOnSession(&binding)
	select {}
}
