/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package main

import (
        "dlib/dbus"
)

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		MANAGER_DEST,
		MANAGER_PATH,
		MANAGER_IFC,
	}
}

func (m *Manager) setPropThemeInfo (name string) {
        switch name {
                case "AvailableFontTheme": {
                        m.AvailableBackground = getBackgroundFiles()
                        dbus.NotifyChange(m, name)
                }
                break
                case "AvailableBackground": {
                        m.AvailableFontTheme = getFontThemes()
                        dbus.NotifyChange(m, name)
                }
                break
                case "AvailableIconTheme": {
                        for _, v := range systemThemes {
		icon := ThemeType{Name: v.IconTheme, Type: "system"}
		m.AvailableIconTheme = append(m.AvailableIconTheme, icon)
                        }
                        dbus.NotifyChange(m, name)
                }
                break
                case "AvailableGtkTheme": {
                        for _, v := range systemThemes {
		gtk := ThemeType{Name: v.GtkTheme, Type: "system"}
		m.AvailableGtkTheme = append(m.AvailableGtkTheme, gtk)
                        }
                        dbus.NotifyChange(m, name)
                }
                break
                case "AvailableCursorTheme": {
                        for _, v := range systemThemes {
		cursor := ThemeType{Name: v.CursorTheme, Type: "system"}
		m.AvailableCursorTheme = append(m.AvailableCursorTheme, cursor)
                        }
                        dbus.NotifyChange(m, name)
                }
                break
                case "AvailableWindowTheme": {
                        for _, v := range systemThemes {
		window := ThemeType{Name: v.WindowTheme, Type: "system"}
		m.AvailableWindowTheme = append(m.AvailableWindowTheme, window)
                        }
                        dbus.NotifyChange(m, name)
                }
                break
        default:
                break
        }
}
