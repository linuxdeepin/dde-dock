/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package keybinding

import (
	"pkg.deepin.io/lib/dbus"
)

const (
	dbusDest      = "com.deepin.daemon.Keybinding"
	bindDBusPath  = "/com/deepin/daemon/Keybinding"
	bindDBusIFC   = "com.deepin.daemon.Keybinding"
	mediaDBusPath = "/com/deepin/daemon/Keybinding/Mediakey"
	mediaDBusIFC  = "com.deepin.daemon.Keybinding.Mediakey"
)

func (*Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: bindDBusPath,
		Interface:  bindDBusIFC,
	}
}

func (*Mediakey) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: mediaDBusPath,
		Interface:  mediaDBusIFC,
	}
}
