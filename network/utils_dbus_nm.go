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
	"fmt"
	"sort"
	"strings"

	"pkg.deepin.io/lib/dbusutil"

	nmdbus "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.networkmanager"

	"pkg.deepin.io/dde/daemon/network/nm"
	"pkg.deepin.io/lib/dbus1"
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

var nmPermissions map[string]string

// Helper functions
func isNmObjectPathValid(p dbus.ObjectPath) bool {
	str := string(p)
	if len(str) == 0 || str == "/" {
		return false
	}
	return true
}

func isNmDeviceObjectExists(devPath dbus.ObjectPath) bool {
	// TODO: 这个方法有问题， 只能判断 devPath 是否是合法path
	_, err := nmNewDevice(devPath)
	if err != nil {
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
	if state > nm.NM_DEVICE_STATE_UNMANAGED {
		return true
	}
	return false
}
func isDeviceStateAvailable(state uint32) bool {
	if state > nm.NM_DEVICE_STATE_UNAVAILABLE {
		return true
	}
	return false
}
func isDeviceStateActivated(state uint32) bool {
	if state == nm.NM_DEVICE_STATE_ACTIVATED {
		return true
	}
	return false
}
func isDeviceStateInActivating(state uint32) bool {
	if state >= nm.NM_DEVICE_STATE_PREPARE && state <= nm.NM_DEVICE_STATE_ACTIVATED {
		return true
	}
	return false
}

func isDeviceStateReasonInvalid(reason uint32) bool {
	switch reason {
	case nm.NM_DEVICE_STATE_REASON_UNKNOWN, nm.NM_DEVICE_STATE_REASON_NONE:
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
func isConnectionStateActivated(state uint32) bool {
	if state == nm.NM_ACTIVE_CONNECTION_STATE_ACTIVATED {
		return true
	}
	return false
}
func isConnectionStateInDeactivating(state uint32) bool {
	if state == nm.NM_ACTIVE_CONNECTION_STATE_DEACTIVATING ||
		state == nm.NM_ACTIVE_CONNECTION_STATE_DEACTIVATED {
		return true
	}
	return false
}
func isConnectionStateDeactivate(state uint32) bool {
	if state == nm.NM_ACTIVE_CONNECTION_STATE_DEACTIVATED {
		return true
	}
	return false
}

func isConnectionInActivating(uuid string) bool {
	apaths, err := nmGetActiveConnectionByUuid(uuid)
	if err != nil || len(apaths) == 0 {
		return false
	}
	for _, apath := range apaths {
		if isConnectionStateInActivating(nmGetActiveConnectionState(apath)) {
			return true
		}
	}
	return false
}

// check if vpn connection activating or activated
func isVpnConnectionStateInActivating(state uint32) bool {
	if state >= nm.NM_VPN_CONNECTION_STATE_PREPARE &&
		state <= nm.NM_VPN_CONNECTION_STATE_ACTIVATED {
		return true
	}
	return false
}
func isVpnConnectionStateActivated(state uint32) bool {
	if state == nm.NM_VPN_CONNECTION_STATE_ACTIVATED {
		return true
	}
	return false
}
func isVpnConnectionStateDeactivate(state uint32) bool {
	if state == nm.NM_VPN_CONNECTION_STATE_DISCONNECTED {
		return true
	}
	return false
}
func isVpnConnectionStateFailed(state uint32) bool {
	if state == nm.NM_VPN_CONNECTION_STATE_FAILED {
		return true
	}
	return false
}

var availableValuesSettingSecretFlags []kvalue

func initAvailableValuesSecretFlags() {
	availableValuesSettingSecretFlags = []kvalue{
		kvalue{nm.NM_SETTING_SECRET_FLAG_NONE, Tr("Saved")}, // system saved
		// kvalue{nm.NM_SETTING_SECRET_FLAG_AGENT_OWNED, Tr("Saved")},
		kvalue{nm.NM_SETTING_SECRET_FLAG_NOT_SAVED, Tr("Always Ask")},
		kvalue{nm.NM_SETTING_SECRET_FLAG_NOT_REQUIRED, Tr("Not Required")},
	}
}

func isSettingRequireSecret(flag uint32) bool {
	if flag == nm.NM_SETTING_SECRET_FLAG_NONE || flag == nm.NM_SETTING_SECRET_FLAG_AGENT_OWNED {
		return true
	}
	return false
}

func isVirtualDeviceIfc(dev *nmdbus.Device) bool {
	driver, _ := dev.Driver().Get(0)
	switch driver {
	case "dummy", "veth", "vboxnet", "vmnet", "vmxnet", "vmxnet2", "vmxnet3":
		return true
	case "unknown":
		// sometimes we could not get vmnet dirver name, so check the
		// udi sys path if is prefix with /sys/devices/virtual/net
		devUdi, _ := dev.Udi().Get(0)
		devInterface, _ := dev.Interface().Get(0)
		if strings.HasPrefix(devUdi, "/sys/devices/virtual/net") ||
			strings.HasPrefix(devUdi, "/virtual/device") ||
			strings.HasPrefix(devInterface, "vmnet") {
			return true
		}
	}
	return false
}

// General function wrappers for network manager
func nmGeneralGetAllDeviceHwAddr(devType uint32) (allHwAddr map[string]string) {
	allHwAddr = make(map[string]string)
	for _, devPath := range nmGetDevices() {
		dev, err := nmNewDevice(devPath)
		if err != nil {
			continue
		}
		deviceType, _ := dev.DeviceType().Get(0)

		if deviceType == devType {
			hwAddr, err := nmGeneralGetDeviceHwAddr(devPath, true)
			// filter all virtual devices
			if err == nil && !isVirtualDeviceIfc(dev) {
				devInterface, _ := dev.Interface().Get(0)
				allHwAddr[devInterface] = hwAddr
			}
		}
	}
	return
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
			hwAddr, _ = devWired.PermHwAddress().Get(0)
		}
		if len(hwAddr) == 0 {
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

func nmGetDeviceIdentifiers() (devIds []string) {
	for _, devPath := range nmGetDevices() {
		id, err := nmGeneralGetDeviceIdentifier(devPath)
		if err == nil {
			devIds = append(devIds, id)
		}
	}
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
		speed = bitRate / 1024
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
	sysPath, err := nmGeneralGetDeviceSysPath(devPath)
	if err != nil {
		return false
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
func nmNewAgentManager() (manager *nmdbus.AgentManager, err error) {
	systemBus, err := dbus.SystemBus()
	if err != nil {
		return
	}
	manager = nmdbus.NewAgentManager(systemBus)
	return
}
func nmNewDHCP4Config(path dbus.ObjectPath) (dhcp4 *nmdbus.Dhcp4Config, err error) {
	systemBus, err := dbus.SystemBus()
	if err != nil {
		return
	}
	dhcp4, err = nmdbus.NewDhcp4Config(systemBus, path)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}
func nmNewDHCP6Config(path dbus.ObjectPath) (dhcp6 *nmdbus.Dhcp6Config, err error) {
	systemBus, err := dbus.SystemBus()
	if err != nil {
		return
	}
	dhcp6, err = nmdbus.NewDhcp6Config(systemBus, path)
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
func nmNewVpnConnection(apath dbus.ObjectPath) (vpnConn *nmdbus.VpnConnection, err error) {
	systemBus, err := dbus.SystemBus()
	if err != nil {
		return
	}
	vpnConn, err =
		nmdbus.NewVpnConnection(systemBus, apath)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

// Destroy network manager objects
func nmDestroyManager(m *nmdbus.Manager) {
	if m == nil {
		logger.Error("Manager to destroy is nil")
		return
	}
	m.RemoveHandler(proxy.RemoveAllHandlers)
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

func nmDestroyActiveConnection(aconn *nmdbus.ActiveConnection) {
	if aconn == nil {
		logger.Error("ActiveConnection to destroy is nil")
		return
	}
	aconn.RemoveHandler(proxy.RemoveAllHandlers)
}

func nmDestroyVpnConnection(vpnConn *nmdbus.VpnConnection) {
	if vpnConn == nil {
		logger.Error("ActiveConnection to destroy is nil")
		return
	}
	vpnConn.RemoveHandler(proxy.RemoveAllHandlers)
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

func nmAgentRegister(identifier string) {
	am, err := nmNewAgentManager()
	if err != nil {
		return
	}
	err = am.Register(0, identifier)
	if err != nil {
		logger.Error(err)
	}
}

func nmAgentUnregister() {
	am, err := nmNewAgentManager()
	if err != nil {
		return
	}
	err = am.Unregister(0)
	if err != nil {
		logger.Error(err)
	}
}

func nmGetDevices() (devPaths []dbus.ObjectPath) {
	devPaths, err := nmManager.GetDevices(0)
	if err != nil {
		logger.Error(err)
	}
	return
}

func nmGetDevicesByType(devType uint32) (specDevPaths []dbus.ObjectPath) {
	for _, p := range nmGetDevices() {
		if dev, err := nmNewDevice(p); err == nil {
			deviceType, _ := dev.DeviceType().Get(0)
			if deviceType == devType {
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

	devInterface, _ = dev.Interface().Get(0)
	return
}

func nmGetDeviceModemCapabilities(devPath dbus.ObjectPath) (capabilities uint32) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	devModem := dev.Modem()

	capabilities, _ = devModem.CurrentCapabilities().Get(0)
	return
}

func nmAddAndActivateConnection(data connectionData, devPath dbus.ObjectPath, forced bool) (cpath, apath dbus.ObjectPath, err error) {
	if len(devPath) == 0 {
		devPath = "/"
	} else {
		if !forced && isWiredDevice(devPath) && !nmGetWiredCarrier(devPath) {
			err = fmt.Errorf("%s", deviceErrorTable[CUSTOM_NM_DEVICE_STATE_REASON_CABLE_UNPLUGGED])
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

func nmGetVpnConnectionState(apath dbus.ObjectPath) (state uint32) {
	vpnConn, err := nmNewVpnConnection(apath)
	if err != nil {
		return
	}

	state, _ = vpnConn.VpnState().Get(0)
	return
}

func nmRequestWirelessScan(devPath dbus.ObjectPath) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	devWireless := dev.Wireless()

	options := make(map[string]dbus.Variant)
	err = devWireless.RequestScan(0, options)
	if err != nil {
		logger.Error(err)
	}
	return
}

func nmGetAccessPoints(devPath dbus.ObjectPath) (apPaths []dbus.ObjectPath) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	devWireless := dev.Wireless()

	apPaths, err = devWireless.GetAccessPoints(0)
	if err != nil {
		logger.Error(err)
	}
	return
}

func nmGetAccessPointSsids(devPath dbus.ObjectPath) (ssids []string) {
	for _, apPath := range nmGetAccessPoints(devPath) {
		if ap, err := nmNewAccessPoint(apPath); err == nil {
			ssid, _ := ap.Ssid().Get(0)
			ssids = append(ssids, decodeSsid(ssid))
		}
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

func nmGetActiveConnectionVpn(apath dbus.ObjectPath) (isVpn bool) {
	aconn, err := nmNewActiveConnection(apath)
	if err != nil {
		return
	}

	isVpn, _ = aconn.Vpn().Get(0)
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

func nmUpdateConnectionData(cpath dbus.ObjectPath, data connectionData) (err error) {
	nmConn, err := nmNewSettingsConnection(cpath)
	if err != nil {
		return
	}

	correctConnectionData(data)
	err = nmConn.Update(0, data)
	if err != nil {
		logger.Error(err)
	}
	return
}

func nmGetConnectionSecrets(cpath dbus.ObjectPath, secretField string) (secrets connectionData, err error) {
	nmConn, err := nmNewSettingsConnection(cpath)
	if err != nil {
		return
	}

	secrets, err = nmConn.GetSecrets(0, secretField)
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
	connections, err := nmSettings.ListConnections(0)
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

func nmGetConnectionUuidsByType(connTypes ...string) (uuids []string) {
	for _, cpath := range nmGetConnectionList() {
		if isStringInArray(nmGetConnectionType(cpath), connTypes) {
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

func nmGetOtherConnectionIds(origUuid string) (ids []string) {
	for _, cpath := range nmGetConnectionList() {
		if uuid, _ := nmGetConnectionUuid(cpath); uuid != origUuid {
			ids = append(ids, nmGetConnectionId(cpath))
		}
	}
	return
}

func isConnAutoConnect(uuid string) bool {
	connPath, err := nmGetConnectionByUuid(uuid)
	if err != nil {
		logger.Warning(err)
		return false
	}

	return nmGetConnectionAutoconnect(connPath)
}

// TODO: dispatch connection permission
func nmGetAddressableConnectionIds() (ids []string) {
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

// get wireless connection by ssid, the connection with special hardware address is priority
func nmGetWirelessConnection(ssid []byte, devPath dbus.ObjectPath) (cpath dbus.ObjectPath, ok bool) {
	var hwAddr string
	if len(devPath) != 0 {
		hwAddr, _ = nmGeneralGetDeviceHwAddr(devPath, true)
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

// TODO: remove, use nmGetIp4ConfigInfo instead
func nmGetDhcp4Info(path dbus.ObjectPath) (ip, mask string, routers, nameServers []string) {
	ip = "0.0.0.0"
	mask = "0.0.0.0"
	routers = make([]string, 0)
	nameServers = make([]string, 0)

	dhcp4, err := nmNewDHCP4Config(path)
	if err != nil {
		return
	}

	options, _ := dhcp4.Options().Get(0)
	if ipData, ok := options["ip_address"]; ok {
		ip, _ = ipData.Value().(string)
	}
	if maskData, ok := options["subnet_mask"]; ok {
		mask, _ = maskData.Value().(string)
	}
	if routersData, ok := options["routers"]; ok {
		routersStr, _ := routersData.Value().(string)
		if len(routersStr) > 0 {
			routers = strings.Split(routersStr, " ")
		}
	}
	if nameServersData, ok := options["domain_name_servers"]; ok {
		nameServersStr, _ := nameServersData.Value().(string)
		if len(nameServersStr) > 0 {
			nameServers = strings.Split(nameServersStr, " ")
		}
	}
	return
}

// TODO: remove, use nmGetIp6ConfigInfo instead
func nmGetDhcp6Info(path dbus.ObjectPath) (ip string, routers, nameServers []string) {
	ip = "0::0"
	routers = make([]string, 0)
	nameServers = make([]string, 0)

	dhcp6, err := nmNewDHCP6Config(path)
	if err != nil {
		return
	}

	options, _ := dhcp6.Options().Get(0)
	if ipData, ok := options["ip6_address"]; ok {
		ip, _ = ipData.Value().(string)
	}
	if routersData, ok := options["routers"]; ok {
		routersStr, _ := routersData.Value().(string)
		if len(routersStr) > 0 {
			routers = strings.Split(routersStr, " ")
		}
	}
	if nameServersData, ok := options["dhcp6_name_servers"]; ok {
		nameServersStr, _ := nameServersData.Value().(string)
		if len(nameServersStr) > 0 {
			nameServers = strings.Split(nameServersStr, " ")
		}
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
	ipv6Addresses := wrapNMDBusIpv6Addresses(addressProp)
	if len(ipv6Addresses) > 0 {
		address = ipv6Addresses[0].Address
		prefix = fmt.Sprintf("%d", ipv6Addresses[0].Prefix)
	}
	for _, address := range ipv6Addresses {
		gateways = append(gateways, address.Gateway)
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

func nmGetDeviceAutoconnect(devPath dbus.ObjectPath) (autoConnect bool) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}

	autoConnect, _ = dev.Autoconnect().Get(0)
	return
}

func nmSetDeviceAutoconnect(devPath dbus.ObjectPath, autoConnect bool) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	dev.Autoconnect().Set(0, autoConnect)
	return
}

func nmSetDeviceManaged(devPath dbus.ObjectPath, managed bool) (err error) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	dev.Managed().Set(0, managed)
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

func nmGetDeviceUdi(devPath dbus.ObjectPath) (udi string) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	udi, _ = dev.Udi().Get(0)
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

func nmGetDeviceAvailableConnections(devPath dbus.ObjectPath) (paths []dbus.ObjectPath) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}
	paths, _ = dev.AvailableConnections().Get(0)
	return
}

func nmGetDeviceActiveConnectionUuid(devPath dbus.ObjectPath) (uuid string, err error) {
	acPath := nmGetDeviceActiveConnection(devPath)
	aconn, err := nmNewActiveConnection(acPath)
	if err != nil {
		return
	}

	uuid, err = aconn.Uuid().Get(0)
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

func nmManagerEnable(enable bool) (err error) {
	err = nmManager.Enable(0, enable)
	if err != nil {
		logger.Error(err)
	}
	return
}

func nmGetPrimaryConnection() (cpath dbus.ObjectPath) {
	cpath, _ = nmManager.PrimaryConnection().Get(0)
	return
}

func nmGetNetworkState() uint32 {
	state, _ := nmManager.State().Get(0)
	return state
}
func nmIsNetworkOffline() bool {
	state, _ := nmManager.State().Get(0)
	if state == nm.NM_STATE_DISCONNECTED || state == nm.NM_STATE_ASLEEP {
		return true
	}
	return false
}

func nmGetNetworkEnabled() bool {
	enabled, _ := nmManager.NetworkingEnabled().Get(0)
	return enabled
}
func nmGetWirelessHardwareEnabled() bool {
	enabled, _ := nmManager.WirelessHardwareEnabled().Get(0)
	return enabled
}
func nmGetWirelessEnabled() bool {
	enabled, _ := nmManager.WirelessEnabled().Get(0)
	return enabled
}

func nmSetNetworkingEnabled(enabled bool) {
	if nmGetNetworkEnabled() != enabled {
		nmManagerEnable(enabled)
	} else {
		logger.Warning("NetworkingEnabled already set as", enabled)
	}
	return
}

func nmSetWirelessEnabled(enabled bool) {
	currentEnabled, err := nmManager.WirelessEnabled().Get(0)
	if err != nil {
		logger.Warning(err)
		return
	}
	if currentEnabled != enabled {
		nmManager.WirelessEnabled().Set(0, enabled)
	} else {
		logger.Warning("WirelessEnabled already set as", enabled)
	}
	return
}

func nmSetWwanEnabled(enabled bool) {
	currentEnabled, err := nmManager.WwanEnabled().Get(0)
	if err != nil {
		logger.Warning(err)
		return
	}
	if currentEnabled != enabled {
		nmManager.WwanEnabled().Set(0, enabled)
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
	devRelatedUuid := nmGeneralGetDeviceUniqueUuid(devPath)
	for _, cpath := range nmGetDeviceAvailableConnections(devPath) {
		if cdata, err := nmGetConnectionData(cpath); err == nil {
			uuid := getSettingConnectionUuid(cdata)
			switch getCustomConnectionType(cdata) {
			case connectionWired, connectionMobileGsm, connectionMobileCdma:
				if devRelatedUuid != uuid {
					// ignore connections that not matching the
					// device, etc wired connections that create in
					// other ways
					continue
				}
			}
			if uuid == lastConnectionUuid {
				// the last activated connection will be dispatch
				// specially
				continue
			}
			if getSettingConnectionAutoconnect(cdata) {
				id := getSettingConnectionId(cdata)
				timestamp := getSettingConnectionTimestamp(cdata)
				if timestamp > 0 {
					// only collect connections that connected before
					ac := autoConnectConn{
						id:        id,
						uuid:      uuid,
						timestamp: timestamp,
					}
					acs = append(acs, ac)
				}
			}
		}
	}
	sort.Sort(sort.Reverse(acs))
	logger.Debugf("autoconnect connections for device type %s, %v",
		getCustomDeviceType(nmGetDeviceType(devPath)), acs)
	if len(lastConnectionUuid) > 0 {
		// the last activated connection has the highest priority if
		// exists and the auto-connect property enabled
		if cpath, err := nmGetConnectionByUuid(lastConnectionUuid); err == nil {
			if nmGetConnectionAutoconnect(cpath) {
				uuids = []string{lastConnectionUuid}
			}
		}
	}
	for _, ac := range acs {
		uuids = append(uuids, ac.uuid)
	}
	return
}

func (m *Manager) nmRunOnceUntilDeviceAvailable(devPath dbus.ObjectPath, cb func()) {
	dev, err := nmNewDevice(devPath)
	if err != nil {
		return
	}

	state, _ := dev.State().Get(0)
	if isDeviceStateAvailable(state) {
		cb()
	} else {
		hasRun := false
		dev.InitSignalExt(m.sysSigLoop, true)
		dev.ConnectStateChanged(func(newState uint32, oldState uint32, reason uint32) {
			if !hasRun && isDeviceStateAvailable(newState) {
				cb()
				nmDestroyDevice(dev)
				hasRun = true
			}
		})
	}
}

func nmRunOnceUtilNetworkAvailable(sysSigLoop *dbusutil.SignalLoop, cb func()) {
	manager, err := nmNewManager()
	if err != nil {
		return
	}
	state, _ := manager.State().Get(0)
	const connectedState uint32 = nm.NM_STATE_CONNECTED_LOCAL
	if state >= connectedState {
		cb()
	} else {
		hasRun := false
		manager.InitSignalExt(sysSigLoop, true)
		manager.ConnectStateChanged(func(state uint32) {
			if !hasRun && state >= connectedState {
				cb()
				nmDestroyManager(manager)
				hasRun = true
			}
		})
	}
}
