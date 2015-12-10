package dock

import (
	"fmt"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xprop"
	"github.com/BurntSushi/xgbutil/xwindow"
	"os/exec"
	"pkg.deepin.io/lib/dbus"
)

var (
	activeWindow    xproto.Window = 0
	isLauncherShown bool          = false
	currentViewport []uint        = nil
)

const (
	DDELauncher string = "dde-launcher"
)

// ClientManager用来管理启动程序相关窗口。
type ClientManager struct {
	// ActiveWindowChanged会在焦点窗口被改变时触发，会将最新的焦点窗口id发送给监听者。
	ActiveWindowChanged func(xid uint32)

	// ShowingDesktopChanged会在_NET_SHOWING_DESKTOP改变时被触发。
	ShowingDesktopChanged func()
}

// NewClientManager creates a new client manager.
func NewClientManager() *ClientManager {
	return &ClientManager{}
}

// CurrentActiveWindow会返回当前焦点窗口的窗口id。
func (m *ClientManager) CurrentActiveWindow() uint32 {
	return uint32(activeWindow)
}

func changeWorkspaceIfNeeded(xid xproto.Window) error {
	desktopNum, err := xprop.PropValNum(xprop.GetProperty(XU, xid, "_NET_WM_DESKTOP"))
	if err != nil {
		return fmt.Errorf("Get _NET_WM_DESKTOP failed: %s", err)
	}

	currentDesktop, err := ewmh.CurrentDesktopGet(XU)
	if err != nil {
		return fmt.Errorf("Get _NET_CURRENT_DESKTOP failed: %v", err)
	}

	if currentDesktop == desktopNum {
		logger.Debug("No need to change workspace, the current desktop is already %v", currentDesktop)
		return nil
	}

	timeStamp, err := ewmh.WmUserTimeGet(XU, xid)
	if err != nil {
		logger.Debugf("Get timestamp of 0x%x failed: %v", uint32(xid), err)
	}

	err = ewmh.ClientEvent(XU, XU.RootWin(), "_NET_CURRENT_DESKTOP", int(desktopNum), int(timeStamp))
	if err != nil {
		return fmt.Errorf("Send ClientMessage Failed: %v", err)
	}

	return nil
}

func activateWindow(xid xproto.Window) error {
	err := changeWorkspaceIfNeeded(xid)
	if err != nil {
		logger.Warning(err)
	}
	return ewmh.ActiveWindowReq(XU, xid)
}

// ActiveWindow会激活给定id的窗口，被激活的窗口将通常会程序焦点窗口。(废弃，名字应该是ActivateWindow，当时手残打错了，此接口会在之后被移除，请使用正确的接口)
func (m *ClientManager) ActiveWindow(xid uint32) bool {
	err := activateWindow(xproto.Window(xid))
	if err != nil {
		logger.Warning("Activate window failed:", err)
		return false
	}
	return true
}

// ActivateWindow会激活给定id的窗口，被激活的窗口通常会成为焦点窗口。
func (m *ClientManager) ActivateWindow(xid uint32) bool {
	err := activateWindow(xproto.Window(xid))
	if err != nil {
		logger.Warning("Activate window failed:", err)
		return false
	}
	return true
}

// CloseWindow会将传入id的窗口关闭。
func (m *ClientManager) CloseWindow(xid uint32) bool {
	err := ewmh.CloseWindow(XU, xproto.Window(xid))
	if err != nil {
		logger.Warning("Close window failed:", err)
		return false
	}
	return true
}

// ToggleShowDesktop会触发显示桌面，当桌面显示时，会将窗口恢复，当桌面未显示时，会隐藏窗口以显示桌面。
func (m *ClientManager) ToggleShowDesktop() {
	exec.Command("/usr/lib/deepin-daemon/desktop-toggle").Run()
}

// IsLauncherShown判断launcher是否已经显示。
func (m *ClientManager) IsLauncherShown() bool {
	return isLauncherShown
}

func walkClientList(pre func(xproto.Window) bool) bool {
	list, err := ewmh.ClientListGet(XU)
	if err != nil {
		logger.Debug("Can't get _NET_CLIENT_LIST", err)
		return false
	}

	for _, c := range list {
		if pre(c) {
			return true
		}
	}

	return false
}

func isMaximizeVertClientPre(xid xproto.Window) bool {
	state, _ := ewmh.WmStateGet(XU, xid)
	return contains(state, "_NET_WM_STATE_MAXIMIZED_VERT")
}

func isHiddenPre(xid xproto.Window) bool {
	state, _ := ewmh.WmStateGet(XU, xid)
	return contains(state, "_NET_WM_STATE_HIDDEN")
}

func isCoverWorkspace(workspaces [][]uint, workspace []uint) bool {
	for _, w := range workspaces {
		if workspace[0] == w[0] && workspace[1] == w[1] {
			return true
		}
	}
	return false
}

func updateCurrentViewport() {
	currentViewport, _ = xprop.PropValNums(
		xprop.GetProperty(
			XU,
			XU.RootWin(),
			"DEEPIN_SCREEN_VIEWPORT",
		))
}

// works for old deepin wm.
func checkDeepinWindowViewports(xid xproto.Window) (bool, error) {
	viewports, err := xprop.PropValNums(xprop.GetProperty(XU, xid,
		"DEEPIN_WINDOW_VIEWPORTS"))
	if err != nil {
		return false, err
	}

	workspaces := make([][]uint, 0)
	for i := uint(0); i < viewports[0]; i++ {
		viewport := make([]uint, 2)
		viewport[0] = viewports[i*2+1]
		viewport[1] = viewports[i*2+2]
		workspaces = append(workspaces, viewport)
	}
	if currentViewport == nil {
		updateCurrentViewport()
	}
	return isCoverWorkspace(workspaces, currentViewport), nil
}

// works for new deepin wm.
func checkCurrentDesktop(xid xproto.Window) (bool, error) {
	num, err := xprop.PropValNum(xprop.GetProperty(XU, xid, "_NET_WM_DESKTOP"))
	if err != nil {
		return false, err
	}

	currentDesktop, err := xprop.PropValNum(xprop.GetProperty(XU, XU.RootWin(), "_NET_CURRENT_DESKTOP"))
	if err != nil {
		return false, err
	}

	return num == currentDesktop, nil
}

func onCurrentWorkspacePre(xid xproto.Window) bool {
	isOnCurrentWorkspace, err := checkDeepinWindowViewports(xid)
	if err != nil {
		isOnCurrentWorkspace, err = checkCurrentDesktop(xid)
		if err != nil {
			return false
		}
		return isOnCurrentWorkspace
	}
	return isOnCurrentWorkspace
}

func hasMaximizeClientPre(xid xproto.Window) bool {
	isMax := isMaximizeVertClientPre(xid)
	isHidden := isHiddenPre(xid)
	onCurrentWorkspace := onCurrentWorkspacePre(xid)
	logger.Debug("isMax:", isMax, "isHidden:", isHidden,
		"onCurrentWorkspace:", onCurrentWorkspace)
	return isMax && !isHidden && onCurrentWorkspace
}

func hasMaximizeClient() bool {
	return walkClientList(hasMaximizeClientPre)
}

func isDeepinLauncher(xid xproto.Window) bool {
	res, err := icccm.WmClassGet(XU, xid)
	if err != nil {
		return false
	}

	return res.Instance == DDELauncher
}

func isWindowOnPrimaryScreen(xid xproto.Window) bool {
	var err error

	win := xwindow.New(XU, xid)
	// include shadow
	gemo, err := win.DecorGeometry()
	if err != nil {
		logger.Debug(err)
		return false
	}

	displayRectX := (int)(displayRect.X)
	displayRectY := (int)(displayRect.Y)
	displayRectWidth := (int)(displayRect.Width)
	displayRectHeight := (int)(displayRect.Height)

	SHADOW_OFFSET := 10
	gemoX := gemo.X() + SHADOW_OFFSET
	gemoY := gemo.Y() + SHADOW_OFFSET
	isOnPrimary := gemoX+SHADOW_OFFSET >= displayRectX &&
		gemoX < displayRectX+displayRectWidth &&
		gemoY >= displayRectY &&
		gemoY < displayRectY+displayRectHeight

	logger.Debugf("isWindowOnPrimaryScreen: %dx%d, %dx%d, %v", gemo.X(),
		gemo.Y(), displayRect.X, displayRect.Y, isOnPrimary)

	return isOnPrimary
}

func isWindowOverlapDock(xid xproto.Window) bool {
	win := xwindow.New(XU, xid)
	rect, err := win.DecorGeometry()
	if err != nil {
		logger.Warningf("isWindowOverlapDock GetDecorGeometry of 0x%x failed: %s", xid, err)
		return false
	}

	winX := int32(rect.X())
	winY := int32(rect.Y())
	winWidth := int32(rect.Width())
	winHeight := int32(rect.Height())

	dockX := int32(displayRect.X) + (int32(displayRect.Width)-
		dockProperty.PanelWidth)/2
	dockY := int32(displayRect.Y) + int32(displayRect.Height) -
		dockProperty.Height
	dockWidth := int32(displayRect.Width)
	if DisplayModeType(setting.GetDisplayMode()) == DisplayModeModernMode {
		dockWidth = dockProperty.PanelWidth
	}

	// TODO: dock on the other side like top, left.
	return dockY < winY+winHeight &&
		dockX < winX+winWidth && dockX+dockWidth > winX
}

func (m *ClientManager) listenRootWindow() {
	var update = func() {
		list, err := ewmh.ClientListGet(XU)
		if err != nil {
			logger.Warning("Can't Get _NET_CLIENT_LIST", err)
		}
		isLauncherShown = false
		for _, xid := range list {
			if !isDeepinLauncher(xid) {
				continue
			}

			winProps, err :=
				xproto.GetWindowAttributes(XU.Conn(),
					xid).Reply()
			if err != nil {
				break
			}
			if winProps.MapState == xproto.MapStateViewable {
				isLauncherShown = true
			}
			break
		}
		ENTRY_MANAGER.runtimeAppChanged(list)
	}

	xwindow.New(XU, XU.RootWin()).Listen(xproto.EventMaskPropertyChange)
	xevent.PropertyNotifyFun(func(XU *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
		switch ev.Atom {
		case _NET_CLIENT_LIST:
			update()
		case _NET_ACTIVE_WINDOW:
			var err error
			isLauncherShown = false
			if activeWindow, err = ewmh.ActiveWindowGet(XU); err == nil {
				// loop gets better performance than find_app_id_by_xid.
				// setLeader/updateState will filter invalid xid.
				for _, app := range ENTRY_MANAGER.runtimeApps {
					app.setLeader(activeWindow)
					app.updateState(activeWindow)
				}

				if isDeepinLauncher(activeWindow) {
					isLauncherShown = true
				}

				dbus.Emit(m, "ActiveWindowChanged", uint32(activeWindow))
			}

			hideMode := HideModeType(setting.GetHideMode())
			if hideMode == HideModeSmartHide || hideMode == HideModeKeepHidden {
				hideModemanager.UpdateState()
			}
		case _NET_SHOWING_DESKTOP:
			dbus.Emit(m, "ShowingDesktopChanged")
		case DEEPIN_SCREEN_VIEWPORT:
			updateCurrentViewport()
		}
	}).Connect(XU, XU.RootWin())

	update()
	xevent.Main(XU)
}
