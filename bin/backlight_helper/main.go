package main

import "pkg.linuxdeepin.com/lib/dbus"
import liblogger "pkg.linuxdeepin.com/lib/logger"

//#cgo pkg-config: libudev
//#include "backlight.h"
//#include <stdlib.h>
import "C"
import "os"
import "time"

type BacklightHelper struct {
	SysPath string
}

func NewBacklightHelper() *BacklightHelper {
	C.update_backlight_device()
	return &BacklightHelper{}
}

func (*BacklightHelper) SetBrightness(v float64) {
	if v > 1 || v < 0 {
		logger.Warningf("SetBacklight %v failed\n", v)
		return
	}
	C.set_backlight(C.double(v))
}
func (*BacklightHelper) GetBrightness() float64 {
	return (float64)(C.get_backlight())
}
func (*BacklightHelper) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.helper.BacklightHelper",
		"/com/deepin/daemon/helper/BacklightHelper",
		"com.deepin.daemon.helper.BacklightHelper",
	}
}

var logger = liblogger.NewLogger("com.deepin.daemon.helper.BacklightHelper")

func main() {
	helper := NewBacklightHelper()
	err := dbus.InstallOnSystem(helper)
	if err != nil {
		logger.Errorf("register dbus interface failed: %v", err)
		os.Exit(1)
	}

	dbus.SetAutoDestroyHandler(time.Second*1, nil)

	dbus.DealWithUnhandledMessage()
	if err := dbus.Wait(); err != nil {
		logger.Errorf("lost dbus session: %v", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
