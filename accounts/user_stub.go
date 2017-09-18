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
	userDBusPath = "/com/deepin/daemon/Accounts/User"
	userDBusIFC  = "com.deepin.daemon.Accounts.User"
)

func (u *User) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusSender,
		ObjectPath: userDBusPath + u.Uid,
		Interface:  userDBusIFC,
	}
}

func (u *User) setPropBool(handler *bool, prop string, value bool) {
	if *handler == value {
		return
	}

	*handler = value
	dbus.NotifyChange(u, prop)
}

func (u *User) setPropInt32(handler *int32, prop string, value int32) {
	if *handler == value {
		return
	}

	*handler = value
	dbus.NotifyChange(u, prop)
}

func (u *User) setPropString(handler *string, prop string, value string) {
	if *handler == value {
		return
	}

	*handler = value
	dbus.NotifyChange(u, prop)
}

func (u *User) setPropStrv(handler *[]string, prop string, value []string) {
	*handler = value
	dbus.NotifyChange(u, prop)
}

func (u *User) setPropIconFile(value string) {
	u.setPropString(&u.IconFile, "IconFile", value)
}
