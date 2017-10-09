/*
 * Copyright (C) 2016 ~ 2017 Deepin Technology Co., Ltd.
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

package keybinding

import (
	"pkg.deepin.io/dde/daemon/keybinding/shortcuts"
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

func init() {
	loader.Register(NewDaemon(logger))
	shortcuts.SetLogger(logger)
}

type Daemon struct {
	*loader.ModuleBase
	manager *Manager
}

var (
	logger = log.NewLogger("daemon/keybinding")
)

func NewDaemon(logger *log.Logger) *Daemon {
	var d = new(Daemon)
	d.ModuleBase = loader.NewModuleBase("keybinding", d, logger)
	return d
}

func (*Daemon) GetDependencies() []string {
	return []string{}
}

func (daemon *Daemon) Start() error {
	if daemon.manager != nil {
		return nil
	}
	logger.BeginTracing()
	var err error

	daemon.manager, err = NewManager()
	if err != nil {
		logger.EndTracing()
		return err
	}

	err = dbus.InstallOnSession(daemon.manager)
	if err != nil {
		daemon.manager.destroy()
		daemon.manager = nil
		logger.EndTracing()
		return err
	}

	go func() {
		daemon.manager.init()

		daemon.manager.initHandlers()

		// listen gsetting changed event
		daemon.manager.listenGSettingsChanged(daemon.manager.gsSystem, shortcuts.ShortcutTypeSystem)
		daemon.manager.listenGSettingsChanged(daemon.manager.gsMediaKey, shortcuts.ShortcutTypeMedia)
		daemon.manager.listenGSettingsChanged(daemon.manager.gsGnomeWM, shortcuts.ShortcutTypeWM)

		daemon.manager.shortcutManager.EventLoop()
	}()

	return nil
}

func (daemon *Daemon) Stop() error {
	if daemon.manager == nil {
		return nil
	}
	logger.EndTracing()
	daemon.manager.destroy()
	dbus.UnInstallObject(daemon.manager)
	daemon.manager = nil
	return nil
}
