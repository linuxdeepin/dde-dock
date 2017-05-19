/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package main

import "pkg.deepin.io/lib/log"

import "pkg.deepin.io/lib"
import "pkg.deepin.io/lib/dbus"
import "os"
import _ "pkg.deepin.io/dde/daemon/accounts"
import _ "pkg.deepin.io/dde/daemon/system/power"
import _ "pkg.deepin.io/dde/daemon/system/gesture"
import _ "pkg.deepin.io/dde/daemon/apps"
import "pkg.deepin.io/dde/daemon/loader"
import . "pkg.deepin.io/lib/gettext"
import "gir/glib-2.0"

var logger = log.NewLogger("daemon/dde-system-daemon")

func main() {
	logger.BeginTracing()
	defer logger.EndTracing()

	if !lib.UniqueOnSystem("com.deepin.daemon") {
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

	loader.StartAll()
	defer loader.StopAll()

	dbus.DealWithUnhandledMessage()
	// NOTE: system/power module requires glib loop
	go glib.StartLoop()

	if err := dbus.Wait(); err != nil {
		logger.Errorf("Lost dbus: %v", err)
		os.Exit(-1)
	} else {
		os.Exit(0)
	}
}
