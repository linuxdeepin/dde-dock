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

package sessionwatcher

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
)

var (
	logger   = log.NewLogger("daemon/sessionwatcher")
	_manager *Manager
)

type Daemon struct {
	*loader.ModuleBase
}

func init() {
	loader.Register(NewDaemon(logger))
}

func NewDaemon(logger *log.Logger) *Daemon {
	var d = new(Daemon)
	d.ModuleBase = loader.NewModuleBase("sessionwatcher", d, logger)
	return d
}

func (*Daemon) GetDependencies() []string {
	return []string{}
}

func (*Daemon) Start() error {
	if _manager != nil {
		return nil
	}
	service := loader.GetService()

	var err error
	_manager, err = newManager(service)
	if err != nil {
		return err
	}

	_manager.initUserSessions()

	err = service.Export(dbusPath, _manager)
	if err != nil {
		return err
	}

	err = service.RequestName(dbusServiceName)
	if err != nil {
		return err
	}

	return nil
}

func (*Daemon) Stop() error {
	if _manager == nil {
		return nil
	}

	service := loader.GetService()
	service.StopExport(_manager)
	_manager.destroy()
	_manager = nil
	logger.EndTracing()
	return nil
}
