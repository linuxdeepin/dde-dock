/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package dock

import (
	"fmt"
	"github.com/BurntSushi/xgbutil/ewmh"
	"pkg.deepin.io/lib/dbus"
	"sort"
	"time"
)

func (m *DockManager) allocEntryId() string {
	num := m.entryCount
	m.entryCount++
	return fmt.Sprintf("e%dT%x", num, getCurrentTimestamp())
}

func (m *DockManager) markAppLaunched(appInfo *AppInfo) {
	if appInfo == nil || m.launchedRecorder == nil {
		return
	}
	path := appInfo.GetFileName()
	m.launchedRecorder.MarkLaunched(path)
}

func (m *DockManager) attachOrDetachWindow(winInfo *WindowInfo) {
	win := winInfo.window
	showOnDock := m.isWindowRegistered(win) && m.clientList.Contains(win) &&
		isGoodWindow(win) && winInfo.canShowOnDock()
	logger.Debugf("win %v showOnDock? %v", win, showOnDock)
	entry := winInfo.entry
	if entry != nil {
		if !showOnDock {
			m.detachWindow(winInfo)
		}
	} else {

		if winInfo.entryInnerId == "" {
			winInfo.entryInnerId, winInfo.appInfo = m.identifyWindow(winInfo)
			go m.markAppLaunched(winInfo.appInfo)
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
	} else {
		logger.Debug("entry not existed, newAppEntry")
		entry = newAppEntry(m, entryInnerId, appInfo)
		m.Entries = m.Entries.Insert(entry, index)
		logger.Debugf("insert entry %v at %v", entry.Id, index)
		isNewAdded = true
	}
	return entry, isNewAdded
}

func (m *DockManager) appendDockedApp(app string) {
	logger.Debugf("appendDockedApp %q", app)
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

			entryId := entry.Id
			logger.Info("removeAppEntry id:", entryId)
			m.Entries = m.Entries.Remove(e)
			dbus.Emit(m, "EntryRemoved", entryId)

			go func() {
				time.Sleep(time.Second)
				dbus.UnInstallObject(e)
			}()
			return
		}
	}
	logger.Warning("removeAppEntry failed, entry not found")
}

func (m *DockManager) attachWindow(winInfo *WindowInfo) {
	entry, isNewAdded := m.addAppEntry(winInfo.entryInnerId, winInfo.appInfo, -1)
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
	entry.updateWindowTitles()
	if !entry.hasWindow() && !entry.IsDocked {
		m.removeAppEntry(entry)
		return
	}
	entry.updateIcon()
	entry.updateMenu()
	entry.updateIsActive()
}
