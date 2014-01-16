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
	"fmt"
	"strconv"
	"strings"
)

func (m *BindManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_BINDING_DEST,
		_BINDING_PATH,
		_BINDING_IFC,
	}
}

func (m *BindManager) setPropList(listName string) {
	switch listName {
	case "SystemList":
		m.SystemList = getSystemKeyInfo()
		dbus.NotifyChange(m, listName)
	case "MediaList":
		m.MediaList = getMediaKeyInfo()
		dbus.NotifyChange(m, listName)
	case "WindowList":
		m.WindowList = getWindowKeyInfo()
		dbus.NotifyChange(m, listName)
	case "WorkSpaceList":
		m.WorkSpaceList = getWorkSpaceKeyInfo()
		dbus.NotifyChange(m, listName)
	case "CustomList":
		m.CustomList = getCustomKeyInfo()
		dbus.NotifyChange(m, listName)
	case "ConflictInvalid":
		m.ConflictInvalid = getConflictList(false)
		dbus.NotifyChange(m, listName)
	case "ConflictValid":
		m.ConflictValid = getConflictList(true)
		dbus.NotifyChange(m, listName)
	}
}

func (m *BindManager) listenSystem() {
	systemGSettings.Connect("changed", func(s *gio.Settings, key string) {
		id := getIdByName(key)
		shortcut := getSystemValue(key, false)
		if !m.updateSystemList(id, shortcut) {
			return
		}

		updateCompizValue(id, key, shortcut)
	})
}

func (m *BindManager) updateSystemList(id int32, shortcut string) bool {
	for _, info := range m.SystemList {
		if info.Id == id {
			//info.Shortcut = shortcut
			//dbus.NotifyChange(m, "SystemList")
			m.setPropList("SystemList")
			grabKeyPairs(SystemPrevPairs, false)
			grabKeyPairs(getSystemPairs(), true)
			return true
		}
	}

	for _, info := range m.WindowList {
		if info.Id == id {
			//info.Shortcut = shortcut
			//dbus.NotifyChange(m, "WindowList")
			m.setPropList("WindowList")
			return true
		}
	}

	for _, info := range m.WorkSpaceList {
		if info.Id == id {
			//info.Shortcut = shortcut
			//dbus.NotifyChange(m, "WorkSpaceList")
			m.setPropList("WorkSpaceList")
			return true
		}
	}

	return false
}

func (m *BindManager) listenCompiz() {
	wmGSettings.Connect("changed", func(s *gio.Settings, key string) {
		if !keyIsExist(key) {
			return
		}
		value := getSystemValue(key, false)
		tmps := wmGSettings.GetStrv(key)
		if len(tmps) <= 0 {
			return
		}
		if value == tmps[0] {
			return
		}
		UpdateSystemShortcut(key, tmps[0])
	})

	shiftGSettings.Connect("changed", func(s *gio.Settings, key string) {
		if !keyIsExist(key) {
			return
		}
		tmp := shiftGSettings.GetString(key)
		value := getSystemValue(key, false)
		if tmp == value {
			return
		}
		UpdateSystemShortcut(key, tmp)
	})

	putGSettings.Connect("changed", func(s *gio.Settings, key string) {
		if !keyIsExist(key) {
			return
		}
		tmp := putGSettings.GetString(key)
		value := getSystemValue(key, false)
		if tmp == value {
			return
		}
		UpdateSystemShortcut(key, tmp)
	})
}

func (m *BindManager) listenCustom() {
	customList := getCustomList()
	for _, k := range customList {
		gs := newGSettingsById(k)
		if gs == nil {
			continue
		}
		IdGSettingsMap[k] = gs

		gs.Connect("changed::shortcut", func(s *gio.Settings, key string) {
			m.setPropList("CustomList")
			grabKeyPairs(CustomPrevPairs, false)
			grabKeyPairs(getCustomPairs(), true)
		})
	}
}

func (m *BindManager) listenConflict() {
	bindGSettings.Connect("changed::conflict-valid", func(s *gio.Settings, key string) {
		m.setPropList("ConflictValid")
	})

	bindGSettings.Connect("changed::conflict-invalid", func(s *gio.Settings, key string) {
		m.setPropList("ConflictInvalid")
	})
}

func getSystemKeyInfo() []ShortcutInfo {
	systemInfoList := []ShortcutInfo{}
	for i, n := range SystemIdNameMap {
		if desc, ok := SystemNameDescMap[n]; ok {
			shortcut := getSystemValue(n, false)
			tmp := newShortcutInfo(i, desc,
				formatShortcut(shortcut))
			systemInfoList = append(systemInfoList, tmp)
		}
	}

	return systemInfoList
}

func getMediaKeyInfo() []ShortcutInfo {
	mediaInfoList := []ShortcutInfo{}
	for i, n := range MediaIdNameMap {
		if desc, ok := MediaNameDescMap[n]; ok {
			shortcut := getSystemValue(n, false)
			tmp := newShortcutInfo(i, desc, shortcut)
			mediaInfoList = append(mediaInfoList, tmp)
		}
	}

	return mediaInfoList
}

func getWindowKeyInfo() []ShortcutInfo {
	windowInfoList := []ShortcutInfo{}
	for i, n := range WindowIdNameMap {
		if desc, ok := WindowNameDescMap[n]; ok {
			shortcut := getSystemValue(n, false)
			tmp := newShortcutInfo(i, desc,
				formatShortcut(shortcut))
			windowInfoList = append(windowInfoList, tmp)
		}
	}

	return windowInfoList
}

func getWorkSpaceKeyInfo() []ShortcutInfo {
	workSpaceInfoList := []ShortcutInfo{}
	for i, n := range WorkSpaceIdNameMap {
		if desc, ok := WorkSpaceNameDescMap[n]; ok {
			shortcut := getSystemValue(n, false)
			tmp := newShortcutInfo(i, desc,
				formatShortcut(shortcut))
			workSpaceInfoList = append(workSpaceInfoList, tmp)
		}
	}

	return workSpaceInfoList
}

func getCustomKeyInfo() []ShortcutInfo {
	customList := getCustomList()
	shortcutInfoList := []ShortcutInfo{}

	for _, k := range customList {
		tmp := ShortcutInfo{}
		gs := newGSettingsById(k)
		if gs == nil {
			continue
		}
		tmp.Id = k
		tmp.Desc = getCustomValue(gs, _CUSTOM_KEY_NAME)
		tmp.Shortcut = formatShortcut(
			getCustomValue(gs, _CUSTOM_KEY_SHORTCUT))
		shortcutInfoList = append(shortcutInfoList, tmp)
	}

	return shortcutInfoList
}

func getConflictList(valid bool) []int32 {
	list := []int32{}

	if !valid {
		invalidList := bindGSettings.GetStrv(_BINDING_INVALID_LIST)
		for _, k := range invalidList {
			tmp, err := strconv.ParseInt(k, 10, 64)
			if err != nil {
				continue
			}
			list = append(list, int32(tmp))
		}
	} else {
		validList := bindGSettings.GetStrv(_BINDING_VALID_LIST)
		for _, k := range validList {
			tmp, err := strconv.ParseInt(k, 10, 64)
			if err != nil {
				continue
			}
			list = append(list, int32(tmp))
		}
	}

	return list
}

func getSystemValue(name string, action bool) string {
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

func getIdByName(name string) int32 {
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

/* Update Shortcut */
func UpdateSystemShortcut(key, value string) {
	values := systemGSettings.GetStrv(key)
	if len(values) <= 0 {
		systemGSettings.SetStrv(key, []string{value})
		//gio.SettingsSync()
		return
	}
	values[0] = value

	systemGSettings.SetStrv(key, values)
	//gio.SettingsSync()
}

func updateCompizValue(id int32, key, shortcut string) {
	if id >= 600 && id < 800 {
		if isInvalidConflict(id) {
			wmGSettings.SetStrv(key, []string{})
			return
		}
		values := wmGSettings.GetStrv(key)
		if len(values) >= 1 {
			if values[0] == shortcut {
				return
			}
		}
		wmGSettings.SetStrv(key, []string{shortcut})
	} else if id >= 800 && id < 900 {
		if isInvalidConflict(id) {
			shiftGSettings.SetStrv(key, []string{})
			return
		}
		value := shiftGSettings.GetString(key)
		if value == shortcut {
			return
		}
		shiftGSettings.SetString(key, shortcut)
	} else if id >= 900 && id < 1000 {
		if isInvalidConflict(id) {
			putGSettings.SetStrv(key, []string{})
			return
		}
		value := putGSettings.GetString(key)
		if value == shortcut {
			return
		}
		putGSettings.SetString(key, shortcut)
	}
	//gio.SettingsSync()
}

func keyIsExist(key string) bool {
	for _, v := range IdNameMap {
		if v == key {
			return true
		}
	}

	return false
}

func newGSettingsById(id int32) *gio.Settings {
	if id < _CUSTOM_ID_BASE {
		fmt.Println("not custom id range")
		return nil
	}

	str := strconv.FormatInt(int64(id), 10) + "/"
	gs := gio.NewSettingsWithPath(_CUSTOM_ADD_SCHEMA_ID,
		_CUSTOM_ADD_SCHEMA_PATH+str)

	return gs
}

func setCustomValues(gs *gio.Settings,
	id int32, name, action, shortcut string) {
	gs.SetInt(_CUSTOM_KEY_ID, int(id))
	gs.SetString(_CUSTOM_KEY_NAME, name)
	gs.SetString(_CUSTOM_KEY_ACTION, action)
	gs.SetString(_CUSTOM_KEY_SHORTCUT, shortcut)
	//gio.SettingsSync()
}

func resetCustomValues(gs *gio.Settings) {
	gs.Reset(_CUSTOM_KEY_ID)
	gs.Reset(_CUSTOM_KEY_NAME)
	gs.Reset(_CUSTOM_KEY_ACTION)
	gs.Reset(_CUSTOM_KEY_SHORTCUT)
	//gio.SettingsSync()
}

func getMaxIdFromCustom() int32 {
	max := int32(0)
	customList := getCustomList()

	for _, k := range customList {
		if max < k {
			max = k
		}
	}

	return max
}

func getCustomList() []int32 {
	customList := []int32{}
	strList := bindGSettings.GetStrv(_BINDING_CUSTOM_LIST)

	for _, k := range strList {
		id, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			fmt.Println(err)
			continue
		}
		customList = append(customList, int32(id))
	}

	return customList
}

func getCustomValue(gs *gio.Settings, key string) string {
	return gs.GetString(key)
}

func getCustomAccels() map[int32]string {
	customList := getCustomList()
	customAccels := make(map[int32]string)

	for _, k := range customList {
		gs := newGSettingsById(k)
		if gs == nil {
			continue
		}
		shortcut := getCustomValue(gs, _CUSTOM_KEY_SHORTCUT)
		if len(shortcut) <= 0 {
			continue
		}
		customAccels[k] = shortcut
	}

	return customAccels
}
