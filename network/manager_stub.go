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
			m.setNetworkingEnabled(newValue)
		}
	case "WirelessEnabled":
		if oldBool != m.WirelessEnabled {
			newValue := m.WirelessEnabled
			m.WirelessEnabled = oldBool
			m.setWirelessEnabled(newValue)
		}
	case "WwanEnabled":
		if oldBool != m.WwanEnabled {
			newValue := m.WwanEnabled
			m.WwanEnabled = oldBool
			m.setWwanEnabled(newValue)
		}
	case "WiredEnabled":
		if oldBool != m.WiredEnabled {
			m.setWiredEnabled(m.WiredEnabled)
		}
	case "VpnEnabled":
		if oldBool != m.VpnEnabled {
			m.setVpnEnabled(m.VpnEnabled)
		}
	}
}

func (m *Manager) setNetworkingEnabled(enabled bool) {
	logger.Debug("setNetworkingEnabled", enabled) // TODO test
	nmSetNetworkingEnabled(enabled)
}
func (m *Manager) setWirelessEnabled(enabled bool) {
	if m.NetworkingEnabled {
		nmSetWirelessEnabled(enabled)
	} else {
		// if NetworkingEnabled is off, turn it on, and only keep
		// current global device switch alive
		m.config.setLastGlobalSwithes(false)
		m.config.setLastWirelessEnabled(true)
		m.setNetworkingEnabled(true)
	}
}
func (m *Manager) setWwanEnabled(enabled bool) {
	if m.NetworkingEnabled {
		nmSetWwanEnabled(enabled)
	} else {
		// if NetworkingEnabled is off, turn it on, and only keep
		// current global device switch alive
		m.config.setLastGlobalSwithes(false)
		m.config.setLastWwanEnabled(true)
		m.setNetworkingEnabled(true)
	}
}
func (m *Manager) setWiredEnabled(enabled bool) {
	if m.NetworkingEnabled {
		m.setPropWiredEnabled(enabled)
	} else {
		// if NetworkingEnabled is off, turn it on, and only keep
		// current global device switch alive
		m.config.setLastGlobalSwithes(false)
		m.config.setLastWiredEnabled(true)
		m.setNetworkingEnabled(true)
	}
}
func (m *Manager) setVpnEnabled(enabled bool) {
	if m.NetworkingEnabled {
		m.setPropVpnEnabled(enabled)
	} else {
		// if NetworkingEnabled is off, turn it on, and only keep
		// current global device switch alive
		m.config.setLastGlobalSwithes(false)
		m.config.setLastVpnEnabled(true)
		m.setNetworkingEnabled(true)
	}
}

func (m *Manager) initPropNetworkingEnabled() {
	m.NetworkingEnabled = nmManager.NetworkingEnabled.Get()
	if !m.NetworkingEnabled {
		m.doTurnOffGlobalDeviceSwitches()
	}
	m.doUpdatePropNetworkingEnabled()
}
func (m *Manager) setPropNetworkingEnabled() {
	if m.NetworkingEnabled == nmManager.NetworkingEnabled.Get() {
		return
	}
	m.NetworkingEnabled = nmManager.NetworkingEnabled.Get()
	// setup global device switches
	if m.NetworkingEnabled {
		m.restoreGlobalDeviceSwitches()
	} else {
		m.saveAndTurnOffGlobalDeviceSwitches()
	}
	m.doUpdatePropNetworkingEnabled()
}
func (m *Manager) doUpdatePropNetworkingEnabled() {
	m.NetworkingEnabled = nmManager.NetworkingEnabled.Get()
	dbus.NotifyChange(m, "NetworkingEnabled")
}
func (m *Manager) restoreGlobalDeviceSwitches() {
	nmSetWirelessEnabled(m.config.LastWirelessEnabled)
	nmSetWwanEnabled(m.config.LastWwanEnabled)
	m.setPropWiredEnabled(m.config.LastWiredEnabled)
	m.setPropVpnEnabled(m.config.LastVpnEnabled)
}
func (m *Manager) saveAndTurnOffGlobalDeviceSwitches() {
	m.config.setLastWirelessEnabled(m.WirelessEnabled)
	m.config.setLastWwanEnabled(m.WwanEnabled)
	m.config.setLastWiredEnabled(m.WiredEnabled)
	m.config.setLastVpnEnabled(m.VpnEnabled)
	m.doTurnOffGlobalDeviceSwitches()
}
func (m *Manager) doTurnOffGlobalDeviceSwitches() {
	nmSetWirelessEnabled(false)
	nmSetWwanEnabled(false)
	m.setPropWiredEnabled(false)
	m.setPropVpnEnabled(false)
}

func (m *Manager) initPropWirelessEnabled() {
	m.WirelessEnabled = nmManager.WirelessEnabled.Get()
	for _, devPath := range nmGetSpecialDevices(NM_DEVICE_TYPE_WIFI) {
		if m.WirelessEnabled {
			m.doEnableDevice(devPath, m.config.getDeviceEnabled(devPath))
		} else {
			m.doEnableDevice(devPath, false)
		}
	}
	m.doUpdatePropWirelessEnabled()
}
func (m *Manager) setPropWirelessEnabled() {
	if m.WirelessEnabled == nmManager.WirelessEnabled.Get() {
		return
	}
	m.WirelessEnabled = nmManager.WirelessEnabled.Get()
	logger.Debug("setPropWirelessEnabled", m.WirelessEnabled)
	// setup wireless devices switches
	for _, devPath := range nmGetSpecialDevices(NM_DEVICE_TYPE_WIFI) {
		if m.WirelessEnabled {
			m.restoreDeviceState(devPath)
		} else {
			m.saveAndDisconnectDevice(devPath)
		}
	}
	m.doUpdatePropWirelessEnabled()
}
func (m *Manager) doUpdatePropWirelessEnabled() {
	m.WirelessEnabled = nmManager.WirelessEnabled.Get()
	dbus.NotifyChange(m, "WirelessEnabled")
}

func (m *Manager) initPropWwanEnabled() {
	m.WwanEnabled = nmManager.WwanEnabled.Get()
	for _, devPath := range nmGetSpecialDevices(NM_DEVICE_TYPE_MODEM) {
		if m.WwanEnabled {
			m.doEnableDevice(devPath, m.config.getDeviceEnabled(devPath))
		} else {
			m.doEnableDevice(devPath, false)
		}
	}
	m.doUpdatePropWwanEnabled()
}
func (m *Manager) setPropWwanEnabled() {
	if m.WwanEnabled == nmManager.WwanEnabled.Get() {
		return
	}
	m.WwanEnabled = nmManager.WwanEnabled.Get()
	// setup modem devices switches
	for _, devPath := range nmGetSpecialDevices(NM_DEVICE_TYPE_MODEM) {
		if m.WwanEnabled {
			m.restoreDeviceState(devPath)
		} else {
			m.saveAndDisconnectDevice(devPath)
		}
	}
	m.doUpdatePropWwanEnabled()
}
func (m *Manager) doUpdatePropWwanEnabled() {
	m.WwanEnabled = nmManager.WwanEnabled.Get()
	dbus.NotifyChange(m, "WwanEnabled")
}

func (m *Manager) initPropWiredEnabled() {
	m.WiredEnabled = m.config.WiredEnabled
	for _, devPath := range nmGetSpecialDevices(NM_DEVICE_TYPE_ETHERNET) {
		if m.WiredEnabled {
			m.doEnableDevice(devPath, m.config.getDeviceEnabled(devPath))
		} else {
			m.doEnableDevice(devPath, false)
		}
	}
	m.doUpdatePropWiredEnabled()
}
func (m *Manager) setPropWiredEnabled(enabled bool) {
	if m.config.WiredEnabled == enabled {
		return
	}
	logger.Debug("setPropWiredEnabled", enabled)
	m.WiredEnabled = enabled
	m.config.setWiredEnabled(enabled)
	// setup wired devices switches
	for _, devPath := range nmGetSpecialDevices(NM_DEVICE_TYPE_ETHERNET) {
		if enabled {
			m.restoreDeviceState(devPath)
		} else {
			m.saveAndDisconnectDevice(devPath)
		}
		// FIXME keep wired device state same with WiredEnabled to avoid issue
		m.config.setDeviceEnabled(devPath, enabled)
	}
	m.doUpdatePropWiredEnabled()
}
func (m *Manager) doUpdatePropWiredEnabled() {
	dbus.NotifyChange(m, "WiredEnabled")
}

func (m *Manager) initPropVpnEnabled() {
	m.VpnEnabled = m.config.VpnEnabled
	// TODO
	// for _, uuid := range nmGetSpecialConnectionUuids(NM_SETTING_VPN_SETTING_NAME) {
	// if m.VpnEnabled {
	// 	// enable vpn if is autoconnect
	// 	m.ActivateConnection(uuid, "/")
	// }
	// }
	m.doUpdatePropVpnEnabled()
}
func (m *Manager) setPropVpnEnabled(enabled bool) {
	if m.config.VpnEnabled == enabled {
		return
	}
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
	m.doUpdatePropVpnEnabled()
}
func (m *Manager) doUpdatePropVpnEnabled() {
	dbus.NotifyChange(m, "VpnEnabled")
}

func (m *Manager) setPropActiveConnections() {
	m.ActiveConnections, _ = marshalJSON(m.activeConnections)
	dbus.NotifyChange(m, "ActiveConnections")
	// logger.Debug("ActiveConnections:", m.ActiveConnections) // TODO test
}

func (m *Manager) setPropState() {
	m.State = nmGetManagerState()
	dbus.NotifyChange(m, "State")
}

func (m *Manager) setPropWiredDevices() {
	dbus.NotifyChange(m, "WiredDevices")
}
func (m *Manager) setPropWirelessDevices() {
	dbus.NotifyChange(m, "WirelessDevices")
}
func (m *Manager) setPropDevices() {
	m.Devices, _ = marshalJSON(m.devices)
	dbus.NotifyChange(m, "Devices")
	// logger.Debug("setPropDevices", m.Devices) // TODO test
}

// TODO
func (m *Manager) setPropAccessPoints() {
	// m.AccessPoints, _ = marshalJSON(m.accessPoints)
	// dbus.NotifyChange(m, "AccessPoints")
	// testJSON, _ := marshalJSON(m.accessPoints)
	// logger.Debug("setPropAccessPoints", testJSON) // TODO test
}

func (m *Manager) setPropConnections() {
	m.Connections, _ = marshalJSON(m.connections)
	// logger.Debug(m.Connections) // TODO test
	dbus.NotifyChange(m, "Connections")
}
