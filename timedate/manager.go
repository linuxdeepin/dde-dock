/*
 * Copyright (C) 2013 ~ 2018 Deepin Technology Co., Ltd.
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
	"dbus/com/deepin/daemon/timedated"
	"dbus/org/freedesktop/timedate1"
	"gir/gio-2.0"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/dbus/property"
)

const (
	timedateSchema          = "com.deepin.dde.datetime"
	settingsKey24Hour       = "is-24hour"
	settingsKeyTimezoneList = "user-timezone-list"
	settingsKeyDSTOffset    = "dst-offset"
)

// Manage time settings
type Manager struct {
	// Whether can use NTP service
	CanNTP bool
	// Whether enable NTP service
	NTP bool
	// Whether set RTC to Local standard
	LocalRTC bool

	// Current timezone
	Timezone string

	// Use 24 hour format to display time
	Use24HourFormat *property.GSettingsBoolProperty `access:"readwrite"`
	// DST offset
	DSTOffset *property.GSettingsIntProperty `access:"readwrite"`
	// User added timezone list
	UserTimezones *property.GSettingsStrvProperty

	settings *gio.Settings
	td1      *timedate1.Timedate1
	setter   *timedated.Timedated
}

// Create Manager, if create freedesktop timedate1 failed return error
func NewManager() (*Manager, error) {
	var m = &Manager{}

	var err error
	m.td1, err = timedate1.NewTimedate1("org.freedesktop.timedate1",
		"/org/freedesktop/timedate1")
	if err != nil {
		return nil, err
	}
	m.setter, err = timedated.NewTimedated("com.deepin.daemon.Timedated",
		"/com/deepin/daemon/Timedated")
	if err != nil {
		timedate1.DestroyTimedate1(m.td1)
		return nil, err
	}

	return m, nil
}

func (m *Manager) init() {
	m.setPropBool(&m.CanNTP, "CanNTP", m.td1.CanNTP.Get())
	m.setPropBool(&m.NTP, "NTP", m.td1.NTP.Get())
	m.setPropBool(&m.LocalRTC, "LocalRTC", m.td1.LocalRTC.Get())
	m.setPropString(&m.Timezone, "Timezone", m.td1.Timezone.Get())

	m.settings = gio.NewSettings(timedateSchema)
	m.Use24HourFormat = property.NewGSettingsBoolProperty(
		m, "Use24HourFormat",
		m.settings, settingsKey24Hour)
	m.DSTOffset = property.NewGSettingsIntProperty(
		m, "DSTOffset",
		m.settings, settingsKeyDSTOffset)
	m.UserTimezones = property.NewGSettingsStrvProperty(
		m, "UserTimezones",
		m.settings, settingsKeyTimezoneList)

	newList, hasNil := filterNilString(m.UserTimezones.Get())
	if hasNil {
		m.UserTimezones.Set(newList)
	}
	m.AddUserTimezone(m.Timezone)
}

func (m *Manager) destroy() {
	if m.settings != nil {
		m.settings.Unref()
		m.settings = nil
	}

	if m.td1 != nil {
		timedate1.DestroyTimedate1(m.td1)
		m.td1 = nil
	}

	if m.setter != nil {
		m.setter = nil
	}

	dbus.UnInstallObject(m)
}
