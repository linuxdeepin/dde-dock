package main

import (
	"fmt"
	"testing"
)

func TestManager(t *testing.T) {
	manager := NewExtDevManager()
	if manager == nil {
		t.Errorf("create ext device manager failed\n")
		return
	}

	tmp := ExtDeviceInfo{
		DevicePath: "/com/deepin/daemon/ExtEntry/Keyboard",
		DeviceType: "Keyboard",
	}
	manager.DevInfoList = append(manager.DevInfoList, tmp)
	fmt.Println(manager.DevInfoList)
}

func TestKeyboard(t *testing.T) {
	if !InitGSettings() {
		t.Errorf("init gsettings failed\n")
		return
	}

	key := NewKeyboardEntry()
	if key == nil {
		t.Errorf("create new keyboard failed\n")
		return
	}

	fmt.Println(key)
}

func TestMouse(t *testing.T) {
	if !InitGSettings() {
		t.Errorf("init gsettings failed\n")
		return
	}

	mouse := NewMouseEntry()
	if mouse == nil {
		t.Errorf("create new mouse failed\n")
		return
	}

	fmt.Println(mouse)
}

func TestTPad(t *testing.T) {
	if !InitGSettings() {
		t.Errorf("init gsettings failed\n")
		return
	}

	tpad := NewTPadEntry()
	if tpad == nil {
		t.Errorf("create new tpad failed\n")
		return
	}

	fmt.Println(tpad)
}

func TestProcDevice (t *testing.T) {
	success, list := GetProcDeviceNameList ()
	if !success {
		t.Errorf("get proc device list failed\n")
		return
	}

	fmt.Println(list)
	if DeviceIsExist (list, "mouse") {
		fmt.Println ("Mouse Exist")
	}

	if DeviceIsExist (list, "touchpad") {
		fmt.Println ("TouchPad Exist")
	}

	if DeviceIsExist (list, "keyboard") {
		fmt.Println ("keyboard Exist")
	}
}
