/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package network

import (
	dbusmgr "dbus/org/freedesktop/dbus/system"
	"dbus/org/freedesktop/login1"
	nmdbus "dbus/org/freedesktop/networkmanager"
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
	nmManager    *nmdbus.Manager
	nmSettings   *nmdbus.Settings
	loginManager *login1.Manager
	dbusDaemon   *dbusmgr.DBusDaemon
)

func initDbusObjects() {
	var err error
	if nmManager, err = nmdbus.NewManager(dbusNmDest, dbusNmPath); err != nil {
		logger.Error(err)
	}
	if nmSettings, err = nmdbus.NewSettings(dbusNmDest, dbusNmSettingPath); err != nil {
		logger.Error(err)
	}
	if loginManager, err = login1.NewManager(dbusLoginDest, dbusLoginPath); err != nil {
		logger.Error(err)
	}
}
func destroyDbusObjects() {
	// destroy global dbus objects manually when stopping service is
	// required for that there are multiple signal connected with
	// theme which need to be removed
	login1.DestroyManager(loginManager)
	nmdbus.DestroyManager(nmManager)
	nmdbus.DestroySettings(nmSettings)
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
