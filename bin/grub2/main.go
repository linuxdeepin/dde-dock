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
	"dde-daemon/grub2"
	"flag"
	"fmt"
	"os"
	"pkg.linuxdeepin.com/lib"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/log"
	"time"
)

const dbusGrubDest = "com.deepin.daemon.Grub2"

var (
	argDebug      bool
	argSetup      bool
	argSetupTheme bool
	argConfig     string
	argThemeDir   string
	argGfxmode    string
)

func main() {
	flag.BoolVar(&argDebug, "d", false, "debug mode")
	flag.BoolVar(&argDebug, "debug", false, "debug mode")
	flag.BoolVar(&argSetup, "setup", false, "setup grub and exit")
	flag.BoolVar(&argSetupTheme, "setup-theme", false, "setup grub theme only and exit")
	flag.StringVar(&argConfig, "config", "", "specify an alternative configuration file when setup grub")
	flag.StringVar(&argThemeDir, "theme-dir", "", "specify an alternative theme directory when setup grub")
	flag.StringVar(&argGfxmode, "gfxmode", "", "specify gfxmode when setup grub")
	flag.Parse()

	if argDebug {
		grub2.SetLoggerLevel(log.LEVEL_DEBUG)
	}

	// dispatch optional arguments
	if len(argConfig) != 0 {
		fmt.Println("config:", argConfig)
		grub2.SetDefaultGrubConfigFile(argConfig)
	}
	if len(argThemeDir) != 0 {
		fmt.Println("theme dir:", argThemeDir)
		grub2.SetDefaultThemePath(argThemeDir)
	}
	if len(argGfxmode) != 0 {
		fmt.Println("gfxmode:", argGfxmode)
	}

	g := grub2.NewGrub2()
	if argSetupTheme {
		fmt.Println("setup theme: true")
		g.SetupTheme(argGfxmode)
	} else if argSetup {
		fmt.Println("setup mode: true")
		g.Setup(argGfxmode)
	} else {
		runAsDaemon()
	}
}

func runAsDaemon() {
	logger := log.NewLogger(dbusGrubDest + ".Wrapper")
	logger.BeginTracing()
	defer logger.EndTracing()

	if !lib.UniqueOnSession(dbusGrubDest) {
		logger.Error("dbus unique:", dbusGrubDest)
		return
	}
	grub2.Start()
	dbus.SetAutoDestroyHandler(60*time.Second, func() bool {
		return !grub2.IsUpdating()
	})
	dbus.DealWithUnhandledMessage()
	if err := dbus.Wait(); err != nil {
		logger.Error("lost dbus session:", err)
		os.Exit(1)
	}
}
