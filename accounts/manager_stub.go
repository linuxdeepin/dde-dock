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

package accounts

import (
	"pkg.deepin.io/lib/dbus"
)

const (
	dbusSender = "com.deepin.daemon.Accounts"
	dbusPath   = "/com/deepin/daemon/Accounts"
	dbusIFC    = "com.deepin.daemon.Accounts"
)

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusSender,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}

func (m *Manager) setPropUserList(list []string) {
	m.UserList = list
	dbus.NotifyChange(m, "UserList")
}

func (m *Manager) setPropGuestIcon(icon string) {
	if icon == m.GuestIcon {
		return
	}

	m.GuestIcon = icon
	dbus.NotifyChange(m, "GuestIcon")
}

func (m *Manager) setPropAllowGuest(allow bool) {
	if m.AllowGuest == allow {
		return
	}

	m.AllowGuest = allow
	dbus.NotifyChange(m, "AllowGuest")
}
