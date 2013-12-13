/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
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

package main

import (
	"dlib"
	"dlib/dbus"
	/*"fmt"*/)

/*type IndividuateManager struct {}*/

const (
	_INDIVI_DEST     = "com.deepin.daemon.IndividuateManager"
	_BG_MANAGER_PATH = "/com/deepin/Individuate/BackgroundManager"
	_BG_MANAGER_IFC  = "com.deepin.daemon.Individuate.BackgroundManager"
)

func (bgManager *BackgroundManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_INDIVI_DEST,
		_BG_MANAGER_PATH,
		_BG_MANAGER_IFC,
	}
}

func main() {
	bgManager := NewBackgroundManager()
	err := dbus.InstallOnSession(bgManager)
	if err != nil {
		panic(err)
	}

	dlib.StartLoop()
}
