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
	section8021x              = NM_SETTING_802_1X_SETTING_NAME
	sectionConnection         = NM_SETTING_CONNECTION_SETTING_NAME
	sectionGsm                = NM_SETTING_GSM_SETTING_NAME
	sectionCdma               = NM_SETTING_CDMA_SETTING_NAME
	sectionIpv4               = NM_SETTING_IP4_CONFIG_SETTING_NAME
	sectionIpv6               = NM_SETTING_IP6_CONFIG_SETTING_NAME
	sectionPppoe              = NM_SETTING_PPPOE_SETTING_NAME
	sectionPpp                = NM_SETTING_PPP_SETTING_NAME
	sectionSerial             = NM_SETTING_SERIAL_SETTING_NAME
	sectionVpn                = NM_SETTING_VPN_SETTING_NAME
	sectionVpnL2tp            = NM_SETTING_VF_VPN_L2TP_SETTING_NAME
	sectionVpnL2tpPpp         = NM_SETTING_VF_VPN_L2TP_PPP_SETTING_NAME
	sectionVpnL2tpIpsec       = NM_SETTING_VF_VPN_L2TP_IPSEC_SETTING_NAME
	sectionVpnOpenconnect     = NM_SETTING_VF_VPN_OPENCONNECT_SETTING_NAME
	sectionVpnOpenvpn         = NM_SETTING_VF_VPN_OPENVPN_SETTING_NAME
	sectionVpnOpenvpnAdvanced = NM_SETTING_VF_VPN_OPENVPN_ADVANCED_SETTING_NAME
	sectionVpnOpenvpnSecurity = NM_SETTING_VF_VPN_OPENVPN_SECURITY_SETTING_NAME
	sectionVpnOpenvpnTlsauth  = NM_SETTING_VF_VPN_OPENVPN_TLSAUTH_SETTING_NAME
	sectionVpnOpenvpnProxies  = NM_SETTING_VF_VPN_OPENVPN_PROXIES_SETTING_NAME
	sectionVpnPptp            = NM_SETTING_VF_VPN_PPTP_SETTING_NAME
	sectionVpnPptpPpp         = NM_SETTING_VF_VPN_PPTP_PPP_SETTING_NAME
	sectionVpnVpnc            = NM_SETTING_VF_VPN_VPNC_SETTING_NAME
	sectionVpnVpncAdvanced    = NM_SETTING_VF_VPN_VPNC_ADVANCED_SETTING_NAME
	sectionWired              = NM_SETTING_WIRED_SETTING_NAME
	sectionWireless           = NM_SETTING_WIRELESS_SETTING_NAME
	sectionWirelessSecurity   = NM_SETTING_WIRELESS_SECURITY_SETTING_NAME
)

// page is a wrapper of sections for easy to configure
const (
	pageGeneral            = "general"              // -> sectionConnection
	pageEthernet           = "ethernet"             // -> sectionWireed
	pageMobile             = "mobile"               // -> sectionGsm
	pageMobileCdma         = "mobile-cdma"          // -> sectionCdma
	pageWifi               = "wifi"                 // -> sectionWireless
	pageIPv4               = "ipv4"                 // -> sectionIpv4
	pageIPv6               = "ipv6"                 // -> sectionIpv6
	pageSecurity           = "security"             // -> section8021x, sectionWirelessSecurity
	pagePppoe              = "pppoe"                // -> sectionPppoe
	pagePpp                = "ppp"                  // -> sectionPpp
	pageVpnL2tp            = "vpn-l2tp"             // -> sectionVpnL2tp
	pageVpnL2tpPpp         = "vpn-l2tp-ppp"         // -> sectionVpnL2tpPpp
	pageVpnL2tpIpsec       = "vpn-l2tp-ipsec"       // -> sectionVpnL2tpIpsec
	pageVpnOpenconnect     = "vpn-openconnect"      // -> sectionVpnOpenconnect
	pageVpnOpenvpn         = "vpn-openvpn"          // -> sectionVpnOpenvpn
	pageVpnOpenvpnAdvanced = "vpn-openvpn-advanced" // -> sectionVpnOpenVpnAdvanced
	pageVpnOpenvpnSecurity = "vpn-openvpn-security" // -> sectionVpnOpenVpnSecurity
	pageVpnOpenvpnTlsauth  = "vpn-openvpn-tlsauth"  // -> sectionVpnOpenVpnTlsauth
	pageVpnOpenvpnProxies  = "vpn-openvpn-proxies"  // -> sectionVpnOpenVpnProxies
	pageVpnPptp            = "vpn-pptp"             // -> sectionVpnPptp
	pageVpnPptpPpp         = "vpn-pptp-ppp"         // -> sectionVpnPptpPpp
	pageVpnVpnc            = "vpn-vpnc"             // -> sectionVpnVpnc
	pageVpnVpncAdvanced    = "vpn-vpnc-advanced"    // -> sectionVpnVpncAdvanced
)

// Virtual sections, used for vpn connectionns.
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

func getRealSectionName(name string) (realName string) {
	realName = name
	switch name {
	case NM_SETTING_VF_VPN_L2TP_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_VF_VPN_L2TP_PPP_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_VF_VPN_L2TP_IPSEC_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_VF_VPN_OPENCONNECT_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_VF_VPN_OPENVPN_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_VF_VPN_OPENVPN_ADVANCED_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_VF_VPN_OPENVPN_SECURITY_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_VF_VPN_OPENVPN_TLSAUTH_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_VF_VPN_OPENVPN_PROXIES_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_VF_VPN_PPTP_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_VF_VPN_PPTP_PPP_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_VF_VPN_VPNC_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_VF_VPN_VPNC_ADVANCED_SETTING_NAME:
		realName = sectionVpn
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
