package dock

import (
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/log"
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
	state               HideStateType
	toggleShowTimer     <-chan time.Time
	cleanToggleShowChan chan bool

	ChangeState func(int32)
}

func NewHideStateManager(mode HideModeType) *HideStateManager {
	h := &HideStateManager{}
	h.toggleShowTimer = nil
	h.cleanToggleShowChan = make(chan bool, 1)

	if mode == HideModeKeepHidden {
		h.state = HideStateHidden
	} else {
		h.state = HideStateShown
	}

	return h
}

func (e *HideStateManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Dock",
		"/dde/dock/HideStateManager",
		"dde.dock.HideStateManager",
	}
}

func (m *HideStateManager) SetState(s int32) int32 {
	state := HideStateType(s)
	logger.Debug("SetState m.state:", m.state, "new state:", state)
	if m.state == state {
		logger.Info("New HideState is the same as the old:", state)
		return s
	}

	m.state = state

	return s
}

func (m *HideStateManager) UpdateState() {
	if m.toggleShowTimer != nil && !isLauncherShown {
		logger.Info("in ToggleShow")
		return
	}
	trigger := TriggerShow
	switch HideModeType(setting.GetHideMode()) {
	case HideModeKeepShowing:
		logger.Debug("KeepShowing Mode")
	case HideModeAutoHide:
		logger.Debug("AutoHide Mode")

		<-time.After(time.Millisecond * 100)
		if region.mouseInRegion() {
			logger.Debug("MouseInDockRegion")
			break
		}

		if logger.GetLogLevel() == log.LevelDebug {
			hasMax := hasMaximizeClientPre(activeWindow)
			isOnPrimary := isWindowOnPrimaryScreen(activeWindow)
			logger.Infof("hasMax: %v, isOnPrimary: %v", hasMax,
				isOnPrimary)
		}

		if isWindowOnPrimaryScreen(activeWindow) &&
			hasMaximizeClientPre(activeWindow) {
			logger.Debug("active window is maximized client")
			trigger = TriggerHide
		}
	case HideModeKeepHidden:
		logger.Debug("KeepHidden Mode")
		<-time.After(time.Millisecond * 100)
		if region.mouseInRegion() {
			logger.Debug("MouseInDockRegion")
			break
		}

		trigger = TriggerHide
	case HideModeSmartHide:
		logger.Debug("SmartHide Mode")

		<-time.After(time.Millisecond * 100)
		if region.mouseInRegion() {
			logger.Debug("mouse in region")
			break
		}

		if isWindowOnPrimaryScreen(activeWindow) &&
			hasMaximizeClientPre(activeWindow) {
			logger.Debug("active window is maximized client")
			trigger = TriggerHide
			break
		}

		for _, app := range ENTRY_MANAGER.runtimeApps {
			for _, winInfo := range app.xids {
				if winInfo.OverlapDock {
					logger.Warning("overlap dock")
					trigger = TriggerHide
					break
				}
			}
		}
	}

	if isLauncherShown {
		m.CancelToggleShow()
		logger.Info("launcher is opend, show dock")
		trigger = TriggerShow
	}

	if m.ChangeState != nil &&
		(m.state != HideStateShown && trigger == TriggerShow) ||
		(m.state != HideStateHidden && trigger == TriggerHide) {
		m.ChangeState(trigger)
	}
}

func (m *HideStateManager) CancelToggleShow() {
	if m.toggleShowTimer != nil {
		logger.Info("Cancel ToggleShow")
		m.cleanToggleShowChan <- true
		m.toggleShowTimer = nil
	}
}

func (m *HideStateManager) ToggleShow() {
	logger.Info("cancel ToggleShow on ToggleShow")
	m.CancelToggleShow()

	if m.state == HideStateHidden || m.state == HideStateHidding {
		m.ChangeState(TriggerShow)
	} else if m.state == HideStateShown || m.state == HideStateShowing {
		m.ChangeState(TriggerHide)
	}

	m.toggleShowTimer = time.After(time.Second * 3)
	go func() {
		select {
		case <-m.toggleShowTimer:
			logger.Info("ToggleShow is done")
			m.toggleShowTimer = nil
			m.UpdateState()
		case <-m.cleanToggleShowChan:
			logger.Info("ToggleShow is cancelled")
			return
		}
	}()
}
