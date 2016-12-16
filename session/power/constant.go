/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package power

const (
	// settingKeys
	settingSchema                     = "com.deepin.dde.power"
	settingKeyBatteryScreenBlackDelay = "battery-screen-black-delay"
	settingKeyBatterySleepDelay       = "battery-sleep-delay"

	settingKeyLinePowerScreenBlackDelay = "line-power-screen-black-delay"
	settingKeyLinePowerSleepDelay       = "line-power-sleep-delay"

	settingKeyScreenBlackLock        = "screen-black-lock"
	settingKeySleepLock              = "sleep-lock"
	settingKeyLidClosedSleep         = "lid-closed-sleep"
	settingKeyPowerButtonPressedExec = "power-button-pressed-exec"

	settingKeyFullscreenWorkaroundEnabled     = "fullscreen-workaround-enabled"
	settingKeyMultiScreenPreventLidClosedExec = "multi-screen-prevent-lid-closed-exec"
	settingKeyUsePercentageForPolicy          = "use-percentage-for-policy"

	settingKeyPowerModuleInitialized = "power-module-initialized"
	settingKeyLowPercentage          = "percentage-low"
	settingKeyCriticalPercentage     = "percentage-critical"
	settingKeyActionPercentage       = "percentage-action"

	settingKeyLowTime      = "time-to-empty-low"
	settingKeyCriticalTime = "time-to-empty-critical"
	settingKeyActionTime   = "time-to-empty-action"

	// dbus info
	dbusDisplayDest = "com.deepin.daemon.Display"
	dbusDisplayPath = "/com/deepin/daemon/Display"

	// cmd
	cmdLowPower = "/usr/lib/deepin-daemon/dde-lowpower"

	batteryDisplay = "Display"
)
