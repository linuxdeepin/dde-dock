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

package screenedge

/*
import (
	libarea "dbus/com/deepin/api/xmousearea"
	libdsp "dbus/com/deepin/daemon/display"
	"dbus/com/deepin/dde/launcher"
	C "launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) {
	C.TestingT(t)
}

func init() {
	var err error
	dspObj, err = libdsp.NewDisplay("com.deepin.daemon.Display",
		"/com/deepin/daemon/Display")
	if err != nil {
		logger.Fatal("New Display Failed: ", err)
	}

	areaObj, err = libarea.NewXMouseArea("com.deepin.api.XMouseArea",
		"/com/deepin/api/XMouseArea")
	if err != nil {
		logger.Fatal("New XMouseArea Failed: ", err)
	}

	launchObj, err = launcher.NewLauncher("com.deepin.dde.launcher",
		"/com/deepin/dde/launcher")
	if err != nil {
		logger.Fatal("New DDE Launcher Failed: ", err)
	}

	C.Suite(newManager())
}

func (m *Manager) TestWindowState(c *C.C) {
	xu := getXUtil()
	if xu == nil {
		c.Error("Get XUtil Failed")
		return
	}

	var xid uint32
	var err error

	if xid, err = getActiveWindow(xu); err != nil {
		c.Error("getActiveWindow failed:", err)
		return
	}

	if _, err = getWindowName(xu, xid); err != nil {
		c.Error("getWindowName failed:", err)
		return
	}

	if _, err = getWindowState(xu, xid); err != nil {
		c.Error("getWindowState failed:", err)
		return
	}

	if _, err = getWindowPid(xu, xid); err != nil {
		c.Error("getWindowPid failed:", err)
		return
	}
}
*/
