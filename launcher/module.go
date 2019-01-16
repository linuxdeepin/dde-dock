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
	"time"

	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger("daemon/launcher")

func init() {
	loader.Register(NewModule(logger))
}

type Module struct {
	*loader.ModuleBase
	manager *Manager
}

func NewModule(logger *log.Logger) *Module {
	daemon := new(Module)
	daemon.ModuleBase = loader.NewModuleBase("launcher", daemon, logger)
	return daemon
}

func (d *Module) GetDependencies() []string {
	return []string{}
}

func (d *Module) start() error {
	service := loader.GetService()

	var err error
	d.manager, err = NewManager(service)
	if err != nil {
		return err
	}

	err = service.Export(dbusObjPath, d.manager)
	if err != nil {
		return err
	}

	err = service.RequestName(dbusServiceName)
	return err
}

func (d *Module) Start() error {
	go func() {
		t0 := time.Now()
		err := d.start()
		if err != nil {
			logger.Warning(err)
		}
		logger.Info("start launcher module cost", time.Since(t0))
	}()
	return nil
}

func (d *Module) Stop() error {
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

	d.manager.destroy()
	d.manager = nil
	return nil
}
