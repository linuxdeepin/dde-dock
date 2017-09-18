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
	"errors"
	"fmt"
	"sync"
	"time"

	// dbus interfaces:
	libApps "dbus/com/deepin/daemon/apps"
	"dbus/com/deepin/dde/daemon/launcher"
	"dbus/com/deepin/sessionmanager"
	"dbus/com/deepin/wm"
	ddbus "pkg.deepin.io/dde/daemon/dbus"

	"gir/gio-2.0"
	"pkg.deepin.io/lib/appinfo"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/dbus/property"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/ewmh"
)

type DockManager struct {
	clientList         windowSlice
	windowInfoMap      map[xproto.Window]*WindowInfo
	windowInfoMapMutex sync.RWMutex

	Entries AppEntries

	settings    *gio.Settings
	HideMode    *property.GSettingsEnumProperty `access:"readwrite"`
	DisplayMode *property.GSettingsEnumProperty `access:"readwrite"`
	Position    *property.GSettingsEnumProperty `access:"readwrite"`
	IconSize    *property.GSettingsUintProperty `access:"readwrite"`
	ShowTimeout *property.GSettingsUintProperty `access:"readwrite"`
	HideTimeout *property.GSettingsUintProperty `access:"readwrite"`
	DockedApps  *property.GSettingsStrvProperty

	activeWindow xproto.Window

	HideState HideStateType

	// TODO remove launchContext
	launchContext      *appinfo.AppLaunchContext
	smartHideModeTimer *time.Timer
	smartHideModeMutex sync.Mutex

	entryCount         uint
	FrontendWindowRect *Rect
	identifyWindowFuns []*IdentifyWindowFunc
	windowPatterns     WindowPatterns

	// dbus objects:
	launcher         *launcher.Launcher
	wm               *wm.Wm
	launchedRecorder *libApps.LaunchedRecorder
	startManager     *sessionmanager.StartManager

	// Signals
	ServiceRestarted func()
	EntryAdded       func(dbus.ObjectPath, int32)
	EntryRemoved     func(string)
}

const (
	dockSchema                     = "com.deepin.dde.dock"
	settingKeyHideMode             = "hide-mode"
	settingKeyDisplayMode          = "display-mode"
	settingKeyPosition             = "position"
	settingKeyIconSize             = "icon-size"
	settingKeyDockedApps           = "docked-apps"
	settingKeyShowTimeout          = "show-timeout"
	settingKeyHideTimeout          = "hide-timeout"
	settingKeyWinIconPreferredApps = "win-icon-preferred-apps"

	frontendWindowWmClass = "dde-dock"
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

	if m.wm != nil {
		wm.DestroyWm(m.wm)
		m.wm = nil
	}

	if m.launcher != nil {
		launcher.DestroyLauncher(m.launcher)
		m.launcher = nil
	}

	if m.launchedRecorder != nil {
		libApps.DestroyLaunchedRecorder(m.launchedRecorder)
		m.launchedRecorder = nil
	}

	if m.startManager != nil {
		sessionmanager.DestroyStartManager(m.startManager)
		m.startManager = nil
	}

	dbus.UnInstallObject(m)
}

func (m *DockManager) launch(file string, timestamp uint32) {
	_, err := m.startManager.LaunchWithTimestamp(file, timestamp)
	if err != nil {
		logger.Warningf("launch %q failed: %v", file, err)
	}
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

func (m *DockManager) MaximizeWindow(win uint32) error {
	err := m.ActivateWindow(win)
	if err != nil {
		return err
	}
	err = maximizeWindow(XU, xproto.Window(win))
	if err != nil {
		logger.Warning("maximize window failed:", err)
		return err
	}
	return nil
}

func (m *DockManager) MinimizeWindow(win uint32) error {
	err := minimizeWindow(XU, xproto.Window(win))
	if err != nil {
		logger.Warning("minimize window failed:", err)
		return err
	}
	return nil
}

func (m *DockManager) MakeWindowAbove(win uint32) error {
	err := m.ActivateWindow(win)
	if err != nil {
		return err
	}

	err = makeWindowAbove(XU, xproto.Window(win))
	if err != nil {
		logger.Warning("make window above failed:", err)
		return err
	}
	return nil
}

func (m *DockManager) MoveWindow(win uint32) error {
	err := m.ActivateWindow(win)
	if err != nil {
		return err
	}

	err = moveWindow(XU, xproto.Window(win))
	if err != nil {
		logger.Warning("move window failed:", err)
		return err
	}
	return nil
}

func (m *DockManager) PreviewWindow(win uint32) error {
	if !ddbus.IsSessionBusActivated(m.wm.DestName) {
		logger.Warning("Deepin window manager not running, unsupported this operation")
		return nil
	}
	return m.wm.PreviewWindow(win)
}

func (m *DockManager) CancelPreviewWindow() error {
	if !ddbus.IsSessionBusActivated(m.wm.DestName) {
		logger.Warning("Deepin window manager not running, unsupported this operation")
		return nil
	}
	return m.wm.CancelPreviewWindow()
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
	if m.FrontendWindowRect.X == x &&
		m.FrontendWindowRect.Y == y &&
		m.FrontendWindowRect.Width == width &&
		m.FrontendWindowRect.Height == height {
		logger.Debug("SetFrontendWindowRect no changed")
		return
	}
	m.FrontendWindowRect.X = x
	m.FrontendWindowRect.Y = y
	m.FrontendWindowRect.Width = width
	m.FrontendWindowRect.Height = height
	dbus.NotifyChange(m, "FrontendWindowRect")
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

func (m *DockManager) IsOnDock(desktopFilePath string) (bool, error) {
	entry, err := m.Entries.GetByDesktopFilePath(desktopFilePath)
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

func (m *DockManager) QueryWindowIdentifyMethod(wid uint32) (string, error) {
	for _, entry := range m.Entries {
		winInfo, ok := entry.windows[xproto.Window(wid)]
		if ok {
			if winInfo.appInfo != nil {
				return winInfo.appInfo.identifyMethod, nil
			} else {
				return "Failed", nil
			}
		}
	}
	return "", fmt.Errorf("window %d not found", wid)
}
