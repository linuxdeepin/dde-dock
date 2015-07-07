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

package keybinding

import (
	"pkg.deepin.io/lib/dbus"
)

func (obj *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       KEYBIND_DEST,
		ObjectPath: MANAGER_PATH,
		Interface:  MANAGER_IFC,
	}
}

func (obj *MediaKeyManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       KEYBIND_DEST,
		ObjectPath: MEDIAKEY_PATH,
		Interface:  MEDIAKEY_IFC,
	}
}

func (obj *Manager) setPropSystemList(list []ShortcutInfo) {
	if !compareShortcutInfoList(list, obj.SystemList) {
		obj.SystemList = list
		dbus.NotifyChange(obj, "SystemList")
	}
}

func (obj *Manager) setPropWindowList(list []ShortcutInfo) {
	if !compareShortcutInfoList(obj.WindowList, list) {
		obj.WindowList = list
		dbus.NotifyChange(obj, "WindowList")
	}
}

func (obj *Manager) setPropWorkspaceList(list []ShortcutInfo) {
	if !compareShortcutInfoList(obj.WorkspaceList, list) {
		obj.WorkspaceList = list
		dbus.NotifyChange(obj, "WorkspaceList")
	}
}

func (obj *Manager) setPropCustomList(list []ShortcutInfo) {
	if !compareShortcutInfoList(obj.CustomList, list) {
		obj.CustomList = list
		dbus.NotifyChange(obj, "CustomList")
	}
}

func (obj *Manager) setPropConflictValid(list []int32) {
	if !compareInt32List(obj.ConflictValid, list) {
		obj.ConflictValid = list
		dbus.NotifyChange(obj, "ConflictValid")
	}
}

func (obj *Manager) setPropConflictInvalid(list []int32) {
	if !compareInt32List(obj.ConflictInvalid, list) {
		obj.ConflictInvalid = list
		dbus.NotifyChange(obj, "ConflictInvalid")
	}
}

func (obj *Manager) updateProps() {
	obj.setPropSystemList(getSystemListInfo())
	obj.setPropWindowList(getWindowListInfo())
	obj.setPropWorkspaceList(getWorkspaceListInfo())
	obj.setPropCustomList(getCustomListInfo())
	obj.setPropConflictValid(getValidConflictList())
	obj.setPropConflictInvalid(getInvalidConflictList())
}
