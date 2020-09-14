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
	dbus "github.com/godbus/dbus"
	"pkg.deepin.io/lib/dbusutil"
)

const (
	dbusServiceName = "com.deepin.system.Power"
	dbusPath        = "/com/deepin/system/Power"
	dbusInterface   = dbusServiceName
)

func (*Manager) GetInterfaceName() string {
	return dbusInterface
}

func (m *Manager) GetBatteries() ([]dbus.ObjectPath, *dbus.Error) {
	m.batteriesMu.Lock()

	result := make([]dbus.ObjectPath, len(m.batteries))
	idx := 0
	for _, bat := range m.batteries {
		result[idx] = bat.getObjPath()
		idx++
	}

	m.batteriesMu.Unlock()
	return result, nil
}

func (m *Manager) refreshBatteries() {
	logger.Debug("RefreshBatteries")
	m.batteriesMu.Lock()
	for _, bat := range m.batteries {
		bat.Refresh()
	}
	m.batteriesMu.Unlock()
}

func (m *Manager) RefreshBatteries() *dbus.Error {
	m.refreshBatteries()
	return nil
}

func (m *Manager) RefreshMains() *dbus.Error {
	logger.Debug("RefreshMains")
	if m.ac == nil {
		return nil
	}

	device := m.ac.newDevice()
	if device == nil {
		logger.Warning("RefreshMains: ac.newDevice failed")
		return nil
	}
	m.refreshAC(device)
	return nil
}

func (m *Manager) Refresh() *dbus.Error {
	err := m.RefreshMains()
	if err != nil {
		return err
	}
	err = m.RefreshBatteries()
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) SetCpuBoost(enabled bool) *dbus.Error {
	err := m.cpus.SetBoostEnabled(enabled)
	if err == nil {
		m.CpuBoost = enabled
	}
	return dbusutil.ToError(err)
}

func (m *Manager) SetCpuGovernor(governor string) *dbus.Error {
	err := m.cpus.SetGovernor(governor)
	if err == nil {
		m.CpuGovernor = governor
	}
	return dbusutil.ToError(err)
}
