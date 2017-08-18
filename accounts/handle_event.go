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
	"pkg.deepin.io/lib/fsnotify"
	"pkg.deepin.io/dde/daemon/accounts/users"
	"pkg.deepin.io/lib/dbus"
	"strings"
	"time"
)

const (
	userFilePasswd  = "/etc/passwd"
	userFileGroup   = "/etc/group"
	userFileShadow  = "/etc/shadow"
	userFileSudoers = "/etc/sudoers"

	lightdmConfig = "/etc/lightdm/lightdm.conf"
	kdmConfig     = "/usr/share/config/kdm/kdmrc"
	gdmConfig     = "/etc/gdm/custom.conf"

	permModeDir = 0755
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
	userListNotChange int = iota + 1
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
	switch {
	case strings.Contains(ev.Name, userFilePasswd):
		if task, _ := m.delayTasker.GetTask(taskNamePasswd); task != nil {
			err = task.Start()
		}
	case strings.Contains(ev.Name, userFileGroup), strings.Contains(ev.Name, userFileSudoers):
		if task, _ := m.delayTasker.GetTask(taskNameGroup); task != nil {
			err = task.Start()
		}
	case strings.Contains(ev.Name, userFileShadow):
		if task, _ := m.delayTasker.GetTask(taskNameShadow); task != nil {
			err = task.Start()
		}
	case strings.Contains(ev.Name, lightdmConfig),
		strings.Contains(ev.Name, kdmConfig),
		strings.Contains(ev.Name, gdmConfig):
		if task, _ := m.delayTasker.GetTask(taskNameDM); task != nil {
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
	m.mapLocker.Lock()
	defer m.mapLocker.Unlock()
	for _, u := range m.usersMap {
		u.updatePropAccountType()
	}
}

func (m *Manager) handleFileShadowChanged() {
	//Update the property 'Locked'
	m.mapLocker.Lock()
	defer m.mapLocker.Unlock()
	for _, u := range m.usersMap {
		u.updatePropLocked()
	}
}

func (m *Manager) handleDMConfigChanged() {
	for _, u := range m.usersMap {
		u.setPropBool(&u.AutomaticLogin, "AutomaticLogin",
			users.IsAutoLoginUser(u.UserName))
	}
}

func (m *Manager) refreshUserList() bool {
	m.userListMutex.Lock()
	defer m.userListMutex.Unlock()

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
	return freshed
}

func (m *Manager) handleUserAdded(list []string) {
	var paths = m.UserList
	for _, p := range list {
		err := m.installUserByPath(p)
		if err != nil {
			logger.Errorf("Install user '%s' failed: %v", p, err)
			continue
		}

		paths = append(paths, p)
		dbus.Emit(m, "UserAdded", p)
		m.copyUserDatas(p)
	}

	m.setPropUserList(paths)
}

func (m *Manager) handleUserDeleted(list []string) {
	var paths = m.UserList
	for _, p := range list {
		m.uninstallUser(p)
		paths = deleteStrFromList(p, paths)
		dbus.Emit(m, "UserDeleted", p)
	}

	m.setPropUserList(paths)
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
