/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
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

func (d *Daemon) Start() error {
	if d.manager != nil {
		return nil
	}
	var err error

	service := loader.GetService()

	d.manager, err = newManager(service)
	if err != nil {
		return err
	}

	err = service.Export(dbusPath, d.manager)
	if err != nil {
		d.manager.destroy()
		d.manager = nil
		return err
	}

	err = service.RequestName(dbusServiceName)
	if err != nil {
		d.manager.destroy()
		d.manager = nil
		return err
	}

	go func() {
		m := d.manager
		m.init()

		m.initHandlers()

		// listen gsettings changed event
		m.listenGSettingsChanged(gsSchemaSystem, d.manager.gsSystem, shortcuts.ShortcutTypeSystem)
		m.listenGSettingsChanged(gsSchemaMediaKey, d.manager.gsMediaKey, shortcuts.ShortcutTypeMedia)
		m.listenGSettingsChanged(gsSchemaGnomeWM, d.manager.gsGnomeWM, shortcuts.ShortcutTypeWM)

		m.eliminateKeystrokeConflict()
		m.shortcutManager.EventLoop()
	}()

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
