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

type activeConnection struct {
	path dbus.ObjectPath

	Devices []dbus.ObjectPath
	Id      string
	Uuid    string
	State   uint32
	Vpn     bool
}

func (m *Manager) initActiveConnectionManage() {
	m.initActiveConnections()

	// custom dbus watcher to catch all signals about active
	// connection, including vpn connection
	senderNm := "org.freedesktop.NetworkManager"
	interfaceDbusProperties := "org.freedesktop.DBus.Properties"
	interfaceActive := "org.freedesktop.NetworkManager.Connection.Active"
	interfaceVpn := "org.freedesktop.NetworkManager.VPN.Connection"
	memberProperties := "PropertiesChanged"
	memberVpnState := "VpnStateChanged"
	m.dbusWatcher.watch("type=signal,sender=" + senderNm + ",interface=" + interfaceDbusProperties + ",member=" + memberProperties)
	m.dbusWatcher.watch("type=signal,sender=" + senderNm + ",interface=" + interfaceActive + ",member=" + memberProperties)
	m.dbusWatcher.watch("type=signal,sender=" + senderNm + ",interface=" + interfaceVpn + ",member=" + memberVpnState)

	// update active connection properties
	m.dbusWatcher.connect(func(s *dbus.Signal) {
		m.activeConnectionsLocker.Lock()
		defer m.activeConnectionsLocker.Unlock()

		var props map[string]dbus.Variant
		if s.Name == interfaceDbusProperties+"."+memberProperties && len(s.Body) >= 2 {
			// compatible with old dbus signal
			if realName, ok := s.Body[0].(string); ok &&
				realName == interfaceActive {
				props, _ = s.Body[1].(map[string]dbus.Variant)
			}
		} else if s.Name == interfaceActive+"."+memberProperties && len(s.Body) >= 1 {
			props, _ = s.Body[0].(map[string]dbus.Variant)
		}
		if props != nil {
			aconn, ok := m.activeConnections[s.Path]
			if !ok {
				aconn = m.newActiveConnection(s.Path)
			}

			// query each properties that changed
			for k, vv := range props {
				if k == "State" {
					aconn.State, _ = vv.Value().(uint32)
				} else if k == "Devices" {
					aconn.Devices, _ = vv.Value().([]dbus.ObjectPath)
				} else if k == "Uuid" {
					aconn.Uuid, _ = vv.Value().(string)
					if cpath, err := nmGetConnectionByUuid(aconn.Uuid); err == nil {
						aconn.Id = nmGetConnectionId(cpath)
					}
				} else if k == "Vpn" {
					aconn.Vpn, _ = vv.Value().(bool)
				} else if k == "Connection" { // ignore
				} else if k == "SpecificObject" { // ignore
				} else if k == "Default" { // ignore
				} else if k == "Default6" { // ignore
				} else if k == "Master" { // ignore
				}
			}

			// use "State" to determine if the active connection is
			// adding or removing, if "State" property is not changed
			// is current sequence, it also means that the active
			// connection already exits
			if isConnectionStateInDeactivating(aconn.State) {
				delete(m.activeConnections, s.Path)
			} else {
				m.activeConnections[s.Path] = aconn
			}
			m.setPropActiveConnections()
		}
	})

	// handle notifications for vpn connection
	m.dbusWatcher.connect(func(s *dbus.Signal) {
		m.activeConnectionsLocker.Lock()
		defer m.activeConnectionsLocker.Unlock()

		if s.Name == interfaceVpn+"."+memberVpnState && len(s.Body) >= 2 {
			state, _ := s.Body[0].(uint32)
			reason, _ := s.Body[1].(uint32)

			// get the corresponding active connection
			aconn, ok := m.activeConnections[s.Path]
			if !ok {
				return
			}

			// update vpn config
			m.config.setVpnConnectionActivated(aconn.Uuid, isVpnConnectionStateInActivating(state))

			// notification for vpn
			if isVpnConnectionStateActivated(state) {
				notifyVpnConnected(aconn.Id)
			} else if isVpnConnectionStateDeactivate(state) {
				notifyVpnDisconnected(aconn.Id)
				delete(m.activeConnections, s.Path)
			} else if isVpnConnectionStateFailed(state) {
				notifyVpnFailed(aconn.Id, reason)
				delete(m.activeConnections, s.Path)
			}
		}
	})
}

func (m *Manager) initActiveConnections() {
	m.activeConnectionsLocker.Lock()
	defer m.activeConnectionsLocker.Unlock()

	m.activeConnections = make(map[dbus.ObjectPath]*activeConnection)
	for _, path := range nmGetActiveConnections() {
		m.activeConnections[path] = m.newActiveConnection(path)
	}

	m.setPropActiveConnections()
}

func (m *Manager) newActiveConnection(path dbus.ObjectPath) (aconn *activeConnection) {
	aconn = &activeConnection{path: path}

	nmAConn, err := nmNewActiveConnection(path)
	if err != nil {
		return
	}

	aconn.State = nmAConn.State.Get()
	aconn.Devices = nmAConn.Devices.Get()
	aconn.Uuid = nmAConn.Uuid.Get()
	aconn.Vpn = nmAConn.Vpn.Get()

	if cpath, err := nmGetConnectionByUuid(aconn.Uuid); err == nil {
		aconn.Id = nmGetConnectionId(cpath)
	}

	return
}
