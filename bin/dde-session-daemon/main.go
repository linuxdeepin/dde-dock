/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package main

//#cgo pkg-config:gtk+-3.0
//#include <gtk/gtk.h>
//void init(){gtk_init(0,0);}
import "C"
import (
	// "gopkg.in/alecthomas/kingpin.v2"
	"os"
	"pkg.deepin.io/dde/daemon/loader"
	dapp "pkg.deepin.io/lib/app"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/proxy"
)

var logger = log.NewLogger("daemon/dde-session-daemon")

func main() {
	InitI18n()
	Textdomain("dde-daemon")

	cmd := dapp.New("dde-session-dameon", "dde session daemon", "version "+__VERSION__)

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
	if err := app.register(); err != nil {
		logger.Info(err)
		os.Exit(0)
	}

	app.logLevel = cmd.LogLevel()
	loader.SetLogLevel(app.logLevel)
	logger.Info("LogLevel is", app.logLevel)

	var err error
	needRunMainLoop := true
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
