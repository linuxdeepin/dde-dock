/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package keybinding

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

type Daemon struct {
	*loader.ModuleBase
}

var (
	_m     *Manager
	logger = log.NewLogger("daemon/keybinding")
)

func NewDaemon(logger *log.Logger) *Daemon {
	var d = new(Daemon)
	d.ModuleBase = loader.NewModuleBase("keybinding", d, logger)
	return d
}

// Check 'touchpad' whether exist
func (*Daemon) GetDependencies() []string {
	return []string{}
}

func (*Daemon) Start() error {
	if _m != nil {
		return nil
	}
	logger.BeginTracing()
	var err error
	_m, err = NewManager()
	if err != nil {
		logger.EndTracing()
		return err
	}

	err = dbus.InstallOnSession(_m)
	if err != nil {
		_m.destroy()
		_m = nil
		logger.EndTracing()
		return err
	}

	dbus.InstallOnSession(_m.media)

	_m.initGrabedList()
	_m.listenGSettingChanged()
	go _m.startLoop()
	return nil
}

func (*Daemon) Stop() error {
	if _m == nil {
		return nil
	}

	logger.EndTracing()
	_m.destroy()
	return nil
}
