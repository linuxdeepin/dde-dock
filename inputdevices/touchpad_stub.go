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

package inputdevices

import (
	"pkg.linuxdeepin.com/lib/dbus"
)

const (
	DBUS_PATH_TPAD = "/com/deepin/daemon/InputDevice/TouchPad"
	DBUS_IFC_TPAD  = "com.deepin.daemon.InputDevice.TouchPad"
)

func (tManager *TouchpadManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		DBUS_SENDER,
		DBUS_PATH_TPAD,
		DBUS_IFC_TPAD,
	}
}

func (tManager *TouchpadManager) setPropDeviceList(devList []PointerDeviceInfo) {
	tManager.DeviceList = devList
	dbus.NotifyChange(tManager, "DeviceList")
}

func (tManager *TouchpadManager) setPropExist(exist bool) {
	if tManager.Exist != exist {
		tManager.Exist = exist
		dbus.NotifyChange(tManager, "Exist")
	}
}
