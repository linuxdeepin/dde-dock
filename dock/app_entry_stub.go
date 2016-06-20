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
	"gir/gio-2.0"
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

func (entry *AppEntry) Activate() error {
	timestamp := getCurrentTimestamp()
	logger.Debug("Activate timestamp:", timestamp)
	if !entry.hasWindow() {
		entry.launchApp(timestamp)
		return nil
	}

	if entry.current == nil {
		logger.Warning("entry.current is nil")
		return nil
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
	return nil
}

func (e *AppEntry) HandleMenuItem(id string) {
	timestamp := getCurrentTimestamp()
	logger.Debugf("HandleMenuItem id: %q timestamp: %v", id, timestamp)
	if e.coreMenu != nil {
		e.coreMenu.HandleAction(id, timestamp)
	} else {
		logger.Warning("HandleMenuItem failed: entry.coreMenu is nil")
	}
}

func (entry *AppEntry) HandleDragDrop(files []string) {
	logger.Debugf("handle drag drop files: %v", files)
	ai := entry.appInfo
	timestamp := getCurrentTimestamp()
	appLaunchContext := gio.GetGdkAppLaunchContext().SetTimestamp(timestamp)
	if ai.DesktopAppInfo != nil {
		_, err := ai.LaunchUris(files, appLaunchContext)
		if err != nil {
			logger.Warning("LaunchUris failed")
		}
	} else {
		logger.Warningf("no support!")
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
