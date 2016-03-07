/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
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
		m.AddUserTimezone(m.Timezone)
	})
}
