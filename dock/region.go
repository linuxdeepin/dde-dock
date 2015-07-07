package dock

import (
	"errors"
	"github.com/BurntSushi/xgb/shape"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/icccm"
	"pkg.deepin.io/lib/dbus"
	"sync"
)

// Region表示dock有效的可接受事件区域以及可显示区域。
type Region struct {
}

var initShapeOnce sync.Once

func initShape() {
	initShapeOnce.Do(func() {
		shape.Init(XU.Conn())
	})
}

func NewRegion() *Region {
	r := Region{}

	return &r
}

func (r *Region) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.daemon.Dock",
		ObjectPath: "/dde/dock/DockRegion",
		Interface:  "dde.dock.DockRegion",
	}
}

func (r *Region) destroy() {
	dbus.UnInstallObject(r)
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

// GetDockRegion获取dock有效的可接受事件区域以及可显示区域。
func (r *Region) GetDockRegion() xproto.Rectangle {
	initShape()
	defer func() {
		if err := recover(); err != nil {
			logger.Warning("Region::GetDockRegion", err)
		}
	}()

	var dockRegion xproto.Rectangle
	dockWindow, err := r.getDockWindow()
	if err != nil {
		logger.Warning(err)
		return dockRegion
	}
	cookie := shape.GetRectangles(XU.Conn(), dockWindow, shape.SkInput)
	rep, err := cookie.Reply()
	if err != nil {
		logger.Warning("get rectangles failed:", err)
		return dockRegion
	}

	firstRect := rep.Rectangles[0]
	dockRegion.X = firstRect.X
	dockRegion.Y = firstRect.Y
	dockRegion.Width = firstRect.Width
	dockRegion.Height = firstRect.Height

	for i, rect := range rep.Rectangles {
		logger.Debugf("dock region %d: (%d,%d)->(%d,%d)",
			i, rect.X, rect.Y,
			int32(rect.X)+int32(rect.Width),
			int32(rect.Y)+int32(rect.Height),
		)

		if i == 0 {
			continue
		}

		if dockRegion.X > rect.X {
			dockRegion.X = rect.X
			dockRegion.Width = rect.Width
		}
		if dockRegion.Y > rect.Y {
			dockRegion.Y = rect.Y
		}

		dockRegion.Height += rect.Height
	}

	logger.Debugf("Dock region: (%d, %d)->(%d, %d)",
		dockRegion.X, dockRegion.Y,
		int32(dockRegion.X)+int32(dockRegion.Width),
		int32(dockRegion.Y)+int32(dockRegion.Height),
	)

	return dockRegion
}

func (r *Region) mouseInRegion() bool {
	region := r.GetDockRegion()
	cookie := xproto.QueryPointer(XU.Conn(), XU.RootWin())
	reply, err := cookie.Reply()
	if err != nil {
		return false
	}

	mouseX := int32(reply.RootX - displayRect.X)
	mouseY := int32(reply.RootY - displayRect.Y)

	logger.Debugf("mouse position: (%d, %d)", mouseX, mouseY)

	startX := int32(region.X)
	startY := int32(region.Y)

	endX := startX + int32(region.Width)
	endY := startY + int32(region.Height)

	inHorizontal := startX <= mouseX && mouseX <= endX
	inVertical := startY <= mouseY && mouseY <= endY
	inRegion := inHorizontal && inVertical

	return inRegion
}
