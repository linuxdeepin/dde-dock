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
	"pkg.linuxdeepin.com/dde-daemon/accounts/users"
	"pkg.linuxdeepin.com/lib/dbus"
	dutils "pkg.linuxdeepin.com/lib/utils"
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
	UserList   []string
	GuestIcon  string
	AllowGuest bool

	UserAdded   func(string)
	UserDeleted func(string)
	// Error(pid, action, reason)
	Error func(uint32, string, string)

	watcher   *dutils.WatchProxy
	usersMap  map[string]*User
	mapLocker sync.Mutex
}

func NewManager() *Manager {
	var m = &Manager{}

	m.usersMap = make(map[string]*User)

	m.setPropGuestIcon(userIconGuest)
	m.setPropAllowGuest(isGuestUserEnabled())

	m.watcher = dutils.NewWatchProxy()
	if m.watcher != nil {
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

func (m *Manager) installUsers(list []string) {
	var paths []string
	for _, v := range list {
		err := m.installUser(v)
		if err != nil {
			logger.Errorf("Install user '%s' failed: %v", v, err)
			continue
		}

		paths = append(paths, v)
	}
	m.setPropUserList(paths)
}

func (m *Manager) uninstallUsers(list []string) {
	for _, p := range list {
		m.uninstallUser(p)
	}
}

func (m *Manager) installUser(userPath string) error {
	u, err := NewUser(userPath)
	if err != nil {
		return err
	}

	err = dbus.InstallOnSystem(u)
	if err != nil {
		return err
	}

	m.mapLocker.Lock()
	defer m.mapLocker.Unlock()
	m.usersMap[userPath] = u

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

func (m *Manager) polkitAuthManagerUser(pid uint32, action string) error {
	err := polkitAuthManagerUser(pid)
	if err != nil {
		triggerSigErr(pid, action, err.Error())
		return err
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
