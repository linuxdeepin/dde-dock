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
	"encoding/json"
	"gir/gio-2.0"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/icccm"
	"pkg.deepin.io/lib/dbus"
	"sort"
)

const (
	entryDBusObjPathPrefix = "/dde/dock/entry/"
	entryDBusDestPrefix    = "dde.dock.entry."
	entryDBusInterface     = "dde.dock.Entry"

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
	entryManager *EntryManager
	// hashId string

	Id      string
	innerId string

	Type string
	Data map[string]string
	// Data Fields
	// Menu
	// Icon
	// Title

	// Signal
	DataChanged func(string, string)

	windows map[xproto.Window]*WindowInfo
	current *WindowInfo

	coreMenu *Menu
	exec     string
	path     string
	appInfo  *AppInfo
	isDocked bool
}

func newAppEntry(entryManager *EntryManager, id string, appInfo *AppInfo) *AppEntry {
	entry := &AppEntry{
		entryManager: entryManager,
		Id:           entryManager.allocEntryId(),
		innerId:      id,
		Type:         "App",
		Data:         make(map[string]string),
		windows:      make(map[xproto.Window]*WindowInfo),
		appInfo:      appInfo,
	}
	return entry
}

func newAppEntryWithWindow(entryManager *EntryManager, id string, winInfo *WindowInfo, appInfo *AppInfo) *AppEntry {
	if appInfo != nil {
		appId := appInfo.GetId()
		recordFrequency(appId)
		markAsLaunched(appId)
	}

	entry := newAppEntry(entryManager, id, appInfo)
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

func (entry *AppEntry) isActive() bool {
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

func (entry *AppEntry) Activate(x, y int32, timestamp uint32) (bool, error) {
	// x,y  useless
	logger.Debug("Activate timestamp:", timestamp)
	if !entry.isActive() {
		entry.launchApp(timestamp)
		return true, nil
	}

	if entry.current == nil {
		logger.Warning("entry.current is nil")
		return false, nil
	}
	win := entry.current.window
	state, err := ewmh.WmStateGet(XU, win)
	if err != nil {
		logger.Warning("Get ewmh wmState failed win:", win)
		return false, err
	}

	if contains(state, "_NET_WM_STATE_FOCUSED") {
		s, err := icccm.WmStateGet(XU, win)
		if err != nil {
			logger.Warning("Get icccm WmState failed win:", win)
			return false, err
		}
		switch s.State {
		case icccm.StateIconic:
			s.State = icccm.StateNormal
			logger.Debugf("set window %v state Iconic to Normal", win)
			icccm.WmStateSet(XU, win, s)
		case icccm.StateNormal:
			if len(entry.windows) == 1 {
				iconifyWindow(win)
			} else {
				if dockManager.activeWindow == win {
					nextWin := entry.findNextLeader()
					activateWindow(nextWin)
				}
			}
		}
	} else {
		activateWindow(win)
	}
	return true, nil
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
	entry.updateStatus()
	entry.updateAppXids()
	entry.updateMenu()

	winInfo.entry = entry

	if (entry.entryManager != nil && win == entry.entryManager.activeWindow) ||
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
	entry.entryManager = nil
	if entry.appInfo != nil {
		entry.appInfo.Destroy()
		entry.appInfo = nil
	}
}

func (e *AppEntry) setData(key, value string) {
	if e.Data[key] != value {
		logger.Debugf("setData %q : %v", key, value)
		e.Data[key] = value
		dbus.Emit(e, "DataChanged", key, value)
	}
}

func (e *AppEntry) getData(key string) string {
	return e.Data[key]
}

func (e *AppEntry) setTitle(title string) {
	e.setData(FieldTitle, title)
}

func (e *AppEntry) setIcon(icon string) {
	e.setData(FieldIcon, icon)
}

func (entry *AppEntry) updateTitle() {
	var title string
	if entry.isActive() {
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
	if entry.isActive() {
		icon = entry.current.getIcon()
	} else {
		icon = entry.appInfo.GetIcon()
	}
	entry.setIcon(icon)
}

func (entry *AppEntry) updateStatus() {
	var status string
	if entry.isActive() {
		status = ActiveStatus
	} else {
		status = NormalStatus
	}
	entry.setData(FieldStatus, status)
}

func (entry *AppEntry) updateAppXids() {
	xids := make([]XidInfo, 0)
	for win, winInfo := range entry.windows {
		xids = append(xids, XidInfo{uint32(win), winInfo.Title})
	}
	bytes, _ := json.Marshal(xids)
	entry.setData(FieldAppXids, string(bytes))
}

func (e *AppEntry) HandleMenuItem(id string, timestamp uint32) {
	logger.Debugf("HandleMenuItem id: %q timestamp: %v", id, timestamp)
	if e.coreMenu != nil {
		e.coreMenu.HandleAction(id, timestamp)
	} else {
		logger.Warning("HandleMenuItem failed: entry.coreMenu is nil")
	}
}

func (entry *AppEntry) HandleDragDrop(path string, timestamp uint32) {
	logger.Debugf("handle drag drop path: %q", path)
	ai := entry.appInfo
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

// 暂时无用或废弃
func (e *AppEntry) ContextMenu(x, y int32)                                    {}
func (e *AppEntry) SecondaryActivate(x, y int32, timestamp uint32)            {}
func (e *AppEntry) HandleDragEnter(x, y int32, data string, timestamp uint32) {}
func (e *AppEntry) HandleDragLeave(x, y int32, data string, timestamp uint32) {}
func (e *AppEntry) HandleDragOver(x, y int32, data string, timestamp uint32)  {}
func (e *AppEntry) HandleMouseWheel(x, y, delta int32, timestamp uint32)      {}
