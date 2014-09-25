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

import (
	"pkg.linuxdeepin.com/lib/dbus"
)

type Manager struct{}

const (
	ZONE_DEST = "com.deepin.daemon.Zone"
	ZONE_PATH = "/com/deepin/daemon/Zone"
	ZONE_IFC  = "com.deepin.daemon.Zone"
)

func (op *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		ZONE_DEST,
		ZONE_PATH,
		ZONE_IFC,
	}
}

func (op *Manager) listenSignal() {
	dspObj.ConnectPrimaryChanged(func(argv []interface{}) {
		unregisterZoneArea()
		registerZoneArea()
	})

	areaObj.ConnectCursorInto(func(x, y int32, id string) {
		if id != areaId {
			return
		}

		if isAppInBlackList() {
			return
		}

		if pid, ok := isActiveWindowFullscreen(); ok {
			if !isAppInWhiteList(pid) {
				return
			}
		}

		if isInArea(x, y, topLeftArea) {
			execEdgeAction(EDGE_TOPLEFT)
		} else if isInArea(x, y, bottomLeftArea) {
			execEdgeAction(EDGE_BOTTOMLEFT)
		} else if isInArea(x, y, topRightArea) {
			execEdgeAction(EDGE_TOPRIGHT)
		} else if isInArea(x, y, bottomRightArea) {
			execEdgeAction(EDGE_BOTTOMRIGHT)
		}
	})

	areaObj.ConnectCancelAllArea(func() {
		unregisterZoneArea()
		registerZoneArea()
	})

	launchObj.ConnectShown(func() {
		enableOneEdge(getEdgeForCommand("/usr/bin/dde-launcher"))
	})

	launchObj.ConnectClosed(func() {
		op.enableAllEdge()
	})
}

func (op *Manager) enableAllEdge() {
	op.SetTopLeft(op.TopLeftAction())
	op.SetBottomLeft(op.BottomLeftAction())
	op.SetTopRight(op.TopRightAction())
	op.SetBottomRight(op.BottomRightAction())
}
