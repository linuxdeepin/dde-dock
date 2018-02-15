/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

package loader

import (
	"sync"

	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/log"
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

func SetService(s *dbusutil.Service) {
	l := getLoader()
	l.service = s
}

func GetService() *dbusutil.Service {
	return getLoader().service
}

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
