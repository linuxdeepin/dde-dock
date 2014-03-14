package main

import "fmt"
import "github.com/BurntSushi/xgbutil"
import "github.com/BurntSushi/xgb/xproto"
import "github.com/BurntSushi/xgbutil/ewmh"
import "github.com/BurntSushi/xgbutil/icccm"

var XU, _ = xgbutil.NewConn()

type RuntimeApp struct {
	xid   xproto.Window
	Title string
}

func contains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func NewRuntimeApp(xid xproto.Window) *RuntimeApp {
	types, _ := ewmh.WmWindowTypeGet(XU, xid)
	if !contains(types, "_NET_WM_WINDOW_TYPE_NORMAL") {
		return nil
	}

	app := &RuntimeApp{}
	app.xid = xid
	app.Title, _ = ewmh.WmNameGet(XU, xid)
	app.update_appid()
	return app
}

func (app *RuntimeApp) update_appid() {
	pid, _ := ewmh.WmPidGet(XU, app.xid)
	iconName, _ := ewmh.WmIconNameGet(XU, app.xid)
	name, _ := ewmh.WmNameGet(XU, app.xid)
	wmClass, _ := icccm.WmClassGet(XU, app.xid)
	var wmInstance, wmClassName string
	if wmClass != nil {
		wmInstance = wmClass.Instance
		wmClassName = wmClass.Class
	}
	appid := find_app_id(pid, name, wmInstance, wmClassName, iconName)
	fmt.Println("INFO:", pid, name, wmInstance, wmClassName, iconName, "RESULT:", appid)
}

func monitor() {
	list, _ := ewmh.ClientListGet(XU)
	for _, xid := range list {
		if app := NewRuntimeApp(xid); app != nil {
			fmt.Println(app)
		}
	}
}

func init() {
	monitor()
}
