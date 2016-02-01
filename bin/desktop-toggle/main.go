/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
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
		// '_NET_SHOWING_DESKTOP' not exist, means not showing
		ret = false
	}
	logger.Debug("Desktop showing state:", ret)

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
