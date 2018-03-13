/*
 * Copyright (C) 2013 ~ 2018 Deepin Technology Co., Ltd.
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

package timedate

func (m *Manager) handlePropChanged() {
	if m.td1 == nil {
		return
	}

	m.td1.CanNTP.ConnectChanged(func() {
		m.PropsMu.Lock()
		m.setPropCanNTP(m.td1.CanNTP.Get())
		m.PropsMu.Unlock()
	})
	m.td1.NTP.ConnectChanged(func() {
		m.PropsMu.Lock()
		m.setPropNTP(m.td1.NTP.Get())
		m.PropsMu.Unlock()
	})
	m.td1.LocalRTC.ConnectChanged(func() {
		m.PropsMu.Lock()
		m.setPropLocalRTC(m.td1.LocalRTC.Get())
		m.PropsMu.Unlock()
	})
	m.td1.Timezone.ConnectChanged(func() {
		m.PropsMu.Lock()
		m.setPropTimezone(m.td1.Timezone.Get())
		m.PropsMu.Unlock()

		m.AddUserTimezone(m.Timezone)
	})
}
