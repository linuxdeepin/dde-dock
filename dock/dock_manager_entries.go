/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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
	"sort"

	"github.com/linuxdeepin/go-x11-client/util/wm/ewmh"
	"pkg.deepin.io/lib/dbus1"
)

func (m *Manager) allocEntryId() string {
	m.PropsMu.Lock()

	num := m.entryCount
	m.entryCount++

	m.PropsMu.Unlock()

	return fmt.Sprintf("e%dT%x", num, getCurrentTimestamp())
}

func (m *Manager) markAppLaunched(appInfo *AppInfo) {
	if appInfo == nil || m.launchedRecorder == nil {
		return
	}
	file := appInfo.GetFileName()
	logger.Debug("markAppLaunched", file)

	go func() {
		err := m.launchedRecorder.MarkLaunched(file)
		if err != nil {
			logger.Debug(err)
		}
	}()
}

func (m *Manager) attachOrDetachWindow(winInfo *WindowInfo) {
	win := winInfo.window

	isReg := m.isWindowRegistered(win)
	clientListContains := m.clientList.Contains(win)
	winInfoCanShow := winInfo.canShowOnDock()
	isGood := isGoodWindow(win)
	logger.Debugf("isReg: %v, client list contains: %v, winInfo can show: %v, isGood: %v",
		isReg, clientListContains, winInfoCanShow, isGood)

	showOnDock := isReg && clientListContains && isGood && winInfoCanShow
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

func (m *Manager) initClientList() {
	clientList, err := ewmh.GetClientList(globalXConn).Reply(globalXConn)
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

func (m *Manager) initDockedApps() {
	dockedApps := uniqStrSlice(m.DockedApps.Get())
	for _, app := range dockedApps {
		m.appendDockedApp(app)
	}
	m.saveDockedApps()
}

func (m *Manager) exportAppEntry(e *AppEntry) error {
	err := m.service.Export(dbus.ObjectPath(entryDBusObjPathPrefix+e.Id), e)
	if err != nil {
		logger.Warning("failed to export AppEntry:", err)
		return err
	}
	return nil
}

func (m *Manager) appendDockedApp(app string) {
	logger.Debugf("appendDockedApp %q", app)
	appInfo := NewDockedAppInfo(app)
	if appInfo == nil {
		logger.Warning("appendDockedApp failed: appInfo is nil")
		return
	}

	entry := newAppEntry(m, appInfo.innerId, appInfo)
	entry.setPropIsDocked(true)
	entry.updateMenu()
	err := m.exportAppEntry(entry)
	if err == nil {
		m.Entries.Append(entry)
	}
}

func (m *Manager) removeAppEntry(e *AppEntry) {
	logger.Info("removeAppEntry id:", e.Id)
	m.Entries.Remove(e)
}

func (m *Manager) attachWindow(winInfo *WindowInfo) {
	entry := m.Entries.GetByInnerId(winInfo.entryInnerId)

	if entry != nil {
		// existed
		entry.attachWindow(winInfo)
	} else {
		entry = newAppEntry(m, winInfo.entryInnerId, winInfo.appInfo)
		ok := entry.attachWindow(winInfo)
		if ok {
			err := m.exportAppEntry(entry)
			if err == nil {
				m.Entries.Append(entry)
			}
		}
	}
}

func (m *Manager) detachWindow(winInfo *WindowInfo) {
	entry := m.Entries.getByWindowId(winInfo.window)
	if entry == nil {
		return
	}
	needRemove := entry.detachWindow(winInfo)
	if needRemove {
		m.removeAppEntry(entry)
	}
}
