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
	"os"
	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
	"time"
)

const (
	DEST = "com.deepin.daemon.Personalization"
)

var (
	Logger = log.NewLogger("dde-daemon/personalization")
)

func main() {
	if !lib.UniqueOnSession(DEST) {
		Logger.Warning("There has an set-fonts running...")
		return
	}

	Logger.BeginTracing()
	defer Logger.EndTracing()
	Logger.SetRestartCommand("/usr/lib/deepin-daemon/personalization")

	StartFont()
	dbus.DealWithUnhandledMessage()

	dbus.SetAutoDestroyHandler(time.Second*5, nil)
	if err := dbus.Wait(); err != nil {
		Logger.Warning("Lost DBus:", err)
		os.Exit(-1)
	} else {
		os.Exit(0)
	}
}
