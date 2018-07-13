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

package screensaver

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
)

func init() {
	loader.Register(NewDaemon(logger))
}

type Daemon struct {
	*loader.ModuleBase
}

func NewDaemon(logger *log.Logger) *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("screensaver", daemon, logger)
	return daemon
}

func (d *Daemon) GetDependencies() []string {
	return []string{}
}

func (d *Daemon) Start() error {
	service := loader.GetService()

	has, err := service.NameHasOwner(dbusServiceName)
	if err != nil {
		return err
	}
	if has {
		logger.Warning("ScreenSaver has been register, exit...")
		return nil
	}

	if _ssaver != nil {
		return nil
	}

	_ssaver, err = newScreenSaver(service)
	if err != nil {
		return err
	}

	err = service.Export(dbusPath, _ssaver)
	if err != nil {
		return err
	}

	err = service.RequestName(dbusServiceName)
	if err != nil {
		_ssaver.destroy()
		_ssaver = nil
		return err
	}

	return nil
}

func (d *Daemon) Stop() error {
	if _ssaver == nil {
		return nil
	}

	service := loader.GetService()
	err := service.ReleaseName(dbusServiceName)
	if err != nil {
		logger.Warning(err)
	}

	err = service.StopExport(_ssaver)
	if err != nil {
		logger.Warning(err)
	}
	_ssaver.destroy()
	_ssaver = nil
	return nil
}
