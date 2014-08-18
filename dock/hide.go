package dock

import (
	"pkg.linuxdeepin.com/lib/dbus"
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

type HideStateManager struct {
	state               HideStateType
	toggleShowTimer     <-chan time.Time
	cleanToggleShowChan chan bool

	StateChanged func(HideStateType)
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

func (m *HideStateManager) SetState(s HideStateType) HideStateType {
	if m.state == s {
		return s
	}

	logger.Debug("SetState m.state:", m.state, "new state:", s)
	m.state = s
	logger.Debug("SetState emit StateChanged signal", m.state)
	m.StateChanged(s)

	return s
}

func (m *HideStateManager) UpdateState() {
	if m.toggleShowTimer != nil {
		logger.Info("in ToggleShow")
		return
	}
	state := m.state
	switch HideModeType(setting.GetHideMode()) {
	case HideModeKeepShowing:
		logger.Debug("KeepShowing Mode")
		state = HideStateShowing
	case HideModeAutoHide:
		logger.Debug("AutoHide Mode")
		state = HideStateShowing

		<-time.After(time.Millisecond * 100)
		if region.mouseInRegion() {
			logger.Debug("MouseInDockRegion")
			break
		}

		if hasMaximizeClientPre(activeWindow) {
			logger.Debug("active window is maximized client")
			state = HideStateHidding
		}
	case HideModeKeepHidden:
		logger.Debug("KeepHidden Mode")
		<-time.After(time.Millisecond * 100)
		if region.mouseInRegion() {
			logger.Debug("MouseInDockRegion")
			state = HideStateShowing
			break
		}

		state = HideStateHidding
	}

	if isLauncherShown {
		logger.Infof("launcher is opend, show dock")
		state = HideStateShowing
	}

	m.SetState(state)
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
		m.SetState(HideStateShowing)
	} else if m.state == HideStateShown || m.state == HideStateShowing {
		m.SetState(HideStateHidding)
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
