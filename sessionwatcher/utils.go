/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package sessionwatcher

import (
	"dbus/org/freedesktop/dbus"
)

func isDBusSenderExist(sender string) bool {
	dbusDaemon, err := dbus.NewDBusDaemon("org.freedesktop.DBus", "/")
	if err != nil {
		logger.Error("Create freedesktop DBus failed:", err)
		return false
	}
	defer dbus.DestroyDBusDaemon(dbusDaemon)

	names, _ := dbusDaemon.ListNames()
	for _, name := range names {
		if name == sender {
			return true
		}
	}

	return false
}
