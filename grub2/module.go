/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package grub2

import (
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

var (
	logger = log.NewLogger("daemon/grub2")
	grub   *Grub2
)

func Start() {
	logger.BeginTracing()

	grub = NewGrub2()
	err := dbus.InstallOnSession(grub)
	if err != nil {
		logger.Errorf("register dbus interface failed: %v", err)
		return
	}
	err = dbus.InstallOnSession(grub.theme)
	if err != nil {
		logger.Errorf("register dbus interface failed: %v", err)
		return
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
