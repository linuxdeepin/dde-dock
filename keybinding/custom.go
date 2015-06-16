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

package keybinding

import (
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"strconv"
	"strings"
)

func addCustomId(id int32) {
	if id < CUSTOM_KEY_ID_BASE {
		return
	}

	idStr := strconv.FormatInt(int64(id), 10)
	list := bindGSettings.GetStrv(BIND_KEY_CUSTOM_LIST)
	if !isStringInList(idStr, list) {
		list = append(list, idStr)
		bindGSettings.SetStrv(BIND_KEY_CUSTOM_LIST, list)
	}
}

func deleteCustomId(id int32) {
	if id < CUSTOM_KEY_ID_BASE {
		return
	}

	idStr := strconv.FormatInt(int64(id), 10)
	list := bindGSettings.GetStrv(BIND_KEY_CUSTOM_LIST)
	if isStringInList(idStr, list) {
		tmpList := []string{}

		for _, v := range list {
			if v == idStr {
				continue
			}
			tmpList = append(tmpList, v)
		}
		bindGSettings.SetStrv(BIND_KEY_CUSTOM_LIST, tmpList)
	}
}

func getCustomIdList() []int32 {
	rets := []int32{}
	list := bindGSettings.GetStrv(BIND_KEY_CUSTOM_LIST)
	for _, v := range list {
		tmp, _ := strconv.ParseInt(v, 10, 64)
		rets = append(rets, int32(tmp))
	}

	return rets
}

func getMaxCustomId() int32 {
	list := getCustomIdList()
	if len(list) < 1 {
		return CUSTOM_KEY_ID_BASE
	}

	max := list[0]
	for _, v := range list {
		if max < v {
			max = v
		}
	}

	return max
}

func getSettingsById(id int32) *gio.Settings {
	if id < CUSTOM_KEY_ID_BASE {
		return nil
	}

	str := strconv.FormatInt(int64(id), 10) + "/"
	gs := gio.NewSettingsWithPath(CUSTOM_KEY_SCHEMA_ID,
		CUSTOM_KEY_BASE_PATH+str)

	return gs
}

func setCustomValue(id int32, key, value string) bool {
	if key == CUSTOM_KEY_NAME ||
		key == CUSTOM_KEY_ACTION ||
		key == CUSTOM_KEY_SHORTCUT {

		logger.Infof("Set id: %d, key : %s, value: %s", id, key, value)
		gs := getSettingsById(id)
		if gs == nil {
			logger.Errorf("Get GSettings Failed For Id: %v", id)
		}
		gs.SetString(key, value)
	}

	return false
}

func getCustomValue(id int32, key string) string {
	ret := ""
	if key == CUSTOM_KEY_NAME ||
		key == CUSTOM_KEY_ACTION ||
		key == CUSTOM_KEY_SHORTCUT {
		gs := getSettingsById(id)
		if gs == nil {
			logger.Errorf("Get GSettings Failed For Id: %v", id)
		} else {
			ret = gs.GetString(key)
		}
	}

	return ret
}

func getCustomKeyAccels() map[int32]string {
	retMap := make(map[int32]string)
	customList := getCustomIdList()

	for _, k := range customList {
		shortcut := getCustomValue(k, CUSTOM_KEY_SHORTCUT)
		if len(shortcut) < 1 {
			continue
		}
		retMap[k] = strings.ToLower(shortcut)
	}

	return retMap
}

func getCustomListInfo() []ShortcutInfo {
	idList := getCustomIdList()
	list := []ShortcutInfo{}

	for _, k := range idList {
		tmp := ShortcutInfo{}

		tmp.Id = k
		tmp.Desc = getCustomValue(k, CUSTOM_KEY_NAME)
		shortcut := getCustomValue(k, CUSTOM_KEY_SHORTCUT)
		tmp.Shortcut = formatShortcut(shortcut)
		list = append(list, tmp)
	}

	return list
}

func (obj *Manager) createCustomShortcut(id int32, name, action, shortcut string) bool {
	gs := getSettingsById(id)
	if gs == nil {
		return false
	}

	gs.SetInt(CUSTOM_KEY_ID, id)
	gs.SetString(CUSTOM_KEY_NAME, name)
	gs.SetString(CUSTOM_KEY_ACTION, action)
	gs.SetString(CUSTOM_KEY_SHORTCUT, shortcut)

	addCustomId(id)

	obj.listenCustomSettings(id)

	return true
}

func (obj *Manager) deleteCustomShortcut(id int32) {
	gs := getSettingsById(id)
	if gs == nil {
		return
	}

	gs.Reset(CUSTOM_KEY_ID)
	gs.Reset(CUSTOM_KEY_NAME)
	gs.Reset(CUSTOM_KEY_ACTION)
	gs.Reset(CUSTOM_KEY_SHORTCUT)

	deleteCustomId(id)
	delete(obj.idSettingsMap, id)
}
