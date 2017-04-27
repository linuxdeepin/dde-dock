/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package power

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger("daemon/session/power")

func init() {
	loader.Register(NewDaemon(logger))
}

type Daemon struct {
	*loader.ModuleBase
	manager *Manager
}

func NewDaemon(logger *log.Logger) *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("power", daemon, logger)
	return daemon
}

func (d *Daemon) GetDependencies() []string {
	return []string{"screensaver", "sessionwatcher"}
}

func (d *Daemon) Start() (err error) {
	logger.BeginTracing()
	d.manager, err = NewManager()
	if err != nil {
		logger.Error(err)
		logger.EndTracing()
		return
	}
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
