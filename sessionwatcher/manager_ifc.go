/**
 * Copyright (c) 2011 ~ 2015 Deepin, Inc.
 *               2013 ~ 2015 jouyouyun
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
	"time"
)

const (
	loopDuration       = time.Second * 20
	admissibleDuration = time.Second * 2

	maxLaunchTimes = 5
)

func (m *Manager) StartLoop() {
	var (
		dockTimes int
		prev      = time.Now().Unix()
	)
	for {
		select {
		case <-time.After(loopDuration):
			if m.launchDockFailed {
				logger.Debug("Launch all programs failed")
				m.QuitLoop()
				return
			}

			duration := time.Now().Unix() - prev
			m.handleDockLaunch(&dockTimes, duration)
			logger.Debug("Handle programs launch end")
		case <-m.quit:
			m.quit = nil
			return
		}
	}
}

func (m *Manager) QuitLoop() {
	if m.quit == nil {
		return
	}

	close(m.quit)
}

func (m *Manager) handleDockLaunch(times *int, duration int64) {
	if !m.canLaunchDock() {
		*times = 0
		logger.Debug("No need to launch dde-dock")
		return
	}
	m.restartDock()

	if duration < int64(loopDuration+admissibleDuration) {
		*times += 1
	} else {
		*times = 0
	}

	logger.Debug("dde-dock launch times:", *times)
	if *times == maxLaunchTimes {
		m.launchDockFailed = true
	}
}
