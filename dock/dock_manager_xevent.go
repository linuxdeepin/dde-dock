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
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xwindow"
)

func (m *DockManager) handleClientListChanged() {
	clientList, err := ewmh.ClientListGet(XU)
	if err != nil {
		logger.Warning("Get client list failed:", err)
		return
	}
	m.entryManager.updateClientList(clientList)
}

func (m *DockManager) handleActiveWindowChanged() {
	logger.Debug("Active window changed")
	var err error
	m.activeWindow, err = ewmh.ActiveWindowGet(XU)
	if err != nil {
		logger.Warning(err)
		return
	}

	m.entryManager.updateActiveWindow(m.activeWindow)
	m.clientManager.updateActiveWindow(m.activeWindow)
	m.hideStateManager.updateActiveWindow(m.activeWindow)
}

func (m *DockManager) listenRootWindowPropertyChange() {
	rootWin := XU.RootWin()
	xwindow.New(XU, rootWin).Listen(xproto.EventMaskPropertyChange)
	xevent.PropertyNotifyFun(func(XU *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
		switch ev.Atom {
		case _NET_CLIENT_LIST:
			m.handleClientListChanged()
		case _NET_ACTIVE_WINDOW:
			m.handleActiveWindowChanged()
		case _NET_SHOWING_DESKTOP:
			m.hideStateManager.updateStateWithoutDelay()
		}
	}).Connect(XU, rootWin)

	m.handleActiveWindowChanged()
	m.handleClientListChanged()
}
