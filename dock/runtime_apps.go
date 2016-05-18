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
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"gir/gio-2.0"
	"gir/glib-2.0"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/xprop"
	. "pkg.deepin.io/lib/gettext"
	"sort"
)

type RuntimeApp struct {
	Id              string
	DesktopID       string
	lock            sync.RWMutex
	windowInfoTable map[xproto.Window]*WindowInfo

	CurrentInfo *WindowInfo
	Menu        string
	coreMenu    *Menu

	exec string
	path string

	appInfo   *AppInfo
	changedCB func()
}

// TODO: 移除此函数
func (app *RuntimeApp) createDesktopAppInfo() *AppInfo {
	core := NewAppInfo(app.DesktopID)

	if core != nil {
		return core
	}

	if newId := guess_desktop_id(app.Id); newId != "" {
		core = NewAppInfo(newId)
		if core != nil {
			return core
		}
	}

	return NewAppInfoFromFile(app.path)
}

func NewRuntimeApp(winInfo *WindowInfo, appInfo *AppInfo, isAppDocked bool) *RuntimeApp {
	appId := appInfo.GetId()
	recordFrequency(appId)
	markAsLaunched(appId)

	logger.Info("NewRuntimeApp for", appId)
	app := &RuntimeApp{
		Id:              strings.ToLower(appId),
		DesktopID:       appId + ".desktop",
		windowInfoTable: make(map[xproto.Window]*WindowInfo),
		appInfo:         appInfo,
	}
	app.path = app.appInfo.GetFilename()
	app.initExec(winInfo)

	app.attachWindow(winInfo)
	app.CurrentInfo = winInfo
	app.updateMenu(isAppDocked)
	app.notifyChanged()
	return app
}

func (app *RuntimeApp) initExec(winInfo *WindowInfo) {
	ai := app.appInfo
	if ai.DesktopAppInfo != nil {
		// NOTE: should NOT use GetExecuable, get wrong result
		// like skype which gets 'env'.
		// TODO: ai.getExec() or ai.getCmdline()
		app.exec = ai.DesktopAppInfo.GetString(glib.KeyFileDesktopKeyExec)
	}

	if winInfo.process != nil {
		app.exec = winInfo.process.GetShellScript()
	}

	logger.Debug("initExec:", app.exec)
}

// TODO: 暂时不删
// func find_exec_by_xid(xid xproto.Window) string {
// 	pid, _ := ewmh.WmPidGet(XU, xid)
// 	if pid == 0 {
// 		name, _ := ewmh.WmNameGet(XU, xid)
// 		if name != "" {
// 			pid = lookthroughProc(name)
// 		}
// 	}
// 	return find_exec_by_pid(pid)
// }

func (app *RuntimeApp) updateMenu(isDocked bool) {
	logger.Debug("Update menu")
	menu := NewMenu()
	menu.AppendItem(app.getMenuItemLaunch())

	desktopActionMenuItems := app.getMenuItemDesktopActions()
	menu.AppendItem(desktopActionMenuItems...)

	menu.AppendItem(app.getMenuItemCloseAll())

	// menu item dock or undock
	logger.Info(app.Id, "Item docked?", isDocked)
	if isDocked {
		menu.AppendItem(app.getMenuItemUndock())
	} else {
		menu.AppendItem(app.getMenuItemDock())
	}

	app.coreMenu = menu
	app.Menu = app.coreMenu.GenerateJSON()
}

func (app *RuntimeApp) getMenuItemDesktopActions() []*MenuItem {
	if app.appInfo == nil {
		return nil
	}

	var menuItems []*MenuItem
	for _, actionName := range app.appInfo.ListActions() {
		//NOTE: don't directly use 'actionName' with closure in an forloop
		actionNameCopy := actionName
		menuItem := NewMenuItem(
			app.appInfo.GetActionName(actionName),
			func(timestamp uint32) {
				logger.Debug("desktop app info launch action:", actionNameCopy)
				app.appInfo.LaunchAction(actionNameCopy,
					gio.GetGdkAppLaunchContext().SetTimestamp(timestamp))
			}, true)
		menuItems = append(menuItems, menuItem)
	}
	return menuItems
}

func (app *RuntimeApp) launchApp(timestamp uint32) {
	var appInfo *gio.AppInfo
	if app.appInfo.DesktopAppInfo != nil {
		logger.Debug("Has AppInfo")
		appInfo = (*gio.AppInfo)(app.appInfo.DesktopAppInfo)
	} else {
		logger.Debug("No AppInfo", app.exec)
		var err error = nil
		appInfo, err = gio.AppInfoCreateFromCommandline(
			app.exec,
			"",
			gio.AppInfoCreateFlagsNone,
		)
		if err != nil {
			logger.Warning("Launch App Falied: ", err)
			return
		}

		defer appInfo.Unref()
	}

	if appInfo == nil {
		logger.Warning("create app info to run program failed")
		return
	}

	_, err := appInfo.Launch(make([]*gio.File, 0), gio.GetGdkAppLaunchContext().SetTimestamp(timestamp))
	if err != nil {
		logger.Warning("Launch App Failed: ", err)
	}
}

func (app *RuntimeApp) getMenuItemLaunch() *MenuItem {
	var appName string
	if app.appInfo != nil {
		appName = strings.Title(app.appInfo.GetDisplayName())
	} else {
		appName = strings.Title(app.Id)
	}
	return NewMenuItem(appName, app.launchApp, true)
}

func (app *RuntimeApp) getMenuItemCloseAll() *MenuItem {
	return NewMenuItem(Tr("_Close All"), func(timestamp uint32) {
		logger.Debug("Close All")
		for win, _ := range app.windowInfoTable {
			ewmh.CloseWindow(XU, win)
		}
	}, true)
}

func (app *RuntimeApp) getMenuItemDock() *MenuItem {
	return NewMenuItem(Tr("_Dock"), func(timestamp uint32) {
		logger.Debug("Dock app:", app.Id)
		var title, icon, exec string
		appInfo := app.appInfo
		// TODO
		if appInfo.DesktopAppInfo != nil {
			title = appInfo.GetDisplayName()
			icon = appInfo.GetIcon()
			// TODO get cmdline
			exec = appInfo.DesktopAppInfo.GetString(glib.KeyFileDesktopKeyExec)
		} else {
			title = app.Id
			// icon
			icon = "application-default-icon"
			for _, v := range app.windowInfoTable {
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
			// exec
			exec = filepath.Join(scratchDir, app.Id+".sh")
			scriptContent := "#!/bin/sh\n" + app.exec
			ioutil.WriteFile(exec, []byte(scriptContent), 0744)
		}
		logger.Debugf("title: %q, icon: %q, exec: %q", title, icon, exec)
		// TODO:
		dockManager.entryManager.dockedAppManager.Dock(app.Id, title, icon, exec)
	}, true)
}

func (app *RuntimeApp) getMenuItemUndock() *MenuItem {
	return NewMenuItem(Tr("_Undock"), func(uint32) {
		// TODO:
		dockManager.entryManager.dockedAppManager.Undock(app.Id)
	}, true)
}

func (app *RuntimeApp) setChangedCB(cb func()) {
	app.changedCB = cb
}

func (app *RuntimeApp) notifyChanged() {
	if app.changedCB != nil {
		logger.Debug("call notifyChanged")
		app.changedCB()
	}
}

func (app *RuntimeApp) HandleMenuItem(id string, timestamp uint32) {
	if app.coreMenu != nil {
		app.coreMenu.HandleAction(id, timestamp)
	}
}

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

// TODO: 废弃
// func find_app_id_by_xid(xid xproto.Window, displayMode DisplayModeType) string {
// 	var appId string
// 	logger.Debugf("find_app_id_by_xid 0x%x", xid)
// 	if displayMode == DisplayModeModernMode {
// 		if id, err := xprop.PropValStr(xprop.GetProperty(XU, xid, "_DDE_DOCK_APP_ID")); err == nil {
// 			appId = getAppIDFromDesktopID(normalizeAppID(id))
// 			if appId != "" {
// 				logger.Info("get app id from _DDE_DOCK_APP_ID", appId)
// 				return appId
// 			}
// 		}
// 	}
//
// 	gtkAppId, err := xprop.PropValStr(xprop.GetProperty(XU, xid, "_GTK_APPLICATION_ID"))
// 	if err != nil {
// 		logger.Debug("get AppId from _GTK_APPLICATION_ID failed:", err)
// 	} else {
// 		appId = gtkAppId
// 		appId = getAppIDFromDesktopID(normalizeAppID(appId))
// 		if appId != "" {
// 			return appId
// 		}
// 	}
//
// 	appId = getAppIDFromXid(xid)
// 	if appId != "" {
// 		logger.Debug("get app id from bamf", appId)
// 		return normalizeAppID(appId)
// 	}
//
// 	wmClass, _ := icccm.WmClassGet(XU, xid)
// 	var wmInstance, wmClassName string
// 	if wmClass != nil {
// 		if utf8.ValidString(wmClass.Instance) {
// 			wmInstance = wmClass.Instance
// 		}
//
// 		// it is possible that getting invalid string which might be xgb implementation's bug.
// 		// for instance: xdemineur's WMClass
// 		if utf8.ValidString(wmClass.Class) {
// 			wmClassName = wmClass.Class
// 		}
// 		logger.Debug("WMClass", wmClassName, ", WMInstance", wmInstance)
// 	}
//
// 	name := getWmName(XU, xid)
// 	logger.Debugf("wmName is %q", name)
//
// 	pid, err := ewmh.WmPidGet(XU, xid)
// 	if err != nil {
// 		logger.Debug("get pid failed for:", xid)
// 		if name != "" {
// 			pid = lookthroughProc(name)
// 			logger.Debugf("lookthroughProc(%q) get pid %v", name, pid)
// 		} else {
// 			newAppId := getAppIDFromDesktopID(normalizeAppID(wmClassName))
// 			if newAppId != "" {
// 				logger.Debugf("get Pid failed, guess app id `%s` by wm class name `%s`",
// 					newAppId, wmClassName)
// 				return newAppId
// 			}
// 		}
// 	}
//
// 	iconName, _ := ewmh.WmIconNameGet(XU, xid)
// 	if pid == 0 {
// 		appId = normalizeAppID(wmClassName)
// 		logger.Debugf("get window name failed, using wm class name as app id %q", appId)
// 		return appId
// 	}
//
// 	logger.Debugf("call c.find_app_id(pid=%v, wmName=%s, wmInstance=%s, wmClassName=%s, iconName=%s)",
// 		pid, name, wmInstance, wmClassName, iconName)
// 	appId = find_app_id(pid, name, wmInstance, wmClassName, iconName)
// 	newAppId := getAppIDFromDesktopID(normalizeAppID(appId))
// 	if newAppId != "" {
// 		appId = newAppId
// 	}
//
// 	appId = specialCaseWorkaround(xid, appId)
//
// 	logger.Debugf("get appid %q", appId)
// 	return appId
// }

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

// TODO: move to winInfo.updateIcon
func (app *RuntimeApp) updateIcon(winInfo *WindowInfo) {
	logger.Debug("update icon for", app.Id)
	if app.appInfo != nil {
		icon := app.appInfo.GetIcon()
		if icon != "" {
			winInfo.Icon = icon
			return
		}
		logger.Warning("get icon from app failed")
	}

	logger.Info(app.Id, "using icon from X")
	winInfo.Icon = getIconFromWindow(XU, winInfo.window)
	if winInfo.Icon == "" {
		winInfo.Icon = "application-default-icon"
	}
}

// func (app *RuntimeApp) updateAppid(xid xproto.Window) {
// 	newAppId := find_app_id_by_xid(
// 		xid,
// 		DisplayModeType(dockManager.setting.GetDisplayMode()),
// 	)
// 	if app.Id != newAppId {
// 		app.detachXid(xid)
// 		if newApp := dockManager.entryManager.createRuntimeApp(xid); newApp != nil {
// 			newApp.attachXid(xid)
// 		}
// 		logger.Debug("APP:", app.Id, "Changed to..", newAppId)
// 		//TODO: Destroy
// 	}
// }

func iconifyWindow(xid xproto.Window) {
	ewmh.ClientEvent(XU, xid, "WM_CHANGE_STATE", icccm.StateIconic)
}

func (app *RuntimeApp) Activate(x, y int32, timestamp uint32) error {
	// x,y, timestamp useless
	xid := app.CurrentInfo.window
	state, err := ewmh.WmStateGet(XU, xid)
	if err != nil {
		logger.Warning("Get ewmh WmState failed xid:", xid)
		return err
	}
	if contains(state, "_NET_WM_STATE_FOCUSED") {
		s, err := icccm.WmStateGet(XU, xid)
		if err != nil {
			logger.Warning("Get icccm WmState failed xid:", xid)
			return err
		}
		switch s.State {
		case icccm.StateIconic:
			s.State = icccm.StateNormal
			icccm.WmStateSet(XU, xid, s)
		case icccm.StateNormal:
			if len(app.windowInfoTable) == 1 {
				iconifyWindow(xid)
			} else {
				if dockManager.activeWindow == xid {
					nextWin := app.findNextLeader()
					ewmh.ActiveWindowReq(XU, nextWin)
				}
			}
		}
	} else {
		activateWindow(xid)
	}
	return nil
}

func (app *RuntimeApp) setLeader(leader xproto.Window) {
	if info, ok := app.windowInfoTable[leader]; ok {
		if app.CurrentInfo != info {
			app.CurrentInfo = info
			app.notifyChanged()
		}
	}
}

func (app *RuntimeApp) findNextLeader() xproto.Window {
	winSlice := make(windowSlice, 0, len(app.windowInfoTable))
	for win, _ := range app.windowInfoTable {
		winSlice = append(winSlice, win)
	}
	sort.Sort(winSlice)
	currentWin := app.CurrentInfo.window
	logger.Debug("sorted window slice:", winSlice)
	logger.Debug("current window:", currentWin)
	currentIndex := -1
	for i, win := range winSlice {
		if win == currentWin {
			currentIndex = i
		}
	}
	if currentIndex == -1 {
		logger.Warning("findNextLeader unexpect, return 0")
		return 0
	}
	// if current window is max, return min: winSlice[0]
	// else return winSlice[currentIndex+1]
	nextIndex := 0
	if currentIndex < len(winSlice)-1 {
		nextIndex = currentIndex + 1
	}
	logger.Debug("next window:", winSlice[nextIndex])
	return winSlice[nextIndex]
}

func (app *RuntimeApp) attachWindow(winInfo *WindowInfo) {
	win := winInfo.window
	logger.Debugf("attach win %v to app %v", win, app.Id)
	if _, ok := app.windowInfoTable[win]; ok {
		logger.Debugf("win %v is already attach to app %v", win, app.Id)
		return
	}

	app.windowInfoTable[win] = winInfo
	winInfo.app = app
	app.updateIcon(winInfo)
	winInfo.Title = winInfo.getTitle()
	app.notifyChanged()
}

func (app *RuntimeApp) detachWindow(winInfo *WindowInfo) {
	win := winInfo.window
	winInfo.app = nil
	if info, ok := app.windowInfoTable[win]; ok {
		delete(app.windowInfoTable, win)
		if len(app.windowInfoTable) == 0 {
			app.setChangedCB(nil)
			return
		}

		if info == app.CurrentInfo {
			// switch to next
			for _, nextInfo := range app.windowInfoTable {
				app.CurrentInfo = nextInfo
			}
		}
		app.notifyChanged()
	}
}

func (app *RuntimeApp) HandleDragDrop(path string, timestamp uint32) {
	logger.Debugf("handle drag drop path: %q", path)
	ai := app.appInfo
	appLaunchContext := gio.GetGdkAppLaunchContext().SetTimestamp(timestamp)
	if ai.DesktopAppInfo != nil {
		paths := []string{path}
		_, err := ai.LaunchUris(paths, appLaunchContext)
		if err != nil {
			logger.Warningf("LaunchUris failed path: %q", path)
		}
	} else {
		logger.Warningf("no support!")
	}
}
