/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package network

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	dbus "github.com/godbus/dbus"
	nmdbus "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.networkmanager"
	"pkg.deepin.io/dde/daemon/network/nm"
	"pkg.deepin.io/lib/dbusutil/proxy"
	. "pkg.deepin.io/lib/gettext"
)

// Wrapper NetworkManger dbus methods to hide
// "go-dbus-factory/org.freedesktop.networkmanager" details for other source
// files.

// Custom device state reasons
const (
	CUSTOM_NM_DEVICE_STATE_REASON_CABLE_UNPLUGGED = iota + 1000
	CUSTOM_NM_DEVICE_STATE_REASON_WIRELESS_DISABLED
	CUSTOM_NM_DEVICE_STATE_REASON_MODEM_NO_SIGNAL
	CUSTOM_NM_DEVICE_STATE_REASON_MODEM_WRONG_PLAN
)

const (
	devWhitelistHuaweiFile = "/lib/vendor/interface"
)

var nmPermissions map[string]string

// Helper functions
func isNmObjectPathValid(p dbus.ObjectPath) bool {
	str := string(p)
	if len(str) == 0 || str == "/" {
		return false
	}
	return true
}

func isDeviceTypeValid(devType uint32) bool {
	switch devType {
	case nm.NM_DEVICE_TYPE_GENERIC, nm.NM_DEVICE_TYPE_UNKNOWN, nm.NM_DEVICE_TYPE_BT, nm.NM_DEVICE_TYPE_TEAM, nm.NM_DEVICE_TYPE_TUN, nm.NM_DEVICE_TYPE_IP_TUNNEL, nm.NM_DEVICE_TYPE_MACVLAN, nm.NM_DEVICE_TYPE_VXLAN, nm.NM_DEVICE_TYPE_VETH, nm.NM_DEVICE_TYPE_PPP:
		return false
	}
	return true
}

// check current device state
func isDeviceStateManaged(state uint32) bool {
	return state > nm.NM_DEVICE_STATE_UNMANAGED
}
func isDeviceStateAvailable(state uint32) bool {
	return state > nm.NM_DEVICE_STATE_UNAVAILABLE
}
func isDeviceStateActivated(state uint32) bool {
	return state == nm.NM_DEVICE_STATE_ACTIVATED
}
func isDeviceStateInActivating(state uint32) bool {
	return state == nm.NM_DEVICE_STATE_ACTIVATED
}

func isDeviceStateReasonInvalid(reason uint32) bool {
	switch reason {
	case nm.NM_DEVICE_STATE_REASON_UNKNOWN:
		return true
	}
	return false
}

// check if connection activating or activated
func isConnectionStateInActivating(state uint32) bool {
	if state == nm.NM_ACTIVE_CONNECTION_STATE_ACTIVATING ||
		state == nm.NM_ACTIVE_CONNECTION_STATE_ACTIVATED {
		return true
	}
	return false
}

func isInDeviceWhitelist(filename string, ifc string) bool {
	if len(ifc) == 0 {
		return false
	}
	fr, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer fr.Close()

	var scanner = bufio.NewScanner(fr)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		if line == ifc {
			return true
		}
	}
	return false

}

func isVirtualDeviceIfc(dev *nmdbus.Device) bool {
	driver, _ := dev.Driver().Get(0)

	//// workaround for huawei pangu
	ifc, _ := dev.Interface().Get(0)
	if isInDeviceWhitelist(devWhitelistHuaweiFile, ifc) {
		return false
	}

	switch driver {
	case "dummy", "veth", "vboxnet", "vmnet":
		return true
	case "unknown", "vmxnet", "vmxnet2", "vmxnet3":
		// sometimes we could not get vmnet dirver name, so check the
		// udi sys path if is prefix with /sys/devices/virtual/net
		devUdi, _ := dev.Udi().Get(0)
		if strings.HasPrefix(devUdi, "/sys/devices/virtual/net") ||
			strings.HasPrefix(devUdi, "/virtual/device") ||
			strings.HasPrefix(ifc, "vmnet") {
			return true
		}
	}
	return false
}

func nmGeneralGetDeviceHwAddr(devPath dbus.ObjectPath, perm bool) (hwAddr string, err error) {
	hwAddr = "00:00:00:00:00:00"
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}

	devType, _ := dev.DeviceType().Get(0)
	switch devType {
	case nm.NM_DEVICE_TYPE_ETHERNET:
		devWired := dev.Wired()
		hwAddr = ""
		if perm {
			hwAddr, err = devWired.PermHwAddress().Get(0)
			hwAddr = strings.ToUpper(hwAddr)
			return
		}
		if hwAddr == "" {
			// may get PermHwAddress failed under NetworkManager 1.4.1
			hwAddr, _ = devWired.HwAddress().Get(0)
		}
	case nm.NM_DEVICE_TYPE_WIFI:
		devWireless := dev.Wireless()
		hwAddr = ""
		if perm {
			hwAddr, _ = devWireless.PermHwAddress().Get(0)
		}
		if len(hwAddr) == 0 {
			hwAddr, _ = devWireless.HwAddress().Get(0)
		}
	case nm.NM_DEVICE_TYPE_BT:
		devBluetooth := dev.Bluetooth()
		hwAddr, _ = devBluetooth.HwAddress().Get(0)
	case nm.NM_DEVICE_TYPE_OLPC_MESH:
		devOlpcMesh := dev.OlpcMesh()
		hwAddr, _ = devOlpcMesh.HwAddress().Get(0)
	case nm.NM_DEVICE_TYPE_WIMAX:
		devWiMax := dev.WiMax()
		hwAddr, _ = devWiMax.HwAddress().Get(0)
	case nm.NM_DEVICE_TYPE_INFINIBAND:
		devInfiniband := dev.Infiniband()
		hwAddr, _ = devInfiniband.HwAddress().Get(0)
	case nm.NM_DEVICE_TYPE_BOND:
		devBond := dev.Bond()
		hwAddr, _ = devBond.HwAddress().Get(0)
	case nm.NM_DEVICE_TYPE_BRIDGE:
		devBridge := dev.Bridge()
		hwAddr, _ = devBridge.HwAddress().Get(0)
	case nm.NM_DEVICE_TYPE_VLAN:
		devVlan := dev.Vlan()
		hwAddr, _ = devVlan.HwAddress().Get(0)
	case nm.NM_DEVICE_TYPE_GENERIC:
		devGeneric := dev.Generic()
		hwAddr, _ = devGeneric.HwAddress().Get(0)
	case nm.NM_DEVICE_TYPE_TEAM:
		devTeam := dev.Team()
		hwAddr, _ = devTeam.HwAddress().Get(0)
	case nm.NM_DEVICE_TYPE_MODEM, nm.NM_DEVICE_TYPE_ADSL, nm.NM_DEVICE_TYPE_TUN, nm.NM_DEVICE_TYPE_IP_TUNNEL, nm.NM_DEVICE_TYPE_MACVLAN, nm.NM_DEVICE_TYPE_VXLAN, nm.NM_DEVICE_TYPE_VETH:
		// there is no hardware address for such devices
		err = fmt.Errorf("there is no hardware address for device modem, adsl, tun")
	default:
		err = fmt.Errorf("unknown device type %d", devType)
		logger.Error(err)
	}
	hwAddr = strings.ToUpper(hwAddr)
	return
}

func nmGeneralGetDeviceIdentifier(devPath dbus.ObjectPath) (devId string, err error) {
	// get device unique identifier, use hardware address if exists
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}

	devType, _ := dev.DeviceType().Get(0)
	switch devType {
	case nm.NM_DEVICE_TYPE_MODEM:
		modemPath, _ := dev.Udi().Get(0)
		devId, err = mmGetModemDeviceIdentifier(dbus.ObjectPath(modemPath))
	case nm.NM_DEVICE_TYPE_ADSL:
		err = fmt.Errorf("could not get adsl device identifier now")
		logger.Error(err)
	case nm.NM_DEVICE_TYPE_ETHERNET:
		// some device the 'hw_addr_perm' unset by driver, so use 'hw_addr' as id
		// PMS Bug ID: 16704
		devId, err = nmGeneralGetDeviceHwAddr(devPath, false)
	default:
		devId, err = nmGeneralGetDeviceHwAddr(devPath, true)
	}
	return
}

// return special unique connection uuid for device, etc wired device
// connection
func nmGeneralGetDeviceUniqueUuid(devPath dbus.ObjectPath) (uuid string) {
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
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}

	t, _ := dev.DeviceType().Get(0)
	switch t {
	case nm.NM_DEVICE_TYPE_ETHERNET:
		devWired := dev.Wired()
		speed, _ = devWired.Speed().Get(0)
	case nm.NM_DEVICE_TYPE_WIFI:
		devWireless := dev.Wireless()
		bitRate, _ := devWireless.Bitrate().Get(0)
		/**
		 * NMSettingWireless:rate:
		 *
		 * If non-zero, directs the device to only use the specified bitrate for
		 * communication with the access point.  Units are in Kb/s, ie 5500 = 5.5
		 * Mbit/s.  This property is highly driver dependent and not all devices
		 * support setting a static bitrate.
		**/
		speed = uint32(math.Trunc((float64(bitRate)/1000.0 + 0.5) * 10 / 10))
	case nm.NM_DEVICE_TYPE_MODEM:
		// TODO: getting device speed for modem device
	default: // ignore speed for other device types
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

	state, _ := dev.State().Get(0)
	if !isDeviceStateManaged(state) {
		return false
	}
	devType, _ := dev.DeviceType().Get(0)
	switch devType {
	case nm.NM_DEVICE_TYPE_WIFI:
		if !nmGetWirelessHardwareEnabled() {
			return false
		}
	}
	return true
}

func nmGeneralGetDeviceSysPath(devPath dbus.ObjectPath) (sysPath string, err error) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}

	deviceType, _ := dev.DeviceType().Get(0)
	devUdi, _ := dev.Udi().Get(0)
	switch deviceType {
	case nm.NM_DEVICE_TYPE_MODEM:
		sysPath, _ = mmGetModemDeviceSysPath(dbus.ObjectPath(devUdi))
	default:
		sysPath = devUdi
	}
	return
}

func nmGeneralGetDeviceDesc(devPath dbus.ObjectPath) (desc string) {
	sysPath, err := nmGeneralGetDeviceSysPath(devPath)
	if err != nil {
		return
	}
	desc, ok := udevGetDeviceDesc(sysPath)
	if !ok {
		desc = nmGetDeviceInterface(devPath)
	}
	return
}

func nmGeneralIsUsbDevice(devPath dbus.ObjectPath) bool {
	var sysPath string
	var err error
	for i := 0; i < 10; i++ {
		sysPath, err = nmGeneralGetDeviceSysPath(devPath)
		if err != nil {
			logger.Warningf("failed to get device %v sys path: %v",
				devPath, err)
			return false
		}
		//logger.Debug("sysPath:", sysPath)
		if strings.HasPrefix(sysPath, "/virtual/device/placeholder/") ||
			sysPath == "" {
			logger.Debug("sleep 500ms")
			time.Sleep(500 * time.Millisecond)
			continue
		} else {
			break
		}
	}

	return udevIsUsbDevice(sysPath)
}

// New network manager objects
func nmNewManager() (m *nmdbus.Manager, err error) {
	systemBus, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	m = nmdbus.NewManager(systemBus)
	return
}
func nmNewDevice(devPath dbus.ObjectPath) (dev *nmdbus.Device, err error) {
	systemBus, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	dev, err = nmdbus.NewDevice(systemBus, devPath)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

func nmNewAccessPoint(apPath dbus.ObjectPath) (ap *nmdbus.AccessPoint, err error) {
	systemBus, err := dbus.SystemBus()
	if err != nil {
		return
	}
	ap, err = nmdbus.NewAccessPoint(systemBus, apPath)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}
func nmNewActiveConnection(apath dbus.ObjectPath) (aconn *nmdbus.ActiveConnection, err error) {
	systemBus, err := dbus.SystemBus()
	if err != nil {
		return
	}
	aconn, err = nmdbus.NewActiveConnection(systemBus, apath)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

func nmNewIP4Config(path dbus.ObjectPath) (ip4config *nmdbus.IP4Config, err error) {
	systemBus, err := dbus.SystemBus()
	if err != nil {
		return
	}
	ip4config, err = nmdbus.NewIP4Config(systemBus, path)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}
func nmNewIP6Config(path dbus.ObjectPath) (ip6config *nmdbus.IP6Config, err error) {
	systemBus, err := dbus.SystemBus()
	if err != nil {
		return
	}

	ip6config, err = nmdbus.NewIP6Config(systemBus, path)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}
func nmNewSettingsConnection(cpath dbus.ObjectPath) (conn *nmdbus.ConnectionSettings, err error) {
	systemBus, err := dbus.SystemBus()
	if err != nil {
		return
	}
	conn, err = nmdbus.NewConnectionSettings(systemBus, cpath)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

func nmDestroyDevice(dev *nmdbus.Device) {
	if dev == nil {
		logger.Error("Device to destroy is nil")
		return
	}
	dev.RemoveHandler(proxy.RemoveAllHandlers)
}

func nmDestroyAccessPoint(ap *nmdbus.AccessPoint) {
	if ap == nil {
		logger.Error("AccessPoint to destroy is nil")
		return
	}
	ap.RemoveHandler(proxy.RemoveAllHandlers)
}

func nmDestroySettingsConnection(conn *nmdbus.ConnectionSettings) {
	if conn == nil {
		logger.Error("SettingsConnection to destroy is nil")
		return
	}
	conn.RemoveHandler(proxy.RemoveAllHandlers)
}

// Operate wrapper for network manager
func nmHasSystemSettingsModifyPermission() (hasPerm bool) {
	permissions := nmGetPermissionsInstance()
	hasPermStr, ok := permissions["org.freedesktop.NetworkManager.settings.modify.system"]
	if !ok {
		hasPermStr = "no"
	}
	if hasPermStr == "yes" {
		hasPerm = true
	} else {
		hasPerm = false
	}
	return
}
func nmGetPermissionsInstance() map[string]string {
	if nmPermissions == nil {
		nmPermissions = nmGetPermissions()
	}
	return nmPermissions
}
func nmGetPermissions() (permissions map[string]string) {
	m, err := nmNewManager()
	if err != nil {
		return
	}

	permissions, err = m.GetPermissions(0)
	if err != nil {
		logger.Error(err)
	}
	return
}

func nmGetDevices() (devPaths []dbus.ObjectPath) {
	devPaths, err := nmManager.GetDevices(0)
	if err != nil {
		logger.Error(err)
	}
	return
}

func nmGetDeviceInterface(devPath dbus.ObjectPath) (devInterface string) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}

	devInterface, _ = dev.Interface().Get(0)
	return
}

func nmAddAndActivateConnection(data connectionData, devPath dbus.ObjectPath, forced bool) (cpath, apath dbus.ObjectPath, err error) {
	if len(devPath) == 0 {
		devPath = "/"
	} else {
		if !forced && isWiredDevice(devPath) && !nmGetWiredCarrier(devPath) {
			err = fmt.Errorf("%s", deviceErrorTable[CUSTOM_NM_DEVICE_STATE_REASON_CABLE_UNPLUGGED])
			notifyWiredCableUnplugged()
			return
		}
	}
	spath := dbus.ObjectPath("/")
	cpath, apath, err = nmManager.AddAndActivateConnection(0, data, devPath, spath)
	if err != nil {
		nmHandleActivatingError(data, devPath)
		logger.Error(err, "devPath:", devPath)
		return
	}
	return
}

func nmActivateConnection(cpath, devPath dbus.ObjectPath) (apath dbus.ObjectPath, err error) {
	if isWiredDevice(devPath) && !nmGetWiredCarrier(devPath) {
		err = fmt.Errorf("%s", deviceErrorTable[CUSTOM_NM_DEVICE_STATE_REASON_CABLE_UNPLUGGED])
		notifyWiredCableUnplugged()
		return
	}
	spath := dbus.ObjectPath("/")
	apath, err = nmManager.ActivateConnection(0, cpath, devPath, spath)
	if err != nil {
		if data, err := nmGetConnectionData(cpath); err == nil {
			nmHandleActivatingError(data, devPath)
		}
		logger.Error(err)
		return
	}
	return
}

func nmHandleActivatingError(data connectionData, devPath dbus.ObjectPath) {
	switch nmGetDeviceType(devPath) {
	case nm.NM_DEVICE_TYPE_ETHERNET:
		// if wired cable unplugged, give a notification
		if !isDeviceStateAvailable(nmGetDeviceState(devPath)) {
			notifyWiredCableUnplugged()
		}
	}
	switch getCustomConnectionType(data) {
	case connectionWirelessAdhoc, connectionWirelessHotspot:
		// if connection type is wireless hotspot, give a notification
		notifyApModeNotSupport()
	}
}

func nmDeactivateConnection(apath dbus.ObjectPath) (err error) {
	err = nmManager.DeactivateConnection(0, apath)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

func nmGetActiveConnections() (apaths []dbus.ObjectPath) {
	apaths, _ = nmManager.ActiveConnections().Get(0)
	return
}

func nmGetVpnActiveConnections() (apaths []dbus.ObjectPath) {
	for _, p := range nmGetActiveConnections() {
		if aconn, err := nmNewActiveConnection(p); err == nil {
			vpn, _ := aconn.Vpn().Get(0)
			if vpn {
				apaths = append(apaths, p)
			}
		}
	}
	return
}

func nmGetAccessPoints(devPath dbus.ObjectPath) (apPaths []dbus.ObjectPath) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	devWireless := dev.Wireless()

	apPaths, err = devWireless.AccessPoints().Get(0)
	if err != nil {
		logger.Error(err)
	}
	return
}

func nmGetManagerState() (state uint32) {
	state, _ = nmManager.State().Get(0)
	return
}

func nmGetActiveConnectionByUuid(uuid string) (apaths []dbus.ObjectPath, err error) {
	for _, apath := range nmGetActiveConnections() {
		if aconn, tmperr := nmNewActiveConnection(apath); tmperr == nil {
			aconnUuid, _ := aconn.Uuid().Get(0)
			if aconnUuid == uuid {
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

	state, _ = aconn.State().Get(0)
	return
}

func nmGetConnectionData(cpath dbus.ObjectPath) (data connectionData, err error) {
	nmConn, err := nmNewSettingsConnection(cpath)
	if err != nil {
		return
	}

	data, err = nmConn.GetSettings(0)
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

func nmGetConnectionList() (connections []dbus.ObjectPath) {
	connections, err := nmSettings.ListConnections(0)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

func nmGetConnectionByUuid(uuid string) (cpath dbus.ObjectPath, err error) {
	cpath, err = nmSettings.GetConnectionByUuid(0, uuid)
	return
}

func isWiredDevice(devPath dbus.ObjectPath) bool {
	device, err := nmNewDevice(devPath)
	if err != nil {
		return false
	}

	deviceType, _ := device.DeviceType().Get(0)
	return deviceType == nm.NM_DEVICE_TYPE_ETHERNET
}

func nmGetWiredCarrier(devPath dbus.ObjectPath) bool {
	device, err := nmNewDevice(devPath)
	if err != nil {
		// TODO: 为什么出错了还返回true？
		return true
	}
	wired := device.Wired()
	hwAddress, _ := wired.HwAddress().Get(0)
	carrier, _ := wired.Carrier().Get(0)

	logger.Debug("--------Check wired available:", hwAddress, carrier)
	return carrier
}

func nmGetWirelessConnectionSsidByUuid(uuid string) (ssid []byte) {
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
	cpath, err = nmSettings.AddConnection(0, data)
	if err != nil {
		logger.Error(err)
	}
	return
}

func nmGetIp4ConfigInfo(path dbus.ObjectPath) (address, mask string, gateways, nameServers []string) {
	address = "0.0.0.0"
	mask = "0.0.0.0"
	ip4config, err := nmNewIP4Config(path)
	if err != nil {
		return
	}
	addressProp, _ := ip4config.Addresses().Get(0)

	ipv4Addresses := wrapIpv4Addresses(addressProp)
	if len(ipv4Addresses) > 0 {
		address = ipv4Addresses[0].Address
		mask = ipv4Addresses[0].Mask
	}
	for _, address := range ipv4Addresses {
		gateways = append(gateways, address.Gateway)
	}

	nameServersProp, _ := ip4config.Nameservers().Get(0)
	nameServers = wrapIpv4Dns(nameServersProp)
	return
}

func nmGetIp6ConfigInfo(path dbus.ObjectPath) (address, prefix string, gateways, nameServers []string) {
	address = "0::0"
	prefix = "0"
	ip6config, err := nmNewIP6Config(path)
	if err != nil {
		return
	}

	addressProp, _ := ip6config.Addresses().Get(0)
	gateway, _ := ip6config.Gateway().Get(0)
	gateways = append(gateways, gateway)
	ipv6Addresses := wrapNMDBusIpv6Addresses(addressProp)
	if len(ipv6Addresses) > 0 {
		address = ipv6Addresses[0].Address
		prefix = fmt.Sprintf("%d", ipv6Addresses[0].Prefix)
	}
	for _, addr := range ipv6Addresses {
		gateways = append(gateways, addr.Gateway)
		if addr.Address[:5] != "FE80:" && // link local
			addr.Address[:5] != "FEC0:" { // site local
			address = addr.Address
			prefix = fmt.Sprintf("%d", addr.Prefix)
		}
	}

	nameServersProp, _ := ip6config.Nameservers().Get(0)
	nameServers = wrapIpv6Dns(nameServersProp)
	return
}

func wrapNMDBusIpv6Addresses(data []nmdbus.IP6Address) (wrapData ipv6AddressesWrapper) {
	for _, d := range data {
		ipv6Addr := ipv6AddressWrapper{}
		ipv6Addr.Address = convertIpv6AddressToString(d.Address)
		ipv6Addr.Prefix = d.Prefix
		ipv6Addr.Gateway = convertIpv6AddressToString(d.Gateway)
		wrapData = append(wrapData, ipv6Addr)
	}
	return
}

func nmGetDeviceState(devPath dbus.ObjectPath) (state uint32) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return nm.NM_DEVICE_STATE_UNKNOWN
	}

	state, _ = dev.State().Get(0)
	return
}

func nmSetDeviceAutoconnect(devPath dbus.ObjectPath, autoConnect bool) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	err = dev.Autoconnect().Set(0, autoConnect)
	if err != nil {
		logger.Warning("failed to set autoconnect:", err)
		return
	}
}

func nmSetDeviceManaged(devPath dbus.ObjectPath, managed bool) (err error) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	err = dev.Managed().Set(0, managed)
	if err != nil {
		logger.Warning("failed to set device managed:", err)
		return
	}
	return
}

func nmGetDeviceType(devPath dbus.ObjectPath) (devType uint32) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return nm.NM_DEVICE_TYPE_UNKNOWN
	}

	devType, _ = dev.DeviceType().Get(0)
	return
}

func nmGetDeviceActiveConnection(devPath dbus.ObjectPath) (acPath dbus.ObjectPath) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	acPath, _ = dev.ActiveConnection().Get(0)
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

	aconnConnection, _ := aconn.Connection().Get(0)
	conn, err := nmNewSettingsConnection(aconnConnection)
	if err != nil {
		return
	}

	data, err = conn.GetSettings(0)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

func nmGetPrimaryConnection() (cpath dbus.ObjectPath) {
	cpath, _ = nmManager.PrimaryConnection().Get(0)
	return
}

func nmGetNetworkEnabled() bool {
	enabled, _ := nmManager.NetworkingEnabled().Get(0)
	return enabled
}
func nmGetWirelessHardwareEnabled() bool {
	enabled, _ := nmManager.WirelessHardwareEnabled().Get(0)
	return enabled
}
