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

package touchpad

import (
	"pkg.deepin.io/dde/daemon/inputdevices/wrapper"
	"pkg.deepin.io/lib/dbus"
)

const (
	DBUS_SENDER    = "com.deepin.daemon.InputDevices"
	DBUS_PATH_TPAD = "/com/deepin/daemon/InputDevice/TouchPad"
	DBUS_IFC_TPAD  = "com.deepin.daemon.InputDevice.TouchPad"
)

func (tpad *Touchpad) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       DBUS_SENDER,
		ObjectPath: DBUS_PATH_TPAD,
		Interface:  DBUS_IFC_TPAD,
	}
}

func (tpad *Touchpad) setPropDeviceList(devList []wrapper.XIDeviceInfo) {
	tpad.DeviceList = devList
	dbus.NotifyChange(tpad, "DeviceList")
}

func (tpad *Touchpad) setPropExist(exist bool) {
	if tpad.Exist != exist {
		tpad.Exist = exist
		dbus.NotifyChange(tpad, "Exist")
	}
}
