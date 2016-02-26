/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package sessionwatcher

import (
	"dbus/com/deepin/dde/desktop"
)

const (
	desktopName = "dde-desktop"
	desktopDest = "com.deepin.dde.desktop"
	desktopPath = "/com/deepin/dde/desktop"
)

func isDesktopRunning() bool {
	return isDBusDestExist(desktopDest)
}

func launchDesktop() error {
	caller, err := desktop.NewDesktop(desktopDest, desktopPath)
	if err != nil {
		return err
	}
	defer desktop.DestroyDesktop(caller)

	return caller.Show()
}

func newDesktopTask() *taskInfo {
	return newTaskInfo(desktopName, isDesktopRunning, launchDesktop)
}
