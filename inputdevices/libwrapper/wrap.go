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

package libwrapper

// #cgo pkg-config: x11 xi
// #cgo CFLAGS: -Wall -g
// #cgo LDFLAGS: -lpthread
// #include <stdlib.h>
// #include "devices.h"
import "C"

import (
	"unsafe"
)

type XIDeviceInfo struct {
	Name     string
	Deviceid int32
	Enabled  bool
}

func GetDevicesList() ([]XIDeviceInfo, []XIDeviceInfo, []XIDeviceInfo) {
	n_devices := C.int(0)
	devices := C.get_device_info_list(&n_devices)
	if n_devices < 1 {
		return nil, nil, nil
	}

	var (
		mouseList []XIDeviceInfo
		tpadList  []XIDeviceInfo
		wacomList []XIDeviceInfo
	)

	tmpList := uintptr(unsafe.Pointer(devices))
	length := unsafe.Sizeof(*devices)
	for i := C.int(0); i < n_devices; i++ {
		devInfo := (*C.DeviceInfo)(unsafe.Pointer(tmpList + uintptr(i)*length))
		if C.is_mouse_device(devInfo.deviceid) == 1 {
			info := XIDeviceInfo{
				C.GoString(devInfo.name),
				int32(devInfo.deviceid),
				false,
			}

			if devInfo.enabled == 1 {
				info.Enabled = true
			}

			mouseList = append(mouseList, info)
		} else if C.is_tpad_device(devInfo.deviceid) == 1 {
			info := XIDeviceInfo{
				C.GoString(devInfo.name),
				int32(devInfo.deviceid),
				false,
			}

			if devInfo.enabled == 1 {
				info.Enabled = true
			}

			tpadList = append(tpadList, info)
		} else if C.is_wacom_device(devInfo.deviceid) == 1 {
			info := XIDeviceInfo{
				C.GoString(devInfo.name),
				int32(devInfo.deviceid),
				false,
			}

			if devInfo.enabled == 1 {
				info.Enabled = true
			}

			wacomList = append(wacomList, info)
		}
	}
	C.free_device_info(devices, n_devices)

	return mouseList, tpadList, wacomList
}

/**
 * Set device pointer move speed
 * range: 0 < accel <= 2
 * defalut value: 1
 **/
func SetMotionAcceleration(deviceid int32, acceleration float64) bool {
	if acceleration <= 0 {
		return false
	}

	ret := C.set_motion_acceleration(C.int(deviceid),
		C.double(acceleration))
	if ret == -1 {
		return false
	}

	return true
}

/**
 * Set device pointer drag threshold
 * range: 1 =< accel <= 5
 * defalut value: 1
 **/
func SetMotionThreshold(deviceid int32, threshold float64) bool {
	if threshold < 1 {
		return false
	}

	ret := C.set_motion_threshold(C.int(deviceid),
		C.double(threshold))
	if ret == -1 {
		return false
	}

	return true
}

func SetLeftHanded(deviceid int32, name string, enabled bool) bool {
	cName := C.CString(name)
	leftHanded := C.int(0)
	if enabled {
		leftHanded = 1
	}

	ret := C.set_left_handed(C.ulong(deviceid), cName, leftHanded)
	C.free(unsafe.Pointer(cName))
	if ret == -1 {
		return false
	}

	return true
}

func SetMouseNaturalScroll(deviceid int32, name string, enabled bool) bool {
	cName := C.CString(name)
	scroll := C.int(0)
	if enabled {
		scroll = 1
	}

	ret := C.set_mouse_natural_scroll(C.ulong(deviceid), cName, scroll)
	C.free(unsafe.Pointer(cName))
	if ret == -1 {
		return false
	}

	return true
}

func SetTouchpadEnabled(deviceid int32, enabled bool) bool {
	tmp := C.int(0)
	if enabled {
		tmp = 1
	}

	ret := C.enabled_touchpad(C.int(deviceid), tmp)
	if ret == -1 {
		return false
	}

	return true
}

func SetTouchpadNaturalScroll(deviceid int32, enabled bool, delta int32) bool {
	scroll := C.int(0)
	if enabled {
		scroll = 1
	}

	ret := C.set_touchpad_natural_scroll(C.int(deviceid), scroll,
		C.int(delta))
	if ret == -1 {
		return false
	}

	return true
}

func SetTouchpadEdgeScroll(deviceid int32, enabled bool) bool {
	scroll := C.int(0)
	if enabled {
		scroll = 1
	}

	ret := C.set_edge_scroll(C.int(deviceid), scroll)
	if ret == -1 {
		return false
	}

	return true
}

func SetTouchpadTwoFingerScroll(deviceid int32,
	vert_enabled, horiz_enabled bool) bool {
	vert := C.int(0)
	if vert_enabled {
		vert = 1
	}

	horiz := C.int(0)
	if horiz_enabled {
		horiz = 1
	}

	ret := C.set_two_finger_scroll(C.int(deviceid), vert, horiz)
	if ret == -1 {
		return false
	}

	return true
}

func SetTouchpadTapToClick(deviceid int32, enabled, leftHanded bool) bool {
	tapEnable := C.int(0)
	if enabled {
		tapEnable = 1
	}

	left := C.int(0)
	if leftHanded {
		left = 1
	}

	ret := C.set_tab_to_click(C.int(deviceid), tapEnable, left)
	if ret == -1 {
		return false
	}

	return true
}

func SetKeyboardRepeat(enabled bool, delay, interval uint32) bool {
	repeat := C.int(0)
	if enabled {
		repeat = 1
	}

	ret := C.set_keyboard_repeat(repeat, C.uint(delay), C.uint(interval))
	if ret == -1 {
		return false
	}

	return true
}
