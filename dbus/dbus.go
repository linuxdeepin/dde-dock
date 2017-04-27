package dbus

import (
	ddbus "dbus/org/freedesktop/dbus"
	dsystem "dbus/org/freedesktop/dbus/system"
	"pkg.deepin.io/lib/strv"
)

func IsSessionBusActivated(dest string) bool {
	bus, err := ddbus.NewDBusDaemon("org.freedesktop.DBus", "/")
	if err != nil {
		return false
	}
	defer ddbus.DestroyDBusDaemon(bus)

	list, err := bus.ListActivatableNames()
	if err != nil {
		return false
	}
	return strv.Strv(list).Contains(dest)
}

func IsSystemBusActivated(dest string) bool {
	bus, err := dsystem.NewDBusDaemon("org.freedesktop.DBus", "/")
	if err != nil {
		return false
	}
	defer dsystem.DestroyDBusDaemon(bus)

	list, err := bus.ListActivatableNames()
	if err != nil {
		return false
	}
	return strv.Strv(list).Contains(dest)
}
