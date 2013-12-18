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
)

func (userManager *AccountUserManager) SetAccountType(accountType int32) {
	userInface := _infaceMap[userManager.ObjectPath]
	if userInface == nil {
		return
	}
	userInface.SetAccountType(accountType)
}

func (userManager *AccountUserManager) SetAutomaticLogin(enabled bool) {
	userInface := _infaceMap[userManager.ObjectPath]
	if userInface == nil {
		return
	}
	userInface.SetAutomaticLogin(enabled)
}

func (userManager *AccountUserManager) SetIconFile(filename string) {
	userInface := _infaceMap[userManager.ObjectPath]
	if userInface == nil {
		return
	}
	userInface.SetIconFile(filename)
	userManager.AddIconToHistory(filename)
}

func (userManager *AccountUserManager) SetLocked(locked bool) {
	userInface := _infaceMap[userManager.ObjectPath]
	if userInface == nil {
		return
	}
	userInface.SetLocked(locked)
}

func (userManager *AccountUserManager) SetPassword(passwd, hint string) {
	userInface := _infaceMap[userManager.ObjectPath]
	if userInface == nil {
		return
	}
	userInface.SetPassword(passwd, hint)
}

func (userManager *AccountUserManager) SetPasswordMode(mode int32) {
	userInface := _infaceMap[userManager.ObjectPath]
	if userInface == nil {
		return
	}
	userInface.SetPasswordMode(mode)
}

func (userManager *AccountUserManager) SetUserName(name string) {
	userInface := _infaceMap[userManager.ObjectPath]
	if userInface == nil {
		return
	}
	userInface.SetUserName(name)
}

func NewAccountUserManager(path string) *AccountUserManager {
	userManager := &AccountUserManager{}

	userManager.ObjectPath = path
	userInface := accounts.GetUser(path)
	_infaceMap[path] = userInface

	GetUserProperties(userManager)
	userInface.ConnectChanged(func() {
		tmpUser := &AccountUserManager{}
		tmpUser.ObjectPath = userManager.ObjectPath
		GetUserProperties(tmpUser)
		CompareUserManager(userManager, tmpUser)
	})

	_userMap[path] = userManager
	return userManager
}

func GetUserProperties(userManager *AccountUserManager) {
	userInface := _infaceMap[userManager.ObjectPath]
	if userInface == nil {
		return
	}
	userManager.AccountType = userInface.AccountType.Get()
	userManager.AutomaticLogin = userInface.AutomaticLogin.Get()
	userManager.IconFile = userInface.IconFile.Get()
	userManager.Locked = userInface.Locked.Get()
	userManager.LoginTime = userInface.LoginTime.Get()
	userManager.PasswordMode = userInface.PasswordMode.Get()
	userManager.UserName = userInface.UserName.Get()
}

func CompareUserManager(src, tmp *AccountUserManager) {
	if src == nil || tmp == nil {
		return
	}

	if src.AccountType != tmp.AccountType {
		src.AccountType = tmp.AccountType
		src.PropertyChanged("AccountType")
	}

	if src.AutomaticLogin != tmp.AutomaticLogin {
		src.AutomaticLogin = tmp.AutomaticLogin
		src.PropertyChanged("AutomaticLogin")
	}

	if src.IconFile != tmp.IconFile {
		src.IconFile = tmp.IconFile
		src.PropertyChanged("IconFile")
	}

	if src.Locked != tmp.Locked {
		src.Locked = tmp.Locked
		src.PropertyChanged("Locked")
	}

	if src.LoginTime != tmp.LoginTime {
		src.LoginTime = tmp.LoginTime
		src.PropertyChanged("LoginTime")
	}

	if src.PasswordMode != tmp.PasswordMode {
		src.PasswordMode = tmp.PasswordMode
		src.PropertyChanged("PasswordMode")
	}

	if src.UserName != tmp.UserName {
		src.UserName = tmp.UserName
		src.PropertyChanged("UserName")
	}
}

func DeleteUserManager(path string) {
	userManager := _userMap[path]
	if userManager != nil {
		return
	}

	dbus.UnInstallObject(userManager)
}
