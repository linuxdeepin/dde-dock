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
	mouseInfos = dxMouses{}
	getMouseInfos(false)
	tpadInfos = dxTouchpads{}
	getTPadInfos(false)
	wacomInfos = dxWacoms{}
	getWacomInfos(false)

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

func isTrackPoint(info *dxutils.DeviceInfo) bool {
	name := strings.ToLower(info.Name)
	return strings.Contains(name, "trackpoint") ||
		strings.Contains(name, "dualpoint stick")
}

func getMouseInfos(force bool) dxMouses {
	if !force && len(mouseInfos) != 0 {
		return mouseInfos
	}

	mouseInfos = dxMouses{}
	for _, info := range getDeviceInfos(force) {
		if info.Type == dxutils.DevTypeMouse {
			mouseInfos = append(mouseInfos, &dxinput.Mouse{
				Id:         info.Id,
				Name:       info.Name,
				TrackPoint: isTrackPoint(info),
			})
		}
	}

	return mouseInfos
}

func getTPadInfos(force bool) dxTouchpads {
	if !force && len(tpadInfos) != 0 {
		return tpadInfos
	}

	tpadInfos = dxTouchpads{}
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

	wacomInfos = dxWacoms{}
	for _, info := range getDeviceInfos(false) {
		if info.Type == dxutils.DevTypeWacom {
			wacomInfo := &dxinput.Wacom{
				Id:   info.Id,
				Name: info.Name,
			}
			wacomInfos = append(wacomInfos, wacomInfo)
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
