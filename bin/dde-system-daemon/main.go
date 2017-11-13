/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"os"

	// modules:
	_ "pkg.deepin.io/dde/daemon/accounts"
	_ "pkg.deepin.io/dde/daemon/apps"
	_ "pkg.deepin.io/dde/daemon/system/gesture"
	_ "pkg.deepin.io/dde/daemon/system/power"
	_ "pkg.deepin.io/dde/daemon/system/swapsched"
	_ "pkg.deepin.io/dde/daemon/system/timedated"

	"gir/glib-2.0"
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/log"
)

type Daemon struct{}

const (
	dbusDest = "com.deepin.daemon.Daemon"
	dbusPath = "/com/deepin/daemon/Daemon"
	dbusIFC  = dbusDest
)

var logger = log.NewLogger("daemon/dde-system-daemon")
var _daemon *Daemon

func main() {
	logger.BeginTracing()
	defer logger.EndTracing()

	if !lib.UniqueOnSystem(dbusDest) {
		logger.Warning("There already has an dde daemon running.")
		return
	}

	// fix no PATH when was launched by dbus
	if os.Getenv("PATH") == "" {
		logger.Warning("No PATH found, manual special")
		os.Setenv("PATH", "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")
	}

	InitI18n()
	Textdomain("dde-daemon")

	logger.SetRestartCommand("/usr/lib/deepin-daemon/dde-system-daemon")

	_daemon = &Daemon{}
	err := dbus.InstallOnSystem(_daemon)
	if err != nil {
		logger.Error("Failed to install daemon bus:", err)
		return
	}
	dbus.DealWithUnhandledMessage()

	loader.StartAll()
	defer loader.StopAll()

	// NOTE: system/power module requires glib loop
	go glib.StartLoop()

	if err := dbus.Wait(); err != nil {
		logger.Errorf("Lost dbus: %v", err)
		os.Exit(-1)
	} else {
		os.Exit(0)
	}
}

func (*Daemon) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}
