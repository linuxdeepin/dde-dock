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
	"dlib/gio-2.0"
	"strings"
)

func InitSystemBind(m *BindManager) {
	m.SystemList = GetSystemKeyInfo()
}

func GetSystemKeyInfo() []*ShortcutInfo {
	systemInfoList := []*ShortcutInfo{}
	for i, n := range SystemIdNameMap {
		if desc, ok := SystemNameDescMap[n]; ok {
			shortcut := GetSystemValue(n, false)
			tmp := NewShortcutInfo(i, desc,
				FormatShortcut(shortcut))
			systemInfoList = append(systemInfoList, tmp)
		}
	}

	return systemInfoList
}

func InitMediaBind(m *BindManager) {
	m.MediaList = GetMediaKeyInfo()
}

func GetMediaKeyInfo() []*ShortcutInfo {
	mediaInfoList := []*ShortcutInfo{}
	for i, n := range MediaIdNameMap {
		if desc, ok := MediaNameDescMap[n]; ok {
			shortcut := GetSystemValue(n, false)
			tmp := NewShortcutInfo(i, desc, shortcut)
			mediaInfoList = append(mediaInfoList, tmp)
		}
	}

	return mediaInfoList
}

func InitWindowBind(m *BindManager) {
	m.WindowList = GetWindowKeyInfo()
}

func GetWindowKeyInfo() []*ShortcutInfo {
	windowInfoList := []*ShortcutInfo{}
	for i, n := range WindowIdNameMap {
		if desc, ok := WindowNameDescMap[n]; ok {
			shortcut := GetSystemValue(n, false)
			tmp := NewShortcutInfo(i, desc,
				FormatShortcut(shortcut))
			windowInfoList = append(windowInfoList, tmp)
		}
	}

	return windowInfoList
}

func InitWorkSpaceBind(m *BindManager) {
	m.WorkSpaceList = GetWorkSpaceKeyInfo()
}

func GetWorkSpaceKeyInfo() []*ShortcutInfo {
	workSpaceInfoList := []*ShortcutInfo{}
	for i, n := range WorkSpaceIdNameMap {
		if desc, ok := WorkSpaceNameDescMap[n]; ok {
			shortcut := GetSystemValue(n, false)
			tmp := NewShortcutInfo(i, desc,
				FormatShortcut(shortcut))
			workSpaceInfoList = append(workSpaceInfoList, tmp)
		}
	}

	return workSpaceInfoList
}

func GetSystemValue(name string, action bool) string {
	values := systemGSettings.GetStrv(name)
	if len(values) <= 0 {
		return ""
	}

	if len(values) == 2 {
		if action {
			return values[1]
		}
	}

	return strings.ToLower(values[0])
}

func ListenSystem(m *BindManager) {
	systemGSettings.Connect("changed", func(s *gio.Settings, key string) {
		id := GetIdByName(key)
		shortcut := GetSystemValue(key, false)
		if !UpdateSystemList(m, id, shortcut) {
			return
		}

		UpdateCompizValue(id, key, shortcut)
	})
}

func GetIdByName(name string) int32 {
	if len(name) <= 0 {
		return -1
	}

	for i, n := range IdNameMap {
		if name == n {
			return i
		}
	}

	return -1
}

func UpdateSystemList(m *BindManager, id int32, shortcut string) bool {
	for _, info := range m.SystemList {
		if info.Id == id {
			info.Shortcut = shortcut
			dbus.NotifyChange(m, "SystemList")
			GrabKeyPairs(SystemPrevPairs, false)
			GrabKeyPairs(GetSystemPairs(), true)
			return true
		}
	}

	for _, info := range m.WindowList {
		if info.Id == id {
			info.Shortcut = shortcut
			dbus.NotifyChange(m, "WindowList")
			return true
		}
	}

	for _, info := range m.WorkSpaceList {
		if info.Id == id {
			info.Shortcut = shortcut
			dbus.NotifyChange(m, "WorkSpaceList")
			return true
		}
	}

	return false
}

/* Update Shortcut */
func UpdateSystemShortcut(key, value string) {
	values := systemGSettings.GetStrv(key)
	values[0] = value

	systemGSettings.SetStrv(key, values)
	gio.SettingsSync()
}

func UpdateCompizValue(id int32, key, shortcut string) {
	if id >= 600 && id < 800 {
		values := wmGSettings.GetStrv(key)
		if values[0] == shortcut {
			return
		}
		wmGSettings.SetStrv(key, []string{shortcut})
	} else if id >= 800 && id < 900 {
		value := shiftGSettings.GetString(key)
		if value == shortcut {
			return
		}
		shiftGSettings.SetString(key, shortcut)
	} else if id >= 900 && id < 1000 {
		value := putGSettings.GetString(key)
		if value == shortcut {
			return
		}
		putGSettings.SetString(key, shortcut)
	}
	gio.SettingsSync()
}

func ListenCompiz(m *BindManager) {
	wmGSettings.Connect("changed", func(s *gio.Settings, key string) {
		value := GetSystemValue(key, false)
		tmps := wmGSettings.GetStrv(key)
		if value == tmps[0] {
			return
		}
		UpdateSystemShortcut(key, tmps[0])
	})

	shiftGSettings.Connect("changed", func(s *gio.Settings, key string) {
		tmp := shiftGSettings.GetString(key)
		value := GetSystemValue(key, false)
		if tmp == value {
			return
		}
		UpdateSystemShortcut(key, tmp)
	})

	putGSettings.Connect("changed", func(s *gio.Settings, key string) {
		tmp := putGSettings.GetString(key)
		value := GetSystemValue(key, false)
		if tmp == value {
			return
		}
		UpdateSystemShortcut(key, tmp)
	})
}
