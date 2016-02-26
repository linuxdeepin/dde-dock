/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package bluetooth

import (
	. "pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

type daemon struct {
	*ModuleBase
}

func newBluetoothDaemon(logger *log.Logger) *daemon {
	var d = new(daemon)
	d.ModuleBase = NewModuleBase("bluetooth", d, logger)
	return d
}

func (*daemon) GetDependencies() []string {
	return []string{}
}

var bluetooth *Bluetooth
var _agent *agent

func initBluetooth() error {
	destroyBluetooth()

	bluetooth = newBluetooth()
	err := dbus.InstallOnSession(bluetooth)
	if err != nil {
		// don't panic or fatal here
		logger.Error("register dbus interface failed: ", err)
		bluetooth = nil
		return err
	}

	_agent = newAgent()
	_agent.b = bluetooth
	bluetooth.agent = _agent
	err = dbus.InstallOnSystem(_agent)
	if err != nil {
		//don't panic or fatal here
		logger.Error("register dbus interface failed: ", err)
		return err
	}

	// initialize bluetooth after dbus interface installed
	bluetooth.init()
	_agent.init()
	return nil
}

func destroyBluetooth() {
	if bluetooth != nil {
		bluetooth.destroy()
		bluetooth = nil
	}

	if _agent != nil {
		_agent.destroy()
		_agent = nil
	}

}

func doStart() {
	initBluetooth()
	bluezWatchRestart()
}

func (*daemon) Start() error {
	if bluetooth != nil {
		return nil
	}

	logger.BeginTracing()

	go doStart()
	return nil
}

func (*daemon) Stop() error {
	logger.EndTracing()
	destroyBluetooth()
	bluezDestroyDbusDaemon(bluezDBusDaemon)
	return nil
}
