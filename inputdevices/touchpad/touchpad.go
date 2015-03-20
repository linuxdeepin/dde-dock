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

package touchpad

import (
	"dbus/com/deepin/sessionmanager"
	"os"
	"os/exec"
	"pkg.linuxdeepin.com/dde-daemon/inputdevices/wrapper"
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
	DeviceList []wrapper.XIDeviceInfo

	logger    *log.Logger
	process   *os.Process
	settings  *gio.Settings
	xsettings *sessionmanager.XSettings
}

var _tpad *Touchpad

func NewTouchpad(l *log.Logger) *Touchpad {
	tpad := &Touchpad{}

	tpad.settings = gio.NewSettings("com.deepin.dde.touchpad")
	tpad.TPadEnable = property.NewGSettingsBoolProperty(
		tpad, "TPadEnable",
		tpad.settings, tpadKeyEnabled)
	tpad.LeftHanded = property.NewGSettingsBoolProperty(
		tpad, "LeftHanded",
		tpad.settings, tpadKeyLeftHanded)
	tpad.DisableIfTyping = property.NewGSettingsBoolProperty(
		tpad, "DisableIfTyping",
		tpad.settings, tpadKeyWhileTyping)
	tpad.NaturalScroll = property.NewGSettingsBoolProperty(
		tpad, "NaturalScroll",
		tpad.settings, tpadKeyNaturalScroll)
	tpad.EdgeScroll = property.NewGSettingsBoolProperty(
		tpad, "EdgeScroll",
		tpad.settings, tpadKeyEdgeScroll)
	tpad.VertScroll = property.NewGSettingsBoolProperty(
		tpad, "VertScroll",
		tpad.settings, tpadKeyVertScroll)
	tpad.HorizScroll = property.NewGSettingsBoolProperty(
		tpad, "HorizScroll",
		tpad.settings, tpadKeyHorizScroll)
	tpad.TapClick = property.NewGSettingsBoolProperty(
		tpad, "TapClick",
		tpad.settings, tpadKeyTapClick)

	tpad.MotionAcceleration = property.NewGSettingsFloatProperty(
		tpad, "MotionAcceleration",
		tpad.settings, tpadKeyAcceleration)
	tpad.MotionThreshold = property.NewGSettingsFloatProperty(
		tpad, "MotionThreshold",
		tpad.settings, tpadKeyThreshold)

	tpad.DeltaScroll = property.NewGSettingsIntProperty(
		tpad, "DeltaScroll",
		tpad.settings, tpadKeyScrollDelta)
	tpad.DoubleClick = property.NewGSettingsIntProperty(
		tpad, "DoubleClick",
		tpad.settings, tpadKeyDoubleClick)
	tpad.DragThreshold = property.NewGSettingsIntProperty(
		tpad, "DragThreshold",
		tpad.settings, tpadKeyDragThreshold)

	_, tpadList, _ := wrapper.GetDevicesList()
	tpad.setPropDeviceList(tpadList)

	if len(tpad.DeviceList) > 0 {
		tpad.setPropExist(true)
	} else {
		tpad.setPropExist(false)
	}

	tpad.logger = l
	var err error
	tpad.xsettings, err = sessionmanager.NewXSettings(
		"com.deepin.SessionManager",
		"/com/deepin/XSettings",
	)
	if err != nil {
		tpad.warningInfo("Create XSettings Failed: %v", err)
		tpad.xsettings = nil
	}

	_tpad = tpad
	tpad.init()
	tpad.handleGSettings()

	return tpad
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

func HandleDeviceChanged(devList []wrapper.XIDeviceInfo) {
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
func HandleDeviceAdded(devInfo wrapper.XIDeviceInfo) {
	if _tpad == nil {
		return
	}
}

// TODO
func HandleDeviceRemoved(devInfo wrapper.XIDeviceInfo) {
	if _tpad == nil {
		return
	}
}

func (tpad *Touchpad) Reset() {
	for _, key := range tpad.settings.ListKeys() {
		tpad.settings.Reset(key)
	}
}

func (tpad *Touchpad) disableTpadWhileTyping(enable bool) {
	if !enable {
		if tpad.process != nil {
			tpad.process.Kill()
			tpad.process = nil
		}
		return
	}

	// Has a syndaemon running...
	if tpad.process != nil {
		return
	}
	cmd := exec.Command("/bin/sh", "-c",
		"syndaemon -i 1 -K -t")
	err := cmd.Start()
	if err != nil {
		tpad.warningInfo("Exec syndaemon failed: %v", err)
		return
	}
	tpad.process = cmd.Process
}

// TODO: set by deviceid
func (tpad *Touchpad) enable(enabled bool) {
	for _, info := range tpad.DeviceList {
		wrapper.SetTouchpadEnabled(info.Deviceid, enabled)
	}
}

// TODO: set by deviceid
func (tpad *Touchpad) leftHanded(enabled bool) {
	for _, info := range tpad.DeviceList {
		wrapper.SetLeftHanded(info.Deviceid, info.Name, enabled)
	}
}

// TODO: set by deviceid
func (tpad *Touchpad) motionAcceleration(accel float64) {
	for _, info := range tpad.DeviceList {
		wrapper.SetMotionAcceleration(info.Deviceid, accel)
	}
}

// TODO: set by deviceid
func (tpad *Touchpad) motionThreshold(thres float64) {
	for _, info := range tpad.DeviceList {
		wrapper.SetMotionThreshold(info.Deviceid, thres)
	}
}

// TODO: set by deviceid
func (tpad *Touchpad) naturalScroll(enabled bool, delta int32) {
	for _, info := range tpad.DeviceList {
		wrapper.SetTouchpadNaturalScroll(info.Deviceid, enabled, delta)
	}
}

// TODO: set by deviceid
func (tpad *Touchpad) edgeScroll(enabled bool) {
	for _, info := range tpad.DeviceList {
		wrapper.SetTouchpadEdgeScroll(info.Deviceid, enabled)
	}
}

// TODO: set by deviceid
func (tpad *Touchpad) twoFingerScroll(vertEnabled,
	horizEnabled bool) {
	for _, info := range tpad.DeviceList {
		wrapper.SetTouchpadTwoFingerScroll(info.Deviceid,
			vertEnabled, horizEnabled)
	}
}

// TODO: set by deviceid
func (tpad *Touchpad) tapToClick(enabled, leftHanded bool) {
	for _, info := range tpad.DeviceList {
		wrapper.SetTouchpadTapToClick(info.Deviceid,
			enabled, leftHanded)
	}
}

// TODO: set by deviceid
func (tpad *Touchpad) doubleClick(value int32) {
	if tpad.xsettings != nil {
		tpad.xsettings.SetInteger("Net/DoubleClickTime",
			uint32(value))
	}
}

// TODO: set by deviceid
func (tpad *Touchpad) dragThreshold(value int32) {
	if tpad.xsettings != nil {
		tpad.xsettings.SetInteger("Net/DndDragThreshold",
			uint32(value))
	}
}

// TODO: set by deviceid
func (tpad *Touchpad) init() {
	if !tpad.Exist {
		return
	}

	enabled := tpad.TPadEnable.Get()
	tpad.enable(enabled)
	if !enabled {
		return
	}

	tpad.leftHanded(tpad.LeftHanded.Get())
	tpad.tapToClick(tpad.TapClick.Get(),
		tpad.LeftHanded.Get())
	if tpad.DisableIfTyping.Get() {
		exec.Command("/bin/sh", "-c", "killall syndaemon").Run()
	}
	tpad.disableTpadWhileTyping(tpad.DisableIfTyping.Get())
	tpad.naturalScroll(tpad.NaturalScroll.Get(),
		tpad.DeltaScroll.Get())
	tpad.edgeScroll(tpad.EdgeScroll.Get())
	tpad.twoFingerScroll(tpad.VertScroll.Get(),
		tpad.HorizScroll.Get())

	tpad.motionAcceleration(tpad.MotionAcceleration.Get())
	tpad.motionThreshold(tpad.MotionThreshold.Get())
}
