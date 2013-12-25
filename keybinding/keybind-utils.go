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
	"fmt"
	"github.com/BurntSushi/xgbutil/keybind"
	"strconv"
	"strings"
)

func NewCustomGSettings(id int32) *gio.Settings {
	customId := strconv.FormatInt(int64(id), 10) + "/"
	gs := gio.NewSettingsWithPath(_CUSTOM_SCHEMA_ADD_ID, _CUSTOM_SCHEMA_ADD_PATH+customId)

	return gs
}

func GetMaxIdFromCustom() int32 {
	customList := GetCustomIdList()
	max := int32(0)

	for _, v := range customList {
		if max < v {
			max = v
		}
	}

	return max
}

func GetGSDPairs() map[string]string {
	gsdPairs := make(map[string]string)

	for k, _ := range gsdMap {
		v := GSDGetValue(k)
		strs := strings.Split(v, ";")
		if len(strs) == 2 {
			shortcut := FormatShortcut(strs[1])
			gsdPairs[shortcut] = strs[0]
		}
	}

	return gsdPairs
}

func GetCustomPairs() map[string]string {
	customPairs := make(map[string]string)

	customList := GetCustomIdList()
	for _, v := range customList {
		gs := NewCustomGSettings(v)
		tmp := gs.GetString(_CUSTOM_KEY_SHORTCUT)
		shortcut := FormatShortcut(tmp)
		action := gs.GetString(_CUSTOM_KEY_ACTION)
		customPairs[shortcut] = action
	}

	return customPairs
}

func GetKeyAccelList() map[int32]*GrabKeyInfo {
	accelList := make(map[int32]*GrabKeyInfo)

	value := ""
	flag := false
	for k, _ := range currentSystemBindings {
		/*else if k >= 300 && k < 600 {
			values := MediaGetValue(k)
			accelList[k] = values
		} */
		if k >= 0 && k < 300 {
			values := GSDGetValue(k)
			strArray := strings.Split(values, ";")
			if len(strArray) == 2 {
				value = strArray[1]
				flag = true
			}
		} else if k >= 600 && k < 800 {
			value = WMGetValue(k)
			flag = true
		} else if k >= 800 && k < 900 {
			value = CompizShiftValue(k)
			flag = true
		} else if k >= 900 && k < 1000 {
			value = CompizPutValue(k)
			flag = true
		}
		if flag {
			flag = false
			tmp := GenericKeyInfo(value)
			if tmp == nil {
				continue
			}
			accelList[k] = tmp
		}
	}

	customList := GetCustomIdList()
	for _, k := range customList {
		value = GetCustomValue(k, _CUSTOM_KEY_SHORTCUT)
		tmp := GenericKeyInfo(value)
		if tmp == nil {
			continue
		}
		accelList[k] = tmp
	}

	return accelList
}

func GenericKeyInfo(shortcut string) *GrabKeyInfo {
	tmp1 := FormatShortcut(strings.ToLower(shortcut))
	tmp2 := GetXGBShortcut(tmp1)
	mod, keys, _ := keybind.ParseString(X, tmp2)
	if len(keys) <= 0 {
		return nil
	}

	return NewKeyInfo(mod, keys[0])
}

func GetCustomIdList() []int32 {
	customIDList := []int32{}
	strList := customGSettings.GetStrv(_CUSTOM_KEY_LIST)

	for _, v := range strList {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			fmt.Println("get custom list failed:", err)
			continue
		}
		customIDList = append(customIDList, int32(id))
	}

	return customIDList
}

func GSDGetValue(id int32) string {
	if _, ok := gsdMap[id]; ok {
		keyName := currentSystemBindings[id]

		return gsdGSettings.GetString(keyName)
	}

	return ""
}

func MediaGetValue(id int32) string {
	if _, ok := mediaMap[id]; ok {
		keyName := currentSystemBindings[id]

		return mediaGSettings.GetString(keyName)
	}

	return ""
}

func WMGetValue(id int32) string {
	if _, ok := wmMap[id]; ok {
		keyName := currentSystemBindings[id]

		values := wmGSettings.GetStrv(keyName)
		strRet := ""

		for _, v := range values {
			strRet += v
		}
		return strRet
	}

	return ""
}

func CompizShiftValue(id int32) string {
	if _, ok := shiftMap[id]; ok {
		keyName := currentSystemBindings[id]
		values := shiftGSettings.GetString(keyName)

		return values
	}

	return ""
}

func CompizPutValue(id int32) string {
	if _, ok := putMap[id]; ok {
		keyName := currentSystemBindings[id]
		values := putGSettings.GetString(keyName)

		return values
	}

	return ""
}
