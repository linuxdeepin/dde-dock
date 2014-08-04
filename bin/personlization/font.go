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
	"dbus/com/deepin/sessionmanager"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/gio-2.0"
)

const (
	FONT_PATH = "/com/deepin/daemon/FontSettings"
	FONT_IFC  = "com.deepin.daemon.FontSettings"
)

type FontSettings struct {
	xs         *sessionmanager.XSettings
	wmSettings *gio.Settings
}

func (fs *FontSettings) initSettings() {
	if fs.xs == nil {
		var err error
		if fs.xs, err = sessionmanager.NewXSettings("com.deepin.SessionManager", "/com/deepin/XSettings"); err != nil {
			Logger.Fatal("New XSettings Failed:", err)
			return
		}
	}

	if fs.wmSettings == nil {
		fs.wmSettings = gio.NewSettings("org.gnome.desktop.wm.preferences")
	}
}

var _fs *FontSettings

func GetFontSettings() *FontSettings {
	if _fs == nil {
		fs := &FontSettings{}
		fs.initSettings()
		_fs = fs
	}

	return _fs
}

func StartFont() {
	if err := dbus.InstallOnSession(GetFontSettings()); err != nil {
		Logger.Fatal("Install DBus Failed:", err)
		return
	}
}
