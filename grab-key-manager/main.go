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

type Manager struct{}

const (
	_GRAB_KEY_DEST = "com.deepin.daemon.GrabKeyManager"
	_GRAB_KEY_PATH = "/com/deepin/daemon/GrabKeyManager"
	_GRAB_KEY_IFC  = "com.deepin.daemon.GrabKeyManager"
)

var (
	X            *xgbutil.XUtil
	_KeyBindings map[*_KeyInfo]string
)

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_GRAB_KEY_DEST,
		_GRAB_KEY_PATH,
		_GRAB_KEY_IFC,
	}
}

func LookupString(mod uint16, keycode xproto.Keycode) string {
	return keybind.LookupString(X, mod, keycode)
}

func InitGrabKey() {
	var err error

	X, err = xgbutil.NewConn()
	if err != nil {
		fmt.Println("Get New Connection Failed:", err)
		return
	}
	keybind.Initialize(X)

	_KeyBindings = make(map[*_KeyInfo]string)
}

func main() {
	InitGrabKey()
	BindingCustomKeys()
	xevent.Main(X)
}
