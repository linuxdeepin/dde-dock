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
	"dlib/dbus"
	liblogger "dlib/logger"
)

var (
	logger     = liblogger.NewLogger(dbusBluetoothDest)
	bluetooth  *Bluetooth
	running    bool
	notifyStop = make(chan int, 100)
)

func Start() {
	logger.BeginTracing()
	defer logger.EndTracing()

	if running {
		logger.Info(dbusBluetoothDest, "already running")
		return
	}
	running = true
	defer func() {
		running = false
	}()

	//if !dlib.UniqueOnSession(dbusBluetoothDest) {
	//logger.Warning("dbus unique:", dbusBluetoothDest)
	//return
	//}

	bluetooth = NewBluetooth()
	err := dbus.InstallOnSession(bluetooth)
	if err != nil {
		// don't panic or fatal here
		logger.Error("register dbus interface failed: ", err)
		return
	}

	// initialize bluetooth after dbus interface installed
	bluetooth.initBluetooth()
	dbus.DealWithUnhandledMessage()

	notifyStop = make(chan int, 100) // reset signal to avoid repeat stop action
	notfiyDbusStop := make(chan int)
	/*
		go func() {
			err := dbus.Wait()
			if err != nil {
				logger.Error("lost dbus session:", err)
			} else {
				logger.Info("dbus session stoped")
			}
			notfiyDbusStop <- 1
		}()
	*/

	select {
	case <-notifyStop:
	case <-notfiyDbusStop:
	}
	DestroyBluetooth(bluetooth)
}

func Stop() {
	if !running {
		logger.Info(dbusBluetoothDest, "already stopped")
		return
	}
	notifyStop <- 1
}
