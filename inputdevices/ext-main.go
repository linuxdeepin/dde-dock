package main

// #cgo pkg-config: gtk+-3.0 gnome-desktop-3.0 x11 xi xkbfile glib-2.0
// #cgo LDFLAGS: -lm
// #include "list-devices-info.h"
// #include "gsd-init.h"
// #include <stdlib.h>
import "C"
import (
        "dlib"
        "dlib/dbus"
        "dlib/logger"
        "os"
        "unsafe"
)

var (
        logObject = logger.NewLogger("daemon/inputdevices")
)

func main() {
        defer logObject.EndTracing()

        if !dlib.UniqueOnSession(_EXT_DEV_NAME) {
                logObject.Warning("There already has an InputDevices daemon running.")
                return
        }
        tpadFlag := false

        logObject.SetRestartCommand("/usr/lib/deepin-daemon/inputdevices")

        go C.gsd_init()
        defer C.gsd_finalize()

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
                        Path:   _EXT_ENTRY_PATH + mouse.DeviceID,
                        Type:   "mouse",
                }
                dev.DevInfoList = append(dev.DevInfoList, tmp)
        }

        if tpadFlag {
                tpad := NewTPadEntry()
                if tpad != nil {
                        dbus.InstallOnSession(tpad)
                        tmp := ExtDeviceInfo{
                                Path:   _EXT_ENTRY_PATH + tpad.DeviceID,
                                Type:   "touchpad",
                        }
                        dev.DevInfoList = append(dev.DevInfoList, tmp)
                }
        }

        keyboard := NewKeyboardEntry()
        if keyboard != nil {
                dbus.InstallOnSession(keyboard)
                tmp := ExtDeviceInfo{
                        Path:   _EXT_ENTRY_PATH + keyboard.DeviceID,
                        Type:   "keyboard",
                }
                dev.DevInfoList = append(dev.DevInfoList, tmp)
        }

        dbus.NotifyChange(dev, "DevInfoList")
        dbus.DealWithUnhandledMessage()
        go dlib.StartLoop()
        if err := dbus.Wait(); err != nil {
                logObject.Info("lost dbus session:", err)
                os.Exit(1)
        } else {
                os.Exit(0)
        }
}
