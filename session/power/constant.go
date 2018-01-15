/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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

const (
	// settingKeys
	settingSchema                     = "com.deepin.dde.power"
	settingKeyBatteryScreenBlackDelay = "battery-screen-black-delay"
	settingKeyBatterySleepDelay       = "battery-sleep-delay"

	settingKeyLinePowerScreenBlackDelay = "line-power-screen-black-delay"
	settingKeyLinePowerSleepDelay       = "line-power-sleep-delay"

	settingKeyAdjustBrightnessEnabled = "adjust-brightness-enabled"
	settingKeyScreenBlackLock         = "screen-black-lock"
	settingKeySleepLock               = "sleep-lock"
	settingKeyLidClosedSleep          = "lid-closed-sleep"
	settingKeyPowerButtonPressedExec  = "power-button-pressed-exec"

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
