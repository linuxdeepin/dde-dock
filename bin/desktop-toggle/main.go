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
	liblogger "pkg.linuxdeepin.com/lib/logger"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
)

func main() {
	logger := liblogger.NewLogger("Desktop Toggle")
	defer logger.EndTracing()

	X, err := xgbutil.NewConn()
	if err != nil {
		logger.Info("New xgbutil connection failed: ", err)
		panic(err)
	}

	if ret, err := ewmh.ShowingDesktopGet(X); err == nil {
		// !!! NOT using ewmh.ShowingDesktopReq
		// because ewmh.ShowingDesktopReq passes a uint argument,
		// and int is used on xevent.NewClientMessage.
		logger.Info("Show Desktop Status: ", ret)
		var showInt int
		if ret {
			showInt = 0
		} else {
			showInt = 1
		}
		ewmh.ClientEvent(X, X.RootWin(), "_NET_SHOWING_DESKTOP", showInt)
	}
	X.Sync()
}
