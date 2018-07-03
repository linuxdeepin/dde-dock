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
	"strings"
	"sync"
	"time"

	"github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/util/wm/ewmh"
	"pkg.deepin.io/lib/gsettings"
	"pkg.deepin.io/lib/procfs"
)

const submodulePSP = "PowerSavePlan"

func init() {
	submoduleList = append(submoduleList, newPowerSavePlan)
}

type powerSavePlan struct {
	manager       *Manager
	tasks         TimeAfterTasks
	doScreenBlack bool
	sleepTimeout  time.Duration
	// key output name, value old brightness
	oldBrightnessTable map[string]float64
	mu                 sync.Mutex
	idleOn             bool

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
		if key == settingKeyLinePowerScreenBlackDelay ||
			key == settingKeyLinePowerSleepDelay {
			if !m.OnBattery {
				logger.Debug("Change OnLinePower plan")
				psp.OnLinePower()
			}
		} else if key == settingKeyBatteryScreenBlackDelay ||
			key == settingKeyBatterySleepDelay {
			if m.OnBattery {
				logger.Debug("Change OnBattery plan")
				psp.OnBattery()
			}
		}
	})
}

func (psp *powerSavePlan) OnBattery() {
	logger.Debug("Use OnBattery plan")
	m := psp.manager
	psp.Update(m.BatteryScreenBlackDelay.Get(), m.BatterySleepDelay.Get())
}

func (psp *powerSavePlan) OnLinePower() {
	logger.Debug("Use OnLinePower plan")
	m := psp.manager
	psp.Update(m.LinePowerScreenBlackDelay.Get(), m.LinePowerSleepDelay.Get())
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
	psp.manager.setDisplayBrightness(brightnessTable)
}

// 取消之前的任务
func (psp *powerSavePlan) interruptTasks() {
	psp.tasks.CancelAll()
	logger.Info("cancel all tasks")
	psp.tasks.Wait(10*time.Millisecond, 200)
	logger.Info("all tasks done!")
}

func (psp *powerSavePlan) Destroy() {
	psp.interruptTasks()
}

/* 更新计划
screenBlackDelay == 0 从不关闭显示屏
sleepDelay == 0 从不待机
*/
func (psp *powerSavePlan) Update(screenBlackDelay, sleepDelay int32) error {
	psp.mu.Lock()
	defer psp.mu.Unlock()

	psp.interruptTasks()
	logger.Debugf("update(screenBlackDelay=%vs, sleepDelay=%vs)",
		screenBlackDelay, sleepDelay)

	var screenSaverTimeout int32 // seconds
	const immediately = time.Millisecond
	const never = time.Duration(0)

	if screenBlackDelay == 0 && sleepDelay != 0 {
		// ex. screenBlack never , sleep 1min
		screenSaverTimeout = sleepDelay
		psp.doScreenBlack = false
		psp.sleepTimeout = immediately
	} else if screenBlackDelay == 0 && sleepDelay == 0 {
		// screenBlack never, sleep never
		screenSaverTimeout = 0 // never
		psp.doScreenBlack = false
		psp.sleepTimeout = never
	} else if screenBlackDelay != 0 && sleepDelay != 0 {
		if screenBlackDelay >= sleepDelay {
			// ex. screenBlack 3min, sleep 1min
			// means no screenBlack
			screenSaverTimeout = sleepDelay
			psp.doScreenBlack = false
			psp.sleepTimeout = immediately
		} else {
			// screenBlackDelay < sleepDelay
			// ex. screenBlack 1min, sleep 3min
			screenSaverTimeout = screenBlackDelay
			psp.doScreenBlack = true
			psp.sleepTimeout = time.Duration(sleepDelay-screenBlackDelay) * time.Second
		}
	} else {
		// is screenBlackDelay != 0 && sleepDelay == 0
		// ex. screenBlack 1min, sleep never
		screenSaverTimeout = screenBlackDelay
		psp.doScreenBlack = true
		psp.sleepTimeout = never
	}

	psp.setScreenSaverTimeout(screenSaverTimeout)
	logger.Debugf("doScreenBlack %v, sleepTimeout %v", psp.doScreenBlack, psp.sleepTimeout)
	return nil
}

func (psp *powerSavePlan) setScreenSaverTimeout(seconds int32) error {
	logger.Debugf("set ScreenSaver timeout to %d", seconds)
	err := psp.manager.helper.ScreenSaver.SetTimeout(0, uint32(seconds), 0, false)
	if err != nil {
		logger.Warningf("set ScreenSaver timeout %d failed: %v", seconds, err)
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

// 降低显示器亮度，最终关闭显示器
func (psp *powerSavePlan) screenBlack() {
	manager := psp.manager
	logger.Info("Start screen black")
	psp.tasks = make(TimeAfterTasks, 0)

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
	taskF := NewTimeAfterTask(fullBlackTime, func() {
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
	psp.tasks = append(psp.tasks, taskF)
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

	if psp.idleOn {
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
	if psp.doScreenBlack {
		psp.screenBlack()
	}

	if psp.sleepTimeout > 0 {
		logger.Infof("sleep after %v", psp.sleepTimeout)
		taskS := NewTimeAfterTask(psp.sleepTimeout, func() {
			logger.Infof("sleep")
			psp.manager.setDPMSModeOn()
			psp.resetBrightness()
			psp.manager.doSuspend()
		})
		psp.tasks = append(psp.tasks, taskS)
	}
}

// 结束 Idle
func (psp *powerSavePlan) HandleIdleOff() {
	psp.mu.Lock()
	defer psp.mu.Unlock()

	if !psp.idleOn {
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
