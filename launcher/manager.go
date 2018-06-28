/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
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
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	libPinyin "github.com/linuxdeepin/go-dbus-factory/com.deepin.api.pinyin"
	libApps "github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.apps"
	libLastore "github.com/linuxdeepin/go-dbus-factory/com.deepin.lastore"

	"gir/gio-2.0"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/gsprop"
	"pkg.deepin.io/lib/fsnotify"
	"pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/notify"
	"pkg.deepin.io/lib/strv"
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

	gsSchemaLauncher        = "com.deepin.dde.launcher"
	gsKeyDisplayMode        = "display-mode"
	gsKeyFullscreen         = "fullscreen"
	gsKeyAppsUseProxy       = "apps-use-proxy"
	gsKeyAppsDisableScaling = "apps-disable-scaling"
	gsKeyAppsHidden         = "apps-hidden"
)

type Manager struct {
	service    *dbusutil.Service
	sysSigLoop *dbusutil.SignalLoop
	items      map[string]*Item
	itemsMutex sync.Mutex

	appsObj        *libApps.Apps
	notification   *notify.Notification
	lastore        *libLastore.Lastore
	pinyin         *libPinyin.Pinyin
	desktopPkgMap  map[string]string
	pkgCategoryMap map[string]CategoryID
	nameMap        map[string]string

	searchTaskStack *searchTaskStack

	itemsChangedHit uint32
	searchMu        sync.Mutex
	currentRunes    []rune
	popPushOpChan   chan *popPushOp

	noPkgItemIDs       map[string]int
	appDirs            []string
	fsWatcher          *fsnotify.Watcher
	fsEventTimers      map[string]*time.Timer
	fsEventTimersMutex sync.Mutex
	settings           *gio.Settings
	appsHidden         []string
	appsHiddenMu       sync.Mutex
	// Properties:
	DisplayMode gsprop.Enum `prop:"access:rw"`
	Fullscreen  gsprop.Bool `prop:"access:rw"`

	signals *struct {
		// SearchDone 返回搜索结果列表
		SearchDone struct {
			apps []string
		}

		ItemChanged struct {
			status     string
			itemInfo   ItemInfo
			categoryID CategoryID
		}

		NewAppLaunched struct {
			appID string
		}

		// UninstallSuccess在卸载程序成功后触发。
		UninstallSuccess struct {
			appID string
		}

		// UninstallFailed在卸载程序失败后触发。
		UninstallFailed struct {
			appId  string
			errMsg string
		}
	}

	methods *struct {
		GetAllItemInfos          func() `out:"itemInfoList"`
		GetItemInfo              func() `in:"id" out:"itemInfo"`
		GetAllNewInstalledApps   func() `out:"apps"`
		IsItemOnDesktop          func() `in:"id" out:"result"`
		RequestRemoveFromDesktop func() `in:"id" out:"ok"`
		RequestSendToDesktop     func() `in:"id" out:"ok"`
		MarkLaunched             func() `in:"id"`
		RequestUninstall         func() `in:"id,purge"`
		Search                   func() `in:"key"`
		GetUseProxy              func() `in:"id" out:"value"`
		SetUseProxy              func() `in:"id,value"`
		GetDisableScaling        func() `in:"id" out:"value"`
		SetDisableScaling        func() `in:"id,value"`
	}
}

func NewManager(service *dbusutil.Service) (*Manager, error) {
	m := &Manager{
		service: service,
	}

	systemBus, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	m.appsObj = libApps.NewApps(systemBus)
	m.lastore = libLastore.NewLastore(systemBus)
	if isZH() {
		m.pinyin = libPinyin.NewPinyin(service.Conn())
	}

	// init fsWatcher
	m.fsWatcher, err = fsnotify.NewWatcher()
	if err == nil {
		err = m.fsWatcher.Watch(lastoreDataDir)
		if err == nil {
			go m.handleFsWatcherEvents()
		} else {
			logger.Warning(err)
		}
	} else {
		logger.Warning("failed to init fsWatcher:", err)
	}

	m.settings = gio.NewSettings(gsSchemaLauncher)
	m.DisplayMode.Bind(m.settings, gsKeyDisplayMode)
	m.Fullscreen.Bind(m.settings, gsKeyFullscreen)

	m.noPkgItemIDs = make(map[string]int)

	m.appsHidden = m.settings.GetStrv(gsKeyAppsHidden)
	logger.Debug("appsHidden: ", m.appsHidden)
	m.listenSettingsChanged()

	// init notification
	notify.Init(dbusServiceName)
	m.notification = notify.NewNotification("", "", "")

	m.appDirs = getAppDirs()
	err = m.loadDesktopPkgMap()
	if err != nil {
		logger.Warning(err)
	}

	err = m.loadPkgCategoryMap()
	if err != nil {
		logger.Warning(err)
	}

	// load name map
	err = m.loadNameMap()
	if err != nil {
		logger.Warning(err)
	}
	m.initItems()

	// init searchTaskStack
	m.searchTaskStack = newSearchTaskStack(m)

	// init popPushOpChan
	m.popPushOpChan = make(chan *popPushOp, 50)
	go m.handlePopPushOps()

	m.fsEventTimers = make(map[string]*time.Timer)

	m.sysSigLoop = dbusutil.NewSignalLoop(systemBus, 100)
	m.sysSigLoop.Start()
	m.appsObj.InitSignalExt(m.sysSigLoop, true)
	m.appsObj.ConnectEvent(func(filename string, _ uint32) {
		if shouldCheckDesktopFile(filename) {
			logger.Debug("DFWatcher event", filename)
			m.delayHandleFileEvent(filename)
		}
	})

	m.appsObj.WatchDirs(0, getDataDirsForWatch())

	m.appsObj.ConnectServiceRestarted(func() {
		if m.appsObj != nil {
			m.appsObj.WatchDirs(0, getDataDirsForWatch())
		}
	})
	m.appsObj.ConnectLaunched(func(path string) {
		item := m.getItemByPath(path)
		if item == nil {
			return
		}
		m.service.Emit(m, "NewAppLaunched", item.ID)
	})
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
	m.service.Emit(m, "SearchDone", ids)
	logger.Debug("emit SearchDone", ids)
}

func (m *Manager) getUseFeature(key, id string) (bool, *dbus.Error) {
	item := m.getItemById(id)
	if item == nil {
		return false, dbusutil.ToError(errorInvalidID)
	}
	apps := strv.Strv(m.settings.GetStrv(key))
	return apps.Contains(id), nil
}

func (m *Manager) setUseFeature(key, id string, val bool) *dbus.Error {
	item := m.getItemById(id)
	if item == nil {
		return dbusutil.ToError(errorInvalidID)
	}
	apps := strv.Strv(m.settings.GetStrv(key))

	var changed bool
	if val {
		apps, changed = apps.Add(id)
	} else {
		apps, changed = apps.Delete(id)
	}

	if !changed {
		return nil
	}

	ok := m.settings.SetStrv(key, []string(apps))
	if !ok {
		return dbusutil.ToError(fmt.Errorf("gsettings set %s failed", key))
	}
	return nil
}
