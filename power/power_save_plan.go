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
	"math"
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
	manager           *Manager
	screenBlackTicker *countTicker
	sleepDelay        int32
	sleepExit         chan int
	oldBrightness     float64
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
	if psp.screenBlackTicker != nil {
		psp.screenBlackTicker.Stop()
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

	if screenBlackDelay > 5 {
		screenBlackDelay -= 5
	} else if screenBlackDelay > 0 {
		// 0 < screenBlackDelay <= 5
		screenBlackDelay = 1
	}

	logger.Debug("screen saver set timeout: ", screenBlackDelay)
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
// 以1分钟关闭显示器为例， 55秒的时候屏幕黑到 50% (以动画的方式从 100% 降到 50%, 动画耗时 500ms, 动画曲线是加速曲线), 5秒钟以后没有用户操作， 关闭显示器。
// 关闭显示器之后，开始休眠记时
func (psp *powerSavePlan) screenBlack() {
	logger.Info("Start screen black")
	psp.saveCurrentBrightness()

	if psp.screenBlackTicker != nil {
		psp.screenBlackTicker.Reset()
		return
	}

	psp.screenBlackTicker = newCountTicker(time.Millisecond*10, func(count int) {
		manager := psp.manager
		if 0 <= count && count <= 50 {
			// 0ms ~ 500ms
			// count: [0,50], brightness ratio:  1 ~ 0.5
			brightnessRatio := 0.5 * (math.Cos(math.Pi*float64(count)/100) + 1)
			manager.setDisplayBrightness(psp.oldBrightness * brightnessRatio)
		} else if count == 500 {
			// 5000ms 5s
			psp.screenBlackTicker.Stop()
			go func() {
				if manager.ScreenBlackLock.Get() {
					manager.lockWaitShow(2 * time.Second)
				}
				// set min brightness
				manager.setDisplayBrightness(0.02)
				manager.setDPMSModeOff()

				// try sleep
				if psp.sleepDelay != 0 {
					psp.sleepExit = make(chan int)
					psp.sleep()
				}

			}()
		}
	})
}

// 等待 psp.sleepDelay 后待机，可被 psp.sleepExit 中断
func (psp *powerSavePlan) sleep() {
	select {
	case <-time.After(time.Second * time.Duration(psp.sleepDelay)):
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
	psp.screenBlack()
}

// 结束 Idle
func (psp *powerSavePlan) HandleIdleOff() {
	if psp.manager.isSuspending {
		logger.Debug("Suspending NOT HandleIdleOff")
		return
	}
	logger.Debug("HandleIdleOff")
	psp.interruptAll()
	time.AfterFunc(50*time.Millisecond, func() {
		// reset brightness
		if psp.oldBrightness != 0.0 {
			psp.manager.setDPMSModeOn()
			logger.Debug("Reset all outputs brightness")
			psp.manager.setDisplayBrightness(psp.oldBrightness)
			psp.oldBrightness = 0.0
		}
	})
}
