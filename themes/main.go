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
	"dbus/com/deepin/daemon/accounts"
	accext "dbus/com/deepin/dde/api/accounts"
	"dlib"
	"dlib/dbus"
	"fmt"
	"os/user"
)

var (
	accountsExtends *accext.Accounts
	userManager     *accounts.User

	currentUid string
)

func InitVariable() {
	var err error

	accountsExtends, err = accext.NewAccounts("/com/deepin/dde/api/Accounts")
	if err != nil {
		fmt.Println("New Accounts Extends Failed.")
		panic(err)
	}

	userInfo, _ := user.Current()
	currentUid = userInfo.Uid

	userManager, err = accounts.NewUser(DACCOUNTS_USER_PATH +
		dbus.ObjectPath(currentUid))
	if err != nil {
		fmt.Println("New User Failed.")
		panic(err)
	}
}

func main() {
	InitVariable()
	ReadThemeDir(THEME_DIR)
	m := NewManager()
	err := dbus.InstallOnSession(m)
	if err != nil {
		panic(err)
	}

	if m.AutoSwitch.Get() {
		go m.switchPictureThread()
	}
	dbus.DealWithUnhandledMessage()
	dlib.StartLoop()
}
