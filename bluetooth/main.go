/**
 * Copyright (c) 2014 Deepin, Inc.
 *               2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package bluetooth

import (
	"dbus/org/bluez"
	. "pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

type Daemon struct {
	*ModuleBase
}

func NewBluetoothDaemon(logger *log.Logger) *Daemon {
	var d = new(Daemon)
	d.ModuleBase = NewModuleBase("bluetooth", d, logger)
	return d
}

func (*Daemon) GetDependencies() []string {
	return []string{}
}

var bluetooth *Bluetooth

func (*Daemon) Start() error {
	if bluetooth != nil {
		return nil
	}

	logger.BeginTracing()

	bluetooth = NewBluetooth()
	err := dbus.InstallOnSession(bluetooth)
	if err != nil {
		// don't panic or fatal here
		logger.Error("register dbus interface failed: ", err)
		bluetooth = nil
		return err
	}

	agent := newAgent()
	agent.b = bluetooth
	err = dbus.InstallOnSystem(agent)
	if err != nil {
		//don't panic or fatal here
		logger.Error("register dbus interface failed: ", err)
		return err
	}

	// initialize bluetooth after dbus interface installed
	bluetooth.initBluetooth()

	agentManager, err := bluez.NewAgentManager1(dbusBluezDest, dbusBluezPath)
	if nil != err {
		logger.Info("get agentmanager failed: ", err)
		return err
	}
	err = agentManager.RegisterAgent(dbusAgentPath, "DisplayYesNo")
	if nil != err {
		logger.Info("register agent failed: ", err)
		return err
	}
	err = agentManager.RequestDefaultAgent(dbusAgentPath)
	if nil != err {
		logger.Info("set defaulet agent failed: ", err)
		return err
	}
	return nil
}

func (*Daemon) Stop() error {
	if bluetooth == nil {
		return nil
	}

	DestroyBluetooth(bluetooth)
	bluetooth = nil
	logger.EndTracing()
	return nil
}
