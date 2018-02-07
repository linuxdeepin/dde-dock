/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

func (m *Manager) setPropString(handler *string, prop, value string) {
	if *handler == value {
		return
	}

	*handler = value
	dbus.NotifyChange(m, prop)
}

func (m *Manager) setPropFontSize(size float64) {
	cur := m.FontSize.Get()
	if cur > size-0.01 && cur < size+0.01 {
		return
	}

	m.FontSize.Set(size)
	dbus.NotifyChange(m, "FontSize")
}
