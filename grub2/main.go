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

package main

import (
	apigrub2ext "dbus/com/deepin/api/grub2"
	"dlib/dbus"
	liblogger "dlib/logger"
	"flag"
	"os"
)

var (
	logger        *liblogger.Logger
	argDebug      bool
	argSetup      bool
	argSetupTheme bool
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			logger.Fatalf("%v", err)
		}
	}()

	flag.BoolVar(&argDebug, "d", false, "debug mode")
	flag.BoolVar(&argDebug, "debug", false, "debug mode")
	flag.BoolVar(&argSetup, "setup", false, "setup grub and exit")
	flag.BoolVar(&argSetupTheme, "setup-theme", false, "setup grub theme only and exit")
	flag.Parse()

	// configure logger
	logger = liblogger.NewLogger("dde-daemon/grub2")
	logger.SetRestartCommand("/usr/lib/deepin-daemon/grub2", "--debug")
	if argDebug {
		logger.SetLogLevel(liblogger.LEVEL_DEBUG)
	}

	grub := NewGrub2()

	// setup grub and exit
	if argSetup {
		grub.setup()
		os.Exit(0)
	}

	// setup grub theme only and exit
	if argSetupTheme {
		grub.setupTheme()
		os.Exit(0)
	}

	grub2ext, _ = apigrub2ext.NewGrub2Ext("com.deepin.api.Grub2", "/com/deepin/api/Grub2")
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

	// load after dbus service installed to ensure property changed
	// signal send success
	grub.load()
	grub.theme.load()
	go grub.resetGfxmodeIfNeed()
	go grub.theme.regenerateBackgroundIfNeed()
	grub.startUpdateLoop()

	dbus.DealWithUnhandledMessage()
	if err := dbus.Wait(); err != nil {
		logger.Errorf("lost dbus session: %v", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
