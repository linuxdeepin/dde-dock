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
	"dbus/com/deepin/dde/dock"
)

const (
	dockName = "dde-dock"
	dockDest = "com.deepin.dde.dock"
	dockPath = "/com/deepin/dde/dock"
)

func isDockRunning() bool {
	return isDBusDestExist(dockDest)
}

func launchDock() error {
	caller, err := dock.NewDock(dockDest, dockPath)
	if err != nil {
		return err
	}

	_, err = caller.Xid()
	return err
}

func newDockTask() *taskInfo {
	return newTaskInfo(dockName, isDockRunning, launchDock)
}
