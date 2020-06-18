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

package power

import (
	"fmt"
	x "github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/ext/dpms"
	"github.com/linuxdeepin/go-x11-client/util/wm/icccm"
	"os"
	"os/exec"
	"pkg.deepin.io/dde/api/soundutils"
	dbus "pkg.deepin.io/lib/dbus1"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/gsettings"
	"pkg.deepin.io/lib/pulse"
	"time"
)

func (m *Manager) findWindow(wmClassInstance, wmClassClass string) x.Window {
	c := m.helper.xConn
	rootWin := c.GetDefaultScreen().Root
	tree, err := x.QueryTree(c, rootWin).Reply(c)
	if err != nil {
		logger.Warning("QueryTree error:", err)
		return 0
	}

	for _, win := range tree.Children {
		wmClass, err := icccm.GetWMClass(c, win).Reply(c)
		if err == nil &&
			wmClass.Instance == wmClassInstance &&
			wmClass.Class == wmClassClass {
			return win
		}
	}
	return 0
}

func (m *Manager) waitLockShowing(timeout time.Duration) {
	ticker := time.NewTicker(time.Millisecond * 300)
	timer := time.NewTimer(timeout)
	for {
		select {
		case _, ok := <-ticker.C:
			if !ok {
				logger.Error("Invalid ticker event")
				return
			}

			logger.Debug("waitLockShowing tick")
			locked, err := m.helper.SessionManager.Locked().Get(0)
			if err != nil {
				logger.Warning(err)
				continue
			}
			if locked {
				logger.Debug("waitLockShowing found")
				ticker.Stop()
				return
			}

		case <-timer.C:
			logger.Debug("waitLockShowing failed timeout!")
			ticker.Stop()
			return
		}
	}
}

func (m *Manager) lockWaitShow(timeout time.Duration, autoStartAuth bool) {
	m.doLock(autoStartAuth)
	if m.UseWayland {
		logger.Debug("In wayland environment, unsupported check lock whether showin")
		return
	}
	m.waitLockShowing(timeout)
}

func (m *Manager) setDPMSModeOn() {
	logger.Info("DPMS On")
	var err error

	if m.UseWayland {
		_, err = exec.Command("dde_wldpms", "-s", "On").Output()
	} else {
		c := m.helper.xConn
		err = dpms.ForceLevelChecked(c, dpms.DPMSModeOn).Check(c)
	}
	if err != nil {
		logger.Warning("Set DPMS on error:", err)
	}
}

func (m *Manager) setDPMSModeOff() {
	logger.Info("DPMS Off")
	var err error
	if m.UseWayland {
		_, err = exec.Command("dde_wldpms", "-s", "Off").Output()
	} else {
		c := m.helper.xConn
		err = dpms.ForceLevelChecked(c, dpms.DPMSModeOff).Check(c)
	}
	if err != nil {
		logger.Warning("Set DPMS off error:", err)
	}
}

const (
	lockFrontServiceName = "com.deepin.dde.lockFront"
	lockFrontIfc         = lockFrontServiceName
	lockFrontObjPath     = "/com/deepin/dde/lockFront"
)

func (m *Manager) doLock(autoStartAuth bool) {
	logger.Info("Lock Screen")
	bus, err := dbus.SessionBus()
	if err != nil {
		logger.Warning(err)
		return
	}
	lockFrontObj := bus.Object(lockFrontServiceName, lockFrontObjPath)
	err = lockFrontObj.Call(lockFrontIfc+".ShowAuth", 0, autoStartAuth).Err
	if err != nil {
		logger.Warning("failed to call lockFront ShowAuth:", err)
	}
}

func (m *Manager) doSuspend() {
	if os.Getenv("POWER_CAN_SLEEP") == "0" {
		logger.Info("can not suspend, env POWER_CAN_SLEEP == 0")
		return
	}

	sessionManager := m.helper.SessionManager
	can, err := sessionManager.CanSuspend(0)
	if err != nil {
		logger.Warning(err)
		return
	}

	if !can {
		logger.Info("can not suspend")
		return
	}

	logger.Debug("suspend")
	err = sessionManager.RequestSuspend(0)
	if err != nil {
		logger.Warning("failed to suspend:", err)
	}
}

func (m *Manager) doShutdown() {
	sessionManager := m.helper.SessionManager
	can, err := sessionManager.CanShutdown(0)
	if err != nil {
		logger.Warning(err)
		return
	}

	if !can {
		logger.Info("can not Shutdown")
		return
	}

	logger.Debug("Shutdown")
	err = sessionManager.RequestShutdown(0)
	if err != nil {
		logger.Warning("failed to Shutdown:", err)
	}
}

func (m *Manager) doHibernate() {
	if os.Getenv("POWER_CAN_SLEEP") == "0" {
		logger.Info("can not hibernate, env POWER_CAN_SLEEP == 0")
		return
	}

	sessionManager := m.helper.SessionManager
	can, err := sessionManager.CanHibernate(0)
	if err != nil {
		logger.Warning(err)
		return
	}

	if !can {
		logger.Info("can not hibernate")
		return
	}

	logger.Debug("hibernate")
	err = sessionManager.RequestHibernate(0)
	if err != nil {
		logger.Warning("failed to hibernate:", err)
	}
}

func (m *Manager) doTurnOffScreen() {
	logger.Info("Turn off screen")
	m.setDPMSModeOff()
}

func (m *Manager) setDisplayBrightness(brightnessTable map[string]float64) {
	display := m.helper.Display
	for output, brightness := range brightnessTable {
		logger.Infof("Change output %q brightness to %.2f", output, brightness)
		err := display.SetBrightness(0, output, brightness)
		if err != nil {
			logger.Warningf("Change output %q brightness to %.2f failed: %v", output, brightness, err)
		} else {
			logger.Infof("Change output %q brightness to %.2f done!", output, brightness)
		}
	}
}

func (m *Manager) setAndSaveDisplayBrightness(brightnessTable map[string]float64) {
	display := m.helper.Display
	for output, brightness := range brightnessTable {
		logger.Infof("Change output %q brightness to %.2f", output, brightness)
		err := display.SetAndSaveBrightness(0, output, brightness)
		if err != nil {
			logger.Warningf("Change output %q brightness to %.2f failed: %v", output, brightness, err)
		} else {
			logger.Infof("Change output %q brightness to %.2f done!", output, brightness)
		}
	}
}

func doShowDDELowPower() {
	logger.Info("Show dde low power")
	go exec.Command(cmdDDELowPower, "--raise").Run()
}

func doCloseDDELowPower() {
	logger.Info("Close low power")
	go exec.Command(cmdDDELowPower, "--quit").Run()
}

func (m *Manager) sendNotify(icon, summary, body string) {
	if !m.LowPowerNotifyEnable.Get() {
		logger.Info("notify switch is off ")
		return
	}
	n := m.helper.Notifications
	_, err := n.Notify(0, "dde-control-center", 0, icon, summary, body, nil, nil, -1)
	if err != nil {
		logger.Warning(err)
	}
}

func (m *Manager) sendChangeNotify(icon, summary, body string) {
	n := m.helper.Notifications
	_, err := n.Notify(0, "dde-control-center", 0, icon, summary, body, nil, nil, -1)
	if err != nil {
		logger.Warning(err)
	}
}

const iconBatteryLow = "notification-battery-low"

func playSound(name string) {
	logger.Debug("play system sound", name)
	go soundutils.PlaySystemSound(name, "")
}

const (
	deepinScreensaverDBusServiceName = "com.deepin.ScreenSaver"
	deepinScreensaverDBusPath        = "/com/deepin/ScreenSaver"
	deepinScreensaverDBusInterface   = deepinScreensaverDBusServiceName
)

func startScreensaver() {
	logger.Info("start screensaver")
	bus, err := dbus.SessionBus()
	if err != nil {
		logger.Warning(err)
		return
	}

	obj := bus.Object(deepinScreensaverDBusServiceName, deepinScreensaverDBusPath)
	err = obj.Call(deepinScreensaverDBusInterface+".Start", 0).Err
	if err != nil {
		logger.Warning(err)
	}
}

func stopScreensaver() {
	logger.Info("stop screensaver")
	bus, err := dbus.SessionBus()
	if err != nil {
		logger.Warning(err)
		return
	}

	obj := bus.Object(deepinScreensaverDBusServiceName, deepinScreensaverDBusPath)
	err = obj.Call(deepinScreensaverDBusInterface+".Stop", 0).Err
	if err != nil {
		logger.Warning(err)
	}
}

// TODO(jouyouyun): move to common library
func suspendPulseSinks(suspend int) {
	var ctx = pulse.GetContext()
	if ctx == nil {
		logger.Error("Failed to connect pulseaudio server")
		return
	}
	for _, sink := range ctx.GetSinkList() {
		ctx.SuspendSinkById(sink.Index, suspend)
	}
}

// TODO(jouyouyun): move to common library
func suspendPulseSources(suspend int) {
	var ctx = pulse.GetContext()
	if ctx == nil {
		logger.Error("Failed to connect pulseaudio server")
		return
	}
	for _, source := range ctx.GetSourceList() {
		ctx.SuspendSourceById(source.Index, suspend)
	}
}

func brightnessRound(x float64) string {
	return fmt.Sprintf("%.2f", x)
}

func (m *Manager) initGSettingsConnectChanged() {
	const powerSettingsIcon = "preferences-system"
	// 监听 session power 的属性的改变,并发送通知
	gsettings.ConnectChanged(gsSchemaPower, "*", func(key string) {
		logger.Debug("Power Settings Changed :", key)
		switch key {
		case settingKeyLinePowerLidClosedAction:
			value := m.LinePowerLidClosedAction.Get()
			notifyString := getNotifyString(settingKeyLinePowerLidClosedAction, value)
			m.sendChangeNotify(powerSettingsIcon, Tr("Power settings changed"), notifyString)
		case settingKeyLinePowerPressPowerBtnAction:
			value := m.LinePowerPressPowerBtnAction.Get()
			notifyString := getNotifyString(settingKeyLinePowerPressPowerBtnAction, value)
			m.sendChangeNotify(powerSettingsIcon, Tr("Power settings changed"), notifyString)
		case settingKeyBatteryLidClosedAction:
			value := m.BatteryLidClosedAction.Get()
			notifyString := getNotifyString(settingKeyBatteryLidClosedAction, value)
			m.sendChangeNotify(powerSettingsIcon, Tr("Power settings changed"), notifyString)
		case settingKeyBatteryPressPowerBtnAction:
			value := m.BatteryPressPowerBtnAction.Get()
			notifyString := getNotifyString(settingKeyBatteryPressPowerBtnAction, value)
			m.sendChangeNotify(powerSettingsIcon, Tr("Power settings changed"), notifyString)
		}
	})
}

func getNotifyString(option string, action int32) string {
	var (
		notifyString string
		firstPart    string
		secondPart   string
	)
	switch option {
	case
		settingKeyLinePowerLidClosedAction,
		settingKeyBatteryLidClosedAction:
		firstPart = Tr("When the lid is closed, ")
	case
		settingKeyLinePowerPressPowerBtnAction,
		settingKeyBatteryPressPowerBtnAction:
		firstPart = Tr("When pressing the power button, ")
	}
	secondPart = getPowerActionString(action)
	notifyString = firstPart + secondPart
	return notifyString
}

func getPowerActionString(action int32) string {
	switch action {
	case powerActionShutdown:
		return Tr("your computer will shut down")
	case powerActionSuspend:
		return Tr("your computer will suspend")
	case powerActionHibernate:
		return Tr("your computer will sleep")
	case powerActionTurnOffScreen:
		return Tr("your monitor will turn off")
	case powerActionDoNothing:
		return Tr("it will do nothing to your computer")
	}
	return ""
}
