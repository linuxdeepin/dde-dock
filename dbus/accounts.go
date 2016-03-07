/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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
