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
	DBUS_PATH_WACOM = "/com/deepin/daemon/InputDevice/Wacom"
	DBUS_IFC_WACOM  = "com.deepin.daemon.InputDevice.Wacom"
)

func (wManager *WacomManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		DBUS_SENDER,
		DBUS_PATH_WACOM,
		DBUS_IFC_WACOM,
	}
}

func (wManager *WacomManager) setPropDeviceList(list []PointerDeviceInfo) {
	if len(wManager.DeviceList) != len(list) {
		wManager.DeviceList = list
		dbus.NotifyChange(wManager, "DeviceList")
	}
}

func (wManager *WacomManager) setPropExist(exist bool) {
	if wManager.Exist != exist {
		wManager.Exist = exist
		dbus.NotifyChange(wManager, "Exist")
	}
}
