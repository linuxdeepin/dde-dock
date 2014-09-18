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
	"os/exec"
	"pkg.linuxdeepin.com/lib/dbus/property"
	"pkg.linuxdeepin.com/lib/gio-2.0"
)

const (
	TPAD_KEY_ENABLE         = "touchpad-enabled"
	TPAD_KEY_LEFT_HAND      = "left-handed"
	TPAD_KEY_W_TYPING       = "disable-while-typing"
	TPAD_KEY_NATURAL_SCROLL = "natural-scroll"
	TPAD_KEY_EDGE_SCROLL    = "edge-scroll-enabled"
	TPAD_KEY_HORIZ_SCROLL   = "horiz-scroll-enabled"
	TPAD_KEY_VERT_SCROLL    = "vert-scroll-enabled"
	TPAD_KEY_ACCEL          = "motion-acceleration"
	TPAD_KEY_THRES          = "motion-threshold"
	TPAD_KEY_TAP_CLICK      = "tap-to-click"
	TPAD_KEY_DELTA          = "delta-scroll"
)

var tpadSettings = gio.NewSettings("com.deepin.dde.touchpad")

type TouchpadManager struct {
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
	DeviceList []PointerDeviceInfo

	typingExitChan chan bool
	typingState    bool

	listenFlag bool
}

var _tpadManager *TouchpadManager

func GetTouchpadManager() *TouchpadManager {
	if _tpadManager == nil {
		_tpadManager = newTouchpadManager()
	}

	return _tpadManager
}

func newTouchpadManager() *TouchpadManager {
	tManager := &TouchpadManager{}

	tManager.TPadEnable = property.NewGSettingsBoolProperty(
		tManager, "TPadEnable",
		tpadSettings, TPAD_KEY_ENABLE)
	tManager.LeftHanded = property.NewGSettingsBoolProperty(
		tManager, "LeftHanded",
		tpadSettings, TPAD_KEY_LEFT_HAND)
	tManager.DisableIfTyping = property.NewGSettingsBoolProperty(
		tManager, "DisableIfTyping",
		tpadSettings, TPAD_KEY_W_TYPING)
	tManager.NaturalScroll = property.NewGSettingsBoolProperty(
		tManager, "NaturalScroll",
		tpadSettings, TPAD_KEY_NATURAL_SCROLL)
	tManager.EdgeScroll = property.NewGSettingsBoolProperty(
		tManager, "EdgeScroll",
		tpadSettings, TPAD_KEY_EDGE_SCROLL)
	tManager.VertScroll = property.NewGSettingsBoolProperty(
		tManager, "VertScroll",
		tpadSettings, TPAD_KEY_VERT_SCROLL)
	tManager.HorizScroll = property.NewGSettingsBoolProperty(
		tManager, "HorizScroll",
		tpadSettings, TPAD_KEY_HORIZ_SCROLL)
	tManager.TapClick = property.NewGSettingsBoolProperty(
		tManager, "TapClick",
		tpadSettings, TPAD_KEY_TAP_CLICK)

	tManager.MotionAcceleration = property.NewGSettingsFloatProperty(
		tManager, "MotionAcceleration",
		tpadSettings, TPAD_KEY_ACCEL)
	tManager.MotionThreshold = property.NewGSettingsFloatProperty(
		tManager, "MotionThreshold",
		tpadSettings, TPAD_KEY_THRES)

	tManager.DeltaScroll = property.NewGSettingsIntProperty(
		tManager, "DeltaScroll",
		tpadSettings, TPAD_KEY_DELTA)
	tManager.DoubleClick = property.NewGSettingsIntProperty(
		tManager, "DoubleClick",
		mouseSettings, MOUSE_KEY_DOUBLE_CLICK)
	tManager.DragThreshold = property.NewGSettingsIntProperty(
		tManager, "DragThreshold",
		mouseSettings, MOUSE_KEY_DRAG_THRES)

	_, tpadList, _ := getPointerDeviceList()
	tManager.setPropDeviceList(tpadList)

	if len(tManager.DeviceList) > 0 {
		tManager.setPropExist(true)
	} else {
		tManager.setPropExist(false)
	}

	tManager.typingExitChan = make(chan bool, 1)

	tManager.listenFlag = false

	tManager.init()

	return tManager
}

func (tManager *TouchpadManager) enableTPadWhileTyping() {
	if tManager.typingState {
		logger.Debug("syndaemon has running...")
		return
	}

	tManager.typingState = true
	go exec.Command("/usr/bin/syndaemon",
		"-i", "1",
		"-K", "-R").Run()
	select {
	case <-tManager.typingExitChan:
		go exec.Command("/usr/bin/killall",
			"/usr/bin/syndaemon").Run()
		tManager.typingState = false
		return
	}
}

func (tManager *TouchpadManager) disableTPadWhileTyping(enable bool) {
	if tpadEnable := tpadSettings.GetBoolean(TPAD_KEY_ENABLE); !tpadEnable {
		if tManager.typingState {
			tManager.typingExitChan <- true
		}
		return
	}

	if !enable {
		if tManager.typingState {
			tManager.typingExitChan <- true
		}
	} else {
		go tManager.enableTPadWhileTyping()
	}
}

func (tManager *TouchpadManager) enable(enabled bool) {
	for _, info := range tManager.DeviceList {
		setTouchpadEnabled(info.Deviceid, enabled)
	}
}

func (tManager *TouchpadManager) leftHanded(enabled bool) {
	for _, info := range tManager.DeviceList {
		setLeftHanded(info.Deviceid, info.Name, enabled)
	}
}

func (tManager *TouchpadManager) motionAcceleration(accel float64) {
	for _, info := range tManager.DeviceList {
		setMotionAcceleration(info.Deviceid, accel)
	}
}

func (tManager *TouchpadManager) motionThreshold(thres float64) {
	for _, info := range tManager.DeviceList {
		setMotionThreshold(info.Deviceid, thres)
	}
}

func (tManager *TouchpadManager) naturalScroll(enabled bool, delta int32) {
	for _, info := range tManager.DeviceList {
		setTouchpadNaturalScroll(info.Deviceid, enabled, delta)
	}
}

func (tManager *TouchpadManager) edgeScroll(enabled bool) {
	for _, info := range tManager.DeviceList {
		setTouchpadEdgeScroll(info.Deviceid, enabled)
	}
}

func (tManager *TouchpadManager) twoFingerScroll(vertEnabled,
	horizEnabled bool) {
	for _, info := range tManager.DeviceList {
		setTouchpadTwoFingerScroll(info.Deviceid,
			vertEnabled, horizEnabled)
	}
}

func (tManager *TouchpadManager) tapToClick(enabled, leftHanded bool) {
	for _, info := range tManager.DeviceList {
		setTouchpadTapToClick(info.Deviceid,
			enabled, leftHanded)
	}
}

func (tManager *TouchpadManager) init() {
	if !tManager.Exist {
		return
	}

	enabled := tManager.TPadEnable.Get()
	tManager.enable(enabled)
	if !enabled {
		return
	}

	logger.Debug("Set leftHanded")
	tManager.leftHanded(tManager.LeftHanded.Get())
	logger.Debug("Set tap click")
	tManager.tapToClick(tManager.TapClick.Get(),
		tManager.LeftHanded.Get())
	logger.Debug("Set while typing")
	tManager.disableTPadWhileTyping(tManager.DisableIfTyping.Get())
	logger.Debug("Set natural scroll")
	tManager.naturalScroll(tManager.NaturalScroll.Get(),
		tManager.DeltaScroll.Get())
	logger.Debug("Set edge scroll")
	tManager.edgeScroll(tManager.EdgeScroll.Get())
	logger.Debug("Set two finger scroll")
	tManager.twoFingerScroll(tManager.VertScroll.Get(),
		tManager.HorizScroll.Get())

	logger.Debug("Set acceleration")
	tManager.motionAcceleration(tManager.MotionAcceleration.Get())
	logger.Debug("Set threshold")
	tManager.motionThreshold(tManager.MotionThreshold.Get())

	if !tManager.listenFlag {
		tManager.listenGSettings()
	}
}
