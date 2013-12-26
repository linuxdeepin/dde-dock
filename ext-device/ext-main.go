package main

// #cgo pkg-config: x11 xi glib-2.0
// #include "list-devices-info.h"
// #include <stdlib.h>
import "C"
import (
	"dlib"
	"dlib/dbus"
	"unsafe"
)

func main() {
	tpadFlag := false

	if !InitGSettings() {
		return
	}
	dev := NewExtDevManager()
	if dev == nil {
		return
	}
	dbus.InstallOnSession(dev)

	success, nameList := GetProcDeviceNameList()
	if success {
		if DeviceIsExist(nameList, "touchpad") {
			tpadFlag = true
		}
	} else {
		t := C.CString("TOUCHPAD")
		defer C.free(unsafe.Pointer(t))
		if C.find_device_by_name(t) == 1 {
			tpadFlag = true
		}
	}

	mouse := NewMouseEntry()
	if mouse != nil {
		dbus.InstallOnSession(mouse)
		tmp := ExtDeviceInfo{
			Path: _EXT_ENTRY_PATH + mouse.DeviceID,
			Type: "mouse",
		}
		dev.DevInfoList = append(dev.DevInfoList, tmp)
	}

	if tpadFlag {
		tpad := NewTPadEntry()
		if tpad != nil {
			dbus.InstallOnSession(tpad)
			tmp := ExtDeviceInfo{
				Path: _EXT_ENTRY_PATH + tpad.DeviceID,
				Type: "TouchPad",
			}
			dev.DevInfoList = append(dev.DevInfoList, tmp)
		}
	}

	keyboard := NewKeyboardEntry()
	if keyboard != nil {
		dbus.InstallOnSession(keyboard)
		tmp := ExtDeviceInfo{
			Path: _EXT_ENTRY_PATH + keyboard.DeviceID,
			Type: "keyboard",
		}
		dev.DevInfoList = append(dev.DevInfoList, tmp)
	}

	dbus.NotifyChange(dev, "DevInfoList")
	dlib.StartLoop()
}
