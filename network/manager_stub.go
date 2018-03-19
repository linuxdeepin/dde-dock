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
	"pkg.deepin.io/dde/daemon/network/nm"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
)

func (m *Manager) networkingEnabledWriteCb(write *dbusutil.PropertyWrite) *dbus.Error {
	enabled, ok := write.Value.(bool)
	if !ok {
		logger.Warning("type of value is not bool")
		return nil
	}
	m.switchHandler.setNetworkingEnabled(enabled)
	return nil
}

func (m *Manager) vpnEnabledWriteCb(write *dbusutil.PropertyWrite) *dbus.Error {
	enabled, ok := write.Value.(bool)
	if !ok {
		logger.Warning("type of value is not bool")
		return nil
	}
	m.switchHandler.setVpnEnabled(enabled)
	return nil
}

func (m *Manager) setPropNetworkingEnabled(value bool) {
	m.NetworkingEnabled = value
	m.service.EmitPropertyChanged(m, "NetworkingEnabled", value)
}
func (m *Manager) setPropWirelessEnabled(value bool) {
	m.wirelessEnabled = value
}
func (m *Manager) setPropWwanEnabled(value bool) {
	m.wwanEnabled = value
}
func (m *Manager) setPropWiredEnabled(value bool) {
	m.wiredEnabled = value
}

func (m *Manager) setPropVpnEnabled(value bool) {
	m.VpnEnabled = value
	m.service.EmitPropertyChanged(m, "VpnEnabled", value)
}

func (m *Manager) updatePropActiveConnections() {
	m.ActiveConnections, _ = marshalJSON(m.activeConnections)
	m.service.EmitPropertyChanged(m, "ActiveConnections", m.ActiveConnections)
}

func (m *Manager) updatePropState() {
	m.State = nmGetManagerState()
	m.service.EmitPropertyChanged(m, "State", m.State)
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
	m.Devices, _ = marshalJSON(filteredDevices)
	m.service.EmitPropertyChanged(m, "Devices", m.Devices)
}

func (m *Manager) updatePropConnections() {
	m.Connections, _ = marshalJSON(m.connections)
	m.service.EmitPropertyChanged(m, "Connections", m.Connections)
}
