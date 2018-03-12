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
	"pkg.deepin.io/lib/log"
)

var (
	_imageBlur *ImageBlur
	logger     = log.NewLogger("daemon/accounts")
)

func init() {
	loader.Register(NewDaemon())
}

type Daemon struct {
	*loader.ModuleBase
	manager        *Manager
	loginedManager *logined.Manager
	imageBlur      *ImageBlur
}

func NewDaemon() *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("accounts", daemon, logger)
	return daemon
}

func (*Daemon) GetDependencies() []string {
	return []string{}
}

func (d *Daemon) Start() error {
	if d.manager != nil {
		return nil
	}

	service := loader.GetService()
	d.manager = NewManager(service)

	err := service.Export(dbusPath, d.manager)
	if err != nil {
		if d.manager.watcher != nil {
			d.manager.watcher.EndWatch()
			d.manager.watcher = nil
		}
		return err
	}

	d.manager.exportUsers()

	d.imageBlur = newImageBlur(service)
	_imageBlur = d.imageBlur
	err = service.Export(imageBlurDBusPath, d.imageBlur)
	if err != nil {
		d.imageBlur = nil
		return err
	}

	d.loginedManager, err = logined.Register(logger, service)
	if err != nil {
		logger.Error("Failed to create logined manager:", err)
		return err
	}
	err = service.Export(logined.DBusPath, d.loginedManager)
	if err != nil {
		logined.Unregister(d.loginedManager)
		d.loginedManager = nil
		return err
	}

	err = service.RequestName(dbusServiceName)
	if err != nil {
		return err
	}
	return nil
}

func (d *Daemon) Stop() error {
	if d.manager != nil {
		d.manager.destroy()
		d.manager = nil
	}

	service := loader.GetService()

	if d.imageBlur != nil {
		service.StopExport(d.imageBlur)
		d.imageBlur = nil
		_imageBlur = nil
	}

	if d.loginedManager != nil {
		service.StopExport(d.loginedManager)
		d.loginedManager = nil
	}

	return nil
}
