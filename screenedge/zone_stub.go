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
	"pkg.deepin.io/lib/dbus"
)

type Manager struct {
	lTopTimer    *edgeTimer
	lBottomTimer *edgeTimer
	rTopTimer    *edgeTimer
	rBottomTimer *edgeTimer
}

const (
	ZONE_DEST = "com.deepin.daemon.Zone"
	ZONE_PATH = "/com/deepin/daemon/Zone"
	ZONE_IFC  = "com.deepin.daemon.Zone"
)

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       ZONE_DEST,
		ObjectPath: ZONE_PATH,
		Interface:  ZONE_IFC,
	}
}

func (m *Manager) listenSignal() {
	dspObj.ConnectPrimaryChanged(func(argv []interface{}) {
		unregisterZoneArea()
		registerZoneArea()
	})

	areaObj.ConnectCursorInto(func(x, y int32, id string) {
		m.handleCursorSignal(x, y, id, true)
	})

	areaObj.ConnectCursorOut(func(x, y int32, id string) {
		m.handleCursorSignal(x, y, id, false)
	})

	areaObj.ConnectCancelAllArea(func() {
		unregisterZoneArea()
		registerZoneArea()
	})

	launchObj.ConnectShown(func() {
		enableSpecialEdges(getEdgeForCommand("dde-launcher"))
	})

	launchObj.ConnectClosed(func() {
		m.enableAllEdge()
	})
}

func (m *Manager) enableAllEdge() {
	m.initEdgeActionMap()
}

func (m *Manager) filterCursorSignal(id string) bool {
	if id != areaId {
		return true
	}

	if isAppInBlackList() {
		return true
	}

	if pid, ok := isActiveWindowFullscreen(); ok {
		if !isAppInWhiteList(pid) {
			return true
		}
	}

	return false
}

func (m *Manager) handleCursorSignal(x, y int32, id string, into bool) {
	if m.filterCursorSignal(id) {
		return
	}

	if !into {
		m.lTopTimer.StopTimer()
		m.lBottomTimer.StopTimer()
		m.rTopTimer.StopTimer()
		m.rBottomTimer.StopTimer()
		return
	}

	setting := zoneSettings()
	delay := setting.GetInt("delay")
	if isInArea(x, y, topLeftArea) {
		m.lTopTimer.DoAction(leftTopEdge, delay)
	} else if isInArea(x, y, bottomLeftArea) {
		m.lBottomTimer.DoAction(leftBottomEdge, delay)
	} else if isInArea(x, y, topRightArea) {
		m.rTopTimer.DoAction(rightTopEdge, delay)
	} else if isInArea(x, y, bottomRightArea) {
		m.rBottomTimer.DoAction(rightBottomEdge, delay)
	}
}
