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
	"pkg.deepin.io/lib/dbus"
)

const (
	dbusSender = "com.deepin.daemon.Timedate"
	dbusPath   = "/com/deepin/daemon/Timedate"
	dbusIFC    = "com.deepin.daemon.Timedate"
)

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusSender,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}

func (m *Manager) setPropBool(handler *bool, prop string, value bool) {
	if *handler == value {
		return
	}

	*handler = value
	dbus.NotifyChange(m, prop)
}

func (m *Manager) setPropString(handler *string, prop, value string) {
	if *handler == value {
		return
	}

	*handler = value
	dbus.NotifyChange(m, prop)
}
