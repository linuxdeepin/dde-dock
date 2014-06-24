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

package keybinding

import (
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	libLogger "pkg.linuxdeepin.com/lib/logger"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
)

var (
	Logger = libLogger.NewLogger("daemon/keybinding")
	X      *xgbutil.XUtil

	grabKeyBindsMap = make(map[KeycodeInfo]string)
	PrevSystemPairs = make(map[string]string)
	PrevCustomPairs = make(map[string]string)

	bindGSettings  = gio.NewSettings("com.deepin.dde.keybinding")
	sysGSettings   = gio.NewSettings("com.deepin.dde.keybinding.system")
	mediaGSettings = gio.NewSettings("com.deepin.dde.keybinding.mediakey")
	wmGSettings    = gio.NewSettings("org.gnome.desktop.wm.keybindings")
	putGSettings   = gio.NewSettingsWithPath("org.compiz.put",
		"/org/compiz/profiles/put/")
	shiftGSettings = gio.NewSettingsWithPath("org.compiz.shift",
		"/org/compiz/profiles/shift/")
)

func StartKeyBinding() {
	var err error

	if X, err = xgbutil.NewConn(); err != nil {
		Logger.Warning("New XGB Util Failed:", err)
		panic(err)
	}
	keybind.Initialize(X)
	initXRecord()

	initSystemNameDescMap()
	initWindowNameDescMap()
	initMediaNameDescMap()
	initWorkspaceNameDescMap()

	grabKeyPairs(getSystemKeyPairs(), true)
	grabKeyPairs(getCustomKeyPairs(), true)
	grabMediaKeys()
}

func Start() {
	Logger.BeginTracing()

	StartKeyBinding()

	if err := dbus.InstallOnSession(GetManager()); err != nil {
		Logger.Error("Install DBus Failed:", err)
		panic(err)
	}

	if err := dbus.InstallOnSession(GetMediaManager()); err != nil {
		Logger.Error("Install DBus Failed:", err)
		panic(err)
	}

	go xevent.Main(X)
}

func Stop() {
	Logger.EndTracing()

	stopXRecord()
	xevent.Quit(X)
	dbus.UnInstallObject(GetManager())
}
