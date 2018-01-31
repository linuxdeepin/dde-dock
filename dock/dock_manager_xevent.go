/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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

package dock

import (
	"sort"
	"time"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xwindow"
)

func (m *DockManager) registerWindow(win xproto.Window) {

	logger.Debug("register window", win)
	registered := m.isWindowRegistered(win)
	if registered {
		logger.Debugf("register window %v failed, window existed", win)
		return
	}

	winInfo := NewWindowInfo(win)
	m.listenWindowXEvent(winInfo)

	m.windowInfoMapMutex.Lock()
	m.windowInfoMap[win] = winInfo
	m.windowInfoMapMutex.Unlock()
}

func (m *DockManager) isWindowRegistered(win xproto.Window) bool {
	m.windowInfoMapMutex.RLock()
	_, ok := m.windowInfoMap[win]
	m.windowInfoMapMutex.RUnlock()
	return ok
}

func (m *DockManager) unregisterWindow(win xproto.Window) {
	logger.Debugf("unregister window %v", win)
	xevent.Detach(XU, win)
	m.windowInfoMapMutex.Lock()
	delete(m.windowInfoMap, win)
	m.windowInfoMapMutex.Unlock()
}

func (m *DockManager) handleClientListChanged() {
	clientList, err := ewmh.ClientListGet(XU)
	if err != nil {
		logger.Warning("Get client list failed:", err)
		return
	}
	newClientList := windowSlice(clientList)
	sort.Sort(newClientList)

	add, remove := diffSortedWindowSlice(m.clientList, newClientList)
	if len(add) > 0 {
		logger.Debug("client list add:", add)
		for _, win := range add {
			m.registerWindow(win)
		}
	}

	if len(remove) > 0 {
		logger.Debug("client list remove:", remove)
	}
	m.clientList = newClientList
}

func (m *DockManager) handleActiveWindowChanged() {
	activeWindow, err := ewmh.ActiveWindowGet(XU)
	if err != nil {
		logger.Warning(err)
		return
	}
	m.activeWindowMu.Lock()
	if m.activeWindow == activeWindow {
		m.activeWindowMu.Unlock()
		logger.Debug("Active window no change")
		return
	}
	// try handle activeWindow == 0
	m.activeWindowOld = m.activeWindow
	m.activeWindow = activeWindow
	m.activeWindowMu.Unlock()

	logger.Debug("Active window changed", activeWindow)

	for _, entry := range m.Entries {
		winInfo, ok := entry.windows[activeWindow]
		if ok {
			entry.setIsActive(true)
			entry.setCurrentWindowInfo(winInfo)
			entry.current.updateWmName()
			entry.updateIcon()
		} else {
			entry.setIsActive(false)
		}
	}

	m.updateHideState(true)
}

func (m *DockManager) listenRootWindowPropertyChange() {
	rootWin := XU.RootWin()
	xwindow.New(XU, rootWin).Listen(xproto.EventMaskPropertyChange | xproto.EventMaskSubstructureNotify)
	xevent.PropertyNotifyFun(func(XU *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
		switch ev.Atom {
		case _NET_CLIENT_LIST:
			m.handleClientListChanged()
		case _NET_ACTIVE_WINDOW:
			m.handleActiveWindowChanged()
		case _NET_SHOWING_DESKTOP:
			m.updateHideState(false)
		}
	}).Connect(XU, rootWin)

	xevent.MapNotifyFun(func(XU *xgbutil.XUtil, ev xevent.MapNotifyEvent) {
		win := ev.Window
		logger.Debugf("rootWin MapNotifyEvent window: %v", win)

		m.registerWindow(win)
	}).Connect(XU, rootWin)

	m.handleActiveWindowChanged()
	m.handleClientListChanged()
}

func (m *DockManager) listenWindowXEvent(winInfo *WindowInfo) {
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
		winInfo.handlePropertyNotifyEvent(ev)
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

func (m *DockManager) handleVisibilityNotifyEvent(winInfo *WindowInfo, ev xevent.VisibilityNotifyEvent) {
	logger.Debug(ev)
	winInfo.updateMapState()
}

func (m *DockManager) handleConfigureNotifyEvent(winInfo *WindowInfo, ev xevent.ConfigureNotifyEvent) {
	if HideModeType(m.HideMode.Get()) != HideModeSmartHide {
		return
	}
	if winInfo.wmClass != nil && winInfo.wmClass.Class == frontendWindowWmClass {
		// ignore frontend window ConfigureNotify event
		return
	}

	winInfo.lastConfigureNotifyEvent = &ev
	const configureNotifyDelay = 100 * time.Millisecond
	if winInfo.updateConfigureTimer != nil {
		winInfo.updateConfigureTimer.Reset(configureNotifyDelay)
	} else {
		winInfo.updateConfigureTimer = time.AfterFunc(configureNotifyDelay, func() {
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
			// if xywh changed ,update hide state without delay
			m.updateHideState(!isXYWHChange)
		})
	}
}

func (m *DockManager) handleDestroyNotifyEvent(winInfo *WindowInfo, ev xevent.DestroyNotifyEvent) {
	logger.Debug(ev)
	m.unregisterWindow(winInfo.window)
	m.detachWindow(winInfo)
}

func (m *DockManager) handleUnmapNotifyEvent(winInfo *WindowInfo, ev xevent.UnmapNotifyEvent) {
	logger.Debug(ev)
	winInfo.updateMapState()
}
