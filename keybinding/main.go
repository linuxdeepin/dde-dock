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
	bindGSettings = gio.NewSettings("com.deepin.dde.keybinding")
	sysGSettings = gio.NewSettings("com.deepin.dde.keybinding.system")
	mediaGSettings = gio.NewSettings("com.deepin.dde.keybinding.mediakey")
}

func finiGSettings() {
	if bindGSettings != nil {
		bindGSettings.Unref()
	}

	if sysGSettings != nil {
		sysGSettings.Unref()
	}

	if mediaGSettings != nil {
		mediaGSettings.Unref()
	}
}

func StartKeyBinding() error {
	var err error
	X, err = xgbutil.NewConn()
	if err != nil {
		return err
	}
	keybind.Initialize(X)
	initXRecord()

	initSystemIdDescList()
	//initMediaIdDescList()
	initWindowIdDescList()
	initWorkspaceIdDescList()

	grabKeyPairs(getSystemKeyPairs(), true)
	grabKeyPairs(getCustomKeyPairs(), true)
	grabMediaKeys(true)

	return nil
}

func endKeyBinding() {
	if X == nil {
		return
	}

	stopXRecord()
	grabMediaKeys(false)
	grabKeyPairs(getSystemKeyPairs(), false)
	grabKeyPairs(getCustomKeyPairs(), false)
	xevent.Quit(X)
	X = nil
}

var (
	_manager *Manager
)

func finalize() {
	finiGSettings()
	endKeyBinding()

	dbus.UnInstallObject(_manager)
	_manager = nil
	logger.EndTracing()
}

func Start() {
	if _manager != nil {
		return
	}

	logger.BeginTracing()
	initGSettings()

	err := StartKeyBinding()
	if err != nil {
		logger.Error("failed start keybinding:", err)
		logger.EndTracing()
		finiGSettings()
		return
	}

	_manager = newManager()
	if _manager == nil {
		logger.Error("Create keybinding manager failed")
		finiGSettings()
		endKeyBinding()
		return
	}

	err = dbus.InstallOnSession(_manager)
	if err != nil {
		logger.Error("Install DBus Failed:", err)
		finalize()
		return
	}

	err = dbus.InstallOnSession(_manager.mediaKey)
	if err != nil {
		logger.Error("Install DBus Failed:", err)
		finalize()
		return
	}

	go xevent.Main(X)
}

func Stop() {
	if _manager == nil {
		return
	}

	finalize()
}
