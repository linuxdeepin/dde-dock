/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package keybinding

import (
	"github.com/BurntSushi/xgbutil/xevent"
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
		daemon.manager.shortcuts.ListenXEvents()

		// listen gsetting changed event
		daemon.manager.listenGSettingsChanged(daemon.manager.sysSetting, shortcuts.ShortcutTypeSystem)
		daemon.manager.listenGSettingsChanged(daemon.manager.mediaSetting, shortcuts.ShortcutTypeMedia)
		daemon.manager.listenGSettingsChanged(daemon.manager.wmSetting, shortcuts.ShortcutTypeWM)

		xevent.Main(daemon.manager.xu)
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
