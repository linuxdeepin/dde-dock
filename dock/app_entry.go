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
	"github.com/BurntSushi/xgb/xproto"
	"pkg.deepin.io/lib/dbus"
	"sort"
	"sync"
)

const (
	entryDBusObjPathPrefix = dockManagerDBusObjPath + "/entries/"
	entryDBusInterface     = dockManagerDBusInterface + ".Entry"

	FieldTitle   = "title"
	FieldIcon    = "icon"
	FieldMenu    = "menu"
	FieldAppXids = "app-xids"

	FieldStatus   = "app-status"
	ActiveStatus  = "active"
	NormalStatus  = "normal"
	InvalidStatus = "invalid"
)

type XidInfo struct {
	Xid   uint32
	Title string
}

type AppEntry struct {
	dockManager *DockManager

	Id      string
	innerId string

	IsActive bool
	Title    string
	Icon     string
	Menu     string

	WindowTitles windowTitlesType
	windows      map[xproto.Window]*WindowInfo
	current      *WindowInfo

	coreMenu  *Menu
	exec      string
	path      string
	appInfo   *AppInfo
	IsDocked  bool
	dockMutex sync.Mutex
}

func newAppEntry(dockManager *DockManager, id string, appInfo *AppInfo) *AppEntry {
	entry := &AppEntry{
		dockManager:  dockManager,
		Id:           dockManager.allocEntryId(),
		innerId:      id,
		WindowTitles: newWindowTitles(),
		windows:      make(map[xproto.Window]*WindowInfo),
		appInfo:      appInfo,
	}
	return entry
}

func newAppEntryWithWindow(dockManager *DockManager, id string, winInfo *WindowInfo, appInfo *AppInfo) *AppEntry {
	if appInfo != nil {
		appId := appInfo.GetId()
		recordFrequency(appId)
		markAsLaunched(appId)
	}

	entry := newAppEntry(dockManager, id, appInfo)
	entry.initExec(winInfo)

	entry.current = winInfo
	entry.attachWindow(winInfo)
	winInfo.updateWmName()
	winInfo.updateIcon()
	return entry
}

func (entry *AppEntry) setAppInfo(newAppInfo *AppInfo) {
	if newAppInfo == nil {
		logger.Debug("setAppInfo failed: newAppInfo is nil")
		return
	}
	entry.appInfo.Destroy()
	entry.appInfo = newAppInfo
}

func (entry *AppEntry) hasWindow() bool {
	return len(entry.windows) != 0
}

func (entry *AppEntry) initExec(winInfo *WindowInfo) {
	ai := entry.appInfo
	if ai != nil && ai.DesktopAppInfo != nil {
		entry.exec = ai.DesktopAppInfo.GetCommandline()
	}

	if winInfo.process != nil {
		entry.exec = winInfo.process.GetShellScript()
	}

	logger.Debug("initExec:", entry.exec)
}

func (entry *AppEntry) getDisplayName() string {
	if entry.appInfo != nil {
		return entry.appInfo.GetDisplayName()
	}
	if entry.current != nil {
		return entry.current.getDisplayName()
	}
	return ""
}

func (entry *AppEntry) setLeader(leader xproto.Window) {
	if info, ok := entry.windows[leader]; ok {
		entry.current = info
	}
}

func (entry *AppEntry) findNextLeader() xproto.Window {
	winSlice := make(windowSlice, 0, len(entry.windows))
	for win, _ := range entry.windows {
		winSlice = append(winSlice, win)
	}
	sort.Sort(winSlice)
	currentWin := entry.current.window
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

func (entry *AppEntry) attachWindow(winInfo *WindowInfo) {
	win := winInfo.window
	logger.Debugf("attach win %v to entry", win)

	if _, ok := entry.windows[win]; ok {
		logger.Debugf("win %v is already attach to entry", win)
		return
	}

	entry.windows[win] = winInfo
	entry.updateWindowTitles()
	entry.updateMenu()
	entry.updateIsActive()

	winInfo.entry = entry

	if (entry.dockManager != nil && win == entry.dockManager.activeWindow) ||
		entry.current == nil {
		entry.current = winInfo
		winInfo.updateWmName()
		winInfo.updateIcon()
	}
}

func (entry *AppEntry) detachWindow(winInfo *WindowInfo) {
	win := winInfo.window
	if _, ok := entry.windows[win]; ok {
		if len(entry.windows) > 1 {
			// switch current to next window
			entry.setLeader(entry.findNextLeader())
		}
		delete(entry.windows, win)
	}
}

func (entry *AppEntry) destroy() {
	entry.dockManager = nil
	if entry.appInfo != nil {
		entry.appInfo.Destroy()
		entry.appInfo = nil
	}
}

func (e *AppEntry) setTitle(title string) {
	if e.Title != title {
		e.Title = title
		dbus.NotifyChange(e, "Title")
	}
}

func (e *AppEntry) setIcon(icon string) {
	if e.Icon != icon {
		e.Icon = icon
		dbus.NotifyChange(e, "Icon")
	}
}

func (e *AppEntry) setIsActive(isActive bool) {
	if e.IsActive != isActive {
		e.IsActive = isActive
		dbus.NotifyChange(e, "IsActive")
	}
}

func (e *AppEntry) setIsDocked(isDocked bool) {
	if e.IsDocked != isDocked {
		e.IsDocked = isDocked
		dbus.NotifyChange(e, "IsDocked")
	}
}

func (entry *AppEntry) updateTitle() {
	var title string
	if entry.hasWindow() {
		title = entry.current.getTitle()
	} else if entry.appInfo != nil {
		title = entry.appInfo.GetDisplayName()
	} else {
		logger.Debug("updateTitle failed, entry is not active and entry.appInfo is nil")
		return
	}
	entry.setTitle(title)
}

func (entry *AppEntry) updateIcon() {
	var icon string
	if entry.hasWindow() {
		icon = entry.current.getIcon()
	} else {
		icon = entry.appInfo.GetIcon()
	}
	entry.setIcon(icon)
}

func (e *AppEntry) updateWindowTitles() {
	windowTitles := newWindowTitles()
	for win, winInfo := range e.windows {
		windowTitles[win] = winInfo.Title
	}
	if !e.WindowTitles.Equal(windowTitles) {
		e.WindowTitles = windowTitles
		dbus.NotifyChange(e, "WindowTitles")
	}
}

func (e *AppEntry) updateIsActive() {
	if e.dockManager == nil {
		return
	}
	_, ok := e.windows[e.dockManager.activeWindow]
	e.setIsActive(ok)
}
