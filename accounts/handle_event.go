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
	"time"

	"pkg.deepin.io/dde/daemon/accounts/users"
	"pkg.deepin.io/lib/fsnotify"
)

const (
	userFilePasswd  = "/etc/passwd"
	userFileGroup   = "/etc/group"
	userFileShadow  = "/etc/shadow"
	userFileSudoers = "/etc/sudoers"

	lightdmConfig = "/etc/lightdm/lightdm.conf"
	kdmConfig     = "/usr/share/config/kdm/kdmrc"
	gdmConfig     = "/etc/gdm/custom.conf"
)

const (
	maxDuration    = time.Second * 1
	deltaDuration  = time.Millisecond * 500
	fileEventDelay = time.Millisecond * 500

	taskNamePasswd = "passwd"
	taskNameGroup  = "group"
	taskNameShadow = "shadow"
	taskNameDM     = "dm"
)

const (
	userListNotChange = iota + 1
	userListAdded
	userListDeleted
)

func (m *Manager) getWatchFiles() []string {
	var list []string
	dmConfig, err := users.GetDMConfig()
	if err == nil {
		list = append(list, dmConfig)
	}

	list = append(list, []string{userFilePasswd, userFileGroup,
		userFileShadow, userFileSudoers}...)
	return list
}

func (m *Manager) handleFileChanged(ev *fsnotify.FileEvent) {
	if ev == nil {
		return
	}

	logger.Debug("File changed:", ev)
	var err error
	switch ev.Name {
	case userFilePasswd:
		if task, _ := m.delayTaskManager.GetTask(taskNamePasswd); task != nil {
			err = task.Start()
		}
	case userFileGroup, userFileSudoers:
		if task, _ := m.delayTaskManager.GetTask(taskNameGroup); task != nil {
			err = task.Start()
		}
	case userFileShadow:
		if task, _ := m.delayTaskManager.GetTask(taskNameShadow); task != nil {
			err = task.Start()
		}
	case lightdmConfig, kdmConfig, gdmConfig:
		if task, _ := m.delayTaskManager.GetTask(taskNameDM); task != nil {
			err = task.Start()
		}
	default:
		logger.Debug("Unknow event, ignore")
		return
	}
	if err != nil {
		logger.Warning("Failed to start task:", err, ev)
	}
	m.watcher.ResetFileListWatch()
}

func (m *Manager) handleFilePasswdChanged() {
	waitDuration := time.Second * 0
	for waitDuration <= maxDuration {
		if m.refreshUserList() {
			break
		}

		waitDuration += deltaDuration
		<-time.After(waitDuration)
	}
}

func (m *Manager) handleFileGroupChanged() {
	m.usersMapMu.Lock()
	defer m.usersMapMu.Unlock()
	for _, u := range m.usersMap {
		u.updatePropAccountType()
		u.updatePropCanNoPasswdLogin()
		u.updatePropGroups()
	}
}

func (m *Manager) handleFileShadowChanged() {
	//Update the property 'PasswordStatus' and 'Locked'
	m.usersMapMu.Lock()
	defer m.usersMapMu.Unlock()
	for _, u := range m.usersMap {
		u.updatePropPasswordStatus()
	}
}

func (m *Manager) handleDMConfigChanged() {
	for _, u := range m.usersMap {
		u.updatePropAutomaticLogin()
	}
}

func (m *Manager) refreshUserList() bool {
	m.UserListMu.Lock()

	var freshed bool
	ret, status := compareUserList(m.UserList, getUserPaths())
	switch status {
	case userListAdded:
		freshed = true
		m.handleUserAdded(ret)
	case userListDeleted:
		freshed = true
		m.handleUserDeleted(ret)
	}

	defer m.UserListMu.Unlock()
	return freshed
}

func (m *Manager) handleUserAdded(list []string) {
	var userList = m.UserList
	for _, p := range list {
		err := m.exportUserByPath(p)
		if err != nil {
			logger.Errorf("Install user '%s' failed: %v", p, err)
			continue
		}

		userList = append(userList, p)
		m.service.Emit(m, "UserAdded", p)
		m.copyUserDatas(p)
	}

	m.UserList = userList
	m.service.EmitPropertyChanged(m, "UserList", userList)
}

func (m *Manager) handleUserDeleted(list []string) {
	var userList = m.UserList
	for _, p := range list {
		m.stopExportUser(p)
		userList = deleteStrFromList(p, userList)
		m.service.Emit(m, "UserDeleted", p)
	}

	m.UserList = userList
	m.service.EmitPropertyChanged(m, "UserList", userList)
}

func compareUserList(oldList, newList []string) ([]string, int) {
	var (
		ret    []string
		oldLen = len(oldList)
		newLen = len(newList)
	)

	if oldLen < newLen {
		for _, v := range newList {
			if isStrInArray(v, oldList) {
				continue
			}
			ret = append(ret, v)
		}
		return ret, userListAdded
	} else if oldLen > newLen {
		for _, v := range oldList {
			if isStrInArray(v, newList) {
				continue
			}
			ret = append(ret, v)
		}
		return ret, userListDeleted
	}

	return ret, userListNotChange
}

func deleteStrFromList(str string, list []string) []string {
	var ret []string
	for _, v := range list {
		if v == str {
			continue
		}
		ret = append(ret, v)
	}

	return ret
}
