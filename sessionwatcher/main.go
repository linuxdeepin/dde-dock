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

package sessionwatcher

import (
	"dbus/org/freedesktop/dbus"
	"pkg.linuxdeepin.com/lib/log"
	"time"
)

var (
	Logger     = log.NewLogger("dde-daemon/sessionwatcher")
	dbusDaemon *dbus.DBusDaemon

	exitTimer chan bool
)

func startTimer() {
	for {
		timer := time.NewTimer(time.Second * 5)
		select {
		case <-timer.C:
			go GetDockApplet_T().restartDockApplet()
			go GetDDeDock_T().restartDock()
		case <-exitTimer:
			close(exitTimer)
			return
		}
	}
	Logger.Debug("Exit Timer...")
}

func Start() {
	Logger.BeginTracing()

	var err error
	if dbusDaemon, err = dbus.NewDBusDaemon("org.freedesktop.DBus", "/"); err != nil {
		Logger.Fatal("New DBusDaemon Failed:", err)
		return
	}

	exitTimer = make(chan bool)
	go startTimer()
}

func Stop() {
	Logger.EndTracing()
	if exitTimer != nil {
		exitTimer <- true
	}
	dbus.DestroyDBusDaemon(dbusDaemon)
	Logger.Debug("Exit sessionwatcher...")
}
