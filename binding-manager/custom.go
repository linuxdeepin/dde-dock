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
	"dlib/gio-2.0"
	"dlib/dbus"
	"fmt"
	"strconv"
)

func InitCustomBind(m *BindManager) {
	m.CustomList = GetCustomKeyInfo()
}

func GetCustomKeyInfo() []*ShortcutInfo {
	customList := GetCustomList()
	shortcutInfoList := []*ShortcutInfo{}

	for _, k := range customList {
		tmp := &ShortcutInfo{}
		gs := NewGSettingsById(k)
		if gs == nil {
			continue
		}
		tmp.Id = k
		tmp.Desc = GetCustomValue(gs, _CUSTOM_KEY_NAME)
		tmp.Shortcut = FormatShortcut(
			GetCustomValue(gs, _CUSTOM_KEY_SHORTCUT))
		shortcutInfoList = append(shortcutInfoList, tmp)
	}

	return shortcutInfoList
}

func NewGSettingsById(id int32) *gio.Settings {
	if id < _CUSTOM_ID_BASE {
		fmt.Println("not custom id range")
		return nil
	}

	str := strconv.FormatInt(int64(id), 10) + "/"
	gs := gio.NewSettingsWithPath(_CUSTOM_ADD_SCHEMA_ID,
		_CUSTOM_ADD_SCHEMA_PATH+str)

	return gs
}

func SetCustomValues(gs *gio.Settings,
	id int32, name, action, shortcut string) {
	gs.SetInt(_CUSTOM_KEY_ID, int(id))
	gs.SetString(_CUSTOM_KEY_NAME, name)
	gs.SetString(_CUSTOM_KEY_ACTION, action)
	gs.SetString(_CUSTOM_KEY_SHORTCUT, shortcut)
	gio.SettingsSync()
}

func ResetCustomValues(gs *gio.Settings) {
	gs.Reset(_CUSTOM_KEY_ID)
	gs.Reset(_CUSTOM_KEY_NAME)
	gs.Reset(_CUSTOM_KEY_ACTION)
	gs.Reset(_CUSTOM_KEY_SHORTCUT)
	gio.SettingsSync()
}

func GetMaxIdFromCustom() int32 {
	max := int32(0)
	customList := GetCustomList()

	for _, k := range customList {
		if max < k {
			max = k
		}
	}

	return max
}

func GetCustomList() []int32 {
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

func GetCustomValue(gs *gio.Settings, key string) string {
	return gs.GetString(key)
}

func GetCustomAccels() map[int32]string {
	customList := GetCustomList()
	customAccels := make(map[int32]string)

	for _, k := range customList {
		gs := NewGSettingsById(k)
		if gs == nil {
			continue
		}
		shortcut := GetCustomValue(gs, _CUSTOM_KEY_SHORTCUT)
		if len(shortcut) <= 0 {
			continue
		}
		customAccels[k] = shortcut
	}

	return customAccels
}

func ListenCustom(m *BindManager) {
	customList := GetCustomList()
	for _, k := range customList {
		gs := NewGSettingsById(k)
		if gs == nil {
			continue
		}
		IdGSettingsMap[k] = gs

		gs.Connect("changed::shortcut", func(s *gio.Settings, key string) {
			m.CustomList = GetCustomKeyInfo()
			dbus.NotifyChange(m, "CustomList")
                        GrabKeyPairs(CustomPrevPairs, false)
                        GrabKeyPairs(GetCustomPairs(), true)
		})
	}
}
