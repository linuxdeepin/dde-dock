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

package power

import (
	"pkg.deepin.io/dde/daemon/loader"
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
	service := loader.GetService()
	d.manager, err = newManager(service)
	if err != nil {
		return
	}

	err = service.Export(dbusPath, d.manager,
		d.manager.warnLevelConfig, d.manager.syncConfig)
	if err != nil {
		return err
	}

	err = service.RequestName(dbusServiceName)
	if err != nil {
		return err
	}

	go d.manager.init()
	return
}

func (d *Daemon) Stop() error {
	if d.manager == nil {
		return nil
	}
	service := loader.GetService()
	err := service.ReleaseName(dbusServiceName)
	if err != nil {
		logger.Warning(err)
	}
	err = service.StopExport(d.manager)
	if err != nil {
		logger.Warning(err)
	}
	err = service.StopExport(d.manager.syncConfig)
	if err != nil {
		logger.Warning(err)
	}

	d.manager.destroy()
	d.manager = nil
	return nil
}
