package main

import "github.com/BurntSushi/xgb/xproto"
import "github.com/BurntSushi/xgbutil/ewmh"
import "bytes"
import "encoding/base64"
import "github.com/BurntSushi/xgbutil"
import "github.com/BurntSushi/xgbutil/icccm"
import "github.com/BurntSushi/xgbutil/xwindow"
import "github.com/BurntSushi/xgbutil/xevent"
import "github.com/BurntSushi/xgbutil/xgraphics"

type RuntimeApp struct {
	Id string
	//TODO: multiple xid window
	xwin  *xwindow.Window
	Title string
	Icon  string

	state     []string
	changedCB func()
}

func (app *RuntimeApp) setChangedCB(cb func()) {
	app.changedCB = cb
}
func (app *RuntimeApp) notifyChanged() {
	if app.changedCB != nil {
		app.changedCB()
	}
}

//func find_app_id(pid uint, instanceName, wmName, wmClass, iconName string) string { return "" }

func find_app_id_by_xid(xid xproto.Window) string {
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
	return find_app_id(pid, name, wmInstance, wmClassName, iconName)
}

func contains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func isNormalWindow(xid xproto.Window) bool {
	types, _ := ewmh.WmWindowTypeGet(XU, xid)
	if contains(types, "_NET_WM_WINDOW_TYPE_NORMAL") {
		return true
	} else {
		return false
	}
}

func NewRuntimeApp(xid xproto.Window, appId string) *RuntimeApp {
	if !isNormalWindow(xid) {
		return nil
	}

	app := &RuntimeApp{Id: appId}
	app.xwin = xwindow.New(XU, xid)
	app.xwin.Listen(xproto.EventMaskPropertyChange | xproto.EventMaskStructureNotify | xproto.EventMaskVisibilityChange)
	app.Title, _ = ewmh.WmNameGet(XU, xid)
	xevent.DestroyNotifyFun(func(XU *xgbutil.XUtil, ev xevent.DestroyNotifyEvent) {
		MANAGER.destroyRuntimeApp(app)
	}).Connect(XU, xid)
	xevent.PropertyNotifyFun(func(XU *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
		switch ev.Atom {
		case ATOM_WINDOW_ICON:
			app.updateIcon()
		case ATOM_WINDOW_NAME:
			app.updateWmClass()
		case ATOM_WINDOW_STATE:
			app.updateState()
		case ATOM_WINDOW_TYPE:
			if !isNormalWindow(ev.Window) {
				MANAGER.destroyRuntimeApp(app)
			}
		default:
			return
		}
		app.notifyChanged()
	}).Connect(XU, xid)
	app.updateIcon()
	app.updateWmClass()
	app.updateState()
	app.notifyChanged()
	return app
}
func (app *RuntimeApp) updateIcon() {
	icon, err := xgraphics.FindIcon(XU, app.xwin.Id, 48, 48)
	if err == nil {
		buf := bytes.NewBuffer(nil)
		icon.WritePng(buf)
		app.Icon = "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
	} else {
		name, _ := ewmh.WmIconNameGet(XU, app.xwin.Id)
		app.Icon = name
	}
}
func (app *RuntimeApp) updateWmClass() {
	if name, err := ewmh.WmNameGet(XU, app.xwin.Id); err == nil {
		app.Title = name
	}

}
func (app *RuntimeApp) updateState() {
	app.state, _ = ewmh.WmStateGet(XU, app.xwin.Id)
}

func (app *RuntimeApp) Activate(x, y int32) {
	switch {
	case !contains(app.state, "_NET_WM_STATE_FOCUSED"):
		ewmh.ActiveWindowSet(XU, app.xwin.Id)
	case contains(app.state, "_NET_WM_STATE_FOCUSED"):
		s, _ := icccm.WmStateGet(XU, app.xwin.Id)
		switch s.State {
		case icccm.StateIconic:
			s.State = icccm.StateNormal
			icccm.WmStateSet(XU, app.xwin.Id, s)
		case icccm.StateNormal:
			activeXid, _ := ewmh.ActiveWindowGet(XU)
			if activeXid == app.xwin.Id {
				s.State = icccm.StateIconic
				icccm.WmStateSet(XU, app.xwin.Id, s)
			}
		}

	}
}
