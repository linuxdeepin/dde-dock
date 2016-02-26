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
