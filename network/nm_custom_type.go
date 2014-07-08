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
	. "pkg.linuxdeepin.com/lib/gettext"
	"strconv"
)

// Custom device types, use sting instead of number, used by front-end
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
)

func getCustomDeviceType(devType uint32) (customDevType string) {
	switch devType {
	case NM_DEVICE_TYPE_ETHERNET:
		return deviceEthernet
	case NM_DEVICE_TYPE_WIFI:
		return deviceWifi
	case NM_DEVICE_TYPE_UNUSED1:
		return deviceUnused1
	case NM_DEVICE_TYPE_UNUSED2:
		return deviceUnused2
	case NM_DEVICE_TYPE_BT:
		return deviceBt
	case NM_DEVICE_TYPE_OLPC_MESH:
		return deviceOlpcMesh
	case NM_DEVICE_TYPE_WIMAX:
		return deviceWimax
	case NM_DEVICE_TYPE_MODEM:
		return deviceModem
	case NM_DEVICE_TYPE_INFINIBAND:
		return deviceInfiniband
	case NM_DEVICE_TYPE_BOND:
		return deviceBond
	case NM_DEVICE_TYPE_VLAN:
		return deviceVlan
	case NM_DEVICE_TYPE_ADSL:
		return deviceAdsl
	case NM_DEVICE_TYPE_BRIDGE:
		return deviceBridge
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
	connectionVpnPptp         = "vpn-pptp"
	connectionVpnVpnc         = "vpn-vpnc"
)

// wrapper for custom connection types
const (
	connectionMobile = "mobile" // wrapper for gsm and cdma
	connectionVpn    = "vpn"    // wrapper for all vpn types
)

// key-map values for internationalization
type connectionType struct {
	Value, Text string
}

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
	connectionVpnVpnc,
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
