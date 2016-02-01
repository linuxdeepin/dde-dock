/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package screenedge

import (
	"fmt"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"io/ioutil"
	dutils "pkg.deepin.io/lib/utils"
	"strings"
)

func isAppInWhiteList(pid uint32) bool {
	filename := fmt.Sprintf("/proc/%v/cmdline", pid)
	if !dutils.IsFileExist(filename) {
		return false
	}

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Warningf("ReadFile '%s' failed: %v", filename, err)
		return false
	}

	whiteList := zoneSettings().GetStrv("white-list")
	for _, v := range whiteList {
		if strings.Contains(string(contents), v) {
			return true
		}
	}

	return false
}

func isAppInBlackList() bool {
	X := getXUtil()
	if X == nil {
		return false
	}

	xid, _ := getActiveWindow(X)
	pid, _ := getWindowPid(X, xid)

	filename := fmt.Sprintf("/proc/%v/cmdline", pid)
	if !dutils.IsFileExist(filename) {
		return false
	}

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Warningf("ReadFile '%s' failed: %v", filename, err)
		return false
	}

	blackList := zoneSettings().GetStrv("black-list")
	for _, target := range blackList {
		if strings.Contains(string(contents), target) {
			return true
		}
	}

	return false
}

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

func getActiveWindow(X *xgbutil.XUtil) (uint32, error) {
	xid, err := ewmh.ActiveWindowGet(X)
	if err != nil {
		logger.Warning("Get active window failed:", err)
		return 0, err
	}

	return uint32(xid), nil
}

func getWindowState(X *xgbutil.XUtil, xid uint32) ([]string, error) {
	list, err := ewmh.WmStateGet(X, xproto.Window(xid))
	if err != nil {
		logger.Warning("Get window state failed:", err)
		return []string{}, err
	}

	return list, nil
}

func getWindowName(X *xgbutil.XUtil, xid uint32) (string, error) {
	name, err := ewmh.WmNameGet(X, xproto.Window(xid))
	if err != nil {
		logger.Warning("Get window name failed:", err)
		return "", err
	}

	return name, nil
}

func getWindowPid(X *xgbutil.XUtil, xid uint32) (uint32, error) {
	pid, err := ewmh.WmPidGet(X, xproto.Window(xid))
	if err != nil {
		logger.Warning("Get window pid failed:", err)
		return 0, err
	}

	return uint32(pid), nil
}

func isActiveWindowFullscreen() (uint32, bool) {
	X := getXUtil()
	if X == nil {
		return 0, false
	}

	xid, _ := getActiveWindow(X)
	list, _ := getWindowState(X, xid)

	if strIsInInList("_NET_WM_STATE_FULLSCREEN", list) {
		pid, _ := getWindowPid(X, xid)
		return pid, true
	}

	return 0, false
}

func strIsInInList(key string, list []string) bool {
	for _, v := range list {
		if key == v {
			return true
		}
	}

	return false
}
