/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

import (
	"sync"

	"pkg.deepin.io/gir/gio-2.0"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/gsprop"
)

const (
	mouseSchema = "com.deepin.dde.mouse"

	mouseKeyLeftHanded      = "left-handed"
	mouseKeyDisableTouchpad = "disable-touchpad"
	mouseKeyMiddleButton    = "middle-button-enabled"
	mouseKeyNaturalScroll   = "natural-scroll"
	mouseKeyAcceleration    = "motion-acceleration"
	mouseKeyThreshold       = "motion-threshold"
	mouseKeyScaling         = "motion-scaling"
	mouseKeyDoubleClick     = "double-click"
	mouseKeyDragThreshold   = "drag-threshold"
	mouseKeyAdaptiveAccel   = "adaptive-accel-profile"
)

type Mouse struct {
	service    *dbusutil.Service
	PropsMu    sync.RWMutex
	DeviceList string
	Exist      bool

	// dbusutil-gen: ignore-below
	LeftHanded            gsprop.Bool `prop:"access:rw"`
	DisableTpad           gsprop.Bool `prop:"access:rw"`
	NaturalScroll         gsprop.Bool `prop:"access:rw"`
	MiddleButtonEmulation gsprop.Bool `prop:"access:rw"`
	AdaptiveAccelProfile  gsprop.Bool `prop:"access:rw"`

	MotionAcceleration gsprop.Double `prop:"access:rw"`
	MotionThreshold    gsprop.Double `prop:"access:rw"`
	MotionScaling      gsprop.Double `prop:"access:rw"`

	DoubleClick   gsprop.Int `prop:"access:rw"`
	DragThreshold gsprop.Int `prop:"access:rw"`

	devInfos dxMouses
	setting  *gio.Settings
	touchPad *Touchpad
}

func newMouse(service *dbusutil.Service, touchPad *Touchpad) *Mouse {
	var m = new(Mouse)

	m.service = service
	m.touchPad = touchPad
	m.setting = gio.NewSettings(mouseSchema)
	m.LeftHanded.Bind(m.setting, mouseKeyLeftHanded)
	m.DisableTpad.Bind(m.setting, mouseKeyDisableTouchpad)
	m.NaturalScroll.Bind(m.setting, mouseKeyNaturalScroll)
	m.MiddleButtonEmulation.Bind(m.setting, mouseKeyMiddleButton)
	m.MotionAcceleration.Bind(m.setting, mouseKeyAcceleration)
	m.MotionThreshold.Bind(m.setting, mouseKeyThreshold)
	m.MotionScaling.Bind(m.setting, mouseKeyScaling)
	m.DoubleClick.Bind(m.setting, mouseKeyDoubleClick)
	m.DragThreshold.Bind(m.setting, mouseKeyDragThreshold)
	m.AdaptiveAccelProfile.Bind(m.setting, mouseKeyAdaptiveAccel)

	m.updateDXMouses()

	return m
}

func (m *Mouse) init() {
	if !m.Exist {
		tpad := m.touchPad
		if tpad.Exist && !tpad.TPadEnable.Get() {
			tpad.TPadEnable.Set(true)
		}
		return
	}

	m.enableLeftHanded()
	m.enableMidBtnEmu()
	m.enableNaturalScroll()
	m.enableAdaptiveAccelProfile()
	m.motionAcceleration()
	m.motionThreshold()
	if m.DisableTpad.Get() {
		m.disableTouchPad()
	}
}

func (m *Mouse) handleDeviceChanged() {
	m.updateDXMouses()
	m.init()
}

func (m *Mouse) updateDXMouses() {
	m.devInfos = dxMouses{}
	for _, info := range getMouseInfos(false) {
		if info.TrackPoint {
			continue
		}

		tmp := m.devInfos.get(info.Id)
		if tmp != nil {
			continue
		}
		m.devInfos = append(m.devInfos, info)
	}

	m.PropsMu.Lock()
	var v string
	if len(m.devInfos) == 0 {
		m.setPropExist(false)
	} else {
		m.setPropExist(true)
		v = m.devInfos.string()
	}
	m.setPropDeviceList(v)
	m.PropsMu.Unlock()
}

func (m *Mouse) disableTouchPad() {
	m.PropsMu.RLock()
	mouseExist := m.Exist
	m.PropsMu.RUnlock()
	if !mouseExist {
		return
	}

	touchPad := m.touchPad
	touchPad.PropsMu.RLock()
	touchPadExist := touchPad.Exist
	touchPad.PropsMu.RUnlock()
	if !touchPadExist {
		return
	}

	if !m.DisableTpad.Get() && touchPad.TPadEnable.Get() {
		touchPad.enable(true)
		return
	}

	touchPad.enable(false)
}

func (m *Mouse) enableLeftHanded() {
	enabled := m.LeftHanded.Get()
	for _, v := range m.devInfos {
		err := v.EnableLeftHanded(enabled)
		if err != nil {
			logger.Debugf("Enable left handed for '%d - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
	setWMMouseBoolKey(wmTPadKeyLeftHanded, enabled)
}

func (m *Mouse) enableNaturalScroll() {
	enabled := m.NaturalScroll.Get()
	for _, v := range m.devInfos {
		err := v.EnableNaturalScroll(enabled)
		if err != nil {
			logger.Debugf("Enable natural scroll for '%d - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
	setWMMouseBoolKey(wmTPadKeyNaturalScroll, enabled)
}

func (m *Mouse) enableMidBtnEmu() {
	enabled := m.MiddleButtonEmulation.Get()
	for _, v := range m.devInfos {
		if v.TrackPoint {
			continue
		}

		err := v.EnableMiddleButtonEmulation(enabled)
		if err != nil {
			logger.Debugf("Enable mid btn emulation for '%d - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (m *Mouse) enableAdaptiveAccelProfile() {
	enabled := m.AdaptiveAccelProfile.Get()
	for _, v := range m.devInfos {
		if !v.CanChangeAccelProfile() {
			continue
		}

		err := v.SetUseAdaptiveAccelProfile(enabled)
		if err != nil {
			logger.Debugf("Enable adaptive accel profile for '%d - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (m *Mouse) motionAcceleration() {
	accel := float32(m.MotionAcceleration.Get())
	for _, v := range m.devInfos {
		if v.TrackPoint {
			continue
		}

		err := v.SetMotionAcceleration(accel)
		if err != nil {
			logger.Debugf("Set acceleration for '%d - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (m *Mouse) motionThreshold() {
	thres := float32(m.MotionThreshold.Get())
	for _, v := range m.devInfos {
		if v.TrackPoint {
			continue
		}

		err := v.SetMotionThreshold(thres)
		if err != nil {
			logger.Debugf("Set threshold for '%d - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (m *Mouse) motionScaling() {
	scaling := float32(m.MotionScaling.Get())
	for _, v := range m.devInfos {
		if v.TrackPoint {
			continue
		}

		err := v.SetMotionScaling(scaling)
		if err != nil {
			logger.Debugf("Set scaling for '%d - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (m *Mouse) doubleClick() {
	xsSetInt32(xsPropDoubleClick, m.DoubleClick.Get())
}

func (m *Mouse) dragThreshold() {
	xsSetInt32(xsPropDragThres, m.DragThreshold.Get())
}
