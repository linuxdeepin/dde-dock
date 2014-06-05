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
	NM_SETTING_VS_GENERAL              = "vs-general"
	NM_SETTING_VS_ETHERNET             = "vs-ethernet"
	NM_SETTING_VS_MOBILE               = "vs-mobile"
	NM_SETTING_VS_MOBILE_GSM           = "vs-mobile-gsm"
	NM_SETTING_VS_MOBILE_CDMA          = "vs-mobile-cdma"
	NM_SETTING_VS_WIFI                 = "vs-wifi"
	NM_SETTING_VS_IPV4                 = "vs-ipv4"
	NM_SETTING_VS_IPV6                 = "vs-ipv6"
	NM_SETTING_VS_SECURITY             = "vs-security"
	NM_SETTING_VS_PPPOE                = "vs-pppoe"
	NM_SETTING_VS_PPP                  = "vs-ppp"
	NM_SETTING_VS_VPN                  = "vs-vpn"
	NM_SETTING_VS_VPN_L2TP             = "vs-vpn-l2tp"
	NM_SETTING_VS_VPN_L2TP_PPP         = "vs-vpn-l2tp-ppp"
	NM_SETTING_VS_VPN_L2TP_IPSEC       = "vs-vpn-l2tp-ipsec"
	NM_SETTING_VS_VPN_OPENCONNECT      = "vs-vpn-openconnect"
	NM_SETTING_VS_VPN_OPENVPN          = "vs-vpn-openvpn"
	NM_SETTING_VS_VPN_OPENVPN_ADVANCED = "vs-vpn-openvpn-advanced"
	NM_SETTING_VS_VPN_OPENVPN_SECURITY = "vs-vpn-openvpn-security"
	NM_SETTING_VS_VPN_OPENVPN_TLSAUTH  = "vs-vpn-openvpn-tlsauth"
	NM_SETTING_VS_VPN_OPENVPN_PROXIES  = "vs-vpn-openvpn-proxies"
	NM_SETTING_VS_VPN_PPTP             = "vs-vpn-pptp"
	NM_SETTING_VS_VPN_PPTP_PPP         = "vs-vpn-pptp-ppp"
	NM_SETTING_VS_VPN_VPNC             = "vs-vpn-vpnc"
	NM_SETTING_VS_VPN_VPNC_ADVANCED    = "vs-vpn-vpnc-advanced"
)
const (
	vsectionGeneral  = NM_SETTING_VS_GENERAL  // -> sectionConnection
	vsectionEthernet = NM_SETTING_VS_ETHERNET // -> sectionWired
	vsectionMobile   = NM_SETTING_VS_MOBILE   // -> sectionGsm, sectionCdma
	// TODO
	vsectionMobileGsm  = NM_SETTING_VS_MOBILE_GSM  // -> sectionGsm
	vsectionMobileCdma = NM_SETTING_VS_MOBILE_CDMA // -> sectionCdma
	vsectionWifi       = NM_SETTING_VS_WIFI        // -> sectionWireless
	vsectionIpv4       = NM_SETTING_VS_IPV4        // -> sectionIpv4
	vsectionIpv6       = NM_SETTING_VS_IPV6        // -> sectionIpv6
	vsectionSecurity   = NM_SETTING_VS_SECURITY    // -> section8021x, sectionWirelessSecurity
	vsectionPppoe      = NM_SETTING_VS_PPPOE       // -> sectionPppoe
	vsectionPpp        = NM_SETTING_VS_PPP         // -> sectionPpp
	vsectionVpn        = NM_SETTING_VS_VPN         // -> sectionVpnL2tp, sectionVpnOpenconnect, sectionVpnOpenvpn, sectionVpnPptp, sectionVpnVpnc
	// TODO
	vsectionVpnL2tp            = NM_SETTING_VS_VPN_L2TP             // -> sectionVpnL2tp
	vsectionVpnL2tpPpp         = NM_SETTING_VS_VPN_L2TP_PPP         // -> sectionVpnL2tpPpp
	vsectionVpnL2tpIpsec       = NM_SETTING_VS_VPN_L2TP_IPSEC       // -> sectionVpnL2tpIpsec
	vsectionVpnOpenconnect     = NM_SETTING_VS_VPN_OPENCONNECT      // -> sectionVpnOpenconnect
	vsectionVpnOpenvpn         = NM_SETTING_VS_VPN_OPENVPN          // -> sectionVpnOpenvpn
	vsectionVpnOpenvpnAdvanced = NM_SETTING_VS_VPN_OPENVPN_ADVANCED // -> sectionVpnOpenVpnAdvanced
	vsectionVpnOpenvpnSecurity = NM_SETTING_VS_VPN_OPENVPN_SECURITY // -> sectionVpnOpenVpnSecurity
	vsectionVpnOpenvpnTlsauth  = NM_SETTING_VS_VPN_OPENVPN_TLSAUTH  // -> sectionVpnOpenVpnTlsauth
	vsectionVpnOpenvpnProxies  = NM_SETTING_VS_VPN_OPENVPN_PROXIES  // -> sectionVpnOpenVpnProxies
	vsectionVpnPptp            = NM_SETTING_VS_VPN_PPTP             // -> sectionVpnPptp
	vsectionVpnPptpPpp         = NM_SETTING_VS_VPN_PPTP_PPP         // -> sectionVpnPptpPpp
	vsectionVpnVpnc            = NM_SETTING_VS_VPN_VPNC             // -> sectionVpnVpnc
	vsectionVpnVpncAdvanced    = NM_SETTING_VS_VPN_VPNC_ADVANCED    // -> sectionVpnVpncAdvanced
)

func isVirtualSection(section string) bool {
	switch section {
	case vsectionGeneral:
		return true
	case vsectionEthernet:
		return true
	case vsectionMobile:
		return true
	case vsectionMobileGsm:
		return true
	case vsectionMobileCdma:
		return true
	case vsectionWifi:
		return true
	case vsectionIpv4:
		return true
	case vsectionIpv6:
		return true
	case vsectionSecurity:
		return true
	case vsectionPppoe:
		return true
	case vsectionPpp:
		return true
	case vsectionVpn:
		return true
	case vsectionVpnL2tp:
		return true
	case vsectionVpnL2tpPpp:
		return true
	case vsectionVpnL2tpIpsec:
		return true
	case vsectionVpnOpenconnect:
		return true
	case vsectionVpnOpenvpn:
		return true
	case vsectionVpnOpenvpnAdvanced:
		return true
	case vsectionVpnOpenvpnSecurity:
		return true
	case vsectionVpnOpenvpnTlsauth:
		return true
	case vsectionVpnOpenvpnProxies:
		return true
	case vsectionVpnPptp:
		return true
	case vsectionVpnPptpPpp:
		return true
	case vsectionVpnVpnc:
		return true
	case vsectionVpnVpncAdvanced:
		return true
	}
	return false
}

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
		sections = append(sections, vsectionMobile)
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
			if isSettingSectionExists(data, section8021x) {
				sections = []string{section8021x}
			}
			sections = append(sections, vsectionSecurity)
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
		sections = append(sections, vsectionVpn)
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
