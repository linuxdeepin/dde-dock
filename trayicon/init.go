package trayicon

import (
	"github.com/BurntSushi/xgb/composite"
	"github.com/BurntSushi/xgb/damage"
	"github.com/BurntSushi/xgb/xfixes"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xprop"
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
)

func init() {
	loader.Register(NewDaemon(logger))
}

var (
	XU                      *xgbutil.XUtil
	logger                  = log.NewLogger("daemon/trayicon")
	_NET_SYSTEM_TRAY_S0     xproto.Atom
	_NET_SYSTEM_TRAY_OPCODE xproto.Atom
	ATOM_MANAGER            xproto.Atom
)

func initX() {
	composite.Init(XU.Conn())
	composite.QueryVersion(XU.Conn(), 0, 4)
	damage.Init(XU.Conn())
	damage.QueryVersion(XU.Conn(), 1, 1)
	xfixes.Init(XU.Conn())
	xfixes.QueryVersion(XU.Conn(), 5, 0)

	_NET_SYSTEM_TRAY_S0, _ = xprop.Atm(XU, "_NET_SYSTEM_TRAY_S0")
	_NET_SYSTEM_TRAY_OPCODE, _ = xprop.Atm(XU, "_NET_SYSTEM_TRAY_OPCODE")
	ATOM_MANAGER, _ = xprop.Atm(XU, "MANAGER")
}
