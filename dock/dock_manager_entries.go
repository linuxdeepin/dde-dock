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
	"github.com/BurntSushi/xgbutil/ewmh"
	"pkg.deepin.io/lib/dbus"
	"sort"
)

func (m *DockManager) allocEntryId() string {
	num := m.entryCount
	m.entryCount++
	return fmt.Sprintf("e%dT%x", num, getCurrentTimestamp())
}

func (m *DockManager) markAppLaunched(appInfo *AppInfo) {
	if appInfo == nil {
		return
	}
	id := appInfo.GetId()
	if id == "" {
		logger.Warning("markAppLaunched failed, appInfo %v no id", appInfo)
		return
	}
	go func() {
		if m.launcher == nil {
			return
		}
		logger.Infof("mark app %q launched", id)
		m.launcher.MarkLaunched(id)
		recordFrequency(id)
	}()
}

func (m *DockManager) attachOrDetachWindow(winInfo *WindowInfo) {
	win := winInfo.window
	showOnDock := m.isWindowRegistered(win) && m.clientList.Contains(win) &&
		winInfo.canShowOnDock()
	logger.Debugf("win %v showOnDock? %v", win, showOnDock)
	entry := winInfo.entry
	if entry != nil {
		if !showOnDock {
			m.detachWindow(winInfo)
		}
	} else {

		if winInfo.entryInnerId == "" {
			winInfo.entryInnerId, winInfo.appInfo = m.identifyWindow(winInfo)
			m.markAppLaunched(winInfo.appInfo)
		} else {
			logger.Debugf("win %v identified", win)
		}

		if showOnDock {
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
		m.registerWindow(win)
	}
}

func (m *DockManager) initDockedApps() {
	dockedApps := uniqStrSlice(m.DockedApps.Get())
	for _, app := range dockedApps {
		m.appendDockedApp(app)
	}
	m.saveDockedApps()
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
	index := m.Entries.IndexOf(e)
	if index >= 0 {
		dbus.Emit(m, "EntryAdded", entryObjPath, int32(index))
	}
}

func (m *DockManager) addAppEntry(entryInnerId string, appInfo *AppInfo, index int) (*AppEntry, bool) {
	logger.Debug("addAppEntry innerId:", entryInnerId)

	var entry *AppEntry
	isNewAdded := false
	if e := m.Entries.GetFirstByInnerId(entryInnerId); e != nil {
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
		m.Entries = m.Entries.Insert(entry, index)
		logger.Debugf("insert entry %v at %v", entry.Id, index)
		isNewAdded = true
	}
	return entry, isNewAdded
}

func (m *DockManager) appendDockedApp(app string) {
	logger.Infof("appendDockedApp %q", app)
	appInfo := NewDockedAppInfo(app)
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
			m.Entries = m.Entries.Remove(e)
			e.destroy()
			dbus.Emit(m, "EntryRemoved", entryId)
			return
		}
	}
	logger.Warning("removeAppEntry failed, entry not found")
}

func (m *DockManager) attachWindow(winInfo *WindowInfo) {
	var appInfoCopy *AppInfo
	if winInfo.appInfo != nil {
		appInfoCopy = NewAppInfoFromFile(winInfo.appInfo.GetFilePath())
	}
	entry, isNewAdded := m.addAppEntry(winInfo.entryInnerId, appInfoCopy, -1)
	entry.windowMutex.Lock()
	defer entry.windowMutex.Unlock()

	entry.attachWindow(winInfo)
	entry.updateMenu()
	if isNewAdded {
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
	winInfo.entry = nil
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
