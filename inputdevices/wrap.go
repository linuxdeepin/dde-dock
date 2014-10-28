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

// #cgo pkg-config: x11 xi
// #cgo CFLAGS: -Wall -g
// #cgo LDFLAGS: -lpthread
// #include "devices.h"
import "C"

import (
	"pkg.linuxdeepin.com/dde-daemon/inputdevices/libmouse"
	"pkg.linuxdeepin.com/dde-daemon/inputdevices/libtouchpad"
	"pkg.linuxdeepin.com/dde-daemon/inputdevices/libwacom"
	"pkg.linuxdeepin.com/dde-daemon/inputdevices/libwrapper"
)

func initDeviceChangedWatcher() bool {
	ret := C.listen_device_changed()
	if ret == -1 {
		return false
	}

	return true
}

func endDeviceListenThread() {
	C.end_device_listen_thread()
}

//export handleDeviceAdded
func handleDeviceAdded(deviceid C.int) {
	mouseList, tpadList, wacomList := libwrapper.GetDevicesList()

	libmouse.HandleDeviceChanged(mouseList)
	libtouchpad.HandleDeviceChanged(tpadList)
	libwacom.HandleDeviceChanged(wacomList)
}

//export handleDeviceRemoved
func handleDeviceRemoved(deviceid C.int) {
	mouseList, tpadList, wacomList := libwrapper.GetDevicesList()

	libmouse.HandleDeviceChanged(mouseList)
	libtouchpad.HandleDeviceChanged(tpadList)
	libwacom.HandleDeviceChanged(wacomList)
}
