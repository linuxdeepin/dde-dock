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
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/xrect"
	"github.com/BurntSushi/xgbutil/xwindow"
	"pkg.deepin.io/lib/dbus"
	"time"
)

const (
	TriggerShow int32 = iota
	TriggerHide
)

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

func (m *DockManager) isWindowDockOverlap(win xproto.Window) bool {
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

	dockWindow := xwindow.New(XU, m.frontendWindow)
	dockRect, err := dockWindow.DecorGeometry()
	if err != nil {
		logger.Warning(err)
		return false
	}

	logger.Debug("window rect:", winRect)
	logger.Debug("dock rect:", dockRect)
	result := hasIntersection(winRect, dockRect)
	logger.Debug("window dock overlap:", result)
	return result
}

const (
	DDELauncher = "dde-launcher"
)

func (m *DockManager) isDeepinLauncherShown() bool {
	winClass, err := icccm.WmClassGet(XU, m.ActiveWindow)
	if err != nil {
		logger.Debug(err)
		return false
	}
	return winClass.Instance == DDELauncher
}

func (m *DockManager) shouldHideOnSmartHideMode() bool {
	if m.isDeepinLauncherShown() {
		logger.Debug("launcher is shown")
		return false
	}
	return m.isWindowDockOverlap(m.ActiveWindow)
}

func (m *DockManager) smartHideModeTimerExpired() {
	logger.Debug("smartHideModeTimer expired!")
	if m.shouldHideOnSmartHideMode() {
		m.emitSignalChangeHideState(TriggerHide)
	} else {
		m.emitSignalChangeHideState(TriggerShow)
	}
}

func (m *DockManager) resetSmartHideModeTimer(delay time.Duration) {
	m.smartHideModeMutex.Lock()
	defer m.smartHideModeMutex.Unlock()

	m.smartHideModeTimer.Reset(delay)
	logger.Debug("reset smart hide mode timer ", delay)
}

func (m *DockManager) cancelSmartHideModeTimer() {
	m.smartHideModeMutex.Lock()
	defer m.smartHideModeMutex.Unlock()

	m.smartHideModeTimer.Stop()
	logger.Debug("cancel smart hide mode timer ")
}

func (m *DockManager) smartHideModeDelayHandle() {
	hideState := m.HideState.Get()
	switch hideState {
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

func (m *DockManager) updateHideState(delay bool) {
	if m.isDeepinLauncherShown() {
		logger.Debug("updateHideState: launcher is shown, show dock")
		m.emitSignalChangeHideState(TriggerShow)
		return
	}

	hideMode := HideModeType(m.HideMode.Get())
	logger.Debug("updateHideState: mode is", hideMode)
	switch hideMode {
	case HideModeKeepShowing:
		m.emitSignalChangeHideState(TriggerShow)

	case HideModeKeepHidden:
		m.emitSignalChangeHideState(TriggerHide)

	case HideModeSmartHide:
		if delay {
			m.smartHideModeDelayHandle()
		} else {
			m.smartHideModeTimer.Reset(0)
		}
	}
}

func (m *DockManager) updateHideStateWithDelay() {
	m.updateHideState(true)
}

func (m *DockManager) updateHideStateWithoutDelay() {
	m.updateHideState(false)
}

func (m *DockManager) emitSignalChangeHideState(trigger int32) {
	hideState := m.HideState.Get()
	if (hideState == HideStateShown && trigger == TriggerShow) ||
		(hideState == HideStateHidden && trigger == TriggerHide) {
		logger.Debug("No need emit signal ChangeState")
		return
	}

	triggerStr := "TriggerUnknown"
	if trigger == TriggerShow {
		triggerStr = "TriggerShow"
	} else if trigger == TriggerHide {
		triggerStr = "TriggerHide"
	}
	logger.Debugf("Emit signal ChangeHideState: %v", triggerStr)
	dbus.Emit(m, "ChangeHideState", trigger)
}
