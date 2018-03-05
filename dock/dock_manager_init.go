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
	libApps "dbus/com/deepin/daemon/apps"
	"dbus/com/deepin/dde/daemon/launcher"
	libDDELauncher "dbus/com/deepin/dde/launcher"
	"dbus/com/deepin/sessionmanager"
	"dbus/com/deepin/wm"

	"gir/gio-2.0"
	"pkg.deepin.io/lib/gsettings"

	"github.com/BurntSushi/xgb/xproto"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil/gsprop"
)

const (
	ddeDataDir                = "/usr/share/dde/data"
	windowPatternsFile        = ddeDataDir + "/window_patterns.json"
	daemonLauncherServiceName = "com.deepin.dde.daemon.Launcher"
	daemonLauncherObjPath     = "/com/deepin/dde/daemon/Launcher"
	ddeLauncherServiceName    = "com.deepin.dde.Launcher"
	ddeLauncherObjPath        = "/com/deepin/dde/Launcher"
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
			m.service.StopExport(entry.GetDBusExportInfo())
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
	m.launcher.ConnectItemChanged(func(status string, itemInfo []interface{}, cid int64) {
		if len(itemInfo) > 2 && status == "deleted" {
			logger.Debugf("launcher item deleted %#v", itemInfo)
			// try desktop file path
			desktopFile, ok := itemInfo[0].(string)
			if !ok {
				logger.Warning("get item desktop file failed")
				return
			}
			dockedEntries := m.Entries.FilterDocked()
			for _, entry := range dockedEntries {
				file := entry.appInfo.GetFileName()
				if file == desktopFile {
					m.undockEntry(entry)
					return
				}
			}

			// try app id
			appId, ok := itemInfo[2].(string)
			if !ok {
				logger.Warning("get item app id failed")
				return
			}
			entry := getByAppId(dockedEntries, appId)
			if entry != nil {
				m.undockEntry(entry)
			}
		}
	})

	m.ddeLauncher.ConnectVisibleChanged(func(visible bool) {
		logger.Debug("dde launcher visible changed", visible)
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
	var err error
	m.settings = gio.NewSettings(dockSchema)
	m.HideMode = gsprop.NewEnum(m.settings, settingKeyHideMode)
	m.DisplayMode = gsprop.NewEnum(m.settings, settingKeyDisplayMode)
	m.Position = gsprop.NewEnum(m.settings, settingKeyPosition)
	m.IconSize = gsprop.NewUint(m.settings, settingKeyIconSize)
	m.ShowTimeout = gsprop.NewUint(m.settings, settingKeyShowTimeout)
	m.HideTimeout = gsprop.NewUint(m.settings, settingKeyHideTimeout)
	m.DockedApps = gsprop.NewStrv(m.settings, settingKeyDockedApps)

	m.FrontendWindowRect = NewRect()
	m.smartHideModeTimer = time.AfterFunc(10*time.Second, m.smartHideModeTimerExpired)
	m.smartHideModeTimer.Stop()

	m.listenSettingsChanged()

	m.windowInfoMap = make(map[xproto.Window]*WindowInfo)
	m.windowPatterns, err = loadWindowPatterns(windowPatternsFile)
	if err != nil {
		logger.Warning("loadWindowPatterns failed:", err)
	}

	m.wm, err = wm.NewWm("com.deepin.wm", "/com/deepin/wm")
	if err != nil {
		return err
	}

	m.launchedRecorder, err = libApps.NewLaunchedRecorder("com.deepin.daemon.Apps", "/com/deepin/daemon/Apps")
	if err != nil {
		return err
	}

	m.launcher, err = launcher.NewLauncher(daemonLauncherServiceName, daemonLauncherObjPath)
	if err != nil {
		return err
	}
	m.ddeLauncher, err = libDDELauncher.NewLauncher(ddeLauncherServiceName, ddeLauncherObjPath)
	if err != nil {
		return err
	}
	m.listenLauncherSignal()

	m.startManager, err = sessionmanager.NewStartManager("com.deepin.SessionManager", "/com/deepin/StartManager")
	if err != nil {
		return err
	}

	m.registerIdentifyWindowFuncs()
	m.initEntries()

	err = m.service.Export(m)
	if err != nil {
		return err
	}

	// 强制将 ClassicMode 转为 EfficientMode
	if m.DisplayMode.Get() == int32(DisplayModeClassicMode) {
		m.DisplayMode.Set(int32(DisplayModeEfficientMode))
	}

	return nil
}
