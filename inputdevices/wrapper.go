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
	"encoding/json"
	"pkg.deepin.io/dde/api/dxinput"
	dxutils "pkg.deepin.io/dde/api/dxinput/utils"
	"strings"
)

type dxMouses []*dxinput.Mouse
type dxTouchpads []*dxinput.Touchpad
type dxWacoms []*dxinput.Wacom

var (
	devInfos   dxutils.DeviceInfos
	mouseInfos dxMouses
	tpadInfos  dxTouchpads
	wacomInfos dxWacoms
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

func getMouseInfos(force bool) dxMouses {
	if !force && len(mouseInfos) != 0 {
		return mouseInfos
	}

	mouseInfos = nil
	for _, info := range getDeviceInfos(force) {
		if info.Type == dxutils.DevTypeMouse {
			mouseInfos = append(mouseInfos, &dxinput.Mouse{
				Id:   info.Id,
				Name: info.Name,
				TrackPoint: strings.Contains(
					strings.ToLower(info.Name),
					"trackpoint"),
			})
		}
	}

	return mouseInfos
}

func getTPadInfos(force bool) dxTouchpads {
	if !force && len(tpadInfos) != 0 {
		return tpadInfos
	}

	tpadInfos = nil
	for _, info := range getDeviceInfos(false) {
		if info.Type == dxutils.DevTypeTouchpad {
			tpadInfos = append(tpadInfos, &dxinput.Touchpad{
				Id:   info.Id,
				Name: info.Name,
			})
		}
	}

	return tpadInfos
}

func getWacomInfos(force bool) dxWacoms {
	if !force && len(wacomInfos) != 0 {
		return wacomInfos
	}

	wacomInfos = nil
	for _, info := range getDeviceInfos(false) {
		if info.Type == dxutils.DevTypeWacom {
			wacomInfos = append(wacomInfos, &dxinput.Wacom{
				Id:   info.Id,
				Name: info.Name,
			})
		}
	}

	return wacomInfos
}

func (infos dxMouses) get(id int32) *dxinput.Mouse {
	for _, info := range infos {
		if info.Id == id {
			return info
		}
	}
	return nil
}

func (infos dxMouses) string() string {
	return toJSON(infos)
}

func (infos dxTouchpads) get(id int32) *dxinput.Touchpad {
	for _, info := range infos {
		if info.Id == id {
			return info
		}
	}
	return nil
}

func (infos dxTouchpads) string() string {
	return toJSON(infos)
}

func (infos dxWacoms) get(id int32) *dxinput.Wacom {
	for _, info := range infos {
		if info.Id == id {
			return info
		}
	}
	return nil
}

func (infos dxWacoms) string() string {
	return toJSON(infos)
}

func toJSON(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}
