package main

// #cgo pkg-config: glib-2.0 gdk-3.0
// #include "get-devices.h"
// #include <stdlib.h>
// #include <gdk/gdk.h>
import "C"
import (
	"dlib/dbus"
	"fmt"
	"unsafe"
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

	m := C.CString("Mouse")
	defer C.free(unsafe.Pointer(m))
	if C.DeviceIsExist(m) == 1 {
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

	t := C.CString("TouchPad")
	defer C.free(unsafe.Pointer(t))
	if C.DeviceIsExist(t) == 1 {
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

	k := C.CString("keyboard")
	defer C.free(unsafe.Pointer(k))
	if C.DeviceIsExist(k) == 1 {
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
