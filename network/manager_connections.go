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
	"fmt"
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

type activeConnection struct {
	nmAConn   *nm.ActiveConnection
	nmVpnConn *nm.VPNConnection
	path      dbus.ObjectPath

	Devices []dbus.ObjectPath
	// SpecificObject dbus.ObjectPath // TODO
	Id    string
	Uuid  string
	State uint32
	Vpn   bool
	// VpnState uint32 // TODO
}

type activeConnectionInfo struct {
	IsPrimaryConnection bool
	ConnectionType      string
	ConnectionName      string
	Security            string
	DeviceType          string
	DeviceInterface     string
	HwAddress           string
	Speed               string
	Ip4                 ip4ConnectionInfo
	Ip6                 ip6ConnectionInfo
}
type ip4ConnectionInfo struct {
	Address string
	Mask    string
	Route   string
	Dns1    string
	Dns2    string
	Dns3    string
}
type ip6ConnectionInfo struct {
	Address string
	Prefix  string
	Route   string
	Dns1    string
	Dns2    string
	Dns3    string
}

func (m *Manager) initConnectionManage() {
	m.connections = make(map[string][]*connection)

	// TODO create special wired connection if need
	// m.updatePropWiredConnections()

	for _, c := range nmGetConnectionList() {
		m.handleConnectionChanged(opAdded, c)
	}
	nmSettings.ConnectNewConnection(func(path dbus.ObjectPath) {
		m.handleConnectionChanged(opAdded, path)
	})
	// TODO
	// nmSettings.ConnectPropertiesChanged(func(path dbus.ObjectPath) {
	// }
}

func (m *Manager) handleConnectionChanged(operation int32, path dbus.ObjectPath) {
	// logger.Debugf("handleConnectionChanged: operation %d, path %s", operation, path) // TODO test
	conn := &connection{Path: path}
	switch operation {
	case opAdded:
		nmConn, _ := nmNewSettingsConnection(path)
		nmConn.ConnectRemoved(func() {
			m.handleConnectionChanged(opRemoved, path)
			nmDestroySettingsConnection(nmConn)
		})

		cdata, err := nmConn.GetSettings()
		if err != nil {
			return
		}
		uuid := getSettingConnectionUuid(cdata)
		conn.Uuid = uuid
		conn.Id = getSettingConnectionId(cdata)

		switch getSettingConnectionType(cdata) {
		case NM_SETTING_WIRED_SETTING_NAME:
			// wired connection will be treatment specially
			// TODO
			// m.WiredConnections = append(m.WiredConnections, uuid)
			// dbus.NotifyChange(m, "WiredConnections")
		case NM_SETTING_WIRELESS_SETTING_NAME:
			conn.Ssid = string(getSettingWirelessSsid(cdata))
			if isSettingWirelessMacAddressExists(cdata) {
				conn.HwAddress = convertMacAddressToString(getSettingWirelessMacAddress(cdata))
			}
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
		m.updatePropConnections()
	case opRemoved:
		for k, conns := range m.connections {
			if m.isConnectionExists(conns, conn) {
				m.connections[k] = m.removeConnection(conns, conn)
			}
		}
		m.updatePropConnections()
	}
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

// TODO GetWiredConnectionUuid return connection uuid for target wired device.
func (m *Manager) GetWiredConnectionUuid(wiredDevPath dbus.ObjectPath) (uuid string) {
	// check if target wired connection exists, if not, create one
	id := Tr("Wired Connection") + " " + nmGetDeviceInterface(wiredDevPath)
	// TODO check connection type, read only
	cpath, ok := nmGetConnectionById(id)
	if ok {
		uuid, _ = nmGetConnectionUuid(cpath)
	} else {
		uuid = newWiredConnection(id)
	}
	return
}

func (m *Manager) GetActiveConnectionInfo() (acinfosJSON string, err error) {
	var acinfos []activeConnectionInfo
	// get activated devices' connection information
	for _, devPath := range nmGetDevices() {
		if isDeviceStateActivated(nmGetDeviceState(devPath)) {
			if info, err := m.doGetActiveConnectionInfo(nmGetDeviceActiveConnection(devPath), devPath); err == nil {
				acinfos = append(acinfos, info)
			}
		}
	}
	// get activated vpn connection information
	for _, apath := range nmGetVpnActiveConnections() {
		if nmAConn, err := nmNewActiveConnection(apath); err == nil {
			if devs := nmAConn.Devices.Get(); len(devs) > 0 {
				devPath := devs[0]
				if info, err := m.doGetActiveConnectionInfo(apath, devPath); err == nil {
					acinfos = append(acinfos, info)
				}
			}
		}
	}
	acinfosJSON, err = marshalJSON(acinfos)
	return
}
func (m *Manager) doGetActiveConnectionInfo(apath, devPath dbus.ObjectPath) (acinfo activeConnectionInfo, err error) {
	var connType, connName, security, devType, devIfc, hwAddress, speed string
	var ip4Address, ip4Mask, ip4Route, ip4Dns1, ip4Dns2, ip4Dns3 string
	var ip6Address, ip6Route, ip6Dns1, ip6Dns2, ip6Dns3 string
	var ip4Info ip4ConnectionInfo
	var ip6Info ip6ConnectionInfo

	// active connection
	nmAConn, err := nmNewActiveConnection(apath)
	if err != nil {
		return
	}
	nmConn, err := nmNewSettingsConnection(nmAConn.Connection.Get())
	if err != nil {
		return
	}

	// device
	nmDev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	devType = getCustomDeviceType(nmDev.DeviceType.Get())
	devIfc = nmDev.Interface.Get()

	// connection data
	hwAddress, _ = nmGeneralGetDeviceHwAddr(devPath)
	speed = nmGeneralGetDeviceSpeed(devPath)

	cdata, err := nmConn.GetSettings()
	if err != nil {
		return
	}
	connName = getSettingConnectionId(cdata)
	connType = getCustomConnectionType(cdata)

	// security
	use8021xSecurity := false
	switch getSettingConnectionType(cdata) {
	case NM_SETTING_WIRED_SETTING_NAME:
		if getSettingVk8021xEnable(cdata) {
			use8021xSecurity = true
		} else {
			security = Tr("None")
		}
	case NM_SETTING_WIRELESS_SETTING_NAME:
		switch getSettingVkWirelessSecurityKeyMgmt(cdata) {
		case "none":
			security = Tr("None")
		case "wep":
			security = Tr("WEP 40/128-bit Key")
		case "wpa-psk":
			security = Tr("WPA/WPA2 Personal")
		case "wpa-eap":
			use8021xSecurity = true
		}
	}
	if use8021xSecurity {
		switch getSettingVk8021xEap(cdata) {
		case "tls":
			security = "EAP/" + Tr("TLS")
		case "md5":
			security = "EAP/" + Tr("MD5")
		case "leap":
			security = "EAP/" + Tr("LEAP")
		case "fast":
			security = "EAP/" + Tr("FAST")
		case "ttls":
			security = "EAP/" + Tr("Tunneled TLS")
		case "peap":
			security = "EAP/" + Tr("Protected EAP")
		}
	}

	// ipv4
	switch getSettingIp4ConfigMethod(cdata) {
	case NM_SETTING_IP4_CONFIG_METHOD_AUTO:
		ip4Address, ip4Mask, ip4Route, ip4Dns1 = nmGetDhcp4Info(nmDev.Dhcp4Config.Get())
		ip4Dns2 = getSettingVkIp4ConfigDns(cdata)
		// ip4Dns2 = getSettingVkIp4ConfigDns2(cdata)
	case NM_SETTING_IP4_CONFIG_METHOD_MANUAL:
		ip4Address = getSettingVkIp4ConfigAddressesAddress(cdata)
		ip4Mask = getSettingVkIp4ConfigAddressesMask(cdata)
		ip4Route = getSettingVkIp4ConfigAddressesGateway(cdata)
		ip4Dns1 = getSettingVkIp4ConfigDns(cdata)
	}
	ip4Info = ip4ConnectionInfo{
		Address: ip4Address,
		Mask:    ip4Mask,
		Route:   ip4Route,
		Dns1:    ip4Dns1,
		Dns2:    ip4Dns2,
		Dns3:    ip4Dns3,
	}

	// ipv6
	if isSettingSectionExists(cdata, sectionIpv6) {
		switch getSettingIp6ConfigMethod(cdata) {
		case NM_SETTING_IP6_CONFIG_METHOD_AUTO, NM_SETTING_IP6_CONFIG_METHOD_DHCP:
			dhcp6Path := nmDev.Dhcp6Config.Get()
			if len(dhcp6Path) > 0 && string(dhcp6Path) != "/" {
				ip6Address, ip6Route, ip6Dns1 = nmGetDhcp6Info(dhcp6Path)
				ip6Dns2 = getSettingVkIp6ConfigDns(cdata)
				ip6Info = ip6ConnectionInfo{
					Address: ip6Address,
					Route:   ip6Route,
					Dns1:    ip6Dns1,
					Dns2:    ip6Dns2,
					Dns3:    ip6Dns3,
				}
			}
		case NM_SETTING_IP6_CONFIG_METHOD_MANUAL:
			ip6Address = getSettingVkIp6ConfigAddressesAddress(cdata)
			ip6Prefix := getSettingVkIp6ConfigAddressesPrefix(cdata)
			ip6Route = getSettingVkIp6ConfigAddressesGateway(cdata)
			ip6Dns1 = getSettingVkIp6ConfigDns(cdata)
			ip6Info = ip6ConnectionInfo{
				Address: fmt.Sprintf("%s/%d", ip6Address, ip6Prefix),
				Route:   ip6Route,
				Dns1:    ip6Dns1,
				Dns2:    ip6Dns2,
				Dns3:    ip6Dns3,
			}
		}
	}

	acinfo = activeConnectionInfo{
		IsPrimaryConnection: nmGetPrimaryConnection() == apath,
		ConnectionType:      connType,
		ConnectionName:      connName,
		Security:            security,
		DeviceType:          devType,
		DeviceInterface:     devIfc,
		HwAddress:           hwAddress,
		Speed:               speed,
		Ip4:                 ip4Info,
		Ip6:                 ip6Info,
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
	// install dbus session
	err := dbus.InstallOnSession(session)
	if err != nil {
		logger.Error(err)
		return
	}
	m.connectionSessions = append(m.connectionSessions, session)
}
func (m *Manager) removeConnectionSession(session *ConnectionSession) {
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

func (m *Manager) ActivateConnection(uuid string, devPath dbus.ObjectPath) (err error) {
	logger.Debugf("ActivateConnection: uuid=%s, devPath=%s", uuid, devPath)
	cpath, err := nmGetConnectionByUuid(uuid)
	if err != nil {
		return
	}
	_, err = nmActivateConnection(cpath, devPath)
	return
}

// TODO
func (m *Manager) DeactivateConnection(uuid string) (err error) {
	apath, ok := nmGetActiveConnectionByUuid(uuid)
	if !ok {
		// not found active connection with uuid
		return
	}
	logger.Debug("DeactivateConnection:", uuid, apath)
	if isConnectionStateInActivating(nmGetActiveConnectionState(apath)) {
		err = nmDeactivateConnection(apath)
	}
	return
}

// DisconnectDevice will disconnect all connection in target device.
func (m *Manager) DisconnectDevice(devPath dbus.ObjectPath) (err error) {
	m.config.setDeviceLastConnectionUuid(devPath, "")
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
