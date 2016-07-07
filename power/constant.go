/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package power

// setting key in com.deepin.dde.power
const (
	settingKeyBatteryScreenBlackDelay = "battery-screen-black-delay"
	settingKeyBatterySleepDelay       = "battery-sleep-delay"

	settingKeyLinePowerScreenBlackDelay = "line-power-screen-black-delay"
	settingKeyLinePowerSleepDelay       = "line-power-sleep-delay"

	settingKeyScreenBlackLock        = "screen-black-lock"
	settingKeySleepLock              = "sleep-lock"
	settingKeyLidClosedExec          = "lid-closed-exec"
	settingKeyPowerButtonPressedExec = "power-button-pressed-exec"

	settingKeyFullscreenWorkaroundEnabled     = "fullscreen-workaround-enabled"
	settingKeyMultiScreenPreventLidClosedExec = "multi-screen-prevent-lid-closed-exec"
	settingKeyUsePercentageForPolicy          = "use-percentage-for-policy"

	settingKeyPowerModuleInitialized = "power-module-initialized"
	settingKeyPercentageLow          = "percentage-low"
	settingKeyPercentageVeryLow      = "percentage-critical"
	settingKeyPercentageExhausted    = "percentage-action"

	settingKeyTimeToEmptyLow       = "time-to-empty-low"
	settingKeyTimeToEmptyVeryLow   = "time-to-empty-critical"
	settingKeyTimeToEmptyExhausted = "time-to-empty-action"
)

var (
	upowerDBusDest    = "org.freedesktop.UPower"
	upowerDBusObjPath = "/org/freedesktop/UPower"
)

const (
	//defined at http://upower.freedesktop.org/docs/Device.html#Device:Type
	DeviceTypeUnknow    = 0
	DeviceTypeLinePower = 1
	DeviceTypeBattery   = 2
	DeviceTypeUps       = 3
	DeviceTypeMonitor   = 4
	DeviceTypeMouse     = 5
	DeviceTypeKeyboard  = 6
	DeviceTypePda       = 7
	DeviceTypePhone     = 8
)

const (
	//internal used
	batteryPowerLevelUnknown = iota
	batteryPowerLevelSufficient
	batteryPowerLevelAbnormal
	batteryPowerLevelLow
	batteryPowerLevelVeryLow
	batteryPowerLevelExhausted
)

var batteryPowerLevelNameMap = map[uint32]string{
	0: "Unknown",
	1: "Sufficient",
	2: "Abnormal",
	3: "Low",
	4: "VeryLow",
	5: "Exhausted",
}

const (
	batteryPercentageAbnormal = 1.0
	timeToEmptyAbnormal       = 1
)

const (
	//sync with com.deepin.dde.power.schemas
	//
	// 按下电源键和合上笔记本盖时支持的操作
	//
	// 关闭显示器
	ActionBlank int32 = 0
	// 挂起
	ActionSuspend = 1
	// 关机
	ActionShutdown = 2
	// 休眠
	ActionHibernate = 3
	// 询问
	ActionInteractive = 4
	// 无
	ActionNothing = 5
	// 注销
	ActionLogout = 6
)

const (
	powerSupplyDataBackendUPower = 0
	powerSupplyDataBackendPoll   = 1
)

const (
	batteryDisplay    = "Display"
	cmdLowPower       = "/usr/lib/deepin-daemon/dde-lowpower"
	sysPowerSupplyDir = "/sys/class/power_supply"
)
