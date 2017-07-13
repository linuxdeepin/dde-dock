/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package power

import (
	"gir/gio-2.0"
	"sync"
	"time"
)

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
}

func newPowerSavePlan(manager *Manager) (string, submodule, error) {
	p := new(powerSavePlan)
	p.manager = manager
	return "PowerSavePlan", p, nil
}

// 监听 GSettings 值改变, 更新节电计划
func (psp *powerSavePlan) initSettingsChangedHandler() {
	m := psp.manager
	m.settings.Connect("changed", func(s *gio.Settings, key string) {
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
	power.OnBattery.ConnectChanged(psp.Reset)

	screenSaver.ConnectIdleOn(psp.HandleIdleOn)
	screenSaver.ConnectIdleOff(psp.HandleIdleOff)
	return nil
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
	err := psp.manager.helper.ScreenSaver.SetTimeout(uint32(seconds), 0, false)
	if err != nil {
		logger.Warningf("set ScreenSaver timeout %d failed: %v", seconds, err)
	}
	return err
}

func (psp *powerSavePlan) saveCurrentBrightness() {
	if psp.oldBrightnessTable == nil {
		psp.oldBrightnessTable = psp.manager.helper.Display.Brightness.Get()
		logger.Info("saveCurrentBrightness", psp.oldBrightnessTable)
	} else {
		logger.Debug("saveCurrentBrightness failed")
	}
}

func (psp *powerSavePlan) resetBrightness() {
	if psp.oldBrightnessTable != nil {
		psp.manager.setDPMSModeOn()
		logger.Debug("Reset all outputs brightness")
		psp.manager.setDisplayBrightness(psp.oldBrightnessTable)
		psp.oldBrightnessTable = nil
	}
}

// 降低显示器亮度，最终关闭显示器
func (psp *powerSavePlan) screenBlack() {
	manager := psp.manager
	logger.Info("Start screen black")
	psp.saveCurrentBrightness()
	psp.tasks = make(TimeAfterTasks, 0)

	// half black
	{
		brightnessTable := make(map[string]float64)
		brightnessRatio := 0.5
		logger.Debug("brightnessRatio:", brightnessRatio)
		for output, oldBrightness := range psp.oldBrightnessTable {
			brightnessTable[output] = oldBrightness * brightnessRatio
		}
		manager.setDisplayBrightness(brightnessTable)
	}

	// full black
	const fullBlackTime = 5000 // ms
	taskF := NewTimeAfterTask(fullBlackTime*time.Millisecond, func() {
		logger.Info("Screen full black")
		if manager.ScreenBlackLock.Get() {
			manager.lockWaitShow(2 * time.Second)
		}
		// set min brightness for all outputs
		brightnessTable := make(map[string]float64)
		for output, _ := range psp.oldBrightnessTable {
			brightnessTable[output] = 0.02
		}
		manager.setDisplayBrightness(brightnessTable)
		manager.setDPMSModeOff()

	})
	psp.tasks = append(psp.tasks, taskF)
}

// 开始 Idle
func (psp *powerSavePlan) HandleIdleOn() {
	psp.mu.Lock()
	defer psp.mu.Unlock()

	if psp.manager.isSuspending {
		logger.Info("Suspending NOT HandleIdleOn")
		return
	}

	if isActive, err := psp.manager.isX11SessionActive(); err != nil {
		logger.Warning(err)
		return
	} else if !isActive {
		logger.Info("X11 session is inactive, don't HandleIdleOn")
		return
	}

	logger.Info("HandleIdleOn")
	if psp.doScreenBlack {
		psp.screenBlack()
	}

	if psp.sleepTimeout > 0 {
		logger.Infof("sleep after %v", psp.sleepTimeout)
		taskS := NewTimeAfterTask(psp.sleepTimeout, func() {
			logger.Infof("sleep")
			psp.manager.doSuspend()
		})
		psp.tasks = append(psp.tasks, taskS)
	}
}

// 结束 Idle
func (psp *powerSavePlan) HandleIdleOff() {
	psp.mu.Lock()
	defer psp.mu.Unlock()

	if psp.manager.isSuspending {
		logger.Info("Suspending NOT HandleIdleOff")
		return
	}
	logger.Info("HandleIdleOff")
	psp.interruptTasks()
	psp.resetBrightness()
}
