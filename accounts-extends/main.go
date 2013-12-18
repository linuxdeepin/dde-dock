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
	"fmt"
	"os"
	"os/user"
)

const (
	_BLUR_PICT_DEST = "com.deepin.Accounts"
	_BLUR_PICT_PATH = "/com/deepin/Accounts"
	_BLUR_PICT_IFC  = "com.deepin.Accounts"
)

type AccountExtendsManager struct {
	BlurPictChanged func(string, string)
}

func (blur *AccountExtendsManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_BLUR_PICT_DEST,
		_BLUR_PICT_PATH,
		_BLUR_PICT_IFC,
	}
}

func GetHomeDirById(uid string) (string, error) {
	userInfo, err := user.LookupId(uid)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return userInfo.HomeDir, nil
}

func FileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("File '%s' not exist:%s\n", filename, err)
		return false
	}

	return true
}

func main() {
	accountExt := &AccountExtendsManager{}
	err := dbus.InstallOnSystem(accountExt)
	if err != nil {
		panic(err)
	}

	select {}
}
