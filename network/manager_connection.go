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
	. "pkg.linuxdeepin.com/lib/gettext"
)

// TODO different connection structures for different types
type connection struct {
	Path dbus.ObjectPath
	Uuid string
	Id   string

	// if not empty, the connection will only apply to special device,
	// works for wired, wireless, infiniband, wimax devices
	HwAddress string

	// works for wireless, olpc-mesh connections
	Ssid string
}

func (m *Manager) initConnectionManage() {
	m.initConnections()
	nmSettings.ConnectNewConnection(func(path dbus.ObjectPath) {
		m.handleConnectionChanged(opAdded, path)
	})
}
func (m *Manager) initConnections() {
	m.connectionsLocker.Lock()
	m.connections = make(map[string][]*connection)
	m.connectionsLocker.Unlock()
	for _, c := range nmGetConnectionList() {
		m.handleConnectionChanged(opAdded, c)
	}
}

func (m *Manager) handleConnectionChanged(operation int32, path dbus.ObjectPath) {
	m.connectionsLocker.Lock()
	defer m.connectionsLocker.Unlock()

	switch operation {
	case opAdded:
		nmConn, _ := nmNewSettingsConnection(path)
		nmConn.ConnectRemoved(func() {
			m.handleConnectionChanged(opRemoved, path)
			nmDestroySettingsConnection(nmConn)
		})
		nmConn.ConnectUpdated(func() {
			m.handleConnectionChanged(opUpdated, path)
		})

		conn := m.newConnection(path)
		cdata, err := nmConn.GetSettings()
		if err != nil {
			return
		}
		switch getSettingConnectionType(cdata) {
		case NM_SETTING_WIRED_SETTING_NAME:
			// wired connection will be treatment specially
		case NM_SETTING_WIRELESS_SETTING_NAME:
			switch getCustomConnectionType(cdata) {
			case connectionWireless:
				m.connections[connectionWireless] = m.addConnection(m.connections[connectionWireless], conn)
			case connectionWirelessAdhoc:
				m.connections[connectionWirelessAdhoc] = m.addConnection(m.connections[connectionWirelessAdhoc], conn)
			case connectionWirelessHotspot:
				m.connections[connectionWirelessHotspot] = m.addConnection(m.connections[connectionWirelessHotspot], conn)
			}
		case NM_SETTING_PPPOE_SETTING_NAME:
			m.connections[connectionPppoe] = m.addConnection(m.connections[connectionPppoe], conn)
		case NM_SETTING_GSM_SETTING_NAME, NM_SETTING_CDMA_SETTING_NAME:
			m.connections[connectionMobile] = m.addConnection(m.connections[connectionMobile], conn)
		case NM_SETTING_VPN_SETTING_NAME:
			m.connections[connectionVpn] = m.addConnection(m.connections[connectionVpn], conn)
		}
	case opRemoved:
		conn := &connection{Path: path}
		for k, conns := range m.connections {
			if m.isConnectionExists(conns, conn) {
				m.connections[k] = m.removeConnection(conns, conn)
			}
		}
	case opUpdated:
		conn := m.newConnection(path)
		for k, conns := range m.connections {
			if m.isConnectionExists(conns, conn) {
				m.connections[k] = m.updateConnection(conns, conn)
			}
		}
	}
	m.setPropConnections()
}

func (m *Manager) newConnection(path dbus.ObjectPath) (conn *connection) {
	conn = &connection{Path: path}

	cdata, err := nmGetConnectionData(path)
	if err != nil {
		return
	}

	conn.Uuid = getSettingConnectionUuid(cdata)
	conn.Id = getSettingConnectionId(cdata)

	switch getSettingConnectionType(cdata) {
	case NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_PPPOE_SETTING_NAME:
		if isSettingWiredMacAddressExists(cdata) {
			conn.HwAddress = convertMacAddressToString(getSettingWiredMacAddress(cdata))
		}
	case NM_SETTING_WIRELESS_SETTING_NAME:
		conn.Ssid = string(getSettingWirelessSsid(cdata))
		if isSettingWirelessMacAddressExists(cdata) {
			conn.HwAddress = convertMacAddressToString(getSettingWirelessMacAddress(cdata))
		}
	}
	return
}

func (m *Manager) clearConnections() {
	m.connectionsLocker.Lock()
	defer m.connectionsLocker.Unlock()
	m.connections = make(map[string][]*connection)
	m.setPropConnections()
}
func (m *Manager) addConnection(conns []*connection, conn *connection) []*connection {
	if m.isConnectionExists(conns, conn) {
		return conns
	}
	if nmGetConnectionType(conn.Path) == NM_SETTING_VPN_SETTING_NAME {
		m.config.addVpnConfig(conn.Uuid)
	}
	conns = append(conns, conn)
	return conns
}
func (m *Manager) removeConnection(conns []*connection, conn *connection) []*connection {
	i := m.getConnectionIndex(conns, conn)
	if i < 0 {
		return conns
	}
	m.config.removeConnection(conns[i].Uuid)
	copy(conns[i:], conns[i+1:])
	conns = conns[:len(conns)-1]
	return conns
}
func (m *Manager) updateConnection(conns []*connection, conn *connection) []*connection {
	i := m.getConnectionIndex(conns, conn)
	if i < 0 {
		return conns
	}
	conns[i] = conn
	return conns
}
func (m *Manager) isConnectionExists(conns []*connection, conn *connection) bool {
	if m.getConnectionIndex(conns, conn) >= 0 {
		return true
	}
	return false
}
func (m *Manager) getConnectionIndex(conns []*connection, conn *connection) int {
	for i, c := range conns {
		if c.Path == conn.Path {
			return i
		}
	}
	return -1
}

// GetSupportedConnectionTypes return all supported connection types
func (m *Manager) GetSupportedConnectionTypes() (types []string) {
	return supportedConnectionTypes
}

// GetWiredConnectionUuid return connection uuid for target wired device.
func (m *Manager) GetWiredConnectionUuid(wiredDevPath dbus.ObjectPath) (uuid string) {
	// this interface will be called by front-end always if user try
	// to connect or edit the wired connection, so ensure the
	// connection exists here is a good choice
	m.ensureWiredConnectionExists(wiredDevPath)
	uuid = nmGeneralGetDeviceRelatedUuid(wiredDevPath)
	return
}

func (m *Manager) ensureWiredConnectionExists(wiredDevPath dbus.ObjectPath) {
	// check if wired connection for target device exists, if not, create one
	uuid := nmGeneralGetDeviceRelatedUuid(wiredDevPath)
	var id string
	if nmGeneralIsUsbDevice(wiredDevPath) {
		id = nmGeneralGetDeviceVendor(wiredDevPath)
	} else {
		id = Tr("Wired Connection")
	}
	if cpath, err := nmGetConnectionByUuid(uuid); err != nil {
		// connection not exists, create one
		hwAddr, _ := nmGeneralGetDeviceHwAddr(wiredDevPath)
		newWiredConnectionForDevice(id, uuid, hwAddr)
	} else {
		// connection already exists, reset its name to keep
		// consistent with current system's language
		nmSetConnectionId(cpath, id)
	}
	return
}

// CreateConnection create a new connection, return ConnectionSession's dbus object path if success.
func (m *Manager) CreateConnection(connType string, devPath dbus.ObjectPath) (session *ConnectionSession, err error) {
	logger.Debug("CreateConnection", connType, devPath)
	session, err = newConnectionSessionByCreate(connType, devPath)
	if err != nil {
		logger.Error(err)
		return
	}
	m.addConnectionSession(session)
	return
}

// EditConnection open a connection through uuid, return ConnectionSession's dbus object path if success.
func (m *Manager) EditConnection(uuid string, devPath dbus.ObjectPath) (session *ConnectionSession, err error) {
	session, err = newConnectionSessionByOpen(uuid, devPath)
	if err != nil {
		logger.Error(err)
		return
	}
	m.addConnectionSession(session)
	return
}

func (m *Manager) addConnectionSession(session *ConnectionSession) {
	m.connectionSessionsLocker.Lock()
	defer m.connectionSessionsLocker.Unlock()

	// install dbus session
	err := dbus.InstallOnSession(session)
	if err != nil {
		logger.Error(err)
		return
	}
	m.connectionSessions = append(m.connectionSessions, session)
}
func (m *Manager) removeConnectionSession(session *ConnectionSession) {
	m.connectionSessionsLocker.Lock()
	defer m.connectionSessionsLocker.Unlock()

	dbus.UnInstallObject(session)

	i := m.getConnectionSessionIndex(session)
	if i < 0 {
		logger.Warning("connection session index is -1", session.sessionPath)
		return
	}

	copy(m.connectionSessions[i:], m.connectionSessions[i+1:])
	newlen := len(m.connectionSessions) - 1
	m.connectionSessions[newlen] = nil
	m.connectionSessions = m.connectionSessions[:newlen]
}
func (m *Manager) getConnectionSessionIndex(session *ConnectionSession) int {
	for i, s := range m.connectionSessions {
		if s.sessionPath == session.sessionPath {
			return i
		}
	}
	return -1
}
func (m *Manager) clearConnectionSessions() {
	m.connectionSessionsLocker.Lock()
	defer m.connectionSessionsLocker.Unlock()

	for _, session := range m.connectionSessions {
		dbus.UnInstallObject(session)
	}
	m.connectionSessions = nil
}

// DeleteConnection delete a connection through uuid.
func (m *Manager) DeleteConnection(uuid string) (err error) {
	//TODO: remove(uninstall dbus) editing connection_session object
	cpath, err := nmGetConnectionByUuid(uuid)
	if err != nil {
		return
	}
	conn, err := nmNewSettingsConnection(cpath)
	if err != nil {
		return err
	}
	return conn.Delete()
}

func (m *Manager) ActivateConnection(uuid string, devPath dbus.ObjectPath) (cpath dbus.ObjectPath, err error) {
	logger.Debugf("ActivateConnection: uuid=%s, devPath=%s", uuid, devPath)
	cpath, err = nmGetConnectionByUuid(uuid)
	if err != nil {
		return
	}
	_, err = nmActivateConnection(cpath, devPath)
	return
}

func (m *Manager) DeactivateConnection(uuid string) (err error) {
	apaths, err := nmGetActiveConnectionByUuid(uuid)
	if err != nil {
		// not found active connection with uuid, ignore error here
		return
	}
	for _, apath := range apaths {
		logger.Debug("DeactivateConnection:", uuid, apath)
		if isConnectionStateInActivating(nmGetActiveConnectionState(apath)) {
			if tmpErr := nmDeactivateConnection(apath); tmpErr != nil {
				err = tmpErr
			}
		}
	}
	return
}

// DisconnectDevice will disconnect all connection in target device.
func (m *Manager) DisconnectDevice(devPath dbus.ObjectPath) (err error) {
	return m.doDisconnectDevice(devPath)
}
func (m *Manager) doDisconnectDevice(devPath dbus.ObjectPath) (err error) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	devState := dev.State.Get()
	if isDeviceStateInActivating(devState) {
		err = dev.Disconnect()
		if err != nil {
			logger.Error(err)
		}
	}
	return
}
