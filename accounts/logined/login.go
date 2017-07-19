/**
 * Copyright (C) 2017 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package logined

import (
	"dbus/org/freedesktop/login1"
	"encoding/json"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
	"sync"
)

// Manager manager logined user list
type Manager struct {
	core   *login1.Manager
	logger *log.Logger
	users  UserInfos
	locker sync.Mutex

	UserList string
}

const (
	dbusLogin1Dest = "org.freedesktop.login1"
	dbusLogin1Path = "/org/freedesktop/login1"
)

// Register register and install loginedManager on dbus
func Register(logger *log.Logger) (*Manager, error) {
	core, err := login1.NewManager(dbusLogin1Dest, dbusLogin1Path)
	if err != nil {
		return nil, err
	}

	var m = &Manager{
		core:   core,
		logger: logger,
	}

	go m.init()
	m.handleChanged()
	return m, nil
}

// Unregister destroy and free Manager object
func Unregister(m *Manager) {
	if m == nil {
		return
	}

	if m.core != nil {
		login1.DestroyManager(m.core)
	}
}

func (m *Manager) init() {
	// the result struct: {uid, username, path}
	list, err := m.core.ListUsers()
	if err != nil {
		m.logger.Warning("Failed to list logined user:", err)
		return
	}

	for _, value := range list {
		if len(value) != 3 {
			continue
		}

		m.logger.Debug("Create user info for:", value[2].(dbus.ObjectPath))
		info, err := newUserInfo(value[2].(dbus.ObjectPath))
		if err != nil {
			m.logger.Error("Failed to new user info:", err)
			continue
		}

		m.locker.Lock()
		m.users, _ = m.users.Add(info)
		m.locker.Unlock()
	}
	m.setPropUserList()
}

func (m *Manager) handleChanged() {
	m.core.ConnectUserNew(func(uid uint32, userPath dbus.ObjectPath) {
		m.logger.Debug("[Event] user new:", uid, userPath)
		info, err := newUserInfo(userPath)
		if err != nil {
			m.logger.Error("Failed to new user info:", err)
			return
		}

		m.locker.Lock()
		var added = false
		m.users, added = m.users.Add(info)
		m.locker.Unlock()
		if added {
			m.setPropUserList()
		}
	})

	m.core.ConnectUserRemoved(func(uid uint32, userPath dbus.ObjectPath) {
		m.logger.Debug("[Event] user removed:", uid, userPath)
		m.locker.Lock()
		var deleted = false
		m.users, deleted = m.users.Delete(uid)
		m.locker.Unlock()
		if deleted {
			m.setPropUserList()
		}
	})

	m.core.ConnectSessionNew(func(id string, sessionPath dbus.ObjectPath) {
		m.logger.Debug("[Event] session new:", id, sessionPath)
		session, err := login1.NewSession(dbusLogin1Dest, sessionPath)
		if err != nil {
			m.logger.Error("Failed to connect session:", err)
			return
		}

		list := session.User.Get()
		login1.DestroySession(session)
		if len(list) != 2 {
			return
		}

		m.logger.Debug("Create user info for:", list[1].(dbus.ObjectPath))
		info, err := newUserInfo(list[1].(dbus.ObjectPath))
		if err != nil {
			m.logger.Error("Failed to new user info:", err)
			return
		}
		m.locker.Lock()
		var added = false
		m.users, added = m.users.Add(info)
		m.locker.Unlock()
		if added {
			m.setPropUserList()
		}
	})
}

func (m *Manager) setPropUserList() {
	m.locker.Lock()
	defer m.locker.Unlock()

	if len(m.users) == 0 {
		return
	}

	data, err := json.Marshal(m.users)
	if err != nil {
		m.logger.Error("Failed to marshal users:", err)
		return
	}

	if m.UserList == string(data) {
		return
	}
	m.UserList = string(data)
	dbus.NotifyChange(m, "UserList")
}

// GetDBusInfo dbus session interface
func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.daemon.Accounts",
		ObjectPath: "/com/deepin/daemon/Logined",
		Interface:  "com.deepin.daemon.Logined",
	}
}
