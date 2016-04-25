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
	"github.com/BurntSushi/xgbutil/ewmh"
	"pkg.deepin.io/lib/dbus"
)

// ClientManager用来管理启动程序相关窗口。
type ClientManager struct {
	activeWindow xproto.Window
	// ActiveWindowChanged会在焦点窗口被改变时触发，会将最新的焦点窗口id发送给监听者。
	ActiveWindowChanged func(win uint32)
}

// NewClientManager creates a new client manager.
func NewClientManager() *ClientManager {
	m := new(ClientManager)
	var err error
	m.activeWindow, err = ewmh.ActiveWindowGet(XU)
	if err != nil {
		logger.Warning(err)
	}
	return m
}

// CurrentActiveWindow会返回当前焦点窗口的窗口id。
func (m *ClientManager) CurrentActiveWindow() uint32 {
	return uint32(m.activeWindow)
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
	return m.ActivateWindow(win)
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

func (m *ClientManager) updateActiveWindow(win xproto.Window) {
	dbus.Emit(m, "ActiveWindowChanged", uint32(m.activeWindow))
}
