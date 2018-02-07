/*
 * Copyright (C) 2015 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package dbus

import (
	"dbus/com/deepin/daemon/accounts"
	"pkg.deepin.io/lib/dbus"
)

const (
	accountsDBusDest = "com.deepin.daemon.Accounts"
	accountsDBusPath = "/com/deepin/daemon/Accounts"
)

func NewAccounts() (*accounts.Accounts, error) {
	return accounts.NewAccounts(accountsDBusDest, accountsDBusPath)
}

func NewUserByName(name string) (*accounts.User, error) {
	m, err := NewAccounts()
	if err != nil {
		return nil, err
	}

	p, err := m.FindUserByName(name)
	if err != nil {
		return nil, err
	}
	return accounts.NewUser(accountsDBusDest, dbus.ObjectPath(p))
}

func NewUserByUid(uid string) (*accounts.User, error) {
	m, err := NewAccounts()
	if err != nil {
		return nil, err
	}

	p, err := m.FindUserById(uid)
	if err != nil {
		return nil, err
	}
	return accounts.NewUser(accountsDBusDest, dbus.ObjectPath(p))
}

func DestroyAccounts(act *accounts.Accounts) {
	if act == nil {
		return
	}
	accounts.DestroyAccounts(act)
}

func DestroyUser(u *accounts.User) {
	if u == nil {
		return
	}
	accounts.DestroyUser(u)
}
