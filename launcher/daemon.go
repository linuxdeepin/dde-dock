/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package launcher

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger("daemon/launcher")

func init() {
	loader.Register(NewDaemon(logger))
}

type Daemon struct {
	*loader.ModuleBase
	manager *Manager
}

func NewDaemon(logger *log.Logger) *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("launcher", daemon, logger)
	return daemon
}

func (d *Daemon) GetDependencies() []string {
	return []string{}
}

func (d *Daemon) Start() (err error) {
	logger.BeginTracing()
	d.manager, err = NewManager()
	if err != nil {
		logger.Error("Failed to new manager:", err)
		logger.EndTracing()
		return
	}
	err = dbus.InstallOnSession(d.manager)
	if err != nil {
		logger.Error("Failed to install dbus:", err)
		logger.EndTracing()
		d.manager.destroy()
		return
	}
	go d.manager.init()
	return
}

func (d *Daemon) Stop() error {
	if d.manager == nil {
		return nil
	}
	d.manager.destroy()
	d.manager = nil
	logger.EndTracing()
	return nil
}
