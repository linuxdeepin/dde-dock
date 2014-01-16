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
	//"dlib/gio-2.0"
	"fmt"
	"github.com/BurntSushi/xgbutil/keybind"
	"strconv"
	"strings"
)

func newShortcutInfo(id int32, desc, shortcut string) ShortcutInfo {
	return ShortcutInfo{Id: id, Desc: desc, Shortcut: shortcut}
}

func newKeyCodeInfo(shortcut string) *KeyCodeInfo {
	mods, keys, err := keybind.ParseString(X, shortcut)
	if err != nil {
		fmt.Println("keybind parse string failed: ", err)
		return nil
	}

	if len(keys) <= 0 {
		fmt.Println("no key")
		return nil
	}

	state, detail := keybind.DeduceKeyInfo(mods, keys[0])
	return &KeyCodeInfo{State: state, Detail: detail}
}

func keyCodeInfoEqual(keyInfo1, keyInfo2 *KeyCodeInfo) bool {
	if keyInfo1 == nil || keyInfo2 == nil {
		return false
	}

	if keyInfo1.State == keyInfo2.State &&
		keyInfo1.Detail == keyInfo2.Detail {
		return true
	}

	return false
}

func modifyShortcutById(id int32, shortcut string) {
        tmpKey := strings.ToLower(shortcut)
        var accel string

        if tmpKey == "super" || tmpKey == "super_l" ||
        tmpKey == "super-super_l" {
                accel = "Super"
        } else {
                accel = shortcut
        }
	if id >= _CUSTOM_ID_BASE {
		gs := newGSettingsById(id)
		if gs != nil {
			gs.SetString(_CUSTOM_KEY_SHORTCUT,accel)
			//gio.SettingsSync()
		}

		return
	}

	if key, ok := IdNameMap[id]; ok {
		UpdateSystemShortcut(key,accel)
		return
	}
}

func getShortcutById(id int32) string {
	if id >= _CUSTOM_ID_BASE {
		gs := newGSettingsById(id)
		if gs != nil {
			return gs.GetString(_CUSTOM_KEY_SHORTCUT)
		}
	}

	value := ""
	if key, ok := IdNameMap[id]; ok {
		value = getSystemValue(key, false)
	}

	return value
}

func getAllShortcuts() map[int32]string {
	allShortcuts := make(map[int32]string)

	for i, n := range IdNameMap {
		shortcut := getSystemValue(n, false)
		if len(shortcut) <= 0 {
			continue
		}
		if i >= 300 && i < 500 {
			continue
		}
		allShortcuts[i] = strings.ToLower(shortcut)
	}

	customShortcuts := getCustomAccels()
	for k, v := range customShortcuts {
		allShortcuts[k] = strings.ToLower(v)
	}

	return allShortcuts
}

func conflictChecked(id int32, shortcut string) ConflictInfo {
	tmpKey := strings.ToLower(shortcut)
	var info *KeyCodeInfo

	if tmpKey == "super-super_l" || tmpKey == "super_l" ||
        tmpKey == "super" {
		info = newKeyCodeInfo("Super_L")
	} else {
		info = newKeyCodeInfo(getXGBShortcut(formatShortcut(shortcut)))
	}
	if info == nil {
		fmt.Println("shortcut invalid. ", shortcut)
		return ConflictInfo{}
	}

	conflict := ConflictInfo{}
	conflict.IsConflict = false

	allShortcuts := getAllShortcuts()
	for i, k := range allShortcuts {
		if i == id {
			continue
		}
		var tmp *KeyCodeInfo

		if k == "super" {
			tmp = newKeyCodeInfo("Super_L")
		} else {
			tmp = newKeyCodeInfo(getXGBShortcut(formatShortcut(k)))
		}
		if tmp == nil {
			continue
		}

		if keyCodeInfoEqual(info, tmp) {
			conflict.IsConflict = true
			conflict.IdList = append(conflict.IdList, i)
		}
	}

	return conflict
}

func isValidConflict(id int32) bool {
	validList := bindGSettings.GetStrv(_BINDING_VALID_LIST)

	for _, k := range validList {
		tmp, _ := strconv.ParseInt(k, 10, 64)
		if id == int32(tmp) {
			return true
		}
	}

	return false
}

func insertConflictValidList(idList []int32) {
	validList := bindGSettings.GetStrv(_BINDING_VALID_LIST)
	for _, k := range idList {
		if isValidConflict(k) || isInvalidConflict(k) {
			continue
		}
		tmp := strconv.FormatInt(int64(k), 10)
		validList = append(validList, tmp)
	}

	bindGSettings.SetStrv(_BINDING_VALID_LIST, validList)
	//gio.SettingsSync()
}

func deleteConflictValidId(id int32) {
	if !isValidConflict(id) {
		return
	}

	tmpList := []string{}
	tmp := strconv.FormatInt(int64(id), 10)
	validList := bindGSettings.GetStrv(_BINDING_VALID_LIST)
	for _, k := range validList {
		if k == tmp {
			continue
		}
		tmpList = append(tmpList, k)
	}
	bindGSettings.SetStrv(_BINDING_VALID_LIST, tmpList)
	//gio.SettingsSync()
}

func isInvalidConflict(id int32) bool {
	invalidList := bindGSettings.GetStrv(_BINDING_INVALID_LIST)

	for _, k := range invalidList {
		tmp, _ := strconv.ParseInt(k, 10, 64)
		if id == int32(tmp) {
			return true
		}
	}

	return false
}

func deleteConflictInvalidId(id int32) {
	if !isInvalidConflict(id) {
		return
	}

	tmpList := []string{}
	tmp := strconv.FormatInt(int64(id), 10)
	invalidList := bindGSettings.GetStrv(_BINDING_INVALID_LIST)
	for _, k := range invalidList {
		if k == tmp {
			continue
		}
		tmpList = append(tmpList, k)
	}
	bindGSettings.SetStrv(_BINDING_INVALID_LIST, tmpList)
	//gio.SettingsSync()
}

func insertConflictInvalidList(id int32) {
	if isInvalidConflict(id) {
		return
	}
	invalidList := bindGSettings.GetStrv(_BINDING_INVALID_LIST)
	tmp := strconv.FormatInt(int64(id), 10)
	invalidList = append(invalidList, tmp)

	bindGSettings.SetStrv(_BINDING_INVALID_LIST, invalidList)
	//gio.SettingsSync()
}

func idIsExist(id int32, idList []int32) bool {
	for _, v := range idList {
		if id == v {
			return true
		}
	}

	return false
}

func getXGBShortcut(shortcut string) string {
	/*str := formatShortcut(shortcut)
	if len(str) <= 0 {
		return ""
	}*/

	value := ""
	array := strings.Split(shortcut, "-")
	for i, v := range array {
		if i != 0 {
			value += "-"
		}

		if v == "alt" || v == "super" ||
			v == "meta" || v == "num_lock" ||
			v == "caps_lock" || v == "hyper" {
			modStr, _ := _ModifierMap[v]
			value += modStr
		} else {
			value += v
		}
	}

	return value
}

/*
 * Input string format: '<Control><Alt>T'
 * Output string format: 'control-alt-t'
 */

func formatShortcut(shortcut string) string {
	l := len(shortcut)

	if l <= 0 {
		fmt.Println("format args null")
		return ""
	}

	str := strings.ToLower(shortcut)
	value := ""
	flag := false
	start := 0
	end := 0

	for i, ch := range str {
		if ch == '<' {
			flag = true
			start = i
		}

		if ch == '>' && flag {
			end = i
			flag = false
			if start != 0 {
				value += "-"
			}

			for j := start + 1; j < end; j++ {
				value += string(str[j])
			}
		}
	}

	if end != l {
		i := 0
		if end > 0 {
			i = end + 1
			value += "-"
		}
		for ; i < l; i++ {
			value += string(str[i])
		}
	}

	array := strings.Split(value, "-")
	value = ""
	for i, v := range array {
		if v == "primary" || v == "control" {
			if !strings.Contains(value, "control") {
				if i != 0 {
					value += "-"
				}

				value += "control"
			}
			continue
		}

		if i != 0 {
			value += "-"
		}

		value += v
	}

	return value
}
