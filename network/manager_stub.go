/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

package network

import (
	"errors"

	dbus "github.com/godbus/dbus"
	"pkg.deepin.io/dde/daemon/network/nm"
	"pkg.deepin.io/lib/dbusutil"
)

func (m *Manager) networkingEnabledWriteCb(write *dbusutil.PropertyWrite) *dbus.Error {
	// currently not need
	return nil
}

func (m *Manager) vpnEnabledWriteCb(write *dbusutil.PropertyWrite) *dbus.Error {
	enabled, ok := write.Value.(bool)
	if !ok {
		err := errors.New("type of value is not bool")
		logger.Warning(err)
		return dbusutil.ToError(err)
	}
	err := m.sysNetwork.VpnEnabled().Set(0, enabled)
	if err != nil {
		logger.Warning(err)
		return dbusutil.ToError(err)
	}
	return nil
}

func (m *Manager) updatePropActiveConnections() {
	activeConnections, _ := marshalJSON(m.activeConnections)
	m.setPropActiveConnections(activeConnections)
}

func (m *Manager) updatePropState() {
	state := nmGetManagerState()
	m.setPropState(state)
}

func (m *Manager) updatePropConnectivity() {
	connectivity, _ := nmManager.Connectivity().Get(0)
	m.setPropConnectivity(connectivity)
}

func (m *Manager) updatePropDevices() {
	filteredDevices := make(map[string][]*device)
	for key, devices := range m.devices {
		filteredDevices[key] = make([]*device, 0)
		for _, d := range devices {
			ignoreIphoneUsbDevice := d.UsbDevice &&
				d.State <= nm.NM_DEVICE_STATE_UNAVAILABLE &&
				d.Driver == "ipheth"
			if !ignoreIphoneUsbDevice {
				filteredDevices[key] = append(filteredDevices[key], d)
			}
		}
	}
	devices, _ := marshalJSON(filteredDevices)
	m.setPropDevices(devices)
}

func (m *Manager) updatePropConnections() {
	connections, _ := marshalJSON(m.connections)
	m.setPropConnections(connections)
}

func (m *Manager) updatePropWirelessAccessPoints() {
	aps, _ := marshalJSON(m.accessPoints)
	m.setPropWirelessAccessPoints(aps)
}
