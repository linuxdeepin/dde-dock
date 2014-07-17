/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
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

package themes

import (
	"pkg.linuxdeepin.com/lib/dbus"
)

func (obj *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		MANAGER_DEST,
		MANAGER_PATH,
		MANAGER_IFC,
	}
}

func (obj *Manager) setPropCurrentTheme(theme string) {
	//if obj.CurrentTheme.GetValue().(string) != theme {
	//obj.CurrentTheme = theme
	//dbus.NotifyChange(obj, "CurrentTheme")
	themeSettings.SetString(GS_KEY_CURRENT_THEME, theme)
	//}
}

func (obj *Manager) setPropThemeList(list []string) {
	if !isStrListEqual(obj.ThemeList, list) {
		obj.ThemeList = list
	}
	dbus.NotifyChange(obj, "ThemeList")
}

func (obj *Manager) setPropGtkThemeList(list []string) {
	if !isStrListEqual(obj.GtkThemeList, list) {
		obj.GtkThemeList = list
	}
	dbus.NotifyChange(obj, "GtkThemeList")
}

func (obj *Manager) setPropIconThemeList(list []string) {
	if !isStrListEqual(obj.IconThemeList, list) {
		obj.IconThemeList = list
	}
	dbus.NotifyChange(obj, "IconThemeList")
}

func (obj *Manager) setPropSoundThemeList(list []string) {
	if !isStrListEqual(obj.SoundThemeList, list) {
		obj.SoundThemeList = list
	}
	dbus.NotifyChange(obj, "SoundThemeList")
}

func (obj *Manager) setPropCursorThemeList(list []string) {
	if !isStrListEqual(obj.CursorThemeList, list) {
		obj.CursorThemeList = list
	}
	dbus.NotifyChange(obj, "CursorThemeList")
}

func (obj *Manager) setPropBackgroundList(list []string) {
	if !isStrListEqual(obj.BackgroundList, list) {
		obj.BackgroundList = list
	}
	dbus.NotifyChange(obj, "BackgroundList")
}

func (obj *Manager) setPropGreeterList(list []string) {
	if !isStrListEqual(obj.GreeterThemeList, list) {
		obj.GreeterThemeList = list
	}
	dbus.NotifyChange(obj, "GreeterThemeList")
}
