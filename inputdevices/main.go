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
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/log"
)

const (
	dbusSender = "com.deepin.daemon.InputDevices"
)

var _m *Manager

func Start() {
	var logger = log.NewLogger("com.deepin.daemon.InputDevices")
	logger.BeginTracing()

	if !initDeviceChangedWatcher() {
		logger.Fatal("Init device changed wacher failed")
		return
	}

	_m := NewManager(logger)
	err := dbus.InstallOnSession(_m)
	if err != nil {
		logger.Fatal("Install Manager DBus Failed:", err)
	}

	err = dbus.InstallOnSession(_m.mouse)
	if err != nil {
		logger.Fatal("Install Mouse DBus Failed:", err)
	}

	err = dbus.InstallOnSession(_m.touchpad)
	if err != nil {
		logger.Fatal("Install Touchpad DBus Failed:", err)
	}

	err = dbus.InstallOnSession(_m.kbd)
	if err != nil {
		logger.Fatal("Install Keyboard DBus Failed:", err)
	}

	err = dbus.InstallOnSession(_m.wacom)
	if err != nil {
		logger.Fatal("Install Wacom DBus Failed:", err)
	}
}

func Stop() {
	if _m == nil {
		return
	}

	if _m.logger != nil {
		_m.logger.EndTracing()
	}

	endDeviceListenThread()

	dbus.UnInstallObject(_m.mouse)
	dbus.UnInstallObject(_m.touchpad)
	dbus.UnInstallObject(_m.kbd)
	dbus.UnInstallObject(_m.wacom)
	dbus.UnInstallObject(_m)
	_m = nil
}
