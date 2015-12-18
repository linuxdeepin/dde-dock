package main

//#cgo pkg-config: libudev
//#include "backlight.h"
//#include <stdlib.h>
import "C"

import (
	"os"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
	"sync"
	"time"
	"unsafe"
)

type BacklightHelper struct {
	SysPath string
	lock    sync.Mutex
}

func NewBacklightHelper() *BacklightHelper {
	return &BacklightHelper{}
}

// If driver is empty, auto detect
func (h *BacklightHelper) SetBrightness(v float64, driver string) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if v > 1 || v < 0 {
		logger.Warningf("SetBacklight %v failed\n", v)
		return
	}
	cDriver := C.CString(driver)
	defer C.free(unsafe.Pointer(cDriver))
	C.set_backlight(C.double(v), cDriver)
}

// If driver is empty, auto detect
func (h *BacklightHelper) GetBrightness(driver string) float64 {
	h.lock.Lock()
	defer h.lock.Unlock()
	cDriver := C.CString(driver)
	defer C.free(unsafe.Pointer(cDriver))
	return (float64)(C.get_backlight(cDriver))
}
func (*BacklightHelper) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.daemon.helper.Backlight",
		ObjectPath: "/com/deepin/daemon/helper/Backlight",
		Interface:  "com.deepin.daemon.helper.Backlight",
	}
}

var logger = log.NewLogger("com.deepin.daemon.helper.Backlight")

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
