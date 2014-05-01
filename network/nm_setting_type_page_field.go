package main

import (
	"dlib"
)

const (
	typeUnknown         = "unknown"
	typeWired           = "wired"
	typeWireless        = "wireless"
	typeWirelessAdhoc   = "wireless-adhoc"
	typeWirelessHotspot = "wireless-hotspot"
	typePppoe           = "pppoe"
	typeVpn             = "vpn"
	typeVpnL2tp         = "vpn-l2tp"
	typeVpnOpenconnect  = "vpn-openconnect"
	typeVpnOpenvpn      = "vpn-openvpn"
	typeVpnPptp         = "vpn-pptp"
	typeVpnVpnc         = "vpn-vpnc"
)

// TODO key-map values for internationalization
type connectionType struct {
	Value, Text string
}

var supportedConnectionTypes = []string{
	// typeWired,// don't support multiple wired connections since now
	typeWireless,
	typePppoe,
	typeVpnL2tp,
	typeVpnOpenconnect,
	typeVpnOpenvpn,
	typeVpnPptp,
	typeVpnVpnc,
}
var supportedConnectionTypesInfo = []connectionType{
	// connectionType{typeWired, dlib.Tr("Ethernet")},// don't support multiple wired connections since now
	connectionType{typeWireless, dlib.Tr("Wi-Fi")},
	connectionType{typePppoe, dlib.Tr("PPPoE")},
	connectionType{typeVpnL2tp, dlib.Tr("VPN-L2TP (Layer 2 Tunneling Protocol)")},
	connectionType{typeVpnOpenconnect, dlib.Tr("VPN-OpenConnect (Cisco AnyConnect Compatible VPN)")},
	connectionType{typeVpnOpenvpn, dlib.Tr("VPN-OpenVPN")},
	connectionType{typeVpnPptp, dlib.Tr("VPN-PPTP (Point-to-Point Tunneling Protocol))")},
	connectionType{typeVpnVpnc, dlib.Tr("VPN-VPNC (Cisco Compatible VPN)")},
}

const (
	field8021x              = NM_SETTING_802_1X_SETTING_NAME
	fieldConnection         = NM_SETTING_CONNECTION_SETTING_NAME
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
