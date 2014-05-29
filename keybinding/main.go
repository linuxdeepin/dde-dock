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
	"dlib"
	"dlib/dbus"
	"dlib/gio-2.0"
	libLogger "dlib/logger"
	libUtils "dlib/utils"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"os"
)

var (
	utilsObj = libUtils.NewUtils()
	logObj   = libLogger.NewLogger("daemon/keybinding")
	X        *xgbutil.XUtil

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
		logObj.Warning("New XGB Util Failed:", err)
		panic(err)
	}
	keybind.Initialize(X)
	initXRecord()

	grabKeyPairs(getSystemKeyPairs(), true)
	grabKeyPairs(getCustomKeyPairs(), true)
	grabMediaKeys()
}

func main() {
	if !dlib.UniqueOnSession(KEYBIND_DEST) {
		logObj.Error("keybinding has running")
		return
	}

	defer logObj.EndTracing()
	StartKeyBinding()

	if err := dbus.InstallOnSession(GetManager()); err != nil {
		logObj.Error("Install DBus Failed:", err)
		panic(err)
	}

	if err := dbus.InstallOnSession(GetMediaManager()); err != nil {
		logObj.Error("Install DBus Failed:", err)
		panic(err)
	}

	dbus.DealWithUnhandledMessage()

	go dlib.StartLoop()
	go xevent.Main(X)

	if err := dbus.Wait(); err != nil {
		logObj.Errorf("Lost DBus: %v", err)
		os.Exit(-1)
	} else {
		os.Exit(0)
	}
}
