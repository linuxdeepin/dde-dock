/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
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

package main

import (
	"flag"
	"fmt"
	"os"

	"pkg.deepin.io/dde/daemon/grub2"
	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger("daemon/grub2")

func init() {
	grub2.SetLogger(logger)
}

var (
	argDebug      bool
	argSetup      bool
	argSetupTheme bool
	argResolution string
)

func main() {
	flag.BoolVar(&argDebug, "debug", false, "debug mode")
	flag.BoolVar(&argSetup, "setup", false, "setup grub and exit")
	flag.BoolVar(&argSetupTheme, "setup-theme", false, "setup grub theme only and exit")
	//flag.StringVar(&argGrubSettingFile, "setting-file", "", "specify an alternative setting file instead of /etc/default/grub when setup grub")
	// TODO --grub-config, --backend, [grub, efi]
	//flag.StringVar(&argThemeDir, "theme-dir", "", "specify an alternative theme directory instead of /boot/grub/themes/deepin when setup grub")
	flag.StringVar(&argResolution, "gfxmode", "auto", "specify gfxmode when setup grub")
	flag.Parse()
	if argDebug {
		logger.SetLogLevel(log.LevelDebug)
	}

	if argSetupTheme {
		fmt.Println("mode: setup theme")
		err := grub2.SetupTheme()
		if err != nil {
			logger.Warning(err)
			os.Exit(1)
		}
	} else if argSetup {
		fmt.Println("mode: setup")
		err := grub2.Setup(argResolution)
		if err != nil {
			logger.Warning(err)
			os.Exit(2)
		}
	} else {
		fmt.Println("mode: daemon")
		grub2.RunAsDaemon()
	}
}
