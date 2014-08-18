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

/*
func (op *KbdEntry) OnPropertiesChanged(key string, old interface{}) {
        switch key {
        case KBD_KEY_LAYOUT:
                if v, ok := old.(string); ok && v != op.CurrentLayout {
                        op.AddUserLayout(op.CurrentLayout, op.LayoutOption)
                }
        case KBD_KEY_LAYOUT_OPTION:
                if v, ok := old.(string); ok && v != op.LayoutOption {
                        op.AddUserLayout(op.CurrentLayout, op.LayoutOption)
                }
        }
}
*/

func (op *MouseEntry) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		DEVICE_DEST,
		DEVICE_PATH + op.deviceId,
		DEVICE_IFC + op.deviceId,
	}
}

func (mouse *MouseEntry) setPropExist(exist bool) {
	if mouse.Exist != exist {
		mouse.Exist = exist
		dbus.NotifyChange(mouse, "Exist")
	}
}

func (op *TPadEntry) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		DEVICE_DEST,
		DEVICE_PATH + op.deviceId,
		DEVICE_IFC + op.deviceId,
	}
}

func (op *KbdEntry) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		DEVICE_DEST,
		DEVICE_PATH + op.deviceId,
		DEVICE_IFC + op.deviceId,
	}
}
