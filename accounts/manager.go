/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package accounts

import (
	"pkg.deepin.io/dde/daemon/accounts/users"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/tasker"
	dutils "pkg.deepin.io/lib/utils"
	"sync"
)

const (
	userConfigDir      = "/var/lib/AccountsService/users"
	userIconsDir       = "/var/lib/AccountsService/icons"
	userCustomIconsDir = "/var/lib/AccountsService/icons/local"

	userIconGuest       = "/var/lib/AccountsService/icons/guest.png"
	actConfigFile       = "/var/lib/AccountsService/accounts.ini"
	actConfigGroupGroup = "Accounts"
	actConfigKeyGuest   = "AllowGuest"
)

type Manager struct {
	// 用户 ObjectPath 列表
	UserList      []string
	userListMutex sync.Mutex
	GuestIcon     string
	AllowGuest    bool

	watcher   *dutils.WatchProxy
	usersMap  map[string]*User
	mapLocker sync.Mutex

	delayTasker *tasker.DelayTaskManager

	// Signals:
	UserAdded   func(string)
	UserDeleted func(string)

	userAddedChans map[string]chan string // key: user name
}

func NewManager() *Manager {
	var m = &Manager{}

	m.usersMap = make(map[string]*User)
	m.userAddedChans = make(map[string]chan string)

	m.setPropGuestIcon(userIconGuest)
	m.setPropAllowGuest(isGuestUserEnabled())
	m.newUsers(getUserPaths())

	m.watcher = dutils.NewWatchProxy()
	if m.watcher != nil {
		m.delayTasker = tasker.NewDelayTaskManager()
		m.delayTasker.AddTask(taskNamePasswd, fileEventDelay, m.handleFilePasswdChanged)
		m.delayTasker.AddTask(taskNameGroup, fileEventDelay, m.handleFileGroupChanged)
		m.delayTasker.AddTask(taskNameShadow, fileEventDelay, m.handleFileShadowChanged)
		m.delayTasker.AddTask(taskNameDM, fileEventDelay, m.handleDMConfigChanged)

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

	m.uninstallUsers(m.UserList)
	dbus.UnInstallObject(m)
}

func (m *Manager) newUsers(list []string) {
	var paths []string
	for _, p := range list {
		u, err := NewUser(p)
		if err != nil {
			logger.Errorf("New user '%s' failed: %v", p, err)
			continue
		}

		paths = append(paths, p)

		m.mapLocker.Lock()
		m.usersMap[p] = u
		m.mapLocker.Unlock()
	}
	m.setPropUserList(paths)
}

func (m *Manager) installUsers() {
	m.mapLocker.Lock()
	defer m.mapLocker.Unlock()
	for _, u := range m.usersMap {
		err := dbus.InstallOnSystem(u)
		if err != nil {
			logger.Errorf("Install user '%s' failed: %v",
				u.Uid, err)
			continue
		}
	}
}

func (m *Manager) uninstallUsers(list []string) {
	for _, p := range list {
		m.uninstallUser(p)
	}
}

func (m *Manager) installUserByPath(userPath string) error {
	u, err := NewUser(userPath)
	if err != nil {
		return err
	}

	m.mapLocker.Lock()
	ch := m.userAddedChans[u.UserName]
	m.mapLocker.Unlock()

	err = dbus.InstallOnSystem(u)
	logger.Debugf("install user %q err: %v", userPath, err)
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

	m.mapLocker.Lock()
	m.usersMap[userPath] = u
	m.mapLocker.Unlock()

	return nil
}

func (m *Manager) uninstallUser(userPath string) {
	m.mapLocker.Lock()
	defer m.mapLocker.Unlock()
	u, ok := m.usersMap[userPath]
	if !ok {
		logger.Debug("Invalid user path:", userPath)
		return
	}

	delete(m.usersMap, userPath)
	u.destroy()
}

func (m *Manager) getUserByName(name string) *User {
	m.mapLocker.Lock()
	defer m.mapLocker.Unlock()

	for _, user := range m.usersMap {
		if user.UserName == name {
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
		paths = append(paths, userDBusPath+info.Uid)
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
