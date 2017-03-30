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
	libApps "dbus/com/deepin/daemon/apps"
	"dbus/com/deepin/dde/daemon/launcher"
	"dbus/com/deepin/wm"
	"gir/gio-2.0"
	"github.com/BurntSushi/xgb/xproto"
	"pkg.deepin.io/lib/appinfo"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/dbus/property"
	"time"
)

const (
	ddeDataDir         = "/usr/share/dde/data"
	windowPatternsFile = ddeDataDir + "/window_patterns.json"
	launcherDest       = "com.deepin.dde.daemon.Launcher"
	launcherObjPath    = "/com/deepin/dde/daemon/Launcher"
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
		m.updateHideStateWithoutDelay()
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
	m.listenLauncherSignal()

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
