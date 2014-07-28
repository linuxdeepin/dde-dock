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
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"io/ioutil"
	"path"
	dutils "pkg.linuxdeepin.com/lib/utils"
	"strconv"
	"strings"
)

func isAppInWhiteList(pid uint32) bool {
	pidStr := strconv.FormatUint(uint64(pid), 10)
	filename := path.Join("/proc", pidStr, "cmdline")
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

func getWindowPid(X *xgbutil.XUtil, xid uint32) uint32 {
	pid, err := ewmh.WmPidGet(X, xproto.Window(xid))
	if err != nil {
		logger.Warning("Get window pid failed:", err)
	}

	return uint32(pid)
}

func isActiveWindowFullscreen() (uint32, bool) {
	X, err := xgbutil.NewConn()
	if err != nil {
		logger.Warning("New xgbutil failed:", err)
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
