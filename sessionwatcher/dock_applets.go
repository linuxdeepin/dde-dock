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

package sessionwatcher

import (
	"os/exec"
)

type DockApplet struct{}

const (
	_DOCK_APPLET_SENDER = "dde.dock.entry.AppletManager"
	_DOCK_APPLET_CMD    = "/usr/bin/dde-dock-applets"
)

func NewDockApplet() *DockApplet {
	applet := &DockApplet{}

	return applet
}

//TODO:
// Through the dock applet xid to determine whether dde-dock-applet normal
// handle when no dock applet plugin show
func (applet *DockApplet) restartDockApplet() {
	if isDBusSenderExist(_DOCK_APPLET_SENDER) {
		return
	}

	if _, err := exec.Command("/usr/bin/killall", _DOCK_APPLET_CMD).Output(); err != nil {
		logger.Warning("killall dde-dock-applets failed:", err)
	}

	if err := exec.Command(_DOCK_APPLET_CMD, "").Run(); err != nil {
		logger.Warning("launch dde-dock-applets failed:", err)
		return
	}
}
