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

package main

import (
	"github.com/BurntSushi/xgbutil/keybind"
	"strconv"
	"strings"
)

func isValidConflict(id int32) bool {
	validList := bindGSettings.GetStrv(BIND_KEY_VALID_LIST)

	for _, v := range validList {
		tmp, _ := strconv.ParseInt(v, 10, 64)
		if id == int32(tmp) {
			return true
		}
	}

	return false
}

func addValidConflictId(id int32) {
	if isValidConflict(id) {
		return
	}

	if isInvalidConflict(id) {
		deleteInvalidConflictId(id)
	}
	validList := bindGSettings.GetStrv(BIND_KEY_VALID_LIST)
	idStr := strconv.FormatInt(int64(id), 10)
	validList = append(validList, idStr)
	bindGSettings.SetStrv(BIND_KEY_VALID_LIST, validList)
}

func deleteValidConflictId(id int32) {
	if !isValidConflict(id) {
		return
	}

	validList := bindGSettings.GetStrv(BIND_KEY_VALID_LIST)
	tmpList := []string{}
	for _, v := range validList {
		tmp, _ := strconv.ParseInt(v, 10, 64)
		if id == int32(tmp) {
			continue
		}
		tmpList = append(tmpList, v)
	}
	bindGSettings.SetStrv(BIND_KEY_VALID_LIST, tmpList)
}

func isInvalidConflict(id int32) bool {
	invalidList := bindGSettings.GetStrv(BIND_KEY_INVALID_LIST)

	for _, v := range invalidList {
		tmp, _ := strconv.ParseInt(v, 10, 64)
		if id == int32(tmp) {
			return true
		}
	}

	return false
}

func addInvalidConflictId(id int32) {
	if isInvalidConflict(id) {
		return
	}

	if isValidConflict(id) {
		deleteValidConflictId(id)
	}
	invalidList := bindGSettings.GetStrv(BIND_KEY_INVALID_LIST)
	idStr := strconv.FormatInt(int64(id), 10)
	invalidList = append(invalidList, idStr)
	bindGSettings.SetStrv(BIND_KEY_INVALID_LIST, invalidList)
}

func deleteInvalidConflictId(id int32) {
	if !isInvalidConflict(id) {
		return
	}

	invalidList := bindGSettings.GetStrv(BIND_KEY_INVALID_LIST)
	tmpList := []string{}
	for _, v := range invalidList {
		tmp, _ := strconv.ParseInt(v, 10, 64)
		if id == int32(tmp) {
			continue
		}
		tmpList = append(tmpList, v)
	}
	bindGSettings.SetStrv(BIND_KEY_INVALID_LIST, tmpList)
}

func conflictChecked(id int32, shortcut string) (bool, []int32) {
	if len(shortcut) < 1 {
		return false, []int32{}
	}

	tmpKey := strings.ToLower(shortcut)
	var (
		info KeycodeInfo
		ok   bool
	)

	logObj.Info("Check Conflict:", shortcut)
	if tmpKey == "super-super_l" || tmpKey == "super-super_r" ||
		tmpKey == "super" {
		info, ok = newKeycodeInfo("Super_L")
	} else {
		info, ok = newKeycodeInfo(formatXGBShortcut(shortcut))
	}

	if !ok {
		logObj.Warning("new keycode failed:", shortcut)
		return false, []int32{}
	}

	isConflict := false
	idList := []int32{}
	allShortcut := getAllAccels()
	for i, k := range allShortcut {
		if i == id {
			continue
		}

		var tmp KeycodeInfo
		if k == "super" {
			tmp, ok = newKeycodeInfo("Super_L")
		} else {
			tmp, ok = newKeycodeInfo(formatXGBShortcut(k))
		}
		if !ok {
			continue
		}

		if isKeycodeInfoEqual(&info, &tmp) {
			isConflict = true
			idList = append(idList, i)
		}
	}

	return isConflict, idList
}

func getSystemListInfo() []ShortcutInfo {
	list := []ShortcutInfo{}

	for i, n := range SystemIdNameMap {
		if desc, ok := SystemNameDescMap[n]; ok {
			shortcut := getSystemKeyValue(n, false)
			tmp := newShortcutInfo(i, desc,
				formatShortcut(shortcut))
			tmp.index = SystemIdIndexMap[i]
			list = append(list, tmp)
		}
	}

	return list
}

func getWindowListInfo() []ShortcutInfo {
	list := []ShortcutInfo{}

	for i, n := range WindowIdNameMap {
		if desc, ok := WindowNameDescMap[n]; ok {
			shortcut := getSystemKeyValue(n, false)
			tmp := newShortcutInfo(i, desc,
				formatShortcut(shortcut))
			tmp.index = WindowIdIndexMap[i]
			list = append(list, tmp)
		}
	}

	return list
}

func getWorkspaceListInfo() []ShortcutInfo {
	list := []ShortcutInfo{}

	for i, n := range WorkspaceIdNameMap {
		if desc, ok := WorkspaceNameDescMap[n]; ok {
			shortcut := getSystemKeyValue(n, false)
			tmp := newShortcutInfo(i, desc,
				formatShortcut(shortcut))
			tmp.index = WorkspaceIdIndexMap[i]
			list = append(list, tmp)
		}
	}

	return list
}

func getValidConflictList() []int32 {
	list := []int32{}

	validList := bindGSettings.GetStrv(BIND_KEY_VALID_LIST)
	for _, v := range validList {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			continue
		}
		list = append(list, int32(id))
	}

	return list
}

func getInvalidConflictList() []int32 {
	list := []int32{}

	invalidList := bindGSettings.GetStrv(BIND_KEY_INVALID_LIST)
	for _, v := range invalidList {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			continue
		}
		list = append(list, int32(id))
	}

	return list
}

func isValidShortcut(shortcut string) bool {
	tmp := formatShortcut(shortcut)
	if len(tmp) == 0 || strings.Contains(tmp, "-") {
		as := strings.Split(tmp, "-")
		str := as[len(as)-1]
		// 修饰键作为单按键的情况
		if strings.Contains(str, "alt") ||
			strings.Contains(str, "shift") ||
			strings.Contains(str, "control") {
			return false
		} else {
			return true
		}
	}

	switch tmp {
	case "f1", "f2", "f3", "f4", "f5", "f6", "f7", "f8", "f9", "f10", "f11", "f12", "print", "super_l", "super_r":
		return true
	}

	return false
}

func getShortcutById(id int32) string {
	if id >= CUSTOM_KEY_ID_BASE {
		return getCustomValue(id, CUSTOM_KEY_SHORTCUT)
	}

	if key, ok := IdNameMap[id]; ok {
		return getSystemKeyValue(key, false)
	}

	return ""
}

func setSystemValue(id int32, value string, action bool) {
	if key, ok := IdNameMap[id]; ok {
		list := sysGSettings.GetStrv(key)
		if len(list) > 1 && action {
			list[1] = value
		} else if len(list) > 0 && !action {
			list[0] = value
		}
		sysGSettings.SetStrv(key, list)
	}
}

func modifyShortcutById(id int32, shortcut string) {
	logObj.Infof("Id: %d, shortcut: %s", id, shortcut)
	if id >= CUSTOM_KEY_ID_BASE {
		logObj.Info("---Set Custom key")
		setCustomValue(id, CUSTOM_KEY_SHORTCUT, shortcut)
		return
	}

	setSystemValue(id, shortcut, false)
}

func newShortcutInfo(id int32, desc, shortcut string) ShortcutInfo {
	return ShortcutInfo{Id: id, Desc: desc, Shortcut: shortcut}
}

func newKeycodeInfo(shortcut string) (KeycodeInfo, bool) {
	if len(shortcut) < 1 {
		return KeycodeInfo{}, false
	}

	mod, keys, err := keybind.ParseString(X, shortcut)
	if err != nil {
		logObj.Warningf("newKeycodeInfo parse '%s' failed: %v",
			shortcut, err)
		return KeycodeInfo{}, false
	}

	if len(keys) < 1 {
		logObj.Warningf("'%s' no details", shortcut)
		return KeycodeInfo{}, false
	}

	state, detail := keybind.DeduceKeyInfo(mod, keys[0])

	return KeycodeInfo{State: state, Detail: detail}, true
}

func isKeycodeInfoEqual(info1, info2 *KeycodeInfo) bool {
	if info1 == nil || info2 == nil {
		return false
	}

	if (info1.Detail == info2.Detail) && (info1.State == info2.State) {
		logObj.Infof("Info1: %v -- %v", info1.Detail, info1.State)
		logObj.Infof("Info2: %v -- %v", info2.Detail, info2.State)
		return true
	}

	return false
}

func getSystemKeyValue(key string, action bool) string {
	values := sysGSettings.GetStrv(key)
	l := len(values)
	if l < 1 {
		return ""
	}

	if l == 2 {
		if action {
			return values[1]
		}
	}

	return strings.ToLower(values[0])
}

func getAllAccels() map[int32]string {
	allMap := make(map[int32]string)

	for i, n := range IdNameMap {
		// Mediakey
		if i >= 300 && i < 500 {
			continue
		}

		shortcut := getSystemKeyValue(n, false)
		if len(shortcut) < 1 {
			continue
		}
		allMap[i] = strings.ToLower(shortcut)
	}

	customMap := getCustomKeyAccels()
	for k, v := range customMap {
		allMap[k] = v
	}

	return allMap
}

func getAccelIdByName(name string) (int32, bool) {
	if len(name) < 1 {
		return -1, false
	}

	for i, n := range IdNameMap {
		if name == n {
			return i, true
		}
	}

	return -1, false
}

func compareInt32List(l1, l2 []int32) bool {
	len1 := len(l1)
	len2 := len(l2)

	if len1 != len2 {
		return false
	}

	for i := 0; i < len1; i++ {
		if l1[i] != l2[i] {
			return false
		}
	}

	return true
}

func compareShortcutInfo(info1, info2 *ShortcutInfo) bool {
	if info1 == nil || info2 == nil {
		return false
	}

	if info1.Desc != info2.Desc ||
		info1.Id != info2.Id ||
		info1.Shortcut != info2.Shortcut ||
		info1.index != info2.index {
		return false
	}

	return true
}

func compareShortcutInfoList(l1, l2 []ShortcutInfo) bool {
	len1 := len(l1)
	len2 := len(l2)

	if len1 != len2 {
		return false
	}

	for i := 0; i < len1; i++ {
		if !compareShortcutInfo(&l1[i], &l2[i]) {
			return false
		}
	}

	return true
}

func sortShortcutList(list []ShortcutInfo) []ShortcutInfo {
	l := len(list)
	orderList := make([]ShortcutInfo, l)

	for _, v := range list {
		orderList[v.index] = v
	}

	return orderList
}

func isStringInList(str string, list []string) bool {
	for _, v := range list {
		if str == v {
			return true
		}
	}

	return false
}
