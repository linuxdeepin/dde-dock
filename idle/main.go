/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
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

package main

import (
	"dlib/dbus"
	"dlib/logger"
	"sync"
)

var (
	lock sync.Mutex

	genID = func() func() int32 {
		id := int32(0)
		return func() int32 {
			lock.Lock()
			tmp := id
			id += 1
			lock.Unlock()
			return tmp
		}
	}()
)

func (m *Manager) SimulateUserActivity() {
}

func (m *Manager) UnsimulateUserActivity() {
}

/*
 * name: the property of signal 'IdleTimeOut'
 * timeout: idle timeout
 *
 * if timeout, the signal 'IdleTimeOut' will be send, and
 * append param name
 */
func (m *Manager) RegisterIdleTick(name string, timeout int32) int32 {
	return genID()
}

func (m *Manager) UnregisterIdleTick(cookie int32) {
}

func NewManager() *Manager {
	return &Manager{}
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			logger.Println("recover err:", err)
		}
	}()

	m := NewManager()
	err := dbus.InstallOnSession(m)
	if err != nil {
		logger.Println("Install Session DBus Failed:", err)
		panic(err)
	}
	dbus.DealWithUnhandledMessage()

	select {}
}
