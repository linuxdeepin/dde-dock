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
	"time"

	// dbus interfaces
	libApps "dbus/com/deepin/daemon/apps"
	"dbus/com/deepin/dde/daemon/launcher"
	libDDELauncher "dbus/com/deepin/dde/launcher"
	"dbus/com/deepin/sessionmanager"
	"dbus/com/deepin/wm"

	"gir/gio-2.0"
	"pkg.deepin.io/lib/appinfo"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/dbus/property"

	"github.com/BurntSushi/xgb/xproto"
)

const (
	ddeDataDir           = "/usr/share/dde/data"
	windowPatternsFile   = ddeDataDir + "/window_patterns.json"
	launcherDest         = "com.deepin.dde.daemon.Launcher"
	launcherObjPath      = "/com/deepin/dde/daemon/Launcher"
	ddeLauncherDest      = "com.deepin.dde.Launcher"
	ddeLauncherInterface = ddeLauncherDest
	ddeLauncherObjPath   = "/com/deepin/dde/Launcher"
)

func (m *DockManager) initEntries() {
	m.initDockedApps()
	m.initClientList()
}

func (m *DockManager) connectSettingKeyChanged(key string, handler func(*gio.Settings, string)) {
	m.settings.Connect("changed::"+key, handler)
}

func (m *DockManager) listenSettingsChanged() {
	// listen hide mode change
	m.connectSettingKeyChanged(settingKeyHideMode, func(g *gio.Settings, key string) {
		mode := HideModeType(g.GetEnum(key))
		logger.Debug(key, "changed to", mode)
		m.updateHideState(false)
	})

	// listen display mode change
	m.connectSettingKeyChanged(settingKeyDisplayMode, func(g *gio.Settings, key string) {
		mode := DisplayModeType(g.GetEnum(key))
		logger.Debug(key, "changed to", mode)
	})

	// listen position change
	m.connectSettingKeyChanged(settingKeyPosition, func(g *gio.Settings, key string) {
		position := positionType(g.GetEnum(key))
		logger.Debug(key, "changed to", position)
	})
}

func (m *DockManager) listenLauncherSignal() {
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
			entry := dockedEntries.GetByAppId(appId)
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

func (m *DockManager) isDDELauncherVisible() bool {
	m.ddeLauncherVisibleMu.Lock()
	result := m.ddeLauncherVisible
	m.ddeLauncherVisibleMu.Unlock()
	return result
}

func (m *DockManager) getWinIconPreferredApps() []string {
	return m.settings.GetStrv(settingKeyWinIconPreferredApps)
}

func (m *DockManager) init() error {
	var err error

	m.launchContext = appinfo.NewAppLaunchContext(XU)
	m.settings = gio.NewSettings(dockSchema)

	m.HideMode = property.NewGSettingsEnumProperty(m, "HideMode", m.settings, settingKeyHideMode)
	m.DisplayMode = property.NewGSettingsEnumProperty(m, "DisplayMode", m.settings, settingKeyDisplayMode)
	m.Position = property.NewGSettingsEnumProperty(m, "Position", m.settings, settingKeyPosition)
	m.IconSize = property.NewGSettingsUintProperty(m, "IconSize", m.settings, settingKeyIconSize)
	m.ShowTimeout = property.NewGSettingsUintProperty(m, "ShowTimeout", m.settings, settingKeyShowTimeout)
	m.HideTimeout = property.NewGSettingsUintProperty(m, "HideTimeout", m.settings, settingKeyHideTimeout)
	m.DockedApps = property.NewGSettingsStrvProperty(m, "DockedApps", m.settings, settingKeyDockedApps)

	m.FrontendWindowRect = NewRect()
	m.smartHideModeTimer = time.AfterFunc(10*time.Second, m.smartHideModeTimerExpired)
	m.smartHideModeTimer.Stop()

	m.listenSettingsChanged()

	m.windowInfoMap = make(map[xproto.Window]*WindowInfo)
	m.windowPatterns, err = loadWindowPatterns(windowPatternsFile)
	if err != nil {
		logger.Warning("loadWindowPatterns failed:", err)
	}
	m.registerIdentifyWindowFuncs()
	m.initEntries()

	m.wm, err = wm.NewWm("com.deepin.wm", "/com/deepin/wm")
	if err != nil {
		return err
	}

	m.launchedRecorder, err = libApps.NewLaunchedRecorder("com.deepin.daemon.Apps", "/com/deepin/daemon/Apps")
	if err != nil {
		return err
	}

	m.launcher, err = launcher.NewLauncher(launcherDest, launcherObjPath)
	if err != nil {
		return err
	}
	m.ddeLauncher, err = libDDELauncher.NewLauncher(ddeLauncherDest, ddeLauncherObjPath)
	if err != nil {
		return err
	}
	m.listenLauncherSignal()

	m.startManager, err = sessionmanager.NewStartManager("com.deepin.SessionManager", "/com/deepin/StartManager")
	if err != nil {
		return err
	}

	err = dbus.InstallOnSession(m)
	if err != nil {
		return err
	}

	// 强制将 ClassicMode 转为 EfficientMode
	if m.DisplayMode.Get() == int32(DisplayModeClassicMode) {
		m.DisplayMode.Set(int32(DisplayModeEfficientMode))
	}

	return nil
}
