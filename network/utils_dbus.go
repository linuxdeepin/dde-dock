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

package network

import (
	dbusmgr "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.dbus"
	"github.com/linuxdeepin/go-dbus-factory/org.freedesktop.login1"
	nmdbus "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.networkmanager"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/proxy"
)

var (
	nmManager    *nmdbus.Manager
	nmSettings   *nmdbus.Settings
	loginManager *login1.Manager
	dbusDaemon   *dbusmgr.DBus // system dbus daemon
)

func (m *Manager) initDbusObjects() {
	systemBus, err := dbus.SystemBus()
	if err != nil {
		logger.Error(err)
		return
	}

	nmManager = nmdbus.NewManager(systemBus)
	nmManager.InitSignalExt(m.sysSigLoop, true)

	nmSettings = nmdbus.NewSettings(systemBus)
	nmSettings.InitSignalExt(m.sysSigLoop, true)

	loginManager = login1.NewManager(systemBus)
	loginManager.InitSignalExt(m.sysSigLoop, true)

}

var sysSigLoop *dbusutil.SignalLoop

func initSysSignalLoop() {
	systemBus, err := dbus.SystemBus()
	if err != nil {
		logger.Error(err)
		return
	}
	sysSigLoop = dbusutil.NewSignalLoop(systemBus, 50)
	sysSigLoop.Start()
}

func initDBusDaemon() {
	systemBus, err := dbus.SystemBus()
	if err != nil {
		logger.Error(err)
		return
	}
	dbusDaemon = dbusmgr.NewDBus(systemBus)
	dbusDaemon.InitSignalExt(sysSigLoop, true)
}

func destroyDBusDaemon() {
	dbusDaemon.RemoveHandler(proxy.RemoveAllHandlers)
}

func destroyDbusObjects() {
	// destroy global dbus objects manually when stopping service is
	// required for that there are multiple signal connected with
	// theme which need to be removed
	nmManager.RemoveHandler(proxy.RemoveAllHandlers)
	nmSettings.RemoveHandler(proxy.RemoveAllHandlers)
	loginManager.RemoveHandler(proxy.RemoveAllHandlers)
}
