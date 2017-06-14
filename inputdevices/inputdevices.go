/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package inputdevices

import (
	"fmt"
	. "pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

var (
	_manager *Manager
	logger   = log.NewLogger("daemon/inputdevices")
)

type Daemon struct {
	*ModuleBase
}

func init() {
	Register(NewInputdevicesDaemon(logger))
}
func NewInputdevicesDaemon(logger *log.Logger) *Daemon {
	var d = new(Daemon)
	d.ModuleBase = NewModuleBase("inputdevices", d, logger)
	return d
}

func (*Daemon) GetDependencies() []string {
	return []string{}
}

func (*Daemon) Start() error {
	if _manager != nil {
		return nil
	}

	logger.BeginTracing()
	_manager = NewManager()

	err := installSessionBus(_manager)
	if err != nil {
		return err
	}

	err = installSessionBus(_manager.kbd)
	if err != nil {
		return err
	}

	err = installSessionBus(_manager.wacom)
	if err != nil {
		return err
	}

	err = installSessionBus(_manager.tpad)
	if err != nil {
		return err
	}
	err = installSessionBus(_manager.mouse)
	if err != nil {
		return err
	}
	err = installSessionBus(_manager.trackPoint)
	if err != nil {
		return err
	}

	go func() {
		_manager.init()
		startDeviceListener()
	}()
	return nil
}

func (*Daemon) Stop() error {
	if _manager == nil {
		return nil
	}
	_manager = nil

	logger.EndTracing()
	getWacom().destroy()
	// TODO endDeviceListener will be stuck
	endDeviceListener()
	return nil
}

func installSessionBus(obj dbus.DBusObject) error {
	if obj == nil {
		logger.Error("Invalid dbus object: empty")
		return fmt.Errorf("Invalid dbus object")
	}

	err := dbus.InstallOnSession(obj)
	if err != nil {
		logger.Error("Install session bus failed:", err)
		return err
	}
	return nil
}
