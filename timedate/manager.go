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
	"os/user"
	"sync"

	"github.com/godbus/dbus"
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.accounts"
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.timedated"
	"github.com/linuxdeepin/go-dbus-factory/org.freedesktop.timedate1"
	ddbus "pkg.deepin.io/dde/daemon/dbus"
	"pkg.deepin.io/dde/daemon/session/common"
	"pkg.deepin.io/gir/gio-2.0"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/gsprop"
	"pkg.deepin.io/lib/dbusutil/proxy"
)

const (
	timeDateSchema             = "com.deepin.dde.datetime"
	settingsKey24Hour          = "is-24hour"
	settingsKeyTimezoneList    = "user-timezone-list"
	settingsKeyDSTOffset       = "dst-offset"
	settingsKeyWeekdayFormat   = "weekday-format"
	settingsKeyShortDateFormat = "short-date-format"
	settingsKeyLongDateFormat  = "long-date-format"
	settingsKeyShortTimeFormat = "short-time-format"
	settingsKeyLongTimeFormat  = "long-time-format"
	settingsKeyWeekBegins      = "week-begins"

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
	Timezone  string
	NTPServer string

	// dbusutil-gen: ignore-below
	// Use 24 hour format to display time
	Use24HourFormat gsprop.Bool `prop:"access:rw"`
	// DST offset
	DSTOffset gsprop.Int `prop:"access:rw"`
	// User added timezone list
	UserTimezones gsprop.Strv

	// weekday shows format
	WeekdayFormat gsprop.Int `prop:"access:rw"`

	// short date shows format
	ShortDateFormat gsprop.Int `prop:"access:rw"`

	// long date shows format
	LongDateFormat gsprop.Int `prop:"access:rw"`

	// short time shows format
	ShortTimeFormat gsprop.Int `prop:"access:rw"`

	// long time shows format
	LongTimeFormat gsprop.Int `prop:"access:rw"`

	WeekBegins gsprop.Int `prop:"access:rw"`

	settings *gio.Settings
	td       *timedate1.Timedate
	setter   *timedated.Timedated
	userObj  *accounts.User
	//nolint
	methods *struct {
		SetDate             func() `in:"year,month,day,hour,min,sec,nsec"`
		SetTime             func() `in:"usec,relative"`
		SetNTP              func() `in:"useNTP"`
		SetNTPServer        func() `in:"server"`
		GetSampleNTPServers func() `out:"servers"`
		SetLocalRTC         func() `in:"localeRTC,fixSystem"`
		SetTimezone         func() `in:"zone"`
		AddUserTimezone     func() `in:"zone"`
		DeleteUserTimezone  func() `in:"zone"`
		GetZoneInfo         func() `in:"zone" out:"zone_info"`
		GetZoneList         func() `out:"zone_list"`
	}
}

// Create Manager, if create freedesktop timedate1 failed return error
func NewManager(service *dbusutil.Service) (*Manager, error) {
	sysBus, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}

	var m = &Manager{
		service: service,
	}

	m.systemSigLoop = dbusutil.NewSignalLoop(sysBus, 10)
	m.td = timedate1.NewTimedate(sysBus)
	m.setter = timedated.NewTimedated(sysBus)

	m.settings = gio.NewSettings(timeDateSchema)
	m.Use24HourFormat.Bind(m.settings, settingsKey24Hour)
	m.DSTOffset.Bind(m.settings, settingsKeyDSTOffset)
	m.UserTimezones.Bind(m.settings, settingsKeyTimezoneList)

	m.WeekdayFormat.Bind(m.settings, settingsKeyWeekdayFormat)
	m.ShortDateFormat.Bind(m.settings, settingsKeyShortDateFormat)
	m.LongDateFormat.Bind(m.settings, settingsKeyLongDateFormat)
	m.ShortTimeFormat.Bind(m.settings, settingsKeyShortTimeFormat)
	m.LongTimeFormat.Bind(m.settings, settingsKeyLongTimeFormat)
	m.WeekBegins.Bind(m.settings, settingsKeyWeekBegins)

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
	err = m.AddUserTimezone(m.Timezone)
	if err != nil {
		logger.Warning("AddUserTimezone error:", err)
	}

	err = common.ActivateSysDaemonService(m.setter.ServiceName_())
	if err != nil {
		logger.Warning(err)
	} else {
		ntpServer, err := m.setter.NTPServer().Get(0)
		if err != nil {
			logger.Warning(err)
		} else {
			m.NTPServer = ntpServer
		}
	}

	sysBus := m.systemSigLoop.Conn()
	m.initUserObj(sysBus)
	m.handleGSettingsChanged()
	m.systemSigLoop.Start()
	m.listenPropChanged()
}

func (m *Manager) destroy() {
	m.settings.Unref()
	m.td.RemoveHandler(proxy.RemoveAllHandlers)
	m.systemSigLoop.Stop()
}

func (*Manager) GetInterfaceName() string {
	return dbusInterface
}

func (m *Manager) initUserObj(systemConn *dbus.Conn) {
	cur, err := user.Current()
	if err != nil {
		logger.Warning("failed to get current user:", err)
		return
	}

	err = common.ActivateSysDaemonService("com.deepin.daemon.Accounts")
	if err != nil {
		logger.Warning(err)
	}

	m.userObj, err = ddbus.NewUserByUid(systemConn, cur.Uid)
	if err != nil {
		logger.Warning("failed to new user object:", err)
		return
	}

	// sync use 24 hour format
	use24hourFormat := m.settings.GetBoolean(settingsKey24Hour)
	err = m.userObj.SetUse24HourFormat(0, use24hourFormat)
	if err != nil {
		logger.Warning(err)
	}

	weekdayFormat := m.settings.GetInt(settingsKeyWeekdayFormat)
	err = m.userObj.SetWeekdayFormat(0, weekdayFormat)
	if err != nil {
		logger.Warning(err)
	}

	shortDateFormat := m.settings.GetInt(settingsKeyShortDateFormat)
	err = m.userObj.SetShortDateFormat(0, shortDateFormat)
	if err != nil {
		logger.Warning(err)
	}

	longDateFormat := m.settings.GetInt(settingsKeyLongDateFormat)
	err = m.userObj.SetLongDateFormat(0, longDateFormat)
	if err != nil {
		logger.Warning(err)
	}

	shortTimeFormat := m.settings.GetInt(settingsKeyShortTimeFormat)
	err = m.userObj.SetShortTimeFormat(0, shortTimeFormat)
	if err != nil {
		logger.Warning(err)
	}

	longTimeFormat := m.settings.GetInt(settingsKeyLongTimeFormat)
	err = m.userObj.SetLongDateFormat(0, longTimeFormat)
	if err != nil {
		logger.Warning(err)
	}

	weekBegins := m.settings.GetInt(settingsKeyWeekBegins)
	err = m.userObj.SetWeekBegins(0, weekBegins)
	if err != nil {
		logger.Warning(err)
	}
}
