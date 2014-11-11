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
	"sync"
)

const (
	dbusNetworkDest = "com.deepin.daemon.Network"
	dbusNetworkPath = "/com/deepin/daemon/Network"
	dbusNetworkIfs  = "com.deepin.daemon.Network"
)

const (
	opAdded = iota
	opRemoved
	opUpdated
)

type connectionData map[string]map[string]dbus.Variant

type Manager struct {
	config *config

	// update by manager.go
	State uint32 // global networking state

	NetworkingEnabled bool `access:"readwrite"` // airplane mode for NetworkManager
	WirelessEnabled   bool `access:"readwrite"`
	WwanEnabled       bool `access:"readwrite"`
	WiredEnabled      bool `access:"readwrite"`
	VpnEnabled        bool `access:"readwrite"`

	// update by manager_devices.go
	devicesLocker sync.Mutex
	devices       map[string][]*device
	Devices       string // array of device objects and marshaled by json

	accessPointsLocker sync.Mutex
	accessPoints       map[dbus.ObjectPath][]*accessPoint

	// update by manager_connections.go
	connectionsLocker sync.Mutex
	connections       map[string][]*connection
	Connections       string // array of connection information and marshaled by json

	connectionSessionsLocker sync.Mutex
	connectionSessions       []*ConnectionSession

	// update by manager_active.go
	activeConnectionsLocker sync.Mutex
	activeConnections       map[dbus.ObjectPath]*activeConnection
	ActiveConnections       string // array of connections that activated and marshaled by json

	// signals
	NeedSecrets                  func(string, string, string)
	DeviceStateChanged           func(devPath string, newState uint32) // TODO remove
	AccessPointAdded             func(devPath, apJSON string)
	AccessPointRemoved           func(devPath, apJSON string)
	AccessPointPropertiesChanged func(devPath, apJSON string)
	DeviceEnabled                func(devPath string, enabled bool)

	agent         *agent
	stateNotifier *stateNotifier
	dbusWatcher   *dbusWatcher
}

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusNetworkDest,
		ObjectPath: dbusNetworkPath,
		Interface:  dbusNetworkIfs,
	}
}

// initialize slice code
func initSlices() {
	initNmDbusObjects()
	initProxyGsettings()
	initAvailableValuesSecretFlags()
	initAvailableValuesNmPptpSecretFlags()
	initAvailableValuesNmL2tpSecretFlags()
	initAvailableValuesNmVpncSecretFlags()
	initAvailableValuesNmOpenvpnSecretFlags()
	initAvailableValuesWirelessChannel()
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
	destroyAgent(m.agent)
	destroyStateNotifier(m.stateNotifier)
	destroyDbusWatcher(m.dbusWatcher)
	m.clearConnectionSessions()
	dbus.UnInstallObject(m)
}

func (m *Manager) initManager() {
	m.dbusWatcher = newDbusWatcher(true)

	// setup global switches
	m.initPropNetworkingEnabled()
	nmManager.NetworkingEnabled.ConnectChanged(func() {
		m.setPropNetworkingEnabled()
	})
	m.initPropWirelessEnabled()
	nmManager.WirelessEnabled.ConnectChanged(func() {
		m.setPropWirelessEnabled()
	})
	m.initPropWwanEnabled()
	nmManager.WwanEnabled.ConnectChanged(func() {
		m.setPropWwanEnabled()
	})

	// load virtual global switches information from configuration file
	m.initPropWiredEnabled()
	m.initPropVpnEnabled()

	// initialize device and connection handlers
	m.initDeviceManage()
	m.initConnectionManage()
	m.initActiveConnectionManage()

	// update property "State"
	nmManager.State.ConnectChanged(func() {
		m.setPropState()
	})
	m.setPropState()

	m.stateNotifier = newStateNotifier()
	m.agent = newAgent()
}
