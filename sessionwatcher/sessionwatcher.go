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
	"pkg.deepin.io/lib/dbus"
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

	m, err := newManager()
	if err != nil {
		return err
	}
	_manager = m
	_manager.initUserSessions()

	err = dbus.InstallOnSession(_manager)
	if err != nil {
		return err
	}

	return nil
}

func (*Daemon) Stop() error {
	if _manager == nil {
		return nil
	}

	dbus.UnInstallObject(_manager)
	_manager.destroy()
	_manager = nil
	logger.EndTracing()
	return nil
}
