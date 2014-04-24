package main

import (
	"dlib/dbus"
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/shape"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/icccm"
)

var (
	C, _ = xgb.NewConn()
)

type Region struct {
	dockWindow xproto.Window
}

func NewRegion() *Region {
	shape.Init(C)
	r := Region{0}

	windows, _ := ewmh.ClientListGet(XU)
	for _, xid := range windows {
		res, err := icccm.WmClassGet(XU, xid)
		if err == nil && res.Instance == "dde-dock" {
			r.dockWindow = xid
			break
		}
	}

	return &r
}

func (r *Region) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Dock",
		"/dde/dock/DockRegion",
		"dde.dock.DockRegion",
	}
}

func (r *Region) GetDockRegion() xproto.Rectangle {
	dockRegion := xproto.Rectangle{0, 0, 0, 0}
	cookie := shape.GetRectangles(C, r.dockWindow, shape.SkInput)
	rep, _ := cookie.Reply()
	for _, rect := range rep.Rectangles {
		logger.Infof("dock region: %dx%d(%d,%d)", rect.Width, rect.Height, rect.X, rect.Y)
		if dockRegion.X == 0 || dockRegion.X > rect.X {
			dockRegion.X = rect.X
			dockRegion.Width = rect.Width
		}
		if dockRegion.Y == 0 || dockRegion.Y > rect.Y {
			dockRegion.Y = rect.Y
		}

		dockRegion.Height += rect.Height
	}

	return dockRegion
}
