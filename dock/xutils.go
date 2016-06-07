package dock

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/ewmh"
)

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
	logger.Debug("activateWindow", win)
	err := changeCurrentWorkspaceToWindowWorkspace(win)
	if err != nil {
		return err
	}
	return ewmh.ActiveWindowReq(XU, win)
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
