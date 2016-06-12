/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package mounts

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

var (
	_manager *Manager
	logger   = log.NewLogger("daemon/mounts")
)

func init() {
	loader.Register(NewDaemon(logger))
}

type Daemon struct {
	*loader.ModuleBase
}

func NewDaemon(logger *log.Logger) *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("mounts", daemon, logger)
	return daemon
}

func (d *Daemon) GetDependencies() []string {
	return []string{}
}

func (d *Daemon) Start() error {
	if _manager != nil {
		return nil
	}

	_manager = newManager()
	logger.BeginTracing()
	err := dbus.InstallOnSession(_manager)
	if err != nil {
		logger.Error("Install mounts dbus failed:", err)
		_manager.destroy()
		_manager = nil
		return err
	}
	go _manager.updateDiskInfo()
	_manager.handleEvent()
	return nil
}

func (d *Daemon) Stop() error {
	if _manager == nil {
		return nil
	}

	_manager.destroy()
	dbus.UnInstallObject(_manager)
	_manager = nil
	return nil
}
