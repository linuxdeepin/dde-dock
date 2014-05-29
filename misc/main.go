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

import (
	"dlib"
	"dlib/dbus"
	Logger "dlib/logger"
	libutils "dlib/utils"
	"os"
)

//import "dde-daemon/deepin-daemon/inputdevices"

var (
	logObj   = Logger.NewLogger("deepin/daemon")
	utilsObj = libutils.NewUtils()
)

type Manager struct{}

func (op *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.DeepinDaemon",
		"/com/deepin/daemon/DeepinDaemon",
		"com.deepin.daemon.DeepinDaemon",
	}
}

func NewManager() *Manager {
	m := &Manager{}

	return m
}

func main() {
	if !dlib.UniqueOnSession("com.deepin.daemon.DeepinDaemon") {
		logObj.Warning("deepin-daemon has running")
		return
	}
	defer logObj.EndTracing()

	logObj.SetRestartCommand("/usr/lib/deepin-daemon/deepin-daemon")
	//enableTouchPad()
	//listenDevices()
	//inputdevices.StartInputDevices()

	m := NewManager()
	if err := dbus.InstallOnSession(m); err != nil {
		logObj.Warning("Install Session Bus Failed: ", err)
		return
	}
	dbus.DealWithUnhandledMessage()

	startMprisDaemon()
	go dlib.StartLoop()

	if err := dbus.Wait(); err != nil {
		os.Exit(0)
	} else {
		logObj.Warning("Lost DBus")
		os.Exit(-1)
	}
}
