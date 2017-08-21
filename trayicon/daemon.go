/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package trayicon

import (
	x "github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/util/wm/ewmh"
	"github.com/linuxdeepin/go-x11-client/util/wm/icccm"
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
)

type Daemon struct {
	*loader.ModuleBase
	manager *TrayManager
}

func NewDaemon(logger *log.Logger) *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("trayicon", daemon, logger)
	return daemon
}

func (d *Daemon) GetDependencies() []string {
	return []string{}
}

func (d *Daemon) Name() string {
	return "trayicon"
}

func (d *Daemon) Start() error {
	logger.BeginTracing()

	var err error
	// init x conn
	XConn, err = x.NewConn()
	if err != nil {
		d.startFailed(err)
		return err
	}

	ewmhConn, err = ewmh.NewConn(XConn)
	if err != nil {
		d.startFailed(err)
		return err
	}

	icccmConn, err = icccm.NewConn(XConn)
	if err != nil {
		d.startFailed(err)
		return err
	}

	initX()
	d.manager = NewTrayManager()

	return nil
}

func (d *Daemon) Stop() error {
	if XConn != nil {
		XConn.Close()
		XConn = nil
	}

	if d.manager != nil {
		d.manager.destroy()
		d.manager = nil
	}

	logger.EndTracing()
	return nil
}

func (d *Daemon) startFailed(args ...interface{}) {
	logger.Error(args...)
	d.Stop()
}
