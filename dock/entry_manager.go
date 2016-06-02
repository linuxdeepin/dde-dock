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
	"gir/gio-2.0"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/ewmh"
	"io/ioutil"
	"os"
	"path/filepath"
	"pkg.deepin.io/lib/dbus"
	"sort"
	"strconv"
	"strings"
	"time"
)

// EntryManager为驻留程序以及打开程序的管理器。
type EntryManager struct {
	activeWindow xproto.Window

	dockedAppManager *DockedAppManager
	clientList       windowSlice
	appIdFilterGroup *AppIdFilterGroup
	windowDesktopMap *windowDesktopMap
	desktopIdFileMap *desktopIdFileMap

	Entries []*AppEntry
	// Added在程序需要在前端显示时被触发。
	Added func(dbus.ObjectPath)
	// Removed会在程序不再需要在dock前端显示时触发。
	Removed func(string)
	// 废弃：TrayInited在trayicon相关内容初始化完成后触发。
	TrayInited func()
	entryCount uint
}

func (m *EntryManager) allocEntryId() string {
	num := m.entryCount
	m.entryCount++
	timeNow := time.Now()
	timeNowUnixSeconds := timeNow.Unix()
	return fmt.Sprintf("e%dT%x", num, timeNowUnixSeconds)
}

func (m *EntryManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.daemon.Dock",
		ObjectPath: "/dde/dock/EntryManager",
		Interface:  "dde.dock.EntryManager",
	}
}

func NewEntryManager() (*EntryManager, error) {
	m := &EntryManager{
		appIdFilterGroup: NewAppIdFilterGroup(),
	}
	var err error
	m.windowDesktopMap, err = newWindowDesktopMapFromFile(filepath.Join(cacheDir, "windowDesktopMap.gob"))
	if err != nil {
		return nil, err
	}
	m.desktopIdFileMap, err = newDesktopIdFileMapFromFile(filepath.Join(cacheDir, "desktopIdFileMap.gob"))
	if err != nil {
		return nil, err
	}

	m.dockedAppManager = NewDockedAppManager(m)
	return m, nil
}

// Reorder 重排序dock上的app项目
// 参数entryIDs为dock上app项目的新顺序id列表，要求与当前app项目是同一个集合，只是顺序不同。
func (m *EntryManager) Reorder(entryIDs []string) error {
	logger.Debugf("Reorder entryIDs %#v", entryIDs)
	if len(entryIDs) != len(m.Entries) {
		logger.Warning("Reorder: len(entryIDs) != len(m.Entries)")
		return errors.New("length of incomming entryIDs not equal length of m.Entries")
	}
	var orderedEntries []*AppEntry
	for _, id := range entryIDs {
		// TODO: 优化
		entry := m.getAppEntryByEntryId(id)
		if entry != nil {
			orderedEntries = append(orderedEntries, entry)
		} else {
			logger.Warningf("Reorder: invaild entry id %q", id)
			return fmt.Errorf("Invaild entry id %q", id)
		}
	}
	m.Entries = orderedEntries
	m.dockedAppManager.saveDockedAppList()
	return nil
}

func (m *EntryManager) GetEntryIDs() []string {
	list := make([]string, 0, len(m.Entries))
	for _, entry := range m.Entries {
		var appId string
		if entry.appInfo != nil {
			appId = entry.appInfo.GetId()
		} else {
			appId = entry.innerId
		}
		list = append(list, appId)
	}
	return list
}

func (m *EntryManager) getAppEntryByWindow(win xproto.Window) *AppEntry {
	for _, entry := range m.Entries {
		_, ok := entry.windows[win]
		if ok {
			return entry
		}
	}
	return nil
}

func (m *EntryManager) getAppEntryByAppId(id string) *AppEntry {
	for _, entry := range m.Entries {
		if entry.appInfo != nil && id == entry.appInfo.GetId() {
			return entry
		}
	}
	return nil
}

func (m *EntryManager) getAppEntryByEntryId(id string) *AppEntry {
	for _, entry := range m.Entries {
		if entry.Id == id {
			return entry
		}
	}
	return nil
}

func (m *EntryManager) getAppEntryByInnerId(id string) *AppEntry {
	for _, entry := range m.Entries {
		if entry.innerId == id {
			return entry
		}
	}
	return nil
}

func (m *EntryManager) updateActiveWindow(win xproto.Window) {
	m.activeWindow = win
	entry := m.getAppEntryByWindow(win)
	if entry != nil {
		entry.setLeader(win)
		entry.current.updateWmName()
		entry.current.updateIcon()
	}
}

func (m *EntryManager) attachOrDetachWindow(winInfo *WindowInfo) {
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

func (m *EntryManager) initClientList() {
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

func (m *EntryManager) initDockedApps() {
	for _, id := range m.dockedAppManager.DockedAppList() {
		m.addDockedAppEntry(id)
	}
}

func (m *EntryManager) addAppEntry(e *AppEntry) {
	if m.getAppEntryByInnerId(e.innerId) != nil {
		logger.Debug("addAppEntry: entry exist, return")
		return
	}

	logger.Debug("addAppEntry: entry not exist, add new entry")

	m.Entries = append(m.Entries, e)

	// install on session D-Bus
	err := dbus.InstallOnSession(e)
	if err != nil {
		logger.Warning("Install AppEntry to dbus failed:", err)
		return
	}

	// emit signal Added
	entryObjPath := dbus.ObjectPath(entryDBusObjPathPrefix + e.Id)
	logger.Debugf("addAppEntry %v", entryObjPath)
	dbus.Emit(m, "Added", entryObjPath)
}

func (m *EntryManager) addDockedAppEntry(id string) *AppEntry {
	logger.Infof("Add docked app entry id: %q", id)
	appInfo := NewAppInfo(id)
	if appInfo == nil {
		logger.Warning("addDockedAppEntry failed: appInfo is nil")
		return nil
	}
	entryInnerId := appInfo.innerId
	var entry *AppEntry

	if e := m.getAppEntryByInnerId(entryInnerId); e != nil {
		entry = e
		appInfo.Destroy()
	} else {
		entry = newAppEntry(m, entryInnerId, appInfo)
		m.addAppEntry(entry)
	}

	entry.isDocked = true
	entry.updateMenu()
	entry.updateStatus()
	entry.updateTitle()
	entry.updateIcon()
	return entry
}

func (m *EntryManager) removeAppEntry(e *AppEntry) {
	for _, entry := range m.Entries {
		if entry == e {
			dbus.ReleaseName(e)
			dbus.UnInstallObject(e)

			entryId := entry.Id
			logger.Info("removeAppEntry id:", entryId)
			m.Entries = entrySliceRemove(m.Entries, e)
			e.destroy()
			// emit signal Removed
			dbus.Emit(m, "Removed", entryId)
			return
		}
	}
	logger.Warning("removeAppEntry failed, entry not found")
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

func (m *EntryManager) identifyWindow(winInfo *WindowInfo) (string, *AppInfo) {
	logger.Debugf("identifyWindow: window id: %v, window hash %v", winInfo.window, winInfo.innerId)
	desktopHash := m.windowDesktopMap.FindRel(winInfo.innerId)
	logger.Debug("identifyWindow: FindRel desktop hash:", desktopHash)
	var appInfo *AppInfo
	if desktopHash != "" {
		appInfo = m.desktopIdFileMap.FindRelAppInfo(desktopHash)
		logger.Debug("identifyWindow: FindRelAppInfo:", appInfo)
	}

	if appInfo == nil {
		// cache faild
		if desktopHash != "" {
			logger.Warning("winHash->DesktopHash success, but DesktopHash->appInfo fail, delRel")
			m.windowDesktopMap.DelRel(winInfo.innerId, desktopHash)
		}
		var canCache bool
		appInfo, canCache = m.getAppInfoFromWindow(winInfo)
		logger.Debug("identifyWindow: getAppInfoFromWindow:", appInfo)
		if appInfo != nil && canCache {
			m.windowDesktopMap.NewRel(winInfo.innerId, appInfo.innerId)
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

	m.windowDesktopMap.AutoSave()
	return entryInnerId, appInfo
}

func (m *EntryManager) attachWindow(winInfo *WindowInfo) {
	entryInnerId, appInfo := m.identifyWindow(winInfo)
	entry := m.getAppEntryByInnerId(entryInnerId)
	if entry != nil {
		logger.Debug("entry innerId exist")
		entry.attachWindow(winInfo)
		if appInfo != nil {
			appInfo.Destroy()
		}
		return
	}

	logger.Debug("entry innerId not exist, add new entry")
	entry = newAppEntryWithWindow(m, entryInnerId, winInfo, appInfo)
	m.addAppEntry(entry)
}

func (m *EntryManager) detachWindow(winInfo *WindowInfo) {
	entry := winInfo.entry
	if entry == nil {
		return
	}

	entry.detachWindow(winInfo)
	if !entry.isActive() && !entry.isDocked {
		m.removeAppEntry(entry)
		return
	}
	entry.updateIcon()
	entry.updateStatus()
	entry.updateAppXids()
	entry.updateMenu()
	entry.updateTitle()
}

func (m *EntryManager) getDockedAppList() []string {
	var list []string
	for _, entry := range m.Entries {
		if entry.appInfo != nil && entry.isDocked {
			appId := entry.appInfo.GetId()
			list = append(list, appId)
		}
	}
	return list
}

func createScratchDesktopFileWithAppEntry(entry *AppEntry) string {
	appId := "docked:" + entry.innerId

	if entry.appInfo != nil {
		desktopFile := entry.appInfo.GetFilePath()
		newPath := filepath.Join(scratchDir, appId+".desktop")
		// try link
		err := os.Link(desktopFile, newPath)
		if err != nil {
			logger.Warning("link failed try copy file contents")
			err = copyFileContents(desktopFile, newPath)
		}
		if err == nil {
			return appId
		} else {
			logger.Warning(err)
		}
	}

	title := entry.current.getDisplayName()
	// icon
	icon := entry.current.getIcon()
	if strings.HasPrefix(icon, "data:image") {
		path, err := dataUriToFile(icon, filepath.Join(scratchDir, appId+".png"))
		if err != nil {
			logger.Warning(err)
			icon = ""
		} else {
			icon = path
		}
	}
	if icon == "" {
		icon = "application-default-icon"
	}

	// cmd
	scriptContent := "#!/bin/sh\n" + entry.exec
	scriptFile := filepath.Join(scratchDir, appId+".sh")
	ioutil.WriteFile(scriptFile, []byte(scriptContent), 0744)
	cmd := scriptFile + " %U"

	err := createScratchDesktopFile(appId, title, icon, cmd)
	if err != nil {
		logger.Warning("createScratchDesktopFile failed:", err)
		return ""
	}
	return appId
}

func (m *EntryManager) requestDock(appId, title, icon, cmd string) bool {
	// create entry
	entry := m.addDockedAppEntry(appId)
	if entry == nil {
		err := createScratchDesktopFile(appId, title, icon, cmd)
		if err != nil {
			return false
		}
		entry = m.addDockedAppEntry(appId)
		if entry == nil {
			logger.Warning("addDockedAppEntry failed with scratch desktop")
			return false
		}
	}
	m.dockEntry(entry)
	return true
}

func (m *EntryManager) dockEntry(entry *AppEntry) {
	needScratchDesktop := false
	if entry.appInfo == nil {
		logger.Debug("dockEntry: entry.appInfo is nil")
		needScratchDesktop = true
	} else {
		// try create appInfo by desktopId
		desktopId := entry.appInfo.GetDesktopId()
		appInfo := gio.NewDesktopAppInfo(desktopId)
		if appInfo != nil {
			appInfo.Unref()
		} else {
			logger.Debugf("dockEntry: gio.NewDesktopAppInfo failed: desktop id %q", desktopId)
			needScratchDesktop = true
		}
	}

	logger.Debug("dockEntry: need scratch desktop?", needScratchDesktop)
	if needScratchDesktop {
		appId := createScratchDesktopFileWithAppEntry(entry)
		if appId != "" {
			entry.appInfo = NewAppInfo(appId)
			entryOldInnerId := entry.innerId
			entry.innerId = entry.appInfo.innerId
			logger.Debug("dockEntry: createScratchDesktopFile successed, entry use new innerId", entry.innerId)
			entryManager := entry.entryManager
			entryManager.windowDesktopMap.NewRel(entryOldInnerId, entry.innerId)
			entryManager.windowDesktopMap.AutoSave()
			entryManager.desktopIdFileMap.NewRel(entry.appInfo.GetFilePath(), entry.innerId)
			entryManager.desktopIdFileMap.AutoSave()
		} else {
			logger.Warning("createScratchDesktopFileWithAppEntry failed")
			return
		}
	}

	entry.isDocked = true
	entry.updateMenu()
	m.dockedAppManager.dockAppEntry(entry)
}

func (m *EntryManager) undockEntry(entry *AppEntry) {
	if entry.appInfo == nil {
		logger.Warning("undockEntry failed, entry.appInfo is nil")
		return
	}
	appId := entry.appInfo.GetId()

	if !entry.isActive() {
		m.removeAppEntry(entry)
	} else {
		dir := filepath.Dir(entry.appInfo.GetFilePath())
		if dir == scratchDir {
			removeScratchFiles(appId)
			// Re-identify Window
			if entry.current != nil {
				var newAppInfo *AppInfo
				entry.innerId, newAppInfo = m.identifyWindow(entry.current)
				entry.setAppInfo(newAppInfo)
			}
		}

		entry.isDocked = false
		entry.updateMenu()
	}
	m.dockedAppManager.undockAppEntry(appId)
}

func (m *EntryManager) undockEntryByAppId(appId string) bool {
	entry := m.getAppEntryByAppId(appId)
	if entry != nil {
		m.undockEntry(entry)
		return true
	}
	return false
}

func (m *EntryManager) getAppInfoFromWindow(winInfo *WindowInfo) (*AppInfo, bool) {
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
