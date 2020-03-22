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
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"sync"

	"pkg.deepin.io/gir/gio-2.0"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/gsprop"
	"pkg.deepin.io/lib/strv"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	tpadSchema = "com.deepin.dde.touchpad"

	tpadKeyEnabled            = "touchpad-enabled"
	tpadKeyLeftHanded         = "left-handed"
	tpadKeyDisableWhileTyping = "disable-while-typing"
	tpadKeyNaturalScroll      = "natural-scroll"
	tpadKeyEdgeScroll         = "edge-scroll-enabled"
	tpadKeyHorizScroll        = "horiz-scroll-enabled"
	tpadKeyVertScroll         = "vert-scroll-enabled"
	tpadKeyAcceleration       = "motion-acceleration"
	tpadKeyThreshold          = "motion-threshold"
	tpadKeyScaling            = "motion-scaling"
	tpadKeyTapClick           = "tap-to-click"
	tpadKeyScrollDelta        = "delta-scroll"
	tpadKeyWhileTypingCmd     = "disable-while-typing-cmd"
	tpadKeyPalmDetect         = "palm-detect"
	tpadKeyPalmMinWidth       = "palm-min-width"
	tpadKeyPalmMinZ           = "palm-min-pressure"
)

const (
	syndaemonPidFile = "/tmp/syndaemon.pid"
)

type Touchpad struct {
	service    *dbusutil.Service
	PropsMu    sync.RWMutex
	Exist      bool
	DeviceList string

	// dbusutil-gen: ignore-below
	TPadEnable      gsprop.Bool `prop:"access:rw"`
	LeftHanded      gsprop.Bool `prop:"access:rw"`
	DisableIfTyping gsprop.Bool `prop:"access:rw"`
	NaturalScroll   gsprop.Bool `prop:"access:rw"`
	EdgeScroll      gsprop.Bool `prop:"access:rw"`
	HorizScroll     gsprop.Bool `prop:"access:rw"`
	VertScroll      gsprop.Bool `prop:"access:rw"`
	TapClick        gsprop.Bool `prop:"access:rw"`
	PalmDetect      gsprop.Bool `prop:"access:rw"`

	MotionAcceleration gsprop.Double `prop:"access:rw"`
	MotionThreshold    gsprop.Double `prop:"access:rw"`
	MotionScaling      gsprop.Double `prop:"access:rw"`

	DoubleClick   gsprop.Int `prop:"access:rw"`
	DragThreshold gsprop.Int `prop:"access:rw"`
	DeltaScroll   gsprop.Int `prop:"access:rw"`
	PalmMinWidth  gsprop.Int `prop:"access:rw"`
	PalmMinZ      gsprop.Int `prop:"access:rw"`

	devInfos     dxTouchpads
	setting      *gio.Settings
	mouseSetting *gio.Settings
}

func newTouchpad(service *dbusutil.Service) *Touchpad {
	var tpad = new(Touchpad)

	tpad.service = service
	tpad.setting = gio.NewSettings(tpadSchema)
	tpad.TPadEnable.Bind(tpad.setting, tpadKeyEnabled)
	tpad.LeftHanded.Bind(tpad.setting, tpadKeyLeftHanded)
	tpad.DisableIfTyping.Bind(tpad.setting, tpadKeyDisableWhileTyping)
	tpad.NaturalScroll.Bind(tpad.setting, tpadKeyNaturalScroll)
	tpad.EdgeScroll.Bind(tpad.setting, tpadKeyEdgeScroll)
	tpad.VertScroll.Bind(tpad.setting, tpadKeyVertScroll)
	tpad.HorizScroll.Bind(tpad.setting, tpadKeyHorizScroll)
	tpad.TapClick.Bind(tpad.setting, tpadKeyTapClick)
	tpad.PalmDetect.Bind(tpad.setting, tpadKeyPalmDetect)
	tpad.MotionAcceleration.Bind(tpad.setting, tpadKeyAcceleration)
	tpad.MotionThreshold.Bind(tpad.setting, tpadKeyThreshold)
	tpad.MotionScaling.Bind(tpad.setting, tpadKeyScaling)
	tpad.DeltaScroll.Bind(tpad.setting, tpadKeyScrollDelta)
	tpad.PalmMinWidth.Bind(tpad.setting, tpadKeyPalmMinWidth)
	tpad.PalmMinZ.Bind(tpad.setting, tpadKeyPalmMinZ)

	tpad.mouseSetting = gio.NewSettings(mouseSchema)
	tpad.DoubleClick.Bind(tpad.mouseSetting, mouseKeyDoubleClick)
	tpad.DragThreshold.Bind(tpad.mouseSetting, mouseKeyDragThreshold)

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
	tpad.enablePalmDetect()
	tpad.setPalmDimensions()
}

func (tpad *Touchpad) handleDeviceChanged() {
	tpad.updateDXTpads()
	tpad.init()
}

func (tpad *Touchpad) updateDXTpads() {
	tpad.devInfos = dxTouchpads{}
	for _, info := range getTPadInfos(false) {
		if !globalWayland {
			tmp := tpad.devInfos.get(info.Id)
			if tmp != nil {
				continue
			}
		}
		tpad.devInfos = append(tpad.devInfos, info)
	}

	tpad.PropsMu.Lock()
	var v string
	if len(tpad.devInfos) == 0 {
		tpad.setPropExist(false)
	} else {
		tpad.setPropExist(true)
		v = tpad.devInfos.string()
	}
	tpad.setPropDeviceList(v)
	tpad.PropsMu.Unlock()
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
	setWMTPadBoolKey(wmTPadKeyLeftHanded, enabled)
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
	setWMTPadBoolKey(wmTPadKeyNaturalScroll, enabled)
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
	setWMTPadBoolKey(wmTPadKeyEdgeScroll, enabled)
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
	setWMTPadBoolKey(wmTPadKeyTapClick, enabled)
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

	syncmd := tpad.setting.GetString(tpadKeyWhileTypingCmd)
	if syncmd == "" {
		logger.Warning("Failed to start syndaemon, because no cmd is specified")
		return
	}
	logger.Debug("[startSyndaemon] will exec:", syncmd)
	args := strings.Split(syncmd, " ")
	argsLen := len(args)
	var cmd *exec.Cmd
	if argsLen == 1 {
		// pidfile will be created only in daemon mode
		cmd = exec.Command(args[0], "-d", "-p", syndaemonPidFile)
	} else {
		list := strv.Strv(args)
		if !list.Contains("-p") {
			if !list.Contains("-d") {
				args = append(args, "-d")
			}
			args = append(args, "-p", syndaemonPidFile)
		}
		argsLen = len(args)
		cmd = exec.Command(args[0], args[1:argsLen]...)
	}
	err := cmd.Start()
	if err != nil {
		os.Remove(syndaemonPidFile)
		logger.Debug("[disableWhileTyping] start syndaemon failed:", err)
		return
	}

	go cmd.Wait()
}

func (tpad *Touchpad) stopSyndaemon() {
	out, err := exec.Command("killall", "syndaemon").CombinedOutput()
	if err != nil {
		logger.Warning("[stopSyndaemon] failed:", string(out), err)
	}
	os.Remove(syndaemonPidFile)
}

func (tpad *Touchpad) enablePalmDetect() {
	enabled := tpad.PalmDetect.Get()
	for _, dev := range tpad.devInfos {
		err := dev.EnablePalmDetect(enabled)
		if err != nil {
			logger.Warning("[enablePalmDetect] failed to enable:", dev.Id, enabled, err)
		}
	}
}

func (tpad *Touchpad) setPalmDimensions() {
	width := tpad.PalmMinWidth.Get()
	z := tpad.PalmMinZ.Get()
	for _, dev := range tpad.devInfos {
		err := dev.SetPalmDimensions(width, z)
		if err != nil {
			logger.Warning("[setPalmDimensions] failed to set:", dev.Id, width, z, err)
		}
	}
}

func isSyndaemonExist(pidFile string) bool {
	if !dutils.IsFileExist(pidFile) {
		out, err := exec.Command("pgrep", "syndaemon").CombinedOutput()
		if err != nil || len(out) < 2 {
			return false
		}
		return true
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
