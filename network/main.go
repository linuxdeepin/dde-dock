/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package network

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

var (
	logger  = log.NewLogger("daemon/network")
	manager *Manager
)

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
	manager = nil
	logger.EndTracing()
	return nil
}
