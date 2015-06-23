/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package main

import (
	_ "pkg.linuxdeepin.com/dde-daemon/appearance"
	_ "pkg.linuxdeepin.com/dde-daemon/audio"
	_ "pkg.linuxdeepin.com/dde-daemon/bluetooth"
	_ "pkg.linuxdeepin.com/dde-daemon/clipboard"
	_ "pkg.linuxdeepin.com/dde-daemon/dock"
	_ "pkg.linuxdeepin.com/dde-daemon/dsc"
	_ "pkg.linuxdeepin.com/dde-daemon/inputdevices"
	_ "pkg.linuxdeepin.com/dde-daemon/keybinding"
	_ "pkg.linuxdeepin.com/dde-daemon/launcher"
	_ "pkg.linuxdeepin.com/dde-daemon/mime"
	_ "pkg.linuxdeepin.com/dde-daemon/mounts"
	_ "pkg.linuxdeepin.com/dde-daemon/mpris"
	_ "pkg.linuxdeepin.com/dde-daemon/network"
	_ "pkg.linuxdeepin.com/dde-daemon/power"
	_ "pkg.linuxdeepin.com/dde-daemon/screenedge"
	_ "pkg.linuxdeepin.com/dde-daemon/screensaver"
	_ "pkg.linuxdeepin.com/dde-daemon/sessionwatcher"
	_ "pkg.linuxdeepin.com/dde-daemon/systeminfo"
	_ "pkg.linuxdeepin.com/dde-daemon/timedate"
	_ "pkg.linuxdeepin.com/dde-daemon/xsettings"
	"pkg.linuxdeepin.com/lib/gio-2.0"
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
