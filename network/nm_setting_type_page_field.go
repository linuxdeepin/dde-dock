package main

import (
	"dlib"
	"strconv"
)

// device type
const (
	deviceTypeUnknown    = "unknown"
	deviceTypeEthernet   = "wired"
	deviceTypeWifi       = "wireless"
	deviceTypeUnused1    = "unused1"
	deviceTypeUnused2    = "unused2"
	deviceTypeBt         = "bt"
	deviceTypeOlpcMesh   = "olpc-mesh"
	deviceTypeWimax      = "wimax"
	deviceTypeModem      = "modem"
	deviceTypeInfiniband = "infiniband"
	deviceTypeBond       = "bond"
	deviceTypeVlan       = "vlan"
	deviceTypeAdsl       = "adsl"
	deviceTypeBridge     = "bridge"
)

func getDeviceTypeName(devType uint32) (devName string) {
	switch devType {
	case NM_DEVICE_TYPE_ETHERNET:
		return deviceTypeEthernet
	case NM_DEVICE_TYPE_WIFI:
		return deviceTypeWifi
	case NM_DEVICE_TYPE_UNUSED1:
		return deviceTypeUnused1
	case NM_DEVICE_TYPE_UNUSED2:
		return deviceTypeUnused2
	case NM_DEVICE_TYPE_BT:
		return deviceTypeBt
	case NM_DEVICE_TYPE_OLPC_MESH:
		return deviceTypeOlpcMesh
	case NM_DEVICE_TYPE_WIMAX:
		return deviceTypeWimax
	case NM_DEVICE_TYPE_MODEM:
		return deviceTypeModem
	case NM_DEVICE_TYPE_INFINIBAND:
		return deviceTypeInfiniband
	case NM_DEVICE_TYPE_BOND:
		return deviceTypeBond
	case NM_DEVICE_TYPE_VLAN:
		return deviceTypeVlan
	case NM_DEVICE_TYPE_ADSL:
		return deviceTypeAdsl
	case NM_DEVICE_TYPE_BRIDGE:
		return deviceTypeBridge
	default:
		logger.Error("unknown device type", devType)
	}
	return deviceTypeUnknown
}

// connection type
const (
	typeUnknown         = "unknown"
	typeWired           = "wired"
	typeWireless        = "wireless"
	typeWirelessAdhoc   = "wireless-adhoc"
	typeWirelessHotspot = "wireless-hotspot"
	typePppoe           = "pppoe"
	typeMobile          = "mobile"
	typeVpn             = "vpn"
	typeVpnL2tp         = "vpn-l2tp"
	typeVpnOpenconnect  = "vpn-openconnect"
	typeVpnOpenvpn      = "vpn-openvpn"
	typeVpnPptp         = "vpn-pptp"
	typeVpnVpnc         = "vpn-vpnc"
)

// key-map values for internationalization
type connectionType struct {
	Value, Text string
}

var supportedConnectionTypes = []string{
	// typeWired,// don't support multiple wired connections since now
	typeWireless,
	typeWirelessAdhoc,
	typeWirelessHotspot,
	typePppoe,
	typeMobile,
	typeVpnL2tp,
	typeVpnOpenconnect,
	typeVpnOpenvpn,
	typeVpnPptp,
	typeVpnVpnc,
}
var supportedConnectionTypesInfo = []connectionType{
	// connectionType{typeWired, dlib.Tr("Ethernet")},// don't support multiple wired connections since now
	connectionType{typeWireless, dlib.Tr("Wi-Fi")},
	connectionType{typeWirelessAdhoc, dlib.Tr("Wi-Fi Ad-Hoc")},
	connectionType{typeWirelessHotspot, dlib.Tr("Wi-Fi Hotspot")},
	connectionType{typePppoe, dlib.Tr("PPPoE")},
	connectionType{typeMobile, dlib.Tr("Mobile 2G/3G/4G-LTE")},
	connectionType{typeVpnL2tp, dlib.Tr("VPN-L2TP (Layer 2 Tunneling Protocol)")},
	connectionType{typeVpnOpenconnect, dlib.Tr("VPN-OpenConnect (Cisco AnyConnect Compatible VPN)")},
	connectionType{typeVpnOpenvpn, dlib.Tr("VPN-OpenVPN")},
	connectionType{typeVpnPptp, dlib.Tr("VPN-PPTP (Point-to-Point Tunneling Protocol))")},
	connectionType{typeVpnVpnc, dlib.Tr("VPN-VPNC (Cisco Compatible VPN)")},
}

const (
	field8021x              = NM_SETTING_802_1X_SETTING_NAME
	fieldConnection         = NM_SETTING_CONNECTION_SETTING_NAME
	fieldGsm                = NM_SETTING_GSM_SETTING_NAME
	fieldIpv4               = NM_SETTING_IP4_CONFIG_SETTING_NAME
	fieldIpv6               = NM_SETTING_IP6_CONFIG_SETTING_NAME
	fieldPppoe              = NM_SETTING_PPPOE_SETTING_NAME
	fieldPpp                = NM_SETTING_PPP_SETTING_NAME
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
	pageGsm                = "gsm"                  // -> fieldGsm
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
	case typeWired:
		idPrefix = dlib.Tr("Wired Connection")
	case typeWireless:
		idPrefix = dlib.Tr("Wireless Connection")
	case typeWirelessAdhoc:
		idPrefix = dlib.Tr("Wireless Ad-Hoc")
	case typeWirelessHotspot:
		idPrefix = dlib.Tr("Wireless Ap-Hotspot")
	case typePppoe:
		idPrefix = dlib.Tr("PPPoE Connection")
	case typeMobile:
		idPrefix = dlib.Tr("Mobile Connection")
	case typeVpn:
		idPrefix = dlib.Tr("VPN Connection")
	case typeVpnL2tp:
		idPrefix = dlib.Tr("VPN L2TP")
	case typeVpnOpenconnect:
		idPrefix = dlib.Tr("VPN OpenConnect")
	case typeVpnOpenvpn:
		idPrefix = dlib.Tr("VPN OpenVPN")
	case typeVpnPptp:
		idPrefix = dlib.Tr("VPN PPTP")
	case typeVpnVpnc:
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
