/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package main

import "os"
import . "pkg.deepin.io/dde/daemon/langselector"
import "pkg.deepin.io/lib/gettext"
import "pkg.deepin.io/lib/dbus"
import "time"

func main() {
	gettext.InitI18n()
	gettext.Textdomain("dde-daemon")

	lang := Start()
	if lang == nil {
		return
	}

	dbus.DealWithUnhandledMessage()

	dbus.SetAutoDestroyHandler(time.Minute*5, func() bool {
		if lang.LocaleState == LocaleStateChanging {
			return false
		} else {
			return true
		}
	})

	if err := dbus.Wait(); err != nil {
		Stop()
		os.Exit(-1)
	}

	Stop()
	os.Exit(0)
}
