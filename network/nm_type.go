package main

import (
	"dlib"
	"strconv"
)

// Custom device types
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

func getDeviceName(devType uint32) (devName string) {
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
	connectionMobile          = "mobile"
	connectionMobileGsm       = "mobile-gsm"
	connectionMobileCdma      = "mobile-cdma"
	connectionVpn             = "vpn"
	connectionVpnL2tp         = "vpn-l2tp"
	connectionVpnOpenconnect  = "vpn-openconnect"
	connectionVpnOpenvpn      = "vpn-openvpn"
	connectionVpnPptp         = "vpn-pptp"
	connectionVpnVpnc         = "vpn-vpnc"
)

// key-map values for internationalization
type connectionType struct {
	Value, Text string
}

var supportedConnectionTypes = []string{
	// connectionWired,// don't support multiple wired connections since now
	connectionWireless,
	connectionWirelessAdhoc,
	connectionWirelessHotspot,
	connectionPppoe,
	connectionMobileGsm,
	connectionMobileCdma,
	connectionVpnL2tp,
	connectionVpnOpenconnect,
	connectionVpnOpenvpn,
	connectionVpnPptp,
	connectionVpnVpnc,
}
var supportedConnectionTypesInfo = []connectionType{
	// connectionType{connectionWired, dlib.Tr("Ethernet")},// don't support multiple wired connections since now
	connectionType{connectionWireless, dlib.Tr("Wi-Fi")},
	connectionType{connectionWirelessAdhoc, dlib.Tr("Wi-Fi Ad-Hoc")},
	connectionType{connectionWirelessHotspot, dlib.Tr("Wi-Fi Hotspot")},
	connectionType{connectionPppoe, dlib.Tr("PPPoE")},
	connectionType{connectionMobileGsm, dlib.Tr("Mobile GSM (GPRS, EDGE, UMTS, HSPA)")},
	connectionType{connectionMobileCdma, dlib.Tr("Mobile CDMA (1xRTT, EVDO)")},
	connectionType{connectionVpnL2tp, dlib.Tr("VPN-L2TP (Layer 2 Tunneling Protocol)")},
	connectionType{connectionVpnOpenconnect, dlib.Tr("VPN-OpenConnect (Cisco AnyConnect Compatible VPN)")},
	connectionType{connectionVpnOpenvpn, dlib.Tr("VPN-OpenVPN")},
	connectionType{connectionVpnPptp, dlib.Tr("VPN-PPTP (Point-to-Point Tunneling Protocol))")},
	connectionType{connectionVpnVpnc, dlib.Tr("VPN-VPNC (Cisco Compatible VPN)")},
}

// generate connection id when creating a new connection
func genConnectionId(connType string) (id string) {
	var idPrefix string
	switch connType {
	default:
		idPrefix = dlib.Tr("Connection")
	case connectionWired:
		idPrefix = dlib.Tr("Wired Connection")
	case connectionWireless:
		idPrefix = dlib.Tr("Wireless Connection")
	case connectionWirelessAdhoc:
		idPrefix = dlib.Tr("Wireless Ad-Hoc")
	case connectionWirelessHotspot:
		idPrefix = dlib.Tr("Wireless Ap-Hotspot")
	case connectionPppoe:
		idPrefix = dlib.Tr("PPPoE Connection")
	case connectionMobileGsm:
		idPrefix = dlib.Tr("Mobile GSM Connection")
	case connectionMobileCdma:
		idPrefix = dlib.Tr("Mobile CDMA Connection")
	case connectionVpn:
		idPrefix = dlib.Tr("VPN Connection")
	case connectionVpnL2tp:
		idPrefix = dlib.Tr("VPN L2TP")
	case connectionVpnOpenconnect:
		idPrefix = dlib.Tr("VPN OpenConnect")
	case connectionVpnOpenvpn:
		idPrefix = dlib.Tr("VPN OpenVPN")
	case connectionVpnPptp:
		idPrefix = dlib.Tr("VPN PPTP")
	case connectionVpnVpnc:
		idPrefix = dlib.Tr("VPN VPNC")
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
