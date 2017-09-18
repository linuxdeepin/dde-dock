/*
 * Copyright (C) 2016 ~ 2017 Deepin Technology Co., Ltd.
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

const (
	dbusDest = "com.deepin.system.Power"
	dbusPath = "/com/deepin/system/Power"
	dbusIFC  = dbusDest
)

func (m *Manager) GetBatteries() []*Battery {
	ret := make([]*Battery, 0, len(m.batteries))
	for _, bat := range m.batteries {
		ret = append(ret, bat)
	}
	return ret
}

func (m *Manager) RefreshBatteries() {
	logger.Debug("RefreshBatteries")
	for _, bat := range m.batteries {
		bat.Refresh()
	}
}

func (m *Manager) RefreshMains() {
	logger.Debug("RefreshMains")
	if m.ac == nil {
		return
	}

	device := m.ac.newDevice()
	if device == nil {
		logger.Warning("RefreshMains: ac.newDevice failed")
		return
	}
	m.refreshAC(device)
}

func (m *Manager) Refresh() {
	m.RefreshMains()
	m.RefreshBatteries()
}
