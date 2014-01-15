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

// #cgo pkg-config: x11 xtst glib-2.0
// #include "grab-xrecord.h"
import "C"

import (
	"dlib"
	"dlib/dbus"
	"dlib/gio-2.0"
	"fmt"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"strconv"
	"strings"
)

type AddAccelRet struct {
	Id    int32
	Check ConflictInfo
}

func (m *BindManager) AddKeyBind(name, action, shortcut string) AddAccelRet {
	id := getMaxIdFromCustom() + 1
	gs := newGSettingsById(id)
	if gs == nil {
		return AddAccelRet{}
	}
	IdGSettingsMap[id] = gs

	setCustomValues(gs, id, name, action, "")
	gs.Connect("changed::shortcut", func(s *gio.Settings, key string) {
		fmt.Printf("key: %s, value: %s\n", key, gs.GetString(key))
		m.setPropList("CustomList")
		grabKeyPairs(CustomPrevPairs, false)
		grabKeyPairs(getCustomPairs(), true)
	})
	ret := AddAccelRet{}
	ret.Id = id
	ret.Check = m.ChangeShortcut(id, shortcut)

	idStr := strconv.FormatInt(int64(id), 10)
	customList := bindGSettings.GetStrv(_BINDING_CUSTOM_LIST)
	customList = append(customList, idStr)
	bindGSettings.SetStrv(_BINDING_CUSTOM_LIST, customList)
	//gio.SettingsSync()

	grabKeyPairs(CustomPrevPairs, false)
	grabKeyPairs(getCustomPairs(), true)
	return ret
}

func KeyIsValid(key string) bool {
	tmp := formatShortcut(key)
	if strings.Contains(tmp, "-") {
		return true
	}

	fmt.Println("KeyIsValid : ", tmp)
	switch tmp {
	case "f1", "f2", "f3", "f4", "f5", "f6", "f7", "f8", "f9", "f10", "f11", "f12", "print":
		return true
	default:
		return false
	}

	return false
}

func (m *BindManager) ChangeShortcut(id int32, shortcut string) ConflictInfo {
	check := conflictChecked(id, shortcut)

	tmpKeys := getShortcutById(id)
	tmpConflict := conflictChecked(id, tmpKeys)
	if check.IsConflict {
		insertConflictInvalidList(id)
		insertConflictValidList(check.IdList)

		if tmpConflict.IsConflict {
			for _, k := range tmpConflict.IdList {
				if k == id {
					continue
				}
				if !idIsExist(k, check.IdList) {
					deleteConflictValidId(k)
				}
			}
		}
	} else {
		deleteConflictInvalidId(id)
		if tmpConflict.IsConflict {
			for _, k := range tmpConflict.IdList {
				if k == id {
					continue
				}
				deleteConflictValidId(k)
			}
		}
	}

	if !KeyIsValid(shortcut) {
		insertConflictInvalidList(id)
	}
	modifyShortcutById(id, shortcut)

	return check
}

func (m *BindManager) DeleteCustomBind(id int32) {
	gs, ok := IdGSettingsMap[id]
	if !ok {
		return
	}

	tmpKeys := getShortcutById(id)
	tmpConflict := conflictChecked(id, tmpKeys)
	if tmpConflict.IsConflict {
		for _, k := range tmpConflict.IdList {
			if k == id {
				continue
			}
			deleteConflictValidId(k)
		}
	}
	deleteConflictValidId(id)
	deleteConflictInvalidId(id)

	resetCustomValues(gs)
	gs.Unref()
	delete(IdGSettingsMap, id)

	tmpList := []string{}
	idStr := strconv.FormatInt(int64(id), 10)
	customList := bindGSettings.GetStrv(_BINDING_CUSTOM_LIST)
	for _, k := range customList {
		if idStr == k {
			continue
		}
		tmpList = append(tmpList, k)
	}
	bindGSettings.SetStrv(_BINDING_CUSTOM_LIST, tmpList)
	//gio.SettingsSync()
}

func InitVariable() {
	var err error

	X, err = xgbutil.NewConn()
	if err != nil {
		fmt.Println("Unable to connect to X server:", err)
		return
	}
	keybind.Initialize(X)

	bindGSettings = gio.NewSettings(_BINDING_SCHEMA_ID)
	systemGSettings = gio.NewSettings(_SYSTEM_SCHEMA_ID)
	wmGSettings = gio.NewSettings(_WM_SCHEMA_ID)
	mediaGSettings = gio.NewSettings(_MEDIA_SCHEMA_ID)
	shiftGSettings = gio.NewSettingsWithPath(_COMPIZ_SHIFT_SCHEMA_ID,
		_COMPIZ_SHIFT_SCHEMA_PATH)
	putGSettings = gio.NewSettingsWithPath(_COMPIZ_PUT_SCHEMA_ID,
		_COMPIZ_PUT_SCHEMA_PATH)

	GrabKeyBinds = make(map[*KeyCodeInfo]string)
	IdGSettingsMap = make(map[int32]*gio.Settings)
	CustomPrevPairs = make(map[string]string)
	SystemPrevPairs = make(map[string]string)
}

func NewBindManager() *BindManager {
	m := &BindManager{}

	m.setPropList("SystemList")
	m.setPropList("MediaList")
	m.setPropList("WindowList")
	m.setPropList("WorkSpaceList")
	m.setPropList("CustomList")
	m.setPropList("ConflictValid")
	m.setPropList("ConflictInvalid")

	m.listenCustom()
	m.listenSystem()
	m.listenCompiz()
	m.listenConflict()

	return m
}

func main() {
	InitVariable()
	C.grab_xrecord_init()
	defer C.grab_xrecord_finalize()

	bm := NewBindManager()
	dbus.InstallOnSession(bm)

	gm := &GrabManager{}
	dbus.InstallOnSession(gm)

	grabKeyPairs(getSystemPairs(), true)
	grabKeyPairs(getCustomPairs(), true)
	listenKeyPressEvent()
	dbus.DealWithUnhandledMessage()

	go dlib.StartLoop()
	xevent.Main(X)
}
