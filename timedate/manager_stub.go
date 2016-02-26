/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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
