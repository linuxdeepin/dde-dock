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
	"pkg.deepin.io/dde/api/soundutils"
	"pkg.deepin.io/lib/dbus"
	. "pkg.deepin.io/lib/gettext"
	"time"
)

// 处理有线电源插入拔出事件
func (m *Manager) initOnBatteryChangedHandler() {
	power := m.helper.Power
	power.OnBattery.ConnectChanged(func() {
		logger.Debug("property OnBattery changed")
		m.setPropOnBattery(power.OnBattery.Get())

		if m.OnBattery {
			playSound(soundutils.EventPowerUnplug)
		} else {
			playSound(soundutils.EventPowerPlug)
		}
	})
}

func (m *Manager) handleBeforeSuspend() {
	logger.Debug("before sleep")
	if m.SleepLock.Get() || m.ScreenBlackLock.Get() {
		m.setDPMSModeOn()
		m.lockWaitShow(4 * time.Second)
	}
	m.setDPMSModeOff()
}

func (m *Manager) handleWakeup() {
	logger.Debug("wakeup")
	m.helper.Power.RefreshBatteries()
	playSound(soundutils.EventWakeup)
	m.setDPMSModeOn()
}

func (m *Manager) initBatteryDisplayUpdateHandler() {
	power := m.helper.Power
	power.ConnectBatteryDisplayUpdate(func(timestamp int64) {
		logger.Debug("BatteryDisplayUpdate", timestamp)
		m.handleBatteryDisplayUpdate()
	})

	m.warnLevelConfig.setChangeCallback(m.handleBatteryDisplayUpdate)
}

func (m *Manager) handleBatteryDisplayUpdate() {
	logger.Debug("handleBatteryDisplayUpdate")
	power := m.helper.Power
	hasBattery := power.HasBattery.Get()
	if hasBattery {
		m.setPropBatteryIsPresent(true)
		percentage := power.BatteryPercentage.Get()
		timeToEmpty := power.BatteryTimeToEmpty.Get()
		m.setPropBatteryPercentage(percentage)
		m.setPropBatteryState(power.BatteryStatus.Get())
		warnLevel := m.getWarnLevel(percentage, timeToEmpty)
		m.setPropWarnLevel(warnLevel)

	} else {
		m.setPropWarnLevel(WarnLevelNone)
		delete(m.BatteryIsPresent, batteryDisplay)
		delete(m.BatteryPercentage, batteryDisplay)
		delete(m.BatteryState, batteryDisplay)
		dbus.NotifyChange(m, "BatteryIsPresent")
		dbus.NotifyChange(m, "BatteryPercentage")
		dbus.NotifyChange(m, "BatteryState")
	}
}

func (m *Manager) disableWarnLevelCountTicker() {
	if m.warnLevelCountTicker != nil {
		m.warnLevelCountTicker.Stop()
		m.warnLevelCountTicker = nil
	}
}

func (m *Manager) handleWarnLevelChanged() {
	logger.Debug("handleWarnLevelChanged")
	m.disableWarnLevelCountTicker()

	switch m.WarnLevel {
	case WarnLevelAction:
		playSound(soundutils.EventBatteryLow)
		m.sendNotify("battery_empty", Tr("Battery Critical Low"), Tr("Computer has been in suspend mode, please plug in"))
		m.warnLevelCountTicker = newCountTicker(time.Second, func(count int) {
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
				m.disableWarnLevelCountTicker()
				m.doSuspend()
			}
		})

	case WarnLevelCritical:
		m.warnLevelCountTicker = newCountTicker(time.Second, func(count int) {
			// notify every 60 seconds
			if count%60 == 0 {
				playSound(soundutils.EventBatteryLow)
				m.sendNotify("battery_low", Tr("Battery Critical Low"), Tr("Computer has been in suspend mode, please plug in"))
			}
		})

	case WarnLevelLow:
		playSound(soundutils.EventBatteryLow)
		m.sendNotify("battery_caution", Tr("Battery Low"), Tr("Computer will be in suspend mode, please plug in now"))

	case WarnLevelNone:
		logger.Debug("Power sufficient")
		doCloseLowpower()
		// 由 低电量 到 电量充足，必然需要有线电源插入
	}
}
