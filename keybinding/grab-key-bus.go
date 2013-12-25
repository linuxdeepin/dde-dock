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
	"dlib/dbus"
	"fmt"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
)

func (m *GrabManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_KEY_BINDING_NAME,
		_GRAB_KEY_PATH,
		_GRAB_KEY_IFC,
	}
}

func (m *GrabManager) GrabShortcut(wid xproto.Window,
	shortcut, action string) bool {
	if wid == 0 {
		wid = X.RootWin()
	}

	key := GetXGBShortcut(shortcut)
	mod, keycodes, _ := keybind.ParseString(X, key)
	if len(keycodes) <= 0 {
		return false
	}
	if !GrabKeyPress(wid, shortcut) {
		return false
	}
	keyInfo := NewKeyInfo(mod, keycodes[0])
	GrabKeyBinds[keyInfo] = action

	return true
}

func (m *GrabManager) UngrabShortcut(wid xproto.Window,
	shortcut string) bool {
	if wid == 0 {
		wid = X.RootWin()
	}

	return UngrabKey(wid, shortcut)
}

func (m *GrabManager) GrabKeyboard() {
	X, err := xgbutil.NewConn()
	if err != nil {
		fmt.Println("Get New Connection Failed:", err)
		return
	}
	keybind.Initialize(X)

	err = keybind.GrabKeyboard(X, X.RootWin())
	if err != nil {
		fmt.Println("Grab Keyboard Failed:", err)
		return
	}

	xevent.KeyReleaseFun(
		func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
			modStr := keybind.ModifierString(e.State)
			keyStr := keybind.LookupString(X, e.State, e.Detail)
			value := ""
			if len(modStr) > 0 {
				value = modStr + "-" + keyStr
			} else {
				value = keyStr
			}
			m.GrabKeyEvent(value)
			keybind.UngrabKeyboard(X)
			fmt.Printf("Key: %s\n", value)
			xevent.Quit(X)
		}).Connect(X, X.RootWin())

	xevent.Main(X)
}
