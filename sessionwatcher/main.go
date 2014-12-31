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

package sessionwatcher

import (
	"pkg.linuxdeepin.com/lib/log"
	"time"
)

var (
	logger = log.NewLogger("dde-daemon/sessionwatcher")
)

type Manager struct {
	dock      *Dock
	appelt    *DockApplet
	exitTimer chan struct{}
}

func NewManager() *Manager {
	m := &Manager{}

	m.dock = NewDock()
	m.appelt = NewDockApplet()
	m.exitTimer = make(chan struct{})

	return m
}

func (m *Manager) destroy() {
	close(m.exitTimer)
}

func (m *Manager) startTimer() {
	for {
		timer := time.NewTimer(time.Second * 5)
		select {
		case <-timer.C:
			go m.appelt.restartDockApplet()
			go m.dock.restartDock()
		case <-m.exitTimer:
			return
		}
	}
}

var _m *Manager

func Start() {
	if _m != nil {
		return
	}

	logger.BeginTracing()

	_m = NewManager()
	go _m.startTimer()
}

func Stop() {
	if _m == nil {
		return
	}

	logger.EndTracing()
	_m.destroy()
	_m = nil
}
