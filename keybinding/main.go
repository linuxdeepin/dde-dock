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
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"pkg.linuxdeepin.com/lib/log"
)

var (
	logger = log.NewLogger("daemon/keybinding")
	X      *xgbutil.XUtil

	grabKeyBindsMap = make(map[KeycodeInfo]string)
	PrevSystemPairs = make(map[string]string)
	PrevCustomPairs = make(map[string]string)

	bindGSettings  *gio.Settings
	sysGSettings   *gio.Settings
	mediaGSettings *gio.Settings
)


func initGSettings() {
	bindGSettings  = gio.NewSettings("com.deepin.dde.keybinding")
	sysGSettings   = gio.NewSettings("com.deepin.dde.keybinding.system")
	mediaGSettings = gio.NewSettings("com.deepin.dde.keybinding.mediakey")
}

func StartKeyBinding() {
	var err error

	if X, err = xgbutil.NewConn(); err != nil {
		logger.Warning("New XGB Util Failed:", err)
		panic(err)
	}
	keybind.Initialize(X)
	initXRecord()

	initSystemIdDescList()
	//initMediaIdDescList()
	initWindowIdDescList()
	initWorkspaceIdDescList()

	grabKeyPairs(getSystemKeyPairs(), true)
	grabKeyPairs(getCustomKeyPairs(), true)
	grabMediaKeys()
}

func Start() {
	logger.BeginTracing()
	initGSettings()

	StartKeyBinding()

	if err := dbus.InstallOnSession(GetManager()); err != nil {
		logger.Error("Install DBus Failed:", err)
		panic(err)
	}

	if err := dbus.InstallOnSession(GetMediaManager()); err != nil {
		logger.Error("Install DBus Failed:", err)
		panic(err)
	}

	go xevent.Main(X)
}

func Stop() {
	logger.EndTracing()

	stopXRecord()
	xevent.Quit(X)
	dbus.UnInstallObject(GetManager())
}
