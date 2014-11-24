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
	"sort"
	"strings"
)

// Helper function
func isNmObjectPathValid(p dbus.ObjectPath) bool {
	str := string(p)
	if len(str) == 0 || str == "/" {
		return false
	}
	return true
}

// General function wrappers for network manager
func nmGeneralGetAllDeviceHwAddr(devType uint32) (allHwAddr map[string]string) {
	allHwAddr = make(map[string]string)
	for _, devPath := range nmGetDevices() {
		if dev, err := nmNewDevice(devPath); err == nil && dev.DeviceType.Get() == devType {
			hwAddr, err := nmGeneralGetDeviceHwAddr(devPath)
			if err == nil {
				allHwAddr[dev.Interface.Get()] = hwAddr
			}
		}
	}
	return
}
func nmGeneralGetDeviceHwAddr(devPath dbus.ObjectPath) (hwAddr string, err error) {
	hwAddr = "00:00:00:00:00:00"
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}

	devType := dev.DeviceType.Get()
	switch devType {
	case NM_DEVICE_TYPE_ETHERNET:
		var devWired *nm.DeviceWired
		devWired, err = nmNewDeviceWired(devPath)
		if err == nil {
			hwAddr = devWired.PermHwAddress.Get()
		}
	case NM_DEVICE_TYPE_WIFI:
		var devWireless *nm.DeviceWireless
		devWireless, err = nmNewDeviceWireless(devPath)
		if err == nil {
			hwAddr = devWireless.PermHwAddress.Get()
		}
	case NM_DEVICE_TYPE_BT:
		var devBluetooth *nm.DeviceBluetooth
		devBluetooth, err = nmNewDeviceBluetooth(devPath)
		if err == nil {
			hwAddr = devBluetooth.HwAddress.Get()
		}
	case NM_DEVICE_TYPE_OLPC_MESH:
		var devOlpcMesh *nm.DeviceOlpcMesh
		devOlpcMesh, err = nmNewDeviceOlpcMesh(devPath)
		if err == nil {
			hwAddr = devOlpcMesh.HwAddress.Get()
		}
	case NM_DEVICE_TYPE_WIMAX:
		var devWiMax *nm.DeviceWiMax
		devWiMax, err = nmNewDeviceWiMax(devPath)
		if err == nil {
			hwAddr = devWiMax.HwAddress.Get()
		}
	case NM_DEVICE_TYPE_INFINIBAND:
		var devInfiniband *nm.DeviceInfiniband
		devInfiniband, err = nmNewDeviceInfiniband(devPath)
		if err == nil {
			hwAddr = devInfiniband.HwAddress.Get()
		}
	case NM_DEVICE_TYPE_BOND:
		var devBond *nm.DeviceBond
		devBond, err = nmNewDeviceBond(devPath)
		if err == nil {
			hwAddr = devBond.HwAddress.Get()
		}
	case NM_DEVICE_TYPE_BRIDGE:
		var devBridge *nm.DeviceBridge
		devBridge, err = nmNewDeviceBridge(devPath)
		if err == nil {
			hwAddr = devBridge.HwAddress.Get()
		}
	case NM_DEVICE_TYPE_VLAN:
		var devVlan *nm.DeviceVlan
		devVlan, err = nmNewDeviceVlan(devPath)
		if err == nil {
			hwAddr = devVlan.HwAddress.Get()
		}
	case NM_DEVICE_TYPE_MODEM, NM_DEVICE_TYPE_ADSL:
		// there is no hardware address for such devices
		err = fmt.Errorf("there is no hardware address for device modem and adsl")
		logger.Error(err)
	default:
		err = fmt.Errorf("unknown device type %d", devType)
		logger.Error(err)
	}
	hwAddr = strings.ToUpper(hwAddr)
	return
}

func nmGetDeviceIdentifiers() (devIds []string) {
	for _, devPath := range nmGetDevices() {
		id, _ := nmGeneralGetDeviceIdentifier(devPath)
		devIds = append(devIds, id)
	}
	return
}
func nmGeneralGetDeviceIdentifier(devPath dbus.ObjectPath) (devId string, err error) {
	// get device unique identifier, use hardware address if exists
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	devType := dev.DeviceType.Get()
	switch devType {
	case NM_DEVICE_TYPE_MODEM:
		modemPath := dev.Udi.Get()
		devId, err = mmGetModemDeviceIdentifier(dbus.ObjectPath(modemPath))
	case NM_DEVICE_TYPE_ADSL:
		err = fmt.Errorf("could not get adsl device identifier now")
		logger.Error(err)
	default:
		devId, err = nmGeneralGetDeviceHwAddr(devPath)
	}
	return
}

func nmGeneralGetDeviceRelatedUuid(devPath dbus.ObjectPath) (uuid string) {
	devId, err := nmGeneralGetDeviceIdentifier(devPath)
	if err != nil {
		return
	}
	return strToUuid(devId)
}

// get device network speed (Mb/s)
func nmGeneralGetDeviceSpeed(devPath dbus.ObjectPath) (speedStr string) {
	speed := uint32(0)
	speedStr = Tr("Unknown")
	nmDev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	switch t := nmDev.DeviceType.Get(); t {
	case NM_DEVICE_TYPE_ETHERNET:
		devWired, _ := nmNewDeviceWired(devPath)
		speed = devWired.Speed.Get()
	case NM_DEVICE_TYPE_WIFI:
		devWireless, _ := nmNewDeviceWireless(devPath)
		speed = devWireless.Bitrate.Get() / 1024
	default:
		err = fmt.Errorf("not support to get device speedStr for device type %d", t)
		logger.Error(err)
	}
	if speed != 0 {
		speedStr = fmt.Sprintf("%d Mb/s", speed)
	}
	return
}

func nmGeneralIsDeviceManaged(devPath dbus.ObjectPath) bool {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return false
	}
	if !isDeviceStateManaged(dev.State.Get()) {
		return false
	}
	devType := dev.DeviceType.Get()
	switch devType {
	case NM_DEVICE_TYPE_WIFI:
		if !nmGetWirelessHardwareEnabled() {
			return false
		}
	}
	return true
}

func nmGeneralGetDeviveSysPath(devPath dbus.ObjectPath) (sysPath string, err error) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	switch dev.DeviceType.Get() {
	case NM_DEVICE_TYPE_MODEM:
		sysPath, _ = mmGetModemDeviceSysPath(dbus.ObjectPath(dev.Udi.Get()))
	default:
		sysPath = dev.Udi.Get()
	}
	return
}

func nmGeneralGetDeviceVendor(devPath dbus.ObjectPath) (vendor string) {
	sysPath, err := nmGeneralGetDeviveSysPath(devPath)
	if err != nil {
		return
	}
	vendor = udevGetDeviceVendor(sysPath)
	return
}

func nmGeneralIsUsbDevice(devPath dbus.ObjectPath) bool {
	sysPath, err := nmGeneralGetDeviveSysPath(devPath)
	if err != nil {
		return false
	}
	return udevIsUsbDevice(sysPath)
}

func nmGeneralGetConnectionAutoconnect(cpath dbus.ObjectPath) (autoConnect bool) {
	switch nmGetConnectionType(cpath) {
	case NM_SETTING_VPN_SETTING_NAME:
		uuid, _ := nmGetConnectionUuid(cpath)
		autoConnect = manager.config.isVpnConnectionAutoConnect(uuid)
	default:
		autoConnect = nmGetConnectionAutoconnect(cpath)
	}
	return
}
func nmGeneralSetConnectionAutoconnect(cpath dbus.ObjectPath, autoConnect bool) {
	switch nmGetConnectionType(cpath) {
	case NM_SETTING_VPN_SETTING_NAME:
		uuid, _ := nmGetConnectionUuid(cpath)
		manager.config.setVpnConnectionAutoConnect(uuid, autoConnect)
	default:
		nmSetConnectionAutoconnect(cpath, autoConnect)
	}
}

// New network manager objects
func nmNewManager() (m *nm.Manager, err error) {
	m, err = nm.NewManager(dbusNmDest, dbusNmPath)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}
func nmNewDevice(devPath dbus.ObjectPath) (dev *nm.Device, err error) {
	dev, err = nm.NewDevice(dbusNmDest, devPath)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}
func nmNewDeviceWired(devPath dbus.ObjectPath) (dev *nm.DeviceWired, err error) {
	dev, err = nm.NewDeviceWired(dbusNmDest, devPath)
	if err != nil {
		logger.Error(err)
	}
	return
}
func nmNewDeviceWireless(devPath dbus.ObjectPath) (dev *nm.DeviceWireless, err error) {
	dev, err = nm.NewDeviceWireless(dbusNmDest, devPath)
	if err != nil {
		logger.Error(err)
	}
	return
}
func nmNewDeviceModem(devPath dbus.ObjectPath) (dev *nm.DeviceModem, err error) {
	dev, err = nm.NewDeviceModem(dbusNmDest, devPath)
	if err != nil {
		logger.Error(err)
	}
	return
}
func nmNewDeviceBluetooth(devPath dbus.ObjectPath) (dev *nm.DeviceBluetooth, err error) {
	dev, err = nm.NewDeviceBluetooth(dbusNmDest, devPath)
	if err != nil {
		logger.Error(err)
	}
	return
}
func nmNewDeviceOlpcMesh(devPath dbus.ObjectPath) (dev *nm.DeviceOlpcMesh, err error) {
	dev, err = nm.NewDeviceOlpcMesh(dbusNmDest, devPath)
	if err != nil {
		logger.Error(err)
	}
	return
}
func nmNewDeviceWiMax(devPath dbus.ObjectPath) (dev *nm.DeviceWiMax, err error) {
	dev, err = nm.NewDeviceWiMax(dbusNmDest, devPath)
	if err != nil {
		logger.Error(err)
	}
	return
}
func nmNewDeviceInfiniband(devPath dbus.ObjectPath) (dev *nm.DeviceInfiniband, err error) {
	dev, err = nm.NewDeviceInfiniband(dbusNmDest, devPath)
	if err != nil {
		logger.Error(err)
	}
	return
}
func nmNewDeviceBond(devPath dbus.ObjectPath) (dev *nm.DeviceBond, err error) {
	dev, err = nm.NewDeviceBond(dbusNmDest, devPath)
	if err != nil {
		logger.Error(err)
	}
	return
}
func nmNewDeviceBridge(devPath dbus.ObjectPath) (dev *nm.DeviceBridge, err error) {
	dev, err = nm.NewDeviceBridge(dbusNmDest, devPath)
	if err != nil {
		logger.Error(err)
	}
	return
}
func nmNewDeviceVlan(devPath dbus.ObjectPath) (dev *nm.DeviceVlan, err error) {
	dev, err = nm.NewDeviceVlan(dbusNmDest, devPath)
	if err != nil {
		logger.Error(err)
	}
	return
}
func nmNewDeviceAdsl(devPath dbus.ObjectPath) (dev *nm.DeviceAdsl, err error) {
	dev, err = nm.NewDeviceAdsl(dbusNmDest, devPath)
	if err != nil {
		logger.Error(err)
	}
	return
}
func nmNewAccessPoint(apPath dbus.ObjectPath) (ap *nm.AccessPoint, err error) {
	ap, err = nm.NewAccessPoint(dbusNmDest, apPath)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}
func nmNewActiveConnection(apath dbus.ObjectPath) (ac *nm.ActiveConnection, err error) {
	ac, err = nm.NewActiveConnection(dbusNmDest, apath)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}
func nmNewAgentManager() (manager *nm.AgentManager, err error) {
	manager, err = nm.NewAgentManager(dbusNmDest, "/org/freedesktop/NetworkManager/AgentManager")
	if err != nil {
		logger.Error(err)
		return
	}
	return
}
func nmNewDHCP4Config(path dbus.ObjectPath) (dhcp4 *nm.DHCP4Config, err error) {
	dhcp4, err = nm.NewDHCP4Config(dbusNmDest, path)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}
func nmNewDHCP6Config(path dbus.ObjectPath) (dhcp6 *nm.DHCP6Config, err error) {
	dhcp6, err = nm.NewDHCP6Config(dbusNmDest, path)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}
func nmNewSettingsConnection(cpath dbus.ObjectPath) (conn *nm.SettingsConnection, err error) {
	conn, err = nm.NewSettingsConnection(dbusNmDest, cpath)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}
func nmNewVpnConnection(apath dbus.ObjectPath) (vpnConn *nm.VPNConnection, err error) {
	vpnConn, err = nm.NewVPNConnection(dbusNmDest, apath)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

// Destroy network manager objects
func nmDestroyManager(m *nm.Manager) {
	if m == nil {
		logger.Error("Manager to destroy is nil")
		return
	}
	nm.DestroyManager(m)
}
func nmDestroyDevice(dev *nm.Device) {
	if dev == nil {
		logger.Error("Device to destroy is nil")
		return
	}
	nm.DestroyDevice(dev)
}
func nmDestroyDeviceWired(dev *nm.DeviceWired) {
	if dev == nil {
		logger.Error("DeviceWired to destroy is null")
		return
	}
	nm.DestroyDeviceWired(dev)
}
func nmDestroyDeviceWireless(dev *nm.DeviceWireless) {
	if dev == nil {
		logger.Error("DeviceWireless to destroy is nil")
		return
	}
	nm.DestroyDeviceWireless(dev)
}
func nmDestroyAccessPoint(ap *nm.AccessPoint) {
	if ap == nil {
		logger.Error("AccessPoint to destroy is nil")
		return
	}
	nm.DestroyAccessPoint(ap)
}
func nmDestroySettingsConnection(conn *nm.SettingsConnection) {
	if conn == nil {
		logger.Error("SettingsConnection to destroy is nil")
		return
	}
	nm.DestroySettingsConnection(conn)
}
func nmDestroyActiveConnection(aconn *nm.ActiveConnection) {
	if aconn == nil {
		logger.Error("ActiveConnection to destroy is nil")
		return
	}
	nm.DestroyActiveConnection(aconn)
}
func nmDestroyVpnConnection(vpnConn *nm.VPNConnection) {
	if vpnConn == nil {
		logger.Error("ActiveConnection to destroy is nil")
		return
	}
	nm.DestroyVPNConnection(vpnConn)
}

// Operate wrapper for network manager
func nmAgentRegister(identifier string) {
	am, err := nmNewAgentManager()
	if err != nil {
		return
	}
	err = am.Register(identifier)
	if err != nil {
		logger.Error(err)
	}
}

func nmAgentUnregister() {
	am, err := nmNewAgentManager()
	if err != nil {
		return
	}
	err = am.Unregister()
	if err != nil {
		logger.Error(err)
	}
}

func nmGetDevices() (devPaths []dbus.ObjectPath) {
	devPaths, err := nmManager.GetDevices()
	if err != nil {
		logger.Error(err)
	}
	return
}

func nmGetDevicesByType(devType uint32) (specDevPaths []dbus.ObjectPath) {
	for _, p := range nmGetDevices() {
		if dev, err := nmNewDevice(p); err == nil {
			if dev.DeviceType.Get() == devType {
				specDevPaths = append(specDevPaths, p)
			}
		}
	}
	return
}

func nmGetDeviceInterface(devPath dbus.ObjectPath) (devInterface string) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	devInterface = dev.Interface.Get()
	return
}

func nmAddAndActivateConnection(data connectionData, devPath dbus.ObjectPath) (cpath, apath dbus.ObjectPath, err error) {
	if len(devPath) == 0 {
		devPath = "/"
	}
	spath := dbus.ObjectPath("/")
	cpath, apath, err = nmManager.AddAndActivateConnection(data, devPath, spath)
	if err != nil {
		// if connection type is wireless hotspot, give a notification
		switch getCustomConnectionType(data) {
		case connectionWirelessAdhoc, connectionWirelessHotspot:
			notifyApModeNotSupport()
		}
		logger.Error(err, "devPath:", devPath)
		return
	}
	return
}

func nmActivateConnection(cpath, devPath dbus.ObjectPath) (apath dbus.ObjectPath, err error) {
	spath := dbus.ObjectPath("/")
	apath, err = nmManager.ActivateConnection(cpath, devPath, spath)
	if err != nil {
		// if connection type is wireless hotspot, give a notification
		if data, err := nmGetConnectionData(cpath); err == nil {
			switch getCustomConnectionType(data) {
			case connectionWirelessAdhoc, connectionWirelessHotspot:
				notifyApModeNotSupport()
			}
		}
		logger.Error(err)
		return
	}
	return
}

func nmDeactivateConnection(apath dbus.ObjectPath) (err error) {
	err = nmManager.DeactivateConnection(apath)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

func nmGetActiveConnections() (apaths []dbus.ObjectPath) {
	apaths = nmManager.ActiveConnections.Get()
	return
}

func nmGetVpnActiveConnections() (apaths []dbus.ObjectPath) {
	for _, p := range nmGetActiveConnections() {
		if aconn, err := nmNewActiveConnection(p); err == nil {
			if aconn.Vpn.Get() {
				apaths = append(apaths, p)
			}
		}
	}
	return
}

func nmGetVpnConnectionState(apath dbus.ObjectPath) (state uint32) {
	vpnConn, err := nmNewVpnConnection(apath)
	if err != nil {
		return
	}
	state = vpnConn.VpnState.Get()
	return
}

func nmGetAccessPoints(devPath dbus.ObjectPath) (apPaths []dbus.ObjectPath) {
	dev, err := nmNewDeviceWireless(devPath)
	if err != nil {
		return
	}
	apPaths, err = dev.GetAccessPoints()
	if err != nil {
		logger.Error(err)
	}
	return
}

func nmGetManagerState() (state uint32) {
	state = nmManager.State.Get()
	return
}

func nmGetActiveConnectionByUuid(uuid string) (apaths []dbus.ObjectPath, err error) {
	for _, apath := range nmGetActiveConnections() {
		if ac, tmperr := nmNewActiveConnection(apath); tmperr == nil {
			if ac.Uuid.Get() == uuid {
				apaths = append(apaths, apath)
				return
			}
		}
	}
	err = fmt.Errorf("not found active connection with uuid %s", uuid)
	return
}

func nmGetActiveConnectionState(apath dbus.ObjectPath) (state uint32) {
	aconn, err := nmNewActiveConnection(apath)
	if err != nil {
		return
	}
	state = aconn.State.Get()
	return
}

func nmGetActiveConnectionVpn(apath dbus.ObjectPath) (isVpn bool) {
	aconn, err := nmNewActiveConnection(apath)
	if err != nil {
		return
	}
	isVpn = aconn.Vpn.Get()
	return
}

func nmGetConnectionData(cpath dbus.ObjectPath) (data connectionData, err error) {
	nmConn, err := nm.NewSettingsConnection(dbusNmDest, cpath)
	if err != nil {
		logger.Error(err)
		return
	}
	data, err = nmConn.GetSettings()
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

func nmUpdateConnectionData(cpath dbus.ObjectPath, data connectionData) (err error) {
	nmConn, err := nmNewSettingsConnection(cpath)
	if err != nil {
		return
	}
	err = nmConn.Update(data)
	if err != nil {
		logger.Error(err)
	}
	return
}

func nmGetConnectionSecrets(cpath dbus.ObjectPath, secretField string) (secrets connectionData, err error) {
	nmConn, err := nm.NewSettingsConnection(dbusNmDest, cpath)
	if err != nil {
		logger.Error(err)
		return
	}
	secrets, err = nmConn.GetSecrets(secretField)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

func nmSetConnectionAutoconnect(cpath dbus.ObjectPath, autoConnect bool) (err error) {
	data, err := nmGetConnectionData(cpath)
	if err != nil {
		return
	}
	setSettingConnectionAutoconnect(data, autoConnect)
	return nmUpdateConnectionData(cpath, data)
}
func nmGetConnectionAutoconnect(cpath dbus.ObjectPath) (autoConnect bool) {
	data, err := nmGetConnectionData(cpath)
	if err != nil {
		return
	}
	autoConnect = getSettingConnectionAutoconnect(data)
	return
}

func nmGetConnectionId(cpath dbus.ObjectPath) (id string) {
	data, err := nmGetConnectionData(cpath)
	if err != nil {
		return
	}
	id = getSettingConnectionId(data)
	if len(id) == 0 {
		logger.Error("get Id of connection failed, id is empty")
	}
	return
}
func nmSetConnectionId(cpath dbus.ObjectPath, id string) (err error) {
	data, err := nmGetConnectionData(cpath)
	if err != nil {
		return
	}
	setSettingConnectionId(data, id)
	return nmUpdateConnectionData(cpath, data)
}

func nmGetConnectionUuid(cpath dbus.ObjectPath) (uuid string, err error) {
	data, err := nmGetConnectionData(cpath)
	if err != nil {
		return
	}
	uuid = getSettingConnectionUuid(data)
	return
}

func nmGetConnectionType(cpath dbus.ObjectPath) (ctype string) {
	data, err := nmGetConnectionData(cpath)
	if err != nil {
		return
	}
	ctype = getSettingConnectionType(data)
	if len(ctype) == 0 {
		logger.Error("get type of connection failed, type is empty")
	}
	return
}

func nmGetConnectionList() (connections []dbus.ObjectPath) {
	connections, err := nmSettings.ListConnections()
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

func nmGetConnectionUuids() (uuids []string) {
	for _, cpath := range nmGetConnectionList() {
		if uuid, err := nmGetConnectionUuid(cpath); err == nil {
			uuids = append(uuids, uuid)
		}
	}
	return
}

func nmGetConnectionUuidsByType(connType string) (uuids []string) {
	for _, cpath := range nmGetConnectionList() {
		if nmGetConnectionType(cpath) == connType {
			if uuid, err := nmGetConnectionUuid(cpath); err == nil {
				uuids = append(uuids, uuid)
			}
		}
	}
	return
}

func nmGetConnectionIds() (ids []string) {
	for _, cpath := range nmGetConnectionList() {
		ids = append(ids, nmGetConnectionId(cpath))
	}
	return
}

func nmGetConnectionById(id string) (cpath dbus.ObjectPath, err error) {
	for _, cpath = range nmGetConnectionList() {
		data, tmperr := nmGetConnectionData(cpath)
		if tmperr != nil {
			continue
		}
		if getSettingConnectionId(data) == id {
			return
		}
	}
	err = fmt.Errorf("connection with id %s not found", id)
	return
}

func nmGetConnectionByUuid(uuid string) (cpath dbus.ObjectPath, err error) {
	cpath, err = nmSettings.GetConnectionByUuid(uuid)
	return
}

// get wireless connection by ssid, the connection with special hardware address is priority
// TODO: use available connections instead
func nmGetWirelessConnection(ssid []byte, devPath dbus.ObjectPath) (cpath dbus.ObjectPath, ok bool) {
	var hwAddr string
	if len(devPath) != 0 {
		hwAddr, _ = nmGeneralGetDeviceHwAddr(devPath)
	}
	ok = false
	for _, p := range nmGetWirelessConnectionListBySsid(ssid) {
		data, err := nmGetConnectionData(p)
		if err != nil {
			continue
		}
		if isSettingWirelessMacAddressExists(data) {
			if hwAddr == convertMacAddressToString(getSettingWirelessMacAddress(data)) {
				cpath = p
				ok = true
				return
			}
		} else if !ok {
			cpath = p
			ok = true
		}
	}
	return
}

func nmGetWirelessConnectionListBySsid(ssid []byte) (cpaths []dbus.ObjectPath) {
	for _, p := range nmGetConnectionList() {
		data, err := nmGetConnectionData(p)
		if err != nil {
			continue
		}
		if getCustomConnectionType(data) != connectionWireless {
			continue
		}
		if isSettingWirelessSsidExists(data) && string(getSettingWirelessSsid(data)) == string(ssid) {
			cpaths = append(cpaths, p)
		}
	}
	return
}

func nmGetConnectionSsidByUuid(uuid string) (ssid []byte) {
	cpath, err := nmGetConnectionByUuid(uuid)
	if err != nil {
		return
	}
	data, err := nmGetConnectionData(cpath)
	if err != nil {
		return
	}
	ssid = getSettingWirelessSsid(data)
	return
}

func nmAddConnection(data connectionData) (cpath dbus.ObjectPath, err error) {
	cpath, err = nmSettings.AddConnection(data)
	if err != nil {
		logger.Error(err)
	}
	return
}

func nmGetDhcp4Info(path dbus.ObjectPath) (ip, mask, route, dns string) {
	ip = "0.0.0.0"
	mask = "0.0.0.0"
	route = "0.0.0.0"
	dns = "0.0.0.0"
	dhcp4, err := nmNewDHCP4Config(path)
	if err != nil {
		return
	}
	options := dhcp4.Options.Get()
	if ipData, ok := options["ip_address"]; ok {
		ip, _ = ipData.Value().(string)
	}
	if maskData, ok := options["subnet_mask"]; ok {
		mask, _ = maskData.Value().(string)
	}
	if routeData, ok := options["routers"]; ok {
		route, _ = routeData.Value().(string)
	}
	if dnsData, ok := options["domain_name_servers"]; ok {
		dns, _ = dnsData.Value().(string)
	}
	return
}

func nmGetDhcp6Info(path dbus.ObjectPath) (ip, route, dns string) {
	ip = "0::0"
	route = ""
	dns = "0::0"
	dhcp6, err := nmNewDHCP6Config(path)
	if err != nil {
		return
	}
	options := dhcp6.Options.Get()
	if ipData, ok := options["ip6_address"]; ok {
		ip, _ = ipData.Value().(string)
	}
	// FIXME how
	if routeData, ok := options["routers"]; ok {
		route, _ = routeData.Value().(string)
	}
	if dnsData, ok := options["dhcp6_name_servers"]; ok {
		dns, _ = dnsData.Value().(string)
	}
	return
}

func nmGetDeviceState(devPath dbus.ObjectPath) (state uint32) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return NM_DEVICE_STATE_UNKNOWN
	}
	state = dev.State.Get()
	return
}

func nmGetDeviceType(devPath dbus.ObjectPath) (devType uint32) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return NM_DEVICE_TYPE_UNKNOWN
	}
	devType = dev.DeviceType.Get()
	return
}

func nmGetDeviceUdi(devPath dbus.ObjectPath) (udi string) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	udi = dev.Udi.Get()
	return
}

func nmGetDeviceActiveConnection(devPath dbus.ObjectPath) (acPath dbus.ObjectPath) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	acPath = dev.ActiveConnection.Get()
	return
}

func nmGetDeviceAvailableConnections(devPath dbus.ObjectPath) (paths []dbus.ObjectPath) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	paths = dev.AvailableConnections.Get()
	return
}

func nmGetDeviceActiveConnectionUuid(devPath dbus.ObjectPath) (uuid string, err error) {
	acPath := nmGetDeviceActiveConnection(devPath)
	aconn, err := nmNewActiveConnection(acPath)
	if err != nil {
		return
	}
	uuid = aconn.Uuid.Get()
	return
}

func nmGetDeviceActiveConnectionData(devPath dbus.ObjectPath) (data connectionData, err error) {
	if !isDeviceStateInActivating(nmGetDeviceState(devPath)) {
		err = fmt.Errorf("device is inactivated %s", devPath)
		return
	}
	acPath := nmGetDeviceActiveConnection(devPath)
	aconn, err := nmNewActiveConnection(acPath)
	if err != nil {
		return
	}
	conn, err := nmNewSettingsConnection(aconn.Connection.Get())
	if err != nil {
		return
	}
	data, err = conn.GetSettings()
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

func nmManagerEnable(enable bool) (err error) {
	err = nmManager.Enable(enable)
	if err != nil {
		logger.Error(err)
	}
	return
}

func nmGetPrimaryConnection() (cpath dbus.ObjectPath) {
	// TODO need update dbus-factory
	// cpath = nmManager.PrimaryConnection.Get()
	cpath = ""
	return
}

func nmGetNetworkState() uint32 {
	return nmManager.State.Get()
}
func nmIsNetworkOffline() bool {
	state := nmManager.State.Get()
	if state == NM_STATE_DISCONNECTED || state == NM_STATE_ASLEEP {
		return true
	}
	return false
}

func nmGetNetworkEnabled() bool {
	return nmManager.NetworkingEnabled.Get()
}
func nmGetWirelessHardwareEnabled() bool {
	return nmManager.WirelessHardwareEnabled.Get()
}
func nmGetWirelessEnabled() bool {
	return nmManager.WirelessEnabled.Get()
}

func nmSetNetworkingEnabled(enabled bool) {
	if nmManager.NetworkingEnabled.Get() != enabled {
		nmManagerEnable(enabled)
	} else {
		logger.Warning("NetworkingEnabled already set as", enabled)
	}
	return
}
func nmSetWirelessEnabled(enabled bool) {
	if nmManager.WirelessEnabled.Get() != enabled {
		nmManager.WirelessEnabled.Set(enabled)
	} else {
		logger.Warning("WirelessEnabled already set as", enabled)
	}
	return
}
func nmSetWwanEnabled(enabled bool) {
	if nmManager.WwanEnabled.Get() != enabled {
		nmManager.WwanEnabled.Set(enabled)
	} else {
		logger.Warning("WwanEnabled already set as", enabled)
	}
}

type autoConnectConn struct {
	id        string
	uuid      string
	timestamp uint64
}
type autoConnectConns []autoConnectConn

func (acs autoConnectConns) Len() int {
	return len(acs)
}
func (acs autoConnectConns) Swap(i, j int) {
	acs[i], acs[j] = acs[j], acs[i]
}
func (acs autoConnectConns) Less(i, j int) bool {
	return acs[i].timestamp < acs[j].timestamp
}
func nmGetConnectionUuidsForAutoConnect(devPath dbus.ObjectPath, lastConnectionUuid string) (uuids []string) {
	acs := make(autoConnectConns, 0)
	devRelatedUuid := nmGeneralGetDeviceRelatedUuid(devPath)
	for _, cpath := range nmGetDeviceAvailableConnections(devPath) {
		if cdata, err := nmGetConnectionData(cpath); err == nil {
			uuid := getSettingConnectionUuid(cdata)
			switch getCustomConnectionType(cdata) {
			case connectionWired, connectionMobileGsm, connectionMobileCdma:
				if devRelatedUuid != uuid {
					// ignore connections that not matching the
					// device, etc other wired connections that create
					// in other ways
					continue
				}
			}
			if uuid == lastConnectionUuid {
				continue
			}
			if getSettingConnectionAutoconnect(cdata) {
				ac := autoConnectConn{
					id:        getSettingConnectionId(cdata),
					uuid:      uuid,
					timestamp: getSettingConnectionTimestamp(cdata),
				}
				acs = append(acs, ac)
			}
		}
	}
	sort.Sort(sort.Reverse(acs))
	logger.Debugf("device type: %s, auto connect connections: %v",
		getCustomDeviceType(nmGetDeviceType(devPath)), acs)
	if len(lastConnectionUuid) > 0 {
		// the last activated connection has the highest priority if
		// exists
		uuids = []string{lastConnectionUuid}
	}
	for _, ac := range acs {
		uuids = append(uuids, ac.uuid)
	}
	return
}

func nmRunOnceUntilDeviceAvailable(devPath dbus.ObjectPath, cb func()) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	state := dev.State.Get()
	if isDeviceStateAvailable(state) {
		cb()
	} else {
		hasRun := false
		dev.ConnectStateChanged(func(newState uint32, oldState uint32, reason uint32) {
			if !hasRun && isDeviceStateAvailable(newState) {
				cb()
				nmDestroyDevice(dev)
				hasRun = true
			}
		})
	}
}

func nmRunOnceUtilNetworkAvailable(cb func()) {
	nm, err := nmNewManager()
	if err != nil {
		return
	}
	state := nm.State.Get()
	if state >= NM_STATE_CONNECTED_LOCAL {
		cb()
	} else {
		hasRun := false
		nm.ConnectStateChanged(func(state uint32) {
			if !hasRun && state >= NM_STATE_CONNECTED_LOCAL {
				cb()
				nmDestroyManager(nm)
				hasRun = true
			}
		})
	}
}
