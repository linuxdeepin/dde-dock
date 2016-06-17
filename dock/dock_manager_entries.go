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
	"fmt"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/ewmh"
	"pkg.deepin.io/lib/dbus"
	"sort"
	"strconv"
)

func (m *DockManager) allocEntryId() string {
	num := m.entryCount
	m.entryCount++
	return fmt.Sprintf("e%dT%x", num, getCurrentTimestamp())
}

func (m *DockManager) getAppEntryByWindow(win xproto.Window) *AppEntry {
	for _, entry := range m.Entries {
		_, ok := entry.windows[win]
		if ok {
			return entry
		}
	}
	return nil
}

func (m *DockManager) getAppEntryByAppId(id string) *AppEntry {
	for _, entry := range m.Entries {
		if entry.appInfo != nil && id == entry.appInfo.GetId() {
			return entry
		}
	}
	return nil
}

func (m *DockManager) getAppEntryByEntryId(id string) *AppEntry {
	for _, entry := range m.Entries {
		if entry.Id == id {
			return entry
		}
	}
	return nil
}

func (m *DockManager) getAppEntryByInnerId(id string) *AppEntry {
	for _, entry := range m.Entries {
		if entry.innerId == id {
			return entry
		}
	}
	return nil
}

func (m *DockManager) attachOrDetachWindow(winInfo *WindowInfo) {
	win := winInfo.window
	canShowOnDock := winInfo.canShowOnDock()
	logger.Debugf("win %v canShowOnDock? %v", win, canShowOnDock)
	entry := winInfo.entry
	if entry != nil {
		if !canShowOnDock {
			m.detachWindow(winInfo)
		}
	} else {
		if canShowOnDock && m.clientList.Contains(win) {
			m.attachWindow(winInfo)
		}
	}
}

func (m *DockManager) initClientList() {
	clientList, err := ewmh.ClientListGet(XU)
	if err != nil {
		logger.Warning("Get client list failed:", err)
		return
	}
	winSlice := windowSlice(clientList)
	sort.Sort(winSlice)
	m.clientList = winSlice
	for _, win := range winSlice {
		winInfo := NewWindowInfo(win)
		m.listenWindowXEvent(winInfo)
		m.attachOrDetachWindow(winInfo)
	}
}

func (m *DockManager) initDockedApps() {
	for _, app := range m.DockedApps.Get() {
		m.appendDockedApp(app)
	}
}

func (m *DockManager) installAppEntry(e *AppEntry) {
	// install on session D-Bus
	err := dbus.InstallOnSession(e)
	if err != nil {
		logger.Warning("Install AppEntry to dbus failed:", err)
		return
	}

	entryObjPath := dbus.ObjectPath(entryDBusObjPathPrefix + e.Id)
	logger.Debugf("insertAndInstallAppEntry %v", entryObjPath)
	index := -1
	for i, entry := range m.Entries {
		if e.Id == entry.Id {
			index = i
			break
		}
	}
	if index >= 0 {
		dbus.Emit(m, "EntryAdded", entryObjPath, int32(index))
	}
}

func (m *DockManager) insertAppEntry(e *AppEntry, index int) {
	var insertIndex int
	m.Entries, insertIndex = entrySliceInsert(m.Entries, e, index)
	logger.Debugf("insertAppEntry entry: %v insert at: %v", e.Id, insertIndex)
}

func (m *DockManager) addAppEntry(entryInnerId string, appInfo *AppInfo, index int) (*AppEntry, bool) {
	logger.Debug("addAppEntry innerId:", entryInnerId)

	var entry *AppEntry
	isNewAdded := false
	if e := m.getAppEntryByInnerId(entryInnerId); e != nil {
		logger.Debug("entry existed")
		entry = e
		if appInfo != nil {
			appInfo.Destroy()
		}
	} else {
		// cache desktop hash => desktop file path
		if appInfo != nil {
			m.desktopHashFileMapCacheManager.SetKeyValue(appInfo.innerId, appInfo.GetFilePath())
			m.desktopHashFileMapCacheManager.AutoSave()
		}
		logger.Debug("entry not existed, newAppEntry")
		entry = newAppEntry(m, entryInnerId, appInfo)
		m.insertAppEntry(entry, index)
		isNewAdded = true
	}
	return entry, isNewAdded
}

func (m *DockManager) appendDockedApp(appId string) {
	logger.Infof("appendDockedApp %q", appId)
	appInfo := NewAppInfo(appId)
	if appInfo == nil {
		logger.Warning("appendDockedApp failed: appInfo is nil")
		return
	}
	entry, isNewAdded := m.addAppEntry(appInfo.innerId, appInfo, -1)
	entry.setIsDocked(true)
	entry.updateMenu()
	if isNewAdded {
		entry.updateName()
		entry.updateIcon()
		m.installAppEntry(entry)
	}
}

func (m *DockManager) removeAppEntry(e *AppEntry) {
	for _, entry := range m.Entries {
		if entry == e {
			dbus.UnInstallObject(e)

			entryId := entry.Id
			logger.Info("removeAppEntry id:", entryId)
			m.Entries = entrySliceRemove(m.Entries, e)
			e.destroy()
			dbus.Emit(m, "EntryRemoved", entryId)
			return
		}
	}
	logger.Warning("removeAppEntry failed, entry not found")
}

func entrySliceInsert(slice []*AppEntry, entry *AppEntry, index int) ([]*AppEntry, int) {
	logger.Debug("entrySliceInsert index:", index)
	sliceLen := len(slice)
	if index < 0 || index >= len(slice) {
		logger.Debug("entrySliceInsert: append")
		return append(slice, entry), sliceLen
	}
	// insert
	return append(slice[:index],
		append([]*AppEntry{entry}, slice[index:]...)...), index
}

func entrySliceRemove(slice []*AppEntry, entry *AppEntry) []*AppEntry {
	var index int = -1
	for i, v := range slice {
		if v.Id == entry.Id {
			index = i
		}
	}
	if index != -1 {
		return append(slice[:index], slice[index+1:]...)
	}
	return slice
}

func (m *DockManager) identifyWindow(winInfo *WindowInfo) (string, *AppInfo) {
	logger.Debugf("identifyWindow: window id: %v, window hash %v", winInfo.window, winInfo.innerId)
	desktopHash := m.desktopWindowsMapCacheManager.GetKeyByValue(winInfo.innerId)
	logger.Debug("identifyWindow: get desktop hash:", desktopHash)
	var appInfo *AppInfo
	if desktopHash != "" {
		appInfo = m.desktopHashFileMapCacheManager.GetAppInfo(desktopHash)
		logger.Debug("identifyWindow: get AppInfo by desktop hash:", appInfo)
	}

	if appInfo == nil {
		// cache fail
		if desktopHash != "" {
			logger.Warning("winHash->DesktopHash success, but DesktopHash->appInfo fail")
			m.desktopHashFileMapCacheManager.DeleteKey(desktopHash)
			m.desktopWindowsMapCacheManager.DeleteKeyValue(desktopHash, winInfo.innerId)
		}

		var canCache bool
		appInfo, canCache = m.getAppInfoFromWindow(winInfo)
		logger.Debug("identifyWindow: getAppInfoFromWindow:", appInfo)
		if appInfo != nil && canCache {
			m.desktopWindowsMapCacheManager.AddKeyValue(appInfo.innerId, winInfo.innerId)
			m.desktopHashFileMapCacheManager.SetKeyValue(appInfo.innerId, appInfo.GetFilePath())
		}
	}

	var entryInnerId string
	if appInfo != nil {
		entryInnerId = appInfo.innerId
		logger.Debug("Set entryInnerId to desktop hash")
	} else {
		entryInnerId = winInfo.innerId
		logger.Debug("Set entryInnerId to window hash")
	}

	m.desktopWindowsMapCacheManager.AutoSave()
	m.desktopHashFileMapCacheManager.AutoSave()
	return entryInnerId, appInfo
}

func (m *DockManager) attachWindow(winInfo *WindowInfo) {
	entryInnerId, appInfo := m.identifyWindow(winInfo)
	entry, isNewAdded := m.addAppEntry(entryInnerId, appInfo, -1)
	entry.windowMutex.Lock()
	defer entry.windowMutex.Unlock()

	entry.attachWindow(winInfo)
	entry.updateMenu()
	if isNewAdded {
		entry.initExec(winInfo)
		entry.updateName()
		entry.updateIcon()
		m.installAppEntry(entry)
	}
}

func (m *DockManager) detachWindow(winInfo *WindowInfo) {
	entry := winInfo.entry
	if entry == nil {
		return
	}
	entry.windowMutex.Lock()
	defer entry.windowMutex.Unlock()

	detached := entry.detachWindow(winInfo)
	if !detached {
		return
	}
	if !entry.hasWindow() && !entry.IsDocked {
		m.removeAppEntry(entry)
		return
	}
	entry.updateWindowTitles()
	entry.updateIcon()
	entry.updateMenu()
	entry.updateIsActive()
}

func (m *DockManager) getAppInfoFromWindow(winInfo *WindowInfo) (*AppInfo, bool) {
	win := winInfo.window
	var ai *AppInfo

	gtkAppId := winInfo.gtkAppId
	logger.Debug("Try gtkAppId", gtkAppId)
	if gtkAppId != "" {
		ai = NewAppInfo(gtkAppId)
		if ai != nil {
			logger.Debugf("Get AppInfo success gtk app id: %q", gtkAppId)
			return ai, true
		}
	}

	// env GIO_LAUNCHED_DESKTOP_FILE
	var launchedDesktopFile string
	logger.Debug("Try process env")
	if winInfo.process != nil {
		envVars, err := getProcessEnvVars(winInfo.process.pid)
		if err == nil {
			launchedDesktopFile = envVars["GIO_LAUNCHED_DESKTOP_FILE"]
			pidStr := envVars["GIO_LAUNCHED_DESKTOP_FILE_PID"]
			launchedDesktopFilePid, _ := strconv.ParseUint(pidStr, 10, 32)
			logger.Debugf("launchedDesktopFile: %q, pid: %v", launchedDesktopFile, launchedDesktopFilePid)
			if winInfo.process.pid != 0 &&
				uint(launchedDesktopFilePid) == winInfo.process.pid {
				ai = NewAppInfoFromFile(launchedDesktopFile)
				if ai != nil {
					logger.Debugf("Get AppInfo success pid equal launchedDesktopFile: %q", launchedDesktopFile)
					return ai, true
				}
			}
		}
	}

	// bamf
	desktop := getDesktopFromWindowByBamf(win)
	logger.Debug("Try bamf")
	if desktop != "" {
		ai = NewAppInfoFromFile(desktop)
		if ai != nil {
			logger.Debugf("Get AppInfo success bamf desktop: %q", desktop)
			return ai, true
		}
	}

	// 通常不由 desktop 文件启动的应用 bamf 识别容易失败
	winGuessAppId := winInfo.guessAppId(m.appIdFilterGroup)
	logger.Debug("Try filter group", winGuessAppId)
	if winGuessAppId != "" {
		ai = NewAppInfo(winGuessAppId)
		if ai != nil {
			logger.Debugf("Get AppInfo success winGuessAppId: %q", winGuessAppId)
			return ai, true
		}
	}

	logger.Debug("Try env var launchedDesktopFile")
	if launchedDesktopFile != "" {
		ai = NewAppInfoFromFile(launchedDesktopFile)
		if ai != nil {
			logger.Debugf("Get AppInfo success launchedDesktopFile %q", launchedDesktopFile)
			return ai, false
		}
	}

	logger.Debug("Get AppInfo failed")
	return nil, false
}
