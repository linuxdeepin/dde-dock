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
	"pkg.deepin.io/dde/daemon/common/dsync"
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
)

func init() {
	loader.Register(newModule(logger))
}

type Module struct {
	sSaver     *ScreenSaver
	syncConfig *dsync.Config
	*loader.ModuleBase
}

func newModule(logger *log.Logger) *Module {
	m := new(Module)
	m.ModuleBase = loader.NewModuleBase("screensaver", m, logger)
	return m
}

func (m *Module) GetDependencies() []string {
	return []string{}
}

func (m *Module) Start() error {
	service := loader.GetService()

	has, err := service.NameHasOwner(dbusServiceName)
	if err != nil {
		return err
	}
	if has {
		logger.Warning("ScreenSaver has been register, exit...")
		return nil
	}

	if m.sSaver != nil {
		return nil
	}

	m.sSaver, err = newScreenSaver(service)
	if err != nil {
		return err
	}

	err = service.Export(dbusPath, m.sSaver)
	if err != nil {
		return err
	}

	err = service.RequestName(dbusServiceName)
	if err != nil {
		return err
	}

	m.syncConfig = dsync.NewConfig("screensaver", &syncConfig{}, m.sSaver.sigLoop, dScreenSaverPath, logger)
	err = service.Export(dScreenSaverPath, m.syncConfig)
	if err != nil {
		return err
	}

	err = service.RequestName(dScreenSaverServiceName)
	if err != nil {
		return err
	}

	err = m.syncConfig.Register()
	if err != nil {
		logger.Warning("failed to register for deepin sync:", err)
	}

	return nil
}

func (m *Module) Stop() error {
	if m.sSaver == nil {
		return nil
	}

	service := loader.GetService()
	err := service.ReleaseName(dbusServiceName)
	if err != nil {
		logger.Warning(err)
	}

	err = service.StopExport(m.sSaver)
	if err != nil {
		logger.Warning(err)
	}
	m.sSaver.destroy()
	m.sSaver = nil

	err = service.ReleaseName(dScreenSaverServiceName)
	if err != nil {
		logger.Warning(err)
	}

	err = service.StopExport(m.syncConfig)
	if err != nil {
		logger.Warning(err)
	}
	m.syncConfig.Destroy()
	if m.sSaver.xConn != nil {
		m.sSaver.xConn.Close()
	}
	return nil
}
