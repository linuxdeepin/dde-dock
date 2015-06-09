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
	"pkg.linuxdeepin.com/dde-daemon"
	"pkg.linuxdeepin.com/lib/gio-2.0"
)

var (
	// modules should be loaded in following order
	orderedModules = []string{
		"inputdevices",
		"screensaver",
		"power",
		"audio",
		"appearance",
		"clipboard",
		"timedate",
		"mime",
		"screenedge",
		"bluetooth",
		"network",
		"mounts",
		"dock",
		"launcher",
		"keybinding",
		"dsc",
		"mpris",
		"systeminfo",
		"sessionwatcher",
	}
	daemonSettings = gio.NewSettings("com.deepin.dde.daemon")
)

func initModules() {
	for _, name := range orderedModules {
		loader.Enable(name, daemonSettings.GetBoolean(name))
	}
}

func listenDaemonSettings() {
	daemonSettings.Connect("changed", func(s *gio.Settings, name string) {
		// gsettings key names must keep consistent with module names
		enable := daemonSettings.GetBoolean(name)
		loader.Enable(name, enable)
		if enable {
			loader.Start(name)
		} else {
			loader.Stop(name)
		}
	})
	daemonSettings.GetBoolean("mounts")
}
