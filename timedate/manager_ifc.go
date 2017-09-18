/*
 * Copyright (C) 2013 ~ 2017 Deepin Technology Co., Ltd.
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

package timedate

import (
	"time"

	"pkg.deepin.io/dde/daemon/timedate/zoneinfo"
)

func (m *Manager) Reset() error {
	return m.SetNTP(true)
}

// SetDate Set the system clock to the specified.
//
// The time may be specified in the format '2015' '1' '1' '18' '18' '18' '8'.
func (m *Manager) SetDate(year, month, day, hour, min, sec, nsec int32) error {
	loc, err := time.LoadLocation(m.Timezone)
	if err != nil {
		logger.Debugf("Load location '%s' failed: %v", m.Timezone, err)
		return err
	}
	ns := time.Date(int(year), time.Month(month), int(day),
		int(hour), int(min), int(sec), int(nsec), loc).UnixNano()
	return m.SetTime(ns/1000, false)
}

// Set the system clock to the specified.
//
// usec: pass a value of microseconds since 1 Jan 1970 UTC.
//
// relative: if true, the passed usec value will be added to the current system time; if false, the current system time will be set to the passed usec value.
func (m *Manager) SetTime(usec int64, relative bool) error {
	err := m.setter.SetTime(usec, relative)
	if err != nil {
		logger.Debug("SetTime failed:", err)
	}

	return err
}

// To control whether the system clock is synchronized with the network.
//
// useNTP: if true, enable ntp; else disable
func (m *Manager) SetNTP(useNTP bool) error {
	err := m.setter.SetNTP(useNTP)
	if err != nil {
		logger.Debug("SetNTP failed:", err)
	}

	return err
}

// To control whether the RTC is the local time or UTC.
// Standards are divided into: localtime and UTC.
// UTC standard will automatically adjust the daylight saving time.
//
// 实时时间(RTC)是否使用 local 时间标准。时间标准分为 local 和 UTC。
// UTC 时间标准会自动根据夏令时而调整系统时间。
//
// localRTC: whether to use local time.
//
// fixSystem: if true, will use the RTC time to adjust the system clock; if false, the system time is written to the RTC taking the new setting into account.
func (m *Manager) SetLocalRTC(localRTC, fixSystem bool) error {
	err := m.setter.SetLocalRTC(localRTC, fixSystem)
	if err != nil {
		logger.Debug("SetLocalRTC failed:", err)
	}

	return err
}

// Set the system time zone to the specified value.
// timezones you may parse from /usr/share/zoneinfo/zone.tab.
//
// zone: pass a value like "Asia/Shanghai" to set the timezone.
func (m *Manager) SetTimezone(zone string) error {
	if !zoneinfo.IsZoneValid(zone) {
		logger.Debug("Invalid zone:", zone)
		return zoneinfo.ErrZoneInvalid
	}

	err := m.setter.SetTimezone(zone)
	if err != nil {
		logger.Debug("SetTimezone failed:", err)
		return err
	}

	return m.AddUserTimezone(zone)
}

// Add the specified time zone to user time zone list.
func (m *Manager) AddUserTimezone(zone string) error {
	if !zoneinfo.IsZoneValid(zone) {
		logger.Debug("Invalid zone:", zone)
		return zoneinfo.ErrZoneInvalid
	}

	oldList, hasNil := filterNilString(m.UserTimezones.Get())
	newList, added := addItemToList(zone, oldList)
	if added || hasNil {
		m.UserTimezones.Set(newList)
	}
	return nil
}

// Delete the specified time zone from user time zone list.
func (m *Manager) DeleteUserTimezone(zone string) error {
	if !zoneinfo.IsZoneValid(zone) {
		logger.Debug("Invalid zone:", zone)
		return zoneinfo.ErrZoneInvalid
	}

	oldList, hasNil := filterNilString(m.UserTimezones.Get())
	newList, deleted := deleteItemFromList(zone, oldList)
	if deleted || hasNil {
		m.UserTimezones.Set(newList)
	}
	return nil
}

// GetZoneInfo returns the information of the specified time zone.
func (m *Manager) GetZoneInfo(zone string) (zoneinfo.ZoneInfo, error) {
	info, err := zoneinfo.GetZoneInfo(zone)
	if err != nil {
		logger.Debugf("Get zone info for '%s' failed: %v", zone, err)
		return zoneinfo.ZoneInfo{}, err
	}

	return *info, nil
}

// GetZoneList returns all the valid timezones.
func (m *Manager) GetZoneList() []string {
	return zoneinfo.GetAllZones()
}
