/**
 * Copyright (c) 2013 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 Xu FaSheng
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

package grub2

import (
	"dlib"
	"dlib/dbus"
	liblogger "dlib/logger"
	"os"
	"time"
)

var (
	logger = liblogger.NewLogger(dbusGrubDest)
	grub   *Grub2
)

func RunAsDaemon() {
	if !dlib.UniqueOnSession(dbusGrubDest) {
		logger.Warning("dbus unique:", dbusGrubDest)
		return
	}
	Start()
	// TODO
	dbus.SetAutoDestroyHandler(1*time.Second, func() bool {
		if grub.Updating || grub.theme.Updating {
			return false
		} else {
			return true
		}
	})
	if err := dbus.Wait(); err != nil {
		logger.Error("lost dbus session:", err)
		os.Exit(1)
	}
}

func Start() {
	logger.BeginTracing()
	defer logger.EndTracing()

	grub = NewGrub2()
	err := dbus.InstallOnSession(grub)
	if err != nil {
		logger.Errorf("register dbus interface failed: %v", err)
		os.Exit(1)
	}
	err = dbus.InstallOnSession(grub.theme)
	if err != nil {
		logger.Errorf("register dbus interface failed: %v", err)
		os.Exit(1)
	}

	// initialize grub2 after dbus service installed to ensure
	// property changed signal send success
	grub.initGrub2()
	dbus.DealWithUnhandledMessage()
}

func Stop() {
	DestroyGrub2(grub)
}

func SetLoggerLevel(level liblogger.Priority) {
	logger.SetLogLevel(level)
}
