/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dock

import (
// "github.com/BurntSushi/xgb/xproto"
// "github.com/BurntSushi/xgbutil/ewmh"
// "github.com/BurntSushi/xgbutil/xprop"
)

// return is mode changed
func (m *DockManager) setDisplayMode(mode DisplayModeType) bool {
	if mode == m.displayMode {
		return false
	}
	m.displayMode = mode

	// for _, rApp := range m.DockManager.runtimeApps {
	// 	rebuildXids := []xproto.Window{}
	// 	for xid, _ := range rApp.xids {
	// 		if _, err := xprop.PropValStr(
	// 			xprop.GetProperty(
	// 				XU,
	// 				xid,
	// 				"_DDE_DOCK_APP_ID",
	// 			),
	// 		); err != nil {
	// 			continue
	// 		}
	//
	// 		rebuildXids = append(rebuildXids, xid)
	// 		rApp.detachXid(xid)
	// 	}
	//
	// 	l := len(rebuildXids)
	// 	if l == 0 {
	// 		continue
	// 	}
	//
	// 	if len(rApp.xids) == 0 {
	// 		m.DockManager.destroyRuntimeApp(rApp)
	// 	}
	//
	// 	newApp := m.DockManager.createRuntimeApp(rebuildXids[0])
	// 	for i := 0; i < l; i++ {
	// 		newApp.attachXid(rebuildXids[i])
	// 	}
	//
	// 	activeXid, err := ewmh.ActiveWindowGet(XU)
	// 	if err != nil {
	// 		continue
	// 	}
	//
	// 	for xid, _ := range newApp.xids {
	// 		logger.Debugf("through new app xids")
	// 		if activeXid == xid {
	// 			logger.Debugf("0x%x(a), 0x%x(x)",
	// 				activeXid, xid)
	// 			newApp.setLeader(xid)
	// 			ewmh.ActiveWindowSet(XU, xid)
	// 			break
	// 		}
	// 	}
	// }

	m.dockHeight = getDockHeightByDisplayMode(mode)
	m.updateDockRect()
	return true
}

// return is mode changed
func (m *DockManager) setHideMode(mode HideModeType) bool {
	if mode == m.hideMode {
		return false
	}
	m.hideMode = mode
	m.hideStateManager.updateHideMode(mode)
	return true
}
