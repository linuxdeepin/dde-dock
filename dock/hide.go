package dock

import (
	"dlib/dbus"
)

const (
	HideStateShowing int32 = iota
	HideStateShown
	HideStateHidding
	HideStateHidden
)

var (
	HideStateMap map[int32]string = map[int32]string{
		HideStateShowing: "HideStateShowing",
		HideStateShown:   "HideStateShown",
		HideStateHidding: "HideStateHidding",
		HideStateHidden:  "HideStateHidden",
	}
)

type HideStateManager struct {
	state int32

	StateChanged func(int32)
}

func NewHideStateManager(mode string) *HideStateManager {
	h := &HideStateManager{}

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
	if m.state == s {
		return s
	}

	logger.Info(HideStateMap[m.state], HideStateMap[s], HideStateMap[s])
	m.state = s
	logger.Debug("SetState emit StateChanged signal", HideStateMap[m.state])
	m.StateChanged(s)

	return s
}

func (m *HideStateManager) UpdateState() {
	logger.Info("UpdateState, HideState:", HideStateMap[m.state])
	switch setting.GetHideMode() {
	case HideModeKeepShowing:
		logger.Debug("KeepShowing Mode")
		m.state = HideStateShowing
	case HideModeAutoHide:
		logger.Debug("AutoHide Mode")
		m.state = HideStateShowing

		if region.mouseInRegion() {
			break
		}

		if hasMaximizeClient() {
			logger.Debug("has maximized client")
			m.state = HideStateHidding
		}
	case HideModeKeepHidden:
		logger.Debug("KeepHidden Mode")
		if region.mouseInRegion() {
			m.state = HideStateShowing
			break
		}

		m.state = HideStateHidding
	}

	logger.Debug("UpdateState emit StateChanged signal",
		HideStateMap[m.state])
	m.StateChanged(m.state)
}
