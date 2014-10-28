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

package libmouse

import (
	"pkg.linuxdeepin.com/dde-daemon/inputdevices/libwrapper"
	"pkg.linuxdeepin.com/lib/dbus"
)

const (
	DBUS_SENDER     = "com.deepin.daemon.InputDevices"
	DBUS_PATH_MOUSE = "/com/deepin/daemon/InputDevice/Mouse"
	DBUS_IFC_MOUSE  = "com.deepin.daemon.InputDevice.Mouse"
)

func (mouse *Mouse) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		DBUS_SENDER,
		DBUS_PATH_MOUSE,
		DBUS_IFC_MOUSE,
	}
}

func (mouse *Mouse) setPropDeviceList(devList []libwrapper.XIDeviceInfo) {
	mouse.DeviceList = devList
	dbus.NotifyChange(mouse, "DeviceList")
}

func (mouse *Mouse) setPropExist(exist bool) {
	if mouse.Exist != exist {
		mouse.Exist = exist
		dbus.NotifyChange(mouse, "Exist")
	}
}
