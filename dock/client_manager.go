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
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xwindow"
	"pkg.deepin.io/lib/dbus"
)

var (
	activeWindow    xproto.Window = 0
	isLauncherShown bool          = false
)

const (
	DDELauncher string = "dde-launcher"
)

// ClientManager用来管理启动程序相关窗口。
type ClientManager struct {
	// ActiveWindowChanged会在焦点窗口被改变时触发，会将最新的焦点窗口id发送给监听者。
	ActiveWindowChanged func(win uint32)
}

// NewClientManager creates a new client manager.
func NewClientManager() *ClientManager {
	return &ClientManager{}
}

// CurrentActiveWindow会返回当前焦点窗口的窗口id。
func (m *ClientManager) CurrentActiveWindow() uint32 {
	return uint32(activeWindow)
}

func getWindowUserTime(win xproto.Window) (uint, error) {
	timestamp, err := ewmh.WmUserTimeGet(XU, win)
	if err != nil {
		userTimeWindow, err := ewmh.WmUserTimeWindowGet(XU, win)
		if err != nil {
			return 0, err
		}
		timestamp, err = ewmh.WmUserTimeGet(XU, userTimeWindow)
		if err != nil {
			return 0, err
		}
	}
	return timestamp, nil
}

func changeCurrentWorkspaceToWindowWorkspace(win xproto.Window) error {
	winWorkspace, err := ewmh.WmDesktopGet(XU, win)
	if err != nil {
		return err
	}

	currentWorkspace, err := ewmh.CurrentDesktopGet(XU)
	if err != nil {
		return err
	}

	if currentWorkspace == winWorkspace {
		logger.Debugf("No need to change workspace, the current desktop is already %v", currentWorkspace)
		return nil
	}
	logger.Debug("Change workspace")

	winUserTime, err := getWindowUserTime(win)
	logger.Debug("window user time:", winUserTime)
	if err != nil {
		// only warning not return
		logger.Warning("getWindowUserTime failed:", err)
	}
	err = ewmh.CurrentDesktopReqExtra(XU, int(winWorkspace), xproto.Timestamp(winUserTime))
	if err != nil {
		return err
	}
	return nil
}

func activateWindow(win xproto.Window) error {
	err := changeCurrentWorkspaceToWindowWorkspace(win)
	if err != nil {
		return err
	}
	return ewmh.ActiveWindowReq(XU, win)
}

// ActiveWindow会激活给定id的窗口，被激活的窗口将通常会程序焦点窗口。(废弃，名字应该是ActivateWindow，当时手残打错了，此接口会在之后被移除，请使用正确的接口)
func (m *ClientManager) ActiveWindow(win uint32) bool {
	err := activateWindow(xproto.Window(win))
	if err != nil {
		logger.Warning("Activate window failed:", err)
		return false
	}
	return true
}

// ActivateWindow会激活给定id的窗口，被激活的窗口通常会成为焦点窗口。
func (m *ClientManager) ActivateWindow(win uint32) bool {
	err := activateWindow(xproto.Window(win))
	if err != nil {
		logger.Warning("Activate window failed:", err)
		return false
	}
	return true
}

// CloseWindow会将传入id的窗口关闭。
func (m *ClientManager) CloseWindow(win uint32) bool {
	err := ewmh.CloseWindow(XU, xproto.Window(win))
	if err != nil {
		logger.Warning("Close window failed:", err)
		return false
	}
	return true
}

// IsLauncherShown判断launcher是否已经显示。
func (m *ClientManager) IsLauncherShown() bool {
	return isLauncherShown
}

func isHiddenPre(win xproto.Window) bool {
	state, _ := ewmh.WmStateGet(XU, win)
	return contains(state, "_NET_WM_STATE_HIDDEN")
}

// works for new deepin wm.
func isWindowOnCurrentWorkspace(win xproto.Window) (bool, error) {
	winWorkspace, err := ewmh.WmDesktopGet(XU, win)
	if err != nil {
		return false, err
	}

	currentWorkspace, err := ewmh.CurrentDesktopGet(XU)
	if err != nil {
		return false, err
	}

	return winWorkspace == currentWorkspace, nil
}

func onCurrentWorkspacePre(win xproto.Window) bool {
	isOnCurrentWorkspace, err := isWindowOnCurrentWorkspace(win)
	if err != nil {
		logger.Warning(err)
		return false
	}
	return isOnCurrentWorkspace
}

func isDeepinLauncher(win xproto.Window) (bool, error) {
	winClass, err := icccm.WmClassGet(XU, win)
	if err != nil {
		return false, err
	}
	return winClass.Instance == DDELauncher, nil
}

func max(a, b int32) int32 {
	if a < b {
		return b
	}
	return a
}

func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func hasIntersection(x, y, w, h, x1, y1, w1, h1 int32) bool {
	ax := max(x, x1)
	ay := max(y, y1)
	bx := min(x+w, x1+w1)
	by := min(y+h, y1+h1)
	return ax <= bx && ay <= by
}

func isWindowDockOverlap(win xproto.Window) bool {
	// overlap 的条件： window showing + window on current workspace + window dock rect has intersection
	window := xwindow.New(XU, win)
	if isHiddenPre(win) || (!onCurrentWorkspacePre(win)) {
		return false
	}

	winRect, err := window.DecorGeometry()
	if err != nil {
		logger.Warningf("isWindowDockOverlap GetDecorGeometry of 0x%x failed: %s", win, err)
		return false
	}

	winX := int32(winRect.X())
	winY := int32(winRect.Y())
	winW := int32(winRect.Width())
	winH := int32(winRect.Height())

	dockX := int32(displayRect.X) + (int32(displayRect.Width)-dockProperty.PanelWidth)/2
	dockY := int32(displayRect.Y) + int32(displayRect.Height) - dockProperty.Height
	dockW := dockProperty.PanelWidth
	dockH := dockProperty.Height

	logger.Debugf("displayRect (%v,%v | %v x %v)",
		displayRect.X, displayRect.Y, displayRect.Width, displayRect.Height)
	logger.Debugf("dockProperty (%v x %v)", dockProperty.PanelWidth, dockProperty.Height)

	result := hasIntersection(winX, winY, winW, winH, dockX, dockY, dockW, dockH)
	logger.Debugf("window: (%v,%v | %v x %v), dock: (%v,%v | %v x %v), hasIntersection: %v",
		winX, winY, winW, winH, dockX, dockY, dockW, dockH, result)
	return result
}

func (m *ClientManager) handleClientListChanged() {
	clientList, err := ewmh.ClientListGet(XU)
	if err != nil {
		logger.Warning("Get client list failed:", err)
		return
	}
	ENTRY_MANAGER.runtimeAppChanged(clientList)
}

func (m *ClientManager) handleActiveWindowChanged() {
	logger.Debug("Active window changed")
	var err error
	isLauncherShown = false
	activeWindow, err = ewmh.ActiveWindowGet(XU)
	if err != nil {
		logger.Warning(err)
		return
	}
	// loop gets better performance than find_app_id_by_xid.
	// setLeader/updateState will filter invalid xid.
	for _, app := range ENTRY_MANAGER.runtimeApps {
		app.setLeader(activeWindow)
		app.updateState(activeWindow)
	}

	isLauncherShown, err = isDeepinLauncher(activeWindow)
	if err != nil {
		logger.Debug(err)
	}
	dbus.Emit(m, "ActiveWindowChanged", uint32(activeWindow))
	hideModemanager.updateStateWithDelay()
}

func (m *ClientManager) listenRootWindowPropertyChange() {
	rootWin := XU.RootWin()
	xwindow.New(XU, rootWin).Listen(xproto.EventMaskPropertyChange)
	xevent.PropertyNotifyFun(func(XU *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
		switch ev.Atom {
		case _NET_CLIENT_LIST:
			m.handleClientListChanged()
		case _NET_ACTIVE_WINDOW:
			m.handleActiveWindowChanged()
		}
	}).Connect(XU, rootWin)

	m.handleClientListChanged()
	xevent.Main(XU)
}
