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
	"strconv"

	"pkg.deepin.io/dde/daemon/network/nm"
	. "pkg.deepin.io/lib/gettext"
)

// Custom device types, use string instead of number, used by front-end
const (
	deviceUnknown    = "unknown"
	deviceEthernet   = "wired"
	deviceWifi       = "wireless"
	deviceUnused1    = "unused1"
	deviceUnused2    = "unused2"
	deviceBt         = "bt"
	deviceOlpcMesh   = "olpc-mesh"
	deviceWimax      = "wimax"
	deviceModem      = "modem"
	deviceInfiniband = "infiniband"
	deviceBond       = "bond"
	deviceVlan       = "vlan"
	deviceAdsl       = "adsl"
	deviceBridge     = "bridge"
	deviceGeneric    = "generic"
	deviceTeam       = "team"
	deviceTun        = "tun"
)

func getCustomDeviceType(devType uint32) (customDevType string) {
	switch devType {
	case nm.NM_DEVICE_TYPE_ETHERNET:
		return deviceEthernet
	case nm.NM_DEVICE_TYPE_WIFI:
		return deviceWifi
	case nm.NM_DEVICE_TYPE_UNUSED1:
		return deviceUnused1
	case nm.NM_DEVICE_TYPE_UNUSED2:
		return deviceUnused2
	case nm.NM_DEVICE_TYPE_BT:
		return deviceBt
	case nm.NM_DEVICE_TYPE_OLPC_MESH:
		return deviceOlpcMesh
	case nm.NM_DEVICE_TYPE_WIMAX:
		return deviceWimax
	case nm.NM_DEVICE_TYPE_MODEM:
		return deviceModem
	case nm.NM_DEVICE_TYPE_INFINIBAND:
		return deviceInfiniband
	case nm.NM_DEVICE_TYPE_BOND:
		return deviceBond
	case nm.NM_DEVICE_TYPE_VLAN:
		return deviceVlan
	case nm.NM_DEVICE_TYPE_ADSL:
		return deviceAdsl
	case nm.NM_DEVICE_TYPE_BRIDGE:
		return deviceBridge
	case nm.NM_DEVICE_TYPE_GENERIC:
		return deviceGeneric
	case nm.NM_DEVICE_TYPE_TEAM:
		return deviceTeam
	case nm.NM_DEVICE_TYPE_TUN:
		return deviceTun
	case nm.NM_DEVICE_TYPE_UNKNOWN:
	default:
		logger.Error("unknown device type", devType)
	}
	return deviceUnknown
}

// Custom connection types
const (
	connectionUnknown         = "unknown"
	connectionWired           = "wired"
	connectionWireless        = "wireless"
	connectionWirelessAdhoc   = "wireless-adhoc"
	connectionWirelessHotspot = "wireless-hotspot"
	connectionPppoe           = "pppoe"
	connectionMobileGsm       = "mobile-gsm"
	connectionMobileCdma      = "mobile-cdma"
	connectionVpnL2tp         = "vpn-l2tp"
	connectionVpnOpenconnect  = "vpn-openconnect"
	connectionVpnOpenvpn      = "vpn-openvpn"
	connectionVpnStrongswan   = "vpn-strongswan"
	connectionVpnPptp         = "vpn-pptp"
	connectionVpnVpnc         = "vpn-vpnc"
)

// wrapper for custom connection types
const (
	connectionMobile = "mobile" // wrapper for gsm and cdma
	connectionVpn    = "vpn"    // wrapper for all vpn types
)

var supportedConnectionTypes = []string{
	connectionWired,
	connectionWireless,
	connectionWirelessAdhoc,
	connectionWirelessHotspot,
	connectionPppoe,
	connectionMobile,
	connectionMobileGsm,
	connectionMobileCdma,
	connectionVpn,
	connectionVpnL2tp,
	connectionVpnOpenconnect,
	connectionVpnOpenvpn,
	connectionVpnPptp,
	connectionVpnStrongswan,
	connectionVpnVpnc,
}

func getCustomConnectionTypeForUuid(uuid string) (connType string) {
	connType = connectionUnknown
	cpath, err := nmGetConnectionByUuid(uuid)
	if err != nil {
		return
	}
	data, err := nmGetConnectionData(cpath)
	if err != nil {
		return
	}
	return getCustomConnectionType(data)
}

// return custom connection type, and the wrapper types will be ignored, e.g. connectionMobile.
func getCustomConnectionType(data connectionData) (connType string) {
	t := getSettingConnectionType(data)
	switch t {
	case nm.NM_SETTING_WIRED_SETTING_NAME:
		connType = connectionWired
	case nm.NM_SETTING_WIRELESS_SETTING_NAME:
		if isSettingWirelessModeExists(data) {
			switch getSettingWirelessMode(data) {
			case nm.NM_SETTING_WIRELESS_MODE_INFRA:
				connType = connectionWireless
			case nm.NM_SETTING_WIRELESS_MODE_ADHOC:
				connType = connectionWirelessAdhoc
			case nm.NM_SETTING_WIRELESS_MODE_AP:
				connType = connectionWirelessHotspot
			}
		} else {
			connType = connectionWireless
		}
	case nm.NM_SETTING_PPPOE_SETTING_NAME:
		connType = connectionPppoe
	case nm.NM_SETTING_GSM_SETTING_NAME:
		connType = connectionMobileGsm
	case nm.NM_SETTING_CDMA_SETTING_NAME:
		connType = connectionMobileCdma
	case nm.NM_SETTING_VPN_SETTING_NAME:
		switch getSettingVpnServiceType(data) {
		case nm.NM_DBUS_SERVICE_L2TP:
			connType = connectionVpnL2tp
		case nm.NM_DBUS_SERVICE_OPENCONNECT:
			connType = connectionVpnOpenconnect
		case nm.NM_DBUS_SERVICE_OPENVPN:
			connType = connectionVpnOpenvpn
		case nm.NM_DBUS_SERVICE_PPTP:
			connType = connectionVpnPptp
		case nm.NM_DBUS_SERVICE_STRONGSWAN:
			connType = connectionVpnStrongswan
		case nm.NM_DBUS_SERVICE_VPNC:
			connType = connectionVpnVpnc
		}
	}
	if len(connType) == 0 {
		connType = connectionUnknown
	}
	return
}

func isWirelessConnection(data connectionData) (isWireless bool) {
	if getSettingConnectionType(data) == nm.NM_SETTING_WIRELESS_SETTING_NAME {
		return true
	}
	return false
}

func isVpnConnection(data connectionData) (isVpn bool) {
	if getSettingConnectionType(data) == nm.NM_SETTING_VPN_SETTING_NAME {
		return true
	}
	return false
}

func isCreatedManuallyConnection(data connectionData) (isCreateManual bool) {
	if isVpnConnection(data) {
		return true
	}
	switch getCustomConnectionType(data) {
	case connectionPppoe:
		return true
	}
	return false
}

// generate connection id when creating a new connection
func genConnectionId(connType string) (id string) {
	var idPrefix string
	switch connType {
	default:
		idPrefix = Tr("Connection")
	case connectionWired:
		idPrefix = Tr("Wired Connection")
	case connectionWireless:
		idPrefix = Tr("Wireless Connection")
	case connectionWirelessAdhoc:
		idPrefix = Tr("Wireless Ad-Hoc")
	case connectionWirelessHotspot:
		idPrefix = Tr("Wireless Ap-Hotspot")
	case connectionPppoe:
		idPrefix = Tr("PPPoE Connection")
	case connectionMobile:
		idPrefix = Tr("Mobile Connection")
	case connectionMobileGsm:
		idPrefix = Tr("Mobile GSM Connection")
	case connectionMobileCdma:
		idPrefix = Tr("Mobile CDMA Connection")
	case connectionVpn:
		idPrefix = Tr("VPN Connection")
	case connectionVpnL2tp:
		idPrefix = Tr("VPN L2TP")
	case connectionVpnOpenconnect:
		idPrefix = Tr("VPN OpenConnect")
	case connectionVpnOpenvpn:
		idPrefix = Tr("VPN OpenVPN")
	case connectionVpnPptp:
		idPrefix = Tr("VPN PPTP")
	case connectionVpnStrongswan:
		idPrefix = Tr("VPN StrongSwan")
	case connectionVpnVpnc:
		idPrefix = Tr("VPN VPNC")
	}
	allIds := nmGetConnectionIds()
	for i := 1; ; i++ {
		id = idPrefix + " " + strconv.Itoa(i)
		if !isStringInArray(id, allIds) {
			break
		}
	}
	return
}

const (
	nmKeyErrorInvalidValue = "invalid value"
)
