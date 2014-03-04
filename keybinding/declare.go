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
        "github.com/BurntSushi/xgb/xproto"
        "github.com/BurntSushi/xgbutil"
)

type ShortcutInfo struct {
        Id       int32
        Desc     string
        Shortcut string
        index    int32
}

type BindManager struct {
        SystemList    []ShortcutInfo
        MediaList     []ShortcutInfo
        WindowList    []ShortcutInfo
        WorkSpaceList []ShortcutInfo
        CustomList    []ShortcutInfo

        ConflictValid   []int32
        ConflictInvalid []int32
}

type KeyCodeInfo struct {
        State   uint16
        Detail  xproto.Keycode
}

type ConflictInfo struct {
        IsConflict bool
        IdList     []int32
}

const (
        _BINDING_DEST = "com.deepin.daemon.KeyBinding"
        _BINDING_PATH = "/com/deepin/daemon/KeyBinding"
        _BINDING_IFC  = "com.deepin.daemon.KeyBinding"

        _BINDING_SCHEMA_ID      = "com.deepin.dde.key-binding"
        _SYSTEM_SCHEMA_ID       = "com.deepin.dde.key-binding.system"
        _CUSTOM_ADD_SCHEMA_ID   = "com.deepin.dde.key-binding.custom"
        _CUSTOM_ADD_SCHEMA_PATH = "/com/deepin/dde/key-binding/profiles/"

        _WM_SCHEMA_ID             = "org.gnome.desktop.wm.keybindings"
        _COMPIZ_SHIFT_SCHEMA_ID   = "org.compiz.shift"
        _COMPIZ_SHIFT_SCHEMA_PATH = "/org/compiz/profiles/shift/"
        _COMPIZ_PUT_SCHEMA_ID     = "org.compiz.put"
        _COMPIZ_PUT_SCHEMA_PATH   = "/org/compiz/profiles/put/"

        _CUSTOM_ID_BASE      = 10000
        _CUSTOM_KEY_ID       = "id"
        _CUSTOM_KEY_NAME     = "name"
        _CUSTOM_KEY_SHORTCUT = "shortcut"
        _CUSTOM_KEY_ACTION   = "action"

        _BINDING_CUSTOM_LIST  = "custom-list"
        _BINDING_VALID_LIST   = "conflict-valid"
        _BINDING_INVALID_LIST = "conflict-invalid"
)

var (
        systemGSettings *gio.Settings
        bindGSettings   *gio.Settings
        wmGSettings     *gio.Settings
        shiftGSettings  *gio.Settings
        putGSettings    *gio.Settings

        X              *xgbutil.XUtil
        GrabKeyBinds   map[*KeyCodeInfo]string
        IdGSettingsMap map[int32]*gio.Settings

        CustomPrevPairs map[string]string
        SystemPrevPairs map[string]string
)

var _ModifierMap = map[string]string{
        "caps_lock": "lock",
        "alt":       "mod1",
        "meta":      "mod1",
        "num_lock":  "mod2",
        "super":     "mod4",
        "hyper":     "mod4",
}

var _ModKeyMap = map[string]string{
        "mod1": "alt",
        "mod4": "super",
        "lock": "caps_lock",
        "mod2": "num_lock",
}
