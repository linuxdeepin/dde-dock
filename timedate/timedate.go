/**
 * Copyright (c) 2011 ~ 2015 Deepin, Inc.
 *               2013 ~ 2015 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package timedate

import (
	"pkg.linuxdeepin.com/dde-daemon/loader"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/log"
)

var (
	_manager *Manager

	logger = log.NewLogger("dde-daemon/timedate")
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
	_manager.handlePropChanged()
	return nil
}

func (d *Daemon) Stop() error {
	if _manager == nil {
		return nil
	}

	logger.EndTracing()
	_manager.destroy()
	_manager = nil
	return nil
}
