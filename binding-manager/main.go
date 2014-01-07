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
)

type AddAccelRet struct {
	Id    int32
	Check *ConflictInfo
}

func (m *BindManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_BINDING_DEST,
		_BINDING_PATH,
		_BINDING_IFC,
	}
}

func (m *BindManager) AddKeyBind(name, action, shortcut string) *AddAccelRet {
	id := GetMaxIdFromCustom() + 1
	gs := NewGSettingsById(id)
	if gs == nil {
		return nil
	}
	IdGSettingsMap[id] = gs

	SetCustomValues(gs, id, name, action, "")
	gs.Connect("changed::shortcut", func(s *gio.Settings, key string) {
		fmt.Printf("key: %s, value: %s\n", key, gs.GetString(key))
		m.CustomList = GetCustomKeyInfo()
		dbus.NotifyChange(m, "CustomList")
		GrabKeyPairs(CustomPrevPairs, false)
		GrabKeyPairs(GetCustomPairs(), true)
	})
	ret := &AddAccelRet{}
	ret.Id = id
	ret.Check = m.ChangeShortcut(id, shortcut)

	idStr := strconv.FormatInt(int64(id), 10)
	customList := bindGSettings.GetStrv(_BINDING_CUSTOM_LIST)
	customList = append(customList, idStr)
	bindGSettings.SetStrv(_BINDING_CUSTOM_LIST, customList)
	gio.SettingsSync()

	GrabKeyPairs(CustomPrevPairs, false)
	GrabKeyPairs(GetCustomPairs(), true)
	return ret
}

func (m *BindManager) ChangeShortcut(id int32, shortcut string) *ConflictInfo {
	check := ConflictChecked(id, shortcut)
	if check == nil {
		return nil
	}

	tmpKeys := GetShortcutById(id)
	tmpConflict := ConflictChecked(id, tmpKeys)
	if check.IsConflict {
		InsertConflictInvalidList(id)
		InsertConflictValidList(check.IdList)

		if tmpConflict != nil && tmpConflict.IsConflict {
			for _, k := range tmpConflict.IdList {
				if k == id {
					continue
				}
				if !IdIsExist(k, check.IdList) {
					DeleteConflictValidId(k)
				}
			}
		}
	} else {
		DeleteConflictInvalidId(id)
		if tmpConflict != nil && tmpConflict.IsConflict {
			for _, k := range tmpConflict.IdList {
				if k == id {
					continue
				}
				DeleteConflictValidId(k)
			}
		}
	}
	ModifyShortcutById(id, shortcut)

	return check
}

func (m *BindManager) DeleteCustomBind(id int32) {
	gs, ok := IdGSettingsMap[id]
	if !ok {
		return
	}

	tmpKeys := GetShortcutById(id)
	tmpConflict := ConflictChecked(id, tmpKeys)
	if tmpConflict != nil && tmpConflict.IsConflict {
		for _, k := range tmpConflict.IdList {
			if k == id {
				continue
			}
			DeleteConflictValidId(k)
		}
	}
	DeleteConflictValidId(id)
	DeleteConflictInvalidId(id)

	ResetCustomValues(gs)
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
	gio.SettingsSync()
}

func InitConflictList(m *BindManager) {
	validList := bindGSettings.GetStrv(_BINDING_VALID_LIST)
	invalidList := bindGSettings.GetStrv(_BINDING_INVALID_LIST)

	for _, k := range validList {
		tmp, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			continue
		}
		m.ConflictValid = append(m.ConflictValid, int32(tmp))
	}

	for _, k := range invalidList {
		tmp, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			continue
		}
		m.ConflictInvalid = append(m.ConflictInvalid, int32(tmp))
	}
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

func InitListen(m *BindManager) {
	ListenCustom(m)
	ListenSystem(m)
	ListenCompiz(m)

	bindGSettings.Connect("changed::conflict-valid", func(s *gio.Settings, key string) {
		validList := s.GetStrv(_BINDING_VALID_LIST)
		fmt.Println("chaned valid: ", validList)
		tmpList := []int32{}
		for _, k := range validList {
			tmp, err := strconv.ParseInt(k, 10, 64)
			if err != nil {
				continue
			}
			tmpList = append(tmpList, int32(tmp))
		}
		m.ConflictValid = tmpList
		dbus.NotifyChange(m, "ConflictValid")
	})

	bindGSettings.Connect("changed::conflict-invalid", func(s *gio.Settings, key string) {
		invalidList := s.GetStrv(_BINDING_INVALID_LIST)
		fmt.Println("changed invalid: ", invalidList)
		tmpList := []int32{}
		for _, k := range invalidList {
			tmp, err := strconv.ParseInt(k, 10, 64)
			if err != nil {
				continue
			}
			tmpList = append(tmpList, int32(tmp))
		}
		m.ConflictInvalid = tmpList
		dbus.NotifyChange(m, "ConflictInvalid")
	})
}

func NewBindManager() *BindManager {
	m := &BindManager{}

	InitSystemBind(m)
	InitMediaBind(m)
	InitWindowBind(m)
	InitWorkSpaceBind(m)
	InitCustomBind(m)
	InitConflictList(m)

	return m
}

func main() {
	InitVariable()
	C.grab_xrecord_init()
	defer C.grab_xrecord_finalize()

	bm := NewBindManager()
	InitListen(bm)
	dbus.InstallOnSession(bm)

	gm := &GrabManager{}
	dbus.InstallOnSession(gm)

	GrabKeyPairs(GetSystemPairs(), true)
	GrabKeyPairs(GetCustomPairs(), true)
	ListenKeyPressEvent()
	dbus.DealWithUnhandledMessage()

        go dlib.StartLoop()
	xevent.Main(X)
}
