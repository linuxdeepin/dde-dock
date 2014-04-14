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
)

var (
	XU, _                        = xgbutil.NewConn()
	_NET_ACTIVE_WINDOW, _        = xprop.Atm(XU, "_NET_ACTIVE_WINDOW")
	lastActive            string = ""
)

const (
	DDELauncher string = "dde-launcher"
)

type SpecialWindowManager struct {
	RequireDockHide func()
	RequireDockShow func()
}

func NewSpecialWindowManager() *SpecialWindowManager {
	return &SpecialWindowManager{}
}

func (m *SpecialWindowManager) listenRootWindow() {
	xwindow.New(XU, XU.RootWin()).Listen(xproto.EventMaskPropertyChange)
	xevent.PropertyNotifyFun(func(XU *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
		switch ev.Atom {
		case _NET_ACTIVE_WINDOW:
			if activedWindow, err := ewmh.ActiveWindowGet(XU); err == nil {
				appId := find_app_id_by_xid(activedWindow)
				logger.Info("current active window:", appId)
				if appId == DDELauncher {
					logger.Info("active window is launcher")
					// TODO: hide dock
					if m.RequireDockHide != nil {
						m.RequireDockHide()
					}
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
					if m.RequireDockShow != nil &&
						lastActive == DDELauncher {
						m.RequireDockShow()
					}
				}
				lastActive = appId
			}
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
