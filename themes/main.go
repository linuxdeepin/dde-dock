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

package themes

import (
	"dbus/com/deepin/api/greeterutils"
	"dbus/com/deepin/sessionmanager"
	"path"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"pkg.linuxdeepin.com/lib/logger"
	dutils "pkg.linuxdeepin.com/lib/utils"
)

var (
	objXS      *sessionmanager.XSettings
	greeterObj *greeterutils.GreeterUtils
	Logger     = logger.NewLogger(MANAGER_DEST)

	themeSettings = gio.NewSettings("com.deepin.dde.personalization")
	gnmSettings   = gio.NewSettings("org.gnome.desktop.background")

	homeDir = dutils.GetHomeDir()
)

func Start() {
	Logger.BeginTracing()

	var err error
	objXS, err = sessionmanager.NewXSettings("com.deepin.SessionManager",
		"/com/deepin/XSettings")
	if err != nil {
		Logger.Fatal("New XSettings Failed:", err)
	}

	if greeterObj, err = greeterutils.NewGreeterUtils("com.deepin.api.GreeterUtils", "/com/deepin/api/GreeterUtils"); err != nil {
		Logger.Fatal("New GreeterUtils Failed:", err)
	}

	if err = dbus.InstallOnSession(GetManager()); err != nil {
		Logger.Fatal("Install DBus Failed", err)
	}

	username := dutils.GetUserName()
	homeDir := dutils.GetHomeDir()
	if homeDir != path.Join("/tmp", username) {
		if len(username) > 0 {
			greeterTheme := GetManager().GreeterTheme.GetValue().(string)
			greeterObj.SetGreeterTheme(username, greeterTheme)
		}
	}
}

func Stop() {
	obj := GetManager()
	obj.destroyAllTheme()
	obj.quitFlag <- true
	obj.watcher.Close()
	obj.bgQuitFlag <- true
	obj.bgWatcher.Close()
	dbus.UnInstallObject(obj)

	Logger.EndTracing()
}
