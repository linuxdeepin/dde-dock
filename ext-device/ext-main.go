package main

// #cgo pkg-config: glib-2.0 gdk-3.0
// #include "get-devices.h"
// #include <stdlib.h>
// #include <gdk/gdk.h>
import "C"
import (
	"dlib"
	"dlib/dbus"
	"unsafe"
)

func main() {
	mouseFlag := false
	tpadFlag := false
	keyboardFlag := false

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
		if DeviceIsExist(nameList, "mouse") {
			mouseFlag = true
		}

		if DeviceIsExist(nameList, "touchpad") {
			tpadFlag = true
		}

		if DeviceIsExist(nameList, "keyboard") {
			keyboardFlag = true
		}
	} else {
		C.gdk_init(nil, nil)
		m := C.CString("Mouse")
		defer C.free(unsafe.Pointer(m))
		if C.DeviceIsExist(m) == 1 {
			mouseFlag = true
		}

		t := C.CString("TouchPad")
		defer C.free(unsafe.Pointer(t))
		if C.DeviceIsExist(t) == 1 {
			tpadFlag = true
		}

		k := C.CString("keyboard")
		defer C.free(unsafe.Pointer(k))
		if C.DeviceIsExist(k) == 1 {
			keyboardFlag = true
		}
	}

	if mouseFlag {
		mouse := NewMouseEntry()
		if mouse != nil {
			dbus.InstallOnSession(mouse)
			tmp := ExtDeviceInfo{
				Path: _EXT_ENTRY_PATH + mouse.DeviceID,
				Type: "mouse",
			}
			dev.DevInfoList = append(dev.DevInfoList, tmp)
		}
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

	if keyboardFlag {
		keyboard := NewKeyboardEntry()
		if keyboard != nil {
			dbus.InstallOnSession(keyboard)
			tmp := ExtDeviceInfo{
				Path: _EXT_ENTRY_PATH + keyboard.DeviceID,
				Type: "keyboard",
			}
			dev.DevInfoList = append(dev.DevInfoList, tmp)
		}
	}

	dbus.NotifyChange(dev, "DevInfoList")
	dlib.StartLoop()
}
