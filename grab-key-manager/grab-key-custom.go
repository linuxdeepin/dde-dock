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
	"dbus/com/deepin/daemon/keybinding"
	"dlib/gio-2.0"
	"fmt"
	"github.com/BurntSushi/xgbutil/keybind"
)

const (
	_CUSTOM_BUS_PATH   = "/com/deepin/daemon/KeyBinding"
	_CUSTOM_SCHEMA_ID  = "com.deepin.dde.key-binding"
	_CUSTOM_SCHEMA_KEY = "key-list"
)

var (
	_customBus       = keybinding.GetKeyBinding(_CUSTOM_BUS_PATH)
	_customGSettings = gio.NewSettings(_CUSTOM_SCHEMA_ID)
)

func BindingCustomKeys() {
	customList := _customBus.GetCustomList()

	for _, v := range customList {
		shortcut := _customBus.GetBindingAccel(v)
		key := GetXGBShortcut(shortcut)
		if GrabKeyPress(X.RootWin(), key) {
			mod, keycodes, _ := keybind.ParseString(X, key)
			action := _customBus.GetBindingExec(v)

			fmt.Println("mod:", mod)
			fmt.Println("keycodes:", keycodes)
			_KeyBindings[NewKeyInfo(mod, keycodes[1])] = action
		}
	}

	fmt.Println(_KeyBindings)
}

func ListenCustomGSettings () {
	_customGSettings.Connect("changed::key-list", 
	func (s *gio.Settings, key string) {
		
	})
}
