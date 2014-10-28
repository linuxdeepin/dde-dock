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

package libtouchpad

import (
	"pkg.linuxdeepin.com/dde-daemon/inputdevices/libwrapper"
	"pkg.linuxdeepin.com/lib/dbus"
)

const (
	DBUS_SENDER    = "com.deepin.daemon.InputDevices"
	DBUS_PATH_TPAD = "/com/deepin/daemon/InputDevice/TouchPad"
	DBUS_IFC_TPAD  = "com.deepin.daemon.InputDevice.TouchPad"
)

func (touchpad *Touchpad) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		DBUS_SENDER,
		DBUS_PATH_TPAD,
		DBUS_IFC_TPAD,
	}
}

func (touchpad *Touchpad) setPropDeviceList(devList []libwrapper.XIDeviceInfo) {
	touchpad.DeviceList = devList
	dbus.NotifyChange(touchpad, "DeviceList")
}

func (touchpad *Touchpad) setPropExist(exist bool) {
	if touchpad.Exist != exist {
		touchpad.Exist = exist
		dbus.NotifyChange(touchpad, "Exist")
	}
}
