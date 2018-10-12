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
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/dde/daemon/network/proxychains"
	"pkg.deepin.io/lib/log"
)

var (
	logger  = log.NewLogger("daemon/network")
	manager *Manager
)

func init() {
	proxychains.SetLogger(logger)
}

type Daemon struct {
	*loader.ModuleBase
}

func NewDaemon(logger *log.Logger) *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("network", daemon, logger)
	return daemon
}

func (d *Daemon) GetDependencies() []string {
	return []string{}
}

func (d *Daemon) Start() error {
	if manager != nil {
		return nil
	}

	initSlices() // initialize slice code
	service := loader.GetService()

	manager = NewManager(service)
	manager.init()

	managerServerObj, err := service.NewServerObject(dbusPath, manager)
	if err != nil {
		return err
	}

	managerServerObj.SetWriteCallback(manager, "NetworkingEnabled", manager.networkingEnabledWriteCb)
	managerServerObj.SetWriteCallback(manager, "VpnEnabled", manager.vpnEnabledWriteCb)

	err = managerServerObj.Export()
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

	go func() {
		manager.initDBusDaemon()
		watchNetworkManagerRestart(manager)
	}()
	return nil
}

func (d *Daemon) Stop() error {
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
	manager.sysSigLoop.Stop()
	service.StopExport(manager)

	if manager.proxyChainsManager != nil {
		service.StopExport(manager.proxyChainsManager)
		manager.proxyChainsManager = nil
	}

	manager = nil
	return nil
}
