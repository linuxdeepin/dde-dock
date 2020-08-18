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
	"os"

	dbus "github.com/godbus/dbus"
	x "github.com/linuxdeepin/go-x11-client"
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/log"
)

type Daemon struct {
	*loader.ModuleBase
	manager *TrayManager
	snw     *StatusNotifierWatcher
	sigLoop *dbusutil.SignalLoop // session bus signal loop
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

	sessionBus, err := dbus.SessionBus()
	if err != nil {
		return err
	}

	d.sigLoop = dbusutil.NewSignalLoop(sessionBus, 10)
	d.sigLoop.Start()

	err = service.Export(dbusPath, d.manager)
	if err != nil {
		return err
	}

	err = d.manager.sendClientMsgMANAGER()
	if err != nil {
		return err
	}

	err = service.RequestName(dbusServiceName)
	if err != nil {
		return err
	}

	err = service.Emit(d.manager, "Inited")
	if err != nil {
		return err
	}

	if os.Getenv("DDE_DISABLE_STATUS_NOTIFIER_WATCHER") != "1" {
		d.snw = newStatusNotifierWatcher(service, d.sigLoop)
		d.snw.listenDBusNameOwnerChanged()
		err = service.Export(snwDBusPath, d.snw)
		if err != nil {
			return err
		}
		err = service.RequestName(snwDBusServiceName)
		if err != nil {
			logger.Warning("failed to request name:", err)
			return nil
		}
	} else {
		logger.Info("disable status notifier watcher")
	}
	return nil
}

func (d *Daemon) Stop() error {
	if XConn != nil {
		XConn.Close()
	}
	return nil
}
