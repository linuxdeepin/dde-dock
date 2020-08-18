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
	"github.com/godbus/dbus"
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.accounts"
)

func NewAccounts(systemConn *dbus.Conn) *accounts.Accounts {
	return accounts.NewAccounts(systemConn)
}

func NewUserByName(systemConn *dbus.Conn, name string) (*accounts.User, error) {
	m := NewAccounts(systemConn)
	userPath, err := m.FindUserByName(0, name)
	if err != nil {
		return nil, err
	}
	return accounts.NewUser(systemConn, dbus.ObjectPath(userPath))
}

func NewUserByUid(systemConn *dbus.Conn, uid string) (*accounts.User, error) {
	m := NewAccounts(systemConn)
	userPath, err := m.FindUserById(0, uid)
	if err != nil {
		return nil, err
	}
	return accounts.NewUser(systemConn, dbus.ObjectPath(userPath))
}
