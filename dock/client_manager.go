package dock

import (
	"dbus/com/deepin/dde/launcher"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xprop"
	"github.com/BurntSushi/xgbutil/xwindow"
	"os/exec"
	"pkg.linuxdeepin.com/lib/dbus"
)

var (
	lastActive      string        = ""
	activeWindow    xproto.Window = 0
	isLauncherShown bool          = false
	currentViewport []uint        = nil
)

const (
	DDELauncher string = "dde-launcher"
)

type ClientManager struct {
	ActiveWindowChanged   func(xid uint32)
	ShowingDesktopChanged func()
}

func NewClientManager() *ClientManager {
	return &ClientManager{}
}

func (m *ClientManager) CurrentActiveWindow() uint32 {
	return uint32(activeWindow)
}

// maybe move to apps-builder
func (m *ClientManager) ActiveWindow(xid uint32) bool {
	err := ewmh.ActiveWindowReq(XU, xproto.Window(xid))
	if err != nil {
		logger.Warning("Actice window failed:", err)
		return false
	}
	return true
}

// maybe move to apps-builder
func (m *ClientManager) CloseWindow(xid uint32) bool {
	err := ewmh.CloseWindow(XU, xproto.Window(xid))
	if err != nil {
		logger.Warning("Actice window failed:", err)
		return false
	}
	return true
}

func (m *ClientManager) ToggleShowDesktop() {
	exec.Command("/usr/lib/deepin-daemon/desktop-toggle").Run()
}

func (m *ClientManager) IsLauncherShown() bool {
	return isLauncherShown
}

func walkClientList(pre func(xproto.Window) bool) bool {
	list, err := ewmh.ClientListGet(XU)
	if err != nil {
		logger.Warning("Can't get _NET_CLIENT_LIST", err)
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

func onCurrentWorkspacePre(xid xproto.Window) bool {
	viewports, err := xprop.PropValNums(xprop.GetProperty(XU, xid,
		"DEEPIN_WINDOW_VIEWPORTS"))
	if err != nil {
		logger.Warning("get DEEPIN_WINDOW_VIEWPORTS failed", err)
		return false
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
	return isCoverWorkspace(workspaces, currentViewport)
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
		logger.Warning(err)
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
				appId := find_app_id_by_xid(activeWindow,
					DisplayModeType(setting.GetDisplayMode()))
				logger.Debug("current active window:", appId)
				if rApp, ok := ENTRY_MANAGER.runtimeApps[appId]; ok {
					logger.Debug("find runtime app")
					rApp.setLeader(activeWindow)
					rApp.updateState(activeWindow)
				}

				logger.Info("active window is", appId)
				if appId != DDELauncher {
					LAUNCHER, err := launcher.NewLauncher(
						"com.deepin.dde.launcher",
						"/com/deepin/dde/launcher",
					)
					if err != nil {
						logger.Debug(err)
					} else {
						LAUNCHER.Hide()
					}
				} else {
					isLauncherShown = true
				}

				lastActive = appId
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
