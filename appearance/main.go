/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
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

package appearance

import (
	. "pkg.linuxdeepin.com/dde-daemon/loader"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/log"
)

type Daemon struct {
	*ModuleBase
}

func NewAppearanceDaemon(logger *log.Logger) *Daemon {
	var d = new(Daemon)
	d.ModuleBase = NewModuleBase("appearance", d, logger)
	return d
}

func (*Daemon) GetDependencies() []string {
	return []string{}
}

var _manager *Manager

func finalize() {
	logger.EndTracing()
	_manager.destroy()
	_manager = nil
}

func (*Daemon) Start() error {
	if _manager != nil {
		return nil
	}

	logger.BeginTracing()
	_manager = NewManager()
	err := dbus.InstallOnSession(_manager)
	if err != nil {
		logger.Error(err)
		finalize()
		return err
	}
	return nil
}

func (*Daemon) Stop() error {
	if _manager == nil {
		return nil
	}

	finalize()
	return nil
}
