/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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
