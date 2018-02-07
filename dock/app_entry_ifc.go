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
	"errors"

	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/icccm"
	"pkg.deepin.io/lib/dbus"
)

func (e *AppEntry) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dockManagerDBusDest,
		ObjectPath: entryDBusObjPathPrefix + e.Id,
		Interface:  entryDBusInterface,
	}
}

func (entry *AppEntry) Activate(timestamp uint32) error {
	logger.Debug("Activate timestamp:", timestamp)
	m := entry.dockManager
	if HideModeType(m.HideMode.Get()) == HideModeSmartHide {
		m.setPropHideState(HideStateShow)
		m.updateHideState(true)
	}

	if !entry.hasWindow() {
		entry.launchApp(timestamp)
		return nil
	}

	if entry.current == nil {
		err := errors.New("entry.current is nil")
		logger.Warning(err)
		return err
	}
	win := entry.current.window
	state, err := ewmh.WmStateGet(XU, win)
	if err != nil {
		logger.Warning("Get ewmh wmState failed win:", win)
		return err
	}

	if strSliceContains(state, "_NET_WM_STATE_FOCUSED") {
		s, err := icccm.WmStateGet(XU, win)
		if err != nil {
			logger.Warning("Get icccm WmState failed win:", win)
			return err
		}
		switch s.State {
		case icccm.StateIconic:
			activateWindow(win)
		case icccm.StateNormal:
			if len(entry.windows) == 1 {
				iconifyWindow(win)
			} else if entry.dockManager.getActiveWindow() == win {
				nextWin := entry.findNextLeader()
				activateWindow(nextWin)
			}
		}
	} else {
		activateWindow(win)
	}
	return nil
}

func (e *AppEntry) HandleMenuItem(timestamp uint32, id string) {
	logger.Debugf("HandleMenuItem id: %q timestamp: %v", id, timestamp)
	if e.coreMenu != nil {
		e.coreMenu.HandleAction(id, timestamp)
	} else {
		logger.Warning("HandleMenuItem failed: entry.coreMenu is nil")
	}
}

func (entry *AppEntry) HandleDragDrop(timestamp uint32, files []string) {
	logger.Debugf("handle drag drop files: %v, timestamp: %v", files, timestamp)
	ai := entry.appInfo
	if ai != nil {
		entry.dockManager.launch(ai.GetFileName(), timestamp, files)
	} else {
		logger.Warning("not supported")
	}
}

// RequestDock 驻留
func (entry *AppEntry) RequestDock() {
	if entry.dockManager != nil {
		entry.dockManager.dockEntry(entry)
	}
}

// RequestUndock 取消驻留
func (entry *AppEntry) RequestUndock() {
	if entry.dockManager != nil {
		entry.dockManager.undockEntry(entry)
	}
}

func (entry *AppEntry) PresentWindows() {
	if entry.dockManager != nil {
		windowIds := entry.getWindowIds()
		if len(windowIds) == 0 {
			return
		}
		entry.dockManager.wm.PresentWindows(windowIds)
	}
}

func (entry *AppEntry) NewInstance(timestamp uint32) {
	entry.launchApp(timestamp)
}

func (entry *AppEntry) Check() {
	for _, winInfo := range entry.windows {
		entry.dockManager.attachOrDetachWindow(winInfo)
	}
}

func (entry *AppEntry) ForceQuit() {
	for _, winInfo := range entry.windows {
		killClient(winInfo.window)
	}
}
