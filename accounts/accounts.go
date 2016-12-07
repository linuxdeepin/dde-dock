/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package accounts

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

var (
	_m     *Manager
	logger = log.NewLogger("daemon/accounts")
)

type Daemon struct {
	*loader.ModuleBase
}

func NewDaemon() *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("accounts", daemon, logger)
	return daemon
}

func (*Daemon) GetDependencies() []string {
	return []string{}
}

func (*Daemon) Start() error {
	if _m != nil {
		return nil
	}

	logger.BeginTracing()
	_m = NewManager()
	err := dbus.InstallOnSystem(_m)
	if err != nil {
		logger.Error("Install manager dbus failed:", err)
		if _m.watcher != nil {
			_m.watcher.EndWatch()
			_m.watcher = nil
		}
		return err
	}

	_m.installUsers()
	return nil
}

func (*Daemon) Stop() error {
	if _m == nil {
		return nil
	}

	_m.destroy()
	_m = nil

	return nil
}
