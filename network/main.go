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
	"dlib/dbus"
	liblogger "dlib/logger"
)

var (
	logger     = liblogger.NewLogger(dbusNetworkDest)
	manager    *Manager
	running    bool
	notifyStop = make(chan int, 100)
)

func Start() {
	logger.BeginTracing()
	defer logger.EndTracing()

	if running {
		logger.Info(dbusNetworkDest, "already running")
		return
	}
	running = true
	defer func() {
		running = false
	}()

	// TODO
	/*
		if !dlib.UniqueOnSession(dbusNetworkDest) {
			logger.Warning("dbus unique:", dbusNetworkDest)
			return
		}
	*/

	initSlices() // initialize slice code

	manager = NewManager()
	err := dbus.InstallOnSession(manager)
	if err != nil {
		logger.Error("register dbus interface failed: ", err)
		return
	}

	// initialize manager after dbus installed
	manager.initManager()
	dbus.DealWithUnhandledMessage()

	notifyStop = make(chan int, 100) // reset signal to avoid repeat stop action
	select {
	case <-notifyStop:
	}
	DestroyManager(manager)
}

func Stop() {
	if !running {
		logger.Info(dbusNetworkDest, "already stopped")
		return
	}
	notifyStop <- 1
}
