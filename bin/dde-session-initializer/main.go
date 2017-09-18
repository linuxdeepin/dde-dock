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

//#cgo pkg-config: x11 gtk+-3.0
//#include <X11/Xlib.h>
//#include <gtk/gtk.h>
//void init(){XInitThreads();gtk_init(0,0);}
import "C"
import (
	"gir/glib-2.0"
	"pkg.deepin.io/dde/api/session"
	_ "pkg.deepin.io/dde/daemon/dock"
	_ "pkg.deepin.io/dde/daemon/launcher"
	"pkg.deepin.io/dde/daemon/loader"
	_ "pkg.deepin.io/dde/daemon/trayicon"
	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/app"
	"pkg.deepin.io/lib/dbus"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/proxy"
	"pkg.deepin.io/lib/utils"

	"os"
	"time"
)

var logger = log.NewLogger("daemon/initializer")

func runMainLoop() {
	logger.Info("register session")
	startTime := time.Now()
	session.Register()
	logger.Info("register session done, cost", time.Now().Sub(startTime))

	logger.Info("DealWithUnhandledMessage")
	startTime = time.Now()
	dbus.DealWithUnhandledMessage()
	logger.Info("DealWithUnhandledMessage done, cost", time.Now().Sub(startTime))
	go glib.StartLoop()

	logger.Info("initialize done")
	if err := dbus.Wait(); err != nil {
		logger.Errorf("Lost dbus: %v", err)
		os.Exit(-1)
	}

	logger.Info("dbus connection is closed by user")
	os.Exit(0)
}

type Initializer struct {
}

func (*Initializer) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.daemon.Initializer",
		ObjectPath: "/com/deepin/daemon/Initializer",
		Interface:  "com.deepin.daemon.Initializer",
	}
}

func main() {
	sessionInitializer := new(Initializer)
	if !lib.UniqueOnSession(sessionInitializer.GetDBusInfo().Dest) {
		logger.Warning("There's a dde-session-initializer instance running.")
		os.Exit(0)
	}

	err := dbus.InstallOnSession(sessionInitializer)
	if err != nil {
		logger.Fatal(err)
	}

	cmd := app.New("dde-session-initializer",
		"daemon for dde-dock dde-launcher", "version "+__VERSION__)
	cmd.ParseCommandLine(os.Args[1:])
	if err := cmd.StartProfile(); err != nil {
		logger.Fatal(err)
	}

	InitI18n()
	Textdomain("dde-daemon")

	C.init()
	proxy.SetupProxy()

	logger.SetLogLevel(log.LevelInfo)
	if cmd.IsLogLevelNone() &&
		(utils.IsEnvExists(log.DebugLevelEnv) || utils.IsEnvExists(log.DebugMatchEnv)) {
		logger.Info("Log level is none and debug env exists, so ignore cmd.loglevel")
	} else {
		appLogLevel := cmd.LogLevel()
		logger.Info("App log level:", appLogLevel)
		// set all modules log level to appLogLevel
		loader.SetLogLevel(appLogLevel)
	}

	loader.EnableModules([]string{"dock", "launcher", "trayicon"}, nil, loader.EnableFlagIgnoreMissingModule)

	runMainLoop()
}
