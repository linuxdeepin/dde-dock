package dock

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"pkg.linuxdeepin.com/dde-daemon"
	"time"
)

func init() {
	loader.Register(&loader.Module{
		Name:   "dock",
		Start:  Start,
		Stop:   Stop,
		Enable: true,
	})
}

var (
	XU     *xgbutil.XUtil
	TrayXU *xgbutil.XUtil

	//There variable must be initialized after the Xu/TrayXU has been
	//created.
	_NET_SHOWING_DESKTOP    xproto.Atom
	DEEPIN_SCREEN_VIEWPORT  xproto.Atom
	_NET_CLIENT_LIST        xproto.Atom
	_NET_ACTIVE_WINDOW      xproto.Atom
	ATOM_WINDOW_ICON        xproto.Atom
	ATOM_WINDOW_NAME        xproto.Atom
	ATOM_WINDOW_STATE       xproto.Atom
	ATOM_WINDOW_TYPE        xproto.Atom
	ATOM_DOCK_APP_ID        xproto.Atom
	_NET_SYSTEM_TRAY_S0     xproto.Atom
	_NET_SYSTEM_TRAY_OPCODE xproto.Atom

	// ATOM_DEEPIN_WINDOW_VIEWPORTS, _ = xprop.Atm(XU, "DEEPIN_WINDOW_VIEWPORTS")

	mouseAreaTimer   *time.Timer
	TOGGLE_HIDE_TIME = time.Millisecond * 400
)
