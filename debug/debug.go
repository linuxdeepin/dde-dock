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

package debug

import (
	_ "net/http/pprof"

	"net/http"
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
)

var (
	logger = log.NewLogger("daemon/keybinding")
)

type Daemon struct {
	pprofExists bool
	*loader.ModuleBase
}

func NewDaemon() *Daemon {
	var d = new(Daemon)
	d.ModuleBase = loader.NewModuleBase("debug", d, logger)
	return d
}

func (*Daemon) GetDependencies() []string {
	return []string{}
}

func (d *Daemon) Start() error {
	if !d.pprofExists {
		go func() {
			d.pprofExists = true
			err := http.ListenAndServe("localhost:6969", nil)
			if err != nil {
				d.pprofExists = false
				logger.Error("Enable pprof failed:", err)
				return
			}
		}()
	}
	if d.LogLevel() != log.LevelDebug {
		loader.ToggleLogDebug(true)
	}

	return nil
}

func (d *Daemon) Stop() error {
	if d.LogLevel() == log.LevelDebug {
		loader.ToggleLogDebug(false)
	}
	return nil
}
