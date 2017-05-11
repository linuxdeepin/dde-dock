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
	libPinyin "dbus/com/deepin/api/pinyin"
	libApps "dbus/com/deepin/daemon/apps"
	libLastore "dbus/com/deepin/lastore"
	"gir/gio-2.0"
	"github.com/howeyc/fsnotify"
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

func (m *Manager) init() error {
	m.settings = gio.NewSettings(gsSchemaLauncher)
	m.DisplayMode = property.NewGSettingsEnumProperty(m, "DisplayMode", m.settings, gsKeyDisplayMode)

	var err error
	// init launchedRecorder
	m.launchedRecorder, err = libApps.NewLaunchedRecorder(appsDBusDest, appsDBusObjectPath)
	if err != nil {
		return err
	}

	// init desktopFileWatcher
	m.desktopFileWatcher, err = libApps.NewDesktopFileWatcher(appsDBusDest, appsDBusObjectPath)
	if err != nil {
		return err
	}

	// init notification
	notify.Init(dbusDest)
	m.notification = notify.NewNotification("", "", "")

	// init system dbus conn
	m.systemDBusConn, err = dbus.SystemBus()
	if err != nil {
		return err
	}

	// init lastoreManager
	m.lastoreManager, err = libLastore.NewManager(lastoreDBusDest, "/com/deepin/lastore")
	if err != nil {
		return err
	}

	// init pinyin if lang is zh*
	if isZH() {
		m.pinyin, err = libPinyin.NewPinyin("com.deepin.api.Pinyin", "/com/deepin/api/Pinyin")
		if err != nil {
			return err
		}
	}

	// init fsWatcher
	m.fsWatcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	m.appDirs = getAppDirs()
	err = m.loadDesktopPkgMap()
	if err != nil {
		logger.Warning(err)
	}
	err = m.fsWatcher.Watch(desktopPkgMapFile)
	if err != nil {
		logger.Warning(err)
	}

	err = m.loadPkgCategoryMap()
	if err != nil {
		logger.Warning(err)
	}
	err = m.fsWatcher.Watch(applicationsFile)
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
	go m.handleFsWatcherEvents()
	m.desktopFileWatcher.ConnectEvent(func(filename string, _ uint32) {
		if shouldCheckDesktopFile(filename) {
			logger.Debug("DFWatcher event", filename)
			m.delayHandleFileEvent(filename)
		}
	})

	m.launchedRecorder.ConnectLaunched(func(path string) {
		item := m.getItemByPath(path)
		if item == nil {
			return
		}
		dbus.Emit(m, "NewAppLaunched", item.ID)
	})
	return nil
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
