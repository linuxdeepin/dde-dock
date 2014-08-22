package dock

import "pkg.linuxdeepin.com/dde-daemon"
import "github.com/BurntSushi/xgbutil/xprop"
import "github.com/BurntSushi/xgbutil"

func init() {
	loader.Register(&loader.Module{"dock", Start, Stop, true})
}

var (
	XU, _     = xgbutil.NewConn()
	TrayXU, _ = xgbutil.NewConn()

	//There variable must be initialized after the Xu/TrayXU has been
	//created.
	_NET_SHOWING_DESKTOP, _   = xprop.Atm(XU, "_NET_SHOWING_DESKTOP")
	DEEPIN_SCREEN_VIEWPORT, _ = xprop.Atm(XU, "DEEPIN_SCREEN_VIEWPORT")
	_NET_CLIENT_LIST, _       = xprop.Atm(XU, "_NET_CLIENT_LIST")
	_NET_ACTIVE_WINDOW, _     = xprop.Atm(XU, "_NET_ACTIVE_WINDOW")
	ATOM_WINDOW_ICON, _       = xprop.Atm(XU, "_NET_WM_ICON")
	ATOM_WINDOW_NAME, _       = xprop.Atm(XU, "_NET_WM_NAME")
	ATOM_WINDOW_STATE, _      = xprop.Atm(XU, "_NET_WM_STATE")
	ATOM_WINDOW_TYPE, _       = xprop.Atm(XU, "_NET_WM_WINDOW_TYPE")
	ATOM_DOCK_APP_ID, _       = xprop.Atm(XU, "_DDE_DOCK_APP_ID")

	_NET_SYSTEM_TRAY_S0, _     = xprop.Atm(TrayXU, "_NET_SYSTEM_TRAY_S0")
	_NET_SYSTEM_TRAY_OPCODE, _ = xprop.Atm(TrayXU, "_NET_SYSTEM_TRAY_OPCODE")

	// ATOM_DEEPIN_WINDOW_VIEWPORTS, _ = xprop.Atm(XU, "DEEPIN_WINDOW_VIEWPORTS")
)
