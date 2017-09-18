/*
 * Copyright (C) 2016 ~ 2017 Deepin Technology Co., Ltd.
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

package power

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger("daemon/session/power")

func init() {
	loader.Register(NewDaemon(logger))
}

type Daemon struct {
	*loader.ModuleBase
	manager *Manager
}

func NewDaemon(logger *log.Logger) *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("power", daemon, logger)
	return daemon
}

func (d *Daemon) GetDependencies() []string {
	return []string{"screensaver", "sessionwatcher"}
}

func (d *Daemon) Start() (err error) {
	logger.BeginTracing()
	d.manager, err = NewManager()
	if err != nil {
		logger.Error(err)
		logger.EndTracing()
		return
	}
	err = dbus.InstallOnSession(d.manager)
	if err != nil {
		d.manager.destroy()
		logger.Error("Failed to install dbus:", err)
		logger.EndTracing()
		return err
	}
	logger.Info("InstallOnSession done")
	go d.manager.init()
	return
}

func (d *Daemon) Stop() error {
	if d.manager == nil {
		return nil
	}
	d.manager.destroy()
	d.manager = nil
	logger.EndTracing()
	return nil
}
