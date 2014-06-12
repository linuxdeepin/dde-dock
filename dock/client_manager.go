package dock

import (
	"dbus/com/deepin/dde/launcher"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	// "github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xprop"
	"github.com/BurntSushi/xgbutil/xwindow"
	"os/exec"
)

var (
	// XU, _                                 = xgbutil.NewConn()
	// _NET_ACTIVE_WINDOW, _                 = xprop.Atm(XU, "_NET_ACTIVE_WINDOW")
	_NET_SHOWING_DESKTOP, _                 = xprop.Atm(XU, "_NET_SHOWING_DESKTOP")
	DEEPIN_SCREEN_VIEWPORT, _               = xprop.Atm(XU, "DEEPIN_SCREEN_VIEWPORT")
	lastActive                string        = ""
	activeWindow              xproto.Window = 0
	currentViewport, _                      = xprop.PropValNums(xprop.GetProperty(XU, XU.RootWin(),
		"DEEPIN_SCREEN_VIEWPORT"))
)

const (
	DDELauncher string = "dde-launcher"
)

type ClientManager struct {
	RequireDockHide func()
	RequireDockShow func()

	RequireDockHideWithAnimation func()
	RequireDockShowWithAnimation func()

	RequireDockHideWithoutChangeWorkarea func()
	RequireDockShowWithoutChangeWorkarea func()

	ActiveWindowChanged   func(xid uint32)
	ShowingDesktopChanged func()
}

func NewClientManager() *ClientManager {
	return &ClientManager{}
}

func (m *ClientManager) ShowDockWithAnimation() {
	m.RequireDockShowWithAnimation()
}

func (m *ClientManager) HideDockWithAnimation() {
	m.RequireDockHideWithAnimation()
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

func isMaximizeClientPre(xid xproto.Window) bool {
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
	return isCoverWorkspace(workspaces, currentViewport)
}

func hasMaximizeClientPre(xid xproto.Window) bool {
	isMax := isMaximizeClientPre(xid)
	isHidden := isHiddenPre(xid)
	onCurrentWorkspace := onCurrentWorkspacePre(xid)
	logger.Debug("isMax:", isMax, "isHidden:", isHidden,
		"onCurrentWorkspace:", onCurrentWorkspace)
	return isMax && !isHidden && onCurrentWorkspace
}

func hasMaximizeClient() bool {
	return walkClientList(hasMaximizeClientPre)
}

func (m *ClientManager) listenRootWindow() {
	var update = func() {
		list, err := ewmh.ClientListGet(XU)
		if err != nil {
			logger.Warning("Can't Get _NET_CLIENT_LIST", err)
		}
		ENTRY_MANAGER.runtimeAppChangged(list)
	}

	xwindow.New(XU, XU.RootWin()).Listen(xproto.EventMaskPropertyChange)
	xevent.PropertyNotifyFun(func(XU *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
		switch ev.Atom {
		case _NET_CLIENT_LIST:
			update()
		case _NET_ACTIVE_WINDOW:
			var err error
			if activeWindow, err = ewmh.ActiveWindowGet(XU); err == nil {
				appId := find_app_id_by_xid(activeWindow)
				logger.Debug("current active window:", appId)
				if rApp, ok := ENTRY_MANAGER.runtimeApps[appId]; ok {
					rApp.setLeader(activeWindow)
				}
				if appId == DDELauncher {
					logger.Info("active window is launcher")
					// TODO: hide dock
					// if m.RequireDockHide != nil {
					// 	m.RequireDockHide()
					// }
				} else {
					logger.Debug("active window is not launcher")
					LAUNCHER, err :=
						launcher.NewLauncher("com.deepin.dde.launcher",
							"/com/deepin/dde/launcher")
					if err != nil {
						logger.Warning(err)
					} else {
						LAUNCHER.Hide()
					}

					// TODO: show dock
					// if m.RequireDockShow != nil &&
					// 	lastActive == DDELauncher {
					// 	m.RequireDockShow()
					// }
				}
				lastActive = appId
				m.ActiveWindowChanged(uint32(activeWindow))
			}
		case _NET_SHOWING_DESKTOP:
			m.ShowingDesktopChanged()
		case DEEPIN_SCREEN_VIEWPORT:
			currentViewport, _ = xprop.PropValNums(xprop.GetProperty(XU, XU.RootWin(),
				"DEEPIN_SCREEN_VIEWPORT"))
		}
	}).Connect(XU, XU.RootWin())

	update()
	hideModemanager.UpdateState()
	xevent.Main(XU)
}
