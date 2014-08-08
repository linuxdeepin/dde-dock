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

import "pkg.linuxdeepin.com/lib/dbus"
import "sync"
import . "pkg.linuxdeepin.com/lib/gettext"
import nm "dbus/org/freedesktop/networkmanager"

const (
	notifyIconNetworkConnected     = "network-transmit-receive"
	notifyIconNetworkDisconnected  = "network-error"
	notifyIconNetworkOffline       = "network-offline"
	notifyIconEthernetConnected    = "notification-network-ethernet-connected"
	notifyIconEthernetDisconnected = "notification-network-ethernet-disconnected"
	notifyIconWirelessConnected    = "notification-network-wireless-full"
	notifyIconWirelessDisconnected = "notification-network-wireless-disconnected"
	notifyIconVpnConnected         = "notification-network-vpn-connected"
	notifyIconVpnDisconnected      = "notification-network-vpn-disconnected"
)

var vpnErrorTable = make(map[uint32]string)
var deviceErrorTable = make(map[uint32]string)

func initNmStateReasons() {
	deviceErrorTable[NM_DEVICE_STATE_REASON_UNKNOWN] = Tr("Device state changed, unknown reason.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_NONE] = Tr("Device state changed, none reason.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_NOW_MANAGED] = Tr("The device is now managed.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_NOW_UNMANAGED] = Tr("The device is no longer managed.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_CONFIG_FAILED] = Tr("The device has not been ready for configuration.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_CONFIG_UNAVAILABLE] = Tr("IP configuration could not be reserved (no available address, timeout, etc).")
	deviceErrorTable[NM_DEVICE_STATE_REASON_CONFIG_EXPIRED] = Tr("The IP configuration is no longer valid.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_NO_SECRETS] = Tr("Passwords were required but not provided.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_SUPPLICANT_DISCONNECT] = Tr("The 802.1X supplicant disconnected from the access point or authentication server.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_SUPPLICANT_CONFIG_FAILED] = Tr("Configuration of the 802.1X supplicant failed.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_SUPPLICANT_FAILED] = Tr("The 802.1X supplicant quitted or failed unexpectedly.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_SUPPLICANT_TIMEOUT] = Tr("The 802.1X supplicant took too long time to authenticate.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_PPP_START_FAILED] = Tr("The PPP service failed to start within the allowed time.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_PPP_DISCONNECT] = Tr("The PPP service disconnected unexpectedly.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_PPP_FAILED] = Tr("The PPP service quitted or failed unexpectedly.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_DHCP_START_FAILED] = Tr("The DHCP service failed to start within the allowed time.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_DHCP_ERROR] = Tr("The DHCP service reported an unexpected error.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_DHCP_FAILED] = Tr("The DHCP service quitted or failed unexpectedly.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_SHARED_START_FAILED] = Tr("The shared connection service failed to start.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_SHARED_FAILED] = Tr("The shared connection service quitted or failed unexpectedly.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_AUTOIP_START_FAILED] = Tr("The AutoIP service failed to start.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_AUTOIP_ERROR] = Tr("The AutoIP service reported an unexpected error.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_AUTOIP_FAILED] = Tr("The AutoIP service quitted or failed unexpectedly.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_MODEM_BUSY] = Tr("Dialing failed due to busy lines.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_MODEM_NO_DIAL_TONE] = Tr("Dialing failed due to no dial tone.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_MODEM_NO_CARRIER] = Tr("Dialing failed due to the carrier.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_MODEM_DIAL_TIMEOUT] = Tr("Dialing timed out.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_MODEM_DIAL_FAILED] = Tr("Dialing failed.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_MODEM_INIT_FAILED] = Tr("Modem initialization failed.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_GSM_APN_FAILED] = Tr("Failed to select the specified GSM APN.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_GSM_REGISTRATION_NOT_SEARCHING] = Tr("No networks searched.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_GSM_REGISTRATION_DENIED] = Tr("Network registration was denied.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_GSM_REGISTRATION_TIMEOUT] = Tr("Network registration timed out.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_GSM_REGISTRATION_FAILED] = Tr("Register to the requested GSM network failed.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_GSM_PIN_CHECK_FAILED] = Tr("PIN check failed.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_FIRMWARE_MISSING] = Tr("Necessary firmware for the device may be missed.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_REMOVED] = Tr("The device was removed.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_SLEEPING] = Tr("NetworkManager went to sleep.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_CONNECTION_REMOVED] = Tr("The device's active connection was removed or disappeared.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_USER_REQUESTED] = Tr("A user or client requested to disconnect.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_CARRIER] = Tr("The device's carrier/link changed.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_CONNECTION_ASSUMED] = Tr("The device's existing connection was assumed.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_SUPPLICANT_AVAILABLE] = Tr("The 802.1x supplicant is now available.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_MODEM_NOT_FOUND] = Tr("The modem could not be found.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_BT_FAILED] = Tr("The Bluetooth connection timed out or failed.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_GSM_SIM_NOT_INSERTED] = Tr("GSM Modem's SIM Card was not inserted.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_GSM_SIM_PIN_REQUIRED] = Tr("GSM Modem's SIM PIN required.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_GSM_SIM_PUK_REQUIRED] = Tr("GSM Modem's SIM PUK required.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_GSM_SIM_WRONG] = Tr("Wrong GSM Modem's SIM")
	deviceErrorTable[NM_DEVICE_STATE_REASON_INFINIBAND_MODE] = Tr("InfiniBand device does not support connected mode.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_DEPENDENCY_FAILED] = Tr("A connection dependency failed.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_BR2684_FAILED] = Tr("Problem with the RFC 2684 Ethernet over ADSL bridge.") // TODO
	deviceErrorTable[NM_DEVICE_STATE_REASON_MODEM_MANAGER_UNAVAILABLE] = Tr("ModemManager was not running or quitted unexpectedly.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_SSID_NOT_FOUND] = Tr("The 802.11 Wi-Fi network could not be found.")
	deviceErrorTable[NM_DEVICE_STATE_REASON_SECONDARY_CONNECTION_FAILED] = Tr("A secondary connection of the base connection failed.")

	vpnErrorTable[NM_VPN_CONNECTION_STATE_REASON_UNKNOWN] = Tr("Activate VPN connection failed, unknown reason.")
	vpnErrorTable[NM_VPN_CONNECTION_STATE_REASON_NONE] = Tr("Activate VPN connection failed.")
	vpnErrorTable[NM_VPN_CONNECTION_STATE_REASON_USER_DISCONNECTED] = Tr("The VPN connection changed state due to disconnection from users.")
	vpnErrorTable[NM_VPN_CONNECTION_STATE_REASON_DEVICE_DISCONNECTED] = Tr("The VPN connection changed state due to disconnection from devices.")
	vpnErrorTable[NM_VPN_CONNECTION_STATE_REASON_SERVICE_STOPPED] = Tr("VPN service stopped.")
	vpnErrorTable[NM_VPN_CONNECTION_STATE_REASON_IP_CONFIG_INVALID] = Tr("The IP config of VPN connection was invalid.")
	vpnErrorTable[NM_VPN_CONNECTION_STATE_REASON_CONNECT_TIMEOUT] = Tr("The connection attempt to VPN service timed out.")
	vpnErrorTable[NM_VPN_CONNECTION_STATE_REASON_SERVICE_START_TIMEOUT] = Tr("The VPN service start timed out.")
	vpnErrorTable[NM_VPN_CONNECTION_STATE_REASON_SERVICE_START_FAILED] = Tr("The VPN service failed to start.")
	vpnErrorTable[NM_VPN_CONNECTION_STATE_REASON_NO_SECRETS] = Tr("Necessary password for the VPN connection was not provided.")
	vpnErrorTable[NM_VPN_CONNECTION_STATE_REASON_LOGIN_FAILED] = Tr("Authentication to VPN server failed.")
	vpnErrorTable[NM_VPN_CONNECTION_STATE_REASON_CONNECTION_REMOVED] = Tr("The connection was deleted from settings.")
}

type stateNotifier struct {
	devices map[dbus.ObjectPath]*deviceStateStruct
	locker  sync.Mutex
}
type deviceStateStruct struct {
	nmDev   *nm.Device
	aconnId string
}

func newStateNotifier() (sn *stateNotifier) {
	sn = &stateNotifier{}
	sn.devices = make(map[dbus.ObjectPath]*deviceStateStruct)

	var watch = func(path dbus.ObjectPath) {
		defer func() {
			if err := recover(); err != nil {
				sn.locker.Lock()
				defer sn.locker.Unlock()
				delete(sn.devices, path)
				logger.Error(err)
			}
		}()
		if dev, err := nmNewDevice(path); err == nil {
			sn.locker.Lock()
			defer sn.locker.Unlock()
			sn.devices[path] = &deviceStateStruct{nmDev: dev}
			if data, err := nmGetDeviceActiveConnectionData(path); err == nil {
				// remember active connection id if exists
				sn.devices[path].aconnId = getSettingConnectionId(data)
			}
			// connect signals
			dev.ConnectStateChanged(func(newState, oldState, reason uint32) {
				switch newState {
				case NM_DEVICE_STATE_PREPARE:
					if data, err := nmGetDeviceActiveConnectionData(path); err == nil {
						sn.devices[path].aconnId = getSettingConnectionId(data)
					}
				case NM_DEVICE_STATE_ACTIVATED:
					if data, err := nmGetDeviceActiveConnectionData(path); err == nil {
						var icon, msg string
						switch getCustomConnectionType(data) {
						case connectionWired:
							icon = notifyIconEthernetConnected
						case connectionWireless:
							icon = notifyIconWirelessConnected
						default:
							icon = notifyIconNetworkConnected
						}
						msg = sn.devices[path].aconnId
						notify(icon, Tr("Connected"), msg)
					}
				case NM_DEVICE_STATE_FAILED, NM_DEVICE_STATE_DISCONNECTED,
					NM_DEVICE_STATE_UNMANAGED, NM_DEVICE_STATE_UNAVAILABLE:
					var icon, msg string
					switch dev.DeviceType.Get() {
					case NM_DEVICE_TYPE_ETHERNET:
						icon = notifyIconEthernetDisconnected
					case NM_DEVICE_TYPE_WIFI:
						icon = notifyIconWirelessDisconnected
					default:
						icon = notifyIconWirelessDisconnected
					}
					if newState == NM_DEVICE_STATE_DISCONNECTED {
						msg = sn.devices[path].aconnId
					} else {
						msg = deviceErrorTable[reason]
					}
					notify(icon, Tr("Disconnected"), msg)
				}
			})
		}
	}
	var remove = func(path dbus.ObjectPath) {
		sn.locker.Lock()
		defer sn.locker.Unlock()
		if dev, ok := sn.devices[path]; ok {
			nmDestroyDevice(dev.nmDev)
		}
	}

	nmManager.ConnectDeviceRemoved(func(path dbus.ObjectPath) {
		remove(path)
	})
	nmManager.ConnectDeviceAdded(func(path dbus.ObjectPath) {
		watch(path)
	})

	for _, path := range nmGetDevices() {
		watch(path)
	}

	nmManager.ConnectStateChanged(func(state uint32) {
		switch state {
		case NM_STATE_DISCONNECTED, NM_STATE_ASLEEP:
			notify(notifyIconNetworkOffline, Tr("Offline"), Tr("Disconnected, you are now offline."))
		}
	})

	return
}

func destroyStateNotifier(n *stateNotifier) {
	n.locker.Lock()
	defer n.locker.Unlock()
	for _, dev := range n.devices {
		nmDestroyDevice(dev.nmDev)
	}
	n.devices = nil
}

func notifyVpnConnected(id string) {
	icon := notifyIconVpnConnected
	notify(icon, Tr("Connected"), id)
}
func notifyVpnDisconnected(id string) {
	icon := notifyIconVpnDisconnected
	notify(icon, Tr("Disconnected"), id)
}
func notifyVpnFailed(id string, reason uint32) {
	icon := notifyIconVpnDisconnected
	msg := vpnErrorTable[reason]
	notify(icon, Tr("Disconnected"), msg)
}

func notifyApModeNotSupport() {
	icon := notifyIconWirelessDisconnected
	notify(icon, Tr("Disconnected"), Tr("Access Point mode is not supported by this device."))
}
