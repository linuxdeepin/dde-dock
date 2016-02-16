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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gir/gio-2.0"
	"gir/glib-2.0"
	"pkg.deepin.io/lib/dbus"
)

const (
	SchemaId       string = "com.deepin.dde.dock"
	DockedApps     string = "docked-apps"
	DockedItemTemp string = `[Desktop Entry]
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
	core  *gio.Settings
	items *list.List

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
	m.core = gio.NewSettings(SchemaId)
	if m.core == nil {
		return
	}

	// TODO:
	// listen changed.
	appList := m.core.GetStrv(DockedApps)
	for _, id := range appList {
		m.items.PushBack(normalizeAppID(id))
	}

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
		createScratchFile(id, title, icon, exec)
	}

	m.core.SetStrv(DockedApps, ids)
	gio.SettingsSync()
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

// DockedAppList返回一个已排序的程序id列表。
func (m *DockedAppManager) DockedAppList() []string {
	if m.core != nil {
		appList := m.core.GetStrv(DockedApps)
		return appList
	}
	return make([]string, 0)
}

// IsDocked通过传入的程序id判断一个程序是否已经驻留。
func (m *DockedAppManager) IsDocked(id string) bool {
	id = normalizeAppID(id)
	item := m.findItem(id)
	if item != nil {
		return true
	}

	if id = trimDesktop(guess_desktop_id(id)); id != "" {
		item = m.findItem(id)
	}
	// logger.Info("IsDocked:", item, item != nil)
	return item != nil
}

type dockedItemInfo struct {
	Name, Icon, Exec string
}

// Dock驻留程序。通常情况下只需要传递程序id即可，在特殊情况下需要传入title，icon以及cmd。
// title表示前端程序的tooltip内容，icon为程序图标，cmd为程序的启动命令。
// 成功后会触发Docked信号。
// （废弃，此接口名并不好，第一反映很难理解，请使用新接口RequestDock)
func (m *DockedAppManager) Dock(id, title, icon, cmd string) bool {
	id = normalizeAppID(id)
	logger.Info("start dock", id)
	idElement := m.findItem(id)
	if idElement != nil {
		logger.Info(id, "is already docked.")
		return false
	}

	id = strings.ToLower(id)
	idElement = m.findItem(id)
	if idElement != nil {
		logger.Info(id, "is already docked.")
		return false
	}

	desktopID := guess_desktop_id(id)
	if desktopID == "" {
		if e := createScratchFile(id, title, icon, cmd); e != nil {
			return false
		}
	} else {
		id = normalizeAppID(trimDesktop(desktopID))
	}
	m.items.PushBack(id)
	dbus.Emit(m, "Docked", id)
	app := ENTRY_MANAGER.runtimeApps[id]
	if app != nil {
		app.buildMenu()
	}

	if _, ok := ENTRY_MANAGER.normalApps[id]; ok {
		logger.Info(id, "is already docked")
		return true
	}
	ENTRY_MANAGER.createNormalApp(id)

	return true
}

// RequestDock驻留程序。通常情况下只需要传递程序id即可，在特殊情况下需要传入title，icon以及cmd。
// title表示前端程序的tooltip内容，icon为程序图标，cmd为程序的启动命令。
// 成功后会触发Docked信号。
func (m *DockedAppManager) ReqeustDock(id, title, icon, cmd string) bool {
	return m.Dock(id, title, icon, cmd)
}

func (m *DockedAppManager) doUndock(id string) bool {
	logger.Info("doUndock", id)
	removeItem := m.findItem(id)
	if removeItem == nil {
		logger.Warning("not find docked app:", id)
		return false
	}

	logger.Info("Undock", id, ", Remove", m.items.Remove(removeItem))
	m.core.SetStrv(DockedApps, m.toSlice())
	gio.SettingsSync()
	os.Remove(filepath.Join(scratchDir, id+".desktop"))
	os.Remove(filepath.Join(scratchDir, id+".sh"))
	os.Remove(filepath.Join(scratchDir, id+".png"))
	dbus.Emit(m, "Undocked", removeItem.Value.(string))
	app := ENTRY_MANAGER.runtimeApps[id]
	if app != nil {
		app.buildMenu()
	}

	if app, ok := ENTRY_MANAGER.normalApps[id]; ok {
		logger.Info("destroy normal app")
		ENTRY_MANAGER.destroyNormalApp(app)
	}

	return true
}

// Undock通过程序id移除已驻留程序。成功后会触发Undocked信号。（废弃，请使用新接口RequestUndock）
func (m *DockedAppManager) Undock(id string) bool {
	id = normalizeAppID(id)
	if m.doUndock(id) {
		return true
	}

	tmpId := ""
	if tmpId = trimDesktop(guess_desktop_id(id)); tmpId != "" {
		logger.Debug("undock guess desktop id:", tmpId)
		m.doUndock(tmpId)
		return true
	}

	tmpId = normalizeAppID(id)
	if m.doUndock(tmpId) {
		logger.Debug("undock replace - to _:", tmpId)
		return true
	}

	return false
}

// RequestUndock益处指定程序id的已驻留程序。成功后会出发Undocked信号。
func (m *DockedAppManager) RequestUndock(id string) bool {
	return m.Undock(id)
}

func (m *DockedAppManager) findItem(id string) *list.Element {
	lowerID := strings.ToLower(id)
	for e := m.items.Front(); e != nil; e = e.Next() {
		if strings.ToLower(e.Value.(string)) == lowerID {
			return e
		}
	}
	return nil
}

// Sort将已驻留的程序按传入的程序id的顺序重新排序，并保存。
func (m *DockedAppManager) Sort(items []string) {
	logger.Debug("sort:", items)
	for _, item := range items {
		item = normalizeAppID(item)
		if i := m.findItem(item); i != nil {
			m.items.PushBack(m.items.Remove(i))
		}
	}
	l := m.toSlice()
	logger.Debug("sorted:", l)
	m.core.SetStrv(DockedApps, l)
	gio.SettingsSync()
}

func (m *DockedAppManager) toSlice() []string {
	appList := make([]string, 0)
	for e := m.items.Front(); e != nil; e = e.Next() {
		appList = append(appList, e.Value.(string))
	}
	return appList
}

func createScratchFile(id, title, icon, cmd string) error {
	logger.Info("create scratch file for %s with cmd %q and title %q", id, cmd, title)

	homeDir := os.Getenv("HOME")
	path := ".config/dock/scratch"
	configDir := filepath.Join(homeDir, path)
	err := os.MkdirAll(configDir, 0775)
	if err != nil {
		logger.Warning("create scratch failed:", err)
		return err
	}
	f, err := os.OpenFile(filepath.Join(configDir, id+".desktop"),
		os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0744)
	if err != nil {
		logger.Warning("OpenScratch to write failed:", err)
		return err
	}
	defer f.Close()
	temp := template.Must(template.New("docked_item_temp").Parse(DockedItemTemp))
	e := temp.Execute(f, dockedItemInfo{title, icon, cmd})
	if e != nil {
		return e
	}
	return nil
}
