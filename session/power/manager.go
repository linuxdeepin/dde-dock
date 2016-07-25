/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package power

import (
	"gir/gio-2.0"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/dbus/property"
)

type Manager struct {
	helper               *Helper
	settings             *gio.Settings
	isSuspending         bool
	warnLevelCountTicker *countTicker
	warnLevelConfig      *WarnLevelConfig
	submodules           map[string]submodule
	isSessionActive      bool
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

	// 按下电源按钮后执行的命令
	PowerButtonAction *property.GSettingsStringProperty `access:"readwrite"`
	// 笔记本电脑关闭盖子后执行的命令
	LidClosedAction *property.GSettingsStringProperty `access:"readwrite"`

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
	err := m.init()
	if err != nil {
		m.destroy()
		return nil, err
	}
	logger.Info("NewManager done")
	return m, nil
}

func (m *Manager) init() error {
	helper, err := NewHelper()
	if err != nil {
		return err
	}
	m.helper = helper

	m.settings = gio.NewSettings("com.deepin.dde.power")

	// warn level config
	m.warnLevelConfig = NewWarnLevelConfig()
	m.warnLevelConfig.connectSettings(m.settings)
	err = dbus.InstallOnSession(m.warnLevelConfig)
	if err != nil {
		return err
	}

	// init sleep inhibitor
	m.inhibitor = newSleepInhibitor(m.helper.Login1Manager)
	m.inhibitor.OnBeforeSuspend = m.handleBeforeSuspend
	m.inhibitor.OnWeakup = m.handleWeakup
	m.inhibitor.block()

	m.LinePowerScreenBlackDelay = property.NewGSettingsIntProperty(m, "LinePowerScreenBlackDelay", m.settings, settingKeyLinePowerScreenBlackDelay)
	m.LinePowerSleepDelay = property.NewGSettingsIntProperty(m, "LinePowerSleepDelay", m.settings, settingKeyLinePowerSleepDelay)
	m.BatteryScreenBlackDelay = property.NewGSettingsIntProperty(m, "BatteryScreenBlackDelay", m.settings, settingKeyBatteryScreenBlackDelay)
	m.BatterySleepDelay = property.NewGSettingsIntProperty(m, "BatterySleepDelay", m.settings, settingKeyBatterySleepDelay)

	m.ScreenBlackLock = property.NewGSettingsBoolProperty(m, "ScreenBlackLock", m.settings, settingKeyScreenBlackLock)
	m.SleepLock = property.NewGSettingsBoolProperty(m, "SleepLock", m.settings, settingKeySleepLock)

	m.PowerButtonAction = property.NewGSettingsStringProperty(m, "PowerButtonAction", m.settings, settingKeyPowerButtonPressedExec)
	m.LidClosedAction = property.NewGSettingsStringProperty(m, "LidClosedAction", m.settings, settingKeyLidClosedExec)

	power := m.helper.Power
	m.LidIsPresent = power.HasLidSwitch.Get()
	m.OnBattery = power.OnBattery.Get()
	logger.Info("LidIsPresent", m.LidIsPresent)

	// init battery display
	m.BatteryIsPresent = make(map[string]bool)
	m.BatteryPercentage = make(map[string]float64)
	m.BatteryState = make(map[string]uint32)
	m.handleBatteryDisplayUpdate()
	m.initBatteryDisplayUpdateHandler()

	m.initPowerModule()

	m.initPowerButtonEventHandler()
	m.initLidSwitchEventHandler()
	m.initOnBatteryChangedHandler()

	sessionWatcher := m.helper.SessionWatcher
	m.isSessionActive = sessionWatcher.IsActive.Get()
	m.initSessionActiveChangedHandler()

	m.initSubmodules()
	m.startSubmodules()

	err = dbus.InstallOnSession(m)
	if err != nil {
		return err
	}
	logger.Info("InstallOnSession done")
	m.StartupNotify()
	return nil
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

func (m *Manager) destroy() {
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
	m.destroySubmodules()
	dbus.UnInstallObject(m)
}
