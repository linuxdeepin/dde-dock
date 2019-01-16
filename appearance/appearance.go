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
	"time"

	"pkg.deepin.io/dde/daemon/appearance/background"
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
)

var (
	_m     *Manager
	logger = log.NewLogger("daemon/appearance")
)

type Module struct {
	*loader.ModuleBase
}

func init() {
	background.SetLogger(logger)
	loader.Register(NewModule(logger))
}

func NewModule(logger *log.Logger) *Module {
	var d = new(Module)
	d.ModuleBase = loader.NewModuleBase("appearance", d, logger)
	return d
}

func (*Module) GetDependencies() []string {
	return []string{}
}

func (*Module) start() error {
	service := loader.GetService()

	_m = newManager(service)
	_m.init()
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

	go _m.listenCursorChanged()
	go _m.handleThemeChanged()
	_m.listenGSettingChanged()
	return nil
}

func (m *Module) Start() error {
	if _m != nil {
		return nil
	}

	go func() {
		t0 := time.Now()
		err := m.start()
		if err != nil {
			logger.Warning(err)
		}
		logger.Info("start appearance module cost", time.Since(t0))
	}()
	return nil
}

func (*Module) Stop() error {
	if _m == nil {
		return nil
	}

	_m.destroy()
	service := loader.GetService()
	service.StopExport(_m)
	_m = nil
	return nil
}
