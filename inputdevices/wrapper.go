/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
	getKeyboard().handleDeviceChanged()
}

func getDeviceInfos(force bool) dxutils.DeviceInfos {
	if force || len(devInfos) == 0 {
		devInfos = dxutils.ListDevice()
	}

	return devInfos
}

func getKeyboardNumber() int {
	var number = 0
	for _, info := range getDeviceInfos(false) {
		// TODO: Improve keyboard device detected by udev property 'ID_INPUT_KEYBOARD'
		if strings.Contains(strings.ToLower(info.Name), "keyboard") {
			number += 1
		}
	}
	return number
}

func getMouseInfos(force bool) dxMouses {
	if !force && len(mouseInfos) != 0 {
		return mouseInfos
	}

	mouseInfos = dxMouses{}
	for _, info := range getDeviceInfos(force) {
		if info.Type == dxutils.DevTypeMouse {
			tmp, _ := dxinput.NewMouseFromDeviceInfo(info)
			mouseInfos = append(mouseInfos, tmp)
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
			tmp, _ := dxinput.NewTouchpadFromDevInfo(info)
			tpadInfos = append(tpadInfos, tmp)
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
			tmp, _ := dxinput.NewWacomFromDevInfo(info)
			wacomInfos = append(wacomInfos, tmp)
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
