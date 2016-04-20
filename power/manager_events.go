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
	"pkg.deepin.io/dde/api/soundutils"
	. "pkg.deepin.io/lib/gettext"
	"time"
)

// 按下电源键
func (m *Manager) initPowerButtonEventHandler() {
	m.mediaKey.ConnectPowerOff(func(press bool) {
		if press {
			logger.Debug("PowerButton pressed")
			cmd := m.PowerButtonAction.Get()
			execCommand(cmd)
		}
	})
}

// 处理开启关闭笔记本盖子事件
func (m *Manager) initLidSwitchEventHandler() {
	upower := m.upower
	upower.LidIsClosed.ConnectChanged(func() {
		m.handleLidSwitch(upower.LidIsClosed.Get())
	})
}

// 处理有线电源插入拔出事件
func (m *Manager) initOnBatteryChangedHandler() {
	upower := m.upower
	upower.OnBattery.ConnectChanged(func() {
		logger.Debug("property OnBattery changed")
		m.setPropOnBattery(upower.OnBattery.Get())
		m.updateBatteryGroupInfo()

		if m.OnBattery {
			playSound(soundutils.EventPowerUnplug)
		} else {
			playSound(soundutils.EventPowerPlug)
		}
	})
}

func (m *Manager) handleLidSwitch(closed bool) {
	if closed {
		logger.Debug("Lid closed")
		// 为了多显示器时关闭笔记本盖子不待机
		if m.isMultiScreen() &&
			m.settings.GetBoolean(settingKeyMultiScreenPreventLidClosedExec) {
			logger.Info("Multi-Screen not exec")
			return
		}
		cmd := m.LidClosedAction.Get()
		execCommand(cmd)
	} else {
		logger.Debug("Lid opened")
		// 盖上盖子后如果待机了，打开盖子将会唤醒机器
	}
}

func (m *Manager) handleBeforeSuspend() {
	logger.Debug("before sleep")
	m.isSuspending = true

	if m.SleepLock.Get() || m.ScreenBlackLock.Get() {
		logger.Debug("* lock screen: DPMS on")
		m.setDPMSModeOn()
		// 此时执行 setDPMSModeOn() 将触发 HandleIdleOff
		// 如果 DPMS 是 off 状态，lock screen 将无法完成
		m.lockWaitShow(2 * time.Second)
	}
}

func (m *Manager) handleWeakup() {
	logger.Debug("weakup")
	m.isSuspending = false
	logger.Debug("Simulate user activity")
	m.screenSaver.SimulateUserActivity()

	m.checkBatteryPowerLevel()
	playSound(soundutils.EventWakeup)
	if m.SleepLock.Get() || m.ScreenBlackLock.Get() {
		m.doLock()
	}
}

func (m *Manager) initEventHandle() {
	m.initPowerButtonEventHandler()
	m.initLidSwitchEventHandler()
	m.initOnBatteryChangedHandler()
}

func (m *Manager) checkBatteryPowerLevel() {
	logger.Debug("checkBatteryPowerLevel")

	newLevel := m.getNewBatteryPowerLevel()
	if m.batteryPowerLevel != newLevel {
		logger.Debugf("batteryPowerLevel changed: %v => %v",
			getBatteryPowerLevelName(m.batteryPowerLevel), getBatteryPowerLevelName(newLevel))
		m.batteryPowerLevel = newLevel
		m.handleBatteryPowerLevelChanged()
	} else {
		logger.Debug("batteryPowerLevel no changed:", getBatteryPowerLevelName(m.batteryPowerLevel))
	}
}

func (m *Manager) disableSecondTicker() {
	if m.secondTicker != nil {
		m.secondTicker.Stop()
		m.secondTicker = nil
	}
}

func (m *Manager) handleBatteryPowerLevelChanged() {
	logger.Debug("handleBatteryPowerLevelChanged")
	m.disableSecondTicker()

	switch m.batteryPowerLevel {
	case batteryPowerLevelAbnormal:
		m.sendNotify("battery_empty", Tr("Abnormal battery power"), Tr("Battery power can not be predicted, please save important documents properly and  not do important operations."))

	case batteryPowerLevelExhausted:
		playSound(soundutils.EventBatteryLow)
		m.sendNotify("battery_empty", Tr("Battery Critical Low"), Tr("Computer has been in suspend mode, please plug in."))
		m.secondTicker = newSecondTicker(func(count int) {
			if count == 3 {
				// after 3 seconds, lock and then show lowpower
				go func() {
					if m.SleepLock.Get() || m.ScreenBlackLock.Get() {
						m.lockWaitShow(2 * time.Second)
					}
					doShowLowpower()
				}()
			} else if count == 5 {
				// after 5 seconds, force suspend
				m.disableSecondTicker()
				m.doSuspend()
			}
		})

	case batteryPowerLevelVeryLow:
		m.secondTicker = newSecondTicker(func(count int) {
			// notify every 60 seconds
			if count%60 == 0 {
				playSound(soundutils.EventBatteryLow)
				m.sendNotify("battery_empty", Tr("Battery Critical Low"), Tr("Computer has been in suspend mode, please plug in."))
			}
		})

	case batteryPowerLevelLow:
		playSound(soundutils.EventBatteryLow)
		m.sendNotify("battery_caution", Tr("Battery Low"), Tr("Computer will be in suspend mode, please plug in now."))

	case batteryPowerLevelSufficient:
		logger.Debug("Power sufficient")
		doCloseLowpower()
		// 由 低电量 到 电量充足，必然需要有线电源插入
	}

}

func (m *Manager) getBatteryGroupTimeToEmpty() int64 {
	var timeToEmpty int64
	for _, batInfo := range m.batteryGroup.InfoMap {
		timeToEmpty += batInfo.TimeToEmpty
	}
	return timeToEmpty
}

func (m *Manager) getBatteryGroupPercentage() float64 {
	var energySum, energyFullSum float64
	for _, batInfo := range m.batteryGroup.InfoMap {
		energySum += batInfo.Energy
		energyFullSum += batInfo.EnergyFull
	}
	return (energySum / energyFullSum) * 100.0
}

func getBatteryPowerLevelByPercentage(percentage float64) uint32 {
	switch {
	case percentage < batteryPercentageAbnormal:
		return batteryPowerLevelAbnormal

	case percentage <= batteryPercentageExhausted:
		return batteryPowerLevelExhausted

	case percentage <= batteryPercentageVeryLow:
		return batteryPowerLevelVeryLow

	case percentage <= batteryPercentageLow:
		return batteryPowerLevelLow

	default:
		return batteryPowerLevelSufficient
	}
}

func getBatteryPowerLevelByTimeToEmpty(time int64) uint32 {
	switch {
	case time < timeToEmptyAbnormal:
		return batteryPowerLevelAbnormal

	case time < timeToEmptyExhausted:
		return batteryPowerLevelExhausted

	case time < timeToEmptyVeryLow:
		return batteryPowerLevelVeryLow

	case time < timeToEmptyLow:
		return batteryPowerLevelLow

	default:
		return batteryPowerLevelSufficient
	}
}

func (m *Manager) getNewBatteryPowerLevel() uint32 {
	if !m.OnBattery {
		return batteryPowerLevelSufficient
	}

	if m.usePercentageForPolicy {
		percentage := m.getBatteryGroupPercentage()
		logger.Debugf("sum percentage: %.2f%%", percentage)
		return getBatteryPowerLevelByPercentage(percentage)
	} else {
		// use time to empty for policy
		timeToEmpty := m.getBatteryGroupTimeToEmpty()
		logger.Debug("sum timeToEmpty(secs):", timeToEmpty)
		powerLevel := getBatteryPowerLevelByTimeToEmpty(timeToEmpty)
		if powerLevel == batteryPowerLevelAbnormal {
			logger.Debug("Try use percentage for policy")
			percentage := m.getBatteryGroupPercentage()
			logger.Debugf("sum percentage: %.2f%%", percentage)
			return getBatteryPowerLevelByPercentage(percentage)
		}
		return powerLevel
	}
}
