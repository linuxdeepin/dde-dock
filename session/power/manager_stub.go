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

import (
	"pkg.deepin.io/dde/api/powersupply/battery"
	"pkg.deepin.io/lib/dbus"
)

const (
	dbusDest = "com.deepin.daemon.Power"
	dbusPath = "/com/deepin/daemon/Power"
	dbusIFC  = dbusDest
)

func (*Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}

func (m *Manager) Reset() {
	logger.Debug("Reset settings")

	var settingKeys = []string{
		settingKeyLinePowerScreenBlackDelay,
		settingKeyLinePowerSleepDelay,
		settingKeyBatteryScreenBlackDelay,
		settingKeyBatterySleepDelay,
		settingKeyScreenBlackLock,
		settingKeySleepLock,
		settingKeyLidClosedSleep,
		settingKeyPowerButtonPressedExec,
	}
	for _, key := range settingKeys {
		logger.Debug("reset setting", key)
		m.settings.Reset(key)
	}
}

func (m *Manager) setPropOnBattery(val bool) {
	if m.OnBattery != val {
		m.OnBattery = val
		dbus.NotifyChange(m, "OnBattery")
	}
}

func (m *Manager) setPropBatteryIsPresent(val bool) {
	old, exist := m.BatteryIsPresent[batteryDisplay]
	if old != val || !exist {
		m.BatteryIsPresent[batteryDisplay] = val
		dbus.NotifyChange(m, "BatteryIsPresent")
	}
}

func (m *Manager) setPropBatteryPercentage(val float64) {
	logger.Debugf("set batteryDisplay percentage %.1f%%", val)
	old, exist := m.BatteryPercentage[batteryDisplay]
	if old != val || !exist {
		m.BatteryPercentage[batteryDisplay] = val
		dbus.NotifyChange(m, "BatteryPercentage")
	}
}

func (m *Manager) setPropBatteryState(val uint32) {
	logger.Debug("set BatteryDisplay status", battery.Status(val), val)
	old, exist := m.BatteryState[batteryDisplay]
	if old != val || !exist {
		m.BatteryState[batteryDisplay] = val
		dbus.NotifyChange(m, "BatteryState")
	}
}

func (m *Manager) setPropWarnLevel(val WarnLevel) {
	logger.Debug("set WarnLevel", val, int(val))
	if m.WarnLevel != val {
		m.WarnLevel = val
		m.handleWarnLevelChanged()
		dbus.NotifyChange(m, "WarnLevel")
	}
}
