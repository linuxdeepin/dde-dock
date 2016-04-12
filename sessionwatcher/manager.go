/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package sessionwatcher

import (
	"dbus/org/freedesktop/login1"
	"pkg.deepin.io/lib/dbus"
	"sync"
)

const (
	login1Dest = "org.freedesktop.login1"
	login1Path = "/org/freedesktop/login1"
)

type Manager struct {
	loginManager  *login1.Manager
	sessionLocker sync.Mutex
	userIds       map[string]uint32
	userSessions  map[uint32][]*login1.Session
}

func newManager() (*Manager, error) {
	loginObj, err := login1.NewManager(login1Dest, login1Path)
	if err != nil {
		logger.Warning("New login1 manager failed:", err)
		return nil, err
	}

	return &Manager{
		loginManager: loginObj,
		userIds:      make(map[string]uint32),
		userSessions: make(map[uint32][]*login1.Session),
	}, nil
}

func (m *Manager) destroy() {
	if m.loginManager != nil {
		m.destroyUserSessions()
	}

	if m.userIds != nil {
		m.userIds = nil
	}

	if m.userSessions != nil {
		m.userSessions = nil
	}
}

func (m *Manager) initUserSessions() {
	list, err := m.loginManager.ListSessions()
	if err != nil {
		logger.Warning("List sessions failed:", err)
		return
	}

	for _, v := range list {
		// v info: (id, uid, username, seat id, session path)
		if len(v) != 5 {
			logger.Warning("Invalid session info:", v)
			continue
		}

		id, ok := v[0].(string)
		if !ok {
			continue
		}

		p, ok := v[4].(dbus.ObjectPath)
		if !ok {
			continue
		}

		m.addSession(id, p)
	}

	m.loginManager.ConnectSessionNew(func(id string, path dbus.ObjectPath) {
		logger.Debug("Session added:", id, path)
		m.addSession(id, path)
	})

	m.loginManager.ConnectSessionRemoved(func(id string, path dbus.ObjectPath) {
		logger.Debug("Session removed:", id, path)
		m.deleteSession(id, path)
	})
}

func (m *Manager) destroyUserSessions() {
	m.sessionLocker.Lock()
	m.sessionLocker.Unlock()
	for _, ss := range m.userSessions {
		for _, s := range ss {
			login1.DestroySession(s)
			s = nil
		}
	}
	m.userSessions = nil
}

func (m *Manager) addSession(id string, path dbus.ObjectPath) {
	logger.Debug("Add session:", path)
	uid, session := newLoginSession(path)
	if session == nil {
		return
	}

	m.sessionLocker.Lock()
	m.userIds[id] = uid
	m.userSessions[uid] = append(m.userSessions[uid], session)
	m.sessionLocker.Unlock()

	session.Active.ConnectChanged(func() {
		if session == nil {
			return
		}
		m.handleSessionChanged(uid)
	})
	m.handleSessionChanged(uid)
}

func (m *Manager) deleteSession(id string, path dbus.ObjectPath) {
	m.sessionLocker.Lock()
	uid, ok := m.userIds[id]
	if !ok {
		m.sessionLocker.Unlock()
		return
	}

	var list []*login1.Session
	for _, s := range m.userSessions[uid] {
		if s.Path != path {
			list = append(list, s)
			continue
		}

		logger.Debug("Delete session:", path)
		login1.DestroySession(s)
		s = nil
	}
	m.userSessions[uid] = list
	m.sessionLocker.Unlock()
	m.handleSessionChanged(uid)
}

func (m *Manager) handleSessionChanged(uid uint32) {
	if m.isActive(uid) {
		suspendPulseSinks(0)
		suspendPulseSources(0)
	} else {
		suspendPulseSinks(1)
		suspendPulseSources(1)
	}
}

func (m *Manager) isActive(uid uint32) bool {
	var active bool = false
	m.sessionLocker.Lock()
	for _, s := range m.userSessions[uid] {
		logger.Debug("[isActive] info:", uid, s.Path, s.Active.Get())
		if s.Active.Get() {
			active = true
			break
		}
	}
	m.sessionLocker.Unlock()

	logger.Debug("Session state:", uid, active)
	return active
}
