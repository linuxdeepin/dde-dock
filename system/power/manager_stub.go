/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
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

func (*Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}

func (m *Manager) setPropOnBattery(val bool) {
	logger.Debug("set OnBattery", val)
	if m.OnBattery != val {
		m.OnBattery = val
		dbus.NotifyChange(m, "OnBattery")
	}
}

func (m *Manager) setPropHasBattery(val bool) {
	if m.HasBattery != val {
		m.HasBattery = val
		dbus.NotifyChange(m, "HasBattery")
	}
}

func (m *Manager) setPropBatteryStatus(val battery.Status) {
	if m.BatteryStatus != val {
		m.BatteryStatus = val
		dbus.NotifyChange(m, "BatteryStatus")
	}
}

func (m *Manager) setPropBatteryPercentage(val float64) {
	if m.BatteryPercentage != val {
		m.BatteryPercentage = val
		dbus.NotifyChange(m, "BatteryPercentage")
	}
}

func (m *Manager) setPropBatteryTimeToEmpty(val uint64) {
	if m.BatteryTimeToEmpty != val {
		m.BatteryTimeToEmpty = val
		dbus.NotifyChange(m, "BatteryTimeToEmpty")
	}
}

func (m *Manager) setPropBatteryTimeToFull(val uint64) {
	if m.BatteryTimeToFull != val {
		m.BatteryTimeToFull = val
		dbus.NotifyChange(m, "BatteryTimeToFull")
	}
}
