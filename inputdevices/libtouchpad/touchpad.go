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

package libtouchpad

import (
	"dbus/com/deepin/sessionmanager"
	"os"
	"os/exec"
	"pkg.linuxdeepin.com/dde-daemon/inputdevices/libwrapper"
	"pkg.linuxdeepin.com/lib/dbus/property"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"pkg.linuxdeepin.com/lib/log"
)

const (
	tpadKeyEnabled       = "touchpad-enabled"
	tpadKeyLeftHanded    = "left-handed"
	tpadKeyWhileTyping   = "disable-while-typing"
	tpadKeyNaturalScroll = "natural-scroll"
	tpadKeyEdgeScroll    = "edge-scroll-enabled"
	tpadKeyHorizScroll   = "horiz-scroll-enabled"
	tpadKeyVertScroll    = "vert-scroll-enabled"
	tpadKeyAcceleration  = "motion-acceleration"
	tpadKeyThreshold     = "motion-threshold"
	tpadKeyTapClick      = "tap-to-click"
	tpadKeyScrollDelta   = "delta-scroll"
	tpadKeyDoubleClick   = "double-click"
	tpadKeyDragThreshold = "drag-threshold"
)

type Touchpad struct {
	TPadEnable      *property.GSettingsBoolProperty `access:"readwrite"`
	LeftHanded      *property.GSettingsBoolProperty `access:"readwrite"`
	DisableIfTyping *property.GSettingsBoolProperty `access:"readwrite"`
	NaturalScroll   *property.GSettingsBoolProperty `access:"readwrite"`
	EdgeScroll      *property.GSettingsBoolProperty `access:"readwrite"`
	HorizScroll     *property.GSettingsBoolProperty `access:"readwrite"`
	VertScroll      *property.GSettingsBoolProperty `access:"readwrite"`
	TapClick        *property.GSettingsBoolProperty `access:"readwrite"`

	MotionAcceleration *property.GSettingsFloatProperty `access:"readwrite"`
	MotionThreshold    *property.GSettingsFloatProperty `access:"readwrite"`

	DoubleClick   *property.GSettingsIntProperty `access:"readwrite"`
	DragThreshold *property.GSettingsIntProperty `access:"readwrite"`
	DeltaScroll   *property.GSettingsIntProperty `access:"readwrite"`

	Exist      bool
	DeviceList []libwrapper.XIDeviceInfo

	logger    *log.Logger
	process   *os.Process
	settings  *gio.Settings
	xsettings *sessionmanager.XSettings
}

var _tpad *Touchpad

func NewTouchpad(l *log.Logger) *Touchpad {
	touchpad := &Touchpad{}

	touchpad.settings = gio.NewSettings("com.deepin.dde.touchpad")
	touchpad.TPadEnable = property.NewGSettingsBoolProperty(
		touchpad, "TPadEnable",
		touchpad.settings, tpadKeyEnabled)
	touchpad.LeftHanded = property.NewGSettingsBoolProperty(
		touchpad, "LeftHanded",
		touchpad.settings, tpadKeyLeftHanded)
	touchpad.DisableIfTyping = property.NewGSettingsBoolProperty(
		touchpad, "DisableIfTyping",
		touchpad.settings, tpadKeyWhileTyping)
	touchpad.NaturalScroll = property.NewGSettingsBoolProperty(
		touchpad, "NaturalScroll",
		touchpad.settings, tpadKeyNaturalScroll)
	touchpad.EdgeScroll = property.NewGSettingsBoolProperty(
		touchpad, "EdgeScroll",
		touchpad.settings, tpadKeyEdgeScroll)
	touchpad.VertScroll = property.NewGSettingsBoolProperty(
		touchpad, "VertScroll",
		touchpad.settings, tpadKeyVertScroll)
	touchpad.HorizScroll = property.NewGSettingsBoolProperty(
		touchpad, "HorizScroll",
		touchpad.settings, tpadKeyHorizScroll)
	touchpad.TapClick = property.NewGSettingsBoolProperty(
		touchpad, "TapClick",
		touchpad.settings, tpadKeyTapClick)

	touchpad.MotionAcceleration = property.NewGSettingsFloatProperty(
		touchpad, "MotionAcceleration",
		touchpad.settings, tpadKeyAcceleration)
	touchpad.MotionThreshold = property.NewGSettingsFloatProperty(
		touchpad, "MotionThreshold",
		touchpad.settings, tpadKeyThreshold)

	touchpad.DeltaScroll = property.NewGSettingsIntProperty(
		touchpad, "DeltaScroll",
		touchpad.settings, tpadKeyScrollDelta)
	touchpad.DoubleClick = property.NewGSettingsIntProperty(
		touchpad, "DoubleClick",
		touchpad.settings, tpadKeyDoubleClick)
	touchpad.DragThreshold = property.NewGSettingsIntProperty(
		touchpad, "DragThreshold",
		touchpad.settings, tpadKeyDragThreshold)

	_, tpadList, _ := libwrapper.GetDevicesList()
	touchpad.setPropDeviceList(tpadList)

	if len(touchpad.DeviceList) > 0 {
		touchpad.setPropExist(true)
	} else {
		touchpad.setPropExist(false)
	}

	touchpad.logger = l
	var err error
	touchpad.xsettings, err = sessionmanager.NewXSettings(
		"com.deepin.SessionManager",
		"/com/deepin/XSettings",
	)
	if err != nil {
		touchpad.warningInfo("Create XSettings Failed: %v", err)
		touchpad.xsettings = nil
	}

	_tpad = touchpad
	touchpad.init()
	touchpad.handleGSettings()

	return touchpad
}

func DeviceEnabled(enable bool) {
	if _tpad == nil {
		return
	}

	if !enable {
		_tpad.enable(false)
		return
	}

	// enable == true
	if !_tpad.TPadEnable.Get() {
		return
	}
	_tpad.enable(true)
}

func HandleDeviceChanged(devList []libwrapper.XIDeviceInfo) {
	if _tpad == nil {
		return
	}

	_tpad.setPropDeviceList(devList)
	if len(devList) == 0 {
		_tpad.setPropExist(false)
	} else {
		_tpad.setPropExist(true)
		_tpad.init()
	}
}

// TODO
func HandleDeviceAdded(devInfo libwrapper.XIDeviceInfo) {
	if _tpad == nil {
		return
	}
}

// TODO
func HandleDeviceRemoved(devInfo libwrapper.XIDeviceInfo) {
	if _tpad == nil {
		return
	}
}

func (touchpad *Touchpad) Reset() {
	for _, key := range touchpad.settings.ListKeys() {
		touchpad.settings.Reset(key)
	}
}

func (touchpad *Touchpad) disableTpadWhileTyping(enable bool) {
	if !enable {
		if touchpad.process != nil {
			touchpad.process.Kill()
			touchpad.process = nil
		}
		return
	}

	// Has a syndaemon running...
	if touchpad.process != nil {
		return
	}
	cmd := exec.Command("/bin/sh", "-c",
		"syndaemon -i 1 -K -t")
	err := cmd.Start()
	if err != nil {
		touchpad.warningInfo("Exec syndaemon failed: %v", err)
		return
	}
	touchpad.process = cmd.Process
}

// TODO: set by deviceid
func (touchpad *Touchpad) enable(enabled bool) {
	for _, info := range touchpad.DeviceList {
		libwrapper.SetTouchpadEnabled(info.Deviceid, enabled)
	}
}

// TODO: set by deviceid
func (touchpad *Touchpad) leftHanded(enabled bool) {
	for _, info := range touchpad.DeviceList {
		libwrapper.SetLeftHanded(info.Deviceid, info.Name, enabled)
	}
}

// TODO: set by deviceid
func (touchpad *Touchpad) motionAcceleration(accel float64) {
	for _, info := range touchpad.DeviceList {
		libwrapper.SetMotionAcceleration(info.Deviceid, accel)
	}
}

// TODO: set by deviceid
func (touchpad *Touchpad) motionThreshold(thres float64) {
	for _, info := range touchpad.DeviceList {
		libwrapper.SetMotionThreshold(info.Deviceid, thres)
	}
}

// TODO: set by deviceid
func (touchpad *Touchpad) naturalScroll(enabled bool, delta int32) {
	for _, info := range touchpad.DeviceList {
		libwrapper.SetTouchpadNaturalScroll(info.Deviceid, enabled, delta)
	}
}

// TODO: set by deviceid
func (touchpad *Touchpad) edgeScroll(enabled bool) {
	for _, info := range touchpad.DeviceList {
		libwrapper.SetTouchpadEdgeScroll(info.Deviceid, enabled)
	}
}

// TODO: set by deviceid
func (touchpad *Touchpad) twoFingerScroll(vertEnabled,
	horizEnabled bool) {
	for _, info := range touchpad.DeviceList {
		libwrapper.SetTouchpadTwoFingerScroll(info.Deviceid,
			vertEnabled, horizEnabled)
	}
}

// TODO: set by deviceid
func (touchpad *Touchpad) tapToClick(enabled, leftHanded bool) {
	for _, info := range touchpad.DeviceList {
		libwrapper.SetTouchpadTapToClick(info.Deviceid,
			enabled, leftHanded)
	}
}

// TODO: set by deviceid
func (touchpad *Touchpad) doubleClick(value int32) {
	if touchpad.xsettings != nil {
		touchpad.xsettings.SetInteger("Net/DoubleClickTime",
			uint32(value))
	}
}

// TODO: set by deviceid
func (touchpad *Touchpad) dragThreshold(value int32) {
	if touchpad.xsettings != nil {
		touchpad.xsettings.SetInteger("Net/DndDragThreshold",
			uint32(value))
	}
}

// TODO: set by deviceid
func (touchpad *Touchpad) init() {
	if !touchpad.Exist {
		return
	}

	enabled := touchpad.TPadEnable.Get()
	touchpad.enable(enabled)
	if !enabled {
		return
	}

	touchpad.leftHanded(touchpad.LeftHanded.Get())
	touchpad.tapToClick(touchpad.TapClick.Get(),
		touchpad.LeftHanded.Get())
	if touchpad.DisableIfTyping.Get() {
		exec.Command("/bin/sh", "-c", "killall syndaemon").Run()
	}
	touchpad.disableTpadWhileTyping(touchpad.DisableIfTyping.Get())
	touchpad.naturalScroll(touchpad.NaturalScroll.Get(),
		touchpad.DeltaScroll.Get())
	touchpad.edgeScroll(touchpad.EdgeScroll.Get())
	touchpad.twoFingerScroll(touchpad.VertScroll.Get(),
		touchpad.HorizScroll.Get())

	touchpad.motionAcceleration(touchpad.MotionAcceleration.Get())
	touchpad.motionThreshold(touchpad.MotionThreshold.Get())
}
