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
	"strings"
)

const dbusNmDest = "org.freedesktop.NetworkManager"

var (
	nmManager, _  = nm.NewManager(dbusNmDest, "/org/freedesktop/NetworkManager")
	nmSettings, _ = nm.NewSettings(dbusNmDest, "/org/freedesktop/NetworkManager/Settings")
)

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
			hwAddr = devWired.HwAddress.Get()
		}
	case NM_DEVICE_TYPE_WIFI:
		var devWireless *nm.DeviceWireless
		devWireless, err = nmNewDeviceWireless(devPath)
		if err == nil {
			hwAddr = devWireless.HwAddress.Get()
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
	default:
		err = fmt.Errorf("unknown device type %d", devType)
		logger.Error(err)
	}
	hwAddr = strings.ToUpper(hwAddr)
	return
}
func nmGetDeviceIdentifier(devPath dbus.ObjectPath) (devId string, err error) {
	// get device unique identifier, use hardware address if exists
	hwAddr, err := nmGeneralGetDeviceHwAddr(devPath)
	if err == nil {
		devId = hwAddr
		return
	}

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
		// TODO
		err = fmt.Errorf("could not get adsl device identifier now")
		logger.Error(err)
	default:
		err = fmt.Errorf("unknown device type %d", devType)
		logger.Error(err)
	}
	return
}
func nmGetDeviceIdentifiers() (devIds []string) {
	for _, devPath := range nmGetDevices() {
		id, _ := nmGetDeviceIdentifier(devPath)
		devIds = append(devIds, id)
	}
	return
}

// New network manager objects
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
func nmNewSettingsConnection(cpath dbus.ObjectPath) (conn *nm.SettingsConnection, err error) {
	conn, err = nm.NewSettingsConnection(dbusNmDest, cpath)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

// Destroy network manager objects
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

// Operate wrapper for network manager
func nmAgentRegister(identifier string) {
	manager, err := nmNewAgentManager()
	if err != nil {
		return
	}
	err = manager.Register(identifier)
	if err != nil {
		logger.Error(err)
	}
}

func nmAgentUnregister() {
	manager, err := nmNewAgentManager()
	if err != nil {
		return
	}
	err = manager.Unregister()
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

func nmGetSpecialDevices(devType uint32) (specDevPaths []dbus.ObjectPath) {
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

func nmGetActiveConnectionByUuid(uuid string) (apath dbus.ObjectPath, ok bool) {
	for _, apath = range nmGetActiveConnections() {
		if ac, err := nmNewActiveConnection(apath); err == nil {
			if ac.Uuid.Get() == uuid {
				ok = true
				return
			}
		}
	}
	ok = false
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

func nmGetSpecialConnectionUuids(connType string) (uuids []string) {
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

func nmGetConnectionById(id string) (cpath dbus.ObjectPath, ok bool) {
	for _, cpath = range nmGetConnectionList() {
		data, err := nmGetConnectionData(cpath)
		if err != nil {
			continue
		}
		if getSettingConnectionId(data) == id {
			ok = true
			return
		}
	}
	ok = false
	return
}

func nmGetConnectionByUuid(uuid string) (cpath dbus.ObjectPath, err error) {
	cpath, err = nmSettings.GetConnectionByUuid(uuid)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

// get wireless connection by ssid, the connection with special hardware address is priority
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

func nmAddConnection(data connectionData) (cpath dbus.ObjectPath, err error) {
	cpath, err = nmSettings.AddConnection(data)
	if err != nil {
		logger.Error(err)
	}
	return
}

func nmGetDHCP4Info(path dbus.ObjectPath) (ip, mask, route, dns string) {
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

func nmGetDeviceState(devPath dbus.ObjectPath) (state uint32, err error) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	state = dev.State.Get()
	return
}

func nmGetDeviceType(devPath dbus.ObjectPath) (devType uint32, err error) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	devType = dev.DeviceType.Get()
	return
}

func nmGetDeviceActiveConnection(devPath dbus.ObjectPath) (acPath dbus.ObjectPath, err error) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	acPath = dev.ActiveConnection.Get()
	if len(acPath) == 0 || acPath == "/" {
		// don't need logger error here
		err = fmt.Errorf("there is no active connection for device", devPath)
		return
	}
	return
}

func nmGetDeviceActiveConnectionUuid(devPath dbus.ObjectPath) (uuid string, err error) {
	acPath, err := nmGetDeviceActiveConnection(devPath)
	if err != nil {
		return
	}
	aconn, err := nmNewActiveConnection(acPath)
	if err != nil {
		return
	}
	uuid = aconn.Uuid.Get()
	return
}

func nmGetDeviceActiveConnectionData(devPath dbus.ObjectPath) (data connectionData, err error) {
	acPath, err := nmGetDeviceActiveConnection(devPath)
	if err != nil {
		return
	}
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

func nmSetNetworkingEnabled(enabled bool) {
	if nmManager.NetworkingEnabled.Get() == enabled {
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
