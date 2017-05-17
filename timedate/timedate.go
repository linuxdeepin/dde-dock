/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package timedate

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

var (
	_manager *Manager

	logger = log.NewLogger("daemon/timedate")
)

type Daemon struct {
	*loader.ModuleBase
}

func NewDaemon(logger *log.Logger) *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("timedate", daemon, logger)
	return daemon
}

func (d *Daemon) GetDependencies() []string {
	return []string{}
}

// Start to run timedate manager
func (d *Daemon) Start() error {
	if _manager != nil {
		return nil
	}

	logger.BeginTracing()

	var err error
	_manager, err = NewManager()
	if err != nil {
		logger.Error("Create Manager failed:", err)
		logger.EndTracing()
		return err
	}

	err = dbus.InstallOnSession(_manager)
	if err != nil {
		logger.Error("Install DBus failed:", err)
		_manager.destroy()
		_manager = nil
		logger.EndTracing()
		return err
	}
	go func() {
		_manager.init()
		_manager.handlePropChanged()
	}()
	return nil
}

// Stop the timedate manager
func (d *Daemon) Stop() error {
	if _manager == nil {
		return nil
	}

	logger.EndTracing()
	_manager.destroy()
	_manager = nil
	return nil
}
