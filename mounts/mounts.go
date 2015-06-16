/**
 * Copyright (c) 2011 ~ 2015 Deepin, Inc.
 *               2013 ~ 2015 jouyouyun
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

package mounts

import (
	"pkg.linuxdeepin.com/dde-daemon"
	"pkg.linuxdeepin.com/lib/dbus"
)

var (
	_manager *Manager
)

func init() {
	loader.Register(&loader.Module{
		Name:   "mounts",
		Start:  Start,
		Stop:   Stop,
		Enable: true,
	})
}

func Start() {
	if _manager != nil {
		return
	}

	_manager = NewManager()
	_manager.logger.BeginTracing()
	err := dbus.InstallOnSession(_manager)
	if err != nil {
		_manager.logger.Error("Install mounts dbus failed:", err)
		_manager.destroy()
		_manager = nil
		return
	}
	_manager.listenDiskChanged()
	go _manager.refrashDiskInfos()
}

func Stop() {
	if _manager == nil {
		return
	}

	_manager.destroy()
	_manager = nil
}
