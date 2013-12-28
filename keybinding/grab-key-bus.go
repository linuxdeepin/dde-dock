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
// #include <stdlib.h>
import "C"

import (
	"dlib/dbus"
	"fmt"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"strings"
	"unsafe"
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
	mousebind.Initialize(X)

	err = keybind.GrabKeyboard(X, X.RootWin())
	if err != nil {
		fmt.Println("Grab Keyboard Failed:", err)
		return
	}

	GrabAllButton(X)

	xevent.ButtonPressFun(
		func(X *xgbutil.XUtil, e xevent.ButtonPressEvent) {
			m.GrabKeyEvent("")
			UngrabAllButton(X)
			keybind.UngrabKeyboard(X)
			fmt.Println("Button Press Event")
			xevent.Quit(X)
		}).Connect(X, X.RootWin())

	xevent.KeyReleaseFun(
		func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
			modStr := keybind.ModifierString(e.State)
			keyStr := keybind.LookupString(X, e.State, e.Detail)
			value := ""
			if len(modStr) > 0 {
				value = ConvertKeyFromMod(modStr) + keyStr
			} else {
				value = keyStr
			}
			m.GrabKeyEvent(value)
			UngrabAllButton(X)
			keybind.UngrabKeyboard(X)
			fmt.Printf("Key: %s\n", value)
			xevent.Quit(X)
		}).Connect(X, X.RootWin())

	xevent.Main(X)
}

func GrabAllButton (X *xgbutil.XUtil) {
	mousebind.Grab(X, X.RootWin(), 0, 1, false)
	mousebind.Grab(X, X.RootWin(), 0, 2, false)
	mousebind.Grab(X, X.RootWin(), 0, 3, false)
}

func UngrabAllButton (X *xgbutil.XUtil) {
			mousebind.Ungrab(X, X.RootWin(), 0, 1)
			mousebind.Ungrab(X, X.RootWin(), 0, 2)
			mousebind.Ungrab(X, X.RootWin(), 0, 3)
}

func ConvertKeyFromMod(mod string) string {
	values := ""
	vals := strings.Split(mod, "-")
	for _, v := range vals {
		if v == "mod1" || v == "mod2" ||
			v == "mod4" || v == "lock" {
			t, ok := _ModKeyMap[v]
			if !ok {
				fmt.Println("Get Key Failed From Modify")
				return ""
			}
			values += t + "-"
		} else {
			values += v + "-"
		}
	}

	return values
}

func (m *GrabManager) GrabSingleKey(key, action string) {
	GrabXRecordKey(key, action)
}

func GrabXRecordKey(key, action string) {
	if len(action) <= 0 {
		fmt.Println("action is null")
		return
	}

	mod, keys, err := keybind.ParseString(X, key)
	if err != nil {
		fmt.Println("ParseString Failed:", err)
		return
	}

	fmt.Printf("mod: %d, key: %d\n", mod, keys[0])
	if mod > 0 {
		fmt.Printf("Not single key\n")
		return
	}

	tmp := C.CString(action)
	defer C.free(unsafe.Pointer(tmp))
	C.grab_xrecord_key(C.int(keys[0]), tmp)
}

func (m *GrabManager) UngrabSingleKey(key string) {
	UngrabXRecordKey(key)
}

func UngrabXRecordKey(key string) {
	mod, keys, err := keybind.ParseString(X, key)
	if err != nil {
		fmt.Println("ParseString Failed:", err)
		return
	}

	if mod > 0 {
		fmt.Printf("Not single key\n")
		return
	}

	C.ungrab_xrecord_key(C.int(keys[0]))
}

func (m GrabManager) GrabSingleFinalize() {
	C.grab_xrecord_finalize()
}
