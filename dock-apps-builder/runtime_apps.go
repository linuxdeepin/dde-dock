package main

import (
	"bytes"
	"dbus/com/deepin/daemon/dock"
	"dlib/gio-2.0"
	"dlib/glib-2.0"
	"encoding/base64"
	"fmt"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xprop"
	"github.com/BurntSushi/xgbutil/xwindow"
	"io/ioutil"
	"path"
	"strings"
)

var (
	XU, _                 = xgbutil.NewConn()
	_NET_CLIENT_LIST, _   = xprop.Atm(XU, "_NET_CLIENT_LIST")
	_NET_ACTIVE_WINDOW, _ = xprop.Atm(XU, "_NET_ACTIVE_WINDOW")
	ATOM_WINDOW_ICON, _   = xprop.Atm(XU, "_NET_WM_ICON")
	ATOM_WINDOW_NAME, _   = xprop.Atm(XU, "_NET_WM_NAME")
	ATOM_WINDOW_STATE, _  = xprop.Atm(XU, "_NET_WM_STATE")
	ATOM_WINDOW_TYPE, _   = xprop.Atm(XU, "_NET_WM_WINDOW_TYPE")
	ATOM_DOCK_APP_ID, _   = xprop.Atm(XU, "_DDE_DOCK_APP_ID")
)

var DOCKED_APP_MANAGER *dock.DockedAppManager

type WindowInfo struct {
	Xid   xproto.Window
	Title string
	Icon  string
}

type RuntimeApp struct {
	Id string
	//TODO: multiple xid window
	xids map[xproto.Window]*WindowInfo

	CurrentInfo *WindowInfo
	Menu        string
	coreMenu    *Menu

	exec string
	core *gio.DesktopAppInfo

	state     []string
	changedCB func()
}

func NewRuntimeApp(xid xproto.Window, appId string) *RuntimeApp {
	if !isNormalWindow(xid) {
		return nil
	}
	app := &RuntimeApp{
		Id:   strings.ToLower(appId),
		xids: make(map[xproto.Window]*WindowInfo),
	}
	app.core = gio.NewDesktopAppInfo(appId + ".desktop")
	if app.core == nil {
		if newId := guess_desktop_id(appId + ".desktop"); newId != "" {
			app.core = gio.NewDesktopAppInfo(newId)
		}
	}
	if app.core != nil {
		LOGGER.Debug(appId, ", Actions:", app.core.ListActions())
	} else {
		LOGGER.Debug(appId, ", Actions:[]")
	}
	app.attachXid(xid)
	app.CurrentInfo = app.xids[xid]
	app.getExec(xid)
	LOGGER.Debug("Exec:", app.exec)
	app.buildMenu()
	return app
}

func find_exec_name_by_xid(xid xproto.Window) string {
	pid, _ := ewmh.WmPidGet(XU, xid)
	return find_exec_name_by_pid(pid)
}
func (app *RuntimeApp) getExec(xid xproto.Window) {
	if app.core != nil {
		LOGGER.Debug(app.Id, " Get Exec from desktop file")
		// should NOT use GetExecuable, get wrong result, like skype
		// which gets 'env'.
		app.exec = app.core.GetString(glib.KeyFileDesktopKeyExec)
		return
	}
	LOGGER.Debug(app.Id, " Get Exec from pid")
	app.exec = find_exec_name_by_xid(xid)
}
func (app *RuntimeApp) buildMenu() {
	app.coreMenu = NewMenu()
	itemName := strings.Title(app.Id)
	if app.core != nil {
		itemName = strings.Title(app.core.GetDisplayName())
	}
	app.coreMenu.AppendItem(NewMenuItem(
		"_"+itemName,
		func() {
			var a *gio.AppInfo
			LOGGER.Info(itemName)
			if app.core != nil {
				LOGGER.Info("DesktopAppInfo")
				a = (*gio.AppInfo)(app.core)
			} else {
				LOGGER.Info("Non-DesktopAppInfo")
				a, err := gio.AppInfoCreateFromCommandline(
					app.exec,
					"",
					gio.AppInfoCreateFlagsNone,
				)
				if err != nil {
					LOGGER.Warning("Launch App Falied: ", err)
					return
				}

				defer a.Unref()
			}

			_, err := a.Launch(make([]*gio.File, 0), nil)
			LOGGER.Warning("Launch App Failed: ", err)
		},
		true,
	))
	app.coreMenu.AddSeparator()
	if app.core != nil {
		for _, actionName := range app.core.ListActions() {
			name := actionName //NOTE: don't directly use 'actionName' with closure in an forloop
			app.coreMenu.AppendItem(NewMenuItem(
				app.core.GetActionName(actionName),
				func() { app.core.LaunchAction(name, nil) },
				true,
			))
		}
		app.coreMenu.AddSeparator()
	}
	closeItem := NewMenuItem(
		"_Close All", // TODO: i18n
		func() {
			LOGGER.Warning("Close All")
			for xid := range app.xids {
				ewmh.CloseWindow(XU, xid)
			}
		},
		true,
	)
	app.coreMenu.AppendItem(closeItem)
	var err error
	if DOCKED_APP_MANAGER == nil {
		DOCKED_APP_MANAGER, err = dock.NewDockedAppManager(
			"com.deepin.daemon.Dock",
			"/dde/dock/DockedAppManager",
		)
		if err != nil {
			LOGGER.Warning("get DockedAppManager failed", err)
			return
		}
	}
	isDocked, err := DOCKED_APP_MANAGER.IsDocked(app.Id) // TODO: status
	if err != nil {
		LOGGER.Error("get docked status failed:", err)
	}
	LOGGER.Debug(app.Id, "Item is docked:", isDocked)
	dockItem := NewMenuItem(
		"_Dock",
		func() {
			LOGGER.Warning("dock item")
			LOGGER.Info("appid:", app.Id)

			var title, icon, exec string
			if app.core == nil {
				title = app.Id
				// TODO:
				icon = ""
				exec = app.exec
			} else {
				title = app.core.GetDisplayName()
				icon =
					get_theme_icon(app.core.GetIcon().ToString(),
						48)
				exec =
					app.core.GetString(glib.KeyFileDesktopKeyExec)
			}

			LOGGER.Info("id", app.Id, "title", title, "icon", icon,
				"exec", exec)
			_, err = DOCKED_APP_MANAGER.Dock(
				app.Id,
				title,
				icon,
				exec,
			)
			if err != nil {
				LOGGER.Error("Docked failed: ", err)
			}
			app.buildMenu()
		},
		!isDocked,
	)
	app.coreMenu.AppendItem(dockItem)

	app.Menu = app.coreMenu.GenerateJSON()
}

func (app *RuntimeApp) setChangedCB(cb func()) {
	app.changedCB = cb
}
func (app *RuntimeApp) notifyChanged() {
	if app.changedCB != nil {
		app.changedCB()
	}
}

func (app *RuntimeApp) HandleMenuItem(id int32) {
	if app.coreMenu != nil {
		app.coreMenu.HandleAction(id)
	}
}

//func find_app_id(pid uint, instanceName, wmName, wmClass, iconName string) string { return "" }

func find_app_id_by_xid(xid xproto.Window) string {
	if id, err := xprop.PropValStr(xprop.GetProperty(XU, xid, "_DDE_DOCK_APP_ID")); err == nil {
		return id
	}
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
	appId := find_app_id(pid, name, wmInstance, wmClassName, iconName)
	return appId
}

func contains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func isSkipTaskbar(xid xproto.Window) bool {
	state, err := ewmh.WmStateGet(XU, xid)
	if err != nil {
		return false
	}

	return contains(state, "_NET_WM_STATE_SKIP_TASKBAR")
}

func canBeMinimized(xid xproto.Window) bool {
	actions, err := ewmh.WmAllowedActionsGet(XU, xid)
	// LOGGER.Infof("%x: %v", xid, actions)
	if err != nil {
		return false
	}
	canBeMin := contains(actions, "_NET_WM_ACTION_MINIMIZE")
	// LOGGER.Infof("%x can be minimized: %v", xid, canBeMin)
	return canBeMin
}

var cannotBeDockedType []string = []string{
	"_NET_WM_WINDOW_TYPE_UTILITY",
	"_NET_WM_WINDOW_TYPE_COMBO",
	"_NET_WM_WINDOW_TYPE_DESKTOP",
	"_NET_WM_WINDOW_TYPE_DND",
	"_NET_WM_WINDOW_TYPE_DOCK",
	"_NET_WM_WINDOW_TYPE_DROPDOWN_MENU",
	"_NET_WM_WINDOW_TYPE_MENU",
	"_NET_WM_WINDOW_TYPE_NOTIFICATION",
	"_NET_WM_WINDOW_TYPE_POPUP_MENU",
	"_NET_WM_WINDOW_TYPE_SPLASH",
	"_NET_WM_WINDOW_TYPE_TOOLBAR",
	"_NET_WM_WINDOW_TYPE_TOOLTIP",
}

func isNormalWindow(xid xproto.Window) bool {
	// LOGGER.Debug("enter isNormalWindow:", xid)
	if wmClass, err := icccm.WmClassGet(XU, xid); err == nil {
		if wmClass.Instance == "explorer.exe" && wmClass.Class == "Wine" {
			return false
		} else if wmClass.Class == "DDELauncher" {
			// FIXME:
			// start_monitor_launcher_window like before?
			return false
		} else if wmClass.Class == "Desktop" {
			// FIXME:
			// get_desktop_pid like before?
			return false
		} else if wmClass.Class == "Dlock" {
			return false
		}
	}
	if isSkipTaskbar(xid) {
		return false
	}
	types, err := ewmh.WmWindowTypeGet(XU, xid)
	if err != nil {
		LOGGER.Debug("Get Window Type failed:", err)
		if _, err := xprop.GetProperty(XU, xid, "_XEMBED_INFO"); err != nil {
			return true
		} else {
			return false
		}
	}
	mayBeDocked := false
	cannotBeDoced := false
	for _, wmType := range types {
		if wmType == "_NET_WM_WINDOW_TYPE_NORMAL" ||
			(wmType == "_NET_WM_WINDOW_TYPE_DIALOG" &&
				canBeMinimized(xid)) {
			mayBeDocked = true
		} else if contains(cannotBeDockedType, wmType) {
			cannotBeDoced = true
		}
	}
	isNormal := mayBeDocked && !cannotBeDoced
	return isNormal
}

func (app *RuntimeApp) updateIcon(xid xproto.Window) {
	if app.core != nil {
		gioIcon := app.core.GetIcon()
		if gioIcon != nil {
			LOGGER.Debug("GetIcon:", gioIcon.ToString())
			icon := get_theme_icon(gioIcon.ToString(), 48)
			if icon != "" {
				LOGGER.Debug("get_theme_icon:", icon)
				// the path.Ext return ".xxx"
				ext := path.Ext(icon)[1:]
				LOGGER.Debug("ext:", ext)
				if strings.EqualFold(ext, "xpm") {
					LOGGER.Debug("change xmp to data uri")
					buf, err := ioutil.ReadFile(icon)
					if err != nil {
						app.xids[xid].Icon = "data:image/png;base64," +
							base64.StdEncoding.EncodeToString(buf)
						return
					}
				} else {
					app.xids[xid].Icon = icon
					return
				}
			}
		}
	}
	icon, err := xgraphics.FindIcon(XU, xid, 48, 48)
	if err == nil {
		buf := bytes.NewBuffer(nil)
		icon.WritePng(buf)
		app.xids[xid].Icon = "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
	} else {
		name, _ := ewmh.WmIconNameGet(XU, xid)
		app.xids[xid].Icon = name
	}
}
func (app *RuntimeApp) updateWmClass(xid xproto.Window) {
	if name, err := ewmh.WmNameGet(XU, xid); err == nil {
		app.xids[xid].Title = name
	}

}
func (app *RuntimeApp) updateState(xid xproto.Window) {
	//TODO: handle state
	app.state, _ = ewmh.WmStateGet(XU, xid)
}

func (app *RuntimeApp) updateAppid(xid xproto.Window) {
	if app.Id != find_app_id_by_xid(xid) {
		app.detachXid(xid)
		if newApp := MANAGER.createRuntimeApp(xid); newApp != nil {
			newApp.attachXid(xid)
		}
		fmt.Println("APP:", app.Id, "Changed to..", find_app_id_by_xid(xid))
		//TODO: Destroy
	}
}

func (app *RuntimeApp) Activate(x, y int32) error {
	//TODO: handle multiple xids
	switch {
	case !contains(app.state, "_NET_WM_STATE_FOCUSED"):
		ewmh.ActiveWindowReq(XU, app.CurrentInfo.Xid)
	case contains(app.state, "_NET_WM_STATE_FOCUSED"):
		s, err := icccm.WmStateGet(XU, app.CurrentInfo.Xid)
		if err != nil {
			LOGGER.Info("WmStateGetError:", s, err)
			return err
		}
		switch s.State {
		case icccm.StateIconic:
			s.State = icccm.StateNormal
			icccm.WmStateSet(XU, app.CurrentInfo.Xid, s)
		case icccm.StateNormal:
			activeXid, _ := ewmh.ActiveWindowGet(XU)
			if len(app.xids) == 1 {
				s.State = icccm.StateIconic
				iconifyWindow(app.CurrentInfo.Xid)
			} else {
				if activeXid == app.CurrentInfo.Xid {
					//ewmh.ActiveWindowReq(XU, app.findNextLeader())

					x := app.findNextLeader()
					ewmh.ActiveWindowReq(XU, x)
				}
			}
		}
	}
	return nil
}

func (app *RuntimeApp) setLeader(leader xproto.Window) {
	if info, ok := app.xids[leader]; ok {
		app.CurrentInfo = info
		app.notifyChanged()
	}
}

func (app *RuntimeApp) findNextLeader() xproto.Window {
	min := app.CurrentInfo

	var afterCurrent []*WindowInfo
	for _, xinfo := range app.xids {
		if xinfo.Xid > app.CurrentInfo.Xid {
			afterCurrent = append(afterCurrent, xinfo)
		}
		if xinfo.Xid < min.Xid {
			min = xinfo
		}
	}

	if len(afterCurrent) == 0 {
		return min.Xid
	} else {
		next := afterCurrent[0].Xid
		for _, xinfo := range afterCurrent {
			if next > xinfo.Xid {
				next = xinfo.Xid
			}
		}
		return next
	}
}

func iconifyWindow(xid xproto.Window) {
	ewmh.ClientEvent(XU, xid, "WM_CHANGE_STATE", icccm.StateIconic)
}

func (app *RuntimeApp) detachXid(xid xproto.Window) {
	if info, ok := app.xids[xid]; ok {
		xwindow.New(XU, xid).Listen(xproto.EventMaskNoEvent)
		xevent.Detach(XU, xid)

		if len(app.xids) == 1 {
			MANAGER.destroyRuntimeApp(app)
		} else {
			delete(app.xids, xid)
			if info == app.CurrentInfo {
				for _, nextInfo := range app.xids {
					if nextInfo != nil {
						app.CurrentInfo = nextInfo
						app.notifyChanged()
					} else {
						MANAGER.destroyRuntimeApp(app)
					}
					break
				}
			}
		}
	}
	if len(app.xids) == 0 {
		app.setChangedCB(nil)
	} else {
		app.notifyChanged()
	}
}

func (app *RuntimeApp) attachXid(xid xproto.Window) {
	if _, ok := app.xids[xid]; ok {
		return
	}
	xwin := xwindow.New(XU, xid)
	xwin.Listen(xproto.EventMaskPropertyChange | xproto.EventMaskStructureNotify | xproto.EventMaskVisibilityChange)
	winfo := &WindowInfo{Xid: xid}
	winfo.Title, _ = ewmh.WmNameGet(XU, xid)
	xevent.UnmapNotifyFun(func(XU *xgbutil.XUtil, ev xevent.UnmapNotifyEvent) {
		app.detachXid(xid)
	}).Connect(XU, xid)
	xevent.DestroyNotifyFun(func(XU *xgbutil.XUtil, ev xevent.DestroyNotifyEvent) {
		app.detachXid(xid)
	}).Connect(XU, xid)
	xevent.PropertyNotifyFun(func(XU *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
		switch ev.Atom {
		case ATOM_WINDOW_ICON:
			app.updateIcon(xid)
			app.updateAppid(xid)
			app.notifyChanged()
		case ATOM_WINDOW_NAME:
			app.updateWmClass(xid)
			app.updateAppid(xid)
			app.notifyChanged()
		case ATOM_WINDOW_STATE:
			app.updateState(xid)
			app.notifyChanged()
		case ATOM_WINDOW_TYPE:
			if !isNormalWindow(ev.Window) {
				app.detachXid(xid)
			}
		case ATOM_DOCK_APP_ID:
			app.updateAppid(xid)
			app.notifyChanged()
		}
	}).Connect(XU, xid)
	app.xids[xid] = winfo
	app.updateIcon(xid)
	app.updateWmClass(xid)
	app.updateState(xid)
	app.notifyChanged()
}

func listenRootWindow() {
	var update = func() {
		list, err := ewmh.ClientListGet(XU)
		if err != nil {
			LOGGER.Warning("Can't Get _NET_CLIENT_LIST", err)
		}
		MANAGER.runtimeAppChangged(list)
	}

	xwindow.New(XU, XU.RootWin()).Listen(xproto.EventMaskPropertyChange)
	xevent.PropertyNotifyFun(func(XU *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
		switch ev.Atom {
		case _NET_CLIENT_LIST:
			update()
		case _NET_ACTIVE_WINDOW:
			if activedWindow, err := ewmh.ActiveWindowGet(XU); err == nil {
				appId := find_app_id_by_xid(activedWindow)
				if rApp, ok := MANAGER.runtimeApps[appId]; ok {
					rApp.setLeader(activedWindow)
				}
			}
		}
	}).Connect(XU, XU.RootWin())
	update()
	xevent.Main(XU)
}
