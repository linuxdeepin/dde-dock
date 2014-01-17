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
	"dlib/dbus"
)

func (dev *ExtDevManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{_EXT_DEV_NAME, _EXT_DEV_PATH, _EXT_DEV_IFC}
}

func (keyboard *KeyboardEntry) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_EXT_DEV_NAME,
		_EXT_ENTRY_PATH + keyboard.DeviceID,
		_EXT_ENTRY_IFC + keyboard.DeviceID,
	}
}

func (mouse *MouseEntry) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_EXT_DEV_NAME,
		_EXT_ENTRY_PATH + mouse.DeviceID,
		_EXT_ENTRY_IFC + mouse.DeviceID,
	}
}

func (tpad *TPadEntry) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_EXT_DEV_NAME,
		_EXT_ENTRY_PATH + tpad.DeviceID,
		_EXT_ENTRY_IFC + tpad.DeviceID,
	}
}
