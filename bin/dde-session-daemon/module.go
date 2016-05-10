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
	_ "pkg.deepin.io/dde/daemon/appearance"
	_ "pkg.deepin.io/dde/daemon/audio"
	_ "pkg.deepin.io/dde/daemon/bluetooth"
	_ "pkg.deepin.io/dde/daemon/clipboard"
	//_ "pkg.deepin.io/dde/daemon/dock"
	"gir/gio-2.0"
	_ "pkg.deepin.io/dde/daemon/inputdevices"
	_ "pkg.deepin.io/dde/daemon/keybinding"
	// _ "pkg.deepin.io/dde/daemon/launcher"
	_ "pkg.deepin.io/dde/daemon/mounts"
	_ "pkg.deepin.io/dde/daemon/mpris"
	_ "pkg.deepin.io/dde/daemon/network"
	_ "pkg.deepin.io/dde/daemon/power"
	_ "pkg.deepin.io/dde/daemon/screenedge"
	_ "pkg.deepin.io/dde/daemon/screensaver"
	_ "pkg.deepin.io/dde/daemon/sessionwatcher"
	_ "pkg.deepin.io/dde/daemon/systeminfo"
	_ "pkg.deepin.io/dde/daemon/timedate"
)

var (
	daemonSettings = gio.NewSettings("com.deepin.dde.daemon")
)

// TODO:
// func listenDaemonSettings() {
// 	daemonSettings.Connect("changed", func(s *gio.Settings, name string) {
// 		// gsettings key names must keep consistent with module names
// 		enable := daemonSettings.GetBoolean(name)
// 		loader.Enable(name, enable)
// 		if enable {
// 			loader.Start(name)
// 		} else {
// 			loader.Stop(name)
// 		}
// 	})
// 	daemonSettings.GetBoolean("mounts")
// }
