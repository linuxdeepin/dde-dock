/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
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

package power

import (
	"pkg.deepin.io/lib/arch"
)

func (m *Manager) initLidSwitch() {
	if arch.Get() == arch.Sunway && isSWLidStateFileExist() {
		m.initLidSwitchSW()
	} else {
		m.initLidSwitchCommon()
	}
	logger.Debug("hasLidSwitch:", m.HasLidSwitch)
}

func (m *Manager) handleLidSwitchEvent(closed bool) {
	if closed {
		logger.Info("Lid Closed")
		err := m.service.Emit(m, "LidClosed")
		if err != nil {
			logger.Warning(err)
		}
	} else {
		logger.Info("Lid Opened")
		err := m.service.Emit(m, "LidOpened")
		if err != nil {
			logger.Warning(err)
		}
	}
}
