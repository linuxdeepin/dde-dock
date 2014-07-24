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
	libLogger "pkg.linuxdeepin.com/lib/logger"
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
	coreSettings   = gio.NewSettingsWithPath(COMPIZ_SETTINGS_CORE,
		COMPIZ_SETTINGS_BASE_PATH+"core/")
	moveSettings = gio.NewSettingsWithPath(COMPIZ_SETTINGS_MOVE,
		COMPIZ_SETTINGS_BASE_PATH+"move/")
	resizeSettings = gio.NewSettingsWithPath(COMPIZ_SETTINGS_RESIZE,
		COMPIZ_SETTINGS_BASE_PATH+"resize/")
	vpswitchSettings = gio.NewSettingsWithPath(COMPIZ_SETTINGS_VPSWITCH,
		COMPIZ_SETTINGS_BASE_PATH+"vpswitch/")
	putGSettings = gio.NewSettingsWithPath(COMPIZ_SETTINGS_PUT,
		COMPIZ_SETTINGS_BASE_PATH+"put/")
	wallSettings = gio.NewSettingsWithPath(COMPIZ_SETTINGS_WALL,
		COMPIZ_SETTINGS_BASE_PATH+"wall/")
	shiftGSettings = gio.NewSettingsWithPath(COMPIZ_SETTINGS_SHIFT,
		COMPIZ_SETTINGS_BASE_PATH+"shift/")
	switcherSettings = gio.NewSettingsWithPath(COMPIZ_SETTINGS_SWITCHER,
		COMPIZ_SETTINGS_BASE_PATH+"switcher/")
)

func StartKeyBinding() {
	var err error

	if X, err = xgbutil.NewConn(); err != nil {
		Logger.Warning("New XGB Util Failed:", err)
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
