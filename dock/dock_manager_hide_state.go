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
	"pkg.deepin.io/lib/dbus"
	"time"
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
	return ax < bx && ay < by
}

func (m *DockManager) getActiveWinGroup() (ret []xproto.Window) {
	ret = []xproto.Window{m.activeWindow}
	list, err := ewmh.ClientListStackingGet(XU)
	if err != nil {
		logger.Warning(err)
		return
	}

	idx := -1
	for i, win := range list {
		if win == m.activeWindow {
			idx = i
			break
		}
	}
	if idx == -1 {
		logger.Warning("getActiveWinGroup: not found active window in clientListStacking")
		return
	} else if idx == 0 {
		return
	}

	aPid := getWmPid(XU, m.activeWindow)
	aWmClass, _ := icccm.WmClassGet(XU, m.activeWindow)
	aLeaderWin, _ := getWmClientLeader(XU, m.activeWindow)

	for i := idx - 1; i >= 0; i-- {
		win := list[i]
		pid := getWmPid(XU, win)
		if aPid != 0 && pid == aPid {
			// ok
			ret = append(ret, win)
			continue
		}

		wmClass, _ := icccm.WmClassGet(XU, win)
		if wmClass != nil && wmClass.Class == frontendWindowWmClass {
			// skip over frontend window
			continue
		}

		if wmClass != nil && aWmClass != nil &&
			wmClass.Class == aWmClass.Class {
			// ok
			ret = append(ret, win)
			continue
		}

		leaderWin, _ := getWmClientLeader(XU, win)
		if aLeaderWin != 0 && leaderWin == aLeaderWin {
			// ok
			ret = append(ret, win)
			continue
		}

		aboveWin := list[i+1]
		aboveWinTransientFor, _ := getWmTransientFor(XU, aboveWin)
		if aboveWinTransientFor != 0 && aboveWinTransientFor == win {
			// ok
			ret = append(ret, win)
			continue
		}

		break
	}
	return
}

func (m *DockManager) isWindowDockOverlap(win xproto.Window) (bool, error) {
	// overlap condition:
	// window type is not desktop
	// window opacity is not zero
	// window showing and  on current workspace,
	// window dock rect has intersection
	windowType, err := ewmh.WmWindowTypeGet(XU, win)
	if err == nil && strSliceContains(windowType, "_NET_WM_WINDOW_TYPE_DESKTOP") {
		return false, nil
	}

	opacity, err := getWmWindowOpacity(XU, win)
	if err == nil && opacity == 0 {
		return false, nil
	}

	if isHiddenPre(win) || (!onCurrentWorkspacePre(win)) {
		logger.Debugf("window %v is hidden or not on current workspace", win)
		return false, nil
	}

	winRect, err := getWindowGeometry(XU, win)
	if err != nil {
		logger.Warning("Get target window geometry failed", err)
		return false, err
	}

	logger.Debug("window rect:", winRect)
	logger.Debug("dock rect:", m.FrontendWindowRect)
	return hasIntersection(winRect, m.FrontendWindowRect.ToXRect()), nil
}

const (
	DDELauncher = "dde-launcher"
)

func (m *DockManager) isDeepinLauncherShown() bool {
	winClass, err := icccm.WmClassGet(XU, m.activeWindow)
	if err != nil {
		logger.Debug(err)
		return false
	}
	return winClass.Instance == DDELauncher
}

func (m *DockManager) shouldHideOnSmartHideMode() (bool, error) {
	if m.activeWindow == 0 {
		logger.Debug("shouldHideOnSmartHideMode: activeWindow is 0")
		return false, errors.New("activeWindow is 0")
	}
	if m.isDeepinLauncherShown() {
		logger.Debug("launcher is shown")
		return false, nil
	}
	list := m.getActiveWinGroup()
	logger.Debug("activeWinGroup:", list)
	for _, win := range list {
		over, err := m.isWindowDockOverlap(win)
		if err != nil {
			return false, err
		}
		logger.Debugf("win %d dock overlap %v", win, over)
		if over {
			return true, nil
		}
	}
	return false, nil
}

func (m *DockManager) smartHideModeTimerExpired() {
	logger.Debug("smartHideModeTimer expired!")
	shouldHide, err := m.shouldHideOnSmartHideMode()
	if err != nil {
		logger.Warning(err)
		m.setPropHideState(HideStateUnknown)
		return
	}

	if shouldHide {
		m.setPropHideState(HideStateHide)
	} else {
		m.setPropHideState(HideStateShow)
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

func (m *DockManager) updateHideState(delay bool) {
	if m.isDeepinLauncherShown() {
		logger.Debug("updateHideState: launcher is shown, show dock")
		m.setPropHideState(HideStateShow)
		return
	}

	hideMode := HideModeType(m.HideMode.Get())
	logger.Debug("updateHideState: mode is", hideMode)
	switch hideMode {
	case HideModeKeepShowing:
		m.setPropHideState(HideStateShow)

	case HideModeKeepHidden:
		m.setPropHideState(HideStateHide)

	case HideModeSmartHide:
		if delay {
			m.resetSmartHideModeTimer(time.Millisecond * 500)
		} else {
			m.resetSmartHideModeTimer(0)
		}
	}
}

func (m *DockManager) updateHideStateWithDelay() {
	m.updateHideState(true)
}

func (m *DockManager) updateHideStateWithoutDelay() {
	m.updateHideState(false)
}

func (m *DockManager) setPropHideState(hideState HideStateType) {
	logger.Debug("setPropHideState", hideState)
	if hideState == HideStateUnknown {
		logger.Warning("try setPropHideState to Unknown")
		return
	}

	if m.HideState != hideState {
		logger.Debugf("HideState %v => %v", m.HideState, hideState)
		m.HideState = hideState
		dbus.NotifyChange(m, "HideState")
	}
}
