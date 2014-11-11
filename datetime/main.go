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
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/log"
)

var _date *DateTime

func Start() {
	var logger = log.NewLogger(dbusSender)

	logger.BeginTracing()

	_date = NewDateTime(logger)
	if _date == nil {
		logger.Error("Create DateTime Failed")
		return
	}
	err := dbus.InstallOnSession(_date)
	if err != nil {
		logger.Error("Install DBus For DateTime Failed")
		return
	}
}

func Stop() {
	if _date == nil {
		return
	}

	_date.Destroy()
	dbus.UnInstallObject(_date)
	_date = nil
}
