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
	"sync"

	"pkg.deepin.io/gir/gio-2.0"
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.timedated"
	"github.com/linuxdeepin/go-dbus-factory/org.freedesktop.timedate1"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/gsprop"
	"pkg.deepin.io/lib/dbusutil/proxy"
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
	service       *dbusutil.Service
	systemSigLoop *dbusutil.SignalLoop
	PropsMu       sync.RWMutex
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
	td       *timedate1.Timedate
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
	systemConn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}

	var m = &Manager{
		service: service,
	}

	m.systemSigLoop = dbusutil.NewSignalLoop(systemConn, 10)
	m.td = timedate1.NewTimedate(systemConn)
	m.setter = timedated.NewTimedated(systemConn)

	m.settings = gio.NewSettings(timeDateSchema)
	m.Use24HourFormat.Bind(m.settings, settingsKey24Hour)
	m.DSTOffset.Bind(m.settings, settingsKeyDSTOffset)
	m.UserTimezones.Bind(m.settings, settingsKeyTimezoneList)

	return m, nil
}

func (m *Manager) init() {
	m.PropsMu.Lock()

	canNTP, err := m.td.CanNTP().Get(0)
	if err != nil {
		logger.Warning(err)
	}
	m.setPropCanNTP(canNTP)

	ntp, err := m.td.NTP().Get(0)
	if err != nil {
		logger.Warning(err)
	}
	m.setPropNTP(ntp)

	localRTC, err := m.td.LocalRTC().Get(0)
	if err != nil {
		logger.Warning(err)
	}
	m.setPropLocalRTC(localRTC)

	timezone, err := m.td.Timezone().Get(0)
	if err != nil {
		logger.Warning(err)
	}
	m.setPropTimezone(timezone)

	m.PropsMu.Unlock()

	newList, hasNil := filterNilString(m.UserTimezones.Get())
	if hasNil {
		m.UserTimezones.Set(newList)
	}
	m.AddUserTimezone(m.Timezone)

	m.systemSigLoop.Start()
}

func (m *Manager) destroy() {
	m.settings.Unref()
	m.td.RemoveHandler(proxy.RemoveAllHandlers)
	m.systemSigLoop.Stop()
}

func (*Manager) GetInterfaceName() string {
	return dbusInterface
}
