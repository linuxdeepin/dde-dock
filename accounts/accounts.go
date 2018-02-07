/*
 * Copyright (C) 2013 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package accounts

import (
	"pkg.deepin.io/dde/daemon/accounts/logined"
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

var (
	_m         *Manager
	_login     *logined.Manager
	_imageBlur *ImageBlur
	logger     = log.NewLogger("daemon/accounts")
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

	_imageBlur = newImageBlur()
	err = dbus.InstallOnSystem(_imageBlur)
	if err != nil {
		logger.Warning("failed to install ImageBlur on system DBus:", err)
		_imageBlur = nil
		return err
	}

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

	if _imageBlur != nil {
		dbus.UnInstallObject(_imageBlur)
		_imageBlur = nil
	}

	if _login != nil {
		logined.Unregister(_login)
		_login = nil
	}

	return nil
}
