/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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

package sessionwatcher

import (
	libdisplay "dbus/com/deepin/daemon/display"
	"dbus/org/freedesktop/login1"
	"pkg.deepin.io/lib/dbus"
	"sync"
)

const (
	login1Dest         = "org.freedesktop.login1"
	login1Path         = "/org/freedesktop/login1"
	displayDBusDest    = "com.deepin.daemon.Display"
	displayDBusObjPath = "/com/deepin/daemon/Display"
)

type Manager struct {
	display           *libdisplay.Display
	loginManager      *login1.Manager
	sessionLocker     sync.Mutex
	sessions          map[string]*login1.Session
	IsActive          bool
	activeSessionType string
}

func newManager() (*Manager, error) {
	manager := &Manager{
		sessions: make(map[string]*login1.Session),
	}
	var err error
	manager.loginManager, err = login1.NewManager(login1Dest, login1Path)
	if err != nil {
		logger.Warning("New login1 manager failed:", err)
		return nil, err
	}

	manager.display, err = libdisplay.NewDisplay(displayDBusDest, displayDBusObjPath)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}

	// default as active
	manager.IsActive = true
	return manager, nil
}

func (m *Manager) destroy() {
	if m.sessions != nil {
		m.destroySessions()
	}

	if m.display != nil {
		libdisplay.DestroyDisplay(m.display)
		m.display = nil
	}

	if m.loginManager != nil {
		login1.DestroyManager(m.loginManager)
		m.loginManager = nil
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
	m.handleSessionChanged()

	m.loginManager.ConnectSessionNew(func(id string, path dbus.ObjectPath) {
		logger.Debug("Session added:", id, path)
		m.addSession(id, path)
		m.handleSessionChanged()
	})

	m.loginManager.ConnectSessionRemoved(func(id string, path dbus.ObjectPath) {
		logger.Debug("Session removed:", id, path)
		m.deleteSession(id, path)
		m.handleSessionChanged()
	})
}

func (m *Manager) destroySessions() {
	m.sessionLocker.Lock()
	for _, s := range m.sessions {
		login1.DestroySession(s)
		s = nil
	}
	m.sessions = nil
	m.sessionLocker.Unlock()
}

func (m *Manager) addSession(id string, path dbus.ObjectPath) {
	uid, session := newLoginSession(path)
	if session == nil {
		return
	}

	logger.Debug("Add session:", id, path, uid)
	if !isCurrentUser(uid) {
		logger.Debug("Not the current user session:", id, path, uid)
		login1.DestroySession(session)
		return
	}

	if session.Remote.Get() {
		logger.Debugf("session %v is remote", id)
		login1.DestroySession(session)
		return
	}

	m.sessionLocker.Lock()
	m.sessions[id] = session
	m.sessionLocker.Unlock()

	session.Active.ConnectChanged(func() {
		if session == nil {
			return
		}
		m.handleSessionChanged()
	})
}

func (m *Manager) deleteSession(id string, path dbus.ObjectPath) {
	m.sessionLocker.Lock()
	session, ok := m.sessions[id]
	if !ok {
		m.sessionLocker.Unlock()
		return
	}

	logger.Debug("Delete session:", id, path)
	login1.DestroySession(session)
	session = nil
	delete(m.sessions, id)
	m.sessionLocker.Unlock()
}

func (m *Manager) handleSessionChanged() {
	if len(m.sessions) == 0 {
		return
	}

	session := m.getActiveSession()
	var isActive bool
	var sessionType string
	if session != nil {
		isActive = true
		sessionType = session.Type.Get()
	}
	m.activeSessionType = sessionType
	changed := m.setIsActive(isActive)
	if !changed {
		return
	}

	if isActive {
		logger.Debug("[handleSessionChanged] Resume pulse")
		// fixed block when unused pulseaudio
		go suspendPulseSinks(0)
		go suspendPulseSources(0)

		logger.Debug("[handleSessionChanged] Refresh Brightness")
		go m.display.RefreshBrightness()
	} else {
		logger.Debug("[handleSessionChanged] Suspend pulse")
		go suspendPulseSinks(1)
		go suspendPulseSources(1)
	}
}

// return is changed?
func (m *Manager) setIsActive(val bool) bool {
	if m.IsActive != val {
		m.IsActive = val
		logger.Debug("[setIsActive] IsActive changed:", val)
		dbus.NotifyChange(m, "IsActive")
		return true
	}
	return false
}

func (m *Manager) getActiveSession() *login1.Session {
	m.sessionLocker.Lock()
	defer m.sessionLocker.Unlock()

	for _, session := range m.sessions {
		active := session.Active.Get()
		if active {
			return session
		}
	}
	return nil
}

func (m *Manager) IsX11SessionActive() bool {
	return m.activeSessionType == "x11"
}

func (m *Manager) GetSessions() (ret []dbus.ObjectPath) {
	m.sessionLocker.Lock()
	ret = make([]dbus.ObjectPath, len(m.sessions))
	i := 0
	for _, session := range m.sessions {
		ret[i] = session.Path
		i++
	}
	m.sessionLocker.Unlock()
	return
}
