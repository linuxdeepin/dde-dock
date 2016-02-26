/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package inputdevices

import (
	"gir/gio-2.0"
	"pkg.deepin.io/dde/api/dxinput"
	dxutils "pkg.deepin.io/dde/api/dxinput/utils"
	"pkg.deepin.io/lib/dbus/property"
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
	// TrackPoint
	mouseKeyTPMidButton        = "trackpoint-mid-button"
	mouseKeyTPMidButtonTimeout = "trackpoint-mid-button-timeout"
	mouseKeyTPWheel            = "trackpoint-wheel"
	mouseKeyTPWheelButton      = "trackpoint-wheel-button"
	mouseKeyTPWheelTimeout     = "trackpoint-wheel-timeout"
	mouseKeyTPWheelHorizScroll = "trackpoint-wheel-horiz-scroll"
	mouseKeyTPAcceleration     = "trackpoint-acceleration"
	mouseKeyTPThreshold        = "trackpoint-threshold"
	mouseKeyTPScaling          = "trackpoint-scaling"
)

type Mouse struct {
	LeftHanded            *property.GSettingsBoolProperty `access:"readwrite"`
	DisableTpad           *property.GSettingsBoolProperty `access:"readwrite"`
	NaturalScroll         *property.GSettingsBoolProperty `access:"readwrite"`
	MiddleButtonEmulation *property.GSettingsBoolProperty `access:"readwrite"`

	MotionAcceleration *property.GSettingsFloatProperty `access:"readwrite"`
	MotionThreshold    *property.GSettingsFloatProperty `access:"readwrite"`
	MotionScaling      *property.GSettingsFloatProperty `access:"readwrite"`

	DoubleClick   *property.GSettingsIntProperty `access:"readwrite"`
	DragThreshold *property.GSettingsIntProperty `access:"readwrite"`

	DeviceList dxutils.DeviceInfos
	Exist      bool

	dxMouses map[int32]*dxinput.Mouse
	setting  *gio.Settings
}

var _mouse *Mouse

func getMouse() *Mouse {
	if _mouse == nil {
		_mouse = NewMouse()

		_mouse.init()
		_mouse.handleGSettings()
	}

	return _mouse
}

func NewMouse() *Mouse {
	var m = new(Mouse)

	m.setting = gio.NewSettings(mouseSchema)
	m.LeftHanded = property.NewGSettingsBoolProperty(
		m, "LeftHanded",
		m.setting, mouseKeyLeftHanded)
	m.DisableTpad = property.NewGSettingsBoolProperty(
		m, "DisableTpad",
		m.setting, mouseKeyDisableTouchpad)
	m.NaturalScroll = property.NewGSettingsBoolProperty(
		m, "NaturalScroll",
		m.setting, mouseKeyNaturalScroll)
	m.MiddleButtonEmulation = property.NewGSettingsBoolProperty(
		m, "MiddleButtonEmulation",
		m.setting, mouseKeyMiddleButton)

	m.MotionAcceleration = property.NewGSettingsFloatProperty(
		m, "MotionAcceleration",
		m.setting, mouseKeyAcceleration)
	m.MotionThreshold = property.NewGSettingsFloatProperty(
		m, "MotionThreshold",
		m.setting, mouseKeyThreshold)
	m.MotionScaling = property.NewGSettingsFloatProperty(
		m, "MotionScaling",
		m.setting, mouseKeyScaling)

	m.DoubleClick = property.NewGSettingsIntProperty(
		m, "DoubleClick",
		m.setting, mouseKeyDoubleClick)
	m.DragThreshold = property.NewGSettingsIntProperty(
		m, "DragThreshold",
		m.setting, mouseKeyDragThreshold)

	m.updateDeviceList()
	m.dxMouses = make(map[int32]*dxinput.Mouse)
	m.updateDXMouses()

	return m
}

func (m *Mouse) init() {
	if !m.Exist {
		if getTouchpad().Exist {
			getTouchpad().enable(true)
		}
		return
	}

	m.enableLeftHanded()
	m.enableMidBtnEmu()
	m.enableNaturalScroll()
	m.motionAcceleration()
	m.motionThreshold()
	if m.DisableTpad.Get() {
		m.disableTouchpad()
	}
	m.setTrackPoint()
}

func (m *Mouse) handleDeviceChanged() {
	m.updateDeviceList()
	m.updateDXMouses()
	m.init()
}

func (m *Mouse) updateDeviceList() {
	m.DeviceList = getMouseInfos(false)
	if len(m.DeviceList) == 0 {
		m.setPropExist(false)
	} else {
		m.setPropExist(true)
	}
}

func (m *Mouse) updateDXMouses() {
	for _, info := range m.DeviceList {
		_, ok := m.dxMouses[info.Id]
		if ok {
			continue
		}

		dxm, err := dxinput.NewMouse(info.Id)
		if err != nil {
			logger.Warning(err)
			continue
		}
		m.dxMouses[info.Id] = dxm
	}
}

func (m *Mouse) disableTouchpad() {
	if !m.Exist {
		return
	}

	tpad := getTouchpad()
	if !tpad.Exist {
		return
	}

	if !m.DisableTpad.Get() && tpad.TPadEnable.Get() {
		tpad.enable(true)
		return
	}

	tpad.enable(false)
}

func (m *Mouse) enableLeftHanded() {
	enabled := m.LeftHanded.Get()
	for _, v := range m.dxMouses {
		err := v.EnableLeftHanded(enabled)
		if err != nil {
			logger.Debugf("Enable left handed for '%d - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (m *Mouse) enableNaturalScroll() {
	enabled := m.NaturalScroll.Get()
	for _, v := range m.dxMouses {
		err := v.EnableNaturalScroll(enabled)
		if err != nil {
			logger.Debugf("Enable natural scroll for '%d - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (m *Mouse) enableMidBtnEmu() {
	enabled := m.MiddleButtonEmulation.Get()
	for _, v := range m.dxMouses {
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

func (m *Mouse) motionAcceleration() {
	accel := float32(m.MotionAcceleration.Get())
	for _, v := range m.dxMouses {
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
	for _, v := range m.dxMouses {
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
	for _, v := range m.dxMouses {
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

func (m *Mouse) setTrackPoint() {
	middle := m.setting.GetBoolean(mouseKeyTPMidButton)
	middleTimeout := int16(m.setting.GetInt(mouseKeyTPMidButtonTimeout))
	wheel := m.setting.GetBoolean(mouseKeyTPWheel)
	wheelButton := int8(m.setting.GetInt(mouseKeyTPWheelButton))
	wheelTimeout := int16(m.setting.GetInt(mouseKeyTPWheelTimeout))
	wheelHorizScroll := m.setting.GetBoolean(mouseKeyTPWheelHorizScroll)
	accel := float32(m.setting.GetDouble(mouseKeyTPAcceleration))
	thres := float32(m.setting.GetDouble(mouseKeyTPThreshold))
	scaling := float32(m.setting.GetDouble(mouseKeyTPScaling))

	for _, v := range m.dxMouses {
		if !v.TrackPoint {
			continue
		}
		v.EnableMiddleButtonEmulation(middle)
		v.SetMiddleButtonEmulationTimeout(middleTimeout)
		v.EnableWheelEmulation(wheel)
		v.SetWheelEmulationButton(wheelButton)
		v.SetWheelEmulationTimeout(wheelTimeout)
		v.EnableWheelHorizScroll(wheelHorizScroll)
		v.SetMotionAcceleration(accel)
		v.SetMotionThreshold(thres)
		v.SetMotionScaling(scaling)
	}
}
