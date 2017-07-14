package dbus

import (
	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
)

// IsSessionBusActivated check the special session bus name whether activated
func IsSessionBusActivated(dest string) bool {
	if !lib.UniqueOnSession(dest) {
		return true
	}

	bus, _ := dbus.SessionBus()
	releaseDBusName(bus, dest)
	return false
}

// IsSystemBusActivated check the special system bus name whether activated
func IsSystemBusActivated(dest string) bool {
	if !lib.UniqueOnSystem(dest) {
		return true
	}

	bus, _ := dbus.SystemBus()
	releaseDBusName(bus, dest)
	return false
}

func releaseDBusName(bus *dbus.Conn, name string) {
	if bus != nil {
		bus.ReleaseName(name)
	}
}
