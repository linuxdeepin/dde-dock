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
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	x "github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/util/wm/ewmh"
	"pkg.deepin.io/lib/gsettings"
	"pkg.deepin.io/lib/procfs"
)

const submodulePSP = "PowerSavePlan"

func init() {
	submoduleList = append(submoduleList, newPowerSavePlan)
}

type powerSavePlan struct {
	manager            *Manager
	screenSaverTimeout int32
	metaTasks          metaTasks
	tasks              delayedTasks
	// key output name, value old brightness
	oldBrightnessTable map[string]float64
	mu                 sync.Mutex
	idleOn             bool
	screensaverRunning bool

	atomNetWMStateFullscreen    x.Atom
	atomNetWMStateFocused       x.Atom
	fullscreenWorkaroundAppList []string
}

func newPowerSavePlan(manager *Manager) (string, submodule, error) {
	p := new(powerSavePlan)
	p.manager = manager

	conn := manager.helper.xConn
	var err error
	p.atomNetWMStateFullscreen, err = conn.GetAtom("_NET_WM_STATE_FULLSCREEN")
	if err != nil {
		return submodulePSP, nil, err
	}
	p.atomNetWMStateFocused, err = conn.GetAtom("_NET_WM_STATE_FOCUSED")
	if err != nil {
		return submodulePSP, nil, err
	}

	p.fullscreenWorkaroundAppList = manager.settings.GetStrv(
		"fullscreen-workaround-app-list")
	return submodulePSP, p, nil
}

// 监听 GSettings 值改变, 更新节电计划
func (psp *powerSavePlan) initSettingsChangedHandler() {
	m := psp.manager
	gsettings.ConnectChanged(gsSchemaPower, "*", func(key string) {
		logger.Debug("setting changed", key)
		switch key {
		case settingKeyLinePowerScreensaverDelay,
			settingKeyLinePowerScreenBlackDelay,
			settingKeyLinePowerSleepDelay:
			if !m.OnBattery {
				logger.Debug("Change OnLinePower plan")
				psp.OnLinePower()
			}

		case settingKeyBatteryScreensaverDelay,
			settingKeyBatteryScreenBlackDelay,
			settingKeyBatterySleepDelay:
			if m.OnBattery {
				logger.Debug("Change OnBattery plan")
				psp.OnBattery()
			}

		case settingKeyAmbientLightAdjuestBrightness:
			psp.manager.claimOrReleaseAmbientLight()
		}
	})
}

func (psp *powerSavePlan) OnBattery() {
	logger.Debug("Use OnBattery plan")
	m := psp.manager
	psp.Update(m.BatteryScreensaverDelay.Get(), m.BatteryScreenBlackDelay.Get(),
		m.BatterySleepDelay.Get())
}

func (psp *powerSavePlan) OnLinePower() {
	logger.Debug("Use OnLinePower plan")
	m := psp.manager
	psp.Update(m.LinePowerScreensaverDelay.Get(), m.LinePowerScreenBlackDelay.Get(),
		m.LinePowerSleepDelay.Get())
}

func (psp *powerSavePlan) Reset() {
	m := psp.manager
	logger.Debug("OnBattery:", m.OnBattery)
	if m.OnBattery {
		psp.OnBattery()
	} else {
		psp.OnLinePower()
	}
}

func (psp *powerSavePlan) Start() error {
	psp.Reset()
	psp.initSettingsChangedHandler()

	helper := psp.manager.helper
	power := helper.Power
	screenSaver := helper.ScreenSaver

	//OnBattery changed will effect current PowerSavePlan
	power.OnBattery().ConnectChanged(func(hasValue bool, value bool) {
		psp.Reset()
	})

	power.PowerSavingModeEnabled().ConnectChanged(psp.handlePowerSavingModeChanged)

	screenSaver.ConnectIdleOn(psp.HandleIdleOn)
	screenSaver.ConnectIdleOff(psp.HandleIdleOff)
	return nil
}

func (psp *powerSavePlan) handlePowerSavingModeChanged(hasValue bool, enabled bool) {
	if !hasValue {
		return
	}
	logger.Debug("power saving mode enabled changed to", enabled)

	psp.manager.PropsMu.RLock()
	hasLightSensor := psp.manager.HasAmbientLightSensor
	psp.manager.PropsMu.RUnlock()

	if hasLightSensor && psp.manager.AmbientLightAdjustBrightness.Get() {
		return
	}

	brightnessTable, err := psp.manager.helper.Display.GetBrightness(0)
	if err != nil {
		logger.Warning(err)
		return
	}
	if enabled {
		// reduce brightness by 20%
		for key, value := range brightnessTable {
			value = value * 0.8
			if value < 0.1 {
				value = 0.1
			}
			brightnessTable[key] = value
		}

	} else {
		// increase brightness by 25%
		for key, value := range brightnessTable {
			value = value * 1.25
			if value > 1 {
				value = 1
			}
			brightnessTable[key] = value
		}
	}
	psp.manager.setAndSaveDisplayBrightness(brightnessTable)
}

// 取消之前的任务
func (psp *powerSavePlan) interruptTasks() {
	psp.tasks.CancelAll()
	psp.tasks.Wait(10*time.Millisecond, 200)
	psp.tasks = nil
}

func (psp *powerSavePlan) Destroy() {
	psp.interruptTasks()
}

func (psp *powerSavePlan) addTaskNoLock(t *delayedTask) {
	psp.tasks = append(psp.tasks, t)
}

func (psp *powerSavePlan) addTask(t *delayedTask) {
	psp.mu.Lock()
	psp.addTaskNoLock(t)
	psp.mu.Unlock()
}

type metaTask struct {
	step   byte
	delay  int32
	name   string
	ignore bool
	fn     func()
}

type metaTasks []metaTask

func (v metaTasks) Len() int {
	return len(v)
}

func (v metaTasks) Less(i, j int) bool {
	if v[i].delay == v[j].delay {
		return v[i].step > v[j].step
	} else {
		return v[i].delay < v[j].delay
	}
}

func (v metaTasks) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (psp *powerSavePlan) Update(screenSaverStartDelay, screenBlackDelay, sleepDelay int32) {
	psp.mu.Lock()
	defer psp.mu.Unlock()

	psp.interruptTasks()
	logger.Debugf("update(screenSaverStartDelay=%vs, screenBlackDelay=%vs, sleepDelay=%vs)",
		screenSaverStartDelay, screenBlackDelay, sleepDelay)

	tasks := make(metaTasks, 0, 3)
	if screenSaverStartDelay > 0 {
		tasks = append(tasks, metaTask{
			name:  "screenSaverStart",
			step:  0,
			delay: screenSaverStartDelay,
			fn:    psp.startScreensaver,
		})
	}

	if screenBlackDelay > 0 {
		tasks = append(tasks, metaTask{
			name:  "screenBlack",
			step:  1,
			delay: screenBlackDelay,
			fn:    psp.screenBlack,
		})
	}

	if sleepDelay > 0 {
		tasks = append(tasks, metaTask{
			name:  "sleep",
			step:  2,
			delay: sleepDelay,
			fn:    psp.makeSystemSleep,
		})
	}

	sort.Sort(tasks)
	psp.metaTasks = tasks

	if len(tasks) == 0 {
		psp.setScreenSaverTimeout(0)
		return
	}
	screenSaverTimeout := tasks[0].delay
	psp.setScreenSaverTimeout(screenSaverTimeout)

	for i := 0; i < len(tasks); i++ {
		a := tasks[i]

		for j := i + 1; j < len(tasks); j++ {
			b := tasks[j]
			if b.ignore {
				continue
			}

			if a.step > b.step {
				tasks[j].ignore = true
			}
		}
	}

	for _, t := range tasks {
		if t.ignore {
			logger.Debugf("task %s will not run", t.name)
		} else {
			logger.Debugf("task %s: %ds", t.name, t.delay)
		}
	}
}

func (psp *powerSavePlan) setScreenSaverTimeout(seconds int32) error {
	psp.screenSaverTimeout = seconds
	logger.Debugf("set ScreenSaver timeout to %d", seconds)
	err := psp.manager.helper.ScreenSaver.SetTimeout(0, uint32(seconds), 0, false)
	if err != nil {
		logger.Warningf("failed to set ScreenSaver timeout %d: %v", seconds, err)
	}
	return err
}

func (psp *powerSavePlan) saveCurrentBrightness() error {
	if psp.oldBrightnessTable == nil {
		var err error
		psp.oldBrightnessTable, err = psp.manager.helper.Display.Brightness().Get(0)
		if err != nil {
			return err
		}
		logger.Info("saveCurrentBrightness", psp.oldBrightnessTable)
		return nil
	}

	return errors.New("oldBrightnessTable is not nil")
}

func (psp *powerSavePlan) resetBrightness() {
	if psp.oldBrightnessTable != nil {
		logger.Debug("Reset all outputs brightness")
		psp.manager.setDisplayBrightness(psp.oldBrightnessTable)
		psp.oldBrightnessTable = nil
	}
}

func (psp *powerSavePlan) startScreensaver() {
	startScreensaver()
	psp.screensaverRunning = true
}

func (psp *powerSavePlan) stopScreensaver() {
	if !psp.screensaverRunning {
		return
	}
	stopScreensaver()
	psp.screensaverRunning = false
}

func (psp *powerSavePlan) makeSystemSleep() {
	psp.stopScreensaver()
	logger.Info("sleep")
	//psp.manager.setDPMSModeOn()
	psp.resetBrightness()
	psp.manager.doSuspend()
}

// 降低显示器亮度，最终关闭显示器
func (psp *powerSavePlan) screenBlack() {
	manager := psp.manager
	logger.Info("Start screen black")

	adjustBrightnessEnabled := manager.settings.GetBoolean(settingKeyAdjustBrightnessEnabled)

	if adjustBrightnessEnabled {
		err := psp.saveCurrentBrightness()
		if err != nil {
			adjustBrightnessEnabled = false
			logger.Warning(err)
		} else {
			// half black
			brightnessTable := make(map[string]float64)
			brightnessRatio := 0.5
			logger.Debug("brightnessRatio:", brightnessRatio)
			for output, oldBrightness := range psp.oldBrightnessTable {
				brightnessTable[output] = oldBrightness * brightnessRatio
			}
			manager.setDisplayBrightness(brightnessTable)
		}
	} else {
		logger.Debug("adjust brightness disabled")
	}

	// full black
	const fullBlackTime = 5000 * time.Millisecond
	taskF := newDelayedTask("screenFullBlack", fullBlackTime, func() {
		psp.stopScreensaver()
		logger.Info("Screen full black")
		if manager.ScreenBlackLock.Get() {
			manager.lockWaitShow(2 * time.Second)
		}

		if adjustBrightnessEnabled {
			// set min brightness for all outputs
			brightnessTable := make(map[string]float64)
			for output := range psp.oldBrightnessTable {
				brightnessTable[output] = 0.02
			}
			manager.setDisplayBrightness(brightnessTable)
		}
		manager.setDPMSModeOff()

	})
	psp.addTask(taskF)
}

func (psp *powerSavePlan) shouldPreventIdle() (bool, error) {
	conn := psp.manager.helper.xConn
	activeWin, err := ewmh.GetActiveWindow(conn).Reply(conn)
	if err != nil {
		return false, err
	}

	isFullscreenAndFocused, err := psp.isWindowFullScreenAndFocused(activeWin)
	if err != nil {
		return false, err
	}

	if !isFullscreenAndFocused {
		return false, nil
	}

	pid, err := ewmh.GetWMPid(conn, activeWin).Reply(conn)
	if err != nil {
		return false, err
	}

	p := procfs.Process(pid)
	cmdline, err := p.Cmdline()
	if err != nil {
		return false, err
	}

	for _, arg := range cmdline {
		for _, app := range psp.fullscreenWorkaroundAppList {
			if strings.Contains(arg, app) {
				logger.Debugf("match %q", app)
				return true, nil
			}
		}
	}
	return false, nil
}

// 开始 Idle
func (psp *powerSavePlan) HandleIdleOn() {
	psp.mu.Lock()
	defer psp.mu.Unlock()

	// Somethimes will emit idle on event when sleep, why?
	if psp.idleOn || psp.manager.getPrepareSuspend() {
		logger.Info("HandleIdleOn : IGNORE =========")
		psp.idleOn = true
		return
	}

	if isActive, err := psp.manager.isX11SessionActive(); err != nil {
		logger.Warning(err)
		return
	} else if !isActive {
		logger.Info("X11 session is inactive, don't HandleIdleOn")
		return
	}

	// check window
	preventIdle, err := psp.shouldPreventIdle()
	if err != nil {
		logger.Warning(err)
	}
	if preventIdle {
		logger.Debug("prevent idle")
		err := psp.manager.helper.ScreenSaver.SimulateUserActivity(0)
		if err != nil {
			logger.Warning(err)
		}
		return
	}
	psp.idleOn = true

	logger.Info("HandleIdleOn")
	for _, t := range psp.metaTasks {
		if t.ignore {
			continue
		}

		var delay time.Duration
		delaySeconds := t.delay - psp.screenSaverTimeout
		if delaySeconds == 0 {
			delay = time.Millisecond
		} else {
			delay = time.Duration(delaySeconds) * time.Second
		}
		logger.Debugf("do %s after %v", t.name, delay)
		task := newDelayedTask(t.name, delay, t.fn)
		psp.addTaskNoLock(task)
	}
}

// 结束 Idle
func (psp *powerSavePlan) HandleIdleOff() {
	psp.mu.Lock()
	defer psp.mu.Unlock()

	if !psp.idleOn || psp.manager.getPrepareSuspend() {
		logger.Info("HandleIdleOff : IGNORE =========")
		psp.idleOn = false
		return
	}
	psp.idleOn = false
	logger.Info("HandleIdleOff")
	psp.interruptTasks()
	psp.manager.setDPMSModeOn()
	psp.resetBrightness()
}

func (psp *powerSavePlan) isWindowFullScreenAndFocused(xid x.Window) (bool, error) {
	conn := psp.manager.helper.xConn
	states, err := ewmh.GetWMState(conn, xid).Reply(conn)
	if err != nil {
		return false, err
	}
	found := 0
	for _, s := range states {
		if s == psp.atomNetWMStateFullscreen {
			found++
		} else if s == psp.atomNetWMStateFocused {
			found++
		}
		if found == 2 {
			return true, nil
		}
	}
	return false, nil
}
