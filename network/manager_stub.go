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
			logger.Warning("NetworkingEnabled already", m.NetworkingEnabled)
		}
	}
}

func (m *Manager) updatePropNetworkingEnabled() {
	m.NetworkingEnabled = nmManager.NetworkingEnabled.Get()
	dbus.NotifyChange(m, "NetworkingEnabled")
}

func (m *Manager) updatePropActiveConnections() {
	m.ActiveConnections, _ = marshalJSON(m.activeConnections)
	dbus.NotifyChange(m, "ActiveConnections")
	// logger.Debug("ActiveConnections:", m.ActiveConnections) // TODO test
}

func (m *Manager) updatePropState() {
	m.State = nmGetState()
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
