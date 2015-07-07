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
	"fmt"
	"math/rand"
	"pkg.deepin.io/dde-daemon/accounts/checkers"
	"pkg.deepin.io/dde-daemon/accounts/users"
	"pkg.deepin.io/lib/dbus"
	dutils "pkg.deepin.io/lib/utils"
	"time"
)

// Create new user.
//
// 如果收到 Error 信号，则创建失败。
//
// name: 用户名
//
// fullname: 全名，可以为空
//
// ty: 用户类型，0 为普通用户，1 为管理员
func (m *Manager) CreateUser(dbusMsg dbus.DMessage,
	name, fullname string, ty int32) error {
	pid := dbusMsg.GetSenderPID()
	err := m.polkitAuthManagerUser(pid, "CreateUser")
	if err != nil {
		return err
	}

	// Avoid dde-control-center UI block
	go func() {
		err := users.CreateUser(name, fullname, "", ty)
		if err != nil {
			logger.Warningf("DoAction: create user '%s' failed: %v\n",
				name, err)
			triggerSigErr(pid, "CreateUser", err.Error())
			return
		}

		err = users.SetUserType(ty, name)
		if err != nil {
			logger.Warningf("DoAction: set user type '%s' failed: %v\n",
				name, err)
		}
	}()

	return nil
}

// Delete a exist user.
//
// name: 用户名
//
// rmFiles: 是否删除用户数据
func (m *Manager) DeleteUser(dbusMsg dbus.DMessage,
	name string, rmFiles bool) (bool, error) {
	pid := dbusMsg.GetSenderPID()
	err := m.polkitAuthManagerUser(pid, "DeleteUser")
	if err != nil {
		return false, err
	}

	go func() {
		err := users.DeleteUser(rmFiles, name)
		if err != nil {
			logger.Warningf("DoAction: delete user '%s' failed: %v\n",
				name, err)
			triggerSigErr(pid, "DeleteUser", err.Error())
			return
		}

		//delete user config and icons
		if !rmFiles {
			return
		}
		clearUserDatas(name)
	}()

	return true, nil
}

func (m *Manager) FindUserById(uid string) (string, error) {
	userPath := userDBusPath + uid
	for _, v := range m.UserList {
		if v == userPath {
			return v, nil
		}
	}

	return "", fmt.Errorf("Invalid uid: %s", uid)
}

func (m *Manager) FindUserByName(name string) (string, error) {
	info, err := users.GetUserInfoByName(name)
	if err != nil {
		return "", err
	}

	return m.FindUserById(info.Uid)
}

// 随机得到一个用户头像
//
// ret0：头像路径，为空则表示获取失败
//
// ret1: 是否获取成功
func (m *Manager) RandUserIcon() (string, bool, error) {
	icons := getUserStandardIcons()
	if len(icons) == 0 {
		return "", false, fmt.Errorf("Did not find any user icons")
	}

	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(icons))
	return icons[idx], true, nil
}

// 检查用户名是否有效
//
// ret0: 是否合法
//
// ret1: 不合法原因
//
// ret2: 不合法代码
func (m *Manager) IsUsernameValid(name string) (bool, string, int32) {
	info := checkers.CheckUsernameValid(name)
	if info == nil {
		return true, "", 0
	}

	return false, info.Error.Error(), int32(info.Code)
}

func (m *Manager) IsPasswordValid(passwd string) bool {
	return true
}

func (m *Manager) AllowGuestAccount(dbusMsg dbus.DMessage, allow bool) error {
	pid := dbusMsg.GetSenderPID()
	err := m.polkitAuthManagerUser(pid, "AllowGuestAccount")
	if err != nil {
		return err
	}

	if allow == isGuestUserEnabled() {
		return nil
	}

	success := dutils.WriteKeyToKeyFile(actConfigFile,
		actConfigGroupGroup, actConfigKeyGuest, allow)
	if !success {
		reason := "Enable guest user failed"
		triggerSigErr(pid, "AllowGuestAccount", reason)
		return fmt.Errorf(reason)
	}
	m.setPropAllowGuest(allow)

	return nil
}

func (m *Manager) CreateGuestAccount() (string, error) {
	name, err := users.CreateGuestUser()
	if err != nil {
		return "", err
	}

	info, err := users.GetUserInfoByName(name)
	if err != nil {
		return "", err
	}

	return userDBusPath + info.Uid, nil
}
