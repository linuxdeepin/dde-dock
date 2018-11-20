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
	"os"

	"pkg.deepin.io/dde/daemon/grub2"
	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger("daemon/grub2")

func init() {
	grub2.SetLogger(logger)
}

var (
	argPrepareGfxmodeDetect bool
	argDebug                bool
)

func main() {
	flag.BoolVar(&argDebug, "debug", false, "debug mode")
	flag.BoolVar(&argPrepareGfxmodeDetect, "prepare-gfxmode-detect", false,
		"prepare gfxmode detect")
	flag.Parse()
	if argDebug {
		logger.SetLogLevel(log.LevelDebug)
	}

	if argPrepareGfxmodeDetect {
		logger.Debug("mode: prepare gfxmode detect")
		err := grub2.PrepareGfxmodeDetect()
		if err != nil {
			logger.Warning(err)
			os.Exit(2)
		}
	} else {
		logger.Debug("mode: daemon")
		grub2.RunAsDaemon()
	}
}
