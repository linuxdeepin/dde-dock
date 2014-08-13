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

package screen_edges

import (
	"fmt"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"io/ioutil"
	dutils "pkg.linuxdeepin.com/lib/utils"
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

	whiteList := zoneSettings.GetStrv("white-list")
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

	xid := getActiveWindow(X)
	pid := getWindowPid(X, xid)

	filename := fmt.Sprintf("/proc/%v/cmdline", pid)
	if !dutils.IsFileExist(filename) {
		return false
	}

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Warningf("ReadFile '%s' failed: %v", filename, err)
		return false
	}

	blackList := zoneSettings.GetStrv("black-list")
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

func getActiveWindow(X *xgbutil.XUtil) uint32 {
	xid, err := ewmh.ActiveWindowGet(X)
	if err != nil {
		logger.Warning("Get active window failed:", err)
		return 0
	}

	return uint32(xid)
}

func getWindowState(X *xgbutil.XUtil, xid uint32) []string {
	list, err := ewmh.WmStateGet(X, xproto.Window(xid))
	if err != nil {
		logger.Warning("Get window state failed:", err)
	}

	return list
}

func getWindowName(X *xgbutil.XUtil, xid uint32) string {
	name, err := ewmh.WmNameGet(X, xproto.Window(xid))
	if err != nil {
		logger.Warning("Get window name failed:", err)
		return ""
	}

	return name
}

func getWindowPid(X *xgbutil.XUtil, xid uint32) uint32 {
	pid, err := ewmh.WmPidGet(X, xproto.Window(xid))
	if err != nil {
		logger.Warning("Get window pid failed:", err)
	}

	return uint32(pid)
}

func isActiveWindowFullscreen() (uint32, bool) {
	X := getXUtil()
	if X == nil {
		return 0, false
	}

	xid := getActiveWindow(X)
	list := getWindowState(X, xid)

	if strIsInInList("_NET_WM_STATE_FULLSCREEN", list) {
		pid := getWindowPid(X, xid)
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
