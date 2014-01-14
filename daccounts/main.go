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
	"dbus/org/freedesktop/accounts"
	"dlib/dbus"
	"fmt"
)

type Manager struct {
	UserAdded   func(string)
	UserDeleted func(string)
}

type User struct {
	AccountType    int32  `access:"readwrite"`
	AutomaticLogin bool   `access:"readwrite"`
	IconFile       string `access:"readwrite"`
	Locked         bool   `access:"readwrite"`
	PasswordMode   int32  `access:"readwrite"`
	UserName       string `access:"readwrite"`
        BackgroundFile string `access:"readwrite"`
	LoginTime      int64
	Uid            uint64
	objectPath     string
	userInface     *accounts.User
}

const (
	_ACCOUNTS_DEST = "com.deepin.daemon.Accounts"
	_ACCOUNTS_PATH = "/com/deepin/daemon/Accounts"
	_ACCOUNTS_IFC  = "com.deepin.daemon.Accounts"

	_ACCOUNTS_USER_IFC = "com.deepin.daemon.Accounts.User"
)

var (
	_accountInface *accounts.Accounts
	_userMap       = make(map[dbus.ObjectPath]*User)
)

func (dam *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_ACCOUNTS_DEST,
		_ACCOUNTS_PATH,
		_ACCOUNTS_IFC,
	}
}

func (u *User) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_ACCOUNTS_DEST,
		u.objectPath,
		_ACCOUNTS_USER_IFC,
	}
}

func main() {
	var err error
	_accountInface, err = accounts.NewAccounts("/org/freedesktop/Accounts")
	if err != nil {
		fmt.Println("New Accounts Failed From Freedesktop:", err)
		return
	}

	account := NewAccountManager()
	err = dbus.InstallOnSession(account)
	if err != nil {
		fmt.Println("Install Manager DBus Failed")
		panic(err)
	}
	userList, _ := _accountInface.ListCachedUsers()
	for _, v := range userList {
		NewAccountUserManager(v)
	}
	dbus.DealWithUnhandledMessage()

	select {}
}
