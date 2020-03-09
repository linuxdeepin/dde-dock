/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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
	"pkg.deepin.io/lib/log"
)

const (
	dbusServiceName = "com.deepin.api.XEventMonitor"
	dbusPath        = "/com/deepin/api/XEventMonitor"
	dbusInterface   = dbusServiceName
	moduleName      = "x-event-monitor"
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
	service := loader.GetService()

	m, err := newManager(service)
	if err != nil {
		return err
	}
	m.initXExtensions()
	go m.handleXEvent()

	err = service.Export(dbusPath, m)
	if err != nil {
		return err
	}

	err = service.RequestName(dbusServiceName)
	if err != nil {
		return err
	}

	service.Emit(m, "CancelAllArea")

	return nil
}

func (d *Daemon) Stop() error {
	return nil
}
