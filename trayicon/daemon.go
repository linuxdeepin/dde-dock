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

package trayicon

import (
	x "github.com/linuxdeepin/go-x11-client"
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
)

type Daemon struct {
	*loader.ModuleBase
	manager *TrayManager
}

const moduleName = "trayicon"

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
	var err error
	// init x conn
	XConn, err = x.NewConn()
	if err != nil {
		return err
	}

	initX()
	service := loader.GetService()
	d.manager = NewTrayManager(service)

	err = service.Export(dbusPath, d.manager)
	if err != nil {
		return err
	}

	d.manager.sendClientMsgMANAGER()

	err = service.RequestName(dbusServiceName)
	if err != nil {
		return err
	}
	service.Emit(d.manager, "Inited")

	return nil
}

func (d *Daemon) Stop() error {
	return nil
}
