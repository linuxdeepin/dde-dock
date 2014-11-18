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

package network

import (
	dbusmgr "dbus/org/freedesktop/dbus/system"
	"dbus/org/freedesktop/login1"
	"dbus/org/freedesktop/modemmanager1"
	nm "dbus/org/freedesktop/networkmanager"
	"dbus/org/freedesktop/notifications"
	"pkg.linuxdeepin.com/lib/dbus"
)

const (
	dbusMmDest        = "org.freedesktop.ModemManager1"
	dbusNmDest        = "org.freedesktop.NetworkManager"
	dbusNmPath        = "/org/freedesktop/NetworkManager"
	dbusNmSettingPath = "/org/freedesktop/NetworkManager/Settings"
	dbusLoginDest     = "org.freedesktop.login1"
	dbusLoginPath     = "/org/freedesktop/login1"
)

var (
	nmManager    *nm.Manager
	nmSettings   *nm.Settings
	loginManager *login1.Manager
	dbusDaemon   *dbusmgr.DBusDaemon
)

func initDbusObjects() {
	var err error
	if nmManager, err = nm.NewManager(dbusNmDest, dbusNmPath); err != nil {
		logger.Error(err)
	}
	if nmSettings, err = nm.NewSettings(dbusNmDest, dbusNmSettingPath); err != nil {
		logger.Error(err)
	}
	if loginManager, err = login1.NewManager(dbusLoginDest, dbusLoginPath); err != nil {
		logger.Error(err)
	}
	if notifier, err = notifications.NewNotifier(dbusNotifyDest, dbusNotifyPath); err != nil {
		logger.Error(err)
	}
}
func destroyDbusObjects() {
	// destroy global dbus objects manually when stopping service is
	// required for that there are multiple signal connected with
	// theme which need to be removed
	login1.DestroyManager(loginManager)
	nm.DestroyManager(nmManager)
	nm.DestroySettings(nmSettings)
	dbusmgr.DestroyDBusDaemon(dbusDaemon)
	notifications.DestroyNotifier(notifier)
}

func initDbusDaemon() {
	var err error
	if dbusDaemon, err = dbusmgr.NewDBusDaemon("org.freedesktop.DBus", "/org/freedesktop/DBus"); err != nil {
		logger.Error(err)
	}
}
func destroyDbusDaemon() {
	dbusmgr.DestroyDBusDaemon(dbusDaemon)
}
func mmNewModem(modemPath dbus.ObjectPath) (modem *modemmanager1.Modem, err error) {
	modem, err = modemmanager1.NewModem(dbusMmDest, modemPath)
	if err != nil {
		logger.Error(err)
	}
	return
}

func mmGetModemDeviceIdentifier(modemPath dbus.ObjectPath) (devId string, err error) {
	modem, err := mmNewModem(modemPath)
	if err != nil {
		return
	}
	devId = modem.DeviceIdentifier.Get()
	return
}

func mmGetModemDeviceSysPath(modemPath dbus.ObjectPath) (sysPath string, err error) {
	modem, err := mmNewModem(modemPath)
	if err != nil {
		return
	}
	sysPath = modem.Device.Get()
	return
}
