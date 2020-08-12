/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
		_, _ = bus.ReleaseName(name)
	}
}
