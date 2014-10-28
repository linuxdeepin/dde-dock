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

package libmouse

import (
	"dbus/com/deepin/sessionmanager"
	"pkg.linuxdeepin.com/dde-daemon/inputdevices/libtouchpad"
	"pkg.linuxdeepin.com/dde-daemon/inputdevices/libwrapper"
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
	LeftHanded    *property.GSettingsBoolProperty `access:"readwrite"`
	DisableTpad   *property.GSettingsBoolProperty `access:"readwrite"`
	NaturalScroll *property.GSettingsBoolProperty `access:"readwrite"`

	MotionAcceleration *property.GSettingsFloatProperty `access:"readwrite"`
	MotionThreshold    *property.GSettingsFloatProperty `access:"readwrite"`

	DoubleClick   *property.GSettingsIntProperty `access:"readwrite"`
	DragThreshold *property.GSettingsIntProperty `access:"readwrite"`

	DeviceList []libwrapper.XIDeviceInfo
	Exist      bool

	logger    *log.Logger
	settings  *gio.Settings
	xsettings *sessionmanager.XSettings
}

var _mouse *Mouse

func NewMouse(l *log.Logger) *Mouse {
	mouse := &Mouse{}

	mouse.settings = gio.NewSettings("com.deepin.dde.mouse")
	mouse.LeftHanded = property.NewGSettingsBoolProperty(
		mouse, "LeftHanded",
		mouse.settings, mouseKeyLeftHanded)
	mouse.DisableTpad = property.NewGSettingsBoolProperty(
		mouse, "DisableTpad",
		mouse.settings, mouseKeyDisableTouchpad)
	mouse.NaturalScroll = property.NewGSettingsBoolProperty(
		mouse, "NaturalScroll",
		mouse.settings, mouseKeyNaturalScroll)

	mouse.MotionAcceleration = property.NewGSettingsFloatProperty(
		mouse, "MotionAcceleration",
		mouse.settings, mouseKeyAcceleration)
	mouse.MotionThreshold = property.NewGSettingsFloatProperty(
		mouse, "MotionThreshold",
		mouse.settings, mouseKeyThreshold)

	mouse.DoubleClick = property.NewGSettingsIntProperty(
		mouse, "DoubleClick",
		mouse.settings, mouseKeyDoubleClick)
	mouse.DragThreshold = property.NewGSettingsIntProperty(
		mouse, "DragThreshold",
		mouse.settings, mouseKeyDragThreshold)

	mouseList, _, _ := libwrapper.GetDevicesList()
	mouse.setPropDeviceList(mouseList)
	if len(mouseList) > 0 {
		mouse.setPropExist(true)
	} else {
		mouse.setPropExist(false)
	}

	mouse.logger = l
	var err error
	mouse.xsettings, err = sessionmanager.NewXSettings(
		"com.deepin.SessionManager",
		"/com/deepin/XSettings",
	)
	if err != nil {
		mouse.warningInfo("Create XSettings Failed: %v", err)
		mouse.xsettings = nil
	}

	_mouse = mouse
	mouse.init()
	mouse.handleGSettings()

	return mouse
}

func HandleDeviceChanged(devList []libwrapper.XIDeviceInfo) {
	if _mouse == nil {
		return
	}

	_mouse.setPropDeviceList(devList)
	if len(devList) == 0 {
		_mouse.setPropExist(false)
		if _mouse.DisableTpad.Get() == true {
			libtouchpad.DeviceEnabled(true)
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

func (mouse *Mouse) Reset() {
	for _, key := range mouse.settings.ListKeys() {
		mouse.settings.Reset(key)
	}
}

func (mouse *Mouse) motionAcceleration(accel float64) {
	for _, info := range mouse.DeviceList {
		libwrapper.SetMotionAcceleration(info.Deviceid, accel)
	}
}

func (mouse *Mouse) motionThreshold(thres float64) {
	for _, info := range mouse.DeviceList {
		libwrapper.SetMotionThreshold(info.Deviceid, thres)
	}
}

func (mouse *Mouse) leftHanded(enabled bool) {
	for _, info := range mouse.DeviceList {
		libwrapper.SetLeftHanded(info.Deviceid, info.Name, enabled)
	}
}

func (mouse *Mouse) naturalScroll(enabled bool) {
	for _, info := range mouse.DeviceList {
		libwrapper.SetMouseNaturalScroll(info.Deviceid, info.Name, enabled)
	}
}

func (mouse *Mouse) doubleClick(value uint32) {
	if mouse.settings != nil {
		mouse.xsettings.SetInteger("Net/DoubleClickTime", value)
	}
}

func (mouse *Mouse) dragThreshold(value uint32) {
	if mouse.settings != nil {
		mouse.xsettings.SetInteger("Net/DndDragThreshold", value)
	}
}

func (mouse *Mouse) disableTouchpad(enabled bool) {
	if enabled && mouse.Exist {
		libtouchpad.DeviceEnabled(false)
	} else {
		libtouchpad.DeviceEnabled(true)
	}
}

// TODO: set by device info
func (mouse *Mouse) init() {
	mouse.disableTouchpad(mouse.DisableTpad.Get())

	if !mouse.Exist {
		return
	}

	mouse.leftHanded(mouse.LeftHanded.Get())

	mouse.motionAcceleration(mouse.MotionAcceleration.Get())
	mouse.motionThreshold(mouse.MotionThreshold.Get())

	mouse.doubleClick(uint32(mouse.DoubleClick.Get()))
	mouse.dragThreshold(uint32(mouse.DragThreshold.Get()))
}
