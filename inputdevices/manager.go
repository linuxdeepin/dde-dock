/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package inputdevices

// #include <stdlib.h>
// #include "utils.h"
import "C"

import (
	"pkg.linuxdeepin.com/lib/dbus"
	"strings"
	"unsafe"
)

type Manager struct {
	Infos    []deviceInfo
	mouseObj *MouseEntry
	tpadObj  *TPadEntry
	kbdObj   *KbdEntry
}

type deviceInfo struct {
	Path string
	Id   string
}

const (
	_PROC_DEVICE_PATH = "/proc/bus/input/devices"
	_PROC_KEY_NAME    = "N: Name"
)

func getDeviceNames() []string {
	names := []string{}

	n_devices := C.int(0)
	list := C.get_device_info_list(&n_devices)
	defer C.free(unsafe.Pointer(list))
	tmp := uintptr(unsafe.Pointer(list))
	l := unsafe.Sizeof(*list)

	for i := C.int(0); i < n_devices; i++ {
		info := (*C.DeviceInfo)(unsafe.Pointer(tmp + uintptr(i)*l))
		name := C.GoString(info.atom_name)
		names = append(names, strings.ToLower(name))
	}

	return names
}

//export parseDeviceAdd
func parseDeviceAdd(devName *C.char) {
	tmp := C.GoString(devName)
	s := strings.ToLower(tmp)
	logObj.Infof("DEVICE CHANGED: %s added", s)
	if strings.Contains(s, "mouse") {
		if managerObj.mouseObj == nil {
			mouse := NewMouse()
			if err := dbus.InstallOnSession(mouse); err != nil {
				logObj.Warning("Mouse DBus Session Failed: ", err)
				panic(err)
			}
			managerObj.mouseObj = mouse
			managerObj.setPropName("Infos")
			disableTPadWhenMouse()
		}
	} else if strings.Contains(s, "touchpad") {
		if managerObj.tpadObj == nil {
			tpad := NewTPad()
			if err := dbus.InstallOnSession(tpad); err != nil {
				logObj.Warning("TPad DBus Session Failed: ", err)
				panic(err)
			}
			managerObj.tpadObj = tpad
			managerObj.setPropName("Infos")
		}
	} else if strings.Contains(s, "keyboard") {
		if managerObj.kbdObj == nil {
			kbd := NewKeyboard()
			if err := dbus.InstallOnSession(kbd); err != nil {
				logObj.Warning("Kbd DBus Session Failed: ", err)
				panic(err)
			}
			managerObj.kbdObj = kbd
			managerObj.setPropName("Infos")
		}
	}
}

//export parseDeviceDelete
func parseDeviceDelete(devName *C.char) {
	tmp := C.GoString(devName)
	s := strings.ToLower(tmp)
	logObj.Infof("DEVICE CHANGED: %s removed", s)
	if strings.Contains(s, "mouse") {
		if managerObj.mouseObj != nil {
			logObj.Info("DELETE mouse DBus")
			dbus.UnInstallObject(managerObj.mouseObj)
			managerObj.mouseObj = nil
			logObj.Info("DELETE mouse DBus end...")
			managerObj.setPropName("Infos")
			for _, info := range managerObj.Infos {
				if info.Id == "touchpad" {
					tpadSettings.SetBoolean(TPAD_KEY_ENABLE, true)
				}
			}
		}
	} else if strings.Contains(s, "touchpad") {
		if managerObj.tpadObj != nil {
			managerObj.setPropName("Infos")
			for _, info := range managerObj.Infos {
				if info.Id == "touchpad" {
					return
				}
			}
			dbus.UnInstallObject(managerObj.tpadObj)
			managerObj.tpadObj = nil
		}
	} else if strings.Contains(s, "keyboard") {
		if managerObj.kbdObj != nil {
			managerObj.setPropName("Infos")
			for _, info := range managerObj.Infos {
				if info.Id == "keyboard" {
					return
				}
			}
			dbus.UnInstallObject(managerObj.kbdObj)
			managerObj.kbdObj = nil
		}
	}
}

func NewManager() *Manager {
	m := &Manager{}

	m.setPropName("Infos")
	m.mouseObj = nil
	m.tpadObj = nil
	m.kbdObj = nil

	return m
}
