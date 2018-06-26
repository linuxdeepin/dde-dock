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
	"time"

	// dbus interfaces
	libApps "github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.apps"
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.dde.daemon.launcher"
	libDDELauncher "github.com/linuxdeepin/go-dbus-factory/com.deepin.dde.launcher"
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.sessionmanager"
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.wm"

	"gir/gio-2.0"
	x "github.com/linuxdeepin/go-x11-client"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/gsettings"
)

const (
	ddeDataDir         = "/usr/share/dde/data"
	windowPatternsFile = ddeDataDir + "/window_patterns.json"
)

func (m *Manager) initEntries() {
	m.initDockedApps()
	m.Entries.insertCb = func(entry *AppEntry, index int) {
		entryObjPath := dbus.ObjectPath(entryDBusObjPathPrefix + entry.Id)
		logger.Debug("entry added", entry.Id, index)
		m.service.Emit(m, "EntryAdded", entryObjPath, int32(index))
	}
	m.Entries.removeCb = func(entry *AppEntry) {
		m.service.Emit(m, "EntryRemoved", entry.Id)
		go func() {
			time.Sleep(time.Second)
			m.service.StopExport(entry)
		}()
	}
	m.initClientList()
}

func (m *Manager) connectSettingKeyChanged(key string, handler func(key string)) {
	gsettings.ConnectChanged(dockSchema, key, handler)
}

func (m *Manager) listenSettingsChanged() {
	// listen hide mode change
	m.connectSettingKeyChanged(settingKeyHideMode, func(key string) {
		mode := HideModeType(m.settings.GetEnum(key))
		logger.Debug(key, "changed to", mode)
		m.updateHideState(false)
	})

	// listen display mode change
	m.connectSettingKeyChanged(settingKeyDisplayMode, func(key string) {
		mode := DisplayModeType(m.settings.GetEnum(key))
		logger.Debug(key, "changed to", mode)
	})

	// listen position change
	m.connectSettingKeyChanged(settingKeyPosition, func(key string) {
		position := positionType(m.settings.GetEnum(key))
		logger.Debug(key, "changed to", position)
	})
}

func (m *Manager) listenLauncherSignal() {
	m.launcher.InitSignalExt(m.sessionSigLoop, true)
	m.launcher.ConnectItemChanged(func(status string, itemInfo launcher.ItemInfo,
		categoryID int64) {
		if status != "deleted" {
			return
		}
		// item deleted
		dockedEntries := m.Entries.FilterDocked()
		for _, entry := range dockedEntries {
			file := entry.appInfo.GetFileName()
			if file == itemInfo.Path {
				m.undockEntry(entry)
				return
			}
		}

		// try app id
		entry := getByAppId(dockedEntries, itemInfo.ID)
		if entry != nil {
			m.undockEntry(entry)
		}
	})

	m.ddeLauncher.InitSignalExt(m.sessionSigLoop, true)
	m.ddeLauncher.ConnectVisibleChanged(func(visible bool) {
		logger.Debug("dde-launcher visible changed", visible)
		m.ddeLauncherVisibleMu.Lock()
		m.ddeLauncherVisible = visible
		m.ddeLauncherVisibleMu.Unlock()

		m.updateHideState(false)
	})
}

func (m *Manager) isDDELauncherVisible() bool {
	m.ddeLauncherVisibleMu.Lock()
	result := m.ddeLauncherVisible
	m.ddeLauncherVisibleMu.Unlock()
	return result
}

func (m *Manager) getWinIconPreferredApps() []string {
	return m.settings.GetStrv(settingKeyWinIconPreferredApps)
}

func (m *Manager) init() error {
	m.rootWindow = globalXConn.GetDefaultScreen().Root

	var err error
	m.settings = gio.NewSettings(dockSchema)
	m.HideMode.Bind(m.settings, settingKeyHideMode)
	m.DisplayMode.Bind(m.settings, settingKeyDisplayMode)
	m.Position.Bind(m.settings, settingKeyPosition)
	m.IconSize.Bind(m.settings, settingKeyIconSize)
	m.ShowTimeout.Bind(m.settings, settingKeyShowTimeout)
	m.HideTimeout.Bind(m.settings, settingKeyHideTimeout)
	m.DockedApps.Bind(m.settings, settingKeyDockedApps)

	m.FrontendWindowRect = NewRect()
	m.smartHideModeTimer = time.AfterFunc(10*time.Second, m.smartHideModeTimerExpired)
	m.smartHideModeTimer.Stop()

	m.listenSettingsChanged()

	m.windowInfoMap = make(map[x.Window]*WindowInfo)
	m.windowPatterns, err = loadWindowPatterns(windowPatternsFile)
	if err != nil {
		logger.Warning("loadWindowPatterns failed:", err)
	}

	sessionBus := m.service.Conn()
	m.wm = wm.NewWm(sessionBus)

	systemBus, err := dbus.SystemBus()
	if err != nil {
		return err
	}
	m.appsObj = libApps.NewApps(systemBus)
	m.launcher = launcher.NewLauncher(sessionBus)
	m.ddeLauncher = libDDELauncher.NewLauncher(sessionBus)
	m.startManager = sessionmanager.NewStartManager(sessionBus)
	m.sessionSigLoop = dbusutil.NewSignalLoop(m.service.Conn(), 10)
	m.sessionSigLoop.Start()
	m.listenLauncherSignal()

	m.registerIdentifyWindowFuncs()
	m.initEntries()

	err = m.service.Export(dbusPath, m)
	if err != nil {
		return err
	}

	// 强制将 ClassicMode 转为 EfficientMode
	if m.DisplayMode.Get() == int32(DisplayModeClassicMode) {
		m.DisplayMode.Set(int32(DisplayModeEfficientMode))
	}

	go m.eventHandleLoop()
	return nil
}
