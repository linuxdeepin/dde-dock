/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
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

package logined

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/linuxdeepin/go-dbus-factory/org.freedesktop.login1"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/proxy"
	"pkg.deepin.io/lib/log"
)

// Manager manager logined user list
type Manager struct {
	service    *dbusutil.Service
	sysSigLoop *dbusutil.SignalLoop
	core       *login1.Manager
	logger     *log.Logger

	userSessions map[uint32]SessionInfos
	locker       sync.Mutex

	UserList string
}

const (
	DBusPath = "/com/deepin/daemon/Logined"
)

// Register register and install loginedManager on dbus
func Register(logger *log.Logger, service *dbusutil.Service) (*Manager, error) {
	systemBus, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	core := login1.NewManager(systemBus)
	sysSigLoop := dbusutil.NewSignalLoop(systemBus, 10)
	sysSigLoop.Start()
	var m = &Manager{
		service:      service,
		core:         core,
		logger:       logger,
		userSessions: make(map[uint32]SessionInfos),
		sysSigLoop:   sysSigLoop,
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

	m.core.RemoveHandler(proxy.RemoveAllHandlers)
	m.sysSigLoop.Stop()

	if m.userSessions != nil {
		m.userSessions = nil
	}

	m = nil
}

func (m *Manager) init() {
	// the result struct: {id, uid, username, seat, path}
	sessions, err := m.core.ListSessions(0)
	if err != nil {
		m.logger.Warning("Failed to list sessions:", err)
		return
	}

	for _, session := range sessions {
		m.addSession(session.Path)
	}
	m.setPropUserList()
}

func (m *Manager) handleChanged() {
	m.core.InitSignalExt(m.sysSigLoop, true)
	m.core.ConnectSessionNew(func(id string, sessionPath dbus.ObjectPath) {
		m.logger.Debug("[Event] session new:", id, sessionPath)
		added := m.addSession(sessionPath)
		if added {
			m.setPropUserList()
		}
	})
	m.core.ConnectSessionRemoved(func(id string, sessionPath dbus.ObjectPath) {
		m.logger.Debug("[Event] session remove:", id, sessionPath)
		deleted := m.deleteSession(sessionPath)
		if deleted {
			m.setPropUserList()
		}
	})
}

func (m *Manager) addSession(sessionPath dbus.ObjectPath) bool {
	m.logger.Debug("Create user session for:", sessionPath)
	info, err := newSessionInfo(sessionPath)
	if err != nil {
		m.logger.Warning("Failed to add session:", sessionPath, err)
		return false
	}

	m.locker.Lock()
	defer m.locker.Unlock()
	infos, ok := m.userSessions[info.Uid]
	if !ok {
		m.userSessions[info.Uid] = SessionInfos{info}
		return true
	}

	var added = false
	infos, added = infos.Add(info)
	m.userSessions[info.Uid] = infos
	return added
}

func (m *Manager) deleteSession(sessionPath dbus.ObjectPath) bool {
	m.logger.Debug("Delete user session for:", sessionPath)
	m.locker.Lock()
	defer m.locker.Unlock()
	var deleted = false
	for uid, infos := range m.userSessions {
		tmp, ok := infos.Delete(sessionPath)
		if !ok {
			continue
		}
		deleted = true
		if len(tmp) == 0 {
			delete(m.userSessions, uid)
		} else {
			m.userSessions[uid] = tmp
		}
		break
	}
	return deleted
}

func (m *Manager) setPropUserList() {
	m.locker.Lock()
	defer m.locker.Unlock()

	if len(m.userSessions) == 0 {
		return
	}

	data := m.marshalUserSessions()
	if m.UserList == string(data) {
		return
	}
	m.UserList = string(data)
	m.service.EmitPropertyChanged(m, "UserList", m.UserList)
}

func (m *Manager) marshalUserSessions() string {
	if len(m.userSessions) == 0 {
		return ""
	}

	var ret = "{"
	for k, v := range m.userSessions {
		data, err := json.Marshal(v)
		if err != nil {
			m.logger.Warning("Failed to marshal:", v, err)
			continue
		}
		ret += fmt.Sprintf("\"%v\":%s,", k, string(data))
	}

	v := []byte(ret)
	v[len(v)-1] = '}'
	return string(v)
}

func (*Manager) GetInterfaceName() string {
	return "com.deepin.daemon.Logined"
}
