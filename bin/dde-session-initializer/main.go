/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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
	"pkg.deepin.io/dde/daemon/loader"
	_ "pkg.deepin.io/dde/daemon/trayicon"
	_ "pkg.deepin.io/dde/daemon/x_event_monitor"
	"pkg.deepin.io/lib/app"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/gsettings"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/proxy"
	"pkg.deepin.io/lib/utils"

	"os"
	"time"

	"pkg.deepin.io/lib/dbusutil"
)

var logger = log.NewLogger("daemon/initializer")

func runMainLoop() {
	err := gsettings.StartMonitor()
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("register session")
	startTime := time.Now()
	session.Register()
	logger.Info("register session done, cost", time.Now().Sub(startTime))

	logger.Info("DealWithUnhandledMessage")
	startTime = time.Now()
	logger.Info("DealWithUnhandledMessage done, cost", time.Now().Sub(startTime))
	go glib.StartLoop()

	logger.Info("initialize done")
}

const dbusServiceName = "com.deepin.daemon.Initializer"

func main() {
	service, err := dbusutil.NewSessionService()
	if err != nil {
		logger.Fatal(err)
	}
	loader.SetService(service)

	err = service.RequestName(dbusServiceName)
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

	loader.EnableModules([]string{"dock", "trayicon", "x_event_monitor"}, nil, loader.EnableFlagIgnoreMissingModule)
	runMainLoop()
	service.Wait()
}
