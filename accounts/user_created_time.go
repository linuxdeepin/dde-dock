/*
 * Copyright (C) 2019 ~ 2019 Deepin Technology Co., Ltd.
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

package accounts

import (
	"path/filepath"
	"syscall"
)

func (u *User) getCreatedTime() (uint64, error) {
	createdTime, err := getCreatedTimeFromBashLogout(u.HomeDir)
	if err == nil {
		return createdTime, nil
	}
	return getCreatedTimeFromUserConfig(u.UserName)
}

// '.bash_logout' from '/etc/skel/.bash_logout'
func getCreatedTimeFromBashLogout(home string) (uint64, error) {
	return getCreatedTimeFromFile(filepath.Join(home, ".bash_logout"))
}

// using deepin accounts dbus created user, will created this configuration file
func getCreatedTimeFromUserConfig(username string) (uint64, error) {
	return getCreatedTimeFromFile(filepath.Join("/var/lib/AccountsService/deepin/users",
		username))
}

func getCreatedTimeFromFile(filename string) (uint64, error) {
	var stat syscall.Stat_t
	err := syscall.Stat(filename, &stat)
	if err != nil {
		return 0, err
	}
	// Ctim: recent change time
	return uint64(stat.Ctim.Sec), nil
}
