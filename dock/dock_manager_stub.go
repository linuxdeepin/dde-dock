/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dock

import "pkg.deepin.io/lib/dbus"

const (
	dockManagerDBusDest      = "com.deepin.dde.daemon.Dock"
	dockManagerDBusObjPath   = "/com/deepin/dde/daemon/Dock"
	dockManagerDBusInterface = dockManagerDBusDest
)

func (m *DockManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dockManagerDBusDest,
		ObjectPath: dockManagerDBusObjPath,
		Interface:  dockManagerDBusInterface,
	}
}
