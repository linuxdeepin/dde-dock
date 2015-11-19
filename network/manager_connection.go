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
	nm "dbus/org/freedesktop/networkmanager"
	"pkg.deepin.io/lib/dbus"
	. "pkg.deepin.io/lib/gettext"
)

// TODO refactor code, different connection structures for different
// types
type connection struct {
	nmConn   *nm.SettingsConnection
	connType string

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
	m.connectionsLock.Lock()
	m.connections = make(map[string][]*connection)
	m.connectionsLock.Unlock()

	for _, cpath := range nmGetConnectionList() {
		m.addConnection(cpath)
	}
	nmSettings.ConnectNewConnection(func(cpath dbus.ObjectPath) {
		logger.Info("add connection", cpath)
		m.addConnection(cpath)
	})
	nmSettings.ConnectConnectionRemoved(func(cpath dbus.ObjectPath) {
		logger.Info("remove connection", cpath)
		m.removeConnection(cpath)
	})
}

func (m *Manager) newConnection(cpath dbus.ObjectPath) (conn *connection, err error) {
	conn = &connection{Path: cpath}
	nmConn, err := nmNewSettingsConnection(cpath)
	if err != nil {
		return
	}

	conn.nmConn = nmConn
	conn.updateProps()

	if conn.connType == connectionVpn {
		m.config.addVpnConfig(conn.Uuid)
	}

	// connect signals
	nmConn.ConnectUpdated(func() {
		m.updateConnection(cpath)
	})

	return
}
func (conn *connection) updateProps() {
	cdata, err := conn.nmConn.GetSettings()
	if err != nil {
		logger.Error(err)
		return
	}

	conn.Uuid = getSettingConnectionUuid(cdata)
	conn.Id = getSettingConnectionId(cdata)

	switch getSettingConnectionType(cdata) {
	case NM_SETTING_GSM_SETTING_NAME, NM_SETTING_CDMA_SETTING_NAME:
		conn.connType = connectionMobile
	case NM_SETTING_VPN_SETTING_NAME:
		conn.connType = connectionVpn
	default:
		conn.connType = getCustomConnectionType(cdata)
	}

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
}

func (m *Manager) destroyConnection(conn *connection) {
	m.config.removeConnection(conn.Uuid)
	nmDestroySettingsConnection(conn.nmConn)
}

func (m *Manager) clearConnections() {
	m.connectionsLock.Lock()
	defer m.connectionsLock.Unlock()
	for _, conns := range m.connections {
		for _, conn := range conns {
			m.destroyConnection(conn)
		}
	}
	m.connections = make(map[string][]*connection)
	m.setPropConnections()
}

func (m *Manager) addConnection(cpath dbus.ObjectPath) {
	if m.isConnectionExists(cpath) {
		logger.Warning("connection already exists", cpath)
		return
	}

	m.connectionsLock.Lock()
	defer m.connectionsLock.Unlock()
	conn, err := m.newConnection(cpath)
	if err != nil {
		return
	}
	logger.Infof("add connection %#v", conn)
	switch conn.connType {
	case connectionWired, connectionUnknown:
		// wired connections will be special treatment, do not shown
		// for front-end here
	default:
		m.connections[conn.connType] = append(m.connections[conn.connType], conn)
	}
	m.setPropConnections()
}

func (m *Manager) removeConnection(cpath dbus.ObjectPath) {
	if !m.isConnectionExists(cpath) {
		logger.Warning("connection not found", cpath)
		return
	}
	connType, i := m.getConnectionIndex(cpath)

	m.connectionsLock.Lock()
	defer m.connectionsLock.Unlock()
	m.connections[connType] = m.doRemoveConnection(m.connections[connType], i)
	m.setPropConnections()
}
func (m *Manager) doRemoveConnection(conns []*connection, i int) []*connection {
	logger.Infof("remove connection %#v", conns[i])
	m.destroyConnection(conns[i])
	copy(conns[i:], conns[i+1:])
	conns = conns[:len(conns)-1]
	return conns
}

func (m *Manager) updateConnection(cpath dbus.ObjectPath) {
	if !m.isConnectionExists(cpath) {
		logger.Warning("connection not found", cpath)
		return
	}
	conn := m.getConnection(cpath)

	m.connectionsLock.Lock()
	defer m.connectionsLock.Unlock()
	conn.updateProps()
	logger.Infof("update connection %#v", conn)
	m.setPropConnections()
}

func (m *Manager) getConnection(cpath dbus.ObjectPath) (conn *connection) {
	connType, i := m.getConnectionIndex(cpath)
	if i < 0 {
		logger.Warning("connection not found", cpath)
		return
	}

	m.connectionsLock.Lock()
	defer m.connectionsLock.Unlock()
	conn = m.connections[connType][i]
	return
}
func (m *Manager) isConnectionExists(cpath dbus.ObjectPath) bool {
	_, i := m.getConnectionIndex(cpath)
	if i >= 0 {
		return true
	}
	return false
}
func (m *Manager) getConnectionIndex(cpath dbus.ObjectPath) (connType string, index int) {
	m.connectionsLock.Lock()
	defer m.connectionsLock.Unlock()
	for t, conns := range m.connections {
		for i, c := range conns {
			if c.Path == cpath {
				return t, i
			}
		}
	}
	return "", -1
}

// GetSupportedConnectionTypes return all supported connection types
func (m *Manager) GetSupportedConnectionTypes() (types []string) {
	return supportedConnectionTypes
}

// TODO: remove, use device.UniqueUuid instead
// GetWiredConnectionUuid return connection uuid for target wired device.
func (m *Manager) GetWiredConnectionUuid(wiredDevPath dbus.ObjectPath) (uuid string) {
	// this interface will be called by front-end always if user try
	// to connect or edit the wired connection, so ensure the
	// connection exists here is a good choice
	m.ensureWiredConnectionExists(wiredDevPath, false)
	uuid = nmGeneralGetDeviceUniqueUuid(wiredDevPath)
	return
}

func (m *Manager) generalEnsureUniqueConnectionExists(devPath dbus.ObjectPath, active bool) (cpath dbus.ObjectPath, exists bool, err error) {
	switch nmGetDeviceType(devPath) {
	case NM_DEVICE_TYPE_ETHERNET:
		cpath, exists, err = m.ensureWiredConnectionExists(devPath, active)
	case NM_DEVICE_TYPE_WIFI:
	case NM_DEVICE_TYPE_MODEM:
		cpath, exists, err = m.ensureMobileConnectionExists(devPath, active)
	}
	return
}

// ensureWiredConnectionExists will check if wired connection for
// target device exists, if not, create one.
func (m *Manager) ensureWiredConnectionExists(wiredDevPath dbus.ObjectPath, active bool) (cpath dbus.ObjectPath, exists bool, err error) {
	uuid := nmGeneralGetDeviceUniqueUuid(wiredDevPath)
	var id string
	if nmGeneralIsUsbDevice(wiredDevPath) {
		id = nmGeneralGetDeviceVendor(wiredDevPath)
	} else {
		id = Tr("Wired Connection")
	}
	if cpath, err := nmGetConnectionByUuid(uuid); err != nil {
		// connection not exists, create one
		exists = false
		cpath, err = newWiredConnectionForDevice(id, uuid, wiredDevPath, active)
	} else {
		// connection already exists, reset its name to keep
		// consistent with current system's language
		exists = true
		nmSetConnectionId(cpath, id)
	}
	return
}

// ensureMobileConnectionExists will check if mobile connection for
// target device exists, if not, create one.
func (m *Manager) ensureMobileConnectionExists(modemDevPath dbus.ObjectPath, active bool) (cpath dbus.ObjectPath, exists bool, err error) {
	uuid := nmGeneralGetDeviceUniqueUuid(modemDevPath)
	_, err = nmGetConnectionByUuid(uuid)
	if err == nil {
		// connection already exists
		exists = true
		return
	}
	// connection id will be reset when setup plan, so here just give
	// an optional name like "mobile"
	cpath, err = newMobileConnectionForDevice("mobile", uuid, modemDevPath, active)
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
	if isNmObjectPathValid(devPath) && nmGeneralGetDeviceUniqueUuid(devPath) == uuid {
		m.generalEnsureUniqueConnectionExists(devPath, false)
	}
	session, err = newConnectionSessionByOpen(uuid, devPath)
	if err != nil {
		logger.Error(err)
		return
	}
	m.addConnectionSession(session)
	return
}

func (m *Manager) addConnectionSession(session *ConnectionSession) {
	m.connectionSessionsLock.Lock()
	defer m.connectionSessionsLock.Unlock()

	// install dbus session
	err := dbus.InstallOnSession(session)
	if err != nil {
		logger.Error(err)
		return
	}
	m.connectionSessions = append(m.connectionSessions, session)
}
func (m *Manager) removeConnectionSession(session *ConnectionSession) {
	m.connectionSessionsLock.Lock()
	defer m.connectionSessionsLock.Unlock()

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
	m.connectionSessionsLock.Lock()
	defer m.connectionSessionsLock.Unlock()
	for _, session := range m.connectionSessions {
		dbus.UnInstallObject(session)
	}
	m.connectionSessions = nil
}

// DeleteConnection delete a connection through uuid.
func (m *Manager) DeleteConnection(uuid string) (err error) {
	// FIXME: uninstall ConnectionSession dbus object if under editing
	cpath, err := nmGetConnectionByUuid(uuid)
	if err != nil {
		return
	}
	nmConn, err := nmNewSettingsConnection(cpath)
	if err != nil {
		return
	}
	defer nmDestroySettingsConnection(nmConn)
	return nmConn.Delete()
}

func (m *Manager) ActivateConnection(uuid string, devPath dbus.ObjectPath) (cpath dbus.ObjectPath, err error) {
	logger.Debugf("ActivateConnection: uuid=%s, devPath=%s", uuid, devPath)
	if isNmObjectPathValid(devPath) && nmGeneralGetDeviceUniqueUuid(devPath) == uuid {
		var exists bool
		cpath, exists, err = m.generalEnsureUniqueConnectionExists(devPath, true)
		if !exists {
			// connection will be activated in
			// generalEnsureUniqueConnectionExists() if not exists
			return
		}
	}
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
	nmDev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	defer nmDestroyDevice(nmDev)

	devState := nmDev.State.Get()
	if isDeviceStateInActivating(devState) {
		err = nmDev.Disconnect()
		if err != nil {
			logger.Error(err)
		}
	}
	return
}
