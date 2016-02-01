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
