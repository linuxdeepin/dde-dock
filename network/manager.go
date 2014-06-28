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

const (
	dbusNetworkDest = "com.deepin.daemon.Network"
	dbusNetworkPath = "/com/deepin/daemon/Network"
	dbusNetworkIfs  = "com.deepin.daemon.Network"
)

// TODO refactor code
const (
	opAdded = iota
	opRemoved
)

type connectionData map[string]map[string]dbus.Variant

type Manager struct {
	config *config

	// update by manager.go
	State             uint32 // networking state
	activeConnections []*activeConnection
	ActiveConnections string // array of connections that activated and marshaled by json

	NetworkingEnabled bool `access:"readwrite"` // airplane mode for NetworkManager
	WirelessEnabled   bool `access:"readwrite"`
	WwanEnabled       bool `access:"readwrite"`
	WiredEnabled      bool `access:"readwrite"`
	VpnEnabled        bool `access:"readwrite"`

	// update by manager_devices.go
	devices      map[string][]*device
	Devices      string // array of device objects and marshaled by json
	accessPoints map[dbus.ObjectPath][]*accessPoint
	// AccessPoints    string // TODO array of access point objects and marshaled by json

	// update by manager_connections.go
	connectionSessions []*ConnectionSession
	connections        map[string][]*connection
	Connections        string // array of connection information and marshaled by json

	// signals
	NeedSecrets                  func(string, string, string)
	DeviceStateChanged           func(devPath string, newState uint32) // TODO remove
	AccessPointAdded             func(devPath, apJSON string)
	AccessPointRemoved           func(devPath, apJSON string)
	AccessPointPropertiesChanged func(devPath, apJSON string)
	DeviceEnabled                func(devPath string, enabled bool)

	agent         *Agent
	stateNotifier *StateNotifier
}

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		dbusNetworkDest,
		dbusNetworkPath,
		dbusNetworkIfs,
	}
}

// initialize slice code
func initSlices() {
	initAvailableValues8021x()
	initAvailableValuesIp4()
	initAvailableValuesIp6()
	initNmStateReasons()
}

func NewManager() (m *Manager) {
	m = &Manager{}
	m.config = newConfig()
	return
}

func DestroyManager(m *Manager) {
	destroyStateNotifier(m.stateNotifier)
	destroyAgent(m.agent)
	m.clearConnectionSessions()
	dbus.UnInstallObject(m)
}

func (m *Manager) initManager() {
	m.initDeviceManage()
	m.initConnectionManage()

	// setup global switches
	m.updatePropNetworkingEnabled()
	nmManager.NetworkingEnabled.ConnectChanged(func() {
		m.updatePropNetworkingEnabled()
	})

	m.updatePropWirelessEnabled()
	nmManager.WirelessEnabled.ConnectChanged(func() {
		m.updatePropWirelessEnabled()
	})

	m.updatePropWwanEnabled()
	nmManager.WwanEnabled.ConnectChanged(func() {
		m.updatePropWwanEnabled()
	})

	// load virtual global switches information from configuration file
	m.updatePropWiredEnabled(m.config.WiredEnabled)
	m.updatePropVpnEnabled(m.config.VpnEnabled)

	// update property "ActiveConnections" after devices initialized
	m.updateActiveConnections()
	nmManager.ActiveConnections.ConnectChanged(func() {
		m.updateActiveConnections()
	})

	// update property "State"
	m.updatePropState()
	nmManager.State.ConnectChanged(func() {
		m.updatePropState()
	})

	m.agent = newAgent()
	m.stateNotifier = newStateNotifier()
}

func (m *Manager) updateActiveConnections() {
	// reset all exists active connection objects
	for i, _ := range m.activeConnections {
		// destroy object to reset all property connects
		nmDestroyActiveConnection(m.activeConnections[i].nmaconn)
		m.activeConnections[i] = nil
	}
	m.activeConnections = make([]*activeConnection, 0)
	for _, acPath := range nmGetActiveConnections() {
		if nmaconn, err := nmNewActiveConnection(acPath); err == nil {
			aconn := &activeConnection{
				nmaconn: nmaconn,
				path:    acPath,
				Devices: nmaconn.Devices.Get(),
				Uuid:    nmaconn.Uuid.Get(),
				State:   nmaconn.State.Get(),
				Vpn:     nmaconn.Vpn.Get(),
			}
			nmaconn.State.ConnectChanged(func() {
				// TODO fix dbus property issue
				logger.Debug("active connection state changed:", aconn.State, nmaconn.State.Get())
				aconn.State = nmaconn.State.Get()
				if aconn.Vpn {
					// TODO notification for vpn
					// update vpn config
					m.config.setVpnConnectionActivated(aconn.Uuid, isConnectionStateActivated(aconn.State))
				}
				m.updatePropActiveConnections()
			})
			m.activeConnections = append(m.activeConnections, aconn)
		}
	}
	m.updatePropActiveConnections()
	// logger.Debug("active connections changed:", m.ActiveConnections) // TODO test
}

// TODO remove
func (m *Manager) isActiveConnectionExists(aconn *activeConnection) bool {
	if aconn == nil {
		return false
	}
	for _, a := range m.activeConnections {
		if &a == &aconn {
			return true
		}
	}
	return false
}
