/*
 * Copyright (C) 2013 ~ 2018 Deepin Technology Co., Ltd.
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
	"errors"
	"fmt"
	"math/rand"
	"pkg.deepin.io/dde/daemon/accounts/checkers"
	"pkg.deepin.io/dde/daemon/accounts/users"
	"pkg.deepin.io/lib/dbus"
	dutils "pkg.deepin.io/lib/utils"
	"time"
)

const nilObjPath = dbus.ObjectPath("/")

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
	name, fullname string, ty int32) (dbus.ObjectPath, error) {

	logger.Debug("[CreateUser] new user:", name, fullname, ty)
	pid := dbusMsg.GetSenderPID()
	if err := polkitAuthManagerUser(pid); err != nil {
		logger.Debug("[CreateUser] access denied:", err)
		return nilObjPath, err
	}

	ch := make(chan string)
	m.mapLocker.Lock()
	m.userAddedChans[name] = ch
	m.mapLocker.Unlock()
	defer func() {
		m.mapLocker.Lock()
		delete(m.userAddedChans, name)
		m.mapLocker.Unlock()
		close(ch)
	}()

	if err := users.CreateUser(name, fullname, "", ty); err != nil {
		logger.Warningf("DoAction: create user '%s' failed: %v\n",
			name, err)
		return nilObjPath, err
	}

	if err := users.SetUserType(ty, name); err != nil {
		logger.Warningf("DoAction: set user type '%s' failed: %v\n",
			name, err)
		return nilObjPath, err
	}

	// create user success
	select {
	case userPath, ok := <-ch:
		if !ok {
			return nilObjPath, errors.New("invalid user path event")
		}

		logger.Debug("receive user path", userPath)
		if userPath == "" {
			return nilObjPath, errors.New("failed to install user on session bus")
		}
		return dbus.ObjectPath(userPath), nil
	case <-time.After(time.Second * 60):
		err := errors.New("wait timeout exceeded")
		logger.Warning(err)
		return nilObjPath, err
	}
}

// Delete a exist user.
//
// name: 用户名
//
// rmFiles: 是否删除用户数据
func (m *Manager) DeleteUser(dbusMsg dbus.DMessage,
	name string, rmFiles bool) error {

	logger.Debug("[DeleteUser] user:", name, rmFiles)
	pid := dbusMsg.GetSenderPID()
	if err := polkitAuthManagerUser(pid); err != nil {
		logger.Debug("[DeleteUser] access denied:", err)
		return err
	}

	user := m.getUserByName(name)
	if user == nil {
		err := fmt.Errorf("user %q not found", name)
		logger.Warning(err)
		return err
	}

	if err := users.DeleteUser(rmFiles, name); err != nil {
		logger.Warningf("DoAction: delete user '%s' failed: %v\n",
			name, err)
		return err
	}

	if users.IsAutoLoginUser(name) {
		users.SetAutoLoginUser("", "")
	}

	//delete user config and icons
	if rmFiles {
		user.clearData()
	}
	return nil
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
func (m *Manager) RandUserIcon() (string, error) {
	icons := getUserStandardIcons()
	if len(icons) == 0 {
		return "", errors.New("Did not find any user icons")
	}

	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(icons))
	return icons[idx], nil
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

// 检测密码是否有效
//
// ret0: 是否合法
//
// ret1: 提示信息
//
// ret2: 不合法代码
func (m *Manager) IsPasswordValid(passwd string) (bool, string, int32) {
	releaseType := getDeepinReleaseType()
	logger.Infof("release type %q", releaseType)
	errCode := checkers.CheckPasswordValid(releaseType, passwd)
	return errCode.IsOk(), errCode.Prompt(), int32(errCode)
}

func (m *Manager) AllowGuestAccount(dbusMsg dbus.DMessage, allow bool) error {
	pid := dbusMsg.GetSenderPID()
	if err := polkitAuthManagerUser(pid); err != nil {
		return err
	}

	if allow == isGuestUserEnabled() {
		return nil
	}

	success := dutils.WriteKeyToKeyFile(actConfigFile,
		actConfigGroupGroup, actConfigKeyGuest, allow)
	if !success {
		return errors.New("Enable guest user failed")
	}
	m.setPropAllowGuest(allow)

	return nil
}

func (m *Manager) CreateGuestAccount(dbusMsg dbus.DMessage) (string, error) {
	pid := dbusMsg.GetSenderPID()
	if err := polkitAuthManagerUser(pid); err != nil {
		return "", err
	}

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
