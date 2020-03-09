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
	"time"

	"pkg.deepin.io/dde/daemon/accounts/checkers"
	"pkg.deepin.io/dde/daemon/accounts/users"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/procfs"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	nilObjPath      = dbus.ObjectPath("/")
	dbusServiceName = "com.deepin.daemon.Accounts"
	dbusPath        = "/com/deepin/daemon/Accounts"
	dbusInterface   = "com.deepin.daemon.Accounts"
)

func (*Manager) GetInterfaceName() string {
	return dbusInterface
}

// Create new user.
//
// 如果收到 Error 信号，则创建失败。
//
// name: 用户名
//
// fullName: 全名，可以为空
//
// ty: 用户类型，0 为普通用户，1 为管理员

func (m *Manager) CreateUser(sender dbus.Sender,
	name, fullName string, accountType int32) (dbus.ObjectPath, *dbus.Error) {

	logger.Debug("[CreateUser] new user:", name, fullName, accountType)

	err := checkAccountType(int(accountType))
	if err != nil {
		return nilObjPath, dbusutil.ToError(err)
	}

	err = m.checkAuth(sender)
	if err != nil {
		logger.Debug("[CreateUser] access denied:", err)
		return nilObjPath, dbusutil.ToError(err)
	}

	ch := make(chan string)
	m.usersMapMu.Lock()
	m.userAddedChanMap[name] = ch
	m.usersMapMu.Unlock()
	defer func() {
		m.usersMapMu.Lock()
		delete(m.userAddedChanMap, name)
		m.usersMapMu.Unlock()
		close(ch)
	}()

	if err := users.CreateUser(name, fullName, ""); err != nil {
		logger.Warningf("DoAction: create user '%s' failed: %v\n",
			name, err)
		return nilObjPath, dbusutil.ToError(err)
	}

	groups := users.GetPresetGroups(int(accountType))
	logger.Debug("groups:", groups)
	err = users.SetGroupsForUser(groups, name)
	if err != nil {
		logger.Warningf("failed to set groups for user %s: %v", name, err)
	}

	// create user success
	select {
	case userPath, ok := <-ch:
		if !ok {
			return nilObjPath, dbusutil.ToError(errors.New("invalid user path event"))
		}

		logger.Debug("receive user path", userPath)
		if userPath == "" {
			return nilObjPath, dbusutil.ToError(errors.New("failed to install user on session bus"))
		}
		return dbus.ObjectPath(userPath), nil
	case <-time.After(time.Second * 60):
		err := errors.New("wait timeout exceeded")
		logger.Warning(err)
		return nilObjPath, dbusutil.ToError(err)
	}
}

// Delete a exist user.
//
// name: 用户名
//
// rmFiles: 是否删除用户数据
func (m *Manager) DeleteUser(sender dbus.Sender,
	name string, rmFiles bool) *dbus.Error {

	logger.Debug("[DeleteUser] user:", name, rmFiles)

	err := m.checkAuth(sender)
	if err != nil {
		logger.Debug("[DeleteUser] access denied:", err)
		return dbusutil.ToError(err)
	}

	user := m.getUserByName(name)
	if user == nil {
		err := fmt.Errorf("user %q not found", name)
		logger.Warning(err)
		return dbusutil.ToError(err)
	}

	if err := users.DeleteUser(rmFiles, name); err != nil {
		logger.Warningf("DoAction: delete user '%s' failed: %v\n",
			name, err)
		return dbusutil.ToError(err)
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

func (m *Manager) FindUserById(uid string) (string, *dbus.Error) {
	userPath := userDBusPathPrefix + uid
	for _, v := range m.UserList {
		if v == userPath {
			return v, nil
		}
	}

	return "", dbusutil.ToError(fmt.Errorf("Invalid uid: %s", uid))
}

func (m *Manager) FindUserByName(name string) (string, *dbus.Error) {
	info, err := users.GetUserInfoByName(name)
	if err != nil {
		return "", dbusutil.ToError(err)
	}

	return m.FindUserById(info.Uid)
}

// 随机得到一个用户头像
//
// ret0：头像路径，为空则表示获取失败
func (m *Manager) RandUserIcon() (string, *dbus.Error) {
	icons := getUserStandardIcons()
	if len(icons) == 0 {
		return "", dbusutil.ToError(errors.New("Did not find any user icons"))
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
func (m *Manager) IsUsernameValid(sender dbus.Sender, name string) (valid bool,
	msg string, code int32, busErr *dbus.Error) {

	var err error
	defer func() {
		busErr = dbusutil.ToError(err)
	}()

	pid, err := m.service.GetConnPID(string(sender))
	if err != nil {
		return
	}

	p := procfs.Process(pid)
	environ, err := p.Environ()
	if err != nil {
		return
	}

	locale := environ.Get("LANG")

	info := checkers.CheckUsernameValid(name)
	if info == nil {
		valid = true
		return
	}

	msg = info.Error.Error()
	logger.Debug("locale:", locale)
	if locale != "" {
		gettext.SetLocale(gettext.LcAll, locale)
		msg = gettext.Tr(msg)
	}
	code = int32(info.Code)
	return
}

// 检测密码是否有效
//
// ret0: 是否合法
//
// ret1: 提示信息
//
// ret2: 不合法代码
func (m *Manager) IsPasswordValid(password string) (bool, string, int32, *dbus.Error) {
	releaseType := getDeepinReleaseType()
	logger.Infof("release type %q", releaseType)
	errCode := checkers.CheckPasswordValid(releaseType, password)
	return errCode.IsOk(), errCode.Prompt(), int32(errCode), nil
}

func (m *Manager) AllowGuestAccount(sender dbus.Sender, allow bool) *dbus.Error {
	err := m.checkAuth(sender)
	if err != nil {
		return dbusutil.ToError(err)
	}

	m.PropsMu.Lock()
	defer m.PropsMu.Unlock()

	if m.AllowGuest == allow {
		return nil
	}

	success := dutils.WriteKeyToKeyFile(actConfigFile,
		actConfigGroupGroup, actConfigKeyGuest, allow)
	if !success {
		return dbusutil.ToError(errors.New("enable guest user failed"))
	}

	m.AllowGuest = allow
	m.emitPropChangedAllowGuest(allow)
	return nil
}

func (m *Manager) CreateGuestAccount(sender dbus.Sender) (string, *dbus.Error) {
	err := m.checkAuth(sender)
	if err != nil {
		return "", dbusutil.ToError(err)
	}

	name, err := users.CreateGuestUser()
	if err != nil {
		return "", dbusutil.ToError(err)
	}

	info, err := users.GetUserInfoByName(name)
	if err != nil {
		return "", dbusutil.ToError(err)
	}

	return userDBusPathPrefix + info.Uid, nil
}

func (m *Manager) GetGroups() ([]string, *dbus.Error) {
	groups, err := users.GetAllGroups()
	return groups, dbusutil.ToError(err)
}

func (m *Manager) GetPresetGroups(accountType int32) ([]string, *dbus.Error) {
	err := checkAccountType(int(accountType))
	if err != nil {
		return nil, dbusutil.ToError(err)
	}

	groups := users.GetPresetGroups(int(accountType))
	return groups, nil
}
