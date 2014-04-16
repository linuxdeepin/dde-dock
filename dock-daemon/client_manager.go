package main

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
)

var (
	XU, _                                 = xgbutil.NewConn()
	_NET_ACTIVE_WINDOW, _                 = xprop.Atm(XU, "_NET_ACTIVE_WINDOW")
	_NET_SHOWING_DESKTOP, _               = xprop.Atm(XU, "_NET_SHOWING_DESKTOP")
	lastActive              string        = ""
	activeWindow            xproto.Window = 0
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
	err := ewmh.ClientEvent(XU, xproto.Window(xid), "_NET_ACTIVE_WINDOW", 2)
	if err != nil {
		logger.Error("Actice window failed:", err)
		return false
	}
	return true
}

// maybe move to apps-builder
func (m *ClientManager) CloseWindow(xid uint32) bool {
	err := ewmh.ClientEvent(XU, xproto.Window(xid), "_NET_CLOSE_WINDOW")
	if err != nil {
		logger.Error("Actice window failed:", err)
		return false
	}
	return true
}

func (m *ClientManager) ToggleShowDesktop() {
	exec.Command("/usr/lib/deepin-daemon/desktop-toggle").Run()
}

func (m *ClientManager) listenRootWindow() {
	xwindow.New(XU, XU.RootWin()).Listen(xproto.EventMaskPropertyChange)
	xevent.PropertyNotifyFun(func(XU *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
		switch ev.Atom {
		case _NET_ACTIVE_WINDOW:
			var err error
			if activeWindow, err = ewmh.ActiveWindowGet(XU); err == nil {
				appId := find_app_id_by_xid(activeWindow)
				logger.Info("current active window:", appId)
				if appId == DDELauncher {
					logger.Info("active window is launcher")
					// TODO: hide dock
					// if m.RequireDockHide != nil {
					// 	m.RequireDockHide()
					// }
				} else {
					logger.Info("active window is not launcher")
					LAUNCHER, err :=
						launcher.Newlauncher("com.deepin.dde.launcher",
							"/com/deepin/dde/launcher")
					if err != nil {
						logger.Error(err)
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
		}
	}).Connect(XU, XU.RootWin())

	xevent.Main(XU)
}

func find_app_id_by_xid(xid xproto.Window) string {
	if id, err := xprop.PropValStr(xprop.GetProperty(XU, xid, "_DDE_DOCK_APP_ID")); err == nil {
		return id
	}
	pid, _ := ewmh.WmPidGet(XU, xid)
	iconName, _ := ewmh.WmIconNameGet(XU, xid)
	name, _ := ewmh.WmNameGet(XU, xid)
	wmClass, _ := icccm.WmClassGet(XU, xid)
	var wmInstance, wmClassName string
	if wmClass != nil {
		wmInstance = wmClass.Instance
		wmClassName = wmClass.Class
	}
	if pid == 0 {
	} else {
	}
	appId := find_app_id(pid, name, wmInstance, wmClassName, iconName)
	return appId
}
