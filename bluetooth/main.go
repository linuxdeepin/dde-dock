/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

package bluetooth

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/log"
)

type daemon struct {
	*loader.ModuleBase
}

func newBluetoothDaemon(logger *log.Logger) *daemon {
	var d = new(daemon)
	d.ModuleBase = loader.NewModuleBase("bluetooth", d, logger)
	return d
}

func (*daemon) GetDependencies() []string {
	return []string{}
}

var bluetooth *Bluetooth
var _agent *agent

func initBluetooth() error {
	destroyBluetooth()

	service := loader.GetService()
	bluetooth = newBluetooth(service)

	err := service.Export(dbusPath, bluetooth)
	if err != nil {
		logger.Warning("failed to export bluetooth:", err)
		bluetooth = nil
		return err
	}

	err = service.RequestName(dbusServiceName)
	if err != nil {
		return err
	}

	sysService, err := dbusutil.NewSystemService()
	if err != nil {
		return err
	}

	_agent = newAgent(sysService)
	_agent.b = bluetooth
	bluetooth.agent = _agent

	err = sysService.Export(agentDBusPath, _agent)
	if err != nil {
		logger.Warning("failed to export agent:", err)
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
