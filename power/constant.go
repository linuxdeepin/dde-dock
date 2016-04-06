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
	//defined at http://upower.freedesktop.org/docs/Device.html#Device:State
	BatteryStateUnknown          = 0
	BatteryStateCharging         = 1
	BatteryStateDischarging      = 2
	BatteryStateEmpty            = 3
	BatteryStateFullyCharged     = 4
	BatteryStatePendingCharge    = 5
	BatteryStatePendingDischarge = 6
)

var batteryStateMap = map[string]uint32{
	"Unknown":          BatteryStateUnknown,
	"Charging":         BatteryStateCharging,
	"Discharging":      BatteryStateDischarging,
	"Empty":            BatteryStateEmpty,
	"FullCharged":      BatteryStateFullyCharged,
	"PendingCharge":    BatteryStatePendingCharge,
	"PendingDischarge": BatteryStatePendingDischarge,
}

const (
	//internal used
	batteryPowerLevelSufficient = iota
	batteryPowerLevelAbnormal
	batteryPowerLevelLow
	batteryPowerLevelVeryLow
	batteryPowerLevelExhausted
)

var batteryPowerLevelNameMap = map[uint32]string{
	0: "Sufficient",
	1: "Abnormal",
	2: "Low",
	3: "VeryLow",
	4: "Exhausted",
}

const (
	batteryPercentageAbnormal  = 1.0
	batteryPercentageLow       = 20.0
	batteryPercentageVeryLow   = 10.0
	batteryPercentageExhausted = 5.0

	timeToEmptyLow       = 1200
	timeToEmptyVeryLow   = 600
	timeToEmptyExhausted = 300
	timeToEmptyAbnormal  = 1
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
	cmdLowPower = "/usr/lib/deepin-daemon/dde-lowpower"
)
