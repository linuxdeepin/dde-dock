/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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
