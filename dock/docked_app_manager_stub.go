/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dock

import (
	"pkg.deepin.io/lib/dbus"
)

func (m *DockedAppManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.daemon.Dock",
		ObjectPath: "/dde/dock/DockedAppManager",
		Interface:  "dde.dock.DockedAppManager",
	}
}

func (m *DockedAppManager) destroy() {
	if m.core != nil {
		m.core.Unref()
	}
	dbus.UnInstallObject(m)
}
