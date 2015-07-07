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
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"pkg.deepin.io/lib/log"
)

func main() {
	logger := log.NewLogger("desktop-toggle")
	logger.BeginTracing()
	defer logger.EndTracing()

	X, err := xgbutil.NewConn()
	if err != nil {
		logger.Error("New xgbutil connection failed: ", err)
		return
	}

	ret, err := ewmh.ShowingDesktopGet(X)
	if err != nil {
		logger.Warning("Get showing desktop state failed:", err)
		return
	}
	logger.Info("Desktop showing state:", ret)

	var showInt int
	if ret {
		showInt = 0
	} else {
		showInt = 1
	}

	// !!! NOT using ewmh.ShowingDesktopReq
	// because ewmh.ShowingDesktopReq passes a uint argument,
	// and int is used on xevent.NewClientMessage.
	err = ewmh.ClientEvent(X, X.RootWin(), "_NET_SHOWING_DESKTOP", showInt)
	if err != nil {
		logger.Warning("Send showing desktop client event failed: ", err)
		return
	}
	X.Sync()
}
