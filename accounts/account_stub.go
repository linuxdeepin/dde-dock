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

package accounts

import (
	"dlib/dbus"
	"sync"
)

func (obj *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		ACCOUNT_DEST,
		ACCOUNT_MANAGER_PATH,
		ACCOUNT_MANAGER_IFC,
	}
}

func (obj *Manager) updatePropUserList(list []string) {
	if !isStrListEqual(obj.UserList, list) {
		obj.UserList = list
		dbus.NotifyChange(obj, "UserList")
	}
}

func (obj *Manager) updatePropAllowGuest(isAllow bool) {
	if obj.AllowGuest != isAllow {
		obj.AllowGuest = isAllow
		dbus.NotifyChange(obj, "AllowGuest")
	}
}

func (obj *Manager) updateUserInfo(path string) {
	if len(path) < 1 {
		return
	}

	u := newUser(path)
	if u == nil {
		return
	}
	if err := dbus.InstallOnSystem(u); err != nil {
		logger.Errorf("Install DBus For %s Failed: %v", path, err)
		panic(err)
	}

	obj.pathUserMap[path] = u
}

func (obj *Manager) destroyAllUser() {
	var mutex sync.Mutex
	mutex.Lock()
	if len(obj.pathUserMap) > 0 {
		for _, v := range obj.pathUserMap {
			v.endFlag <- true
			v.endWatchFlag <- true
			v.watcher.Close()
			dbus.UnInstallObject(v)
		}

		obj.pathUserMap = make(map[string]*User)
	}
	mutex.Unlock()
}

func (obj *Manager) updateAllUserInfo() {
	obj.destroyAllUser()

	for _, path := range obj.UserList {
		obj.updateUserInfo(path)
	}
}
