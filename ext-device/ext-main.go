package main

// #cgo pkg-config: glib-2.0 gdk-3.0
// #include "get-devices.h"
// #include <gdk/gdk.h>
import "C"
import (
	"dlib/dbus"
	"fmt"
)

func main() {
	if !InitGSettings() {
		return
	}
	C.gdk_init(nil, nil)
	dev := NewExtDevManager()
	if dev == nil {
		return
	}
	dbus.InstallOnSession(dev)

	if C.DeviceIsExist(C.CString("Mouse")) == 1 {
		mouse := NewMouseEntry()
		if mouse != nil {
			dbus.InstallOnSession(mouse)
			tmp := ExtDeviceInfo{
				DevicePath: _EXT_ENTRY_PATH + mouse.DeviceID,
				DeviceType: "mouse",
			}
			dev.DevInfoList = append(dev.DevInfoList, tmp)
		}
	}

	if C.DeviceIsExist(C.CString("TouchPad")) == 1 {
		tpad := NewTPadEntry()
		if tpad != nil {
			dbus.InstallOnSession(tpad)
			tmp := ExtDeviceInfo{
				DevicePath: _EXT_ENTRY_PATH + tpad.DeviceID,
				DeviceType: "TouchPad",
			}
			dev.DevInfoList = append(dev.DevInfoList, tmp)
		}
	}

	if C.DeviceIsExist(C.CString("keyboard")) == 1 {
		keyboard := NewKeyboardEntry()
		if keyboard != nil {
			dbus.InstallOnSession(keyboard)
			tmp := ExtDeviceInfo{
				DevicePath: _EXT_ENTRY_PATH + keyboard.DeviceID,
				DeviceType: "keyboard",
			}
			dev.DevInfoList = append(dev.DevInfoList, tmp)
		}
	}
	fmt.Println(dev.DevInfoList)
	select {}
}
