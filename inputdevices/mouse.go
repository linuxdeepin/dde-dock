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

import (
	"pkg.linuxdeepin.com/lib/dbus/property"
	"pkg.linuxdeepin.com/lib/gio-2.0"
)

const (
	MOUSE_KEY_LEFT_HAND      = "left-handed"
	MOUSE_KEY_DISABLE_TPAD   = "disable-touchpad"
	MOUSE_KEY_MID_BUTTON     = "middle-button-enabled"
	MOUSE_KEY_NATURAL_SCROLL = "natural-scroll"
	MOUSE_KEY_ACCEL          = "motion-acceleration"
	MOUSE_KEY_THRES          = "motion-threshold"
	MOUSE_KEY_DOUBLE_CLICK   = "double-click"
	MOUSE_KEY_DRAG_THRES     = "drag-threshold"
)

var mouseSettings = gio.NewSettings("com.deepin.dde.mouse")

type MouseManager struct {
	LeftHanded    *property.GSettingsBoolProperty `access:"readwrite"`
	DisableTpad   *property.GSettingsBoolProperty `access:"readwrite"`
	NaturalScroll *property.GSettingsBoolProperty `access:"readwrite"`

	MotionAcceleration *property.GSettingsFloatProperty `access:"readwrite"`
	MotionThreshold    *property.GSettingsFloatProperty `access:"readwrite"`

	DoubleClick   *property.GSettingsIntProperty `access:"readwrite"`
	DragThreshold *property.GSettingsIntProperty `access:"readwrite"`

	DeviceList []PointerDeviceInfo
	Exist      bool

	listenFlag bool
}

var _mManager *MouseManager

func GetMouseManager() *MouseManager {
	if _mManager == nil {
		_mManager = newMouseManager()
	}

	return _mManager
}

func newMouseManager() *MouseManager {
	mManager := &MouseManager{}

	mManager.LeftHanded = property.NewGSettingsBoolProperty(
		mManager, "LeftHanded",
		mouseSettings, MOUSE_KEY_LEFT_HAND)
	mManager.DisableTpad = property.NewGSettingsBoolProperty(
		mManager, "DisableTpad",
		mouseSettings, MOUSE_KEY_DISABLE_TPAD)
	mManager.NaturalScroll = property.NewGSettingsBoolProperty(
		mManager, "NaturalScroll",
		mouseSettings, MOUSE_KEY_NATURAL_SCROLL)

	mManager.MotionAcceleration = property.NewGSettingsFloatProperty(
		mManager, "MotionAcceleration",
		mouseSettings, MOUSE_KEY_ACCEL)
	mManager.MotionThreshold = property.NewGSettingsFloatProperty(
		mManager, "MotionThreshold",
		mouseSettings, MOUSE_KEY_THRES)

	mManager.DoubleClick = property.NewGSettingsIntProperty(
		mManager, "DoubleClick",
		mouseSettings, MOUSE_KEY_DOUBLE_CLICK)
	mManager.DragThreshold = property.NewGSettingsIntProperty(
		mManager, "DragThreshold",
		mouseSettings, MOUSE_KEY_DRAG_THRES)

	mouseList, _, _ := getPointerDeviceList()
	mManager.setPropDeviceList(mouseList)
	if len(mouseList) > 0 {
		mManager.setPropExist(true)
	} else {
		mManager.setPropExist(false)
	}

	mManager.listenFlag = false

	mManager.init()

	return mManager
}

func (mManager *MouseManager) motionAcceleration(accel float64) {
	for _, info := range mManager.DeviceList {
		setMotionAcceleration(info.Deviceid, accel)
	}
}

func (mManager *MouseManager) motionThreshold(thres float64) {
	for _, info := range mManager.DeviceList {
		setMotionThreshold(info.Deviceid, thres)
	}
}

func (mManager *MouseManager) leftHanded(enabled bool) {
	for _, info := range mManager.DeviceList {
		setLeftHanded(info.Deviceid, info.Name, enabled)
	}
}

func (mManager *MouseManager) naturalScroll(enabled bool) {
	for _, info := range mManager.DeviceList {
		setMouseNaturalScroll(info.Deviceid, info.Name, enabled)
	}
}

func (mManager *MouseManager) doubleClick(value uint32) {
	if xsObj != nil {
		xsObj.SetInterger("Net/DoubleClickTime", value)
	}
}

func (mManager *MouseManager) dragThreshold(value uint32) {
	if xsObj != nil {
		xsObj.SetInterger("Net/DndDragThreshold", value)
	}
}

func (mManager *MouseManager) disableTouchpad(enabled bool) {
	if enabled && mManager.Exist {
		tpadSettings.SetBoolean(TPAD_KEY_ENABLE, false)
	} else {
		tpadSettings.SetBoolean(TPAD_KEY_ENABLE, true)
	}
}

func (mManager *MouseManager) init() {
	logger.Debug("Set disable tpad")
	mManager.disableTouchpad(mManager.DisableTpad.Get())

	if !mManager.Exist {
		return
	}

	logger.Debug("Set left handed")
	mManager.leftHanded(mManager.LeftHanded.Get())

	logger.Debug("Set acceleration")
	mManager.motionAcceleration(mManager.MotionAcceleration.Get())
	logger.Debug("Set threshold")
	mManager.motionThreshold(mManager.MotionThreshold.Get())

	if !mManager.listenFlag {
		mManager.listenGSettings()

		mManager.doubleClick(uint32(mManager.DoubleClick.Get()))
		mManager.dragThreshold(uint32(mManager.DragThreshold.Get()))
	}
}
