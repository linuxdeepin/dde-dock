/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
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

package accounts

import (
	"pkg.linuxdeepin.com/lib/dbus"
	dutils "pkg.linuxdeepin.com/lib/utils"
)

func (obj *Manager) CreateGuestAccount() string {
	args := []string{}

	username := getGuestName()
	passwd := encodePasswd("")
	args = append(args, "-m")
	args = append(args, "-d")
	args = append(args, "/tmp/"+username)
	args = append(args, "-s")
	args = append(args, "/bin/bash")
	args = append(args, "-l")
	args = append(args, "-p")
	args = append(args, passwd)
	args = append(args, username)
	if !execCommand(CMD_USERADD, args) {
		return ""
	}

	info, _ := getUserInfoByName(username)

	return info.Path
}

func (obj *Manager) AllowGuestAccount(dbusMsg dbus.DMessage, allow bool) bool {
	if ok := polkitAuthWithPid(POLKIT_MANAGER_USER,
		dbusMsg.GetSenderPID()); !ok {
		return false
	}

	logger.Infof("Allow guest: %v", allow)
	if ok := dutils.WriteKeyToKeyFile(ACCOUNT_CONFIG_FILE,
		ACCOUNT_GROUP_KEY, ACCOUNT_KEY_GUEST, allow); !ok {
		logger.Error("AllowGuest Failed")
		return false
	}

	return true
}

func (obj *Manager) CreateUser(dbusMsg dbus.DMessage, name, fullname string, accountTyte int32) (string, bool) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Recover Error In CreateUser:%v",
				err)
		}
	}()
	if ok := polkitAuthWithPid(POLKIT_MANAGER_USER,
		dbusMsg.GetSenderPID()); !ok {
		return "", false
	}

	args := []string{}

	args = append(args, "-m")
	args = append(args, "-s")
	args = append(args, "/bin/bash")
	args = append(args, "-c")
	args = append(args, fullname)
	args = append(args, name)
	if !execCommand(CMD_USERADD, args) {
		return "", false
	}

	info, _ := getUserInfoByName(name)
	if u, ok := obj.pathUserMap[info.Path]; ok {
		u.setPropAccountType(accountTyte)
	}

	changeFileOwner(name, name, "/home/"+name)

	return info.Path, true
}

func (obj *Manager) DeleteUser(dbusMsg dbus.DMessage, name string, removeFiles bool) bool {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Recover Error In DeleteUser:%v",
				err)
		}
	}()

	//if ok := opUtils.PolkitAuthWithPid(POLKIT_MANAGER_USER,
	if ok := polkitAuthWithPid(POLKIT_MANAGER_USER,
		dbusMsg.GetSenderPID()); !ok {
		return false
	}

	args := []string{}
	user, ok := obj.pathUserMap[obj.FindUserByName(name)]
	if ok {
		if user.AutomaticLogin {
			user.SetAutomaticLogin(dbusMsg, false)
		}
	}

	if removeFiles {
		args = append(args, "-r")
	}
	args = append(args, name)

	if !execCommand(CMD_USERDEL, args) {
		return false
	}

	return true
}

func (obj *Manager) FindUserById(id string) string {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Recover Error In FindUserById:%v",
				err)
		}
	}()

	path := USER_MANAGER_PATH + id

	for _, v := range obj.UserList {
		if path == v {
			return path
		}
	}

	return ""
}

func (obj *Manager) FindUserByName(name string) string {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("Recover Error In FindUserByName:%v", err)
		}
	}()

	if info, ok := getUserInfoByName(name); ok {
		return info.Path
	}

	return ""
}

func (obj *Manager) RandUserIcon() (string, bool) {
	if icon := getRandUserIcon(); len(icon) > 0 {
		return icon, true
	}

	return "", false
}

func (m *Manager) IsUsernameValid(username string) (bool, string) {
	if !isUsernameValid(username) {
		return false, "The user name is not valid."
	}

	if !isUserExist(username) {
		return false, "The user name already exists."
	}

	return true, ""
}

func (m *Manager) IsPasswordValid(passwd string) bool {
	return isPasswordValid(passwd)
}
