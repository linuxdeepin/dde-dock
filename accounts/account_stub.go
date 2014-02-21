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

type AccountManager struct {
        UserList []string

        UserAdded   func()
        UserDeleted func()
}

func (op *AccountManager) GetDBusInfo() dbus.DBusInfo {
        return dbus.DBusInfo{
                ACCOUNT_DEST,
                ACCOUNT_MANAGER_PATH,
                ACCOUNT_MANAGER_IFC,
        }
}

func (op *AccountManager) setPropName(name string) {
        switch name {
        case "UserList":
                infos := getUserInfoList()
                list := []string{}

                for _, info := range infos {
                        path := USER_MANAGER_PATH + info.Uid
                        list = append(list, path)
                }
                op.UserList = list
                dbus.NotifyChange(op, name)
        }
}

func newAccountManager() *AccountManager {
        m := &AccountManager{}

        m.setPropName("UserList")

        return m
}
