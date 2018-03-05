/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

func (m *Manager) registerWindow(win xproto.Window) {
	logger.Debug("register window", win)

	m.windowInfoMapMutex.RLock()
	winInfo, ok := m.windowInfoMap[win]
	m.windowInfoMapMutex.RUnlock()
	if ok {
		logger.Debugf("register window %v failed, window existed", win)
		m.attachOrDetachWindow(winInfo)
		return
	}

	winInfo = NewWindowInfo(win)
	m.listenWindowXEvent(winInfo)

	m.windowInfoMapMutex.Lock()
	m.windowInfoMap[win] = winInfo
	m.windowInfoMapMutex.Unlock()

	m.attachOrDetachWindow(winInfo)
}

func (m *Manager) isWindowRegistered(win xproto.Window) bool {
	m.windowInfoMapMutex.RLock()
	_, ok := m.windowInfoMap[win]
	m.windowInfoMapMutex.RUnlock()
	return ok
}

func (m *Manager) unregisterWindow(win xproto.Window) {
	logger.Debugf("unregister window %v", win)
	xevent.Detach(XU, win)
	m.windowInfoMapMutex.Lock()
	delete(m.windowInfoMap, win)
	m.windowInfoMapMutex.Unlock()
}

func (m *Manager) handleClientListChanged() {
	clientList, err := ewmh.ClientListGet(XU)
	if err != nil {
		logger.Warning("Get client list failed:", err)
		return
	}
	newClientList := windowSlice(clientList)
	sort.Sort(newClientList)
	add, remove := diffSortedWindowSlice(m.clientList, newClientList)
	m.clientList = newClientList

	if len(add) > 0 {
		logger.Debug("client list add:", add)
		for _, win := range add {
			m.registerWindow(win)
		}
	}

	if len(remove) > 0 {
		logger.Debug("client list remove:", remove)
		for _, win := range remove {

			m.windowInfoMapMutex.RLock()
			winInfo := m.windowInfoMap[win]
			m.windowInfoMapMutex.RUnlock()
			if winInfo != nil {
				m.detachWindow(winInfo)
			}
		}
	}
}

func (m *Manager) handleActiveWindowChanged() {
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

	m.Entries.mu.RLock()
	for _, entry := range m.Entries.items {
		entry.PropsMu.Lock()

		winInfo, ok := entry.windows[activeWindow]
		if ok {
			entry.setPropIsActive(true)
			entry.setCurrentWindowInfo(winInfo)
			entry.updateName()
			entry.updateIcon()
		} else {
			entry.setPropIsActive(false)
		}

		entry.PropsMu.Unlock()
	}
	m.Entries.mu.RUnlock()

	m.updateHideState(true)
}

func (m *Manager) listenRootWindowPropertyChange() {
	rootWin := XU.RootWin()
	xwindow.New(XU, rootWin).Listen(xproto.EventMaskPropertyChange | xproto.EventMaskSubstructureNotify)
	xevent.PropertyNotifyFun(func(XU *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
		switch ev.Atom {
		case atomNetClientList:
			m.handleClientListChanged()
		case atomNetActiveWindow:
			m.handleActiveWindowChanged()
		case atomNetShowingDesktop:
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

func (m *Manager) listenWindowXEvent(winInfo *WindowInfo) {
	win := winInfo.window
	logger.Debugf("start listen window %v x event", win)
	xwin := xwindow.New(XU, win)

	xwin.Listen(xproto.EventMaskPropertyChange | xproto.EventMaskStructureNotify | xproto.EventMaskVisibilityChange)

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
}

func (m *Manager) handlePropertyNotifyEvent(winInfo *WindowInfo, ev xevent.PropertyNotifyEvent) {
	switch ev.Atom {
	case atomNetWMState:
		winInfo.updateWmState()

	case atomNetWMName:
		winInfo.updateWmName()

	case atomNetWMIcon:
		winInfo.updateIcon()
	}

	entry := m.Entries.getByWindowId(ev.Window)
	if entry == nil {
		return
	}

	entry.PropsMu.Lock()
	defer entry.PropsMu.Unlock()

	switch ev.Atom {
	case atomNetWMState:
		entry.updateWindowInfos()

	case atomNetWMIcon:
		if entry.current == winInfo {
			entry.updateIcon()
		}

	case atomNetWMName:
		if entry.current == winInfo {
			entry.updateName()
		}
		entry.updateWindowInfos()
	}
}

func (m *Manager) handleConfigureNotifyEvent(winInfo *WindowInfo, ev xevent.ConfigureNotifyEvent) {
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

func (m *Manager) handleDestroyNotifyEvent(winInfo *WindowInfo, ev xevent.DestroyNotifyEvent) {
	logger.Debug(ev)
	m.unregisterWindow(winInfo.window)
	m.detachWindow(winInfo)
}
