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
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	libPinyin "dbus/com/deepin/api/pinyin"
	libApps "dbus/com/deepin/daemon/apps"
	libLastore "dbus/com/deepin/lastore"
	"github.com/howeyc/fsnotify"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/notify"
)

const (
	lastoreDataDir    = "/var/lib/lastore/"
	desktopPkgMapFile = lastoreDataDir + "desktop_package.json"
	applicationsFile  = lastoreDataDir + "applications.json"

	ddeDataDir              = "/usr/share/dde/data/"
	appNameTranslationsFile = ddeDataDir + "app_name_translations.json"

	AppStatusCreated  = "created"
	AppStatusModified = "updated"
	AppStatusDeleted  = "deleted"
	lastoreDBusDest   = "com.deepin.lastore"
)

type Manager struct {
	items      map[string]*Item
	itemsMutex sync.Mutex

	launchedRecorder   *libApps.LaunchedRecorder
	desktopFileWatcher *libApps.DesktopFileWatcher
	notification       *notify.Notification
	lastoreManager     *libLastore.Manager
	pinyin             *libPinyin.Pinyin
	desktopPkgMap      map[string]string
	pkgCategoryMap     map[string]CategoryID
	nameMap            map[string]string

	searchTaskStack *searchTaskStack

	// TODO
	itemChanged    bool
	searchKeyMutex sync.Mutex
	currentRunes   []rune
	popPushOpChan  chan *popPushOp

	systemDBusConn *dbus.Conn

	appDirs            []string
	fsWatcher          *fsnotify.Watcher
	fsEventTimers      map[string]*time.Timer
	fsEventTimersMutex sync.Mutex
	// Signals:
	// SearchDone 返回搜索结果列表
	SearchDone     func([]string)
	ItemChanged    func(status string, itemInfo ItemInfo, categoryID CategoryID)
	NewAppLaunched func(string)
	// UninstallSuccess在卸载程序成功后触发。
	UninstallSuccess func(string)
	// UninstallFailed在卸载程序失败后触发。
	UninstallFailed func(string, string)
}

func NewManager() (*Manager, error) {
	m := &Manager{}
	err := m.init()
	if err != nil {
		m.destroy()
		return nil, err
	}
	return m, nil
}

func (m *Manager) getItemById(id string) *Item {
	m.itemsMutex.Lock()
	defer m.itemsMutex.Unlock()
	return m.items[id]
}

func (m *Manager) getAppIdByFilePath(file string) string {
	return getAppIdByFilePath(file, m.appDirs)
}

func (m *Manager) getItemByPath(path string) *Item {
	appId := m.getAppIdByFilePath(path)
	item := m.getItemById(appId)
	if item != nil && item.Path == path {
		return item
	}
	return nil
}

func (m *Manager) addItem(item *Item) {
	if item == nil {
		return
	}
	logger.Debugf("addItem path: %q", item.Path)

	item.ID = m.getAppIdByFilePath(item.Path)
	logger.Debugf("addItem id: %q", item.ID)

	// NOTE: change name before call item.setSearchTargets
	if m.nameMap != nil {
		newName := m.nameMap[item.ID]
		if newName != "" {
			item.Name = newName
		}
	}

	item.CategoryID = m.queryCategoryID(item)
	logger.Debug("addItem category", item.CategoryID)
	item.setSearchTargets(m.pinyin)
	logger.Debug("item search targets:", item.searchTargets)
	m.items[item.ID] = item
}

func (m *Manager) addItemWithLock(item *Item) {
	m.itemsMutex.Lock()
	m.addItem(item)
	m.itemsMutex.Unlock()
}

func (m *Manager) removeItem(id string) {
	m.itemsMutex.Lock()
	delete(m.items, id)
	m.itemsMutex.Unlock()
}

func (m *Manager) queryCategoryID(item *Item) CategoryID {
	pkg := m.queryPkgName(item)
	logger.Debugf("queryCategoryID desktopPkgMap item %v -> pkg %v", item, pkg)
	if pkg != "" && m.pkgCategoryMap != nil {
		if cid, ok := m.pkgCategoryMap[pkg]; ok {
			logger.Debugf("queryCategoryID pkgCategoryMap item %v -> category %v", item, cid)
			return cid
		}
	}

	categoryGuess := item.getXCategory()
	logger.Debugf("queryCategoryID categoryGuess item %v -> category %v", item, categoryGuess)
	return categoryGuess
}

func (m *Manager) queryPkgName(item *Item) string {
	if m.desktopPkgMap == nil {
		logger.Warning("queryPkgName failed: Manager.desktopPkgMap is nil")
		return ""
	}

	if pkg, ok := m.desktopPkgMap[item.ID]; ok {
		return pkg
	}
	// fail
	return ""
}

func (m *Manager) loadNameMap() error {
	file, err := os.Open(appNameTranslationsFile)
	if err != nil {
		return err
	}
	defer file.Close()
	dec := json.NewDecoder(bufio.NewReader(file))

	var data map[string](map[string]string)
	err = dec.Decode(&data)
	if err != nil {
		return err
	}

	lang := gettext.QueryLang()
	m.nameMap = data[lang]
	logger.Debugf("loadNameMap lang %v: %v", lang, m.nameMap)
	return nil
}

func (m *Manager) loadDesktopPkgMap() error {
	f, err := os.Open(desktopPkgMapFile)
	if err != nil {
		return err
	}
	defer f.Close()
	dec := json.NewDecoder(bufio.NewReader(f))

	var data map[string]string
	err = dec.Decode(&data)
	if err != nil {
		return err
	}
	m.desktopPkgMap = m.convertDesktopPkgMap(data)
	logger.Debug("loadDesktopPkgMap count:", len(m.desktopPkgMap))
	return nil
}

func (m *Manager) convertDesktopPkgMap(in map[string]string) map[string]string {
	out := make(map[string]string)
	for k, v := range in {
		if !filepath.IsAbs(k) {
			continue
		}
		if appId := m.getAppIdByFilePath(k); appId != "" {
			out[appId] = v
		}
	}
	return out
}

func (m *Manager) loadPkgCategoryMap() error {
	f, err := os.Open(applicationsFile)
	if err != nil {
		return err
	}
	defer f.Close()
	decoder := json.NewDecoder(bufio.NewReader(f))
	var jsonData map[string]struct{ Category string }
	if err := decoder.Decode(&jsonData); err != nil {
		return err
	}

	infos := make(map[string]CategoryID)
	for pkg, v := range jsonData {
		infos[pkg] = parseCategoryString(v.Category)
	}
	//logger.Debugf("loadPkgCategoryMap jsonData %#v", jsonData)
	//logger.Debugf("loadPkgCategoryMap infos %#v", infos)

	m.pkgCategoryMap = infos
	logger.Debug("loadPkgCategoryMap count:", len(infos))
	return nil
}

func (m *Manager) sendNotification(summary, body, icon string) {
	n := m.notification
	n.Update(summary, body, icon)
	go func() {
		err := n.Show()
		logger.Infof("sendNotification summary: %q, body: %q, icon: %q", summary, body, icon)
		if err != nil {
			logger.Warning("sendNotification failed:", err)
		}
	}()
}

func (m *Manager) emitSearchDone(result MatchResults) {
	var ids []string
	if result != nil {
		ids = result.Copy().GetTruncatedOrderedIDs()
	}
	dbus.Emit(m, "SearchDone", ids)
	logger.Debug("emit SearchDone", ids)
}
