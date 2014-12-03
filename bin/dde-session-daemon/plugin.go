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
	_pluginList = []string{
		"inputdevices",
		"screensaver",
		"power",
		"audio",
		"appearance",
		"clipboard",
		"datetime",
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

	_daemonSettings = gio.NewSettings("com.deepin.dde.daemon")
)

func initPlugins() {
	for _, plugin := range _pluginList {
		enable := _daemonSettings.GetBoolean(plugin)
		if !enable {
			loader.Enable(plugin, false)
		}
	}
}

func listenDaemonSettings() {
	_daemonSettings.Connect("changed", func(s *gio.Settings, key string) {
		enable := _daemonSettings.GetBoolean(key)
		if enable {
			logger.Info("Enable plugin:", key)
			loader.StartPlugin(key)
		} else {
			logger.Info("Disable plugin:", key)
			loader.StopPlugin(key)
		}
	})
}
