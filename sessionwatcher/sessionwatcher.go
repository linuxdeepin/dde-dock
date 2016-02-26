/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package sessionwatcher

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
)

var (
	logger   = log.NewLogger("daemon/sessionwatcher")
	_manager *Manager
)

type Daemon struct {
	*loader.ModuleBase
}

func init() {
	loader.Register(NewDaemon(logger))
}

func NewDaemon(logger *log.Logger) *Daemon {
	var d = new(Daemon)
	d.ModuleBase = loader.NewModuleBase("sessionwatcher", d, logger)
	return d
}

func (*Daemon) GetDependencies() []string {
	return []string{}
}

func (*Daemon) Start() error {
	if _manager != nil {
		return nil
	}

	_manager = newManager()
	_manager.AddTask(newDockTask())
	_manager.AddTask(newDesktopTask())
	go _manager.StartLoop()
	return nil
}

func (*Daemon) Stop() error {
	if _manager == nil {
		return nil
	}

	_manager.QuitLoop()
	_manager = nil
	logger.EndTracing()
	return nil
}
