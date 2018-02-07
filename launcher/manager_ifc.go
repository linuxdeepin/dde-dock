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
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/keyfile"
)

const (
	dbusDest           = "com.deepin.dde.daemon.Launcher"
	dbusObjPath        = "/com/deepin/dde/daemon/Launcher"
	dbusIFC            = dbusDest
	desktopMainSection = "Desktop Entry"
)

var errorInvalidID = errors.New("Invalid ID")

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusObjPath,
		Interface:  dbusIFC,
	}
}

func (m *Manager) GetAllItemInfos() []ItemInfo {
	list := make([]ItemInfo, 0, len(m.items))
	for _, item := range m.items {
		list = append(list, item.newItemInfo())
	}
	logger.Debug("GetAllItemInfos list length:", len(list))
	return list
}

func (m *Manager) GetItemInfo(id string) (ItemInfo, error) {
	item := m.getItemById(id)
	if item == nil {
		return ItemInfo{}, errorInvalidID
	}
	return item.newItemInfo(), nil
}

func (m *Manager) GetAllNewInstalledApps() ([]string, error) {
	newApps, err := m.launchedRecorder.GetNew()
	if err != nil {
		return nil, err
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

func (m *Manager) IsItemOnDesktop(id string) (bool, error) {
	item := m.getItemById(id)
	if item == nil {
		return false, errorInvalidID
	}
	file := appInDesktop(m.getAppIdByFilePath(item.Path))
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			// not exist
			return false, nil
		} else {
			return false, err
		}
	} else {
		// exist
		return true, nil
	}
}

func (m *Manager) RequestRemoveFromDesktop(id string) (bool, error) {
	item := m.getItemById(id)
	if item == nil {
		return false, errorInvalidID
	}
	file := appInDesktop(m.getAppIdByFilePath(item.Path))
	err := os.Remove(file)
	return err == nil, err
}

func (m *Manager) RequestSendToDesktop(id string) (bool, error) {
	item := m.getItemById(id)
	if item == nil {
		return false, errorInvalidID
	}
	dest := appInDesktop(m.getAppIdByFilePath(item.Path))
	_, err := os.Stat(dest)
	if err != nil {
		if !os.IsNotExist(err) {
			return false, err
		}
		// dest file not exist
	} else {
		// dest file exist
		return false, os.ErrExist
	}

	kf := keyfile.NewKeyFile()
	if err := kf.LoadFromFile(item.Path); err != nil {
		logger.Warning(err)
		return false, err
	}
	kf.SetString(desktopMainSection, "X-Deepin-CreatedBy", dbusDest)
	kf.SetString(desktopMainSection, "X-Deepin-AppID", id)
	// Desktop files in user desktop direcotry do not require executable permission
	if err := kf.SaveToFile(dest); err != nil {
		logger.Warning("save new desktop file failed:", err)
		return false, err
	}
	// success
	go soundutils.PlaySystemSound(soundutils.EventIconToDesktop, "")
	return true, nil
}

// MarkLaunched 废弃
func (m *Manager) MarkLaunched(id string) error {
	return nil
}

// purge is useless
func (m *Manager) RequestUninstall(id string, purge bool) {
	go func() {
		logger.Infof("RequestUninstall id: %q", id)
		err := m.uninstall(id)
		if err != nil {
			logger.Warningf("uninstall %q failed: %v", id, err)
			dbus.Emit(m, "UninstallFailed", id, err.Error())
			return
		}

		m.removeAutostart(id)
		logger.Infof("uninstall %q success", id)
		dbus.Emit(m, "UninstallSuccess", id)
	}()
}

func (m *Manager) isItemsChanged() bool {
	old := atomic.SwapUint32(&m.itemsChangedHit, 0)
	return old > 0
}

func (m *Manager) Search(key string) {
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
}

func (m *Manager) GetUseProxy(id string) (bool, error) {
	return m.getUseFeature(gsKeyAppsUseProxy, id)
}

func (m *Manager) SetUseProxy(id string, val bool) error {
	return m.setUseFeature(gsKeyAppsUseProxy, id, val)
}

func (m *Manager) GetDisableScaling(id string) (bool, error) {
	return m.getUseFeature(gsKeyAppsDisableScaling, id)
}

func (m *Manager) SetDisableScaling(id string, val bool) error {
	return m.setUseFeature(gsKeyAppsDisableScaling, id, val)
}
