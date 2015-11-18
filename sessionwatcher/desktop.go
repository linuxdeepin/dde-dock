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
