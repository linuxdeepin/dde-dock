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

package main

// #cgo pkg-config: gtk+-3.0 x11 glib-2.0
// #cgo CFLAGS: -Wall -g
// #include "gsd-clipboard-manager.h"
import "C"

import (
	"dlib"
	"dlib/dbus"
	"dlib/logger"
	"os"
)

type Manager struct{}

var (
	logObj = logger.NewLogger("daemon/clipboard")
)

func (op *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Clipboard",
		"/com/deepin/daemon/Clipboard",
		"com.deepin.daemon.Clipboard",
	}
}

func (op *Manager) StopClipboard() {
	C.stop_clip_manager()
}

func main() {
	if !dlib.UniqueOnSession("com.deepin.daemon.Clipboard") {
		logObj.Warning("Clipboard has running...")
		return
	}

	defer logObj.EndTracing()
	logObj.SetRestartCommand("/usr/lib/deepin-daemon/clipboard")

	C.start_clip_manager()

	m := &Manager{}
	if err := dbus.InstallOnSession(m); err != nil {
		logObj.Fatal("Clipboard Install DBus Failed:", err)
		return
	}
	dbus.DealWithUnhandledMessage()

	if err := dbus.Wait(); err != nil {
		logObj.Warning("Clipboard Lost DBus")
		os.Exit(-1)
	} else {
		os.Exit(0)
	}
}
