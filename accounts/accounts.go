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
	"pkg.deepin.io/dde/daemon/accounts/logined"
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

var (
	_m     *Manager
	_login *logined.Manager
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

	_login, err = logined.Register(logger)
	if err != nil {
		logger.Error("Failed to create logined manager:", err)
		return err
	}
	err = dbus.InstallOnSystem(_login)
	if err != nil {
		logined.Unregister(_login)
		_login = nil
		logger.Error("Failed to install logined bus:", err)
		return err
	}

	return nil
}

func (*Daemon) Stop() error {
	if _m != nil {
		_m.destroy()
		_m = nil
	}

	if _login != nil {
		logined.Unregister(_login)
		_login = nil
	}

	return nil
}
