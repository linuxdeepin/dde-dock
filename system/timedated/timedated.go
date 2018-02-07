/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
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
