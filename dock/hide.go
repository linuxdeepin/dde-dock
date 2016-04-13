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
	"pkg.deepin.io/lib/dbus"
	"sync"
	"time"
)

type HideStateType int32

const (
	HideStateShowing HideStateType = iota
	HideStateShown
	HideStateHidding
	HideStateHidden
)

func (s HideStateType) String() string {
	switch s {
	case HideStateShowing:
		return "HideStateShowing"
	case HideStateShown:
		return "HideStateShown"
	case HideStateHidding:
		return "HideStateHidding"
	case HideStateHidden:
		return "HideStateHidden"
	default:
		return "Unknown state"
	}
}

const (
	TriggerShow int32 = iota
	TriggerHide
)

type HideStateManager struct {
	state       HideStateType
	ChangeState func(int32)

	smartHideModeTimer *time.Timer
	smartHideModeMutex sync.Mutex
}

func NewHideStateManager(mode HideModeType) *HideStateManager {
	m := &HideStateManager{}
	if mode == HideModeKeepHidden {
		m.state = HideStateHidden
	} else {
		m.state = HideStateShown
	}

	m.smartHideModeTimer = time.AfterFunc(10*time.Second, m.smartHideModeTimerExpired)
	m.smartHideModeTimer.Stop()
	return m
}

func (e *HideStateManager) destroy() {
	if e.smartHideModeTimer != nil {
		e.smartHideModeTimer.Stop()
		e.smartHideModeTimer = nil
	}
	dbus.UnInstallObject(e)
}

func (e *HideStateManager) GetDBusInfo() dbus.DBusInfo {
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

func shouldHideOnSmartHideMode() bool {
	return isWindowDockOverlap(activeWindow)
}

func (m *HideStateManager) smartHideModeTimerExpired() {
	logger.Debug("smartHideModeTimer expired!")
	if isLauncherShown {
		logger.Debug("launcher is showing, dock show")
		m.emitSignalChangeState(TriggerShow)
		return
	}
	if shouldHideOnSmartHideMode() {
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
		if shouldHideOnSmartHideMode() {
			logger.Debug("smartHideModeDelayHandle: show -> hide")
			m.resetSmartHideModeTimer(time.Millisecond * 500)
		} else {
			logger.Debug("smartHideModeDelayHandle: show -> show")
			m.cancelSmartHideModeTimer()
		}

	case HideStateHidden:
		if shouldHideOnSmartHideMode() {
			logger.Debug("smartHideModeDelayHandle: hide -> hide")
			m.cancelSmartHideModeTimer()
		} else {
			logger.Debug("smartHideModeDelayHandle: hide -> show")
			m.resetSmartHideModeTimer(time.Millisecond * 500)
		}
	}
}

func (m *HideStateManager) updateState(delay bool) {
	if isLauncherShown {
		logger.Debug("updateState: launcher is showing, show dock")
		m.emitSignalChangeState(TriggerShow)
		return
	}

	hideMode := HideModeType(setting.GetHideMode())
	logger.Debug("updateState:", hideMode)
	switch hideMode {
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

// super+H 功能废弃，此接口废弃
func (m *HideStateManager) ToggleShow() {
}
