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
