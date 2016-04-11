/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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
	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/app"
	"pkg.deepin.io/lib/dbus"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/proxy"

	"os"
	"time"
)

var logger = log.NewLogger("daemon/dock-daemon")

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

type DockDaemon struct {
}

func (*DockDaemon) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.daemon.DockDaemon",
		ObjectPath: "/com/deepin/daemon/DockDaemon",
		Interface:  "com.deepin.daemon.DockDaemon",
	}
}

func main() {
	dockDaemon := new(DockDaemon)
	if !lib.UniqueOnSession(dockDaemon.GetDBusInfo().Dest) {
		logger.Warning("There's a dde-dock-daemon instance running.")
		os.Exit(0)
	}

	err := dbus.InstallOnSession(dockDaemon)
	if err != nil {
		logger.Fatal(err)
	}

	cmd := app.New("dde-dock-daemon", "daemon for dde-dock", "version "+__VERSION__)
	cmd.ParseCommandLine(os.Args[1:])
	if err := cmd.StartProfile(); err != nil {
		logger.Fatal(err)
	}

	InitI18n()
	Textdomain("dde-daemon")

	C.init()
	proxy.SetupProxy()

	loader.SetLogLevel(cmd.LogLevel())
	loader.EnableModules([]string{"dock"}, nil, loader.EnableFlagIgnoreMissingModule)

	runMainLoop()
}
