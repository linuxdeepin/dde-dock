/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dock

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xprop"
	"path/filepath"
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
	dockManager *DockManager

	XU *xgbutil.XUtil

	//There variable must be initialized after the Xu/TrayXU has been
	//created.
	_NET_SHOWING_DESKTOP    xproto.Atom
	_NET_CLIENT_LIST        xproto.Atom
	_NET_ACTIVE_WINDOW      xproto.Atom
	ATOM_WINDOW_ICON        xproto.Atom
	ATOM_WINDOW_NAME        xproto.Atom
	ATOM_WINDOW_STATE       xproto.Atom
	ATOM_WINDOW_TYPE        xproto.Atom
	ATOM_DOCK_APP_ID        xproto.Atom
	_NET_SYSTEM_TRAY_S0     xproto.Atom
	_NET_SYSTEM_TRAY_OPCODE xproto.Atom
	ATOM_XEMBED_INFO        xproto.Atom
)

func initDir() {
	homeDir = basedir.GetUserHomeDir()
	scratchDir = filepath.Join(basedir.GetUserConfigDir(), "dock/scratch")
	logger.Debugf("scratch dir: %q", scratchDir)
}

func initAtom() {
	_NET_SHOWING_DESKTOP, _ = xprop.Atm(XU, "_NET_SHOWING_DESKTOP")
	_NET_CLIENT_LIST, _ = xprop.Atm(XU, "_NET_CLIENT_LIST")
	_NET_ACTIVE_WINDOW, _ = xprop.Atm(XU, "_NET_ACTIVE_WINDOW")
	ATOM_WINDOW_ICON, _ = xprop.Atm(XU, "_NET_WM_ICON")
	ATOM_WINDOW_NAME, _ = xprop.Atm(XU, "_NET_WM_NAME")
	ATOM_WINDOW_STATE, _ = xprop.Atm(XU, "_NET_WM_STATE")
	ATOM_WINDOW_TYPE, _ = xprop.Atm(XU, "_NET_WM_WINDOW_TYPE")
	ATOM_XEMBED_INFO, _ = xprop.Atm(XU, "_XEMBED_INFO")
}
