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

package appearance

import "pkg.linuxdeepin.com/lib/dbus"

const (
	managerDBusPath = "/com/deepin/daemon/ThemeManager"
	managerDBusIFC  = "com.deepin.daemon.ThemeManager"
)

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusSender,
		ObjectPath: managerDBusPath,
		Interface:  managerDBusIFC,
	}
}

func (m *Manager) setPropThemeList(list []string) {
	if !m.isStrListEqual(m.ThemeList, list) {
		m.ThemeList = list
		dbus.NotifyChange(m, "ThemeList")
	}
}

func (m *Manager) setPropGtkThemeList(list []string) {
	if !m.isStrListEqual(m.GtkThemeList, list) {
		m.GtkThemeList = sortNameByDeepin(list)
		dbus.NotifyChange(m, "GtkThemeList")
	}
}

func (m *Manager) setPropIconThemeList(list []string) {
	if !m.isStrListEqual(m.IconThemeList, list) {
		m.IconThemeList = sortNameByDeepin(list)
		dbus.NotifyChange(m, "IconThemeList")
	}
}

func (m *Manager) setPropCursorThemeList(list []string) {
	if !m.isStrListEqual(m.CursorThemeList, list) {
		m.CursorThemeList = sortNameByDeepin(list)
		dbus.NotifyChange(m, "CursorThemeList")
	}
}

func (m *Manager) setPropSoundThemeList(list []string) {
	if !m.isStrListEqual(m.SoundThemeList, list) {
		m.SoundThemeList = sortNameByDeepin(list)
		dbus.NotifyChange(m, "SoundThemeList")
	}
}

func (m *Manager) setPropGreeterThemeList(list []string) {
	if !m.isStrListEqual(m.GreeterThemeList, list) {
		m.GreeterThemeList = sortNameByDeepin(list)
		dbus.NotifyChange(m, "GreeterThemeList")
	}
}

func (m *Manager) setPropBackgroundList(list []string) {
	if !m.isStrListEqual(m.BackgroundList, list) {
		m.BackgroundList = list
		dbus.NotifyChange(m, "BackgroundList")
	}
}

func (m *Manager) setPropFontNameList(list []string) {
	if !m.isStrListEqual(m.FontNameList, list) {
		m.FontNameList = list
		dbus.NotifyChange(m, "FontNameList")
	}
}

func (m *Manager) setPropFontMonoList(list []string) {
	if !m.isStrListEqual(m.FontMonoList, list) {
		m.FontMonoList = list
		dbus.NotifyChange(m, "FontMonoList")
	}
}

func (m *Manager) isStrListEqual(list1, list2 []string) bool {
	len1 := len(list1)
	len2 := len(list2)

	if len1 != len2 {
		return false
	}

	for i := 0; i < len1; i++ {
		if list1[i] != list2[i] {
			return false
		}
	}

	return true
}
