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
	// "gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/user"
	"runtime"

	"pkg.deepin.io/dde/daemon/loader"
	dapp "pkg.deepin.io/lib/app"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/proxy"
	"pkg.deepin.io/lib/utils"
)

var logger = log.NewLogger("daemon/dde-session-daemon")

func main() {
	InitI18n()
	Textdomain("dde-daemon")

	cmd := dapp.New("dde-session-dameon", "dde session daemon", "version "+__VERSION__)

	usr, err := user.Current()
	if err == nil {
		os.Chdir(usr.HomeDir)
	}

	flags := new(Flags)
	flags.IgnoreMissingModules = cmd.Flag("Ignore", "ignore missing modules, --no-ignore to revert it.").Short('i').Default("true").Bool()
	flags.ForceStart = cmd.Flag("force", "Force start disabled module.").Short('f').Bool()

	cmd.Command("auto", "Automatically get enabled and disabled modules from settings.").Default()
	enablingModules := cmd.Command("enable", "Enable modules and their dependencies, ignore settings.").Arg("module", "module names.").Required().Strings()
	disableModules := cmd.Command("disable", "Disable modules, ignore settings.").Arg("module", "module names.").Required().Strings()
	listModule := cmd.Command("list", "List all the modules or the dependencies of one module.").Arg("module", "module name.").String()

	subCmd := cmd.ParseCommandLine(os.Args[1:])
	cmd.StartProfile()

	C.init()
	proxy.SetupProxy()

	app := NewSessionDaemon(flags, daemonSettings, logger)
	if err = app.register(); err != nil {
		logger.Info(err)
		os.Exit(0)
	}

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

	needRunMainLoop := true

	// Ensure each module and mainloop in the same thread
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	switch subCmd {
	case "auto":
		app.execDefaultAction()
	case "enable":
		err = app.EnableModules(*enablingModules)
	case "disable":
		err = app.DisableModules(*disableModules)
	case "list":
		err = app.ListModule(*listModule)
		needRunMainLoop = false
	}

	if err != nil {
		logger.Warning(err)
		os.Exit(1)
	}

	if needRunMainLoop {
		runMainLoop()
	}
}
