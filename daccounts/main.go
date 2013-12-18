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

type AccountManager struct {
	UserAdded   func(string)
	UserDeleted func(string)
}

type AccountUserManager struct {
	AccountType    int32
	AutomaticLogin bool
	IconFile       string
	Locked         bool
	PasswordMode   int32
	UserName       string
	LoginTime      int64
	ObjectPath     string

	PropertyChanged func(string)
}

const (
	_ACCOUNTS_DEST = "com.deepin.daemon.Accounts"
	_ACCOUNTS_PATH = "/com/deepin/daemon/Accounts"
	_ACCOUNTS_IFC  = "com.deepin.daemon.Accounts"

	_ACCOUNTS_USER_IFC = "com.deepin.daemon.Accounts.User"
)

var (
	_accountInface = accounts.GetAccounts("/org/freedesktop/Accounts")
	_userMap       = make(map[string]*AccountUserManager)
	_infaceMap     = make(map[string]*accounts.User)
)

func (dam *AccountManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_ACCOUNTS_DEST,
		_ACCOUNTS_PATH,
		_ACCOUNTS_IFC,
	}
}

func (userManager *AccountUserManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_ACCOUNTS_DEST,
		userManager.ObjectPath,
		_ACCOUNTS_USER_IFC,
	}
}

func main() {
	account := NewAccountManager()
	err := dbus.InstallOnSession(account)
	if err != nil {
		fmt.Println("Install AccountManager DBus Failed")
		panic(err)
	}
	userList := account.ListCachedUsers()
	for _, v := range userList {
		userManager := NewAccountUserManager(v)
		dbus.InstallOnSession(userManager)
	}

	select {}
}
