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
	"pkg.deepin.io/lib/dbus"
)

func (m *Manager) OnPropertiesChanged(name string, oldv interface{}) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("%v", err)
		}
	}()
	logger.Debug("OnPropertiesChanged: " + name)
	var oldBool bool
	switch oldv.(type) {
	case bool:
		oldBool, _ = oldv.(bool)
	}

	// NetworkingEnabled, WirelessEnabled and WwanEnabled were managed
	// by NetworkManager, so restore to their old value here and get
	// the right value from NetworkManager.
	switch name {
	case "NetworkingEnabled":
		if oldBool != m.NetworkingEnabled {
			newValue := m.NetworkingEnabled
			m.NetworkingEnabled = oldBool
			m.switchHandler.setNetworkingEnabled(newValue)
		}
	case "WirelessEnabled":
		if oldBool != m.wirelessEnabled {
			newValue := m.wirelessEnabled
			m.wirelessEnabled = oldBool
			m.switchHandler.setWirelessEnabled(newValue)
		}
	case "WwanEnabled":
		if oldBool != m.wwanEnabled {
			newValue := m.wwanEnabled
			m.wwanEnabled = oldBool
			m.switchHandler.setWwanEnabled(newValue)
		}
	case "WiredEnabled":
		if oldBool != m.wiredEnabled {
			m.switchHandler.setWiredEnabled(m.wiredEnabled)
		}
	case "VpnEnabled":
		if oldBool != m.VpnEnabled {
			m.switchHandler.setVpnEnabled(m.VpnEnabled)
		}
	}
}

func (m *Manager) setPropNetworkingEnabled(value bool) {
	m.NetworkingEnabled = value
	dbus.NotifyChange(m, "NetworkingEnabled")
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
	dbus.NotifyChange(m, "VpnEnabled")
}

func (m *Manager) setPropActiveConnections() {
	m.ActiveConnections, _ = marshalJSON(m.activeConnections)
	dbus.NotifyChange(m, "ActiveConnections")
}

func (m *Manager) setPropState() {
	m.State = nmGetManagerState()
	dbus.NotifyChange(m, "State")
}

func (m *Manager) setPropDevices() {
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
	dbus.NotifyChange(m, "Devices")
}

func (m *Manager) setPropConnections() {
	m.Connections, _ = marshalJSON(m.connections)
	dbus.NotifyChange(m, "Connections")
}
