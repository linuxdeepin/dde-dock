/*
 * Copyright (C) 2016 ~ 2017 Deepin Technology Co., Ltd.
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

	"gir/gio-2.0"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/dbus/property"
	"pkg.deepin.io/lib/fsnotify"
	"pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/notify"
)

const (
	lastoreDataDir    = "/var/lib/lastore"
	desktopPkgMapFile = lastoreDataDir + "/desktop_package.json"
	applicationsFile  = lastoreDataDir + "/applications.json"

	ddeDataDir              = "/usr/share/dde/data/"
	appNameTranslationsFile = ddeDataDir + "app_name_translations.json"

	AppStatusCreated  = "created"
	AppStatusModified = "updated"
	AppStatusDeleted  = "deleted"
	lastoreDBusDest   = "com.deepin.lastore"

	gsSchemaLauncher  = "com.deepin.dde.launcher"
	gsKeyDisplayMode  = "display-mode"
	gsKeyFullscreen   = "fullscreen"
	gsKeyAppsUseProxy = "apps-use-proxy"
	gsKeyAppsHidden   = "apps-hidden"
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

	noPkgItemIDs       map[string]int
	appDirs            []string
	fsWatcher          *fsnotify.Watcher
	fsEventTimers      map[string]*time.Timer
	fsEventTimersMutex sync.Mutex
	settings           *gio.Settings
	appsHidden         []string
	appsHiddenMu       sync.Mutex
	// Properties:
	DisplayMode *property.GSettingsEnumProperty `access:"readwrite"`
	Fullscreen  *property.GSettingsBoolProperty `access:"readwrite"`

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

	var err error
	// init launchedRecorder
	m.launchedRecorder, err = libApps.NewLaunchedRecorder(appsDBusDest, appsDBusObjectPath)
	if err != nil {
		return nil, err
	}

	// init desktopFileWatcher
	m.desktopFileWatcher, err = libApps.NewDesktopFileWatcher(appsDBusDest, appsDBusObjectPath)
	if err != nil {
		m.destroy()
		return nil, err
	}

	// init system dbus conn
	m.systemDBusConn, err = dbus.SystemBus()
	if err != nil {
		m.destroy()
		return nil, err
	}

	// init lastoreManager
	m.lastoreManager, err = libLastore.NewManager(lastoreDBusDest, "/com/deepin/lastore")
	if err != nil {
		m.destroy()
		return nil, err
	}

	// init pinyin if lang is zh*
	if isZH() {
		m.pinyin, err = libPinyin.NewPinyin("com.deepin.api.Pinyin", "/com/deepin/api/Pinyin")
		if err != nil {
			m.destroy()
			return nil, err
		}
	}

	// init fsWatcher
	m.fsWatcher, err = fsnotify.NewWatcher()
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

func (m *Manager) setItemID(item *Item) {
	item.ID = m.getAppIdByFilePath(item.Path)
}

func (m *Manager) hiddenByGSettings(id string) bool {
	for _, appID := range m.appsHidden {
		if id == appID {
			return true
		}
	}
	return false
}

func (m *Manager) hiddenByGSettingsWithLock(id string) bool {
	m.appsHiddenMu.Lock()
	defer m.appsHiddenMu.Unlock()
	return m.hiddenByGSettings(id)
}

func (m *Manager) addItem(item *Item) {
	if item == nil {
		return
	}
	logger.Debugf("addItem path: %q, id: %q", item.Path, item.ID)

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
	pkg := m.queryPkgName(item.ID)
	if pkg == "" {
		m.noPkgItemIDs[item.ID] = 1
	}
	return m._queryCategoryID(item, pkg)
}

func (m *Manager) _queryCategoryID(item *Item, pkg string) CategoryID {
	logger.Debugf("queryCategoryID desktopPkgMap %v -> pkg %q", item, pkg)
	if pkg != "" {
		if cid, ok := m.pkgCategoryMap[pkg]; ok {
			logger.Debugf("queryCategoryID pkgCategoryMap %v -> %v", item, cid)
			return cid
		}
	}
	if cid, ok := parseCategoryString(item.xDeepinCategory); ok {
		logger.Debugf("queryCategoryID X-Deepin %v -> %v", item, cid)
		return cid
	}

	categoryGuess := item.getXCategory()
	logger.Debugf("queryCategoryID guess %v -> %v", item, categoryGuess)
	return categoryGuess
}

func (m *Manager) queryPkgName(itemID string) string {
	if m.desktopPkgMap == nil {
		logger.Warning("queryPkgName failed: Manager.desktopPkgMap is nil")
		return ""
	}
	return m.desktopPkgMap[itemID]
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

// get pkg->category map from applicationsFile
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
		cid, ok := parseCategoryString(v.Category)
		if !ok {
			logger.Warning("loadPkgCategoryMap: parse category %q failed", v.Category)
		}
		infos[pkg] = cid
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
