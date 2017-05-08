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
	"fmt"
	"gir/gio-2.0"
	"io/ioutil"
	"os"
	"os/exec"
	"pkg.deepin.io/lib/dbus/property"
	dutils "pkg.deepin.io/lib/utils"
	"strconv"
	"strings"
)

const (
	tpadSchema = "com.deepin.dde.touchpad"

	tpadKeyEnabled       = "touchpad-enabled"
	tpadKeyLeftHanded    = "left-handed"
	tpadKeyWhileTyping   = "disable-while-typing"
	tpadKeyNaturalScroll = "natural-scroll"
	tpadKeyEdgeScroll    = "edge-scroll-enabled"
	tpadKeyHorizScroll   = "horiz-scroll-enabled"
	tpadKeyVertScroll    = "vert-scroll-enabled"
	tpadKeyAcceleration  = "motion-acceleration"
	tpadKeyThreshold     = "motion-threshold"
	tpadKeyScaling       = "motion-scaling"
	tpadKeyTapClick      = "tap-to-click"
	tpadKeyScrollDelta   = "delta-scroll"
)

const (
	syndaemonPidFile = "/tmp/syndaemon.pid"
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
	MotionScaling      *property.GSettingsFloatProperty `access:"readwrite"`

	DoubleClick   *property.GSettingsIntProperty `access:"readwrite"`
	DragThreshold *property.GSettingsIntProperty `access:"readwrite"`
	DeltaScroll   *property.GSettingsIntProperty `access:"readwrite"`

	Exist      bool
	DeviceList string

	devInfos     dxTouchpads
	setting      *gio.Settings
	mouseSetting *gio.Settings
	synProcess   *os.Process
}

var _tpad *Touchpad

func getTouchpad() *Touchpad {
	if _tpad == nil {
		_tpad = NewTouchpad()
	}

	return _tpad
}

func NewTouchpad() *Touchpad {
	var tpad = new(Touchpad)

	tpad.setting = gio.NewSettings(tpadSchema)

	tpad.TPadEnable = property.NewGSettingsBoolProperty(
		tpad, "TPadEnable",
		tpad.setting, tpadKeyEnabled)
	tpad.LeftHanded = property.NewGSettingsBoolProperty(
		tpad, "LeftHanded",
		tpad.setting, tpadKeyLeftHanded)
	tpad.DisableIfTyping = property.NewGSettingsBoolProperty(
		tpad, "DisableIfTyping",
		tpad.setting, tpadKeyWhileTyping)
	tpad.NaturalScroll = property.NewGSettingsBoolProperty(
		tpad, "NaturalScroll",
		tpad.setting, tpadKeyNaturalScroll)
	tpad.EdgeScroll = property.NewGSettingsBoolProperty(
		tpad, "EdgeScroll",
		tpad.setting, tpadKeyEdgeScroll)
	tpad.VertScroll = property.NewGSettingsBoolProperty(
		tpad, "VertScroll",
		tpad.setting, tpadKeyVertScroll)
	tpad.HorizScroll = property.NewGSettingsBoolProperty(
		tpad, "HorizScroll",
		tpad.setting, tpadKeyHorizScroll)
	tpad.TapClick = property.NewGSettingsBoolProperty(
		tpad, "TapClick",
		tpad.setting, tpadKeyTapClick)

	tpad.MotionAcceleration = property.NewGSettingsFloatProperty(
		tpad, "MotionAcceleration",
		tpad.setting, tpadKeyAcceleration)
	tpad.MotionThreshold = property.NewGSettingsFloatProperty(
		tpad, "MotionThreshold",
		tpad.setting, tpadKeyThreshold)
	tpad.MotionScaling = property.NewGSettingsFloatProperty(
		tpad, "MotionScaling",
		tpad.setting, tpadKeyScaling)

	tpad.DeltaScroll = property.NewGSettingsIntProperty(
		tpad, "DeltaScroll",
		tpad.setting, tpadKeyScrollDelta)

	tpad.mouseSetting = gio.NewSettings(mouseSchema)
	tpad.DoubleClick = property.NewGSettingsIntProperty(
		tpad, "DoubleClick",
		tpad.mouseSetting, mouseKeyDoubleClick)
	tpad.DragThreshold = property.NewGSettingsIntProperty(
		tpad, "DragThreshold",
		tpad.mouseSetting, mouseKeyDragThreshold)

	tpad.updateDXTpads()

	return tpad
}

func (tpad *Touchpad) init() {
	if !tpad.Exist {
		return
	}

	tpad.enable(tpad.TPadEnable.Get())
	tpad.enableLeftHanded()
	tpad.enableNaturalScroll()
	tpad.enableEdgeScroll()
	tpad.enableTapToClick()
	tpad.enableTwoFingerScroll()
	tpad.motionAcceleration()
	tpad.motionThreshold()
	tpad.motionScaling()
	tpad.disableWhileTyping()
}

func (tpad *Touchpad) handleDeviceChanged() {
	tpad.updateDXTpads()
	tpad.init()
}

func (tpad *Touchpad) updateDXTpads() {
	tpad.devInfos = dxTouchpads{}
	for _, info := range getTPadInfos(false) {
		tmp := tpad.devInfos.get(info.Id)
		if tmp != nil {
			continue
		}
		tpad.devInfos = append(tpad.devInfos, info)
	}

	var v string
	if len(tpad.devInfos) == 0 {
		tpad.setPropExist(false)
	} else {
		tpad.setPropExist(true)
		v = tpad.devInfos.string()
	}
	setPropString(tpad, &tpad.DeviceList, "DeviceList", v)
}

func (tpad *Touchpad) enable(enabled bool) {
	for _, v := range tpad.devInfos {
		err := v.Enable(enabled)
		if err != nil {
			logger.Warningf("Enable '%v - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
	enableGesture(enabled)
}

func (tpad *Touchpad) enableLeftHanded() {
	enabled := tpad.LeftHanded.Get()
	for _, v := range tpad.devInfos {
		err := v.EnableLeftHanded(enabled)
		if err != nil {
			logger.Debugf("Enable left handed '%v - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (tpad *Touchpad) enableNaturalScroll() {
	enabled := tpad.NaturalScroll.Get()
	for _, v := range tpad.devInfos {
		err := v.EnableNaturalScroll(enabled)
		if err != nil {
			logger.Debugf("Enable natural scroll '%v - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (tpad *Touchpad) setScrollDistance() {
	delta := tpad.DeltaScroll.Get()
	for _, v := range tpad.devInfos {
		err := v.SetScrollDistance(delta, delta)
		if err != nil {
			logger.Debugf("Set natural scroll distance '%v - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (tpad *Touchpad) enableEdgeScroll() {
	enabled := tpad.EdgeScroll.Get()
	for _, v := range tpad.devInfos {
		err := v.EnableEdgeScroll(enabled)
		if err != nil {
			logger.Debugf("Enable edge scroll '%v - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (tpad *Touchpad) enableTwoFingerScroll() {
	vert := tpad.VertScroll.Get()
	horiz := tpad.HorizScroll.Get()
	for _, v := range tpad.devInfos {
		err := v.EnableTwoFingerScroll(vert, horiz)
		if err != nil {
			logger.Debugf("Enable two-finger scroll '%v - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (tpad *Touchpad) enableTapToClick() {
	enabled := tpad.TapClick.Get()
	for _, v := range tpad.devInfos {
		err := v.EnableTapToClick(enabled)
		if err != nil {
			logger.Debugf("Enable tap to click '%v - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (tpad *Touchpad) motionAcceleration() {
	accel := float32(tpad.MotionAcceleration.Get())
	for _, v := range tpad.devInfos {
		err := v.SetMotionAcceleration(accel)
		if err != nil {
			logger.Debugf("Set acceleration for '%d - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (tpad *Touchpad) motionThreshold() {
	thres := float32(tpad.MotionThreshold.Get())
	for _, v := range tpad.devInfos {
		err := v.SetMotionThreshold(thres)
		if err != nil {
			logger.Debugf("Set threshold for '%d - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (tpad *Touchpad) motionScaling() {
	scaling := float32(tpad.MotionScaling.Get())
	for _, v := range tpad.devInfos {
		err := v.SetMotionScaling(scaling)
		if err != nil {
			logger.Debugf("Set scaling for '%d - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (tpad *Touchpad) doubleClick() {
	xsSetInt32(xsPropDoubleClick, tpad.DoubleClick.Get())
}

func (tpad *Touchpad) dragThreshold() {
	xsSetInt32(xsPropDragThres, tpad.DragThreshold.Get())
}

func (tpad *Touchpad) disableWhileTyping() {
	if !tpad.Exist {
		return
	}

	var usedLibinput bool = false
	enabled := tpad.DisableIfTyping.Get()
	for _, v := range tpad.devInfos {
		err := v.EnableDisableWhileTyping(enabled)
		if err != nil {
			continue
		}
		usedLibinput = true
	}
	if usedLibinput {
		return
	}

	if enabled {
		tpad.startSyndaemon()
	} else {
		tpad.stopSyndaemon()
	}
}

func (tpad *Touchpad) startSyndaemon() {
	if isSyndaemonExist(syndaemonPidFile) {
		logger.Debug("Syndaemon has running")
		return
	}

	var cmd = exec.Command("/bin/sh", "-c", "syndaemon -i 1 -K -t")
	err := cmd.Start()
	if err != nil {
		os.Remove(syndaemonPidFile)
		logger.Debug("[disableWhileTyping] start syndaemon failed:", err)
		return
	}
	tpad.synProcess = cmd.Process
	content := fmt.Sprintf("%v", tpad.synProcess.Pid)
	ioutil.WriteFile(syndaemonPidFile, []byte(content), 0777)
}

func (tpad *Touchpad) stopSyndaemon() {
	if tpad.synProcess == nil {
		return
	}

	tpad.synProcess.Kill()
	tpad.synProcess = nil
	os.Remove(syndaemonPidFile)
}

func isSyndaemonExist(pidFile string) bool {
	if !dutils.IsFileExist(pidFile) {
		return false
	}

	context, err := ioutil.ReadFile(pidFile)
	if err != nil {
		return false
	}

	pid, err := strconv.ParseInt(strings.TrimSpace(string(context)), 10, 64)
	if err != nil {
		return false
	}
	var file = fmt.Sprintf("/proc/%v/cmdline", pid)
	return isProcessExist(file, "syndaemon")
}

func isProcessExist(file, name string) bool {
	context, err := ioutil.ReadFile(file)
	if err != nil {
		return false
	}

	return strings.Contains(string(context), name)
}

func enableGesture(enabled bool) {
	s, err := dutils.CheckAndNewGSettings("com.deepin.dde.gesture")
	if err != nil {
		return
	}
	if s.GetBoolean("enabled") == enabled {
		return
	}

	s.SetBoolean("enabled", enabled)
	s.Unref()
}
