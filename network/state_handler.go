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

import "dbus/org/freedesktop/notifications"
import "dlib/dbus"
import "sync"
import nm "dbus/org/freedesktop/networkmanager"
import . "dlib/gettext"

var VpnErrorTable = make(map[uint32]string)
var DeviceErrorTable = make(map[uint32]string)

func initNmStateReasons() {
	DeviceErrorTable[NM_DEVICE_STATE_REASON_NOW_MANAGED] = Tr("The device is now managed.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_NOW_UNMANAGED] = Tr("The device is no longer managed.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_CONFIG_FAILED] = Tr("The device have not been ready for configuration.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_CONFIG_UNAVAILABLE] = Tr("IP configuration could not be reserved (no available address timeout etc).")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_CONFIG_EXPIRED] = Tr("The IP configuration is no longer valid.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_NO_SECRETS] = Tr("Passwords were required but not provided.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_SUPPLICANT_DISCONNECT] = Tr("The 802.1X supplication disconnected from the access point or authentication server.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_SUPPLICANT_CONFIG_FAILED] = Tr("Configuration of the 802.1X supplication failed.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_SUPPLICANT_FAILED] = Tr("The 802.1X supplicant quitted or failed unexpectedly.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_SUPPLICANT_TIMEOUT] = Tr("The 802.1X supplicant took too long time to authenticate.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_PPP_START_FAILED] = Tr("The PPP service failed to start within the allowed time.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_PPP_DISCONNECT] = Tr("The PPP service disconnected unexpectedly.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_PPP_FAILED] = Tr("The PPP service quitted or failed unexpectedly.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_DHCP_START_FAILED] = Tr("The DHCP service failed to start within the allowed time.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_DHCP_ERROR] = Tr("The DHCP service reported an unexpected error.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_DHCP_FAILED] = Tr("The DHCP service quitted or failed unexpectedly.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_SHARED_START_FAILED] = Tr("The shared connection service failed to start.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_SHARED_FAILED] = Tr("The shared connection service quitted or failed unexpectedly.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_AUTOIP_START_FAILED] = Tr("The AutoIP service failed to start.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_AUTOIP_ERROR] = Tr("The AutoIP service reported an unexpected error.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_AUTOIP_FAILED] = Tr("The AutoIP service quitted or failed unexpectedly.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_MODEM_BUSY] = Tr("Dialing failed due to busy lines.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_MODEM_NO_DIAL_TONE] = Tr("Dialing failed due to no dial tone.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_MODEM_NO_CARRIER] = Tr("Dialing failed due to the carrier.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_MODEM_DIAL_TIMEOUT] = Tr("Dialing timed out.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_MODEM_DIAL_FAILED] = Tr("Dialing failed.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_MODEM_INIT_FAILED] = Tr("Modem initialization failed.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_GSM_APN_FAILED] = Tr("Failed to select the specified GSM APN.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_GSM_REGISTRATION_NOT_SEARCHING] = Tr("No networks searched.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_GSM_REGISTRATION_DENIED] = Tr("Network registration was denied.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_GSM_REGISTRATION_TIMEOUT] = Tr("Network registration timed out.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_GSM_REGISTRATION_FAILED] = Tr("Register to the requested GSM network failed.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_GSM_PIN_CHECK_FAILED] = Tr("PIN check failed.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_FIRMWARE_MISSING] = Tr("Necessary firmware for the device may be missed.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_REMOVED] = Tr("The device was removed.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_SLEEPING] = Tr("NetworkManager went to sleep.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_CONNECTION_REMOVED] = Tr("The device's active connection was removed or disappeared.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_USER_REQUESTED] = Tr("An user or client requested to disconnect.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_CARRIER] = Tr("The device's carrier/link changed.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_CONNECTION_ASSUMED] = Tr("The device's existing connection was assumed.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_SUPPLICANT_AVAILABLE] = Tr("The 802.1x supplication is now available.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_MODEM_NOT_FOUND] = Tr("The modem could not be found.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_BT_FAILED] = Tr("The Bluetooth connection timed out or failed.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_GSM_SIM_NOT_INSERTED] = Tr("GSM Modem's SIM Card not inserted.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_GSM_SIM_PIN_REQUIRED] = Tr("GSM Modem's SIM Pin required.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_GSM_SIM_PUK_REQUIRED] = Tr("GSM Modem's SIM Puk required.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_GSM_SIM_WRONG] = Tr("Wrong GSM Modem's SIM")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_INFINIBAND_MODE] = Tr("InfiniBand device does not support connected mode.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_DEPENDENCY_FAILED] = Tr("A connection dependency failed.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_BR2684_FAILED] = Tr("Problem with the RFC 2684 Ethernet over ADSL bridge.") // TODO
	DeviceErrorTable[NM_DEVICE_STATE_REASON_MODEM_MANAGER_UNAVAILABLE] = Tr("ModemManager was not running or quitted unexpectedly.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_SSID_NOT_FOUND] = Tr("The 802.11 Wi-Fi network could not be found.")
	DeviceErrorTable[NM_DEVICE_STATE_REASON_SECONDARY_CONNECTION_FAILED] = Tr("A secondary connection of the base connection failed.")

	VpnErrorTable[NM_VPN_CONNECTION_STATE_REASON_USER_DISCONNECTED] = Tr("VPN connection changed state due to disconnection from users.")
	VpnErrorTable[NM_VPN_CONNECTION_STATE_REASON_DEVICE_DISCONNECTED] = Tr("VPN connection changed state due to disconnection from device.")
	VpnErrorTable[NM_VPN_CONNECTION_STATE_REASON_SERVICE_STOPPED] = Tr("VPN connection service stopped.")
	VpnErrorTable[NM_VPN_CONNECTION_STATE_REASON_IP_CONFIG_INVALID] = Tr("IP config of VPN connection was invalid.")
	VpnErrorTable[NM_VPN_CONNECTION_STATE_REASON_CONNECT_TIMEOUT] = Tr("The connection attempt to VPN  service timed out.")
	VpnErrorTable[NM_VPN_CONNECTION_STATE_REASON_SERVICE_START_TIMEOUT] = Tr("VPN connection service start timed out.")
	VpnErrorTable[NM_VPN_CONNECTION_STATE_REASON_SERVICE_START_FAILED] = Tr("VPN connection service failed to start.")
	VpnErrorTable[NM_VPN_CONNECTION_STATE_REASON_NO_SECRETS] = Tr("Necessary password for VPN connection was not provided.")
	VpnErrorTable[NM_VPN_CONNECTION_STATE_REASON_LOGIN_FAILED] = Tr("Authentication to VPN server failed.")
	VpnErrorTable[NM_VPN_CONNECTION_STATE_REASON_CONNECTION_REMOVED] = Tr("The connection was deleted from settings.")
}

type StateNotifier struct {
	devices map[dbus.ObjectPath]*deviceStateStruct
	locker  sync.Mutex
}
type deviceStateStruct struct {
	nmDev   *nm.Device
	aconnId string
}

func newStateNotifier() (n *StateNotifier) {
	n = &StateNotifier{}
	n.devices = make(map[dbus.ObjectPath]*deviceStateStruct)

	var notify *notifications.Notifier
	notify, _ = notifications.NewNotifier("org.freedesktop.Notifications", "/org/freedesktop/Notifications")

	var watch = func(path dbus.ObjectPath) {
		defer func() {
			if err := recover(); err != nil {
				n.locker.Lock()
				defer n.locker.Unlock()
				delete(n.devices, path)
				logger.Error(err)
			}
		}()
		if dev, err := nmNewDevice(path); err == nil {
			n.locker.Lock()
			defer n.locker.Unlock()
			n.devices[path] = &deviceStateStruct{nmDev: dev}
			if data, err := nmGetDeviceActiveConnectionData(path); err == nil {
				// remember active connection id if exists
				n.devices[path].aconnId = getSettingConnectionId(data)
			}
			// connect signals
			dev.ConnectStateChanged(func(newState, oldState, reason uint32) {
				switch newState {
				case NM_DEVICE_STATE_PREPARE:
					if data, err := nmGetDeviceActiveConnectionData(path); err == nil {
						n.devices[path].aconnId = getSettingConnectionId(data)
					}
				case NM_DEVICE_STATE_ACTIVATED:
					if data, err := nmGetDeviceActiveConnectionData(path); err == nil {
						var icon, msg string
						switch getCustomConnectionType(data) {
						case connectionWired:
							icon = "notification-network-ethernet-connected"
						case connectionWireless:
							// TODO
							// icon = "notification-network-wireless-connected"
							icon = "notification-network-wireless-full"
						default:
							icon = "network-transmit-receive"
						}
						msg = n.devices[path].aconnId
						notify.Notify("Network", 0, icon, Tr("Connected"), msg, nil, nil, 0)
					}
				case NM_DEVICE_STATE_FAILED, NM_DEVICE_STATE_DISCONNECTED,
					NM_DEVICE_STATE_UNMANAGED, NM_DEVICE_STATE_UNAVAILABLE:
					switch oldState {
					case NM_DEVICE_STATE_FAILED, NM_DEVICE_STATE_DISCONNECTED,
						NM_DEVICE_STATE_UNMANAGED, NM_DEVICE_STATE_UNAVAILABLE:
						// this was not a disconnect
					default:
						//TODO: icon name can be different for different device type
						if reason != NM_DEVICE_STATE_REASON_NONE && reason != NM_DEVICE_STATE_REASON_UNKNOWN {
							var icon, msg string
							switch dev.DeviceType.Get() {
							case NM_DEVICE_TYPE_ETHERNET:
								icon = "notification-network-ethernet-disconnected"
							case NM_DEVICE_TYPE_WIFI:
								icon = "notification-network-wireless-disconnected"
							default:
								icon = "network-error"
							}
							if len(n.devices[path].aconnId) > 0 {
								msg = n.devices[path].aconnId
							} else {
								msg = DeviceErrorTable[reason]
							}
							notify.Notify("Network", 0, icon, Tr("Disconnect"), msg, nil, nil, 0)
						}
					}
				}
			})
		}
	}
	var remove = func(path dbus.ObjectPath) {
		n.locker.Lock()
		defer n.locker.Unlock()
		if dev, ok := n.devices[path]; ok {
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
			notify.Notify("Network", 0, "network-offline", Tr("Offline"), Tr("Disconnected - You are now offline."), nil, nil, 0)
		}
	})

	//TODO: VPN state

	return
}

func destroyStateNotifier(n *StateNotifier) {
	n.locker.Lock()
	defer n.locker.Unlock()
	for _, dev := range n.devices {
		nmDestroyDevice(dev.nmDev)
	}
	n.devices = nil
}

func notifyApModeNotSupport() {
	icon := "notification-network-wireless-disconnected"
	notify(icon, Tr("Disconnect"), Tr("Access Point (AP) mode is not supported by this device."))
}
