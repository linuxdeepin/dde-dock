/*
 * Copyright (C) 2017 ~ 2017 Deepin Technology Co., Ltd.
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

package mime

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger("daemon/mime")

func init() {
	loader.Register(NewDaemon(logger))
}

type Daemon struct {
	*loader.ModuleBase
	manager *Manager
}

func NewDaemon(logger *log.Logger) *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("mime", daemon, logger)
	return daemon
}

func (d *Daemon) GetDependencies() []string {
	return []string{}
}

func (d *Daemon) Start() error {
	logger.BeginTracing()
	d.manager = NewManager()

	err := dbus.InstallOnSession(d.manager)
	if err != nil {
		logger.Warning("Install Manager dbus failed:", err)
		return err
	}

	media, err := NewMedia()
	if err != nil {
		logger.Error("New Media failed:", err)
		return err
	}
	d.manager.media = media

	err = dbus.InstallOnSession(media)
	if err != nil {
		logger.Warning("Install Media dbus failed:", err)
		return err
	}

	d.manager.initConfigData()
	return nil
}

func (d *Daemon) Stop() error {
	if d.manager == nil {
		return nil
	}

	d.manager.destroy()
	d.manager = nil
	return nil
}
