package main

import "dlib/dbus"

type KeyBinding struct {
	RepeatDelay    float64
	RepeatSpeed    float64
	CursorFlash    float64
	ForbiddenTPad  float64
	KeyboardLayout string

	RepeatDelayChanged    func(float64)
	RepeatSpeedChanged    func(float64)
	CursorFlashChanged    func(float64)
	ForbiddenTPadChanged  func(float64)
	KeyboardLayoutChanged func(string)
}

func (keybind *KeyBinding) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Keybind",
		"/com/deepin/daemon/Keybind",
		"com.deepin.daemon.Keybind",
	}
}

func (keybind *KeyBinding) GetKeybindList(value string) map[string]string {
	return nil
}

func (keybind *KeyBinding) ChangeKeybind(name, value string) bool {
	return true
}

func (keybind *KeyBinding) AddKeybind(name, desc, value string) bool {
	return true
}

func (keybind *KeyBinding) DeleteKeybind(name string) bool {
	return true
}

func main() {
	kb := KeyBinding{}

	dbus.InstallOnSession(&kb)

	select {}
}
