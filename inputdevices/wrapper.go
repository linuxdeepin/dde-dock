/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package inputdevices

// #cgo pkg-config: x11 xi
// #cgo LDFLAGS: -lpthread
// #include "listen.h"
import "C"

import (
	dxutils "pkg.deepin.io/dde/api/dxinput/utils"
)

var (
	devInfos      dxutils.DeviceInfos
	mouseDevInfos dxutils.DeviceInfos
	tpadDevInfos  dxutils.DeviceInfos
	wacomDevInfos dxutils.DeviceInfos
)

func startDeviceListener() {
	C.start_device_listener()
}

func endDeviceListener() {
	C.end_device_listener()
}

//export handleDeviceChanged
func handleDeviceChanged() {
	logger.Debug("Device changed")

	getDeviceInfos(true)
	getMouseInfos(true)
	getTPadInfos(true)
	getWacomInfos(true)

	getTouchpad().handleDeviceChanged()
	getMouse().handleDeviceChanged()
	getWacom().handleDeviceChanged()
}

func getDeviceInfos(force bool) dxutils.DeviceInfos {
	if force || len(devInfos) == 0 {
		devInfos = dxutils.ListDevice()
	}

	return devInfos
}

func getMouseInfos(force bool) dxutils.DeviceInfos {
	if !force && len(mouseDevInfos) != 0 {
		return mouseDevInfos
	}

	mouseDevInfos = nil
	for _, info := range getDeviceInfos(false) {
		if info.Type == dxutils.DevTypeMouse {
			mouseDevInfos = append(mouseDevInfos, info)
		}
	}

	return mouseDevInfos
}

func getTPadInfos(force bool) dxutils.DeviceInfos {
	if !force && len(tpadDevInfos) != 0 {
		return tpadDevInfos
	}

	tpadDevInfos = nil
	for _, info := range getDeviceInfos(false) {
		if info.Type == dxutils.DevTypeTouchpad {
			tpadDevInfos = append(tpadDevInfos, info)
		}
	}

	return tpadDevInfos
}

func getWacomInfos(force bool) dxutils.DeviceInfos {
	if !force && len(wacomDevInfos) != 0 {
		return wacomDevInfos
	}

	wacomDevInfos = nil
	for _, info := range getDeviceInfos(false) {
		if info.Type == dxutils.DevTypeWacom {
			wacomDevInfos = append(wacomDevInfos, info)
		}
	}

	return wacomDevInfos
}
