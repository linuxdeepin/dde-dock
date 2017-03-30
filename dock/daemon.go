/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dock

import (
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

type Daemon struct {
	*loader.ModuleBase
}

func NewDaemon(logger *log.Logger) *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("dock", daemon, logger)
	return daemon
}

func (d *Daemon) Stop() error {
	if dockManager != nil {
		dockManager.destroy()
		dockManager = nil
	}

	if XU != nil {
		XU.Conn().Close()
		XU = nil
	}

	logger.EndTracing()
	return nil
}

func (d *Daemon) startFailed(args ...interface{}) {
	logger.Error(args...)
	d.Stop()
}

func (d *Daemon) Start() error {
	if dockManager != nil {
		return nil
	}
	logger.BeginTracing()

	var err error
	// init x conn
	XU, err = xgbutil.NewConn()
	if err != nil {
		d.startFailed(err)
		return err
	}

	initAtom()
	initDir()
	initPathDirCodeMap()

	dockManager, err = NewDockManager()
	if err != nil {
		d.startFailed(err)
		return err
	}

	dockManager.listenRootWindowPropertyChange()
	go xevent.Main(XU)

	dbus.Emit(dockManager, "ServiceRestarted")
	return nil
}

func (d *Daemon) GetDependencies() []string {
	return []string{}
}

func (d *Daemon) Name() string {
	return "dock"
}
