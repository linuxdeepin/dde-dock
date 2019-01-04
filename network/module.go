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

package network

import (
	"time"

	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/dde/daemon/network/proxychains"
	"pkg.deepin.io/lib/log"
	libnotify "pkg.deepin.io/lib/notify"
)

var (
	logger  = log.NewLogger("daemon/network")
	manager *Manager
)

func init() {
	loader.Register(newModule(logger))
	proxychains.SetLogger(logger)
}

type Module struct {
	*loader.ModuleBase
}

func newModule(logger *log.Logger) *Module {
	module := new(Module)
	module.ModuleBase = loader.NewModuleBase("network", module, logger)
	return module
}

func (d *Module) GetDependencies() []string {
	return []string{}
}

func (d *Module) start() error {
	service := loader.GetService()
	manager = NewManager(service)
	manager.init()

	managerServerObj, err := service.NewServerObject(dbusPath, manager)
	if err != nil {
		return err
	}

	err = managerServerObj.SetWriteCallback(manager, "NetworkingEnabled", manager.networkingEnabledWriteCb)
	if err != nil {
		return err
	}
	err = managerServerObj.SetWriteCallback(manager, "VpnEnabled", manager.vpnEnabledWriteCb)
	if err != nil {
		return err
	}

	// err = managerServerObj.Export()
	err = service.Export(dbusPath, manager, manager.syncConfig)
	if err != nil {
		logger.Error("failed to export manager:", err)
		manager = nil
		return err
	}

	manager.proxyChainsManager = proxychains.NewManager(service)
	err = service.Export(proxychains.DBusPath, manager.proxyChainsManager)
	if err != nil {
		logger.Warning("failed to export proxyChainsManager:", err)
		manager.proxyChainsManager = nil
		return err
	}

	err = service.RequestName(dbusServiceName)
	if err != nil {
		return err
	}

	initDBusDaemon()
	watchNetworkManagerRestart(manager)
	return nil
}

func (d *Module) Start() error {
	libnotify.Init("dde-session-daemon")
	if manager != nil {
		return nil
	}

	initSlices() // initialize slice code
	initSysSignalLoop()
	go func() {
		t0 := time.Now()
		err := d.start()
		if err != nil {
			logger.Warning(err)
		}
		logger.Info("start network module cost", time.Since(t0))
	}()

	return nil
}

func (d *Module) Stop() error {
	if manager == nil {
		return nil
	}

	service := loader.GetService()

	err := service.ReleaseName(dbusServiceName)
	if err != nil {
		logger.Warning(err)
	}

	manager.destroy()
	destroyDBusDaemon()
	sysSigLoop.Stop()
	err = service.StopExport(manager)
	if err != nil {
		logger.Warning(err)
	}

	if manager.proxyChainsManager != nil {
		err = service.StopExport(manager.proxyChainsManager)
		if err != nil {
			logger.Warning(err)
		}
		manager.proxyChainsManager = nil
	}

	manager = nil
	return nil
}
