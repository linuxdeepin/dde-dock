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

package inputdevices

import (
	"pkg.linuxdeepin.com/dde-daemon/inputdevices/libkeyboard"
	"pkg.linuxdeepin.com/dde-daemon/inputdevices/libmouse"
	"pkg.linuxdeepin.com/dde-daemon/inputdevices/libtouchpad"
	"pkg.linuxdeepin.com/dde-daemon/inputdevices/libwacom"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/log"
)

// TODO: delete this struct && rm this file
type Manager struct {
	Infos []devicePathInfo

	mouse        *libmouse.Mouse
	touchpad     *libtouchpad.Touchpad
	kbd          *libkeyboard.Keyboard
	wacom        *libwacom.Wacom
	logger       *log.Logger
	versionRight bool
}

type devicePathInfo struct {
	Path string
	Type string
}

func NewManager(l *log.Logger) *Manager {
	m := &Manager{}

	m.Infos = []devicePathInfo{
		devicePathInfo{"com.deepin.daemon.InputDevice.Keyboard",
			"keyboard"},
		devicePathInfo{"com.deepin.daemon.InputDevice.Mouse",
			"mouse"},
		devicePathInfo{"com.deepin.daemon.InputDevice.TouchPad",
			"touchpad"},
	}

	m.logger = l
	m.mouse = libmouse.NewMouse(l)
	m.touchpad = libtouchpad.NewTouchpad(l)
	m.kbd = libkeyboard.NewKeyboard(l)
	m.wacom = libwacom.NewWacom(l)
	m.versionRight = m.isVersionRight()

	return m
}

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusSender,
		ObjectPath: "/com/deepin/daemon/InputDevices",
		Interface:  dbusSender,
	}
}
