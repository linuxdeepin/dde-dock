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
	"fmt"
	"pkg.linuxdeepin.com/lib/dbus"
	. "pkg.linuxdeepin.com/lib/gettext"
)

type activeConnection struct {
	path dbus.ObjectPath

	Devices []dbus.ObjectPath
	Id      string
	Uuid    string
	State   uint32
	Vpn     bool
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
			} else if isVpnConnectionStateFailed(state) {
				notifyVpnFailed(aconn.Id, reason)
			}

			if isVpnConnectionStateInActivating(state) {
				m.switchHandler.doEnableVpn(true)
			} else {
				delete(m.activeConnections, s.Path)
				m.switchHandler.turnOffVpnSwitchIfNeed(autoConnectTimeout)
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

func (m *Manager) clearActiveConnections() {
	m.activeConnectionsLocker.Lock()
	defer m.activeConnectionsLocker.Unlock()

	m.activeConnections = make(map[dbus.ObjectPath]*activeConnection)
	m.setPropActiveConnections()
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
