/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package screenedge

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
)

var _X *xgbutil.XUtil

func getXUtil() *xgbutil.XUtil {
	if _X == nil {
		if X, err := xgbutil.NewConn(); err != nil {
			logger.Warning("New xgbutil failed:", err)
			return nil
		} else {
			_X = X
		}
	}

	return _X
}

func getActiveWindow() (uint32, error) {
	X := getXUtil()
	xid, err := ewmh.ActiveWindowGet(X)
	if err != nil {
		logger.Warning("Get active window failed:", err)
		return 0, err
	}

	return uint32(xid), nil
}

func getWindowState(xid uint32) ([]string, error) {
	X := getXUtil()
	list, err := ewmh.WmStateGet(X, xproto.Window(xid))
	if err != nil {
		logger.Warning("Get window state failed:", err)
		return []string{}, err
	}

	return list, nil
}

func getWindowName(xid uint32) (string, error) {
	X := getXUtil()
	name, err := ewmh.WmNameGet(X, xproto.Window(xid))
	if err != nil {
		logger.Warning("Get window name failed:", err)
		return "", err
	}

	return name, nil
}

func getWindowPid(xid uint32) (uint32, error) {
	X := getXUtil()
	pid, err := ewmh.WmPidGet(X, xproto.Window(xid))
	if err != nil {
		logger.Warning("Get window pid failed:", err)
		return 0, err
	}

	return uint32(pid), nil
}

func isWindowFullscreen(xid uint32) (bool, error) {
	X := getXUtil()
	if X == nil {
		return false, nil
	}

	stateList, err := getWindowState(xid)
	if err != nil {
		return false, err
	}
	if isStrInList("_NET_WM_STATE_FULLSCREEN", stateList) {
		return true, nil
	}
	return false, nil
}

func isStrInList(key string, list []string) bool {
	for _, v := range list {
		if key == v {
			return true
		}
	}

	return false
}
