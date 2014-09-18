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
	libsession "dbus/com/deepin/sessionmanager"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/log"
)

const (
	DBUS_SENDER = "com.deepin.daemon.InputDevices"
)

var (
	logger = log.NewLogger("com.deepin.daemon.InputDevices")

	xsObj *libsession.XSettings
)

func Start() {
	logger.BeginTracing()

	var err error
	xsObj, err = libsession.NewXSettings("com.deepin.SessionManager",
		"/com/deepin/XSettings")
	if err != nil {
		logger.Warning("New XSettings Object Failed: ", err)
		xsObj = nil
	}

	if !initDeviceChangedWatcher() {
		logger.Fatal("Init device changed wacher failed")
		return
	}

	if err := dbus.InstallOnSession(GetManager()); err != nil {
		logger.Fatal("Install Manager DBus Failed:", err)
	}

	if err := dbus.InstallOnSession(GetMouseManager()); err != nil {
		logger.Fatal("Install Mouse DBus Failed:", err)
	}

	if err := dbus.InstallOnSession(GetTouchpadManager()); err != nil {
		logger.Fatal("Install Touchpad DBus Failed:", err)
	}

	if err := dbus.InstallOnSession(GetKeyboardManager()); err != nil {
		logger.Fatal("Install Keyboard DBus Failed:", err)
	}

	if err := dbus.InstallOnSession(GetWacomManager()); err != nil {
		logger.Fatal("Install Wacom DBus Failed:", err)
	}
}

func Stop() {
	logger.EndTracing()

	endDeviceListenThread()

	GetTouchpadManager().typingExitChan <- true

	dbus.UnInstallObject(GetKeyboardManager())
	dbus.UnInstallObject(GetTouchpadManager())
	dbus.UnInstallObject(GetMouseManager())
	dbus.UnInstallObject(GetWacomManager())
	dbus.UnInstallObject(GetManager())
}
