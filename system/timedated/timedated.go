/**
 * Copyright (C) 2017 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package timedated

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

type Daemon struct {
	*loader.ModuleBase
}

var (
	logger   = log.NewLogger("timedated")
	_manager *Manager
)

func NewDaemon(logger *log.Logger) *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("timedated", daemon, logger)
	return daemon
}

func (*Daemon) GetDependencies() []string {
	return []string{}
}

func init() {
	loader.Register(NewDaemon(logger))
}

func (*Daemon) Start() error {
	if _manager != nil {
		return nil
	}

	var err error
	_manager, err = NewManager()
	if err != nil {
		logger.Error("Failed to new timedated manager:", err)
		return err
	}

	err = dbus.InstallOnSystem(_manager)
	if err != nil {
		logger.Error("Failed to install system dbus:", err)
		return err
	}
	return nil
}

func (*Daemon) Stop() error {
	if _manager != nil {
		_manager.destroy()
		_manager = nil
	}

	return nil
}
