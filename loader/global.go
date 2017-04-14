/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package loader

import (
	"pkg.deepin.io/lib/log"
	"sync"
)

var loaderInitializer sync.Once

var getLoader = func() func() *Loader {
	var loader *Loader
	return func() *Loader {
		loaderInitializer.Do(func() {
			loader = &Loader{
				modules: Modules{},
				log:     log.NewLogger("daemon/loader"),
			}
		})
		return loader
	}
}()

func Register(m Module) {
	loader := getLoader()
	loader.AddModule(m)
}

func List() []Module {
	return getLoader().List()
}

func GetModule(name string) Module {
	return getLoader().GetModule(name)
}

func SetLogLevel(pri log.Priority) {
	getLoader().SetLogLevel(pri)
}

func EnableModules(enablingModules []string, disableModules []string, flag EnableFlag) error {
	return getLoader().EnableModules(enablingModules, disableModules, flag)
}

func ToggleLogDebug(enabled bool) {
	var priority log.Priority = log.LevelInfo
	if enabled {
		priority = log.LevelDebug
	}
	for _, m := range getLoader().modules {
		m.SetLogLevel(priority)
	}
}

func StartAll() {
	allModules := getLoader().List()
	modules := []string{}
	for _, module := range allModules {
		modules = append(modules, module.Name())
	}
	getLoader().EnableModules(modules, []string{}, EnableFlagNone)
}

// TODO: check dependencies
func StopAll() {
	modules := getLoader().List()
	for _, module := range modules {
		module.Enable(false)
	}
}
