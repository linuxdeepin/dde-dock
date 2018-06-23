/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
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

package gesture

import (
	"github.com/linuxdeepin/go-dbus-factory/org.freedesktop.login1"
	x "github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/util/keybind"
	"github.com/linuxdeepin/go-x11-client/util/wm/ewmh"
	"pkg.deepin.io/lib/dbus1"
)

var (
	xconn  *x.Conn
	_dconn *dbus.Conn
	_self  *login1.Session
)

func isKbdAlreadyGrabbed() bool {
	if xconn == nil {
		conn, err := x.NewConn()
		if err != nil {
			return false
		}
		xconn = conn
	}

	var grabWin x.Window

	rootWin := xconn.GetDefaultScreen().Root
	if activeWin, _ := ewmh.GetActiveWindow(xconn).Reply(xconn); activeWin == 0 {
		grabWin = rootWin
	} else {
		// check viewable
		attrs, err := x.GetWindowAttributes(xconn, activeWin).Reply(xconn)
		if err != nil {
			grabWin = rootWin
		} else if attrs.MapState != x.MapStateViewable {
			// err is nil and activeWin is not viewable
			grabWin = rootWin
		} else {
			// err is nil, activeWin is viewable
			grabWin = activeWin
		}
	}

	err := keybind.GrabKeyboard(xconn, grabWin)
	if err == nil {
		// grab keyboard successful
		keybind.UngrabKeyboard(xconn)
		return false
	}

	logger.Warningf("GrabKeyboard win %d failed: %v", grabWin, err)

	gkErr, ok := err.(keybind.GrabKeyboardError)
	if ok && gkErr.Status == x.GrabStatusAlreadyGrabbed {
		return true
	}
	return false
}

func isSessionActive() bool {
	if _dconn == nil {
		conn, err := dbus.SystemBus()
		if err != nil {
			logger.Error("Failed to new system bus:", err)
			return false
		}
		_dconn = conn
	}

	if _self == nil {
		self, err := login1.NewSession(_dconn, "/org/freedesktop/login1/session/self")
		if err != nil {
			logger.Error("Failed to connect self session:", err)
			return false
		}
		_self = self
	}

	active, err := _self.Active().Get(dbus.FlagNoAutoStart)
	if err != nil {
		logger.Error("Failed to get self active:", err)
		return false
	}
	return active
}
