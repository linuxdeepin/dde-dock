/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
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

package datetime

import (
	"io/ioutil"
	"pkg.linuxdeepin.com/dde-daemon/datetime/ntp"
	"pkg.linuxdeepin.com/dde-daemon/datetime/timezone"
	. "pkg.linuxdeepin.com/dde-daemon/datetime/utils"
	"pkg.linuxdeepin.com/lib/dbus/property"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"pkg.linuxdeepin.com/lib/log"
	"strings"
)

const (
	dbusSender = "com.deepin.daemon.DateAndTime"
	dbusPath   = "/com/deepin/daemon/DateAndTime"
	dbusIFC    = "com.deepin.daemon.DateAndTime"

	gsKeyAutoSetTime  = "is-auto-set"
	gsKey24Hour       = "is-24hour"
	gsKeyUTCOffset    = "utc-offset"
	gsKeyTimezoneList = "user-timezone-list"
	gsKeyDSTOffset    = "dst-offset"

	defaultTimezone     = "UTC"
	defaultTimezoneFile = "/etc/timezone"
	defaultUTCOffset    = "+00:00"
)

type DateTime struct {
	NTPEnabled       *property.GSettingsBoolProperty `access:"readwrite"`
	Use24HourDisplay *property.GSettingsBoolProperty `access:"readwrite"`
	DSTOffset        *property.GSettingsIntProperty  `access:"readwrite"`

	UserTimezones *property.GSettingsStrvProperty

	CurrentTimezone string

	settings *gio.Settings
	logger   *log.Logger
}

func NewDateTime(l *log.Logger) *DateTime {
	date := &DateTime{}

	err := InitSetDateTime()
	if err != nil {
		return nil
	}

	err = ntp.InitNtpModule()
	if err != nil {
		return nil
	}

	date.logger = l
	date.settings = gio.NewSettings("com.deepin.dde.datetime")
	date.NTPEnabled = property.NewGSettingsBoolProperty(
		date, "NTPEnabled",
		date.settings, gsKeyAutoSetTime)
	date.Use24HourDisplay = property.NewGSettingsBoolProperty(
		date, "Use24HourDisplay",
		date.settings, gsKey24Hour)
	date.UserTimezones = property.NewGSettingsStrvProperty(
		date, "UserTimezones",
		date.settings, gsKeyTimezoneList)
	date.DSTOffset = property.NewGSettingsIntProperty(
		date, "DSTOffset",
		date.settings, gsKeyDSTOffset)

	date.setPropString(&date.CurrentTimezone,
		"CurrentTimezone", getDefaultTimezone(defaultTimezoneFile))

	date.AddUserTimezone(date.CurrentTimezone)
	date.enableNTP(date.NTPEnabled.Get())

	return date
}

func (date *DateTime) enableNTP(value bool) {
	ntp.Enabled(value, date.CurrentTimezone)
}

func getDefaultTimezone(config string) string {
	zone, err := getTimezoneFromFile(config)
	if err != nil || !timezone.IsZoneValid(zone) {
		return defaultTimezone
	}

	return zone
}

func getTimezoneFromFile(filename string) (string, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(contents), "\n")
	return lines[0], nil
}
