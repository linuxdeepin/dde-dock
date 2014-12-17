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

const (
	_DDE_DOCK_SENDER = "com.deepin.dde.dock"
	_DDE_DOCK_CMD    = "/usr/bin/dde-dock"
)

type Dock struct{}

func NewDock() *Dock {
	dock := &Dock{}

	return dock
}

func (dock *Dock) fixSwitchMode() {
	_, err := exec.Command("/usr/bin/xdotool", "search", "--onlyvisible", "dde-dock", "windowmove", "%1", "x", "y").Output()
	if err != nil {
		logger.Warning("fixSwitchMode failed:", err)
	}
}

func (dock *Dock) restartDock() {
	if isDBusSenderExist(_DDE_DOCK_SENDER) {
		dock.fixSwitchMode()
		return
	}

	if _, err := exec.Command("/usr/bin/killall", _DDE_DOCK_CMD).Output(); err != nil {
		logger.Warning("killall dde-dock failed:", err)
	}

	if err := exec.Command(_DDE_DOCK_CMD, "").Run(); err != nil {
		logger.Warning("launch dde-dock failed:", err)
		return
	}
}
