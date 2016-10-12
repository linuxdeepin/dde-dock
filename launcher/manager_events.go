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
	"os"
	"path/filepath"
	"pkg.deepin.io/lib/dbus"
	"time"
)

const (
	desktopFilePattern = `[^.]*.desktop`
)

func isDesktopFile(path string) bool {
	basename := filepath.Base(path)
	matched, _ := filepath.Match(desktopFilePattern, basename)
	return matched
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

func (m *Manager) handleFsWatcherEvents() {
	watcher := m.fsWatcher
	for {
		select {
		case ev := <-watcher.Event:
			logger.Debugf("fsWatcher event: %v", ev)
			fileName := ev.Name
			if ev.IsCreate() {
				isDir, err := isDirectory(fileName)
				if err != nil {
					logger.Warning(err)
					break
				}

				if !isDir {
					m.delayHandleFileEvent(fileName)
					break
				}
				// dir created
				m.addAppDir(fileName, true)
			} else {
				m.delayHandleFileEvent(fileName)
			}
		case err := <-watcher.Error:
			logger.Warning("eventHandler error", err)
		}
	}
}

func shouldIgnoreDir(dir string) bool {
	basename := filepath.Base(dir)
	if basename == "menu-xdg" {
		return true
	}
	return false
}

func (m *Manager) addAppDir(path string, loadExisted bool) {
	watcher := m.fsWatcher
	logger.Debugf("add dir %q", path)
	walkfn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Warning(err)
			return nil
		}
		if info.IsDir() {
			if shouldIgnoreDir(path) {
				logger.Debugf("ignore dir %q", path)
				return nil
			}

			logger.Debugf("watch dir %q", path)
			err := watcher.Watch(path)
			if err != nil {
				logger.Warning(err)
			}
		} else if loadExisted {
			dir := filepath.Dir(path)
			if shouldIgnoreDir(dir) {
				return nil
			}

			if isDesktopFile(path) {
				m.delayHandleFileEvent(path)
			}
		}
		return nil
	}
	filepath.Walk(path, walkfn)
}

func (m *Manager) delayHandleFileEvent(name string) {
	m.fsEventTimersMutex.Lock()
	defer m.fsEventTimersMutex.Unlock()

	delay := 2000 * time.Millisecond
	timer, ok := m.fsEventTimers[name]
	if ok {
		timer.Stop()
		timer.Reset(delay)
	} else {
		m.fsEventTimers[name] = time.AfterFunc(delay, func() {
			if name == desktopPkgMapFile {
				err := m.loadDesktopPkgMap()
				if err != nil {
					logger.Warning(err)
				}
				return
			} else if name == applicationsFile {
				err := m.loadPkgCategoryMap()
				if err != nil {
					logger.Warning(err)
				}
				return
			} else if isDesktopFile(name) {
				m.checkDesktopFile(name)
			}
		})
	}
}

func (m *Manager) checkDesktopFile(file string) {
	logger.Debug("checkDesktopFile", file)
	appId := m.getAppIdByFilePath(file)
	logger.Debugf("app id %q", appId)
	if appId == "" {
		logger.Warningf("appId is empty, ignore file %q", file)
		return
	}

	item := m.getItemById(appId)
	appInfo := gio.NewDesktopAppInfo(appId + desktopExt)

	// 如果是新增加的目录里的desktop 用 gio.NewDesktopAppInfo 是找不到的，必须用完整路径
	if appInfo == nil {
		appInfo = gio.NewDesktopAppInfoFromFilename(file)
	}

	if appInfo == nil {
		logger.Warningf("appId %q appInfo is nil", appId)
		if item != nil {
			m.removeItem(appId)
			m.emitItemChanged(item, AppStatusDeleted)
		}
	} else {
		// appInfo is not nil
		shouldShow := appShouldShow(appInfo)
		defer appInfo.Unref()
		newItem := NewItemWithDesktopAppInfo(appInfo)

		// add or update item
		if item != nil {

			if shouldShow {
				// update item
				m.addItemWithLock(newItem)
				m.emitItemChanged(newItem, AppStatusModified)
			} else {
				m.removeItem(appId)
				m.emitItemChanged(newItem, AppStatusDeleted)
			}
		} else {
			if shouldShow {
				m.addItemWithLock(newItem)
				m.emitItemChanged(newItem, AppStatusCreated)
			}
		}
	}
}

func (m *Manager) emitItemChanged(item *Item, status string) {
	m.itemChanged = true
	itemInfo := item.newItemInfo()
	logger.Debugf("emit signal ItemChanged status: %v, itemInfo: %v", status, itemInfo)
	dbus.Emit(m, "ItemChanged", status, itemInfo, itemInfo.CategoryID)
}
