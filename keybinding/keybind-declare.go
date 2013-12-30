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

type Manager struct {
	CustomBindList []int32
	customAccelMap map[string]string
	gsdAccelMap    map[string]string
	/*wmAccelMap     map[string]string*/
	/*mediaAccelMap  map[string]string*/
	/*shiftAccelMap  map[string]string*/
	/*putAccelMap    map[string]string*/
}

type GrabManager struct {
	GrabKeyEvent func(string)
}

type GrabKeyInfo struct {
	State  uint16
	Detail xproto.Keycode
}

type GrabExecInfo struct {
	Id   int32
	Exec string
}

const (
	_KEY_BINDING_NAME = "com.deepin.daemon.KeyBinding"
	_KEY_BINDING_PATH = "/com/deepin/daemon/KeyBinding"
	_KEY_BINDING_IFC  = "com.deepin.daemon.KeyBinding"
	_GRAB_KEY_PATH    = "/com/deepin/daemon/GrabManager"
	_GRAB_KEY_IFC     = "com.deepin.daemon.GrabManager"

	_CUSTOM_SCHEMA_ID       = "com.deepin.dde.key-binding"
	_CUSTOM_SCHEMA_ADD_ID   = "com.deepin.dde.key-binding.custom"
	_CUSTOM_SCHEMA_ADD_PATH = "/com/deepin/dde/key-binding/profiles/"

	_WM_SCHEMA_ID             = "org.gnome.desktop.wm.keybindings"
	_GSD_SCHEMA_ID            = "org.gnome.settings-daemon.plugins.key-bindings"
	_MEDIA_SCHEMA_ID          = "org.gnome.settings-daemon.plugins.media-keys"
	_COMPIZ_SHIFT_SCHEMA_ID   = "org.compiz.shift"
	_COMPIZ_SHIFT_SCHEMA_PATH = "/org/compiz/profiles/shift/"
	_COMPIZ_PUT_SCHEMA_ID     = "org.compiz.put"
	_COMPIZ_PUT_SCHEMA_PATH   = "/org/compiz/profiles/put/"

	_CUSTOM_KEY_BASE     = 10000
	_CUSTOM_KEY_LIST     = "key-list"
	_CUSTOM_KEY_ID       = "id"
	_CUSTOM_KEY_NAME     = "name"
	_CUSTOM_KEY_SHORTCUT = "shortcut"
	_CUSTOM_KEY_ACTION   = "action"
)

var (
	customGSettings = gio.NewSettings(_CUSTOM_SCHEMA_ID)
	gsdGSettings    = gio.NewSettings(_GSD_SCHEMA_ID)
	mediaGSettings  = gio.NewSettings(_MEDIA_SCHEMA_ID)
	wmGSettings     = gio.NewSettings(_WM_SCHEMA_ID)
	shiftGSettings  = gio.NewSettingsWithPath(_COMPIZ_SHIFT_SCHEMA_ID,
		_COMPIZ_SHIFT_SCHEMA_PATH)
	putGSettings = gio.NewSettingsWithPath(_COMPIZ_PUT_SCHEMA_ID,
		_COMPIZ_PUT_SCHEMA_PATH)

	X            *xgbutil.XUtil
	GrabKeyBinds map[*GrabKeyInfo]string
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

/*
 * 0 ~ 299: org.gnome-settings-daemon.plugins.key-bindings
 * 300 ~ 599: org.gnome-settings-daemon.plugins.media-keys
 * 600 ~ 799: org.gnome.desktop.wm.keybindings
 * 800 ~ 899: org.compiz.shift; path: /org/compiz/profiles/shift/
 * 900 ~ 999: org.compiz.put; path: /org/compiz/profiles/put/
 */
var currentSystemBindings = map[int32]string{
	0:   "key1",
	1:   "key2",
	2:   "key3",
	3:   "key4",
	4:   "key5",
	5:   "key6",
	6:   "key7",
	7:   "key8",
	8:   "key9",
	9:   "key10",
	10:  "key11",
	301: "calculator",
	302: "eject",
	303: "email",
	304: "www",
	305: "media",
	306: "play",
	307: "pause",
	308: "stop",
	309: "volume-down",
	310: "volume-mute",
	311: "volume-up",
	312: "previous",
	313: "next",
	600: "activate-window-menu",
	601: "begin-move",
	602: "begin-resize",
	603: "close",
	604: "maximize",
	605: "minimize",
	606: "toggle-shaded",
	607: "unmaximize",
	608: "switch-to-workspace-1",
	609: "switch-to-workspace-2",
	610: "switch-to-workspace-3",
	611: "switch-to-workspace-4",
	612: "switch-to-workspace-down",
	613: "switch-to-workspace-left",
	614: "switch-to-workspace-right",
	615: "switch-to-workspace-up",
	616: "move-to-workspace-down",
	617: "move-to-workspace-left",
	618: "move-to-workspace-right",
	619: "move-to-workspace-up",
	620: "show-desktop",
	621: "switch-windows",
	622: "switch-windows-backward",
	800: "prev-key",           //switch apps with 3D
	801: "next-key",           //reverse switch apps with 3D
	900: "put-viewport-1-key", //Move window to workspace 1
	901: "put-viewport-2-key",
	902: "put-viewport-3-key",
	903: "put-viewport-4-key",
}

var gsdMap = map[int32]string{
	0:  "key1",
	1:  "key2",
	2:  "key3",
	3:  "key4",
	4:  "key5",
	5:  "key6",
	6:  "key7",
	7:  "key8",
	8:  "key9",
	9:  "key10",
	10: "key11",
}

var mediaMap = map[int32]string{
	301: "calculator",
	302: "eject",
	303: "email",
	304: "www",
	305: "media",
	306: "play",
	307: "pause",
	308: "stop",
	309: "volume-down",
	310: "volume-mute",
	311: "volume-up",
	312: "previous",
	313: "next",
}

var wmMap = map[int32]string{
	600: "activate-window-menu", //window
	601: "begin-move",
	602: "begin-resize",
	603: "close",
	604: "maximize",
	605: "minimize",
	606: "toggle-shaded",
	607: "unmaximize",            //window
	608: "switch-to-workspace-1", //workspace
	609: "switch-to-workspace-2",
	610: "switch-to-workspace-3",
	611: "switch-to-workspace-4",
	612: "switch-to-workspace-down",
	613: "switch-to-workspace-left",
	614: "switch-to-workspace-right",
	615: "switch-to-workspace-up",
	616: "move-to-workspace-down",
	617: "move-to-workspace-left",
	618: "move-to-workspace-right",
	619: "move-to-workspace-up",    //workspace
	620: "show-desktop",
	621: "switch-windows",          // gsd
	622: "switch-windows-backward", // gsd
}

//gsd
var shiftMap = map[int32]string{
	800: "prev-key", //switch apps with 3D
	801: "next-key", //reverse switch apps with 3D
}

//workspace
var putMap = map[int32]string{
	900: "put-viewport-1-key", //Move window to workspace 1
	901: "put-viewport-2-key",
	902: "put-viewport-3-key",
	903: "put-viewport-4-key",
}
