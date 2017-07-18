/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package launcher

import (
	"gir/gio-2.0"
	"path/filepath"
	"pkg.deepin.io/lib/appinfo/desktopappinfo"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/dbus/property"
	"pkg.deepin.io/lib/notify"
	"time"
)

const (
	appsDBusDest       = "com.deepin.daemon.Apps"
	appsDBusObjectPath = "/com/deepin/daemon/Apps"
)

func (m *Manager) init() {
	m.noPkgItemIDs = make(map[string]int)
	m.settings = gio.NewSettings(gsSchemaLauncher)
	m.DisplayMode = property.NewGSettingsEnumProperty(m, "DisplayMode", m.settings, gsKeyDisplayMode)
	m.Fullscreen = property.NewGSettingsBoolProperty(m, "Fullscreen", m.settings, gsKeyFullscreen)

	// init notification
	notify.Init(dbusDest)
	m.notification = notify.NewNotification("", "", "")

	m.appDirs = getAppDirs()
	err := m.loadDesktopPkgMap()
	if err != nil {
		logger.Warning(err)
	}

	err = m.loadPkgCategoryMap()
	if err != nil {
		logger.Warning(err)
	}

	err = m.loadNameMap()
	if err != nil {
		logger.Warning(err)
	}

	// load items
	m.items = make(map[string]*Item)

	skipDirs := make(map[string][]string)
	skipDirs["/usr/share/applications"] = []string{"screensavers"}
	userAppsDir := getUserAppDir()
	skipDirs[userAppsDir] = []string{"menu-xdg"}
	allApps := desktopappinfo.GetAll(skipDirs)
	for _, ai := range allApps {
		if !ai.IsExecutableOk() ||
			isDeepinCustomDesktopFile(ai.GetFileName()) {
			continue
		}
		item := NewItemWithDesktopAppInfo(ai)
		m.addItem(item)
	}
	logger.Debug("load items count:", len(m.items))

	// init searchTaskStack
	m.searchTaskStack = newSearchTaskStack(m)

	// init popPushOpChan
	m.popPushOpChan = make(chan *popPushOp, 50)
	go m.handlePopPushOps()

	m.fsEventTimers = make(map[string]*time.Timer)

	err = m.fsWatcher.Watch(lastoreDataDir)
	if err != nil {
		logger.Warning(err)
	}
	go m.handleFsWatcherEvents()
	m.desktopFileWatcher.ConnectEvent(func(filename string, _ uint32) {
		if shouldCheckDesktopFile(filename) {
			logger.Debug("DFWatcher event", filename)
			m.delayHandleFileEvent(filename)
		}
	})

	m.launchedRecorder.WatchDirs(getDataDirsForWatch())

	m.launchedRecorder.ConnectServiceRestarted(func() {
		if m.launchedRecorder != nil {
			m.launchedRecorder.WatchDirs(getDataDirsForWatch())
		}
	})
	m.launchedRecorder.ConnectLaunched(func(path string) {
		item := m.getItemByPath(path)
		if item == nil {
			return
		}
		dbus.Emit(m, "NewAppLaunched", item.ID)
	})
}

func shouldCheckDesktopFile(filename string) bool {
	dir, basename := filepath.Split(filename)
	dir = filepath.Clean(dir)
	matched, _ := filepath.Match(desktopFilePattern, basename)
	if !matched {
		return false
	}

	// ignore $HOME/.local/share/applications/menu-xdg/
	skipDir := filepath.Join(getUserAppDir(), "menu-xdg")
	if dir == skipDir {
		return false
	}
	return true
}

type popPushOp struct {
	popCount  int
	runesPush []rune
}

func (m *Manager) handlePopPushOps() {
	stack := m.searchTaskStack
	for op := range m.popPushOpChan {
		logger.Debug("op:", op)

		for i := 0; i < op.popCount; i++ {
			stack.Pop()
		}
		if len(op.runesPush) == 0 {
			// emit top result
			top := stack.topTask()
			if top == nil {
				logger.Debug("emit SearchDone []")
				dbus.Emit(m, "SearchDone", []string{})
			} else {
				top.emitResult()
			}
		}
		for _, v := range op.runesPush {
			stack.Push(v)
		}
	}
}

func (m *Manager) destroy() {
	// TODO
}
