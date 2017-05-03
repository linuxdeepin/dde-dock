package dock

import (
	"bytes"
	"encoding/base64"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xprop"
	"github.com/BurntSushi/xgbutil/xrect"
	"github.com/BurntSushi/xgbutil/xwindow"
	"strings"
)

func iconifyWindow(win xproto.Window) {
	logger.Debug("iconifyWindow", win)
	ewmh.ClientEvent(XU, win, "WM_CHANGE_STATE", icccm.StateIconic)
}

type windowFrameExtents struct {
	Left, Right, Top, Bottom uint
}

func getWindowFrameExtents(xu *xgbutil.XUtil, win xproto.Window) (*windowFrameExtents, error) {
	reply, err := xprop.GetProperty(xu, win, "_NET_FRAME_EXTENTS")
	if err != nil {
		// try _GTK_FRAME_EXTENTS
		reply, err = xprop.GetProperty(xu, win, "_GTK_FRAME_EXTENTS")
		if err != nil {
			return nil, err
		}
	}
	nums, err := xprop.PropValNums(reply, err)
	if err != nil {
		return nil, err
	}
	extents := &windowFrameExtents{nums[0], nums[1], nums[2], nums[3]}
	return extents, err
}

func getWindowGeometry(xu *xgbutil.XUtil, win xproto.Window) (xrect.Rect, error) {
	window := xwindow.New(xu, win)
	winRect, err := window.DecorGeometry()
	if err != nil {
		return nil, err
	}
	frameExtents, _ := getWindowFrameExtents(xu, win)
	if frameExtents != nil {
		x := winRect.X() + int(frameExtents.Left)
		y := winRect.Y() + int(frameExtents.Top)
		w := winRect.Width() - int(frameExtents.Left+frameExtents.Right)
		h := winRect.Height() - int(frameExtents.Top+frameExtents.Bottom)
		return xrect.New(x, y, w, h), nil
	}
	return winRect, nil
}

func getWmName(xu *xgbutil.XUtil, win xproto.Window) string {
	// get _NET_WM_NAME
	name, err := ewmh.WmNameGet(xu, win)
	if err != nil || name == "" {
		// get WM_NAME
		name, _ = icccm.WmNameGet(xu, win)
	}
	return strings.Replace(name, "\x00", "", -1)
}

func getWmPid(xu *xgbutil.XUtil, win xproto.Window) uint {
	pid, _ := ewmh.WmPidGet(xu, win)
	return pid
}

// WM_CLIENT_LEADER
func getWmClientLeader(xu *xgbutil.XUtil, win xproto.Window) (xproto.Window, error) {
	return xprop.PropValWindow(xprop.GetProperty(xu, win, "WM_CLIENT_LEADER"))
}

// WM_TRANSIENT_FOR
func getWmTransientFor(xu *xgbutil.XUtil, win xproto.Window) (xproto.Window, error) {
	return xprop.PropValWindow(xprop.GetProperty(xu, win, "WM_TRANSIENT_FOR"))
}

// _NET_WM_WINDOW_OPACITY
func getWmWindowOpacity(xu *xgbutil.XUtil, win xproto.Window) (uint, error) {
	return xprop.PropValNum(xprop.GetProperty(xu, win, "_NET_WM_WINDOW_OPACITY"))
}

func getWmCommand(xu *xgbutil.XUtil, win xproto.Window) ([]string, error) {
	command, err := xprop.PropValStrs(xprop.GetProperty(xu, win, "WM_COMMAND"))
	return command, err
}

func getWindowGtkApplicationId(xu *xgbutil.XUtil, win xproto.Window) string {
	gtkAppId, _ := xprop.PropValStr(xprop.GetProperty(xu, win, "_GTK_APPLICATION_ID"))
	return gtkAppId
}

func getWmWindowRole(xu *xgbutil.XUtil, win xproto.Window) string {
	role, _ := xprop.PropValStr(xprop.GetProperty(xu, win, "WM_WINDOW_ROLE"))
	return role
}

func getIconFromWindow(xu *xgbutil.XUtil, win xproto.Window) string {
	// find largest icon
	icon, err := xgraphics.FindIcon(xu, win, 0, 0)
	// FIXME: gets empty icon for minecraft
	if err == nil {
		buf := bytes.NewBuffer(nil)
		icon.WritePng(buf)
		return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
	}
	// NOTE: Do not get icon from property _NET_WM_ICON_NAME
	logger.Debugf("getIconFromWindow failed win %v, err: %v", win, err)
	return ""
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
	logger.Debug("activateWindow", win)
	err := changeCurrentWorkspaceToWindowWorkspace(win)
	if err != nil {
		return err
	}
	return ewmh.ActiveWindowReq(XU, win)
}

func isHiddenPre(win xproto.Window) bool {
	state, _ := ewmh.WmStateGet(XU, win)
	return strSliceContains(state, "_NET_WM_STATE_HIDDEN")
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
		// 也许是窗口跳过窗口管理器了，如 dde-control-center
		return true
	}
	return isOnCurrentWorkspace
}
