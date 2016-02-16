/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package launcher

import (
	storeApi "dbus/com/deepin/store/api"
	"strings"
	"sync"

	"gir/gio-2.0"
	"pkg.deepin.io/dde/daemon/dstore"
	"pkg.deepin.io/dde/daemon/launcher/category"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
	"pkg.deepin.io/dde/daemon/launcher/item"
	"pkg.deepin.io/dde/daemon/launcher/item/search"
	. "pkg.deepin.io/dde/daemon/launcher/log"
	. "pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/gettext"
	. "pkg.deepin.io/lib/initializer"
	"pkg.deepin.io/lib/log"
)

// Daemon is the module interface's implementation.
type Daemon struct {
	*ModuleBase
	launcher *Launcher
}

// NewLauncherDaemon creates a new launcher daemon module.
func NewLauncherDaemon(logger *log.Logger) *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = NewModuleBase("launcher", daemon, logger)
	return daemon
}

// GetDependencies returns the dependencies of this module.
func (d *Daemon) GetDependencies() []string {
	return []string{}
}

// Stop stops the launcher module.
func (d *Daemon) Stop() error {
	if d.launcher == nil {
		return nil
	}

	d.launcher.destroy()
	d.launcher = nil

	Log.EndTracing()
	return nil
}

func loadItemsInfo(im *item.Manager, cm *category.Manager) {
	appChan := make(chan *gio.AppInfo)
	go func() {
		allApps := gio.AppInfoGetAll()
		for _, app := range allApps {
			appChan <- app
		}
		close(appChan)
	}()

	err := cm.LoadCategoryInfo()
	if err != nil {
		Log.Warning("LoadAppCategoryInfo failed:", err)
	}
	var wg sync.WaitGroup
	const N = 20
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			for app := range appChan {
				if !app.ShouldShow() {
					app.Unref()
					continue
				}

				desktopApp := gio.ToDesktopAppInfo(app)
				newItem := item.New(desktopApp)
				cid, err := cm.QueryID(desktopApp)
				Log.Debug("get category", category.ToString(cid), "for", newItem.ID())
				newItem.SetCategoryID(cid)
				if err != nil {
					Log.Debug("QueryCategoryID failed:", err)
				}

				im.AddItem(newItem)
				cm.AddItem(newItem.ID(), newItem.CategoryID())

				app.Unref()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	cm.FreeAppCategoryInfo()
}
func isZH() bool {
	lang := gettext.QueryLang()
	return strings.HasPrefix(lang, "zh")
}

// Start starts the launcher module.
func (d *Daemon) Start() error {
	if d.launcher != nil {
		return nil
	}

	Log.BeginTracing()

	gettext.InitI18n()

	// DesktopAppInfo.ShouldShow does not know deepin.
	gio.DesktopAppInfoSetDesktopEnv("Deepin")

	err := NewInitializer().Init(func(interface{}) (interface{}, error) {
		store, err := dstore.New()
		storeAdapter := NewDStoreAdapter(store)
		f := func(DStore) {}
		f(storeAdapter)
		return storeAdapter, err
	}).InitOnSessionBus(func(store interface{}) (interface{}, error) {
		d.launcher = NewLauncher()

		im := item.NewManager(store.(DStore), DStoreDesktopPkgMapFile, DStoreInstalledTimeFile)
		cm := category.NewManager(store.(DStore), category.GetAllInfos(DStoreAllCategoryInfoFile))

		d.launcher.setItemManager(im)
		d.launcher.setCategoryManager(cm)

		loadItemsInfo(im, cm)

		storeAPI, err := storeApi.NewDStoreDesktop("com.deepin.store.Api", "/com/deepin/store/Api")
		if err == nil {
			d.launcher.setStoreAPI(storeAPI)
		}

		if isZH() {
			Log.Info("enable pinyin search in zh env")
			names := []string{}
			for _, item := range im.GetAllItems() {
				names = append(names, item.LocaleName())
			}

			pinyinObj, err := search.NewPinYinSearchAdapter(names)
			if err != nil {
				Log.Warning("CreatePinYinSearch failed:", err)
			}
			d.launcher.setPinYinObject(pinyinObj)
		}

		d.launcher.listenItemChanged()

		return d.launcher, nil
	}).InitOnSessionBus(func(interface{}) (interface{}, error) {
		coreSetting := gio.NewSettings("com.deepin.dde.launcher")
		setting, err := NewSettings(coreSetting)
		d.launcher.setSetting(setting)
		return setting, err
	}).GetError()

	if err != nil {
		d.Stop()
	}

	return err
}
