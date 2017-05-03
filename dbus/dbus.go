package dbus

import (
	"pkg.deepin.io/lib"
)

func IsSessionBusActivated(dest string) bool {
	return !lib.UniqueOnSession(dest)
}

func IsSystemBusActivated(dest string) bool {
	return !lib.UniqueOnSystem(dest)
}
