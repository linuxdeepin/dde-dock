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
	"strings"
)

var _manager *Manager

func (obj *Manager) Reset() bool {
	list := sysGSettings.ListKeys()
	for _, key := range list {
		sysGSettings.Reset(key)
	}

	for _, id := range obj.ConflictInvalid {
		if id >= CUSTOM_KEY_ID_BASE {
			continue
		}
		obj.ModifyShortcut(id, "")
	}

	return true
}

func (obj *Manager) AddCustomShortcut(name, action string) (int32, bool) {
	id := getMaxCustomId() + 1

	if !obj.createCustomShortcut(id, name, action, "") {
		return -1, false
	}

	return id, true
}

func (obj *Manager) AddCustomShortcutCheck(name, action, shortcut string) (int32, string, []int32) {
	id, ok := obj.AddCustomShortcut(name, action)
	if !ok {
		return -1, "failed", []int32{}
	}

	str, list := obj.CheckShortcutConflict(shortcut)

	return id, str, list
}

func (obj *Manager) CheckShortcutConflict(shortcut string) (string, []int32) {
	if !isValidShortcut(shortcut) {
		return "Invalid", []int32{}
	}

	isConflict, idList := conflictChecked(-1, shortcut)
	if isConflict {
		return "Conflict", idList
	}

	return "Valid", []int32{}
}

func (obj *Manager) ModifyShortcut(id int32, shortcut string) (string, []int32) {
	tmpStr := strings.ToLower(shortcut)
	if tmpStr == "super" || tmpStr == "super_l" || tmpStr == "super_r" ||
		tmpStr == "super-super_l" || tmpStr == "super-super_r" {
		// Compiz 不支持单按键
		if id >= 300 && id < 1000 {
			return "Invalid", []int32{}
		}
	}

	tmpAccel := getShortcutById(id)
	tmpConflict, tmpList := conflictChecked(id, tmpAccel)
	if tmpConflict {
		for _, k := range tmpList {
			deleteValidConflictId(k)
			deleteInvalidConflictId(k)
		}
	}

	retStr := ""
	retList := []int32{}

	if !isValidShortcut(shortcut) {
		addInvalidConflictId(id)
		retStr = "Invalid"
	} else if len(shortcut) < 1 {
		retStr = "Valid"
	} else {
		isConflict, idList := conflictChecked(id, shortcut)
		logger.Infof("'%s' isConflict: %v, idList: %v", shortcut, isConflict, idList)
		if isConflict {
			addInvalidConflictId(id)
			for _, k := range idList {
				addValidConflictId(k)
			}
			retStr = "Conflict"
			retList = idList
		} else {
			deleteValidConflictId(id)
			deleteInvalidConflictId(id)
			retStr = "Valid"
		}
	}

	modifyShortcutById(id, shortcut)

	return retStr, retList
}

func (obj *Manager) DeleteCustomShortcut(id int32) {
	tmpKey := getShortcutById(id)
	tmpConflict, tmpList := conflictChecked(id, tmpKey)
	if tmpConflict {
		for _, k := range tmpList {
			if k == id {
				continue
			}
			deleteValidConflictId(k)
			deleteInvalidConflictId(k)
		}
	}
	deleteValidConflictId(id)
	deleteInvalidConflictId(id)

	obj.deleteCustomShortcut(id)
}

func (obj *Manager) GrabSignalShortcut(shortcut, action string, isGrab bool) {
	grabSignalShortcut(shortcut, action, isGrab)
}

func (obj *Manager) GrabKbdAndMouse() {
	go grabKeyboardAndMouse()
}

func newManager() *Manager {
	obj := &Manager{}
	obj.idSettingsMap = make(map[int32]*gio.Settings)

	obj.listenKeyEvents()
	obj.listenSettings()
	obj.listenCompizSettings()
	obj.listenAllCustomSettings()
	obj.updateProps()

	return obj
}

func GetManager() *Manager {
	if _manager == nil {
		_manager = newManager()
	}

	return _manager
}
