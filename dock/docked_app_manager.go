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
	"container/list"
	"gir/gio-2.0"
	"gir/glib-2.0"
	"io/ioutil"
	"os"
	"path/filepath"
	"pkg.deepin.io/lib/dbus"
	dutils "pkg.deepin.io/lib/utils"
	"text/template"
)

const (
	dockSchema           string = "com.deepin.dde.dock"
	settingKeyDockedApps string = "docked-apps"
	dockedItemTemplate   string = `[Desktop Entry]
Name={{ .Name }}
Exec={{ .Exec }}
Icon={{ .Icon }}
Type=Application
Terminal=false
StartupNotify=false
`
)

var scratchDir string = filepath.Join(os.Getenv("HOME"), ".config/dock/scratch")

// DockedAppManager是管理已驻留程序的管理器。
type DockedAppManager struct {
	settings *gio.Settings
	items    *list.List

	// Docked是信号，在某程序驻留成功后被触发，并将该程序的id发送给信号的接受者。
	Docked func(id string) // find indicator on front-end.
	// Undocked是信号，在某已驻留程序被移除驻留后被触发，将被移除程序id发送给信号接受者。
	Undocked func(id string)
}

func NewDockedAppManager() *DockedAppManager {
	m := &DockedAppManager{}
	m.init()
	return m
}

func (m *DockedAppManager) init() {
	m.items = list.New()
	m.settings = gio.NewSettings(dockSchema)
	if m.settings == nil {
		return
	}
	m.handleOldConfigFile()

	// TODO:
	// listen changed.
	appList := m.settings.GetStrv(settingKeyDockedApps)
	for _, id := range appList {
		m.items.PushBack(normalizeAppID(id))
	}
}

func (m *DockedAppManager) destroy() {
	if m.settings != nil {
		m.settings.Unref()
	}
	dbus.UnInstallObject(m)
}

func (m *DockedAppManager) handleOldConfigFile() {
	conf := glib.NewKeyFile()
	defer conf.Free()

	confFile := filepath.Join(glib.GetUserConfigDir(), "dock/apps.ini")
	_, err := conf.LoadFromFile(confFile, glib.KeyFileFlagsNone)
	if err != nil {
		logger.Debug("Open old dock config file failed:", err)
		return
	}

	inited, err := conf.GetBoolean("__Config__", "inited")
	if err == nil && inited {
		return
	}

	_, ids, err := conf.GetStringList("__Config__", "Position")
	if err != nil {
		logger.Debug("Read docked app from old config file failed:", err)
		return
	}
	for _, id := range ids {
		if a := NewDesktopAppInfo(id + ".desktop"); a != nil {
			a.Unref()
			continue
		}

		exec, _ := conf.GetString(id, "CmdLine")
		icon, _ := conf.GetString(id, "Icon")
		title, _ := conf.GetString(id, "Name")
		createScratchDesktopFile(id, title, icon, exec)
	}

	m.saveAppList(ids)
	conf.SetBoolean("__Config__", "inited", true)

	_, content, err := conf.ToData()
	if err != nil {
		return
	}

	var mode os.FileMode = 0666
	stat, err := os.Lstat(confFile)
	if err == nil {
		mode = stat.Mode()
	}

	err = ioutil.WriteFile(confFile, []byte(content), mode)
	if err != nil {
		logger.Warning("Save Config file failed:", err)
	}
}

// DockedAppList返回程序id列表。
func (m *DockedAppManager) DockedAppList() []string {
	if m.settings != nil {
		appList := m.settings.GetStrv(settingKeyDockedApps)
		return appList
	}
	return nil
}

// IsDocked通过传入的程序id判断一个程序是否已经驻留。
func (m *DockedAppManager) IsDocked(id string) bool {
	_, item := m.fuzzyFindItem(id)
	return item != nil
}

func (m *DockedAppManager) fuzzyFindItem(id string) (string, *list.Element) {
	// return new id and element
	id = normalizeAppID(id)
	if item := m.findItem(id); item != nil {
		return id, item
	}
	// item is nil
	guessId := trimDesktop(normalizeAppID(guess_desktop_id(id)))
	if guessId == id {
		return id, nil
	}
	return guessId, m.findItem(id)
}

type dockedItemInfo struct {
	Name, Icon, Exec string
}

// 废弃，请使用新接口RequestDock
func (m *DockedAppManager) Dock(id, title, icon, cmd string) bool {
	logger.Info("Try dock", id)
	newId, item := m.fuzzyFindItem(id)
	if item != nil {
		logger.Debugf("App %q is already docked.", newId)
		return false
	}
	if newId == "" {
		// item is nil and newId is empty
		// create scratch desktop file
		if err := createScratchDesktopFile(id, title, icon, cmd); err != nil {
			return false
		}
	} else {
		id = newId
	}

	m.items.PushBack(id)
	m.saveAppList(m.getAppList())
	m.emitSignal("Docked", id)
	app := ENTRY_MANAGER.runtimeApps[id]
	if app != nil {
		app.buildMenu()
	}
	ENTRY_MANAGER.createNormalApp(id)
	return true
}

// RequestDock驻留程序。通常情况下只需要传递程序id即可，在特殊情况下需要传入title，icon以及cmd。
// title表示前端程序的tooltip内容，icon为程序图标，cmd为程序的启动命令。
// 成功后会触发Docked信号。
func (m *DockedAppManager) RequestDock(id, title, icon, cmd string) bool {
	return m.Dock(id, title, icon, cmd)
}

// TODO: 删除此函数，因为拼写错误
func (m *DockedAppManager) ReqeustDock(id, title, icon, cmd string) bool {
	return m.Dock(id, title, icon, cmd)
}

// 废弃，请使用新接口RequestUndock
func (m *DockedAppManager) Undock(id string) bool {
	logger.Info("Try Undock", id)
	id, removeItem := m.fuzzyFindItem(id)
	if removeItem == nil {
		logger.Debug("Not find docked app:", id)
		return false
	}
	m.items.Remove(removeItem)
	m.saveAppList(m.getAppList())
	removedApp := removeItem.Value.(string)
	m.emitSignal("Undocked", removedApp)
	removeScratchFiles(id)
	app := ENTRY_MANAGER.runtimeApps[id]
	if app != nil {
		// update menu item undock to dock
		app.buildMenu()
	}
	ENTRY_MANAGER.destroyNormalApp(id)
	return true
}

// RequestUndock 通过程序id移除已驻留程序。成功后会触发Undocked信号。
func (m *DockedAppManager) RequestUndock(id string) bool {
	return m.Undock(id)
}

func (m *DockedAppManager) emitSignal(name string, values ...interface{}) {
	logger.Debugf("Emit Signal %v %v", name, values)
	dbus.Emit(m, name, values...)
}

func (m *DockedAppManager) findItem(id string) *list.Element {
	// logger.Debugf("findItem %q in %v", id, m.getAppList()) // just for debug
	if id == "" {
		return nil
	}
	for e := m.items.Front(); e != nil; e = e.Next() {
		if e.Value.(string) == id {
			return e
		}
	}
	return nil
}

// Sort将已驻留的程序按传入的程序id的顺序重新排序，并保存。
func (m *DockedAppManager) Sort(items []string) {
	logger.Debug("Before sort:", items)
	for _, item := range items {
		item = normalizeAppID(item)
		if i := m.findItem(item); i != nil {
			m.items.PushBack(m.items.Remove(i))
		}
	}

	sortedAppList := m.getAppList()
	logger.Debug("After sort:", sortedAppList)
	m.saveAppList(sortedAppList)
}

func (m *DockedAppManager) getAppList() []string {
	length := m.items.Len()
	list := make([]string, 0, length)
	for e := m.items.Front(); e != nil; e = e.Next() {
		list = append(list, e.Value.(string))
	}
	return list
}

// 保存 apps 到 gsettings
func (m *DockedAppManager) saveAppList(apps []string) {
	m.settings.SetStrv(settingKeyDockedApps, apps)
	gio.SettingsSync()
}

func createScratchDesktopFile(id, title, icon, cmd string) error {
	logger.Debugf("create scratch file for %q", id)
	err := os.MkdirAll(scratchDir, 0775)
	if err != nil {
		logger.Warning("create scratch directory failed:", err)
		return err
	}
	f, err := os.OpenFile(filepath.Join(scratchDir, id+".desktop"),
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0744)
	if err != nil {
		logger.Warning("Open file for write failed:", err)
		return err
	}

	defer f.Close()
	temp := template.Must(template.New("docked_item_temp").Parse(dockedItemTemplate))
	dockedItem := dockedItemInfo{title, icon, cmd}
	logger.Debugf("dockedItem: %#v", dockedItem)
	err = temp.Execute(f, dockedItem)
	if err != nil {
		return err
	}
	return nil
}

func removeScratchFiles(id string) {
	extList := []string{"desktop", "sh", "png"}
	for _, ext := range extList {
		file := filepath.Join(scratchDir, id+"."+ext)
		if dutils.IsFileExist(file) {
			logger.Debugf("remove scratch file %q", file)
			err := os.Remove(file)
			if err != nil {
				logger.Warning("remove scratch file %q failed:", file, err)
			}
		}
	}
}
