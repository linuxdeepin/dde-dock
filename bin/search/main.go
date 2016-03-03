/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package main

import (
	"os"
	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
	"time"
)

type Manager struct {
	writeStart bool
	writeEnd   chan bool
}

const (
	DBUS_DEST = "com.deepin.daemon.Search"
	DBUS_PATH = "/com/deepin/daemon/Search"
	DBUS_IFC  = "com.deepin.daemon.Search"
)

var (
	logger = log.NewLogger("daemon/search")
)

func newManager() *Manager {
	m := Manager{}

	m.writeStart = false

	return &m
}

var _manager *Manager

func GetManager() *Manager {
	if _manager == nil {
		_manager = newManager()
	}

	return _manager
}

func main() {
	if !lib.UniqueOnSession(DBUS_DEST) {
		logger.Warning("There is an Search running")
		return
	}

	logger.BeginTracing()
	defer logger.EndTracing()
	logger.SetRestartCommand("/usr/lib/deepin-daemon/search")

	if err := dbus.InstallOnSession(GetManager()); err != nil {
		logger.Fatal("Search Install DBus Failed:", err)
		return
	}
	dbus.DealWithUnhandledMessage()

	dbus.SetAutoDestroyHandler(time.Second*5, func() bool {
		if GetManager().writeStart {
			select {
			case <-GetManager().writeEnd:
				return true
			}
		}

		return true
	})

	if err := dbus.Wait(); err != nil {
		logger.Warning("Search lost dbus:", err)
		os.Exit(-1)
	} else {
		os.Exit(0)
	}
}
