/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package main

//#cgo pkg-config: glib-2.0 libudev
//#include "backlight.h"
//#include <stdlib.h>
import "C"

import (
	"fmt"
	"os"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
	"sync"
	"time"
	"unsafe"
)

const (
	dbusDest = "com.deepin.daemon.helper.Backlight"
	dbusPath = "/com/deepin/daemon/helper/Backlight"
	dbusIFC  = "com.deepin.daemon.helper.Backlight"
)

var logger = log.NewLogger("backlight_helper")

type Manager struct {
	initFailed bool
	locker     sync.Mutex
	kbdLocker  sync.Mutex
}

// ListSysPath return all the backlight device syspath
func (m *Manager) ListSysPath() []string {
	if m.initFailed {
		return nil
	}

	var cNum = C.int(0)
	cList := C.get_syspath_list(&cNum)
	num := int(cNum)
	if num == 0 {
		return nil
	}
	cSlice := (*[1 << 5]*C.char)(unsafe.Pointer(cList))[:num:num]

	var list []string
	for _, cItem := range cSlice {
		list = append(list, C.GoString(cItem))
	}
	C.free_syspath_list(cList, cNum)
	return list
}

// GetSysPathByType return the special type's syspath
// The type range: raw, platform, firmware
func (m *Manager) GetSysPathByType(ty string) (string, error) {
	if m.initFailed {
		return "", fmt.Errorf("Init udev backlight failed")
	}

	switch ty {
	case "raw", "platform", "firmware":
		break
	default:
		return "", fmt.Errorf("Invalid backlight type: %s", ty)
	}

	cTy := C.CString(ty)
	cSysPath := C.get_syspath_by_type(cTy)
	C.free(unsafe.Pointer(cTy))
	sysPath := C.GoString(cSysPath)
	C.free(unsafe.Pointer(cSysPath))
	return sysPath, nil
}

// GetBrightness return the special syspath's brightness
func (m *Manager) GetBrightness(sysPath string) (int32, error) {
	if m.initFailed {
		return 1, fmt.Errorf("Init udev backlight failed")
	}

	m.locker.Lock()
	defer m.locker.Unlock()
	cSysPath := C.CString(sysPath)
	ret := C.get_brightness(cSysPath)
	C.free(unsafe.Pointer(cSysPath))
	if int(ret) == -1 {
		return 1, fmt.Errorf("Get brightness failed for: %s", sysPath)
	}
	return int32(ret), nil
}

func (m *Manager) GetKbdBrightness() (int32, error) {
	if m.initFailed {
		return 1, fmt.Errorf("Init udev backlight failed")
	}

	m.kbdLocker.Lock()
	defer m.kbdLocker.Unlock()
	ret := C.get_kbd_brightness()
	if int(ret) == -1 {
		return 1, fmt.Errorf("Get keyboard brightness failed")
	}
	return int32(ret), nil
}

// GetBrightness return the special syspath's max brightness
func (m *Manager) GetMaxBrightness(sysPath string) (int32, error) {
	if m.initFailed {
		return 1, fmt.Errorf("Init udev backlight failed")
	}

	cSysPath := C.CString(sysPath)
	ret := C.get_max_brightness(cSysPath)
	C.free(unsafe.Pointer(cSysPath))
	if int(ret) == -1 {
		return 1, fmt.Errorf("Get max brightness failed for: %s",
			sysPath)
	}
	return int32(ret), nil
}

func (m *Manager) GetKbdMaxBrightness() (int32, error) {
	if m.initFailed {
		return 1, fmt.Errorf("Init udev backlight failed")
	}

	ret := C.get_kbd_max_brightness()
	if int(ret) == -1 {
		return 1, fmt.Errorf("Get keyboard brightness failed")
	}
	return int32(ret), nil
}

// SetBrightness set the special syspath's brightness
func (m *Manager) SetBrightness(sysPath string, value int32) error {
	if m.initFailed {
		return fmt.Errorf("Init udev backlight failed")
	}

	m.locker.Lock()
	defer m.locker.Unlock()
	cSysPath := C.CString(sysPath)
	ret := C.set_brightness(cSysPath, C.int(value))
	C.free(unsafe.Pointer(cSysPath))
	if int(ret) != 0 {
		return fmt.Errorf("Set brightness for %s to %d failed",
			sysPath, value)
	}
	return nil
}

func (m *Manager) SetKbdBrightness(value int32) error {
	if m.initFailed {
		return fmt.Errorf("Init udev backlight failed")
	}

	m.kbdLocker.Lock()
	defer m.kbdLocker.Unlock()
	ret := C.set_kbd_brightness(C.int(value))
	if int(ret) != 0 {
		return fmt.Errorf("Set keyboard brightness to %d failed",
			value)
	}
	return nil
}

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}

func main() {
	m := &Manager{
		initFailed: false,
	}

	cRet := C.init_udev()
	m.initFailed = (int(cRet) != 0)
	defer C.finalize_udev()

	err := dbus.InstallOnSystem(m)
	if err != nil {
		logger.Error("Install session bus failed:", err)
		return
	}
	dbus.SetAutoDestroyHandler(time.Second*3, nil)
	dbus.DealWithUnhandledMessage()
	err = dbus.Wait()
	if err != nil {
		logger.Error("Lost dbus connection:", err)
		os.Exit(-1)
	}
	os.Exit(0)
}
