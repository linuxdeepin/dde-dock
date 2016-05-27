/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package main

import (
	"flag"
	"os"
	"pkg.deepin.io/dde/daemon/grub2"
	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
	"time"
)

var logger = log.NewLogger("daemon/grub2ext-runner")
var argDebug bool

func main() {
	logger.BeginTracing()
	defer logger.EndTracing()

	if !lib.UniqueOnSystem(grub2.DbusGrub2ExtDest) {
		logger.Warning("There already has an Grub2 daemon running.")
		return
	}

	flag.BoolVar(&argDebug, "d", false, "debug mode")
	flag.BoolVar(&argDebug, "debug", false, "debug mode")
	flag.Parse()
	if argDebug {
		logger.Info("debug mode")
		logger.SetLogLevel(log.LevelDebug)
	}

	ge := grub2.NewGrub2Ext()
	err := dbus.InstallOnSystem(ge)
	if err != nil {
		logger.Errorf("register dbus interface failed: %v", err)
		os.Exit(1)
	}

	dbus.DealWithUnhandledMessage()
	dbus.SetAutoDestroyHandler(300*time.Second, nil)

	if err := dbus.Wait(); err != nil {
		logger.Error("lost dbus session:", err)
		os.Exit(1)
	}
}
