/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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

package x_event_monitor

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

const (
	dbusDest      = "com.deepin.api.XEventMonitor"
	dbusObjPath   = "/com/deepin/api/XEventMonitor"
	dbusInterface = dbusDest
	moduleName    = "x_event_monitor"
)

var (
	logger = log.NewLogger(moduleName)
)

func init() {
	loader.Register(NewDaemon(logger))
}

type Daemon struct {
	*loader.ModuleBase
}

func NewDaemon(logger *log.Logger) *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase(moduleName, daemon, logger)
	return daemon
}

func (d *Daemon) GetDependencies() []string {
	return []string{}
}

func (d *Daemon) Name() string {
	return moduleName
}

func (d *Daemon) Start() error {
	logger.BeginTracing()
	m := newManager()
	rawEventCallback = m.handleRawEvent
	go startListen()

	err := dbus.InstallOnSession(m)
	if err != nil {
		return err
	}

	dbus.DealWithUnhandledMessage()
	dbus.Emit(m, "CancelAllArea")
	return nil
}

func (d *Daemon) Stop() error {
	logger.EndTracing()
	return nil
}
