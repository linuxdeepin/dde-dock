/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package appearance

import (
	"pkg.deepin.io/lib/dbus"
)

const (
	dbusDest = "com.deepin.daemon.Appearance"
	dbusPath = "/com/deepin/daemon/Appearance"
	dbusIFC  = "com.deepin.daemon.Appearance"
)

func (*Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}

func (m *Manager) setPropTheme(value string) {
	if m.Theme == value {
		return
	}

	m.Theme = value
	dbus.NotifyChange(m, "Theme")
}

func (m *Manager) setPropFontSize(size int32) {
	if m.FontSize == size {
		return
	}

	m.FontSize = size
	dbus.NotifyChange(m, "FontSize")
}
