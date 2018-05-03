/*
 * Copyright (C) 2013 ~ 2018 Deepin Technology Co., Ltd.
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

package main

import (
	"fmt"
	"sync"

	"gir/gio-2.0"
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/gsettings"

	// sort modules
	_ "pkg.deepin.io/dde/daemon/network"

	_ "pkg.deepin.io/dde/daemon/audio"

	_ "pkg.deepin.io/dde/daemon/screensaver"
	_ "pkg.deepin.io/dde/daemon/sessionwatcher"

	// depends: screensaver, sessionwatcher
	_ "pkg.deepin.io/dde/daemon/session/power"

	_ "pkg.deepin.io/dde/daemon/launcher"
	_ "pkg.deepin.io/dde/daemon/service_trigger"

	_ "pkg.deepin.io/dde/daemon/clipboard"
	_ "pkg.deepin.io/dde/daemon/keybinding"

	_ "pkg.deepin.io/dde/daemon/appearance"
	_ "pkg.deepin.io/dde/daemon/inputdevices"

	_ "pkg.deepin.io/dde/daemon/gesture"
	_ "pkg.deepin.io/dde/daemon/housekeeping"
	_ "pkg.deepin.io/dde/daemon/timedate"

	_ "pkg.deepin.io/dde/daemon/bluetooth"
	_ "pkg.deepin.io/dde/daemon/screenedge"

	_ "pkg.deepin.io/dde/daemon/fprintd"
	_ "pkg.deepin.io/dde/daemon/mime"

	// depends: network
	_ "pkg.deepin.io/dde/daemon/miracast"
	_ "pkg.deepin.io/dde/daemon/systeminfo"

	_ "pkg.deepin.io/dde/daemon/calltrace"
	_ "pkg.deepin.io/dde/daemon/debug"
)

var (
	moduleLocker   sync.Mutex
	daemonSchema   = "com.deepin.dde.daemon"
	daemonSettings = gio.NewSettings(daemonSchema)
)

func listenDaemonSettings() {
	gsettings.ConnectChanged(daemonSchema, "*", func(name string) {
		// gsettings key names must keep consistent with module names
		moduleLocker.Lock()
		defer moduleLocker.Unlock()
		module := loader.GetModule(name)
		if module == nil {
			logger.Error("Invalid module name:", name)
			return
		}

		enable := daemonSettings.GetBoolean(name)
		err := checkDependencies(daemonSettings, module, enable)
		if err != nil {
			logger.Error(err)
			return
		}

		err = module.Enable(enable)
		if err != nil {
			logger.Warningf("Enable '%s' failed: %v", name, err)
			return
		}
	})
}

func checkDependencies(s *gio.Settings, module loader.Module, enabled bool) error {
	if enabled {
		depends := module.GetDependencies()
		for _, n := range depends {
			if s.GetBoolean(n) != true {
				return fmt.Errorf("Dependency lose: %v", n)
			}
		}
		return nil
	}

	for _, m := range loader.List() {
		if m == nil || m.Name() == module.Name() {
			continue
		}

		if m.IsEnable() == true && isStrInList(module.Name(), m.GetDependencies()) {
			return fmt.Errorf("Can not diable this module '%s', because of it was depended by'%s'",
				module.Name(), m.Name())
		}
	}
	return nil
}

func isStrInList(item string, list []string) bool {
	for _, v := range list {
		if item == v {
			return true
		}
	}
	return false
}
