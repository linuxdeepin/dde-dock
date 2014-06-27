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
		m.setNetworkingEnabled(m.NetworkingEnabled)
	case "WirelessEnabled":
		m.setWirelessEnabled(m.WirelessEnabled)
	case "WwanEnabled":
		m.setWwanEnabled(m.WwanEnabled)
	case "WiredEnabled":
		m.setWiredEnabled(m.WiredEnabled)
	case "VpnEnabled":
		m.setVpnEnabled(m.VpnEnabled)
	}
}

func (m *Manager) setNetworkingEnabled(enabled bool) {
	m.NetworkingEnabled = enabled
	// setup global device switches
	if enabled {
		m.restoreGlobalDeviceSwitches()
	} else {
		m.disableGlobalDeviceSwitches()
	}
	nmSetNetworkingEnabled(enabled)
}
func (m *Manager) restoreGlobalDeviceSwitches() {
	nmSetWirelessEnabled(m.config.LastWirelessEnabled)
	nmSetWwanEnabled(m.config.LastWwanEnabled)
	m.updatePropWiredEnabled(m.config.LastWiredEnabled)
	m.updatePropVpnEnabled(m.config.LastVpnEnabled)
}
func (m *Manager) disableGlobalDeviceSwitches() {
	m.config.setLastWirelessEnabled(m.WirelessEnabled)
	nmSetWirelessEnabled(false)

	m.config.setLastWwanEnabled(m.WwanEnabled)
	nmSetWwanEnabled(false)

	m.updatePropWiredEnabled(false)
	m.updatePropVpnEnabled(false)
}

func (m *Manager) setWirelessEnabled(enabled bool) {
	m.WirelessEnabled = enabled
	if m.NetworkingEnabled {
		nmSetWirelessEnabled(enabled)
	} else {
		// if NetworkingEnabled is off, turn it on, and only keep
		// current global device alive
		m.config.setLastGlobalSwithes(false)
		m.config.setLastWirelessEnabled(true)
		m.setNetworkingEnabled(true)
	}
}
func (m *Manager) setWwanEnabled(enabled bool) {
	m.WwanEnabled = enabled
	if m.NetworkingEnabled {
		nmSetWwanEnabled(enabled)
	} else {
		// if NetworkingEnabled is off, turn it on, and only keep
		// current global device alive
		m.config.setLastGlobalSwithes(false)
		m.config.setLastWwanEnabled(true)
		m.setNetworkingEnabled(true)
	}
}
func (m *Manager) setWiredEnabled(enabled bool) {
	m.WiredEnabled = enabled
	if m.NetworkingEnabled {
		m.updatePropWiredEnabled(enabled)
	} else {
		// if NetworkingEnabled is off, turn it on, and only keep
		// current global device alive
		m.config.setLastGlobalSwithes(false)
		m.config.setLastWiredEnabled(true)
		m.setNetworkingEnabled(true)
	}
}
func (m *Manager) setVpnEnabled(enabled bool) {
	m.VpnEnabled = enabled
	if m.NetworkingEnabled {
		m.updatePropVpnEnabled(enabled)
	} else {
		// if NetworkingEnabled is off, turn it on, and only keep
		// current global device alive
		m.config.setLastGlobalSwithes(false)
		m.config.setLastVpnEnabled(true)
		m.setNetworkingEnabled(true)
	}
}

func (m *Manager) updatePropNetworkingEnabled() {
	m.NetworkingEnabled = nmManager.NetworkingEnabled.Get()
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
	dbus.NotifyChange(m, "WwanEnabled")
}
func (m *Manager) updatePropWiredEnabled(enabled bool) {
	m.WiredEnabled = enabled
	// setup wired devices switches
	for _, devPath := range nmGetSpecialDevices(NM_DEVICE_TYPE_ETHERNET) {
		if enabled {
			m.restoreDeviceState(devPath)
		} else {
			m.EnableDevice(devPath, false)
		}
	}
	m.config.setWiredEnabled(enabled)
	dbus.NotifyChange(m, "WiredEnabled")
}
func (m *Manager) updatePropVpnEnabled(enabled bool) {
	m.VpnEnabled = enabled
	m.config.setVpnEnabled(enabled)
	// setup vpn connections
	for _, uuid := range nmGetSpecialConnectionUuids(NM_SETTING_VPN_SETTING_NAME) {
		if enabled {
			m.restoreVpnConnectionState(uuid)
		} else {
			m.deactivateVpnConnection(uuid)
		}
	}
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
