/*
 * Copyright (C) 2016 ~ 2017 Deepin Technology Co., Ltd.
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
	"gir/gio-2.0"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/dbus/property"
)

const (
	gsSchemaPower = "com.deepin.dde.power"
)

type Manager struct {
	helper               *Helper
	settings             *gio.Settings
	isSuspending         bool
	warnLevelCountTicker *countTicker
	warnLevelConfig      *WarnLevelConfig
	submodules           map[string]submodule
	inhibitor            *sleepInhibitor

	// 接通电源时，不做任何操作，到关闭屏幕需要的时间
	LinePowerScreenBlackDelay *property.GSettingsIntProperty `access:"readwrite"`
	// 接通电源时，不做任何操作，从黑屏到睡眠的时间
	LinePowerSleepDelay *property.GSettingsIntProperty `access:"readwrite"`

	// 使用电池时，不做任何操作，到关闭屏幕需要的时间
	BatteryScreenBlackDelay *property.GSettingsIntProperty `access:"readwrite"`
	// 使用电池时，不做任何操作，从黑屏到睡眠的时间
	BatterySleepDelay *property.GSettingsIntProperty `access:"readwrite"`

	// 关闭显示器前是否锁定
	ScreenBlackLock *property.GSettingsBoolProperty `access:"readwrite"`
	// 睡眠前是否锁定
	SleepLock *property.GSettingsBoolProperty `access:"readwrite"`

	// 笔记本电脑盖上盖子后是否睡眠
	LidClosedSleep *property.GSettingsBoolProperty `access:"readwrite"`

	// 是否有盖子，一般笔记本电脑才有
	LidIsPresent bool
	// 是否使用电池, 接通电源时为 false, 使用电池时为 true
	OnBattery bool

	// 电池是否可用，是否存在
	BatteryIsPresent map[string]bool
	// 电池电量百分比
	BatteryPercentage map[string]float64
	// 电池状态
	BatteryState map[string]uint32

	// 警告级别
	WarnLevel WarnLevel
}

func NewManager() (*Manager, error) {
	m := &Manager{}
	helper, err := NewHelper()
	if err != nil {
		return nil, err
	}
	m.helper = helper

	m.settings = gio.NewSettings(gsSchemaPower)

	// warn level config
	m.warnLevelConfig = NewWarnLevelConfig()
	m.warnLevelConfig.connectSettings(m.settings)
	err = dbus.InstallOnSession(m.warnLevelConfig)
	if err != nil {
		m.destroy()
		return nil, err
	}

	m.LinePowerScreenBlackDelay = property.NewGSettingsIntProperty(m, "LinePowerScreenBlackDelay", m.settings, settingKeyLinePowerScreenBlackDelay)
	m.LinePowerSleepDelay = property.NewGSettingsIntProperty(m, "LinePowerSleepDelay", m.settings, settingKeyLinePowerSleepDelay)
	m.BatteryScreenBlackDelay = property.NewGSettingsIntProperty(m, "BatteryScreenBlackDelay", m.settings, settingKeyBatteryScreenBlackDelay)
	m.BatterySleepDelay = property.NewGSettingsIntProperty(m, "BatterySleepDelay", m.settings, settingKeyBatterySleepDelay)

	m.ScreenBlackLock = property.NewGSettingsBoolProperty(m, "ScreenBlackLock", m.settings, settingKeyScreenBlackLock)
	m.SleepLock = property.NewGSettingsBoolProperty(m, "SleepLock", m.settings, settingKeySleepLock)

	m.LidClosedSleep = property.NewGSettingsBoolProperty(m, "LidClosedSleep", m.settings, settingKeyLidClosedSleep)

	power := m.helper.Power
	m.LidIsPresent = power.HasLidSwitch.Get()
	m.OnBattery = power.OnBattery.Get()
	logger.Info("LidIsPresent", m.LidIsPresent)

	// init battery display
	m.BatteryIsPresent = make(map[string]bool)
	m.BatteryPercentage = make(map[string]float64)
	m.BatteryState = make(map[string]uint32)

	logger.Info("NewManager done")
	return m, nil
}

func (m *Manager) init() {
	// init sleep inhibitor
	m.inhibitor = newSleepInhibitor(m.helper.Login1Manager)
	m.inhibitor.OnBeforeSuspend = m.handleBeforeSuspend
	m.inhibitor.OnWakeup = m.handleWakeup
	m.inhibitor.block()

	m.handleBatteryDisplayUpdate()
	m.initBatteryDisplayUpdateHandler()

	m.initPowerModule()

	m.initOnBatteryChangedHandler()
	m.initSubmodules()
	m.startSubmodules()

	m.StartupNotify()
}

func (m *Manager) StartupNotify() {
	props := []string{"BatteryIsPresent", "BatteryPercentage", "BatteryState",
		"WarnLevel", "OnBattery", "LidIsPresent"}
	for _, propName := range props {
		dbus.NotifyChange(m, propName)
	}
}

func (m *Manager) initPowerModule() {
	inited := m.settings.GetBoolean(settingKeyPowerModuleInitialized)
	if !inited {
		// TODO: 也许有更好的判断台式机的方法
		power := m.helper.Power
		if !power.HasBattery.Get() {
			// 无电池，判断为台式机, 设置待机为 从不
			m.LinePowerSleepDelay.Set(0)
			m.BatterySleepDelay.Set(0)
		}
		m.settings.SetBoolean(settingKeyPowerModuleInitialized, true)
	}
}

func (m *Manager) isX11SessionActive() (bool, error) {
	return m.helper.SessionWatcher.IsX11SessionActive()
}

func (m *Manager) destroy() {
	m.destroySubmodules()

	if m.helper != nil {
		m.helper.Destroy()
		m.helper = nil
	}
	if m.warnLevelConfig != nil {
		dbus.UnInstallObject(m.warnLevelConfig)
		m.warnLevelConfig = nil
	}

	if m.inhibitor != nil {
		m.inhibitor.unblock()
		m.inhibitor = nil
	}
	dbus.UnInstallObject(m)
}
