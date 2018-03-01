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
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	"pkg.deepin.io/dde/api/soundutils"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/keyfile"
)

const (
	dbusServiceName    = "com.deepin.dde.daemon.Launcher"
	dbusObjPath        = "/com/deepin/dde/daemon/Launcher"
	dbusInterface      = dbusServiceName
	desktopMainSection = "Desktop Entry"
)

var errorInvalidID = errors.New("invalid ID")

func (m *Manager) GetDBusExportInfo() dbusutil.ExportInfo {
	return dbusutil.ExportInfo{
		Path:      dbusObjPath,
		Interface: dbusInterface,
	}
}

func (m *Manager) GetAllItemInfos() ([]ItemInfo, *dbus.Error) {
	list := make([]ItemInfo, 0, len(m.items))
	for _, item := range m.items {
		list = append(list, item.newItemInfo())
	}
	logger.Debug("GetAllItemInfos list length:", len(list))
	return list, nil
}

func (m *Manager) GetItemInfo(id string) (ItemInfo, *dbus.Error) {
	item := m.getItemById(id)
	if item == nil {
		return ItemInfo{}, dbusutil.ToError(errorInvalidID)
	}
	return item.newItemInfo(), nil
}

func (m *Manager) GetAllNewInstalledApps() ([]string, *dbus.Error) {
	newApps, err := m.launchedRecorder.GetNew()
	if err != nil {
		return nil, dbusutil.ToError(err)
	}
	var ids []string
	// newApps type is map[string][]string
	for dir, names := range newApps {
		for _, name := range names {
			path := filepath.Join(dir, name) + desktopExt
			if item := m.getItemByPath(path); item != nil {
				ids = append(ids, item.ID)
			}
		}
	}
	return ids, nil
}

func (m *Manager) IsItemOnDesktop(id string) (bool, *dbus.Error) {
	item := m.getItemById(id)
	if item == nil {
		return false, dbusutil.ToError(errorInvalidID)
	}
	file := appInDesktop(m.getAppIdByFilePath(item.Path))
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			// not exist
			return false, nil
		} else {
			return false, dbusutil.ToError(err)
		}
	} else {
		// exist
		return true, nil
	}
}

func (m *Manager) RequestRemoveFromDesktop(id string) (bool, *dbus.Error) {
	item := m.getItemById(id)
	if item == nil {
		return false, dbusutil.ToError(errorInvalidID)
	}
	file := appInDesktop(m.getAppIdByFilePath(item.Path))
	err := os.Remove(file)
	return err == nil, dbusutil.ToError(err)
}

func (m *Manager) RequestSendToDesktop(id string) (bool, *dbus.Error) {
	item := m.getItemById(id)
	if item == nil {
		return false, dbusutil.ToError(errorInvalidID)
	}
	dest := appInDesktop(m.getAppIdByFilePath(item.Path))
	_, err := os.Stat(dest)
	if err != nil {
		if !os.IsNotExist(err) {
			return false, dbusutil.ToError(err)
		}
		// dest file not exist
	} else {
		// dest file exist
		return false, dbusutil.ToError(os.ErrExist)
	}

	kf := keyfile.NewKeyFile()
	if err := kf.LoadFromFile(item.Path); err != nil {
		logger.Warning(err)
		return false, dbusutil.ToError(err)
	}
	kf.SetString(desktopMainSection, "X-Deepin-CreatedBy", dbusServiceName)
	kf.SetString(desktopMainSection, "X-Deepin-AppID", id)
	// Desktop files in user desktop direcotry do not require executable permission
	if err := kf.SaveToFile(dest); err != nil {
		logger.Warning("save new desktop file failed:", err)
		return false, dbusutil.ToError(err)
	}
	// success
	go soundutils.PlaySystemSound(soundutils.EventIconToDesktop, "")
	return true, nil
}

// MarkLaunched 废弃
func (m *Manager) MarkLaunched(id string) *dbus.Error {
	return nil
}

// purge is useless
func (m *Manager) RequestUninstall(id string, purge bool) *dbus.Error {
	go func() {
		logger.Infof("RequestUninstall id: %q", id)
		err := m.uninstall(id)
		if err != nil {
			logger.Warningf("uninstall %q failed: %v", id, err)
			m.service.Emit(m, "UninstallFailed", id, err.Error())
			return
		}

		m.removeAutostart(id)
		logger.Infof("uninstall %q success", id)
		m.service.Emit(m, "UninstallSuccess", id)
	}()
	return nil
}

func (m *Manager) isItemsChanged() bool {
	old := atomic.SwapUint32(&m.itemsChangedHit, 0)
	return old > 0
}

func (m *Manager) Search(key string) *dbus.Error {
	key = strings.ToLower(key)
	logger.Debug("Search key:", key)

	keyRunes := []rune(key)

	m.searchMu.Lock()

	if m.isItemsChanged() {
		// clear search cache
		m.popPushOpChan <- &popPushOp{popCount: len(m.currentRunes)}
		m.currentRunes = nil
	}

	popCount, runesPush := runeSliceDiff(keyRunes, m.currentRunes)

	logger.Debugf("runeSliceDiff key %v, current %v", keyRunes, m.currentRunes)
	logger.Debugf("runeSliceDiff popCount %v, runesPush %v", popCount, runesPush)

	m.popPushOpChan <- &popPushOp{popCount, runesPush}
	m.currentRunes = keyRunes

	m.searchMu.Unlock()
	return nil
}

func (m *Manager) GetUseProxy(id string) (bool, *dbus.Error) {
	return m.getUseFeature(gsKeyAppsUseProxy, id)
}

func (m *Manager) SetUseProxy(id string, val bool) *dbus.Error {
	return m.setUseFeature(gsKeyAppsUseProxy, id, val)
}

func (m *Manager) GetDisableScaling(id string) (bool, *dbus.Error) {
	return m.getUseFeature(gsKeyAppsDisableScaling, id)
}

func (m *Manager) SetDisableScaling(id string, val bool) *dbus.Error {
	return m.setUseFeature(gsKeyAppsDisableScaling, id, val)
}
