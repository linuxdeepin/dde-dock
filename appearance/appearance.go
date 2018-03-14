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

// Manage desktop appearance
package appearance

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
)

var (
	_m     *Manager
	logger = log.NewLogger("daemon/appearance")
)

type Daemon struct {
	*loader.ModuleBase
}

func init() {
	loader.Register(NewAppearanceDaemon(logger))
}

func NewAppearanceDaemon(logger *log.Logger) *Daemon {
	var d = new(Daemon)
	d.ModuleBase = loader.NewModuleBase("appearance", d, logger)
	return d
}

func (*Daemon) GetDependencies() []string {
	return []string{}
}

func (*Daemon) Start() error {
	if _m != nil {
		return nil
	}

	service := loader.GetService()

	_m = newManager(service)
	err := service.Export(dbusPath, _m)
	if err != nil {
		_m.destroy()
		return err
	}

	err = service.RequestName(dbusServiceName)
	if err != nil {
		_m.destroy()
		service.StopExport(_m)
		return err
	}

	_m.init()

	go _m.listenCursorChanged()
	go _m.handleThemeChanged()
	_m.listenGSettingChanged()

	return nil
}

func (*Daemon) Stop() error {
	if _m == nil {
		return nil
	}

	_m.destroy()
	service := loader.GetService()
	service.StopExport(_m)
	_m = nil
	return nil
}
