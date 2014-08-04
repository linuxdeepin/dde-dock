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
	"os"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/log"
)

var (
	logger = log.NewLogger(DbusGrubDest)
	grub   *Grub2
)

func Start() {
	logger.BeginTracing()

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
	// property changing signal send successful
	grub.initGrub2()
}

func Stop() {
	DestroyGrub2(grub)
	logger.EndTracing()
}

func SetLogLevel(level log.Priority) {
	logger.SetLogLevel(level)
}

func IsUpdating() bool {
	if grub.Updating || grub.theme.Updating {
		return true
	} else {
		return false
	}
}
