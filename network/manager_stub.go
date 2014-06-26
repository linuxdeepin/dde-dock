/**
 * Copyright (c) 2014 Deepin, Inc.
 *               2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
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

package network

import (
	"pkg.linuxdeepin.com/lib/dbus"
)

func (m *Manager) OnPropertiesChanged(name string, oldv interface{}) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("%v", err)
		}
	}()
	logger.Debug("OnPropertiesChanged: " + name)
	switch name {
	case "NetworkingEnabled":
		if m.NetworkingEnabled != nmManager.NetworkingEnabled.Get() {
			nmManagerEnable(m.NetworkingEnabled)
		} else {
			logger.Warning("NetworkingEnabled already set as", m.NetworkingEnabled)
		}
	case "WirelessEnabled":
		if m.WirelessEnabled != nmManager.WirelessEnabled.Get() {
			nmManager.WirelessEnabled.Set(m.WirelessEnabled)
		} else {
			logger.Warning("WirelessEnabled already set as", m.WirelessEnabled)
		}
	case "WwanEnabled":
		if m.WwanEnabled != nmManager.WwanEnabled.Get() {
			nmManager.WwanEnabled.Set(m.WwanEnabled)
		} else {
			logger.Warning("WwanEnabled already set as", m.WwanEnabled)
		}
	case "WiredEnabled":
		m.updatePropWiredEnabled()
	case "VpnEnabled":
		m.updatePropVpnEnabled()
	}
}

func (m *Manager) updatePropNetworkingEnabled() {
	m.NetworkingEnabled = nmManager.NetworkingEnabled.Get()

	// TODO setup other global switches

	dbus.NotifyChange(m, "NetworkingEnabled")
}

func (m *Manager) updatePropWirelessEnabled() {
	m.WirelessEnabled = nmManager.WirelessEnabled.Get()
	// setup wireless devices switches
	for _, devPath := range nmGetSpecialDevices(NM_DEVICE_TYPE_WIFI) {
		if m.WirelessEnabled {
			m.restoreDeviceState(devPath)
		} else {
			m.EnableDevice(devPath, false)
		}
	}
	dbus.NotifyChange(m, "WirelessEnabled")
}

func (m *Manager) updatePropWwanEnabled() {
	m.WwanEnabled = nmManager.WwanEnabled.Get()
	// setup modem devices switches
	for _, devPath := range nmGetSpecialDevices(NM_DEVICE_TYPE_MODEM) {
		if m.WwanEnabled {
			m.restoreDeviceState(devPath)
		} else {
			m.EnableDevice(devPath, false)
		}
	}
	dbus.NotifyChange(m, "WirelessEnabled")
}

func (m *Manager) updatePropWiredEnabled() {
	// setup wired devices switches
	for _, devPath := range nmGetSpecialDevices(NM_DEVICE_TYPE_ETHERNET) {
		if m.WiredEnabled {
			m.restoreDeviceState(devPath)
		} else {
			m.EnableDevice(devPath, false)
		}
	}
	m.config.setWiredEnabled(m.WiredEnabled)
	dbus.NotifyChange(m, "WiredEnabled")
}

func (m *Manager) updatePropVpnEnabled() {
	// TODO
	m.config.setVpnEnabled(m.VpnEnabled)
	dbus.NotifyChange(m, "VpnEnabled")
}

func (m *Manager) updatePropActiveConnections() {
	m.ActiveConnections, _ = marshalJSON(m.activeConnections)
	dbus.NotifyChange(m, "ActiveConnections")
	// logger.Debug("ActiveConnections:", m.ActiveConnections) // TODO test
}

func (m *Manager) updatePropState() {
	m.State = nmGetManagerState()
	dbus.NotifyChange(m, "State")
}

func (m *Manager) updatePropWiredDevices() {
	dbus.NotifyChange(m, "WiredDevices")
}
func (m *Manager) updatePropWirelessDevices() {
	dbus.NotifyChange(m, "WirelessDevices")
}
func (m *Manager) updatePropDevices() {
	m.Devices, _ = marshalJSON(m.devices)
	dbus.NotifyChange(m, "Devices")
	// logger.Debug("updatePropDevices", m.Devices) // TODO test
}

// TODO
func (m *Manager) updatePropAccessPoints() {
	// m.AccessPoints, _ = marshalJSON(m.accessPoints)
	// dbus.NotifyChange(m, "AccessPoints")
	// testJSON, _ := marshalJSON(m.accessPoints)
	// logger.Debug("updatePropAccessPoints", testJSON) // TODO test
}

func (m *Manager) updatePropConnections() {
	m.Connections, _ = marshalJSON(m.connections)
	// logger.Debug(m.Connections) // TODO test
	dbus.NotifyChange(m, "Connections")
}
