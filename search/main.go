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

package search

import (
	"os"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/logger"
	"time"
)

type Manager struct {
	writeStart bool
	writeEnd   chan bool
}

const (
	DBUS_DEST = "com.deepin.daemon.Search"
	DBUS_PATH = "/com/deepin/daemon/Search"
	DBUS_IFC  = "com.deepin.daemon.Search"
)

var (
	Logger = logger.NewLogger(DBUS_DEST)
)

func newManager() *Manager {
	m := Manager{}

	m.writeStart = false

	return &m
}

var _manager *Manager

func GetManager() *Manager {
	if _manager == nil {
		_manager = newManager()
	}

	return _manager
}

func Start() {
	Logger.BeginTracing()
	Logger.SetRestartCommand("/usr/lib/deepin-daemon/dde-session-daemon")

	if err := dbus.InstallOnSession(GetManager()); err != nil {
		Logger.Fatal("SearchReg Install DBus Failed:", err)
		return
	}

	dbus.SetAutoDestroyHandler(time.Second*5, func() bool {
		if GetManager().writeStart {
			select {
			case <-GetManager().writeEnd:
				return true
			}
		}

		return true
	})

	if err := dbus.Wait(); err != nil {
		Logger.Warning("SearchReg lost dbus:", err)
		os.Exit(-1)
	} else {
		os.Exit(0)
	}
}

func Stop() {
	Logger.EndTracing()
}
