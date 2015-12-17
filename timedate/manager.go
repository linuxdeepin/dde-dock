/**
 * Copyright (c) 2011 ~ 2015 Deepin, Inc.
 *               2013 ~ 2015 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package timedate

import (
	"dbus/org/freedesktop/timedate1"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/dbus/property"
	"gir/gio-2.0"
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

	return m, nil
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

	dbus.UnInstallObject(m)
}
