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
	"pkg.deepin.io/dde-daemon/inputdevices/keyboard"
	"pkg.deepin.io/dde-daemon/inputdevices/mouse"
	"pkg.deepin.io/dde-daemon/inputdevices/touchpad"
	"pkg.deepin.io/dde-daemon/inputdevices/wacom"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

// TODO: delete this struct && rm this file
type Manager struct {
	Infos []devicePathInfo

	mouse        *mouse.Mouse
	touchpad     *touchpad.Touchpad
	kbd          *keyboard.Keyboard
	wacom        *wacom.Wacom
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
	// Touchpad must be created Before Mouse
	m.touchpad = touchpad.NewTouchpad(l)
	m.mouse = mouse.NewMouse(l)
	m.kbd = keyboard.NewKeyboard(l)
	m.wacom = wacom.NewWacom(l)
	m.versionRight = m.isVersionRight()

	return m
}

func (m *Manager) destroy() {
	dbus.UnInstallObject(_m.mouse)
	dbus.UnInstallObject(_m.touchpad)
	dbus.UnInstallObject(_m.kbd)
	dbus.UnInstallObject(_m.wacom)
	dbus.UnInstallObject(_m)

	if _m.logger != nil {
		_m.logger.EndTracing()
	}
}

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusSender,
		ObjectPath: "/com/deepin/daemon/InputDevices",
		Interface:  dbusSender,
	}
}
