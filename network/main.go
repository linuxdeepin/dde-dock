/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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
	"pkg.deepin.io/lib/dbus"
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

	logger.BeginTracing()

	initSlices() // initialize slice code

	manager = NewManager()
	err := dbus.InstallOnSession(manager)
	if err != nil {
		logger.Error("register dbus interface failed: ", err)
		manager = nil
		return err
	}

	manager.proxyChainsManager = proxychains.NewManager()
	err = dbus.InstallOnSession(manager.proxyChainsManager)
	if err != nil {
		logger.Warning("register proxychains manager dbus interface failed: ", err)
		manager.proxyChainsManager = nil
		return err
	}

	// initialize manager after dbus installed
	go func() {
		manager.initManager()
		initDbusDaemon()
		watchNetworkManagerRestart(manager)
	}()
	return nil
}

func (d *Daemon) Stop() error {
	if manager == nil {
		return nil
	}

	destroyDbusDaemon()
	DestroyManager(manager)
	dbus.UnInstallObject(manager)

	if manager.proxyChainsManager != nil {
		dbus.UnInstallObject(manager.proxyChainsManager)
		manager.proxyChainsManager = nil
	}

	manager = nil
	logger.EndTracing()
	return nil
}
