/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dock

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"gir/gio-2.0"
	"gir/glib-2.0"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xprop"
	"github.com/BurntSushi/xgbutil/xwindow"
	. "pkg.deepin.io/lib/gettext"
)

type WindowInfo struct {
	Xid         xproto.Window
	Title       string
	Icon        string
	OverlapDock bool
}

type RuntimeApp struct {
	Id        string
	DesktopID string
	lock      sync.RWMutex
	//TODO: multiple xid window
	xids map[xproto.Window]*WindowInfo

	CurrentInfo *WindowInfo
	Menu        string
	coreMenu    *Menu

	exec string
	path string

	state       []string
	isHidden    bool
	isMaximized bool
	// workspaces  [][]uint

	changedCB func()

	updateConfigureTimer *time.Timer
	updateWMNameTimer    *time.Timer
}

func (app *RuntimeApp) updateWMName(xid xproto.Window) {
	app.cancelUpdateWMName()
	app.updateWMNameTimer = time.AfterFunc(time.Millisecond*20, func() {
		app.updateWmName(xid)
		app.updateAppid(xid)
		app.notifyChanged()
		app.updateWMNameTimer = nil
	})
}

func (app *RuntimeApp) cancelUpdateWMName() {
	if app.updateWMNameTimer != nil {
		app.updateWMNameTimer.Stop()
		app.updateWMNameTimer = nil
	}
}

func (app *RuntimeApp) createDesktopAppInfo() *DesktopAppInfo {
	core := NewDesktopAppInfo(app.DesktopID)

	if core != nil {
		return core
	}

	if newId := guess_desktop_id(app.Id); newId != "" {
		core = NewDesktopAppInfo(newId)
		if core != nil {
			return core
		}
	}

	return NewDesktopAppInfoFromFilename(app.path)
}

func NewRuntimeApp(xid xproto.Window, appId string) *RuntimeApp {
	recordFrequency(appId)
	markAsLaunched(appId)

	if !isNormalWindow(xid) {
		return nil
	}
	logger.Info("NewRuntimeApp for", appId)
	app := &RuntimeApp{
		Id:        strings.ToLower(appId),
		DesktopID: appId + ".desktop",
		xids:      make(map[xproto.Window]*WindowInfo),
	}
	core := NewDesktopAppInfo(app.DesktopID)
	if core == nil {
		if newId := guess_desktop_id(app.Id); newId != "" {
			core = NewDesktopAppInfo(newId)
			if core != nil {
				app.DesktopID = newId
			}
		}
	}

	if core != nil {
		logger.Debug(appId, ", Actions:", core.ListActions())
		app.path = core.GetFilename()
		core.Unref()
	} else {
		logger.Debug(appId, ", Actions:[]")
	}
	app.attachXid(xid)
	app.CurrentInfo = app.xids[xid]
	app.getExec(xid)
	logger.Debug("Exec:", app.exec)
	app.buildMenu()
	return app
}

func find_exec_by_xid(xid xproto.Window) string {
	pid, _ := ewmh.WmPidGet(XU, xid)
	if pid == 0 {
		name, _ := ewmh.WmNameGet(XU, xid)
		if name != "" {
			pid = lookthroughProc(name)
		}
	}
	return find_exec_by_pid(pid)
}

func (app *RuntimeApp) getExec(xid xproto.Window) {
	core := app.createDesktopAppInfo()
	if core != nil {
		logger.Debug(app.Id, "Get Exec from desktop file")
		// should NOT use GetExecuable, get wrong result, like skype
		// which gets 'env'.
		app.exec = core.DesktopAppInfo.GetString(glib.KeyFileDesktopKeyExec)
		core.Unref()
		return
	}
	logger.Debug(app.Id, "Get Exec from pid")
	app.exec = find_exec_by_xid(xid)
	logger.Debug("app get exec:", app.exec)
}

func actionGenarator(id string) func(uint32) {
	return func(timestamp uint32) {
		app, ok := ENTRY_MANAGER.runtimeApps[id]
		if !ok {
			return
		}
		logger.Debug("dock item")
		logger.Debug("appid:", app.Id)

		var title, icon, exec string
		core := app.createDesktopAppInfo()
		if core == nil {
			icon = "application-default-icon"
			title = app.Id
			for _, v := range app.xids {
				if v.Icon == "" {
					continue
				}

				if strings.HasPrefix(v.Icon, "data:image") {
					path, err := dataUriToFile(
						v.Icon,
						filepath.Join(scratchDir,
							app.Id+".png"),
					)
					if err != nil {
						logger.Debug("dock", app.Id, "dataUriToFile failed", err)
						break
					}
					icon = path
				} else {
					icon = v.Icon
				}
				break
			}
			execFile := filepath.Join(
				scratchDir,
				app.Id+".sh",
			)
			shExec := "#!/usr/bin/env bash\n\n" + app.exec
			ioutil.WriteFile(execFile, []byte(shExec), 0744)
			exec = execFile
		} else {
			defer core.Unref()
			title = core.GetDisplayName()
			icon = get_theme_icon(core.GetIcon().ToString(), 48)
			exec = core.DesktopAppInfo.GetString(glib.KeyFileDesktopKeyExec)
		}

		logger.Debug("id", app.Id, "title", title, "icon", icon,
			"exec", exec)
		DOCKED_APP_MANAGER.Dock(
			app.Id,
			title,
			icon,
			exec,
		)
	}
}

func (app *RuntimeApp) buildMenu() {
	app.coreMenu = NewMenu()
	itemName := strings.Title(app.Id)
	core := app.createDesktopAppInfo()
	if core != nil {
		itemName = strings.Title(core.GetDisplayName())
		defer core.Unref()
	}
	app.coreMenu.AppendItem(NewMenuItem(
		itemName,
		func(timestamp uint32) {
			var a *gio.AppInfo
			logger.Debug(itemName)
			core := app.createDesktopAppInfo()
			if core != nil {
				logger.Debug("DesktopAppInfo")
				a = (*gio.AppInfo)(core.DesktopAppInfo)
				defer core.Unref()
			} else {
				logger.Debug("Non-DesktopAppInfo", app.exec)
				var err error = nil
				a, err = gio.AppInfoCreateFromCommandline(
					app.exec,
					"",
					gio.AppInfoCreateFlagsNone,
				)
				if err != nil {
					logger.Warning("Launch App Falied: ", err)
					return
				}

				defer a.Unref()
			}

			if a == nil {
				logger.Warning("create app info to run program failed")
				return
			}

			_, err := a.Launch(make([]*gio.File, 0), gio.GetGdkAppLaunchContext().SetTimestamp(timestamp))
			if err != nil {
				logger.Warning("Launch App Failed: ", err)
			}
		},
		true,
	))
	app.coreMenu.AddSeparator()
	if core != nil {
		for _, actionName := range core.ListActions() {
			name := actionName //NOTE: don't directly use 'actionName' with closure in an forloop
			app.coreMenu.AppendItem(NewMenuItem(
				core.GetActionName(actionName),
				func(timestamp uint32) {
					core := app.createDesktopAppInfo()
					if core == nil {
						return
					}
					defer core.Unref()
					core.LaunchAction(name, gio.GetGdkAppLaunchContext().SetTimestamp(timestamp))
				},
				true,
			))
		}
		app.coreMenu.AddSeparator()
	}
	closeItem := NewMenuItem(
		Tr("_Close All"),
		func(timestamp uint32) {
			logger.Debug("Close All")
			for xid := range app.xids {
				ewmh.CloseWindow(XU, xid)
			}
		},
		true,
	)
	app.coreMenu.AppendItem(closeItem)
	isDocked := DOCKED_APP_MANAGER.IsDocked(app.Id)
	logger.Info(app.Id, "Item is docked:", isDocked)
	var message string = ""
	var action func(uint32) = nil
	if isDocked {
		logger.Info(app.Id, "change to undock")
		message = Tr("_Undock")
		action = func(id string) func(uint32) {
			return func(uint32) {
				app, ok := ENTRY_MANAGER.runtimeApps[id]
				if !ok {
					return
				}
				DOCKED_APP_MANAGER.Undock(app.Id)
			}
		}(app.Id)
	} else {
		logger.Info(app.Id, "change to dock")
		message = Tr("_Dock")
		action = actionGenarator(app.Id)
	}

	logger.Debug(app.Id, "New Menu Item:", message)
	dockItem := NewMenuItem(message, action, true)
	app.coreMenu.AppendItem(dockItem)

	app.Menu = app.coreMenu.GenerateJSON()
	app.notifyChanged()
}

func (app *RuntimeApp) setChangedCB(cb func()) {
	app.changedCB = cb
}
func (app *RuntimeApp) notifyChanged() {
	if app.changedCB != nil {
		app.changedCB()
	}
}

func (app *RuntimeApp) HandleMenuItem(id string, timestamp uint32) {
	if app.coreMenu != nil {
		app.coreMenu.HandleAction(id, timestamp)
	}
}

//func find_app_id(pid uint, instanceName, wmName, wmClass, iconName string) string { return "" }

// FIXME:
// to fix mine craft, however, IDLE will get problem.
func lookthroughProc(name string) uint {
	name = strings.ToLower(strings.Split(name, " ")[0])

	dirs, err := ioutil.ReadDir("/proc")
	if err != nil {
		logger.Debug("read /proc dir failed:", err)
		return 0
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		content, err := ioutil.ReadFile(filepath.Join("/proc", dir.Name(), "cmdline"))
		if err != nil {
			logger.Debug("read cmdline failed:", err)
			continue
		}

		if strings.Contains(strings.ToLower(string(content)), name) {
			pid, err := strconv.Atoi(dir.Name())
			if err != nil {
				logger.Debug("change string to int failed:",
					err)
				return 0
			}
			return uint(pid)
		}
	}

	return 0
}

func find_app_id_by_xid(xid xproto.Window, displayMode DisplayModeType) string {
	var appId string
	if displayMode == DisplayModeModernMode {
		if id, err := xprop.PropValStr(xprop.GetProperty(XU, xid, "_DDE_DOCK_APP_ID")); err == nil {
			appId = getAppIDFromDesktopID(normalizeAppID(id))
			if appId != "" {
				logger.Info("get app id from _DDE_DOCK_APP_ID", appId)
				return appId
			}
		}
	}

	gtkAppId, err := xprop.PropValStr(xprop.GetProperty(XU, xid, "_GTK_APPLICATION_ID"))
	if err != nil {
		logger.Debug("get AppId from _GTK_APPLICATION_ID failed:", err)
	} else {
		appId = gtkAppId
		appId = getAppIDFromDesktopID(normalizeAppID(appId))
		if appId != "" {
			return appId
		}
	}

	appId = getAppIDFromXid(xid)
	if appId != "" {
		logger.Debug("get app id from bamf", appId)
		return normalizeAppID(appId)
	}

	wmClass, _ := icccm.WmClassGet(XU, xid)
	var wmInstance, wmClassName string
	if wmClass != nil {
		if utf8.ValidString(wmClass.Instance) {
			wmInstance = wmClass.Instance
		}

		// it is possible that getting invalid string which might be xgb implementation's bug.
		// for instance: xdemineur's WMClass
		if utf8.ValidString(wmClass.Class) {
			wmClassName = wmClass.Class
		}
		logger.Debug("WMClass", wmClassName, ", WMInstance", wmInstance)
	}
	name, _ := ewmh.WmNameGet(XU, xid)
	pid, err := ewmh.WmPidGet(XU, xid)
	if err != nil {
		logger.Debug("get pid failed for:", xid)
		if name != "" {
			pid = lookthroughProc(name)
		} else {
			newAppId := getAppIDFromDesktopID(normalizeAppID(wmClassName))
			if newAppId != "" {
				logger.Debug("get Pid failed, using wm class name as app id", newAppId)
				return newAppId
			}
		}
	}

	iconName, _ := ewmh.WmIconNameGet(XU, xid)
	if pid == 0 {
		appId = normalizeAppID(wmClassName)
		logger.Debug("get window name failed, using wm class name as app id", appId)
		return appId
	} else {
	}
	appId = find_app_id(pid, name, wmInstance, wmClassName, iconName)
	newAppId := getAppIDFromDesktopID(normalizeAppID(appId))
	if newAppId != "" {
		appId = newAppId
	}

	appId = specialCaseWorkaround(xid, appId)

	logger.Debugf("get appid %q", appId)
	return appId
}

func specialCaseWorkaround(xid xproto.Window, appId string) string {
	switch appId {
	case "xwalk":
		desktopID, err := xprop.PropValStr(xprop.GetProperty(XU, xid, "_NET_WM_DESKTOP_FILE"))
		if err != nil {
			logger.Debug("get xwalk AppId from _NET_WM_DESKTOP_FILE failed:", err)
		} else {
			appId = trimDesktop(filepath.Base(desktopID))
			return appId
		}
	}
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
	// logger.Infof("%x: %v", xid, actions)
	if err != nil {
		return false
	}
	canBeMin := contains(actions, "_NET_WM_ACTION_MINIMIZE")
	// logger.Infof("%x can be minimized: %v", xid, canBeMin)
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
	winProps, err := xproto.GetWindowAttributes(XU.Conn(), xid).Reply()
	if err != nil {
		logger.Debug("faild Get WindowAttributes:", xid, err)
		return false
	}
	if winProps.MapState != xproto.MapStateViewable {
		return false
	}
	// logger.Debug("enter isNormalWindow:", xid)
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
		logger.Debug("Get Window Type failed:", err)
		if _, err := xprop.GetProperty(XU, xid, "_XEMBED_INFO"); err != nil {
			logger.Debug("has _XEMBED_INFO")
			return true
		} else {
			logger.Debug("Get _XEMBED_INFO failed:", err)
			return false
		}
	}

	mayBeDocked := false
	cannotBeDoced := false
	for _, wmType := range types {
		if wmType == "_NET_WM_WINDOW_TYPE_NORMAL" {
			mayBeDocked = true
		} else if wmType == "_NET_WM_WINDOW_TYPE_DIALOG" {
			if !canBeMinimized(xid) {
				mayBeDocked = false
				break
			} else {
				mayBeDocked = true
			}
		} else if contains(cannotBeDockedType, wmType) {
			cannotBeDoced = true
		}
	}

	logger.Debug("mayBeDocked:", mayBeDocked, "cannotBeDoced:", cannotBeDoced)
	isNormal := mayBeDocked && !cannotBeDoced
	return isNormal
}

func (app *RuntimeApp) updateIcon(xid xproto.Window) {
	logger.Debug("update icon for", app.Id)
	core := app.createDesktopAppInfo()
	if core != nil {
		defer core.Unref()
		icon := getAppIcon(core.DesktopAppInfo)
		if icon != "" {
			app.xids[xid].Icon = icon
			return
		}
		logger.Warning("get icon from app failed")
	} else {
		logger.Warningf("create desktop app info for %s(window id: 0x%x) failed:", app.DesktopID, xid)
	}

	logger.Info(app.Id, "using icon from X")
	icon, err := xgraphics.FindIcon(XU, xid, 48, 48)
	// logger.Info(icon, err)
	// FIXME: gets empty icon for minecraft
	if err == nil {
		buf := bytes.NewBuffer(nil)
		icon.WritePng(buf)
		app.xids[xid].Icon = "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
		return
	}

	logger.Debug("get icon from X failed:", err)
	logger.Debug("get icon name from _NET_WM_ICON_NAME")
	name, _ := ewmh.WmIconNameGet(XU, xid)
	if name != "" {
		app.xids[xid].Icon = name
	} else {
		app.xids[xid].Icon = "application-default-icon"
	}

}
func (app *RuntimeApp) updateWmName(xid xproto.Window) {
	if _, ok := app.xids[xid]; !ok {
		return
	}
	if name, err := ewmh.WmNameGet(XU, xid); err == nil && name != "" {
		if utf8.ValidString(name) {
			app.xids[xid].Title = name
			return
		}
	}

	if name, err := xprop.PropValStr(xprop.GetProperty(XU, xid,
		"WM_NAME")); err == nil {
		if utf8.ValidString(name) {
			app.xids[xid].Title = name
			return
		}
	}

	if app.xids[xid].Title == "" {
		app.xids[xid].Title = app.Id
	}
}

func (app *RuntimeApp) updateState(xid xproto.Window) {
	if _, ok := app.xids[xid]; !ok {
		return
	}
	logger.Debugf("get state of %s(0x%x)", app.Id, xid)
	//TODO: handle state
	var err error
	app.state, err = ewmh.WmStateGet(XU, xid)
	if err != nil {
		logger.Warningf("get state of %s(0x%x) failed: %s", app.Id, xid, err)
	}
	logger.Debugf("state of %s(0x%x) is %v", app.Id, xid, app.state)
	app.isHidden = contains(app.state, "_NET_WM_STATE_HIDDEN")
	app.isMaximized = contains(app.state, "_NET_WM_STATE_MAXIMIZED_VERT")
}

// TODO: using this instead of walking throught all client
// to get the workspaces
// func (app *RuntimeApp) updateViewports(xid xproto.Window) {
// 	app.workspaces = nil
// 	viewports, err := xprop.PropValNums(xprop.GetProperty(XU, xid,
// 		"DEEPIN_WINDOW_VIEWPORTS"))
// 	if err != nil {
// 		logger.Error("get DEEPIN_WINDOW_VIEWPORTS failed", err)
// 		return
// 	}
// 	app.workspaces = make([][]uint, 0)
// 	for i := uint(0); i < viewports[0]; i++ {
// 		viewport := make([]uint, 0)
// 		viewport[0] = viewports[i+1]
// 		viewport[1] = viewports[i+2]
// 		app.workspaces = append(app.workspaces, viewport)
// 	}
// }

func (app *RuntimeApp) updateAppid(xid xproto.Window) {
	newAppId := find_app_id_by_xid(
		xid,
		DisplayModeType(setting.GetDisplayMode()),
	)
	if app.Id != newAppId {
		app.detachXid(xid)
		if newApp := ENTRY_MANAGER.createRuntimeApp(xid); newApp != nil {
			newApp.attachXid(xid)
		}
		logger.Debug("APP:", app.Id, "Changed to..", newAppId)
		//TODO: Destroy
	}
}

func (app *RuntimeApp) Activate(x, y int32, timestamp uint32) error {
	//TODO: handle multiple xids
	switch {
	case !contains(app.state, "_NET_WM_STATE_FOCUSED"):
		activateWindow(app.CurrentInfo.Xid)
	case contains(app.state, "_NET_WM_STATE_FOCUSED"):
		s, err := icccm.WmStateGet(XU, app.CurrentInfo.Xid)
		if err != nil {
			logger.Info("WmStateGetError:", s, err)
			return err
		}
		switch s.State {
		case icccm.StateIconic:
			s.State = icccm.StateNormal
			icccm.WmStateSet(XU, app.CurrentInfo.Xid, s)
		case icccm.StateNormal:
			activeXid, _ := ewmh.ActiveWindowGet(XU)
			logger.Debugf("%s, 0x%x(c), 0x%x(a), %v", app.Id,
				app.CurrentInfo.Xid, activeXid, app.state)
			if len(app.xids) == 1 {
				s.State = icccm.StateIconic
				iconifyWindow(app.CurrentInfo.Xid)
			} else {
				logger.Debugf("activeXid is 0x%x, current is 0x%x", activeXid,
					app.CurrentInfo.Xid)
				if activeXid == app.CurrentInfo.Xid {
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
			ENTRY_MANAGER.destroyRuntimeApp(app)
		} else {
			delete(app.xids, xid)
			if info == app.CurrentInfo {
				for _, nextInfo := range app.xids {
					if nextInfo != nil {
						app.CurrentInfo = nextInfo
						app.updateState(app.CurrentInfo.Xid)
						app.notifyChanged()
					} else {
						ENTRY_MANAGER.destroyRuntimeApp(app)
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

func (app *RuntimeApp) updateOverlap(xid xproto.Window) {
	if _, ok := app.xids[xid]; ok {
		logger.Debug(app.Id, "isHidden:", isHiddenPre(xid), "isOnCurrentWorkspace:", onCurrentWorkspacePre(xid), "isOverDock:", isWindowOverlapDock(xid))
		overlap := !isHiddenPre(xid) && onCurrentWorkspacePre(xid) && isWindowOverlapDock(xid)
		if overlap != app.xids[xid].OverlapDock {
			app.xids[xid].OverlapDock = overlap
			hideModemanager.UpdateState()
		}
	}
}

func (app *RuntimeApp) attachXid(xid xproto.Window) {
	logger.Debugf("attach 0x%x to %s", xid, app.Id)
	if _, ok := app.xids[xid]; ok {
		logger.Debugf("0x%x is already on %s", xid, app.Id)
		return
	}
	xwin := xwindow.New(XU, xid)
	xwin.Listen(xproto.EventMaskPropertyChange | xproto.EventMaskStructureNotify | xproto.EventMaskVisibilityChange)
	winfo := &WindowInfo{Xid: xid}
	xevent.UnmapNotifyFun(func(XU *xgbutil.XUtil, ev xevent.UnmapNotifyEvent) {
		app.detachXid(xid)
	}).Connect(XU, xid)
	xevent.DestroyNotifyFun(func(XU *xgbutil.XUtil, ev xevent.DestroyNotifyEvent) {
		app.detachXid(xid)
	}).Connect(XU, xid)
	xevent.PropertyNotifyFun(func(XU *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
		app.lock.Lock()
		defer app.lock.Unlock()
		switch ev.Atom {
		case ATOM_WINDOW_ICON:
			app.updateIcon(xid)
			app.updateAppid(xid)
			app.notifyChanged()
		case ATOM_WINDOW_NAME:
			app.updateWMName(xid)
		case ATOM_WINDOW_STATE:
			logger.Debugf("%s(0x%x) WM_STATE is changed", app.Id, xid)
			if app.CurrentInfo.Xid == xid {
				logger.Debug("is current window info changed")
				app.updateState(xid)
			}
			app.notifyChanged()

			if HideModeType(setting.GetHideMode()) != HideModeSmartHide {
				break
			}

			time.AfterFunc(time.Millisecond*20, func() {
				app.updateOverlap(xid)
			})
		// case ATOM_DEEPIN_WINDOW_VIEWPORTS:
		// 	app.updateViewports(xid)
		case ATOM_WINDOW_TYPE:
			if !isNormalWindow(ev.Window) {
				app.detachXid(xid)
			}
		case ATOM_DOCK_APP_ID:
			app.updateAppid(xid)
			app.notifyChanged()
		}
	}).Connect(XU, xid)
	update := func(xid xproto.Window) {
		app.lock.Lock()
		defer app.lock.Unlock()
		app.updateOverlap(xid)
	}
	xevent.ConfigureNotifyFun(func(XU *xgbutil.XUtil, ev xevent.ConfigureNotifyEvent) {
		app.lock.Lock()
		defer app.lock.Unlock()
		if app.updateConfigureTimer != nil {
			app.updateConfigureTimer.Stop()
			app.updateConfigureTimer = nil
		}
		app.updateConfigureTimer = time.AfterFunc(time.Millisecond*20, func() {
			update(ev.Window)
			app.updateConfigureTimer = nil
		})
	}).Connect(XU, xid)
	app.xids[xid] = winfo
	update(xid)
	app.updateIcon(xid)
	app.updateWmName(xid)
	app.updateState(xid)
	// app.updateViewports(xid)
	app.notifyChanged()
}
