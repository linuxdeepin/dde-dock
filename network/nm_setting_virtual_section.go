package network

// Sections, correspondence to "NM_SETTING_XXX" in network manager.
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
	sectionVpnL2tp            = NM_SETTING_ALIAS_VPN_L2TP_SETTING_NAME
	sectionVpnL2tpPpp         = NM_SETTING_ALIAS_VPN_L2TP_PPP_SETTING_NAME
	sectionVpnL2tpIpsec       = NM_SETTING_ALIAS_VPN_L2TP_IPSEC_SETTING_NAME
	sectionVpnOpenconnect     = NM_SETTING_ALIAS_VPN_OPENCONNECT_SETTING_NAME
	sectionVpnOpenvpn         = NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME
	sectionVpnOpenvpnAdvanced = NM_SETTING_ALIAS_VPN_OPENVPN_ADVANCED_SETTING_NAME
	sectionVpnOpenvpnSecurity = NM_SETTING_ALIAS_VPN_OPENVPN_SECURITY_SETTING_NAME
	sectionVpnOpenvpnTlsauth  = NM_SETTING_ALIAS_VPN_OPENVPN_TLSAUTH_SETTING_NAME
	sectionVpnOpenvpnProxies  = NM_SETTING_ALIAS_VPN_OPENVPN_PROXIES_SETTING_NAME
	sectionVpnPptp            = NM_SETTING_ALIAS_VPN_PPTP_SETTING_NAME
	sectionVpnPptpPpp         = NM_SETTING_ALIAS_VPN_PPTP_PPP_SETTING_NAME
	sectionVpnVpnc            = NM_SETTING_ALIAS_VPN_VPNC_SETTING_NAME
	sectionVpnVpncAdvanced    = NM_SETTING_ALIAS_VPN_VPNC_ADVANCED_SETTING_NAME
	sectionWired              = NM_SETTING_WIRED_SETTING_NAME
	sectionWireless           = NM_SETTING_WIRELESS_SETTING_NAME
	sectionWirelessSecurity   = NM_SETTING_WIRELESS_SECURITY_SETTING_NAME
)

// Alias sections, used for vpn connection which is a special key in fact
const (
	NM_SETTING_ALIAS_VPN_L2TP_SETTING_NAME             = "alias-vpn-l2tp"
	NM_SETTING_ALIAS_VPN_L2TP_PPP_SETTING_NAME         = "alias-vpn-l2tp-ppp"
	NM_SETTING_ALIAS_VPN_L2TP_IPSEC_SETTING_NAME       = "alias-vpn-l2tp-ipsec"
	NM_SETTING_ALIAS_VPN_OPENCONNECT_SETTING_NAME      = "alias-vpn-openconnect"
	NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME          = "alias-vpn-openvpn"
	NM_SETTING_ALIAS_VPN_OPENVPN_ADVANCED_SETTING_NAME = "alias-vpn-openvpn-advanced"
	NM_SETTING_ALIAS_VPN_OPENVPN_SECURITY_SETTING_NAME = "alias-vpn-openvpn-security"
	NM_SETTING_ALIAS_VPN_OPENVPN_TLSAUTH_SETTING_NAME  = "alias-vpn-openvpn-tlsauth"
	NM_SETTING_ALIAS_VPN_OPENVPN_PROXIES_SETTING_NAME  = "alias-vpn-openvpn-proxies"
	NM_SETTING_ALIAS_VPN_PPTP_SETTING_NAME             = "alias-vpn-pptp"
	NM_SETTING_ALIAS_VPN_PPTP_PPP_SETTING_NAME         = "alias-vpn-pptp-ppp"
	NM_SETTING_ALIAS_VPN_VPNC_SETTING_NAME             = "alias-vpn-vpnc"
	NM_SETTING_ALIAS_VPN_VPNC_ADVANCED_SETTING_NAME    = "alias-vpn-advanced"
)

func getRealSectionName(name string) (realName string) {
	realName = name
	switch name {
	case NM_SETTING_ALIAS_VPN_L2TP_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_ALIAS_VPN_L2TP_PPP_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_ALIAS_VPN_L2TP_IPSEC_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_ALIAS_VPN_OPENCONNECT_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_ALIAS_VPN_OPENVPN_ADVANCED_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_ALIAS_VPN_OPENVPN_SECURITY_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_ALIAS_VPN_OPENVPN_TLSAUTH_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_ALIAS_VPN_OPENVPN_PROXIES_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_ALIAS_VPN_PPTP_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_ALIAS_VPN_PPTP_PPP_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_ALIAS_VPN_VPNC_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_ALIAS_VPN_VPNC_ADVANCED_SETTING_NAME:
		realName = sectionVpn
	}
	return
}

// Virtual sections for front-end to easy to configure
const (
	vsectionGeneral  = "vs-general"  // -> sectionConnection
	vsectionEthernet = "vs-ethernet" // -> sectionWired
	vsectionMobile   = "vs-mobile"   // -> sectionGsm, sectionCdma
	// TODO
	vsectionMobileGsm  = "vs-mobile-gsm"  // -> sectionGsm
	vsectionMobileCdma = "vs-mobile-cdma" // -> sectionCdma
	vsectionWifi       = "vs-wifi"        // -> sectionWireless
	vsectionIpv4       = "vs-ipv4"        // -> sectionIpv4
	vsectionIpv6       = "vs-ipv6"        // -> sectionIpv6
	vsectionSecurity   = "vs-security"    // -> section8021x, sectionWirelessSecurity
	vsectionPppoe      = "vs-pppoe"       // -> sectionPppoe
	vsectionPpp        = "vs-ppp"         // -> sectionPpp
	vsectionVpn        = "vs-vpn"         // -> sectionVpnL2tp, sectionVpnOpenconnect, sectionVpnOpenvpn, sectionVpnPptp, sectionVpnVpnc
	// TODO
	vsectionVpnL2tp            = "vs-vpn-l2tp"             // -> sectionVpnL2tp
	vsectionVpnL2tpPpp         = "vs-vpn-l2tp-ppp"         // -> sectionVpnL2tpPpp
	vsectionVpnL2tpIpsec       = "vs-vpn-l2tp-ipsec"       // -> sectionVpnL2tpIpsec
	vsectionVpnOpenconnect     = "vs-vpn-openconnect"      // -> sectionVpnOpenconnect
	vsectionVpnOpenvpn         = "vs-vpn-openvpn"          // -> sectionVpnOpenvpn
	vsectionVpnOpenvpnAdvanced = "vs-vpn-openvpn-advanced" // -> sectionVpnOpenVpnAdvanced
	vsectionVpnOpenvpnSecurity = "vs-vpn-openvpn-security" // -> sectionVpnOpenVpnSecurity
	vsectionVpnOpenvpnTlsauth  = "vs-vpn-openvpn-tlsauth"  // -> sectionVpnOpenVpnTlsauth
	vsectionVpnOpenvpnProxies  = "vs-vpn-openvpn-proxies"  // -> sectionVpnOpenVpnProxies
	vsectionVpnPptp            = "vs-vpn-pptp"             // -> sectionVpnPptp
	vsectionVpnPptpPpp         = "vs-vpn-pptp-ppp"         // -> sectionVpnPptpPpp
	vsectionVpnVpnc            = "vs-vpn-vpnc"             // -> sectionVpnVpnc
	vsectionVpnVpncAdvanced    = "vs-vpn-vpnc-advanced"    // -> sectionVpnVpncAdvanced
)

// get available virtual sections for target connection type
func getAvailableVsections(data connectionData) (vsections []string) {
	connectionType := getCustomConnectionType(data)
	switch connectionType {
	case connectionWired:
		vsections = []string{
			vsectionGeneral,
			vsectionEthernet,
			vsectionIpv4,
			vsectionIpv6,
			vsectionSecurity,
		}
	case connectionWireless:
		vsections = []string{
			vsectionGeneral,
			vsectionWifi,
			vsectionIpv4,
			vsectionIpv6,
			vsectionSecurity,
		}
	case connectionWirelessAdhoc:
		vsections = []string{
			vsectionGeneral,
			vsectionWifi,
			vsectionIpv4,
			vsectionIpv6,
			vsectionSecurity,
		}
	case connectionWirelessHotspot:
		vsections = []string{
			vsectionGeneral,
			vsectionWifi,
			vsectionIpv4,
			vsectionIpv6,
			vsectionSecurity,
		}
	case connectionPppoe:
		vsections = []string{
			vsectionGeneral,
			vsectionEthernet,
			vsectionPppoe,
			vsectionPpp,
			vsectionIpv4,
		}
	case connectionVpnL2tp:
		vsections = []string{
			vsectionGeneral,
			vsectionVpn,
			vsectionVpnL2tpPpp,
			vsectionVpnL2tpIpsec,
			vsectionIpv4,
		}
	case connectionVpnOpenconnect:
		vsections = []string{
			vsectionGeneral,
			vsectionVpn,
			vsectionIpv4,
			vsectionIpv6,
		}
	case connectionVpnOpenvpn:
		vsections = []string{
			vsectionGeneral,
			vsectionVpn,
			vsectionVpnOpenvpnAdvanced,
			vsectionVpnOpenvpnSecurity,
			vsectionVpnOpenvpnProxies,
			vsectionIpv4,
			vsectionIpv6,
		}
		// when connection connection is static key, vsectionVpnOpenvpnTlsauth is not available
		if getSettingVpnOpenvpnKeyConnectionType(data) != NM_OPENVPN_CONTYPE_STATIC_KEY {
			vsections = append(vsections, vsectionVpnOpenvpnTlsauth)
		}
	case connectionVpnPptp:
		vsections = []string{
			vsectionGeneral,
			vsectionVpn,
			vsectionVpnPptpPpp,
			vsectionIpv4,
		}
	case connectionVpnVpnc:
		vsections = []string{
			vsectionGeneral,
			vsectionVpn,
			vsectionVpnVpncAdvanced,
			vsectionIpv4,
		}
	case connectionMobileGsm:
		vsections = []string{
			vsectionGeneral,
			vsectionMobile,
			vsectionPpp,
			vsectionIpv4,
		}
	case connectionMobileCdma:
		vsections = []string{
			vsectionGeneral,
			vsectionMobile,
			vsectionPpp,
			vsectionIpv4,
		}
	}
	return
}

func getRelatedSectionsOfVsection(data connectionData, vsection string) (sections []string) {
	connectionType := getCustomConnectionType(data)
	switch vsection {
	default:
		logger.Error("getRelatedSectionsOfVsection: invalid vsection name", vsection)
	case vsectionGeneral:
		sections = []string{sectionConnection}
	case vsectionMobile:
		switch connectionType {
		case connectionMobileGsm:
			sections = []string{sectionGsm}
		case connectionMobileCdma:
			sections = []string{sectionCdma}
		}
	case vsectionMobileGsm: // TODO: remove
		sections = []string{sectionGsm}
	case vsectionMobileCdma:
		sections = []string{sectionCdma}
	case vsectionEthernet:
		sections = []string{sectionWired}
	case vsectionWifi:
		sections = []string{sectionWireless}
	case vsectionIpv4:
		sections = []string{sectionIpv4}
	case vsectionIpv6:
		sections = []string{sectionIpv6}
	case vsectionSecurity:
		switch connectionType {
		case connectionWired:
			sections = []string{section8021x}
		case connectionWireless, connectionWirelessAdhoc, connectionWirelessHotspot:
			if isSettingSectionExists(data, section8021x) {
				sections = []string{sectionWirelessSecurity, section8021x}
			} else {
				sections = []string{sectionWirelessSecurity}
			}
		}
	case vsectionPppoe:
		sections = []string{sectionPppoe}
	case vsectionPpp:
		sections = []string{sectionPpp}
	case vsectionVpn:
		switch connectionType {
		case connectionVpnL2tp:
			sections = []string{sectionVpnL2tp}
		case connectionVpnOpenconnect:
			sections = []string{sectionVpnOpenconnect}
		case connectionVpnOpenvpn:
			sections = []string{sectionVpnOpenvpn}
		case connectionVpnPptp:
			sections = []string{sectionVpnPptp}
		case connectionVpnVpnc:
			sections = []string{sectionVpnVpnc}
		}
	case vsectionVpnL2tp: // TODO
		sections = []string{sectionVpnL2tp}
	case vsectionVpnL2tpPpp:
		sections = []string{sectionVpnL2tpPpp}
	case vsectionVpnL2tpIpsec:
		sections = []string{sectionVpnL2tpIpsec}
	case vsectionVpnOpenconnect: // TODO
		sections = []string{sectionVpnOpenconnect}
	case vsectionVpnOpenvpn: // TODO
		sections = []string{sectionVpnOpenvpn}
	case vsectionVpnOpenvpnAdvanced:
		sections = []string{sectionVpnOpenvpnAdvanced}
	case vsectionVpnOpenvpnSecurity:
		sections = []string{sectionVpnOpenvpnSecurity}
	case vsectionVpnOpenvpnTlsauth:
		sections = []string{sectionVpnOpenvpnTlsauth}
	case vsectionVpnOpenvpnProxies:
		sections = []string{sectionVpnOpenvpnProxies}
	case vsectionVpnPptp: // TODO
		sections = []string{sectionVpnPptp}
	case vsectionVpnPptpPpp:
		sections = []string{sectionVpnPptpPpp}
	case vsectionVpnVpnc: // TODO
		sections = []string{sectionVpnVpnc}
	case vsectionVpnVpncAdvanced:
		sections = []string{sectionVpnVpncAdvanced}
	}
	return
}

// getAvailableSections return all virtual section related real sections
func getAvailableSections(data connectionData) (sections []string) {
	for _, vsection := range getAvailableVsections(data) {
		sections = appendStrArrayUnique(sections, getRelatedSectionsOfVsection(data, vsection)...)
	}
	return
}

// get real section name of target key in virtual section
func getSectionOfKeyInVsection(data connectionData, vsection, key string) (section string) {
	sections := getRelatedSectionsOfVsection(data, vsection)
	for _, section := range sections {
		if generalIsKeyInSettingSection(section, key) {
			return section
		}
	}
	logger.Errorf("get corresponding section of key in virtual section failed, vsection=%s, key=%s", vsection, key)
	return ""
}

// get available keys of virtual section
func getAvailableKeysOfVsection(data connectionData, vsection string) (keys []string) {
	sections := getRelatedSectionsOfVsection(data, vsection)
	for _, section := range sections {
		keys = appendStrArrayUnique(keys, generalGetSettingAvailableKeys(data, section)...)
	}
	if len(keys) == 0 {
		logger.Warning("there is no available keys for virtual section", vsection)
	}
	return
}
