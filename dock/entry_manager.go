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
	"errors"
	"fmt"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xprop"
	"os"
	"path/filepath"
	"pkg.deepin.io/lib/dbus"
	"sort"
)

// EntryManager为驻留程序以及打开程序的管理器。
type EntryManager struct {
	runtimeApps map[string]*RuntimeApp
	normalApps  map[string]*NormalApp
	appEntries  map[string]*AppEntry

	// 保持 entry 的顺序
	entryIDs []string

	dockedAppManager *DockedAppManager
	clientList       windowSlice
	appIdFilterGroup *AppIdFilterGroup
}

func NewEntryManager() *EntryManager {
	m := &EntryManager{
		runtimeApps:      make(map[string]*RuntimeApp),
		normalApps:       make(map[string]*NormalApp),
		appEntries:       make(map[string]*AppEntry),
		dockedAppManager: NewDockedAppManager(),
		appIdFilterGroup: NewAppIdFilterGroup(),
	}
	return m
}

// Reorder 重排序dock上的app项目
// 参数entryIDs为dock上app项目的新顺序id列表，要求与当前app项目是同一个集合，只是顺序不同。
func (m *EntryManager) Reorder(entryIDs []string) error {
	logger.Debugf("Reorder entryIDs %#v", entryIDs)
	if len(entryIDs) != len(m.entryIDs) {
		logger.Warning("Reorder: len(entryIDs) != len(m.entryIDs)")
		return errors.New("length of incomming entryIDs not equal length of m.entryIDs")
	}
	var orderedEntryIDs []string
	for _, id := range entryIDs {
		_, ok := m.appEntries[id]
		if ok {
			orderedEntryIDs = append(orderedEntryIDs, id)
		} else {
			logger.Warningf("Reorder: invaild entry id %q", id)
			return fmt.Errorf("Invaild entry id %q", id)
		}
	}
	m.entryIDs = orderedEntryIDs
	m.dockedAppManager.reorderThenSave(m.entryIDs)
	return nil
}

func (m *EntryManager) GetEntryIDs() []string {
	return m.entryIDs
}

func (m *EntryManager) getRuntimeAppByWindow(win xproto.Window) *RuntimeApp {
	for _, app := range m.runtimeApps {
		_, ok := app.windowInfoTable[win]
		if ok {
			return app
		}
	}
	return nil
}

func (m *EntryManager) updateActiveWindow(win xproto.Window) {
	app := m.getRuntimeAppByWindow(win)
	if app != nil {
		app.setLeader(win)
	}
}

func (m *EntryManager) attachOrDetachRuntimeAppWindow(winInfo *WindowInfo) {
	win := winInfo.window
	canShowOnDock := winInfo.canShowOnDock()
	logger.Debugf("win %v canShowOnDock? %v", win, canShowOnDock)
	app := winInfo.app
	if app != nil {
		if !canShowOnDock {
			m.detachRuntimeAppWindow(winInfo)
		}
	} else {
		// app is nil
		if canShowOnDock && m.clientList.Contains(win) {
			m.attachRuntimeAppWindow(winInfo)
		}
	}
}

func (m *EntryManager) initRuntimeApps() {
	clientList, err := ewmh.ClientListGet(XU)
	if err != nil {
		logger.Warning("Get client list failed:", err)
		return
	}
	winSlice := windowSlice(clientList)
	sort.Sort(winSlice)
	for _, win := range winSlice {
		winInfo := NewWindowInfo(win)
		m.listenWindowXEvent(winInfo)
		m.attachOrDetachRuntimeAppWindow(winInfo)
	}
	m.clientList = winSlice
}

func (m *EntryManager) initDockedApps() {
	for _, id := range m.dockedAppManager.DockedAppList() {
		id = normalizeAppID(id)
		logger.Debug("load docked app", id)
		m.createNormalApp(id)
	}
}

func (m *EntryManager) addAppEntry(id string, e *AppEntry) {
	m.appEntries[id] = e
	m.entryIDs = append(m.entryIDs, id)

	err := dbus.InstallOnSession(e)
	if err != nil {
		logger.Warning("Install AppEntry to dbus failed:", err)
	}
}

func (m *EntryManager) mustGetEntry(nApp *NormalApp, rApp *RuntimeApp) *AppEntry {
	if rApp != nil {
		if e, ok := m.appEntries[rApp.Id]; ok {
			return e
		} else {
			e := NewAppEntryWithRuntimeApp(rApp)
			m.addAppEntry(rApp.Id, e)
			return e
		}
	} else if nApp != nil {
		if e, ok := m.appEntries[nApp.Id]; ok {
			return e
		} else {
			e := NewAppEntryWithNormalApp(nApp)
			m.addAppEntry(nApp.Id, e)
			return e
		}
	}
	panic("mustGetEntry: at least give one app")
}

func (m *EntryManager) destroyEntry(appId string) {
	if e, ok := m.appEntries[appId]; ok {
		e.detachNormalApp()
		e.detachRuntimeApp()
		dbus.ReleaseName(e)
		dbus.UnInstallObject(e)
		logger.Info("destroyEntry:", appId)
	}
	delete(m.appEntries, appId)
	m.entryIDs = strSliceRemove(m.entryIDs, appId)
}

func strSliceRemove(slice []string, str string) []string {
	var index int = -1
	for i, v := range slice {
		if v == str {
			index = i
		}
	}
	if index != -1 {
		return append(slice[:index], slice[index+1:]...)
	}
	return slice
}

func (m *EntryManager) updateEntry(appId string, nApp *NormalApp, rApp *RuntimeApp) {
	switch {
	case nApp == nil && rApp == nil:
		m.destroyEntry(appId)
	case nApp == nil && rApp != nil:
		e := m.mustGetEntry(nApp, rApp)
		e.attachRuntimeApp(rApp)
		e.detachNormalApp()
	case nApp != nil && rApp != nil:
		e := m.mustGetEntry(nApp, rApp)
		e.attachNormalApp(nApp)
		e.attachRuntimeApp(rApp)
	case nApp != nil && rApp == nil:
		e := m.mustGetEntry(nApp, rApp)
		e.attachNormalApp(nApp)
		e.detachRuntimeApp()
	}
}

func (m *EntryManager) getAppInfoFromWindow(winInfo *WindowInfo) *AppInfo {
	win := winInfo.window
	var ai *AppInfo

	// _GTK_APPLICATION_ID
	gtkAppId, err := xprop.PropValStr(xprop.GetProperty(XU, win, "_GTK_APPLICATION_ID"))
	if err != nil {
		logger.Debug("get AppId from _GTK_APPLICATION_ID failed:", err)
	} else {
		ai = NewAppInfo(gtkAppId + ".desktop")
		if ai != nil {
			return ai
		}
		logger.Debugf("NewAppInfo failed gtk app id: %q", gtkAppId)
	}

	// bamf
	desktop := getDesktopFromWindowByBamf(win)
	if desktop == "" {
		logger.Debug("get desktop from bamf failed")
	} else {
		logger.Debugf("bamf desktop: %q", desktop)
		ai = NewAppInfoFromFile(desktop)
		if ai != nil {
			return ai
		}
		logger.Debugf("NewAppInfoFromFile faield, desktop: %q", desktop)
	}

	// 通常不由 desktop 文件启动的应用 bamf 识别容易失败
	winGuessAppId := winInfo.guessAppId(m.appIdFilterGroup)
	if winGuessAppId == "" {
		logger.Debug("guess app id by window info failed")
	} else {
		logger.Debugf("win guess app id: %q", winGuessAppId)
		ai = NewAppInfo(winGuessAppId + ".desktop")
		if ai != nil {
			return ai
		}
		logger.Debugf("NewAppInfo failed win guess app id: %q", winGuessAppId)
	}

	// fail
	winAppInfo := NewAppInfoFromWindow(winInfo)
	return winAppInfo
}

// 给 window 一个 runtimeApp
// 根据 window id 找到 appId， 如果 runtimeApp 已经存, 则 app.attachWindow
// 如果不存在，则 NewRuntimeApp 创建新的 RuntimeApp
func (m *EntryManager) attachRuntimeAppWindow(winInfo *WindowInfo) *RuntimeApp {
	win := winInfo.window
	appInfo := m.getAppInfoFromWindow(winInfo)
	if appInfo == nil {
		logger.Warning("getAppInfoFromWindow failed, win:", win)
		return nil
	}
	appId := appInfo.GetId()

	if v, ok := m.runtimeApps[appId]; ok {
		v.attachWindow(winInfo)
		return v
	}

	isAppDocked := m.dockedAppManager.IsDocked(appId)
	rApp := NewRuntimeApp(winInfo, appInfo, isAppDocked)
	if rApp == nil {
		logger.Warningf("NewRuntimeApp failed win %v app id %v", win, appId)
		return nil
	}

	m.runtimeApps[appId] = rApp
	m.updateEntry(appId, m.mustGetEntry(nil, rApp).nApp, rApp)
	return rApp
}

// 取消绑定 winInfo.app 与 winInfo , 如果 rApp 窗口数量为 0 者销毁 rApp
func (m *EntryManager) detachRuntimeAppWindow(winInfo *WindowInfo) {
	app := winInfo.app
	if app == nil {
		return
	}
	app.detachWindow(winInfo)
	if len(app.windowInfoTable) == 0 {
		m.destroyRuntimeApp(app)
	}
}

func (m *EntryManager) destroyRuntimeApp(rApp *RuntimeApp) {
	logger.Debug("Destory runtime app", rApp.Id)
	delete(m.runtimeApps, rApp.Id)
	m.updateEntry(rApp.Id, m.mustGetEntry(nil, rApp).nApp, nil)
}

func (m *EntryManager) createNormalApp(id string) {
	logger.Info("createNormalApp for", id)
	if _, ok := m.normalApps[id]; ok {
		logger.Debug("normal app for", id, "is exist")
		return
	}

	desktopId := id + ".desktop"
	nApp := NewNormalApp(desktopId)
	if nApp == nil {
		logger.Info("get desktop file failed, create", id, "from scratch file")
		newId := filepath.Join(
			os.Getenv("HOME"),
			".config/dock/scratch",
			desktopId,
		)
		nApp = NewNormalApp(newId)
		if nApp == nil {
			logger.Warning("create normal app failed:", id)
			m.dockedAppManager.Undock(id)
			return
		}
	}

	m.normalApps[id] = nApp
	m.updateEntry(id, nApp, m.mustGetEntry(nApp, nil).rApp)
}

func (m *EntryManager) destroyNormalApp(id string) {
	if nApp, ok := m.normalApps[id]; ok {
		logger.Debugf("destroyNormalApp id: %q", id)
		delete(m.normalApps, nApp.Id)
		m.updateEntry(nApp.Id, nil, m.mustGetEntry(nApp, nil).rApp)
	} else {
		logger.Debugf("no need destroyNormalApp id: %q", id)
	}
}
