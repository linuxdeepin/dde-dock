package main

import (
	"dlib"
	"strconv"
)

// custom device types
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

func getDeviceTypeName(devType uint32) (devName string) {
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

// custom connection types
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

const (
	field8021x              = NM_SETTING_802_1X_SETTING_NAME
	fieldConnection         = NM_SETTING_CONNECTION_SETTING_NAME
	fieldGsm                = NM_SETTING_GSM_SETTING_NAME
	fieldCdma               = NM_SETTING_CDMA_SETTING_NAME
	fieldIpv4               = NM_SETTING_IP4_CONFIG_SETTING_NAME
	fieldIpv6               = NM_SETTING_IP6_CONFIG_SETTING_NAME
	fieldPppoe              = NM_SETTING_PPPOE_SETTING_NAME
	fieldPpp                = NM_SETTING_PPP_SETTING_NAME
	fieldSerial             = NM_SETTING_SERIAL_SETTING_NAME
	fieldVpn                = NM_SETTING_VPN_SETTING_NAME
	fieldVpnL2tp            = NM_SETTING_VF_VPN_L2TP_SETTING_NAME
	fieldVpnL2tpPpp         = NM_SETTING_VF_VPN_L2TP_PPP_SETTING_NAME
	fieldVpnL2tpIpsec       = NM_SETTING_VF_VPN_L2TP_IPSEC_SETTING_NAME
	fieldVpnOpenconnect     = NM_SETTING_VF_VPN_OPENCONNECT_SETTING_NAME
	fieldVpnOpenvpn         = NM_SETTING_VF_VPN_OPENVPN_SETTING_NAME
	fieldVpnOpenvpnAdvanced = NM_SETTING_VF_VPN_OPENVPN_ADVANCED_SETTING_NAME
	fieldVpnOpenvpnSecurity = NM_SETTING_VF_VPN_OPENVPN_SECURITY_SETTING_NAME
	fieldVpnOpenvpnTlsauth  = NM_SETTING_VF_VPN_OPENVPN_TLSAUTH_SETTING_NAME
	fieldVpnOpenvpnProxies  = NM_SETTING_VF_VPN_OPENVPN_PROXIES_SETTING_NAME
	fieldVpnPptp            = NM_SETTING_VF_VPN_PPTP_SETTING_NAME
	fieldVpnPptpPpp         = NM_SETTING_VF_VPN_PPTP_PPP_SETTING_NAME
	fieldVpnVpnc            = NM_SETTING_VF_VPN_VPNC_SETTING_NAME
	fieldVpnVpncAdvanced    = NM_SETTING_VF_VPN_VPNC_ADVANCED_SETTING_NAME
	fieldWired              = NM_SETTING_WIRED_SETTING_NAME
	fieldWireless           = NM_SETTING_WIRELESS_SETTING_NAME
	fieldWirelessSecurity   = NM_SETTING_WIRELESS_SECURITY_SETTING_NAME
)

// page is a wrapper of fields for easy to configure
const (
	pageGeneral            = "general"              // -> fieldConnection
	pageEthernet           = "ethernet"             // -> fieldWireed
	pageMobile             = "mobile"               // -> fieldGsm
	pageMobileCdma         = "mobile-cdma"          // -> fieldCdma
	pageWifi               = "wifi"                 // -> fieldWireless
	pageIPv4               = "ipv4"                 // -> fieldIpv4
	pageIPv6               = "ipv6"                 // -> fieldIpv6
	pageSecurity           = "security"             // -> field8021x, fieldWirelessSecurity
	pagePppoe              = "pppoe"                // -> fieldPppoe
	pagePpp                = "ppp"                  // -> fieldPpp
	pageVpnL2tp            = "vpn-l2tp"             // -> fieldVpnL2tp
	pageVpnL2tpPpp         = "vpn-l2tp-ppp"         // -> fieldVpnL2tpPpp
	pageVpnL2tpIpsec       = "vpn-l2tp-ipsec"       // -> fieldVpnL2tpIpsec
	pageVpnOpenconnect     = "vpn-openconnect"      // -> fieldVpnOpenconnect
	pageVpnOpenvpn         = "vpn-openvpn"          // -> fieldVpnOpenvpn
	pageVpnOpenvpnAdvanced = "vpn-openvpn-advanced" // -> fieldVpnOpenVpnAdvanced
	pageVpnOpenvpnSecurity = "vpn-openvpn-security" // -> fieldVpnOpenVpnSecurity
	pageVpnOpenvpnTlsauth  = "vpn-openvpn-tlsauth"  // -> fieldVpnOpenVpnTlsauth
	pageVpnOpenvpnProxies  = "vpn-openvpn-proxies"  // -> fieldVpnOpenVpnProxies
	pageVpnPptp            = "vpn-pptp"             // -> fieldVpnPptp
	pageVpnPptpPpp         = "vpn-pptp-ppp"         // -> fieldVpnPptpPpp
	pageVpnVpnc            = "vpn-vpnc"             // -> fieldVpnVpnc
	pageVpnVpncAdvanced    = "vpn-vpnc-advanced"    // -> fieldVpnVpncAdvanced
)

// Virtual fields, used for vpn connectionns.
const (
	NM_SETTING_VF_VPN_L2TP_SETTING_NAME             = "vf-vpn-l2tp"
	NM_SETTING_VF_VPN_L2TP_PPP_SETTING_NAME         = "vf-vpn-l2tp-ppp"
	NM_SETTING_VF_VPN_L2TP_IPSEC_SETTING_NAME       = "vf-vpn-l2tp-ipsec"
	NM_SETTING_VF_VPN_OPENCONNECT_SETTING_NAME      = "vf-vpn-openconnect"
	NM_SETTING_VF_VPN_OPENVPN_SETTING_NAME          = "vf-vpn-openvpn"
	NM_SETTING_VF_VPN_OPENVPN_ADVANCED_SETTING_NAME = "vf-vpn-openvpn-advanced"
	NM_SETTING_VF_VPN_OPENVPN_SECURITY_SETTING_NAME = "vf-vpn-openvpn-security"
	NM_SETTING_VF_VPN_OPENVPN_TLSAUTH_SETTING_NAME  = "vf-vpn-openvpn-tlsauth"
	NM_SETTING_VF_VPN_OPENVPN_PROXIES_SETTING_NAME  = "vf-vpn-openvpn-proxies"
	NM_SETTING_VF_VPN_PPTP_SETTING_NAME             = "vf-vpn-pptp"
	NM_SETTING_VF_VPN_PPTP_PPP_SETTING_NAME         = "vf-vpn-pptp-ppp"
	NM_SETTING_VF_VPN_VPNC_SETTING_NAME             = "vf-vpn-vpnc"
	NM_SETTING_VF_VPN_VPNC_ADVANCED_SETTING_NAME    = "vf-vpn-advanced"
)

func getRealFieldName(name string) (realName string) {
	realName = name
	switch name {
	case NM_SETTING_VF_VPN_L2TP_SETTING_NAME:
		realName = fieldVpn
	case NM_SETTING_VF_VPN_L2TP_PPP_SETTING_NAME:
		realName = fieldVpn
	case NM_SETTING_VF_VPN_L2TP_IPSEC_SETTING_NAME:
		realName = fieldVpn
	case NM_SETTING_VF_VPN_OPENCONNECT_SETTING_NAME:
		realName = fieldVpn
	case NM_SETTING_VF_VPN_OPENVPN_SETTING_NAME:
		realName = fieldVpn
	case NM_SETTING_VF_VPN_OPENVPN_ADVANCED_SETTING_NAME:
		realName = fieldVpn
	case NM_SETTING_VF_VPN_OPENVPN_SECURITY_SETTING_NAME:
		realName = fieldVpn
	case NM_SETTING_VF_VPN_OPENVPN_TLSAUTH_SETTING_NAME:
		realName = fieldVpn
	case NM_SETTING_VF_VPN_OPENVPN_PROXIES_SETTING_NAME:
		realName = fieldVpn
	case NM_SETTING_VF_VPN_PPTP_SETTING_NAME:
		realName = fieldVpn
	case NM_SETTING_VF_VPN_PPTP_PPP_SETTING_NAME:
		realName = fieldVpn
	case NM_SETTING_VF_VPN_VPNC_SETTING_NAME:
		realName = fieldVpn
	case NM_SETTING_VF_VPN_VPNC_ADVANCED_SETTING_NAME:
		realName = fieldVpn
	}
	return
}

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
