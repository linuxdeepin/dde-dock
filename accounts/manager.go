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
	"sort"
	"sync"

	"pkg.deepin.io/dde/daemon/accounts/users"
	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/tasker"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	actConfigDir       = "/var/lib/AccountsService"
	userConfigDir      = actConfigDir + "/deepin/users"
	userIconsDir       = actConfigDir + "/icons"
	userCustomIconsDir = actConfigDir + "/icons/local"

	userIconGuest       = actConfigDir + "/icons/guest.png"
	actConfigFile       = actConfigDir + "/accounts.ini"
	actConfigGroupGroup = "Accounts"
	actConfigKeyGuest   = "AllowGuest"
)

//go:generate dbusutil-gen -type Manager,User manager.go user.go

type Manager struct {
	service *dbusutil.Service
	PropsMu sync.RWMutex

	UserList   []string
	UserListMu sync.RWMutex

	// dbusutil-gen: ignore
	GuestIcon  string
	AllowGuest bool

	watcher    *dutils.WatchProxy
	usersMap   map[string]*User
	usersMapMu sync.Mutex

	delayTaskManager *tasker.DelayTaskManager

	userAddedChanMap map[string]chan string
	//                    ^ username
	//nolint
	signals *struct {
		UserAdded struct {
			objPath string
		}

		UserDeleted struct {
			objPath string
		}
	}
	//nolint
	methods *struct {
		CreateUser         func() `in:"name,fullName,accountType" out:"user"`
		DeleteUser         func() `in:"name,rmFiles"`
		FindUserById       func() `in:"uid" out:"user"`
		FindUserByName     func() `in:"name" out:"user"`
		RandUserIcon       func() `out:"iconFile"`
		IsUsernameValid    func() `in:"name" out:"ok,errReason,errCode"`
		IsPasswordValid    func() `in:"password" out:"ok,errReason,errCode"`
		AllowGuestAccount  func() `in:"allow"`
		CreateGuestAccount func() `out:"user"`
		GetGroups          func() `out:"groups"`
		GetPresetGroups    func() `in:"accountType" out:"groups"`
	}
}

func NewManager(service *dbusutil.Service) *Manager {
	var m = &Manager{
		service: service,
	}

	m.usersMap = make(map[string]*User)
	m.userAddedChanMap = make(map[string]chan string)

	m.GuestIcon = userIconGuest
	m.AllowGuest = isGuestUserEnabled()
	m.initUsers(getUserPaths())

	m.watcher = dutils.NewWatchProxy()
	if m.watcher != nil {
		m.delayTaskManager = tasker.NewDelayTaskManager()
		_ = m.delayTaskManager.AddTask(taskNamePasswd, fileEventDelay, m.handleFilePasswdChanged)
		_ = m.delayTaskManager.AddTask(taskNameGroup, fileEventDelay, m.handleFileGroupChanged)
		_ = m.delayTaskManager.AddTask(taskNameShadow, fileEventDelay, m.handleFileShadowChanged)
		_ = m.delayTaskManager.AddTask(taskNameDM, fileEventDelay, m.handleDMConfigChanged)

		m.watcher.SetFileList(m.getWatchFiles())
		m.watcher.SetEventHandler(m.handleFileChanged)
		go m.watcher.StartWatch()
	}

	return m
}

func (m *Manager) destroy() {
	if m.watcher != nil {
		m.watcher.EndWatch()
		m.watcher = nil
	}

	m.stopExportUsers(m.UserList)
	_ = m.service.StopExport(m)
}

func (m *Manager) initUsers(list []string) {
	var userList []string
	for _, p := range list {
		u, err := NewUser(p, m.service)
		if err != nil {
			logger.Errorf("New user '%s' failed: %v", p, err)
			continue
		}

		userList = append(userList, p)

		m.usersMapMu.Lock()
		m.usersMap[p] = u
		m.usersMapMu.Unlock()
	}
	sort.Strings(userList)
	m.UserList = userList
}

func (m *Manager) exportUsers() {
	m.usersMapMu.Lock()

	for _, u := range m.usersMap {
		err := m.service.Export(dbus.ObjectPath(userDBusPathPrefix+u.Uid), u)
		if err != nil {
			logger.Errorf("failed to export user %q: %v",
				u.Uid, err)
			continue
		}

	}

	m.usersMapMu.Unlock()
}

func (m *Manager) stopExportUsers(list []string) {
	for _, p := range list {
		m.stopExportUser(p)
	}
}

func (m *Manager) exportUserByPath(userPath string) error {
	u, err := NewUser(userPath, m.service)
	if err != nil {
		return err
	}

	m.usersMapMu.Lock()
	ch := m.userAddedChanMap[u.UserName]
	m.usersMapMu.Unlock()

	err = m.service.Export(dbus.ObjectPath(userPath), u)
	logger.Debugf("export user %q err: %v", userPath, err)
	if ch != nil {
		if err != nil {
			ch <- ""
		} else {
			ch <- userPath
		}
		logger.Debug("after ch <- userPath")
	}

	if err != nil {
		return err
	}

	m.usersMapMu.Lock()
	m.usersMap[userPath] = u
	m.usersMapMu.Unlock()

	return nil
}

func (m *Manager) stopExportUser(userPath string) {
	m.usersMapMu.Lock()
	defer m.usersMapMu.Unlock()
	u, ok := m.usersMap[userPath]
	if !ok {
		logger.Debug("Invalid user path:", userPath)
		return
	}

	delete(m.usersMap, userPath)
	_ = m.service.StopExport(u)
}

func (m *Manager) getUserByName(name string) *User {
	m.usersMapMu.Lock()
	defer m.usersMapMu.Unlock()

	for _, user := range m.usersMap {
		if user.UserName == name {
			return user
		}
	}
	return nil
}

func (m *Manager) getUserByUid(uid string) *User {
	m.usersMapMu.Lock()
	defer m.usersMapMu.Unlock()

	for _, user := range m.usersMap {
		if user.Uid == uid {
			return user
		}
	}
	return nil
}

func getUserPaths() []string {
	infos, err := users.GetHumanUserInfos()
	if err != nil {
		return nil
	}

	var paths []string
	for _, info := range infos {
		paths = append(paths, userDBusPathPrefix+info.Uid)
	}

	return paths
}

func isGuestUserEnabled() bool {
	v, exist := dutils.ReadKeyFromKeyFile(actConfigFile,
		actConfigGroupGroup, actConfigKeyGuest, true)
	if !exist {
		return false
	}

	ret, ok := v.(bool)
	if !ok {
		return false
	}

	return ret
}

func (m *Manager) checkAuth(sender dbus.Sender) error {
	return checkAuth(polkitActionUserAdministration, string(sender))
}
