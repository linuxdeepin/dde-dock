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
	"time"
)

const (
	dbusDisplayDest = "com.deepin.daemon.Display"
	dbusDisplayPath = "/com/deepin/daemon/Display"
)

func init() {
	submoduleList = append(submoduleList, newPowerSavePlan)
}

type powerSavePlan struct {
	manager         *Manager
	sleepDelay      int32
	screenBlackExit chan int
	sleepExit       chan int
	oldBrightness   float64
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

	//OnBattery changed will effect current PowerSavePlan
	psp.manager.upower.OnBattery.ConnectChanged(psp.Reset)

	screenSaver := psp.manager.screenSaver
	screenSaver.ConnectIdleOn(psp.HandleIdleOn)
	screenSaver.ConnectIdleOff(psp.HandleIdleOff)
	return nil
}

// 取消之前的计划
func (psp *powerSavePlan) interruptScreenBlack() {
	if psp.screenBlackExit != nil {
		close(psp.screenBlackExit)
		psp.screenBlackExit = nil
	}
}

func (psp *powerSavePlan) interruptSleep() {
	if psp.sleepExit != nil {
		close(psp.sleepExit)
		psp.sleepExit = nil
	}
}

func (psp *powerSavePlan) interruptAll() {
	psp.interruptScreenBlack()
	psp.interruptSleep()
}

func (psp *powerSavePlan) Destroy() {
	psp.interruptAll()
}

/* 更新计划
screenBlackDelay == 0 从不关闭显示屏
sleepDelay == 0 从不待机
*/
func (psp *powerSavePlan) Update(screenBlackDelay, sleepDelay int32) error {
	psp.interruptAll()
	logger.Debugf("update(screenBlackDelay=%vs, sleepDelay=%vs)",
		screenBlackDelay, sleepDelay)

	if err := psp.manager.screenSaver.SetTimeout(uint32(screenBlackDelay), 0, false); err != nil {
		logger.Errorf("Failed set ScreenSaver's timeout %v : %v", screenBlackDelay, err)
		return err
	}
	psp.sleepDelay = sleepDelay
	return nil
}

// 获取当前亮度并保存到 psp.oldBrightness
// 假设所有显示器亮度都一样，只获取第一个
func (psp *powerSavePlan) saveCurrentBrightness() {
	for _, br := range psp.manager.display.Brightness.Get() {
		psp.oldBrightness = br
		logger.Debug("Save brightness:", br)
		break
	}
}

// 逐渐降低显示器亮度，最终关闭显示器
// 可被 psp.screenBlackExit 中断
func (psp *powerSavePlan) screenBlack() {
	psp.saveCurrentBrightness()
	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()
	const brightnessDelta = 0.05
	b := psp.oldBrightness - brightnessDelta
	for {
		select {
		case <-psp.screenBlackExit:
			logger.Debug("Interrupt screen black")
			return
		case <-ticker.C:
			if b >= 0.0 {
				psp.manager.setDisplayBrightness(b)
				b -= brightnessDelta
			} else {
				manager := psp.manager
				manager.setDPMSModeOff()
				if manager.ScreenBlackLock.Get() {
					manager.doLock()
				}
				return
			}
		}
	}
}

// 等待 psp.sleepDelay 后待机，可被 psp.sleepExit 中断
func (psp *powerSavePlan) sleep() {
	select {
	case <-time.After(time.Second * time.Duration(psp.sleepDelay)):
		// 打断有可能还在的 idle 过程
		psp.interruptScreenBlack()
		psp.manager.doSuspend()
	case <-psp.sleepExit:
		logger.Debugf("Interrupt sleep")
	}
}

// 开始 Idle
func (psp *powerSavePlan) HandleIdleOn() {
	if psp.manager.isSuspending {
		logger.Debug("Suspending NOT HandleIdleOn")
		return
	}

	logger.Debug("HandleIdleOn")
	psp.screenBlackExit = make(chan int)
	go psp.screenBlack()

	// sleep
	if psp.sleepDelay != 0 {
		psp.sleepExit = make(chan int)
		go psp.sleep()
	}
}

// 结束 Idle
// stop idle/suspend action
func (psp *powerSavePlan) HandleIdleOff() {
	if psp.manager.isSuspending {
		logger.Debug("Suspending NOT HandleIdleOff")
		return
	}
	logger.Debug("HandleIdleOff")
	psp.interruptAll()
	// reset display
	if psp.oldBrightness != 0.0 {
		psp.manager.setDPMSModeOn()
		logger.Debug("Reset all outputs brightness")
		psp.manager.setDisplayBrightness(psp.oldBrightness)
		psp.oldBrightness = 0.0
	}
}
