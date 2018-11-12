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

	x "github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/util/wm/ewmh"
)

func (m *Manager) registerWindow(win x.Window) {
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

func (m *Manager) isWindowRegistered(win x.Window) bool {
	m.windowInfoMapMutex.RLock()
	_, ok := m.windowInfoMap[win]
	m.windowInfoMapMutex.RUnlock()
	return ok
}

func (m *Manager) unregisterWindow(win x.Window) {
	logger.Debugf("unregister window %v", win)
	m.windowInfoMapMutex.Lock()
	delete(m.windowInfoMap, win)
	m.windowInfoMapMutex.Unlock()
}

func (m *Manager) handleClientListChanged() {
	clientList, err := ewmh.GetClientList(globalXConn).Reply(globalXConn)
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
	activeWindow, err := ewmh.GetActiveWindow(globalXConn).Reply(globalXConn)
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

func (m *Manager) listenRootWindowXEvent() {
	const eventMask = x.EventMaskPropertyChange | x.EventMaskSubstructureNotify
	err := x.ChangeWindowAttributesChecked(globalXConn, m.rootWindow, x.CWEventMask,
		[]uint32{eventMask}).Check(globalXConn)
	if err != nil {
		logger.Warning(err)
	}
	m.handleActiveWindowChanged()
	m.handleClientListChanged()
}

func (m *Manager) listenWindowXEvent(winInfo *WindowInfo) {
	const eventMask = x.EventMaskPropertyChange | x.EventMaskStructureNotify | x.EventMaskVisibilityChange
	err := x.ChangeWindowAttributesChecked(globalXConn, winInfo.window, x.CWEventMask,
		[]uint32{eventMask}).Check(globalXConn)
	if err != nil {
		logger.Warning(err)
	}
}

func (m *Manager) handleDestroyNotifyEvent(ev *x.DestroyNotifyEvent) {
	logger.Debug("DestroyNotifyEvent window:", ev.Window)
	winInfo := m.getWindowInfo(ev.Window)
	if winInfo != nil {
		m.detachWindow(winInfo)
	}
	m.unregisterWindow(ev.Window)
}

func (m *Manager) handleMapNotifyEvent(ev *x.MapNotifyEvent) {
	logger.Debug("MapNotifyEvent window:", ev.Window)
	m.registerWindow(ev.Window)
}

func (m *Manager) getWindowInfo(win x.Window) *WindowInfo {
	m.windowInfoMapMutex.RLock()
	v := m.windowInfoMap[win]
	m.windowInfoMapMutex.RUnlock()
	return v
}

func (m *Manager) handleConfigureNotifyEvent(ev *x.ConfigureNotifyEvent) {
	winInfo := m.getWindowInfo(ev.Window)
	if winInfo == nil {
		return
	}

	if HideModeType(m.HideMode.Get()) != HideModeSmartHide {
		return
	}
	if winInfo.wmClass != nil && winInfo.wmClass.Class == frontendWindowWmClass {
		// ignore frontend window ConfigureNotify event
		return
	}

	winInfo.mu.Lock()
	winInfo.lastConfigureNotifyEvent = ev
	winInfo.mu.Unlock()

	const configureNotifyDelay = 100 * time.Millisecond
	if winInfo.updateConfigureTimer != nil {
		winInfo.updateConfigureTimer.Reset(configureNotifyDelay)
	} else {
		winInfo.updateConfigureTimer = time.AfterFunc(configureNotifyDelay, func() {
			logger.Debug("ConfigureNotify: updateConfigureTimer expired")

			winInfo.mu.Lock()
			ev := winInfo.lastConfigureNotifyEvent
			winInfo.mu.Unlock()

			logger.Debugf("in closure: configure notify ev: %#v", ev)
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

func (m *Manager) handleRootWindowPropertyNotifyEvent(ev *x.PropertyNotifyEvent) {
	switch ev.Atom {
	case atomNetClientList:
		m.handleClientListChanged()
	case atomNetActiveWindow:
		m.handleActiveWindowChanged()
	case atomNetShowingDesktop:
		m.updateHideState(false)
	}
}

func (m *Manager) handlePropertyNotifyEvent(ev *x.PropertyNotifyEvent) {
	if ev.Window == m.rootWindow {
		m.handleRootWindowPropertyNotifyEvent(ev)
		return
	}

	winInfo := m.getWindowInfo(ev.Window)
	if winInfo == nil {
		return
	}

	atomName, err := globalXConn.GetAtomName(ev.Atom)
	if err != nil {
		logger.Warning(err)
	}
	logger.Debugf("winInfo %d property %s changed", winInfo.window, atomName)

	switch ev.Atom {
	case atomNetWMState:
		winInfo.updateWmState()
		m.attachOrDetachWindow(winInfo)

	case atomNetWMName:
		winInfo.updateWmName()

	case atomNetWMIcon:
		winInfo.updateIcon()

	case atomNetWmAllowedActions:
		winInfo.updateWmAllowedActions()

	case x.AtomWMClass:
		winInfo.updateWmClass()
		m.attachOrDetachWindow(winInfo)

	case atomXEmbedInfo:
		winInfo.updateHasXEmbedInfo()
		m.attachOrDetachWindow(winInfo)

	case atomNetWMWindowType:
		winInfo.updateWmWindowType()
		m.attachOrDetachWindow(winInfo)

	case x.AtomWMTransientFor:
		winInfo.updateHasWmTransientFor()
		m.attachOrDetachWindow(winInfo)
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

	case atomNetWmAllowedActions:
		entry.updateMenu()
	}
}

func (m *Manager) eventHandleLoop() {
	eventChan := make(chan x.GenericEvent, 500)
	globalXConn.AddEventChan(eventChan)

	for ev := range eventChan {
		switch ev.GetEventCode() {
		case x.MapNotifyEventCode:
			event, _ := x.NewMapNotifyEvent(ev)
			m.handleMapNotifyEvent(event)

		case x.DestroyNotifyEventCode:
			event, _ := x.NewDestroyNotifyEvent(ev)
			m.handleDestroyNotifyEvent(event)

		case x.ConfigureNotifyEventCode:
			event, _ := x.NewConfigureNotifyEvent(ev)
			m.handleConfigureNotifyEvent(event)

		case x.PropertyNotifyEventCode:
			event, _ := x.NewPropertyNotifyEvent(ev)
			m.handlePropertyNotifyEvent(event)
		}
	}
}
