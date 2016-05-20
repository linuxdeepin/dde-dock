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
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xwindow"
	"sort"
	"time"
)

func (m *EntryManager) updateClientList(clientList []xproto.Window) {
	newClientList := windowSlice(clientList)
	sort.Sort(newClientList)

	add, remove := diffSortedWindowSlice(m.clientList, newClientList)
	if len(add) > 0 {
		logger.Debug("client list add:", add)
		for _, win := range add {
			winInfo := NewWindowInfo(win)
			m.listenWindowXEvent(winInfo)
		}
	}

	if len(remove) > 0 {
		logger.Debug("client list remove:", remove)
	}
	m.clientList = newClientList
}

func (m *EntryManager) listenWindowXEvent(winInfo *WindowInfo) {
	win := winInfo.window
	logger.Debugf("start listen window %v x event", win)
	xwin := xwindow.New(XU, win)

	xwin.Listen(xproto.EventMaskPropertyChange | xproto.EventMaskStructureNotify | xproto.EventMaskVisibilityChange)

	// need listen EventMaskVisibilityChange
	xevent.VisibilityNotifyFun(func(XU *xgbutil.XUtil, ev xevent.VisibilityNotifyEvent) {
		m.handleVisibilityNotifyEvent(winInfo, ev)
	}).Connect(XU, win)

	winInfo.initPropertyNotifyEventHandler(m)
	// need listen EventMaskPropertyChange
	xevent.PropertyNotifyFun(func(XU *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
		m.handlePropertyNotifyEvent(winInfo, ev)
	}).Connect(XU, win)

	// need listen EventMaskStructureNotify
	// move resize minimized Maximize window
	xevent.ConfigureNotifyFun(func(XU *xgbutil.XUtil, ev xevent.ConfigureNotifyEvent) {
		m.handleConfigureNotifyEvent(winInfo, ev)
	}).Connect(XU, win)

	xevent.DestroyNotifyFun(func(XU *xgbutil.XUtil, ev xevent.DestroyNotifyEvent) {
		m.handleDestroyNotifyEvent(winInfo, ev)
	}).Connect(XU, win)

	xevent.UnmapNotifyFun(func(XU *xgbutil.XUtil, ev xevent.UnmapNotifyEvent) {
		m.handleUnmapNotifyEvent(winInfo, ev)
	}).Connect(XU, win)
}

func (m *EntryManager) handleVisibilityNotifyEvent(winInfo *WindowInfo, ev xevent.VisibilityNotifyEvent) {
	logger.Debug(ev)
	winInfo.updateMapState()
}

func (m *EntryManager) handlePropertyNotifyEvent(winInfo *WindowInfo, ev xevent.PropertyNotifyEvent) {
	winInfo.propertyNotifyAtomTable[ev.Atom] = true
	if winInfo.propertyNotifyEnabled {
		winInfo.propertyNotifyTimer.Reset(300 * time.Millisecond)
		winInfo.propertyNotifyEnabled = false
	}
}

func (m *EntryManager) handleConfigureNotifyEvent(winInfo *WindowInfo, ev xevent.ConfigureNotifyEvent) {
	if dockManager.hideMode != HideModeSmartHide {
		return
	}
	winInfo.lastConfigureNotifyEvent = &ev
	const configureNotifyDelay = 100 // ms
	if winInfo.updateConfigureTimer != nil {
		winInfo.updateConfigureTimer.Reset(time.Millisecond * configureNotifyDelay)
	} else {
		winInfo.updateConfigureTimer = time.AfterFunc(time.Millisecond*configureNotifyDelay, func() {
			logger.Debug("ConfigureNotify: updateConfigureTimer expired")
			ev := winInfo.lastConfigureNotifyEvent
			logger.Debugf("in closure: configure notify ev %s", ev)
			isXYWHChange := false
			if winInfo.x != ev.X {
				winInfo.x = ev.X
				isXYWHChange = true
			}

			if winInfo.y != ev.Y {
				winInfo.y = ev.Y
				isXYWHChange = true
			}

			if winInfo.width != ev.Width {
				winInfo.width = ev.Width
				isXYWHChange = true
			}

			if winInfo.height != ev.Height {
				winInfo.height = ev.Height
				isXYWHChange = true
			}
			logger.Debug("isXYWHChange", isXYWHChange)
			if isXYWHChange {
				dockManager.hideStateManager.updateStateWithoutDelay()
			} else {
				dockManager.hideStateManager.updateStateWithDelay()
			}
		})
	}
}

func (m *EntryManager) handleDestroyNotifyEvent(winInfo *WindowInfo, ev xevent.DestroyNotifyEvent) {
	logger.Debug(ev)
	xevent.Detach(XU, winInfo.window)
	m.detachRuntimeAppWindow(winInfo)
}

func (m *EntryManager) handleUnmapNotifyEvent(winInfo *WindowInfo, ev xevent.UnmapNotifyEvent) {
	logger.Debug(ev)
	xevent.Detach(XU, winInfo.window)
	m.detachRuntimeAppWindow(winInfo)
}
