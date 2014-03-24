package main

import "github.com/BurntSushi/xgbutil"
import "github.com/BurntSushi/xgbutil/xwindow"
import "github.com/BurntSushi/xgb/xproto"

var (
	XU, _ = xgbutil.NewConn()
)

func hideWindow(xid xproto.Window) {
	xwindow.New(XU, xid).Unmap()
}
