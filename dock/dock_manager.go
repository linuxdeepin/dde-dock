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
	"dbus/com/deepin/wm"
	"errors"
	"gir/gio-2.0"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xrect"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/dbus/property"
	"sync"
	"time"
)

type DockManager struct {
	clientList                     windowSlice
	appIdFilterGroup               *AppIdFilterGroup
	desktopWindowsMapCacheManager  *desktopWindowsMapCacheManager
	desktopHashFileMapCacheManager *desktopHashFileMapCacheManager

	Entries AppEntries

	settings    *gio.Settings
	HideMode    *property.GSettingsEnumProperty `access:"readwrite"`
	DisplayMode *property.GSettingsEnumProperty `access:"readwrite"`
	Position    *property.GSettingsEnumProperty `access:"readwrite"`
	IconSize    *property.GSettingsUintProperty `access:"readwrite"`
	DockedApps  *property.GSettingsStrvProperty

	activeWindow xproto.Window

	HideState HideStateType

	smartHideModeTimer *time.Timer
	smartHideModeMutex sync.Mutex

	entryCount uint
	dockRect   xrect.Rect

	wm *wm.Wm

	// Signals
	ServiceRestarted func()
	EntryAdded       func(dbus.ObjectPath, int32)
	EntryRemoved     func(string)
}

const (
	dockSchema            = "com.deepin.dde.dock"
	settingKeyHideMode    = "hide-mode"
	settingKeyDisplayMode = "display-mode"
	settingKeyPosition    = "position"
	settingKeyIconSize    = "icon-size"
	settingKeyDockedApps  = "docked-apps"
)

func NewDockManager() (*DockManager, error) {
	m := new(DockManager)
	err := m.init()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *DockManager) destroy() {
	if m.smartHideModeTimer != nil {
		m.smartHideModeTimer.Stop()
		m.smartHideModeTimer = nil
	}

	if m.settings != nil {
		m.settings.Unref()
		m.settings = nil
	}

	dbus.UnInstallObject(m)
}

// ActivateWindow会激活给定id的窗口，被激活的窗口通常会成为焦点窗口。
func (m *DockManager) ActivateWindow(win uint32) error {
	err := activateWindow(xproto.Window(win))
	if err != nil {
		logger.Warning("Activate window failed:", err)
		return err
	}
	return nil
}

// CloseWindow会将传入id的窗口关闭。
func (m *DockManager) CloseWindow(win uint32) error {
	err := ewmh.CloseWindow(XU, xproto.Window(win))
	if err != nil {
		logger.Warning("Close window failed:", err)
		return err
	}
	return nil
}

// for debug
func (m *DockManager) GetEntryIDs() []string {
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

func (m *DockManager) SetFrontendWindowRect(x, y int32, width, height uint32) {
	m.dockRect.XSet(int(x))
	m.dockRect.YSet(int(y))
	m.dockRect.WidthSet(int(width))
	m.dockRect.HeightSet(int(height))
	m.updateHideStateWithoutDelay()
}

func (m *DockManager) IsDocked(desktopFilePath string) (bool, error) {
	entry, err := m.getDockedAppEntryByDesktopFilePath(desktopFilePath)
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

func (m *DockManager) RequestDock(desktopFilePath string, index int32) (bool, error) {
	appInfo := NewAppInfoFromFile(desktopFilePath)
	if appInfo == nil {
		return false, errors.New("Invalid desktopFilePath")
	}
	entry, isNewAdded := m.addAppEntry(appInfo.innerId, appInfo, int(index))
	dockResult := m.dockEntry(entry)
	if isNewAdded {
		entry.updateName()
		entry.updateIcon()
		m.installAppEntry(entry)
	}
	return dockResult, nil
}

func (m *DockManager) RequestUndock(desktopFilePath string) (bool, error) {
	entry, err := m.getDockedAppEntryByDesktopFilePath(desktopFilePath)
	if err != nil {
		return false, err
	}
	if entry == nil {
		return false, nil
	}
	m.undockEntry(entry)
	return true, nil
}

func (m *DockManager) MoveEntry(index, newIndex int32) error {
	entries, err := m.Entries.Move(int(index), int(newIndex))
	if err != nil {
		logger.Warning("MoveEntry failed:", err)
		return err
	}
	logger.Debug("MoveEntry ok")
	m.Entries = entries
	m.saveDockedApps()
	return nil
}
