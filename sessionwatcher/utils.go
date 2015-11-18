package sessionwatcher

import (
	"dbus/org/freedesktop/dbus"
)

func isDBusDestExist(dest string) bool {
	daemon, err := dbus.NewDBusDaemon("org.freedesktop.DBus", "/")
	if err != nil {
		return false
	}
	defer dbus.DestroyDBusDaemon(daemon)

	names, err := daemon.ListNames()
	if err != nil {
		return false
	}
	return isItemInList(dest, names)
}

func isItemInList(item string, list []string) bool {
	for _, v := range list {
		if item == v {
			return true
		}
	}
	return false
}
