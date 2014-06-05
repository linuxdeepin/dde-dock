package dock

import (
	"dlib/dbus"
	"errors"
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/shape"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/icccm"
)

var (
	conn, _ = xgb.NewConn()
)

type Region struct {
}

func NewRegion() *Region {
	shape.Init(conn)
	r := Region{}

	return &r
}

func (r *Region) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Dock",
		"/dde/dock/DockRegion",
		"dde.dock.DockRegion",
	}
}

func (r *Region) getDockWindow() (xproto.Window, error) {
	windows, _ := ewmh.ClientListGet(XU)
	for _, xid := range windows {
		res, err := icccm.WmClassGet(XU, xid)
		if err == nil && res.Instance == "dde-dock" {
			return xid, nil
		}
	}
	return 0, errors.New("find dock window failed, it's not existed.")
}

func (r *Region) GetDockRegion() xproto.Rectangle {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Region::GetDockRegion", err)
		}
	}()

	dockRegion := xproto.Rectangle{0, 0, 0, 0}
	dockWindow, err := r.getDockWindow()
	if err != nil {
		logger.Error(err)
		return dockRegion
	}
	cookie := shape.GetRectangles(conn, dockWindow, shape.SkInput)
	rep, err := cookie.Reply()
	if err != nil {
		logger.Error("get rectangles failed:", err)
		return dockRegion
	}
	for _, rect := range rep.Rectangles {
		logger.Debugf("dock region: (%d,%d)->(%d,%d)",
			rect.X, rect.Y,
			int32(rect.X)+int32(rect.Width),
			int32(rect.Y)+int32(rect.Height),
		)
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

func (r *Region) mouseInRegion() bool {
	region := r.GetDockRegion()
	cookie := xproto.QueryPointer(conn, XU.RootWin())
	reply, err := cookie.Reply()
	if err != nil {
		return false
	}

	mouseX := int32(reply.RootX)
	mouseY := int32(reply.RootY)

	startX := int32(region.X)
	startY := int32(region.Y)

	endX := startX + int32(region.Width)
	endY := startY + int32(region.Height)

	inHorizontal := startX <= mouseX && mouseX <= endX
	inVertical := startY <= mouseY && mouseY <= endY
	inRegion := inHorizontal && inVertical

	return inRegion
}
