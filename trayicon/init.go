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
	TrayXU                  *xgbutil.XUtil
	logger                  = log.NewLogger("daemon/trayicon")
	_NET_SYSTEM_TRAY_S0     xproto.Atom
	_NET_SYSTEM_TRAY_OPCODE xproto.Atom
)

func initX() {
	composite.Init(TrayXU.Conn())
	composite.QueryVersion(TrayXU.Conn(), 0, 4)
	damage.Init(TrayXU.Conn())
	damage.QueryVersion(TrayXU.Conn(), 1, 1)
	xfixes.Init(TrayXU.Conn())
	xfixes.QueryVersion(TrayXU.Conn(), 5, 0)

	_NET_SYSTEM_TRAY_S0, _ = xprop.Atm(TrayXU, "_NET_SYSTEM_TRAY_S0")
	_NET_SYSTEM_TRAY_OPCODE, _ = xprop.Atm(TrayXU, "_NET_SYSTEM_TRAY_OPCODE")
}
