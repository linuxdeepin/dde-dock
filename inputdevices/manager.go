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

func getDeviceNames() []string {
	names := []string{}

	n_devices := C.int(0)
	list := C.get_device_info_list(&n_devices)
	defer C.free(unsafe.Pointer(list))
	tmp := uintptr(unsafe.Pointer(list))
	l := unsafe.Sizeof(*list)

	for i := C.int(0); i < n_devices; i++ {
		info := (*C.DeviceInfo)(unsafe.Pointer(tmp + uintptr(i)*l))
		name := C.GoString(info.name)
		atomName := C.GoString(info.atom_name)

		name = strings.ToLower(name)
		atomName = strings.ToLower(atomName)
		if strings.Contains(name, "touchpad") ||
			atomName == "touchpad" {
			names = append(names, "touchpad")
			continue
		}

		names = append(names, atomName)
	}

	return names
}

//export parseDeviceAdd
func parseDeviceAdd(devName *C.char) {
	tmp := C.GoString(devName)
	s := strings.ToLower(tmp)
	logger.Infof("DEVICE CHANGED: %s added", s)
	if strings.Contains(s, "mouse") {
		if managerObj.mouseObj == nil {
			mouse := NewMouse()
			if err := dbus.InstallOnSession(mouse); err != nil {
				logger.Fatal("Mouse DBus Session Failed: ", err)
			}
			managerObj.mouseObj = mouse
			managerObj.setPropName("Infos")
		}
		initMouseSettings()
		disableTPadWhenMouse()
	} else if strings.Contains(s, "touchpad") {
		if managerObj.tpadObj == nil {
			tpad := NewTPad()
			if err := dbus.InstallOnSession(tpad); err != nil {
				logger.Fatal("TPad DBus Session Failed: ", err)
			}
			managerObj.tpadObj = tpad
			managerObj.setPropName("Infos")
		}
		initTPadSettings(true)
	} else if strings.Contains(s, "keyboard") {
		if managerObj.kbdObj == nil {
			kbd := NewKeyboard()
			if err := dbus.InstallOnSession(kbd); err != nil {
				logger.Warning("Kbd DBus Session Failed: ", err)
				panic(err)
			}
			managerObj.kbdObj = kbd
			managerObj.setPropName("Infos")
		}
		initKbdSettings()
	}
}

//export parseDeviceDelete
func parseDeviceDelete(devName *C.char) {
	tmp := C.GoString(devName)
	s := strings.ToLower(tmp)
	logger.Infof("DEVICE CHANGED: %s removed", s)
	if strings.Contains(s, "touchpad") {
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
	} else if strings.Contains(s, "mouse") {
		for _, info := range managerObj.Infos {
			if info.Id == "touchpad" {
				tpadSettings.SetBoolean(TPAD_KEY_ENABLE, true)
			}
		}
	}
	/*
		if strings.Contains(s, "mouse") {
			if managerObj.mouseObj != nil {
				logger.Info("DELETE mouse DBus")
				dbus.UnInstallObject(managerObj.mouseObj)
				managerObj.mouseObj = nil
				logger.Info("DELETE mouse DBus end...")
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
	*/
}

func NewManager() *Manager {
	m := &Manager{}

	m.setPropName("Infos")
	m.mouseObj = nil
	m.tpadObj = nil
	m.kbdObj = nil

	return m
}
