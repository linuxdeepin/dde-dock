/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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
	"sync"

	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/proxy"

	libdisplay "github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.display"
	"github.com/linuxdeepin/go-dbus-factory/org.freedesktop.login1"
)

const (
	dbusServiceName = "com.deepin.daemon.SessionWatcher"
	dbusPath        = "/com/deepin/daemon/SessionWatcher"
	dbusInterface   = dbusServiceName
)

type Manager struct {
	service           *dbusutil.Service
	display           *libdisplay.Display
	loginManager      *login1.Manager
	systemSigLoop     *dbusutil.SignalLoop
	mu                sync.Mutex
	sessions          map[string]*login1.Session
	activeSessionType string

	PropsMu  sync.RWMutex
	IsActive bool
	//nolint
	methods  *struct {
		GetSessions        func() `out:"sessions"`
		IsX11SessionActive func() `out:"is_active"`
	}
}

var (
	_validSessionList = []string{
		"x11",
		"wayland",
	}
)

func newManager(service *dbusutil.Service) (*Manager, error) {
	manager := &Manager{
		service:  service,
		sessions: make(map[string]*login1.Session),
	}
	systemConn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	sessionConn := service.Conn()
	manager.loginManager = login1.NewManager(systemConn)
	manager.display = libdisplay.NewDisplay(sessionConn)

	manager.systemSigLoop = dbusutil.NewSignalLoop(systemConn, 10)
	manager.systemSigLoop.Start()
	manager.loginManager.InitSignalExt(manager.systemSigLoop, true)

	// default as active
	manager.IsActive = true
	return manager, nil
}

func (m *Manager) destroy() {
	m.mu.Lock()
	for _, session := range m.sessions {
		session.RemoveHandler(proxy.RemoveAllHandlers)
	}
	m.mu.Unlock()

	m.loginManager.RemoveHandler(proxy.RemoveAllHandlers)
	m.systemSigLoop.Stop()
}

func (*Manager) GetInterfaceName() string {
	return dbusInterface
}

func (m *Manager) initUserSessions() {
	sessions, err := m.loginManager.ListSessions(0)
	if err != nil {
		logger.Warning("List sessions failed:", err)
		return
	}

	for _, session := range sessions {
		m.addSession(session.SessionId, session.Path)
	}
	m.handleSessionChanged()

	_, err = m.loginManager.ConnectSessionNew(func(id string, path dbus.ObjectPath) {
		logger.Debug("Session added:", id, path)
		m.addSession(id, path)
		m.handleSessionChanged()
	})
	if err != nil {
		logger.Warning("ConnectSessionNew error:",err)
	}

	_, err = m.loginManager.ConnectSessionRemoved(func(id string, path dbus.ObjectPath) {
		logger.Debug("Session removed:", id, path)
		m.deleteSession(id, path)
		m.handleSessionChanged()
	})
	if err != nil {
		logger.Warning("ConnectSessionRemoved error:",err)
	}
}

func (m *Manager) addSession(id string, path dbus.ObjectPath) {
	systemConn := m.systemSigLoop.Conn()
	session, err := login1.NewSession(systemConn, path)
	if err != nil {
		logger.Warning(err)
		return
	}

	userInfo, err := session.User().Get(0)
	if err != nil {
		logger.Warning(err)
		return
	}

	uid := userInfo.UID
	logger.Debug("Add session:", id, path, uid)
	if !isCurrentUser(uid) {
		logger.Debug("Not the current user session:", id, path, uid)
		return
	}
	remote, err := session.Remote().Get(0)
	if err != nil {
		logger.Warning(err)
		return
	}
	if remote {
		logger.Debugf("session %v is remote", id)
		return
	}

	m.mu.Lock()
	m.sessions[id] = session
	m.mu.Unlock()

	session.InitSignalExt(m.systemSigLoop, true)
	err = session.Active().ConnectChanged(func(hasValue bool, value bool) {
		m.handleSessionChanged()
	})
	if err != nil {
		logger.Warning("ConnectChanged error:",err)
	}
}

func (m *Manager) deleteSession(id string, path dbus.ObjectPath) {
	m.mu.Lock()
	session, ok := m.sessions[id]
	if !ok {
		m.mu.Unlock()
		return
	}

	session.RemoveHandler(proxy.RemoveAllHandlers)
	logger.Debug("Delete session:", id, path)
	delete(m.sessions, id)
	m.mu.Unlock()
}

func (m *Manager) handleSessionChanged() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.sessions) == 0 {
		return
	}

	session := m.getActiveSession()
	var isActive bool
	var sessionType string
	if session != nil {
		isActive = true
		var err error
		sessionType, err = session.Type().Get(0)
		if err != nil {
			logger.Warning(err)
		}
	}
	m.activeSessionType = sessionType
	m.PropsMu.Lock()
	changed := m.setIsActive(isActive)
	m.PropsMu.Unlock()
	if !changed {
		return
	}

	if isActive {
		logger.Debug("[handleSessionChanged] Resume pulse")
		// fixed block when unused pulse-audio
		go suspendPulseSinks(0)
		go suspendPulseSources(0)

		logger.Debug("[handleSessionChanged] Refresh Brightness")
		go func() {
			_ = m.display.RefreshBrightness(0)
		}()
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
		err := m.service.EmitPropertyChanged(m, "IsActive", val)
		if err != nil {
			logger.Warning("EmitPropertyChanged error:",err)
		}
		return true
	}
	return false
}

func (m *Manager) getActiveSession() *login1.Session {
	for _, session := range m.sessions {
		active, err := session.Active().Get(0)
		if err != nil {
			logger.Warning(err)
			continue
		}
		if active {
			return session
		}
	}
	return nil
}

func (m *Manager) IsX11SessionActive() (bool, *dbus.Error) {
	m.mu.Lock()
	ty := m.activeSessionType
	m.mu.Unlock()
	for _, session := range _validSessionList {
		if session == ty {
			return true, nil
		}
	}
	return false, nil
}

func (m *Manager) GetSessions() (ret []dbus.ObjectPath, err *dbus.Error) {
	m.mu.Lock()
	ret = make([]dbus.ObjectPath, len(m.sessions))
	i := 0
	for _, session := range m.sessions {
		ret[i] = session.Path_()
		i++
	}
	m.mu.Unlock()
	return
}
