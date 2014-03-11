/**
 * Copyright (c) 2014 Deepin, Inc.
 *               2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
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
	"dlib/dbus"
)

const (
	WALLPAPER_DEST = "com.deepin.daemon.Themes"
	WALLPAPER_PATH = "/com/deepin/daemon/Wallpaper"
	WALLPAPER_IFC  = "com.deepin.daemon.Wallpaper"
)

func (w *Wallpaper) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		WALLPAPER_DEST,
		WALLPAPER_PATH,
		WALLPAPER_IFC,
	}
}
