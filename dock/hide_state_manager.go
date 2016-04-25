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
	"errors"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/xrect"
	"github.com/BurntSushi/xgbutil/xwindow"
	"pkg.deepin.io/lib/dbus"
	"sync"
	"time"
)

const (
	TriggerShow int32 = iota
	TriggerHide
)

type HideStateManager struct {
	state        HideStateType
	ChangeState  func(int32)
	mode         HideModeType
	activeWindow xproto.Window
	dockRect     *xrect.XRect

	smartHideModeTimer *time.Timer
	smartHideModeMutex sync.Mutex
}

func NewHideStateManager() *HideStateManager {
	m := &HideStateManager{}
	m.smartHideModeTimer = time.AfterFunc(10*time.Second, m.smartHideModeTimerExpired)
	m.smartHideModeTimer.Stop()
	return m
}

func getDeepinDock() (xproto.Window, error) {
	clientList, err := ewmh.ClientListGet(XU)
	if err != nil {
		return 0, err
	}
	for _, win := range clientList {
		winClass, err := icccm.WmClassGet(XU, win)
		if err != nil {
			logger.Debug(err)
			continue
		}
		if winClass.Instance == "dde-dock" {
			return win, nil
		}
	}
	return 0, errors.New("not found deepin dock")
}

func (m *HideStateManager) getHideStateByMode() HideStateType {
	switch m.mode {
	case HideModeKeepShowing, HideModeSmartHide:
		return HideStateShown
	case HideModeKeepHidden:
		return HideStateHidden
	}
	return HideStateShown
}

func (m *HideStateManager) initHideState() {
	deepinDockWin, err := getDeepinDock()
	if err != nil {
		logger.Debug(err)
		m.state = m.getHideStateByMode()
		return
	}

	// get dock window rect
	window := xwindow.New(XU, deepinDockWin)
	dockRect, err := window.Geometry()
	if err != nil {
		logger.Warning(err)
		m.state = m.getHideStateByMode()
		return
	}
	logger.Debug("initHideState: dock window rect", dockRect)

	if dockRect.Height() > 1 {
		m.state = HideStateShown
	} else {
		m.state = HideStateHidden
	}
	logger.Debug("initHideState state", m.state)
}

func (m *HideStateManager) destroy() {
	if m.smartHideModeTimer != nil {
		m.smartHideModeTimer.Stop()
		m.smartHideModeTimer = nil
	}
	dbus.UnInstallObject(m)
}

func (m *HideStateManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.daemon.Dock",
		ObjectPath: "/dde/dock/HideStateManager",
		Interface:  "dde.dock.HideStateManager",
	}
}

func (m *HideStateManager) SetState(s int32) int32 {
	newState := HideStateType(s)
	if m.state != newState {
		logger.Debugf("SetState: %v => %v", m.state, newState)
		m.state = newState
	}
	return s
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func hasIntersection(rectA, rectB xrect.Rect) bool {
	if rectA == nil || rectB == nil {
		logger.Warning("hasIntersection rectA or rectB is nil")
		return false
	}
	x, y, w, h := rectA.Pieces()
	x1, y1, w1, h1 := rectB.Pieces()
	ax := max(x, x1)
	ay := max(y, y1)
	bx := min(x+w, x1+w1)
	by := min(y+h, y1+h1)
	return ax <= bx && ay <= by
}

func (m *HideStateManager) isWindowDockOverlap(win xproto.Window) bool {
	// overlap condition:  window showing and  on current workspace,
	// window dock rect has intersection
	window := xwindow.New(XU, win)
	if isHiddenPre(win) || (!onCurrentWorkspacePre(win)) {
		return false
	}

	winRect, err := window.DecorGeometry()
	if err != nil {
		logger.Warningf("isWindowDockOverlap GetDecorGeometry of 0x%x failed: %s", win, err)
		return false
	}

	logger.Debug("window rect:", winRect)
	logger.Debug("dock rect:", m.dockRect)
	result := hasIntersection(winRect, m.dockRect)
	logger.Debug("window dock overlap:", result)
	return result
}

const (
	DDELauncher = "dde-launcher"
)

func (m *HideStateManager) isDeepinLauncherShown() bool {
	winClass, err := icccm.WmClassGet(XU, m.activeWindow)
	if err != nil {
		logger.Debug(err)
		return false
	}
	return winClass.Instance == DDELauncher
}

func (m *HideStateManager) shouldHideOnSmartHideMode() bool {
	if m.isDeepinLauncherShown() {
		logger.Debug("launcher is shown")
		return false
	}
	return m.isWindowDockOverlap(m.activeWindow)
}

func (m *HideStateManager) smartHideModeTimerExpired() {
	logger.Debug("smartHideModeTimer expired!")
	if m.shouldHideOnSmartHideMode() {
		m.emitSignalChangeState(TriggerHide)
	} else {
		m.emitSignalChangeState(TriggerShow)
	}
}

func (m *HideStateManager) resetSmartHideModeTimer(delay time.Duration) {
	m.smartHideModeMutex.Lock()
	defer m.smartHideModeMutex.Unlock()

	m.smartHideModeTimer.Reset(delay)
	logger.Debug("reset smart hide mode timer ", delay)
}

func (m *HideStateManager) cancelSmartHideModeTimer() {
	m.smartHideModeMutex.Lock()
	defer m.smartHideModeMutex.Unlock()

	m.smartHideModeTimer.Stop()
	logger.Debug("cancel smart hide mode timer ")
}

func (m *HideStateManager) smartHideModeDelayHandle() {
	switch m.state {
	case HideStateShown:
		if m.shouldHideOnSmartHideMode() {
			logger.Debug("smartHideModeDelayHandle: show -> hide")
			m.resetSmartHideModeTimer(time.Millisecond * 500)
		} else {
			logger.Debug("smartHideModeDelayHandle: show -> show")
			m.cancelSmartHideModeTimer()
		}

	case HideStateHidden:
		if m.shouldHideOnSmartHideMode() {
			logger.Debug("smartHideModeDelayHandle: hide -> hide")
			m.cancelSmartHideModeTimer()
		} else {
			logger.Debug("smartHideModeDelayHandle: hide -> show")
			m.resetSmartHideModeTimer(time.Millisecond * 500)
		}
	}
}

func (m *HideStateManager) updateHideMode(mode HideModeType) {
	m.mode = mode
	m.updateStateWithoutDelay()
}

func (m *HideStateManager) updateActiveWindow(win xproto.Window) {
	m.activeWindow = win
	// 切换窗口时需要延时
	m.updateStateWithDelay()
}

func (m *HideStateManager) updateState(delay bool) {
	if m.isDeepinLauncherShown() {
		logger.Debug("updateState: launcher is shown, show dock")
		m.emitSignalChangeState(TriggerShow)
		return
	}

	logger.Debug("updateState: mode is", m.mode)
	switch m.mode {
	case HideModeKeepShowing:
		m.emitSignalChangeState(TriggerShow)

	case HideModeKeepHidden:
		m.emitSignalChangeState(TriggerHide)

	case HideModeSmartHide:
		if delay {
			m.smartHideModeDelayHandle()
		} else {
			m.smartHideModeTimer.Reset(0)
		}
	}
}

func (m *HideStateManager) updateStateWithDelay() {
	m.updateState(true)
}

func (m *HideStateManager) updateStateWithoutDelay() {
	m.updateState(false)
}

func (m *HideStateManager) UpdateState() {
	logger.Debug("dbus call UpdateState")
	m.updateState(true)
}

func (m *HideStateManager) emitSignalChangeState(trigger int32) {
	if (m.state == HideStateShown && trigger == TriggerShow) ||
		(m.state == HideStateHidden && trigger == TriggerHide) {
		logger.Debug("No need emit signal ChangeState")
		return
	}

	triggerStr := "TriggerUnknown"
	if trigger == TriggerShow {
		triggerStr = "TriggerShow"
	} else if trigger == TriggerHide {
		triggerStr = "TriggerHide"
	}
	logger.Debugf("Emit signal ChangeState: %v", triggerStr)
	dbus.Emit(m, "ChangeState", trigger)
}
