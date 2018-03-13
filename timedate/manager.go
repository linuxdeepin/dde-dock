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
	"sync"

	"gir/gio-2.0"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/gsprop"
)

const (
	timeDateSchema          = "com.deepin.dde.datetime"
	settingsKey24Hour       = "is-24hour"
	settingsKeyTimezoneList = "user-timezone-list"
	settingsKeyDSTOffset    = "dst-offset"

	dbusServiceName = "com.deepin.daemon.Timedate"
	dbusPath        = "/com/deepin/daemon/Timedate"
	dbusInterface   = dbusServiceName
)

//go:generate dbusutil-gen -type Manager manager.go
// Manage time settings
type Manager struct {
	service *dbusutil.Service
	PropsMu sync.RWMutex
	// Whether can use NTP service
	CanNTP bool
	// Whether enable NTP service
	NTP bool
	// Whether set RTC to Local standard
	LocalRTC bool

	// Current timezone
	Timezone string

	// dbusutil-gen: ignore-below
	// Use 24 hour format to display time
	Use24HourFormat gsprop.Bool `prop:"access:rw"`
	// DST offset
	DSTOffset gsprop.Int `prop:"access:rw"`
	// User added timezone list
	UserTimezones gsprop.Strv

	settings *gio.Settings
	td1      *timedate1.Timedate1
	setter   *timedated.Timedated

	methods *struct {
		SetDate            func() `in:"year,month,day,hour,min,sec,nsec"`
		SetTime            func() `in:"usec,relative"`
		SetNTP             func() `in:"useNTP"`
		SetLocalRTC        func() `in:"localeRTC,fixSystem"`
		SetTimezone        func() `in:"zone"`
		AddUserTimezone    func() `in:"zone"`
		DeleteUserTimezone func() `in:"zone"`
		GetZoneInfo        func() `in:"zone" out:"zone_info"`
		GetZoneList        func() `out:"zone_list"`
	}
}

// Create Manager, if create freedesktop timedate1 failed return error
func NewManager(service *dbusutil.Service) (*Manager, error) {
	var m = &Manager{
		service: service,
	}

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

	m.settings = gio.NewSettings(timeDateSchema)
	m.Use24HourFormat.Bind(m.settings, settingsKey24Hour)
	m.DSTOffset.Bind(m.settings, settingsKeyDSTOffset)
	m.UserTimezones.Bind(m.settings, settingsKeyTimezoneList)

	return m, nil
}

func (m *Manager) init() {
	m.PropsMu.Lock()
	m.setPropCanNTP(m.td1.CanNTP.Get())
	m.setPropNTP(m.td1.NTP.Get())
	m.setPropLocalRTC(m.td1.LocalRTC.Get())
	m.setPropTimezone(m.td1.Timezone.Get())
	m.PropsMu.Unlock()

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

	m.service.StopExport(m)
}

func (*Manager) GetInterfaceName() string {
	return dbusInterface
}
