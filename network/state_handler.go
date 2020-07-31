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
	"sync"

	nmdbus "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.networkmanager"

	"pkg.deepin.io/dde/daemon/network/nm"
	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	. "pkg.deepin.io/lib/gettext"
)

var vpnErrorTable = make(map[uint32]string)
var deviceErrorTable = make(map[uint32]string)

func initNmStateReasons() {
	// device error table
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_NONE] = Tr("Device state changed")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_UNKNOWN] = Tr("Device state changed, reason unknown")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_NOW_MANAGED] = Tr("The device is now managed")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_NOW_UNMANAGED] = Tr("The device is no longer managed")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_CONFIG_FAILED] = Tr("The device has not been ready for configuration")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_IP_CONFIG_UNAVAILABLE] = Tr("IP configuration could not be reserved (no available address, timeout, etc)")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_IP_CONFIG_EXPIRED] = Tr("The IP configuration is no longer valid")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_NO_SECRETS] = Tr("Passwords were required but not provided")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_SUPPLICANT_DISCONNECT] = Tr("The 802.1X supplicant disconnected from the access point or authentication server")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_SUPPLICANT_CONFIG_FAILED] = Tr("Configuration of the 802.1X supplicant failed")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_SUPPLICANT_FAILED] = Tr("The 802.1X supplicant quitted or failed unexpectedly")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_SUPPLICANT_TIMEOUT] = Tr("The 802.1X supplicant took too long time to authenticate")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_PPP_START_FAILED] = Tr("The PPP service failed to start within the allowed time")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_PPP_DISCONNECT] = Tr("The PPP service disconnected unexpectedly")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_PPP_FAILED] = Tr("The PPP service quitted or failed unexpectedly")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_DHCP_START_FAILED] = Tr("The DHCP service failed to start within the allowed time")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_DHCP_ERROR] = Tr("The DHCP service reported an unexpected error")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_DHCP_FAILED] = Tr("The DHCP service quitted or failed unexpectedly")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_SHARED_START_FAILED] = Tr("The shared connection service failed to start")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_SHARED_FAILED] = Tr("The shared connection service quitted or failed unexpectedly")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_AUTOIP_START_FAILED] = Tr("The AutoIP service failed to start")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_AUTOIP_ERROR] = Tr("The AutoIP service reported an unexpected error")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_AUTOIP_FAILED] = Tr("The AutoIP service quitted or failed unexpectedly")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_MODEM_BUSY] = Tr("Dialing failed due to busy lines")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_MODEM_NO_DIAL_TONE] = Tr("Dialing failed due to no dial tone")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_MODEM_NO_CARRIER] = Tr("Dialing failed due to the carrier")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_MODEM_DIAL_TIMEOUT] = Tr("Dialing timed out")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_MODEM_DIAL_FAILED] = Tr("Dialing failed")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_MODEM_INIT_FAILED] = Tr("Modem initialization failed")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_GSM_APN_FAILED] = Tr("Failed to select the specified GSM APN")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_GSM_REGISTRATION_NOT_SEARCHING] = Tr("No networks searched")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_GSM_REGISTRATION_DENIED] = Tr("Network registration was denied")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_GSM_REGISTRATION_TIMEOUT] = Tr("Network registration timed out")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_GSM_REGISTRATION_FAILED] = Tr("Failed to register to the requested GSM network")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_GSM_PIN_CHECK_FAILED] = Tr("PIN check failed")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_FIRMWARE_MISSING] = Tr("Necessary firmware for the device may be missed")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_REMOVED] = Tr("The device was removed")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_SLEEPING] = Tr("NetworkManager went to sleep")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_CONNECTION_REMOVED] = Tr("The device's active connection was removed or disappeared")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_USER_REQUESTED] = Tr("A user or client requested to disconnect")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_CARRIER] = Tr("The device's carrier/link changed") // TODO translate
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_CONNECTION_ASSUMED] = Tr("The device's existing connection was assumed")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_SUPPLICANT_AVAILABLE] = Tr("The 802.1x supplicant is now available") // TODO translate: full stop
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_MODEM_NOT_FOUND] = Tr("The modem could not be found")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_BT_FAILED] = Tr("The Bluetooth connection timed out or failed")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_GSM_SIM_NOT_INSERTED] = Tr("GSM Modem's SIM Card was not inserted")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_GSM_SIM_PIN_REQUIRED] = Tr("GSM Modem's SIM PIN required")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_GSM_SIM_PUK_REQUIRED] = Tr("GSM Modem's SIM PUK required")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_GSM_SIM_WRONG] = Tr("SIM card error in GSM Modem")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_INFINIBAND_MODE] = Tr("InfiniBand device does not support connected mode")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_DEPENDENCY_FAILED] = Tr("A dependency of the connection failed")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_BR2684_FAILED] = Tr("RFC 2684 Ethernet bridging error to ADSL") // TODO translate
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_MODEM_MANAGER_UNAVAILABLE] = Tr("ModemManager did not run or quitted unexpectedly")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_SSID_NOT_FOUND] = Tr("The 802.11 WLAN network could not be found")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_SECONDARY_CONNECTION_FAILED] = Tr("A secondary connection of the base connection failed")

	// works for nm 1.0+
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_DCB_FCOE_FAILED] = Tr("DCB or FCoE setup failed")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_TEAMD_CONTROL_FAILED] = Tr("Network teaming control failed")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_MODEM_FAILED] = Tr("Modem failed to run or not available")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_MODEM_AVAILABLE] = Tr("Modem now ready and available")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_SIM_PIN_INCORRECT] = Tr("SIM PIN is incorrect")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_NEW_ACTIVATION] = Tr("New connection activation is enqueuing")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_PARENT_CHANGED] = Tr("Parent device changed")
	deviceErrorTable[nm.NM_DEVICE_STATE_REASON_PARENT_MANAGED_CHANGED] = Tr("Management status of parent device changed")

	// device error table for custom state reasons
	deviceErrorTable[CUSTOM_NM_DEVICE_STATE_REASON_CABLE_UNPLUGGED] = Tr("Network cable is unplugged")
	deviceErrorTable[CUSTOM_NM_DEVICE_STATE_REASON_MODEM_NO_SIGNAL] = Tr("Please make sure SIM card has been inserted with mobile network signal")
	deviceErrorTable[CUSTOM_NM_DEVICE_STATE_REASON_MODEM_WRONG_PLAN] = Tr("Please make sure a correct plan was selected without arrearage of SIM card")

	// vpn error table
	vpnErrorTable[nm.NM_VPN_CONNECTION_STATE_REASON_UNKNOWN] = Tr("Failed to activate VPN connection, reason unknown")
	vpnErrorTable[nm.NM_VPN_CONNECTION_STATE_REASON_NONE] = Tr("Failed to activate VPN connection")
	vpnErrorTable[nm.NM_VPN_CONNECTION_STATE_REASON_USER_DISCONNECTED] = Tr("The VPN connection state changed due to being disconnected by users")
	vpnErrorTable[nm.NM_VPN_CONNECTION_STATE_REASON_DEVICE_DISCONNECTED] = Tr("The VPN connection state changed due to being disconnected from devices")
	vpnErrorTable[nm.NM_VPN_CONNECTION_STATE_REASON_SERVICE_STOPPED] = Tr("VPN service stopped")
	vpnErrorTable[nm.NM_VPN_CONNECTION_STATE_REASON_IP_CONFIG_INVALID] = Tr("The IP config of VPN connection was invalid")
	vpnErrorTable[nm.NM_VPN_CONNECTION_STATE_REASON_CONNECT_TIMEOUT] = Tr("The connection attempt to VPN service timed out")
	vpnErrorTable[nm.NM_VPN_CONNECTION_STATE_REASON_SERVICE_START_TIMEOUT] = Tr("The VPN service start timed out")
	vpnErrorTable[nm.NM_VPN_CONNECTION_STATE_REASON_SERVICE_START_FAILED] = Tr("The VPN service failed to start")
	vpnErrorTable[nm.NM_VPN_CONNECTION_STATE_REASON_NO_SECRETS] = Tr("The VPN connection password was not provided")
	vpnErrorTable[nm.NM_VPN_CONNECTION_STATE_REASON_LOGIN_FAILED] = Tr("Authentication to VPN server failed")
	vpnErrorTable[nm.NM_VPN_CONNECTION_STATE_REASON_CONNECTION_REMOVED] = Tr("The connection was deleted from settings")
}

type stateHandler struct {
	m       *Manager
	devices map[dbus.ObjectPath]*deviceStateInfo
	locker  sync.Mutex

	sysSigLoop *dbusutil.SignalLoop
}

type deviceStateInfo struct {
	nmDev          *nmdbus.Device
	enabled        bool
	devUdi         string
	devType        uint32
	aconnId        string
	connectionType string
}

func newStateHandler(sysSigLoop *dbusutil.SignalLoop, m *Manager) (sh *stateHandler) {
	sh = &stateHandler{
		m:          m,
		sysSigLoop: sysSigLoop,
		devices:    make(map[dbus.ObjectPath]*deviceStateInfo),
	}

	_, err := nmManager.ConnectDeviceRemoved(func(path dbus.ObjectPath) {
		sh.remove(path)
	})
	if err != nil {
		logger.Warning(err)
	}
	_, err = nmManager.ConnectDeviceAdded(func(path dbus.ObjectPath) {
		sh.watch(path)
	})
	if err != nil {
		logger.Warning(err)
	}
	for _, path := range nmGetDevices() {
		sh.watch(path)
	}

	err = nmManager.NetworkingEnabled().ConnectChanged(func(hasValue bool, value bool) {
		if !nmGetNetworkEnabled() {
			notifyAirplanModeEnabled()
		}
	})
	if err != nil {
		logger.Warning(err)
	}
	_ = nmManager.WirelessHardwareEnabled().ConnectChanged(func(hasValue bool, value bool) {
		if !nmGetWirelessHardwareEnabled() {
			notifyWirelessHardSwitchOff()
		}
	})

	return
}

func destroyStateHandler(sh *stateHandler) {
	for path := range sh.devices {
		sh.remove(path)
	}
	sh.devices = nil
}

func (sh *stateHandler) watch(path dbus.ObjectPath) {
	defer func() {
		if err := recover(); err != nil {
			sh.locker.Lock()
			defer sh.locker.Unlock()
			delete(sh.devices, path)
			logger.Error(err)
		}
	}()

	nmDev, err := nmNewDevice(path)
	if err != nil {
		return
	}

	deviceType, _ := nmDev.DeviceType().Get(0)
	if !isDeviceTypeValid(deviceType) {
		return
	}

	sh.locker.Lock()
	defer sh.locker.Unlock()
	sh.devices[path] = &deviceStateInfo{nmDev: nmDev}
	sh.devices[path].devType = deviceType
	sh.devices[path].devUdi, _ = nmDev.Udi().Get(0)
	enabled, err := sh.m.sysNetwork.IsDeviceEnabled(0, string(path))
	if err == nil {
		sh.devices[path].enabled = enabled
	} else {
		logger.Warning(err)
	}

	if data, err := nmGetDeviceActiveConnectionData(path); err == nil {
		// remember active connection id and type if exists
		sh.devices[path].aconnId = getSettingConnectionId(data)
		sh.devices[path].connectionType = getCustomConnectionType(data)
	}

	// connect signals
	nmDev.InitSignalExt(sh.sysSigLoop, true)
	_, err = nmDev.ConnectStateChanged(func(newState, oldState, reason uint32) {
		logger.Debugf("device state changed, %d => %d, reason[%d] %s", oldState, newState, reason, deviceErrorTable[reason])
		sh.locker.Lock()
		defer sh.locker.Unlock()
		if data, err := nmGetDeviceActiveConnectionData(path); err == nil {
			// update active connection id and type if exists
			sh.devices[path].aconnId = getSettingConnectionId(data)
			sh.devices[path].connectionType = getCustomConnectionType(data)
		}
		dsi, ok := sh.devices[path]
		if !ok {
			// the device already been removed
			return
		}

		switch newState {
		case nm.NM_DEVICE_STATE_PREPARE:
			if data, err := nmGetDeviceActiveConnectionData(path); err == nil {
				dsi.aconnId = getSettingConnectionId(data)
				icon := generalGetNotifyDisconnectedIcon(dsi.devType, path)
				logger.Debug("--------[Prepare] Active connection info:", dsi.aconnId, dsi.connectionType, dsi.nmDev.Path_())
				if dsi.connectionType == connectionWirelessHotspot {
					notify(icon, "", Tr("Enabling hotspot"))
				} else {
					notify(icon, "", fmt.Sprintf(Tr("Connecting %q"), dsi.aconnId))
				}
			}
		case nm.NM_DEVICE_STATE_ACTIVATED:
			icon := generalGetNotifyConnectedIcon(dsi.devType, path)
			msg := dsi.aconnId
			logger.Debug("--------[Activated] Active connection info:", dsi.aconnId, dsi.connectionType, dsi.nmDev.Path_())
			if dsi.connectionType == connectionWirelessHotspot {
				notify(icon, "", Tr("Hotspot enabled"))
			} else {
				notify(icon, "", fmt.Sprintf(Tr("%q connected"), msg))
				if !sh.m.hasSaveSecret {
					if data, err := nmGetDeviceActiveConnectionData(path); err == nil {
						sh.savePasswordByConnectionStatus(data)
					}
					sh.m.hasSaveSecret = true
				}
			}
			go sh.m.doPortalAuthentication()
		case nm.NM_DEVICE_STATE_FAILED, nm.NM_DEVICE_STATE_DISCONNECTED, nm.NM_DEVICE_STATE_NEED_AUTH,
			nm.NM_DEVICE_STATE_UNMANAGED, nm.NM_DEVICE_STATE_UNAVAILABLE:
			logger.Infof("device disconnected, type %s, %d => %d, reason[%d] %s", getCustomDeviceType(dsi.devType), oldState, newState, reason, deviceErrorTable[reason])

			// ignore device removed signals for that could not
			// query related information correct
			if reason == nm.NM_DEVICE_STATE_REASON_REMOVED {
				if dsi.connectionType == connectionWirelessHotspot {
					icon := generalGetNotifyDisconnectedIcon(dsi.devType, path)
					notify(icon, "", Tr("Hotspot disabled"))
				}
				return
			}

			// ignore if device's old state is not available
			if !isDeviceStateAvailable(oldState) {
				logger.Debug("no notify, old state is not available")
				return
			}

			// notify only when network enabled
			if !nmGetNetworkEnabled() {
				logger.Debug("no notify, network disabled")
				return
			}

			// notify only when device enabled
			if oldState == nm.NM_DEVICE_STATE_DISCONNECTED && !dsi.enabled {
				logger.Debug("no notify, notify only when device enabled")
				return
			}

			// fix reasons
			switch dsi.devType {
			case nm.NM_DEVICE_TYPE_ETHERNET:
				if reason == nm.NM_DEVICE_STATE_REASON_CARRIER {
					reason = CUSTOM_NM_DEVICE_STATE_REASON_CABLE_UNPLUGGED
				}
			case nm.NM_DEVICE_TYPE_MODEM:
				if isDeviceStateReasonInvalid(reason) {
					// mobile device is specially, fix its reasons here
					signalQuality, _ := mmGetModemDeviceSignalQuality(dbus.ObjectPath(dsi.devUdi))
					if signalQuality == 0 {
						reason = CUSTOM_NM_DEVICE_STATE_REASON_MODEM_NO_SIGNAL
					} else {
						reason = CUSTOM_NM_DEVICE_STATE_REASON_MODEM_WRONG_PLAN
					}
				}
			}

			// ignore invalid reasons
			if isDeviceStateReasonInvalid(reason) {
				logger.Debug("no notify, device state reason invalid")
				return
			}

			logger.Debug("--------[Disconnect] Active connection info:", dsi.aconnId, dsi.connectionType, dsi.nmDev.Path_())
			var icon, msg string
			icon = generalGetNotifyDisconnectedIcon(dsi.devType, path)
			if len(msg) == 0 {
				switch reason {
				case nm.NM_DEVICE_STATE_REASON_USER_REQUESTED:
					if newState == nm.NM_DEVICE_STATE_DISCONNECTED {
						if dsi.connectionType == connectionWirelessHotspot {
							notify(icon, "", Tr("Hotspot disabled"))
						} else {
							msg = fmt.Sprintf(Tr("%q disconnected"), dsi.aconnId)
						}
					}
				case nm.NM_DEVICE_STATE_REASON_NEW_ACTIVATION:
				case nm.NM_DEVICE_STATE_REASON_IP_CONFIG_UNAVAILABLE:
					if dsi.connectionType == connectionWirelessHotspot {
						msg = Tr("Unable to share hotspot, please check dnsmasq settings")
					} else if dsi.connectionType == connectionWireless {
						msg = fmt.Sprintf(Tr("Unable to connect %q, please keep closer to the wireless router"), dsi.aconnId)
					} else if dsi.connectionType == connectionWired {
						msg = fmt.Sprintf(Tr("Unable to connect %q, please check your router or net cable."), dsi.aconnId)
					}
				case nm.NM_DEVICE_STATE_REASON_NO_SECRETS:
					msg = fmt.Sprintf(Tr("Password is required to connect %q"), dsi.aconnId)
				case nm.NM_DEVICE_STATE_REASON_SUPPLICANT_DISCONNECT:
					if oldState == nm.NM_DEVICE_STATE_CONFIG && newState == nm.NM_DEVICE_STATE_NEED_AUTH {
						msg = fmt.Sprintf(Tr("Connection failed, unable to connect %q, wrong password"), dsi.aconnId)
					}
					sh.m.hasSaveSecret = true
				case CUSTOM_NM_DEVICE_STATE_REASON_CABLE_UNPLUGGED: //disconnected due to cable unplugged
					// if device is ethernet,notify disconnected message

					logger.Debug("Disconnected due to unplugged cable")
					if dsi.devType == nm.NM_DEVICE_TYPE_ETHERNET {
						logger.Debug("unplugged device is ethernet")
						msg = fmt.Sprintf(Tr("%q disconnected"), dsi.aconnId)
					}

					//default:
					//	if dsi.aconnId != "" {
					//		msg = fmt.Sprintf(Tr("%q disconnected"), dsi.aconnId)
					//	}
				}
			}
			if msg != "" {
				notify(icon, "", msg)
			}
		}
	})
	if err != nil {
		logger.Warning(err)
	}
}

func (sh *stateHandler) remove(path dbus.ObjectPath) {
	sh.locker.Lock()
	defer sh.locker.Unlock()
	if dev, ok := sh.devices[path]; ok {
		nmDestroyDevice(dev.nmDev)
		delete(sh.devices, path)
	}
}

//Because the password is saved to keyring by
//default, the storage operation needs to be
//performed according to the status judgment, and
//only when the password is correct.
func (sh *stateHandler) savePasswordByConnectionStatus(data connectionData) {
	connUUID, ok := getConnectionDataString(data, "connection", "uuid")
	if !ok {
		logger.Debug("Failed to save password because can not find connUUID")
		return
	}
	for _, item := range sh.m.items {
		err := sh.m.secretAgent.set(item.label, connUUID, item.settingName, item.settingKey, item.value)
		if err != nil {
			logger.Debug("failed to save secret when status connected")
			return
		}
	}
}
