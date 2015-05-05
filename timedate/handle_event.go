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

package timedate

func (m *Manager) handlePropChanged() {
	if m.td1 == nil {
		return
	}

	m.td1.CanNTP.ConnectChanged(func() {
		m.setPropBool(&m.CanNTP, "CanNTP", m.td1.CanNTP.Get())
	})
	m.td1.NTP.ConnectChanged(func() {
		m.setPropBool(&m.NTP, "NTP", m.td1.NTP.Get())
	})
	m.td1.LocalRTC.ConnectChanged(func() {
		m.setPropBool(&m.LocalRTC, "LocalRTC", m.td1.LocalRTC.Get())
	})
	m.td1.Timezone.ConnectChanged(func() {
		m.setPropString(&m.Timezone, "Timezone", m.td1.Timezone.Get())
	})
}
