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
	"sync"

	x "github.com/linuxdeepin/go-x11-client"
	"pkg.deepin.io/lib/dbusutil"
)

const (
	entryDBusObjPathPrefix = dbusPath + "/entries/"
	entryDBusInterface     = dbusInterface + ".Entry"
)

//go:generate dbusutil-gen -type AppEntry -import=github.com/linuxdeepin/go-x11-client=x app_entry.go

type AppEntry struct {
	PropsMu       sync.RWMutex
	Id            string
	IsActive      bool
	Name          string
	Icon          string
	Menu          string
	DesktopFile   string
	CurrentWindow x.Window
	IsDocked      bool
	// dbusutil-gen: equal=method:Equal
	WindowInfos windowInfosType

	service          *dbusutil.Service
	manager          *Manager
	innerId          string
	windows          map[x.Window]*WindowInfo
	current          *WindowInfo
	coreMenu         *Menu
	appInfo          *AppInfo
	winIconPreferred bool

	methods *struct {
		Activate       func() `in:"timestamp"`
		HandleMenuItem func() `in:"timestamp,id"`
		HandleDragDrop func() `in:"timestamp,files"`
		NewInstance    func() `in:"timestamp"`
	}
}

func newAppEntry(dockManager *Manager, innerId string, appInfo *AppInfo) *AppEntry {
	entry := &AppEntry{
		manager: dockManager,
		service: dockManager.service,
		Id:      dockManager.allocEntryId(),
		innerId: innerId,
		windows: make(map[x.Window]*WindowInfo),
	}
	entry.setAppInfo(appInfo)
	entry.Name = entry.getName()
	entry.Icon = entry.getIcon()
	return entry
}

func (entry *AppEntry) setAppInfo(newAppInfo *AppInfo) {
	if entry.appInfo == newAppInfo {
		logger.Debug("setAppInfo failed: old == new")
		return
	}
	entry.appInfo = newAppInfo

	if newAppInfo == nil {
		entry.winIconPreferred = true
		entry.setPropDesktopFile("")
	} else {
		entry.winIconPreferred = false
		entry.setPropDesktopFile(newAppInfo.GetFileName())
		id := newAppInfo.GetId()
		if strSliceContains(entry.manager.getWinIconPreferredApps(), id) {
			entry.winIconPreferred = true
			return
		}

		icon := newAppInfo.GetIcon()
		if icon == "" {
			entry.winIconPreferred = true
		}
	}
}

func (entry *AppEntry) hasWindow() bool {
	return len(entry.windows) != 0
}

func (entry *AppEntry) getWindowIds() []uint32 {
	list := make([]uint32, 0, len(entry.windows))
	for _, winInfo := range entry.windows {
		list = append(list, uint32(winInfo.window))
	}
	return list
}

func (entry *AppEntry) getWindowInfoSlice() []*WindowInfo {
	winInfoSlice := make([]*WindowInfo, 0, len(entry.windows))
	for _, winInfo := range entry.windows {
		winInfoSlice = append(winInfoSlice, winInfo)
	}
	return winInfoSlice
}

func (entry *AppEntry) getExec(oneLine bool) string {
	if entry.current == nil {
		return ""
	}
	winProcess := entry.current.process
	if winProcess != nil {
		if oneLine {
			return winProcess.GetOneCommandLine()
		} else {
			return winProcess.GetShellScriptLines()
		}
	}
	return ""
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

func (entry *AppEntry) setCurrentWindowInfo(winInfo *WindowInfo) {
	entry.current = winInfo
	if winInfo == nil {
		entry.setPropCurrentWindow(0)
	} else {
		entry.setPropCurrentWindow(winInfo.window)
	}
}

func (entry *AppEntry) findNextLeader() x.Window {
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

func (entry *AppEntry) attachWindow(winInfo *WindowInfo) bool {
	win := winInfo.window
	logger.Debugf("attach win %v to entry", win)

	winInfo.entry = entry

	entry.PropsMu.Lock()
	defer entry.PropsMu.Unlock()

	if _, ok := entry.windows[win]; ok {
		logger.Debugf("win %v is already attach to entry", win)
		return false
	}

	entry.windows[win] = winInfo
	entry.updateWindowInfos()
	entry.updateIsActive()

	if entry.current == nil {
		// from no window to has window
		entry.setCurrentWindowInfo(winInfo)
	}
	entry.updateIcon()
	entry.updateMenu()

	// print window info
	wmClassStr := "-"
	if winInfo.wmClass != nil {
		wmClassStr = fmt.Sprintf("%q %q", winInfo.wmClass.Class, winInfo.wmClass.Instance)
	}
	logger.Infof("attach window id: %d, wmClass: %s, wmState: %v,"+
		" wmWindowType: %v, wmAllowedActions: %v, hasXEmbedInfo: %v, hasWmTransientFor: %v",
		winInfo.window, wmClassStr,
		winInfo.wmState, winInfo.wmWindowType, winInfo.wmAllowedActions, winInfo.hasXEmbedInfo,
		winInfo.hasWmTransientFor)

	return true
}

// return need remove?
func (entry *AppEntry) detachWindow(winInfo *WindowInfo) bool {
	winInfo.entry = nil
	win := winInfo.window
	logger.Debug("detach window ", win)

	entry.PropsMu.Lock()
	defer entry.PropsMu.Unlock()

	delete(entry.windows, win)
	if len(entry.windows) == 0 && !entry.IsDocked {
		// no window and not docked
		return true
	}

	for _, winInfo := range entry.windows {
		// select first
		entry.setCurrentWindowInfo(winInfo)
		break
	}
	entry.updateWindowInfos()
	entry.updateIcon()
	entry.updateIsActive()
	entry.updateMenu()
	return false
}

func (entry *AppEntry) getName() (name string) {
	if entry.appInfo != nil {
		name = entry.appInfo.GetDisplayName()
	} else if entry.current != nil {
		name = entry.current.getDisplayName()
	}
	return
}

func (entry *AppEntry) updateName() {
	name := entry.getName()
	entry.setPropName(name)
}

func (entry *AppEntry) updateIcon() {
	icon := entry.getIcon()
	entry.setPropIcon(icon)
}

func (entry *AppEntry) getIcon() string {
	var icon string
	appInfo := entry.appInfo
	current := entry.current

	if entry.hasWindow() {
		if current == nil {
			logger.Warning("AppEntry.getIcon entry.hasWindow but entry.current is nil")
			return ""
		}

		// has window && current not nil
		if entry.winIconPreferred {
			// try current window icon first
			icon = current.getIcon()
			if icon != "" {
				return icon
			}
		}
		if appInfo != nil {
			icon = appInfo.GetIcon()
			if icon != "" {
				return icon
			}
		}
		return current.getIcon()

	} else if appInfo != nil {
		// no window
		return appInfo.GetIcon()
	}
	return ""
}

func (e *AppEntry) updateWindowInfos() {
	windowInfos := newWindowInfos()
	for win, winInfo := range e.windows {
		windowInfos[win] = ExportWindowInfo{
			Title: winInfo.Title,
			Flash: winInfo.hasWmStateDemandsAttention(),
		}
	}
	e.setPropWindowInfos(windowInfos)
}

func (e *AppEntry) updateIsActive() {
	_, ok := e.windows[e.manager.getActiveWindow()]
	e.setPropIsActive(ok)
}
