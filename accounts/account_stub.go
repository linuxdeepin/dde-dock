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
	"fmt"
	"pkg.linuxdeepin.com/lib/dbus"
)

func (obj *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       ACCOUNT_DEST,
		ObjectPath: ACCOUNT_MANAGER_PATH,
		Interface:  ACCOUNT_MANAGER_IFC,
	}
}

func (obj *Manager) setPropUserList(list []string) {
	if !isStrListEqual(obj.UserList, list) {
		obj.UserList = list
		dbus.NotifyChange(obj, "UserList")
	}
}

func (obj *Manager) setPropAllowGuest(isAllow bool) {
	if obj.AllowGuest != isAllow {
		obj.AllowGuest = isAllow
		dbus.NotifyChange(obj, "AllowGuest")
	}
}

func (m *Manager) setPropGuestIcon(icon string) {
	if m.GuestIcon != icon {
		m.GuestIcon = icon
		dbus.NotifyChange(m, "GuestIcon")
	}
}

func (obj *Manager) newUserByPath(path string) error {
	if len(path) == 0 {
		return fmt.Errorf("Invalid ObjectPath")
	}

	u := newUser(path)
	if u == nil {
		return fmt.Errorf("Create User Object Failed")
	}
	if err := dbus.InstallOnSystem(u); err != nil {
		return fmt.Errorf("Install DBus For %s Failed: %v", path, err)
	}
	u.setProps()

	obj.pathUserMap[path] = u

	return nil
}

func (m *Manager) destroyUser(path string) {
	u, ok := m.pathUserMap[path]
	if !ok {
		return
	}

	if u.watcher != nil {
		u.quitFlag <- true
		u.watcher.Close()
	}
	dbus.UnInstallObject(u)
	u = nil
	delete(m.pathUserMap, path)
}

func (m *Manager) destroyAllUser() {
	for _, path := range m.UserList {
		m.destroyUser(path)
	}
	m.pathUserMap = make(map[string]*User)
}

func (obj *Manager) updateAllUserInfo() {
	obj.destroyAllUser()

	for _, path := range obj.UserList {
		err := obj.newUserByPath(path)
		if err != nil {
			logger.Error(err)
		}
	}
}

func (m *Manager) handleuserAdded(list []string) {
	for _, path := range list {
		err := m.newUserByPath(path)
		if err != nil {
			logger.Error(err)
			continue
		}
		dbus.Emit(m, "UserAdded", path)
	}

	m.setPropUserList(getUserList())
}

func (m *Manager) handleUserRemoved(list []string) {
	for _, path := range list {
		m.destroyUser(path)
		dbus.Emit(m, "UserDeleted", path)
	}

	m.setPropUserList(getUserList())
}

func (obj *Manager) destroy() {
	if obj.listWatcher != nil {
		obj.listQuit <- true
		obj.listWatcher.Close()
	}

	if obj.infoWatcher != nil {
		obj.infoQuit <- true
		obj.infoWatcher.Close()
	}

	obj.destroyAllUser()
}
