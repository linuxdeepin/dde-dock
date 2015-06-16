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

package mouse

import (
	"dbus/com/deepin/sessionmanager"
	"pkg.linuxdeepin.com/dde-daemon/inputdevices/touchpad"
	"pkg.linuxdeepin.com/dde-daemon/inputdevices/wrapper"
	"pkg.linuxdeepin.com/lib/dbus/property"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"pkg.linuxdeepin.com/lib/log"
)

const (
	mouseKeyLeftHanded      = "left-handed"
	mouseKeyDisableTouchpad = "disable-touchpad"
	mouseKeyMiddleButton    = "middle-button-enabled"
	mouseKeyNaturalScroll   = "natural-scroll"
	mouseKeyAcceleration    = "motion-acceleration"
	mouseKeyThreshold       = "motion-threshold"
	mouseKeyDoubleClick     = "double-click"
	mouseKeyDragThreshold   = "drag-threshold"
)

type Mouse struct {
	LeftHanded            *property.GSettingsBoolProperty `access:"readwrite"`
	DisableTpad           *property.GSettingsBoolProperty `access:"readwrite"`
	NaturalScroll         *property.GSettingsBoolProperty `access:"readwrite"`
	MiddleButtonEmulation *property.GSettingsBoolProperty `access:"readwrite"`

	MotionAcceleration *property.GSettingsFloatProperty `access:"readwrite"`
	MotionThreshold    *property.GSettingsFloatProperty `access:"readwrite"`

	DoubleClick   *property.GSettingsIntProperty `access:"readwrite"`
	DragThreshold *property.GSettingsIntProperty `access:"readwrite"`

	DeviceList []wrapper.XIDeviceInfo
	Exist      bool

	logger    *log.Logger
	settings  *gio.Settings
	xsettings *sessionmanager.XSettings
}

var _mouse *Mouse

func NewMouse(l *log.Logger) *Mouse {
	m := &Mouse{}

	m.settings = gio.NewSettings("com.deepin.dde.mouse")
	m.LeftHanded = property.NewGSettingsBoolProperty(
		m, "LeftHanded",
		m.settings, mouseKeyLeftHanded)
	m.DisableTpad = property.NewGSettingsBoolProperty(
		m, "DisableTpad",
		m.settings, mouseKeyDisableTouchpad)
	m.NaturalScroll = property.NewGSettingsBoolProperty(
		m, "NaturalScroll",
		m.settings, mouseKeyNaturalScroll)
	m.MiddleButtonEmulation = property.NewGSettingsBoolProperty(
		m, "MiddleButtonEmulation",
		m.settings, mouseKeyMiddleButton)

	m.MotionAcceleration = property.NewGSettingsFloatProperty(
		m, "MotionAcceleration",
		m.settings, mouseKeyAcceleration)
	m.MotionThreshold = property.NewGSettingsFloatProperty(
		m, "MotionThreshold",
		m.settings, mouseKeyThreshold)

	m.DoubleClick = property.NewGSettingsIntProperty(
		m, "DoubleClick",
		m.settings, mouseKeyDoubleClick)
	m.DragThreshold = property.NewGSettingsIntProperty(
		m, "DragThreshold",
		m.settings, mouseKeyDragThreshold)

	mouseList, _, _ := wrapper.GetDevicesList()
	m.setPropDeviceList(mouseList)
	if len(mouseList) > 0 {
		m.setPropExist(true)
	} else {
		m.setPropExist(false)
	}

	m.logger = l
	var err error
	m.xsettings, err = sessionmanager.NewXSettings(
		"com.deepin.SessionManager",
		"/com/deepin/XSettings",
	)
	if err != nil {
		m.warningInfo("Create XSettings Failed: %v", err)
		m.xsettings = nil
	}

	_mouse = m
	m.init()
	m.handleGSettings()

	return m
}

func HandleDeviceChanged(devList []wrapper.XIDeviceInfo) {
	if _mouse == nil {
		return
	}

	_mouse.setPropDeviceList(devList)
	if len(devList) == 0 {
		_mouse.setPropExist(false)
		if _mouse.DisableTpad.Get() == true {
			touchpad.DeviceEnabled(true)
		}
	} else {
		_mouse.setPropExist(true)
		_mouse.init()
	}
}

/**
 * TODO:
 *	HandleDeviceAdded
 *	HandleDeviceRemoved
 **/

func (m *Mouse) Reset() {
	for _, key := range m.settings.ListKeys() {
		m.settings.Reset(key)
	}
}

func (m *Mouse) motionAcceleration(accel float64) {
	for _, info := range m.DeviceList {
		wrapper.SetMotionAcceleration(info.Deviceid, accel)
	}
}

func (m *Mouse) motionThreshold(thres float64) {
	for _, info := range m.DeviceList {
		wrapper.SetMotionThreshold(info.Deviceid, thres)
	}
}

func (m *Mouse) leftHanded(enabled bool) {
	for _, info := range m.DeviceList {
		wrapper.SetLeftHanded(info.Deviceid, info.Name, enabled)
	}
}

func (m *Mouse) naturalScroll(enabled bool) {
	for _, info := range m.DeviceList {
		wrapper.SetMouseNaturalScroll(info.Deviceid, info.Name, enabled)
	}
}

func (m *Mouse) middleButtonEmulation(enabled bool) {
	for _, info := range m.DeviceList {
		wrapper.SetMiddleButtonEmulation(info.Deviceid, enabled)
	}
}

func (m *Mouse) doubleClick(value uint32) {
	if m.settings != nil {
		m.xsettings.SetInteger("Net/DoubleClickTime", value)
	}
}

func (m *Mouse) dragThreshold(value uint32) {
	if m.settings != nil {
		m.xsettings.SetInteger("Net/DndDragThreshold", value)
	}
}

func (m *Mouse) disableTouchpad(enabled bool) {
	if !enabled {
		touchpad.DeviceEnabled(true)
		return
	}

	if m.Exist {
		touchpad.DeviceEnabled(false)
	} else {
		touchpad.DeviceEnabled(true)
	}
}

// TODO: set by device info
func (m *Mouse) init() {
	m.disableTouchpad(m.DisableTpad.Get())

	if !m.Exist {
		return
	}

	m.leftHanded(m.LeftHanded.Get())
	m.middleButtonEmulation(m.MiddleButtonEmulation.Get())

	m.motionAcceleration(m.MotionAcceleration.Get())
	m.motionThreshold(m.MotionThreshold.Get())

	m.doubleClick(uint32(m.DoubleClick.Get()))
	m.dragThreshold(uint32(m.DragThreshold.Get()))
}
