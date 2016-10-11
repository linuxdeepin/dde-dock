/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package launcher2

import (
	libPinyin "dbus/com/deepin/api/pinyin"
	libLastore "dbus/com/deepin/lastore"
	libStoreApi "dbus/com/deepin/store/api"
	libNotifications "dbus/org/freedesktop/notifications"
	"gir/gio-2.0"
	"github.com/howeyc/fsnotify"
	"os"
	"pkg.deepin.io/lib/dbus"
	"time"
)

func (m *Manager) init() error {
	var err error
	// init store.Api
	m.storeApi, err = libStoreApi.NewDStoreDesktop("com.deepin.store.Api", "/com/deepin/store/Api")
	if err != nil {
		return err
	}

	// init notifications
	m.notifier, err = libNotifications.NewNotifier("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	if err != nil {
		return err
	}

	// init system dbus conn
	m.systemDBusConn, err = dbus.SystemBus()
	if err != nil {
		return err
	}

	// init lastore
	m.lastore, err = libLastore.NewManager(lastoreDBusDest, "/com/deepin/lastore")
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

	allApps := gio.AppInfoGetAll()
	for _, app := range allApps {
		dAppInfo := gio.ToDesktopAppInfo(app)
		if !appShouldShow(dAppInfo) {
			continue
		}

		item := NewItemWithDesktopAppInfo(dAppInfo)
		m.addItem(item)
		app.Unref()
	}
	logger.Debug("load items count:", len(m.items))

	// init popPushOpChan
	m.popPushOpChan = make(chan *popPushOp, 50)
	go m.handlePopPushOps()

	// init searchTaskStack
	m.searchTaskStack = newSearchTaskStack(m)

	// watch item change
	m.watchApplicationsDirs()

	go m.handleFsWatcherEvents()
	return nil
}

func (m *Manager) watchApplicationsDirs() {
	m.fsEventTimers = make(map[string]*time.Timer)

	// create userAppsDir if it not exist
	userAppDir := getUserAppDir()
	if _, err := os.Stat(userAppDir); os.IsNotExist(err) {
		// userAppsDir does not exist
		os.MkdirAll(userAppDir, DirDefaultPerm)
	}

	for _, dir := range m.appDirs {
		m.addAppDir(dir, false)
	}
}

type popPushOp struct {
	popCount  int
	runesPush []rune
}

func (m *Manager) handlePopPushOps() {
	stack := m.searchTaskStack
	if stack == nil {
		logger.Warning("Manager.searchTaskStack is nil")
		return
	}
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
