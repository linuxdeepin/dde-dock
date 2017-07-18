/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package apps

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger("daemon/apps")

func init() {
	loader.Register(NewDaemon(logger))
}

type Daemon struct {
	*loader.ModuleBase
	recorder *ALRecorder
	watcher  *DFWatcher
}

func NewDaemon(logger *log.Logger) *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("apps", daemon, logger)
	return daemon
}

func (d *Daemon) Start() error {
	logger.Debug("apps daemon start")
	if watcher, err := NewDFWachter(); err != nil {
		return err
	} else {
		d.watcher = watcher
	}

	d.recorder = NewALRecorder(d.watcher)

	// install recorder and watcher
	if err := dbus.InstallOnSystem(d.recorder); err != nil {
		return err
	}
	if err := dbus.InstallOnSystem(d.watcher); err != nil {
		return err
	}

	d.recorder.emitServiceRestarted()

	return nil
}

func (d *Daemon) Stop() error {
	// TODO
	return nil
}

func (d *Daemon) GetDependencies() []string {
	return []string{}
}

func (d *Daemon) Name() string {
	return "apps"
}
