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

package main

import (
	"fmt"
	"log"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"github.com/msteinert/pam"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
)

const (
	PromptQuestion uint32 = iota + 1
	PromptSecret
	ErrorMsg
	TextInfo
	Failure
	Success
)

type Manager struct {
	service       *dbusutil.Service
	authLocker    sync.Mutex
	authUserTable map[string]chan string // 'pid+user': 'password'

	methods *struct {
		CurrentUser      func() `out:"username"`
		IsLiveCD         func() `in:"username" out:"result"`
		SwitchToUser     func() `in:"username"`
		AuthenticateUser func() `in:"username"`
		UnlockCheck      func() `in:"username,password"`
	}

	signals *struct {
		Event struct {
			eventType uint32
			pid       uint32
			username  string
			message   string
		}

		UserChanged struct {
			username string
		}
	}
}

const (
	dbusServiceName = "com.deepin.dde.LockService"
	dbusPath        = "/com/deepin/dde/LockService"
	dbusInterface   = "com.deepin.dde.LockService"
)

var _m *Manager

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	service, err := dbusutil.NewSystemService()
	if err != nil {
		log.Fatal("failed to new system service:", err)
	}

	_m = newManager(service)
	err = service.Export(dbusPath, _m)
	if err != nil {
		log.Fatal("failed to export:", err)
	}

	err = service.RequestName(dbusServiceName)
	if err != nil {
		log.Fatal("failed to request name:", err)
	}

	service.SetAutoQuitHandler(time.Minute*2, func() bool {
		_m.authLocker.Lock()
		canQuit := len(_m.authUserTable) == 0
		_m.authLocker.Unlock()
		return canQuit
	})
	service.Wait()
}

func (*Manager) GetInterfaceName() string {
	return dbusInterface
}

func newManager(service *dbusutil.Service) *Manager {
	var m = Manager{
		authUserTable: make(map[string]chan string),
		service:       service,
	}

	return &m
}

func (m *Manager) CurrentUser() (string, *dbus.Error) {
	username, err := getGreeterUser(greeterUserConfig)
	if err != nil {
		return "", dbusutil.ToError(err)
	}
	return username, nil
}

func (m *Manager) IsLiveCD(username string) (bool, *dbus.Error) {
	return isInLiveCD(username), nil
}

func (m *Manager) SwitchToUser(username string) *dbus.Error {
	current, _ := getGreeterUser(greeterUserConfig)
	if current == username {
		return nil
	}

	err := setGreeterUser(greeterUserConfig, username)
	if err != nil {
		return dbusutil.ToError(err)
	}
	if current != "" {
		m.service.Emit(m, "UserChanged", username)
	}
	return nil
}

func (m *Manager) AuthenticateUser(sender dbus.Sender, username string) *dbus.Error {
	if username == "" {
		return dbusutil.ToError(fmt.Errorf("no user to authenticate"))
	}

	m.authLocker.Lock()
	pid, err := m.service.GetConnPID(string(sender))
	if err != nil {
		return dbusutil.ToError(err)
	}

	id := getId(pid, username)
	_, ok := m.authUserTable[id]
	if ok {
		log.Println("In authenticating:", id)
		m.authLocker.Unlock()
		return nil
	}

	m.authUserTable[id] = make(chan string)
	m.authLocker.Unlock()
	go m.doAuthenticate(username, "", pid)
	return nil
}

func (m *Manager) UnlockCheck(sender dbus.Sender, username, password string) *dbus.Error {
	if username == "" {
		return dbusutil.ToError(fmt.Errorf("no user to authenticate"))
	}

	pid, err := m.service.GetConnPID(string(sender))
	if err != nil {
		return dbusutil.ToError(err)
	}
	id := getId(pid, username)

	m.authLocker.Lock()
	v, ok := m.authUserTable[id]
	m.authLocker.Unlock()
	if ok && v != nil {
		log.Println("In authenticating:", id)
		// in authenticate
		if password != "" {
			v <- password
		}
		return nil
	}

	go m.doAuthenticate(username, password, pid)
	return nil
}

func getId(pid uint32, username string) string {
	return strconv.Itoa(int(pid)) + username
}

func (m *Manager) doAuthenticate(username, password string, pid uint32) {
	handler, err := pam.StartFunc("lightdm", username, func(style pam.Style, msg string) (string, error) {
		switch style {
		// case pam.PromptEchoOn:
		// 	if msg != "" {
		// 		fmt.Println("Echo on:", msg)
		// 		m.sendEvent(PromptQuestion, pid, username, msg)
		// 	}
		// 	// TODO: read data from input
		// 	return "", nil
		case pam.PromptEchoOff, pam.PromptEchoOn:
			if password != "" {
				tmp := password
				password = ""
				return tmp, nil
			}

			if msg != "" {
				if style == pam.PromptEchoOff {
					log.Println("Echo off:", msg)
					m.sendEvent(PromptSecret, pid, username, msg)
				} else {
					log.Println("Echo on:", msg)
					m.sendEvent(PromptQuestion, pid, username, msg)
				}
			}

			id := getId(pid, username)
			m.authLocker.Lock()
			v, ok := m.authUserTable[id]
			m.authLocker.Unlock()
			if !ok || v == nil {
				return "", fmt.Errorf("no passwd channel found for %s", username)
			}
			log.Println("Join select:", id)
			select {
			case tmp, ok := <-v:
				if !ok {
					log.Println("Invalid select channel")
					return "", nil
				}

				m.authLocker.Lock()
				delete(m.authUserTable, id)
				m.authLocker.Unlock()
				close(v)
				v = nil
				return tmp, nil
			}
		case pam.ErrorMsg:
			if msg != "" {
				log.Println("ShowError:", msg)
				m.sendEvent(ErrorMsg, pid, username, msg)
			}
			return "", nil
		case pam.TextInfo:
			if msg != "" {
				log.Println("Text info:", msg)
				m.sendEvent(TextInfo, pid, username, msg)
			}
			return "", nil
		}
		return "", fmt.Errorf("unexpected style: %v", style)
	})
	if err != nil {
		log.Println("Failed to start pam:", err)
		m.sendEvent(Failure, pid, username, err.Error())
		return
	}

	err = handler.Authenticate(pam.DisallowNullAuthtok)

	id := getId(pid, username)
	m.authLocker.Lock()
	v, ok := m.authUserTable[id]
	if ok {
		if v != nil {
			close(v)
			v = nil
		}
		delete(m.authUserTable, id)
	}
	m.authLocker.Unlock()
	if err != nil {
		log.Println("Failed to authenticate:", err)
		m.sendEvent(Failure, pid, username, err.Error())
	} else {
		log.Println("Authenticate success")
		m.sendEvent(Success, pid, username, "Authenticated")
	}
	handler = nil
	debug.FreeOSMemory()
}

func (m *Manager) sendEvent(ty, pid uint32, username, msg string) {
	err := m.service.Emit(m, "Event", ty, pid, username, msg)
	if err != nil {
		log.Println("Failed to emit event:", ty, pid, username, msg)
	}
}
