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

package inputdevices

import (
	"strings"
)

func (mManager *MouseManager) Reset() bool {
	list := mouseSettings.ListKeys()
	for _, key := range list {
		mouseSettings.Reset(key)
	}

	return true
}

func (tManager *TouchpadManager) Reset() bool {
	list := tpadSettings.ListKeys()
	for _, key := range list {
		tpadSettings.Reset(key)
	}

	return true
}

func (kbdManager *KeyboardManager) Reset() bool {
	list := kbdSettings.ListKeys()
	for _, key := range list {
		kbdSettings.Reset(key)
	}

	return true
}

func (kbdManager *KeyboardManager) LayoutList() map[string]string {
	defer func() {
		if err := recover(); err != nil {
			logger.Warning("Recover error:", err)
			return
		}
	}()

	if kbdManager.layoutDescMap == nil {
		kbdManager.layoutDescMap = make(map[string]string)
		kbdManager.layoutDescMap = kbdManager.getLayoutList()
	}

	return kbdManager.layoutDescMap
}

func (kbdManager *KeyboardManager) GetLayoutLocale(layout string) string {
	if len(layout) < 1 || !strings.Contains(layout, LAYOUT_DELIM) {
		layout = "us" + LAYOUT_DELIM
	}

	listMap := kbdManager.LayoutList()
	desc, ok := listMap[layout]
	if !ok {
		logger.Debug("Invalid layout:", layout)
		return ""
	}

	return desc
}

func (kbdManager *KeyboardManager) AddLayoutOption(option string) {
	if len(option) < 1 {
		return
	}

	options := kbdSettings.GetStrv(KBD_KEY_LAYOUT_OPTIONS)
	if !kbdManager.isStrInList(option, options) {
		options = append(options, option)
		kbdSettings.SetStrv(KBD_KEY_LAYOUT_OPTIONS, options)
	}
}

func (kbdManager *KeyboardManager) DeleteLayoutOption(option string) {
	if len(option) < 1 {
		return
	}

	options := kbdSettings.GetStrv(KBD_KEY_LAYOUT_OPTIONS)
	if kbdManager.isStrInList(option, options) {
		tmpList := []string{}
		for _, v := range options {
			if v != option {
				tmpList = append(tmpList, v)
			}
		}
		kbdSettings.SetStrv(KBD_KEY_LAYOUT_OPTIONS, tmpList)
	}
}

func (kbdManager *KeyboardManager) ClearLayoutOption() {
	kbdSettings.Reset(KBD_KEY_LAYOUT_OPTIONS)
}

func (kbdManager *KeyboardManager) AddUserLayout(layout string) bool {
	if len(layout) < 1 || !strings.Contains(layout, LAYOUT_DELIM) {
		layout = "us" + LAYOUT_DELIM
	}

	list := kbdManager.UserLayoutList.Get()
	if kbdManager.isStrInList(layout, list) {
		return true
	}

	list = append(list, layout)
	kbdSettings.SetStrv(KBD_KEY_USER_LAYOUT_LIST, list)
	return true
}

func (kbdManager *KeyboardManager) DeleteUserLayout(layout string) bool {
	if len(layout) < 1 || !strings.Contains(layout, LAYOUT_DELIM) {
		return false
	}

	list := kbdManager.UserLayoutList.Get()
	if !kbdManager.isStrInList(layout, list) {
		return false
	}

	tmpList := []string{}
	for _, v := range list {
		if v != layout {
			tmpList = append(tmpList, v)
		}
	}
	kbdSettings.SetStrv(KBD_KEY_USER_LAYOUT_LIST, tmpList)

	return true
}
