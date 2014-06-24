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
	"github.com/BurntSushi/xgb/xproto"
)

const (
	CUSTOM_KEY_ID        = "id"
	CUSTOM_KEY_NAME      = "name"
	CUSTOM_KEY_ACTION    = "action"
	CUSTOM_KEY_SHORTCUT  = "shortcut"
	CUSTOM_KEY_ID_BASE   = 10000
	CUSTOM_KEY_SCHEMA_ID = "com.deepin.dde.keybinding.custom"
	CUSTOM_KEY_BASE_PATH = "/com/deepin/dde/keybinding/profiles/"

	BIND_KEY_VALID_LIST   = "conflict-valid"
	BIND_KEY_INVALID_LIST = "conflict-invalid"
	BIND_KEY_CUSTOM_LIST  = "custom-list"

	KEYBIND_DEST  = "com.deepin.daemon.KeyBinding"
	MANAGER_PATH  = "/com/deepin/daemon/KeyBinding"
	MANAGER_IFC   = "com.deepin.daemon.KeyBinding"
	MEDIAKEY_PATH = "/com/deepin/daemon/MediaKey"
	MEDIAKEY_IFC  = "com.deepin.daemon.MediaKey"
)

var (
	keyToModMap = map[string]string{
		"caps_lock": "lock",
		"alt":       "mod1",
		"meta":      "mod1",
		"num_lock":  "mod2",
		"super":     "mod4",
		"hyper":     "mod4",
	}

	modToKeyMap = map[string]string{
		"mod1": "alt",
		"mod2": "num_lock",
		"mod4": "super",
		"lock": "caps_lock",
	}
)

type ShortcutInfo struct {
	Id       int32
	Desc     string
	Shortcut string
	index    int32
}

type KeycodeInfo struct {
	State  uint16
	Detail xproto.Keycode
}

type ConflictInfo struct {
	IsConflict bool
	IdList     []int32
}

type Manager struct {
	SystemList []ShortcutInfo
	//MediaList     []ShortcutInfo
	WindowList    []ShortcutInfo
	WorkspaceList []ShortcutInfo
	CustomList    []ShortcutInfo

	ConflictValid   []int32
	ConflictInvalid []int32

	idSettingsMap map[int32]*gio.Settings

	KeyReleaseEvent func(string)
}

type MediaKeyManager struct {
	AudioMute        func(bool)
	AudioUp          func(bool)
	AudioDown        func(bool)
	BrightnessUp     func(bool)
	BrightnessDown   func(bool)
	CapsLockOn       func(bool)
	CapsLockOff      func(bool)
	NumLockOn        func(bool)
	NumLockOff       func(bool)
	SwitchMonitors   func(bool)
	TouchPadOn       func(bool)
	TouchPadOff      func(bool)
	PowerOff         func(bool)
	PowerSleep       func(bool)
	SwitchLayout     func(bool)
	AudioPlay        func(bool)
	AudioPause       func(bool)
	AudioStop        func(bool)
	AudioPrevious    func(bool)
	AudioNext        func(bool)
	AudioRewind      func(bool)
	AudioForward     func(bool)
	AudioRepeat      func(bool)
	LaunchEmail      func(bool)
	LaunchBrowser    func(bool)
	LaunchCalculator func(bool)
}
