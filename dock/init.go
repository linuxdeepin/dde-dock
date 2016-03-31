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
	"pkg.deepin.io/dde/daemon/loader"
	"time"
)

func init() {
	loader.Register(NewDaemon(logger))
}

var (
	XU     *xgbutil.XUtil
	TrayXU *xgbutil.XUtil

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


	mouseAreaTimer   *time.Timer
	TOGGLE_HIDE_TIME = time.Millisecond * 400
)
