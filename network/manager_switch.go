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
	"pkg.deepin.io/dde/daemon/network/nm"
	"pkg.deepin.io/lib/dbus1"
)

type switchHandler struct {
	config *config

	NetworkingEnabled bool // airplane mode for NetworkManager
	WirelessEnabled   bool
	WwanEnabled       bool
	WiredEnabled      bool
	VpnEnabled        bool
}

func newSwitchHandler(c *config) (sh *switchHandler) {
	sh = &switchHandler{config: c}
	sh.init()

	// connect global switch signals
	nmManager.NetworkingEnabled().ConnectChanged(func(hasValue bool, value bool) {
		sh.setPropNetworkingEnabled()
	})
	nmManager.WirelessEnabled().ConnectChanged(func(hasValue bool, value bool) {
		sh.setPropWirelessEnabled()
	})
	nmManager.WwanEnabled().ConnectChanged(func(hasValue bool, value bool) {
		// FIXME: when mobile adapter plugin, dbus property
		// WwanEnabled will be set to false automatically, don't known
		// why, so we force it to true here
		wwanEnabled, _ := nmManager.WwanEnabled().Get(0)
		if wwanEnabled == false {
			nmManager.WwanEnabled().Set(0, true)
			return
		}
		sh.setPropWwanEnabled()
	})

	return
}

func (sh *switchHandler) init() {
	// initialize global switches
	sh.initPropNetworkingEnabled()
	sh.initPropWirelessEnabled()
	sh.initPropWwanEnabled()

	// initialize virtual global switches
	sh.initPropWiredEnabled()
	sh.initPropVpnEnabled()
}

func (sh *switchHandler) setNetworkingEnabled(enabled bool) {
	nmSetNetworkingEnabled(enabled)
}

func (sh *switchHandler) setWirelessEnabled(enabled bool) {
	if sh.NetworkingEnabled {
		nmSetWirelessEnabled(enabled)
	} else {
		// if NetworkingEnabled is off, turn it on, and only keep
		// current global device switch alive
		sh.config.setLastGlobalSwithes(false)
		sh.config.setLastWirelessEnabled(true)
		sh.setNetworkingEnabled(true)
	}
}
func (sh *switchHandler) setWwanEnabled(enabled bool) {
	if sh.NetworkingEnabled {
		nmSetWwanEnabled(enabled)
	} else {
		// if NetworkingEnabled is off, turn it on, and only keep
		// current global device switch alive
		sh.config.setLastGlobalSwithes(false)
		sh.config.setLastWwanEnabled(true)
		sh.setNetworkingEnabled(true)
	}
}
func (sh *switchHandler) setWiredEnabled(enabled bool) {
	if sh.NetworkingEnabled {
		sh.setPropWiredEnabled(enabled)
	} else {
		// if NetworkingEnabled is off, turn it on, and only keep
		// current global device switch alive
		sh.config.setLastGlobalSwithes(false)
		sh.config.setLastWiredEnabled(true)
		sh.setNetworkingEnabled(true)
	}
}
func (sh *switchHandler) setVpnEnabled(enabled bool) {
	logger.Debug("switchHandler.setVpnEnabled", enabled)
	if sh.NetworkingEnabled {
		sh.setPropVpnEnabled(enabled)
	} else {
		// if NetworkingEnabled is off, turn it on, and only keep
		// current global device switch alive
		sh.config.setLastGlobalSwithes(false)
		sh.config.setLastVpnEnabled(true)
		sh.setNetworkingEnabled(true)
	}
}

func (sh *switchHandler) initPropNetworkingEnabled() {
	if nmHasSystemSettingsModifyPermission() {
		sh.NetworkingEnabled, _ = nmManager.NetworkingEnabled().Get(0)
		if !sh.NetworkingEnabled {
			sh.doTurnOffGlobalDeviceSwitches()
		}
	}
	sh.NetworkingEnabled, _ = nmManager.NetworkingEnabled().Get(0)
	manager.setPropNetworkingEnabled(sh.NetworkingEnabled)
}
func (sh *switchHandler) setPropNetworkingEnabled() {
	networkingEnabled, _ := nmManager.NetworkingEnabled().Get(0)
	if sh.NetworkingEnabled == networkingEnabled {
		return
	}
	sh.NetworkingEnabled = networkingEnabled
	// setup global device switches
	if sh.NetworkingEnabled {
		sh.restoreGlobalDeviceSwitches()
	} else {
		sh.saveAndTurnOffGlobalDeviceSwitches()
	}
	sh.NetworkingEnabled, _ = nmManager.NetworkingEnabled().Get(0)
	manager.setPropNetworkingEnabled(sh.NetworkingEnabled)
}
func (sh *switchHandler) restoreGlobalDeviceSwitches() {
	nmSetWirelessEnabled(sh.config.getLastWirelessEnabled())
	nmSetWwanEnabled(sh.config.getLastWwanEnabled())
	sh.setPropWiredEnabled(sh.config.getLastWiredEnabled())
	sh.setPropVpnEnabled(sh.config.getLastVpnEnabled())
}
func (sh *switchHandler) saveAndTurnOffGlobalDeviceSwitches() {
	sh.config.setLastWirelessEnabled(sh.WirelessEnabled)
	sh.config.setLastWwanEnabled(sh.WwanEnabled)
	sh.config.setLastWiredEnabled(sh.WiredEnabled)
	sh.config.setLastVpnEnabled(sh.VpnEnabled)
	sh.doTurnOffGlobalDeviceSwitches()
}
func (sh *switchHandler) doTurnOffGlobalDeviceSwitches() {
	nmSetWirelessEnabled(false)
	nmSetWwanEnabled(false)
	sh.setPropWiredEnabled(false)
	sh.setPropVpnEnabled(false)
}

func (sh *switchHandler) initPropWirelessEnabled() {
	if nmHasSystemSettingsModifyPermission() {
		sh.WirelessEnabled, _ = nmManager.WirelessEnabled().Get(0)
		for _, devPath := range nmGetDevicesByType(nm.NM_DEVICE_TYPE_WIFI) {
			if sh.WirelessEnabled {
				sh.doEnableDevice(devPath, sh.config.getDeviceEnabled(devPath))
			} else {
				sh.doEnableDevice(devPath, false)
			}
		}
	}
	sh.WirelessEnabled, _ = nmManager.WirelessEnabled().Get(0)
	manager.setPropWirelessEnabled(sh.WirelessEnabled)
}
func (sh *switchHandler) setPropWirelessEnabled() {
	wirelessEnabled, _ := nmManager.WirelessEnabled().Get(0)
	if sh.WirelessEnabled == wirelessEnabled {
		return
	}
	sh.WirelessEnabled = wirelessEnabled
	logger.Debug("setPropWirelessEnabled", sh.WirelessEnabled)
	// setup wireless devices switches
	for _, devPath := range nmGetDevicesByType(nm.NM_DEVICE_TYPE_WIFI) {
		if sh.WirelessEnabled {
			sh.restoreDeviceState(devPath)
		} else {
			sh.saveAndDisconnectDevice(devPath)
		}
	}
	sh.WirelessEnabled, _ = nmManager.WirelessEnabled().Get(0)
	manager.setPropWirelessEnabled(sh.WirelessEnabled)
}

func (sh *switchHandler) initPropWwanEnabled() {
	if nmHasSystemSettingsModifyPermission() {
		sh.WwanEnabled, _ = nmManager.WwanEnabled().Get(0)
		for _, devPath := range nmGetDevicesByType(nm.NM_DEVICE_TYPE_MODEM) {
			if sh.WwanEnabled {
				sh.doEnableDevice(devPath, sh.config.getDeviceEnabled(devPath))
			} else {
				sh.doEnableDevice(devPath, false)
			}
		}
	}
	sh.WwanEnabled, _ = nmManager.WwanEnabled().Get(0)
	manager.setPropWwanEnabled(sh.WwanEnabled)
}

func (sh *switchHandler) setPropWwanEnabled() {
	wwanEnabled, _ := nmManager.WwanEnabled().Get(0)
	if sh.WwanEnabled == wwanEnabled {
		return
	}
	sh.WwanEnabled = wwanEnabled
	// setup modem devices switches
	for _, devPath := range nmGetDevicesByType(nm.NM_DEVICE_TYPE_MODEM) {
		if sh.WwanEnabled {
			sh.restoreDeviceState(devPath)
		} else {
			sh.saveAndDisconnectDevice(devPath)
		}
	}
	sh.WwanEnabled, _ = nmManager.WwanEnabled().Get(0)
	manager.setPropWwanEnabled(sh.WwanEnabled)
}

func (sh *switchHandler) initPropWiredEnabled() {
	sh.WiredEnabled = sh.config.getWiredEnabled()
	if nmHasSystemSettingsModifyPermission() {
		for _, devPath := range nmGetDevicesByType(nm.NM_DEVICE_TYPE_ETHERNET) {
			if sh.WiredEnabled {
				sh.doEnableDevice(devPath, sh.config.getDeviceEnabled(devPath))
			} else {
				sh.doEnableDevice(devPath, false)
			}
		}
	}
	manager.setPropWiredEnabled(sh.WiredEnabled)
}
func (sh *switchHandler) setPropWiredEnabled(enabled bool) {
	if sh.config.WiredEnabled == enabled {
		return
	}
	logger.Debug("setPropWiredEnabled", enabled)
	sh.WiredEnabled = enabled
	sh.config.setWiredEnabled(enabled)
	// setup wired devices switches
	for _, devPath := range nmGetDevicesByType(nm.NM_DEVICE_TYPE_ETHERNET) {
		if enabled {
			sh.restoreDeviceState(devPath)
		} else {
			sh.saveAndDisconnectDevice(devPath)
		}
	}
	manager.setPropWiredEnabled(sh.WiredEnabled)
}

func (sh *switchHandler) initPropVpnEnabled() {
	sh.VpnEnabled = sh.config.getVpnEnabled()
	if nmHasSystemSettingsModifyPermission() {
		sh.enableVpn(sh.config.getVpnEnabled())
	}
	manager.setPropVpnEnabled(sh.VpnEnabled)
}
func (sh *switchHandler) setPropVpnEnabled(enabled bool) {
	if sh.config.getVpnEnabled() == enabled {
		return
	}
	sh.enableVpn(enabled)
}
func (sh *switchHandler) enableVpn(enabled bool) {
	// setup vpn connections
	for _, uuid := range nmGetConnectionUuidsByType(nm.NM_SETTING_VPN_SETTING_NAME) {
		if enabled {
			sh.restoreVpnConnectionState(uuid)
		} else {
			sh.deactivateVpnConnection(uuid)
		}
	}

	// setup VpnEnabled state after vpn connections dispatched to
	// avoid such issue that if there are activating vpn connections
	// and user toggle off VpnEnabled manually then the code in
	// Manager.doHandleVpnNotification will fix the VpnEnabled state
	// back to true.
	sh.doEnableVpn(enabled)
}
func (sh *switchHandler) doEnableVpn(enabled bool) {
	sh.VpnEnabled = enabled
	sh.config.setVpnEnabled(enabled)
	manager.setPropVpnEnabled(enabled)
}

func (sh *switchHandler) initDeviceState(devPath dbus.ObjectPath) (err error) {
	err = sh.doEnableDevice(devPath, sh.config.getDeviceEnabled(devPath))
	return
}
func (sh *switchHandler) restoreDeviceState(devPath dbus.ObjectPath) (err error) {
	sh.config.restoreDeviceState(devPath)
	err = sh.doEnableDevice(devPath, sh.config.getDeviceEnabled(devPath))
	return
}
func (sh *switchHandler) saveAndDisconnectDevice(devPath dbus.ObjectPath) (err error) {
	sh.config.saveDeviceState(devPath)
	err = sh.doEnableDevice(devPath, false)
	return
}

func (sh *switchHandler) enableDevice(devPath dbus.ObjectPath, enabled bool) (err error) {
	if nmGetDeviceType(devPath) == nm.NM_DEVICE_TYPE_WIFI {
		if !nmGetWirelessHardwareEnabled() {
			notifyWirelessHardSwitchOff()
			return
		}
	}
	return sh.doEnableDevice(devPath, enabled)
}
func (sh *switchHandler) doEnableDevice(devPath dbus.ObjectPath, enabled bool) (err error) {
	if enabled && sh.trunOnGlobalDeviceSwitchIfNeed(devPath) {
		return
	}
	devConfig, err := sh.config.getDeviceConfigForPath(devPath)
	if err != nil {
		return
	}
	logger.Debugf("doEnableDevice %s %v %#v", devPath, enabled, devConfig)

	sh.config.setDeviceEnabled(devPath, enabled)
	if enabled {
		// try to active last connection
		uuids := nmGetConnectionUuidsForAutoConnect(devPath, devConfig.LastConnectionUuid)
		var uuidToActive string

		switch nmGetDeviceType(devPath) {
		case nm.NM_DEVICE_TYPE_WIFI:
			ssids := nmGetAccessPointSsids(devPath)
			logger.Debug("available ssids", ssids)
			for _, uuid := range uuids {
				// ignore the hotspot/adhoc connections
				switch getCustomConnectionTypeForUuid(uuid) {
				case connectionWirelessHotspot, connectionWirelessAdhoc:
					continue
				}

				// if is wireless connection, check if the access
				// point exists around, if not, ignore it
				ssid := string(nmGetWirelessConnectionSsidByUuid(uuid))
				logger.Debug("check ssid", ssid, uuid)
				if !isStringInArray(ssid, ssids) {
					continue
				} else {
					uuidToActive = uuid
					break
				}
			}
		default:
			if len(uuids) > 0 {
				uuidToActive = uuids[0]
			}
		}

		if len(uuidToActive) > 0 {
			activeUuid, _ := nmGetDeviceActiveConnectionUuid(devPath)
			if uuidToActive != activeUuid {
				manager.nmRunOnceUntilDeviceAvailable(devPath, func() {
					manager.activateConnection(uuidToActive, devPath)
				})
			}
		}
	} else {
		err = manager.doDisconnectDevice(devPath)
	}
	return
}

func (sh *switchHandler) restoreVpnConnectionState(uuid string) (err error) {
	vpnConfig, err := sh.config.getVpnConfig(uuid)
	if err != nil {
		return
	}
	if vpnConfig.lastActivated || vpnConfig.AutoConnect {
		sh.activateVpnConnection(uuid)
	} else {
		err = manager.deactivateConnection(uuid)
	}
	return
}
func (sh *switchHandler) activateVpnConnection(uuid string) {
	if _, err := nmGetActiveConnectionByUuid(uuid); err == nil {
		// connection already activated
		return
	}
	nmRunOnceUtilNetworkAvailable(func() {
		manager.activateConnection(uuid, "/")
	})
}
func (sh *switchHandler) deactivateVpnConnection(uuid string) (err error) {
	vpnConfig, err := sh.config.getVpnConfig(uuid)
	if err != nil {
		return
	}
	vpnConfig.lastActivated = vpnConfig.activated
	err = manager.deactivateConnection(uuid)
	sh.config.save()
	return
}

func (sh *switchHandler) trunOnGlobalDeviceSwitchIfNeed(devPath dbus.ObjectPath) (need bool) {
	// if global device switch is off, turn it on, and only keep
	// current device alive
	need = (sh.generalGetGlobalDeviceEnabled(devPath) == false)
	if !need {
		return
	}
	sh.config.setAllDeviceLastEnabled(false)
	sh.config.setDeviceLastEnabled(devPath, true)
	sh.generalSetGlobalDeviceEnabled(devPath, true)
	return
}

func (sh *switchHandler) generalGetGlobalDeviceEnabled(devPath dbus.ObjectPath) (enabled bool) {
	switch devType := nmGetDeviceType(devPath); devType {
	case nm.NM_DEVICE_TYPE_ETHERNET:
		enabled = sh.WiredEnabled
	case nm.NM_DEVICE_TYPE_WIFI:
		enabled = sh.WirelessEnabled
	case nm.NM_DEVICE_TYPE_MODEM:
		enabled = sh.WwanEnabled
	}
	return
}
func (sh *switchHandler) generalSetGlobalDeviceEnabled(devPath dbus.ObjectPath, enabled bool) {
	switch devType := nmGetDeviceType(devPath); devType {
	case nm.NM_DEVICE_TYPE_ETHERNET:
		sh.setWiredEnabled(enabled)
	case nm.NM_DEVICE_TYPE_WIFI:
		sh.setWirelessEnabled(enabled)
	case nm.NM_DEVICE_TYPE_MODEM:
		sh.setWwanEnabled(enabled)
	}
}
