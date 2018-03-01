/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package dock

import (
	"path/filepath"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xprop"
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/xdg/basedir"
)

func init() {
	loader.Register(NewDaemon(logger))
}

var (
	logger      = log.NewLogger("daemon/dock")
	homeDir     string
	scratchDir  string
	dockManager *Manager

	XU *xgbutil.XUtil

	atomNetShowingDesktop xproto.Atom
	atomNetClientList     xproto.Atom
	atomNetActiveWindow   xproto.Atom
	atomNetWMIcon         xproto.Atom
	atomNetWMName         xproto.Atom
	atomNetWMState        xproto.Atom
	atomNetWMWindowType   xproto.Atom
	atomXEmbedInfo        xproto.Atom
)

func initDir() {
	homeDir = basedir.GetUserHomeDir()
	scratchDir = filepath.Join(basedir.GetUserConfigDir(), "dock/scratch")
	logger.Debugf("scratch dir: %q", scratchDir)
}

func initAtom() {
	atomNetShowingDesktop, _ = xprop.Atm(XU, "_NET_SHOWING_DESKTOP")
	atomNetClientList, _ = xprop.Atm(XU, "_NET_CLIENT_LIST")
	atomNetActiveWindow, _ = xprop.Atm(XU, "_NET_ACTIVE_WINDOW")
	atomNetWMIcon, _ = xprop.Atm(XU, "_NET_WM_ICON")
	atomNetWMName, _ = xprop.Atm(XU, "_NET_WM_NAME")
	atomNetWMState, _ = xprop.Atm(XU, "_NET_WM_STATE")
	atomNetWMWindowType, _ = xprop.Atm(XU, "_NET_WM_WINDOW_TYPE")
	atomXEmbedInfo, _ = xprop.Atm(XU, "_XEMBED_INFO")
}
