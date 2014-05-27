package main

// Sections
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
	sectionVpnL2tp            = NM_SETTING_VS_VPN_L2TP_SETTING_NAME
	sectionVpnL2tpPpp         = NM_SETTING_VS_VPN_L2TP_PPP_SETTING_NAME
	sectionVpnL2tpIpsec       = NM_SETTING_VS_VPN_L2TP_IPSEC_SETTING_NAME
	sectionVpnOpenconnect     = NM_SETTING_VS_VPN_OPENCONNECT_SETTING_NAME
	sectionVpnOpenvpn         = NM_SETTING_VS_VPN_OPENVPN_SETTING_NAME
	sectionVpnOpenvpnAdvanced = NM_SETTING_VS_VPN_OPENVPN_ADVANCED_SETTING_NAME
	sectionVpnOpenvpnSecurity = NM_SETTING_VS_VPN_OPENVPN_SECURITY_SETTING_NAME
	sectionVpnOpenvpnTlsauth  = NM_SETTING_VS_VPN_OPENVPN_TLSAUTH_SETTING_NAME
	sectionVpnOpenvpnProxies  = NM_SETTING_VS_VPN_OPENVPN_PROXIES_SETTING_NAME
	sectionVpnPptp            = NM_SETTING_VS_VPN_PPTP_SETTING_NAME
	sectionVpnPptpPpp         = NM_SETTING_VS_VPN_PPTP_PPP_SETTING_NAME
	sectionVpnVpnc            = NM_SETTING_VS_VPN_VPNC_SETTING_NAME
	sectionVpnVpncAdvanced    = NM_SETTING_VS_VPN_VPNC_ADVANCED_SETTING_NAME
	sectionWired              = NM_SETTING_WIRED_SETTING_NAME
	sectionWireless           = NM_SETTING_WIRELESS_SETTING_NAME
	sectionWirelessSecurity   = NM_SETTING_WIRELESS_SECURITY_SETTING_NAME
)

// TODO refactor, alias sections
// Virtual sections, used for vpn connection which is a special key in fact
const (
	NM_SETTING_VS_VPN_L2TP_SETTING_NAME             = "vs-vpn-l2tp"
	NM_SETTING_VS_VPN_L2TP_PPP_SETTING_NAME         = "vs-vpn-l2tp-ppp"
	NM_SETTING_VS_VPN_L2TP_IPSEC_SETTING_NAME       = "vs-vpn-l2tp-ipsec"
	NM_SETTING_VS_VPN_OPENCONNECT_SETTING_NAME      = "vs-vpn-openconnect"
	NM_SETTING_VS_VPN_OPENVPN_SETTING_NAME          = "vs-vpn-openvpn"
	NM_SETTING_VS_VPN_OPENVPN_ADVANCED_SETTING_NAME = "vs-vpn-openvpn-advanced"
	NM_SETTING_VS_VPN_OPENVPN_SECURITY_SETTING_NAME = "vs-vpn-openvpn-security"
	NM_SETTING_VS_VPN_OPENVPN_TLSAUTH_SETTING_NAME  = "vs-vpn-openvpn-tlsauth"
	NM_SETTING_VS_VPN_OPENVPN_PROXIES_SETTING_NAME  = "vs-vpn-openvpn-proxies"
	NM_SETTING_VS_VPN_PPTP_SETTING_NAME             = "vs-vpn-pptp"
	NM_SETTING_VS_VPN_PPTP_PPP_SETTING_NAME         = "vs-vpn-pptp-ppp"
	NM_SETTING_VS_VPN_VPNC_SETTING_NAME             = "vs-vpn-vpnc"
	NM_SETTING_VS_VPN_VPNC_ADVANCED_SETTING_NAME    = "vs-vpn-advanced"
)

func getRealSectionName(name string) (realName string) {
	realName = name
	switch name {
	case NM_SETTING_VS_VPN_L2TP_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_VS_VPN_L2TP_PPP_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_VS_VPN_L2TP_IPSEC_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_VS_VPN_OPENCONNECT_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_VS_VPN_OPENVPN_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_VS_VPN_OPENVPN_ADVANCED_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_VS_VPN_OPENVPN_SECURITY_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_VS_VPN_OPENVPN_TLSAUTH_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_VS_VPN_OPENVPN_PROXIES_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_VS_VPN_PPTP_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_VS_VPN_PPTP_PPP_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_VS_VPN_VPNC_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_VS_VPN_VPNC_ADVANCED_SETTING_NAME:
		realName = sectionVpn
	}
	return
}

// Pages, page is a wrapper of sections for easy to configure
const (
	pageGeneral            = "general"              // -> sectionConnection
	pageEthernet           = "ethernet"             // -> sectionWireed
	pageMobile             = "mobile"               // -> sectionGsm
	pageMobileCdma         = "mobile-cdma"          // -> sectionCdma
	pageWifi               = "wifi"                 // -> sectionWireless
	pageIpv4               = "ipv4"                 // -> sectionIpv4
	pageIpv6               = "ipv6"                 // -> sectionIpv6
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

// TODO rename
// listPages return supported pages for target connection type.
func getAvailablePages(data connectionData) (pages []string) {
	connectionType := getCustomConnectionType(data)
	switch connectionType {
	case connectionWired:
		pages = []string{
			pageGeneral,
			pageEthernet,
			pageIpv4,
			pageIpv6,
			pageSecurity,
		}
	case connectionWireless:
		pages = []string{
			pageGeneral,
			pageWifi,
			pageIpv4,
			pageIpv6,
			pageSecurity,
		}
	case connectionWirelessAdhoc:
		pages = []string{
			pageGeneral,
			pageWifi,
			pageIpv4,
			pageIpv6,
			pageSecurity,
		}
	case connectionWirelessHotspot:
		pages = []string{
			pageGeneral,
			pageWifi,
			pageIpv4,
			pageIpv6,
			pageSecurity,
		}
	case connectionPppoe:
		pages = []string{
			pageGeneral,
			pageEthernet,
			pagePppoe,
			pagePpp,
			pageIpv4,
		}
	case connectionVpnL2tp:
		pages = []string{
			pageGeneral,
			pageVpnL2tp,
			pageVpnL2tpPpp,
			pageVpnL2tpIpsec,
			pageIpv4,
		}
	case connectionVpnOpenconnect:
		pages = []string{
			pageGeneral,
			pageVpnOpenconnect,
			pageIpv4,
			pageIpv6,
		}
	case connectionVpnOpenvpn:
		pages = []string{
			pageGeneral,
			pageVpnOpenvpn,
			pageVpnOpenvpnAdvanced,
			pageVpnOpenvpnSecurity,
			pageVpnOpenvpnProxies,
			pageIpv4,
			pageIpv6,
		}
		// when connection connection is static key, pageVpnOpenvpnTlsauth is not available
		if getSettingVpnOpenvpnKeyConnectionType(data) != NM_OPENVPN_CONTYPE_STATIC_KEY {
			pages = append(pages, pageVpnOpenvpnTlsauth)
		}
	case connectionVpnPptp:
		pages = []string{
			pageGeneral,
			pageVpnPptp,
			pageVpnPptpPpp,
			pageIpv4,
		}
	case connectionVpnVpnc:
		pages = []string{
			pageGeneral,
			pageVpnVpnc,
			pageVpnVpncAdvanced,
			pageIpv4,
		}
	case connectionMobileGsm:
		pages = []string{
			pageGeneral,
			pageMobile,
			pagePpp,
			pageIpv4,
		}
	case connectionMobileCdma:
		pages = []string{
			pageGeneral,
			pageMobileCdma,
			pagePpp,
			pageIpv4,
		}
	}
	return
}

func pageToSections(data connectionData, page string) (sections []string) {
	connectionType := getCustomConnectionType(data)
	switch page {
	default:
		logger.Error("pageToSections: invalid page name", page)
	case pageGeneral:
		sections = []string{sectionConnection}
	case pageMobile:
		sections = []string{sectionGsm}
	case pageMobileCdma:
		sections = []string{sectionCdma}
	case pageEthernet:
		sections = []string{sectionWired}
	case pageWifi:
		sections = []string{sectionWireless}
	case pageIpv4:
		sections = []string{sectionIpv4}
	case pageIpv6:
		sections = []string{sectionIpv6}
	case pageSecurity:
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
	case pagePppoe:
		sections = []string{sectionPppoe}
	case pagePpp:
		sections = []string{sectionPpp}
	case pageVpnL2tp:
		sections = []string{sectionVpnL2tp}
	case pageVpnL2tpPpp:
		sections = []string{sectionVpnL2tpPpp}
	case pageVpnL2tpIpsec:
		sections = []string{sectionVpnL2tpIpsec}
	case pageVpnOpenconnect:
		sections = []string{sectionVpnOpenconnect}
	case pageVpnOpenvpn:
		sections = []string{sectionVpnOpenvpn}
	case pageVpnOpenvpnAdvanced:
		sections = []string{sectionVpnOpenvpnAdvanced}
	case pageVpnOpenvpnSecurity:
		sections = []string{sectionVpnOpenvpnSecurity}
	case pageVpnOpenvpnTlsauth:
		sections = []string{sectionVpnOpenvpnTlsauth}
	case pageVpnOpenvpnProxies:
		sections = []string{sectionVpnOpenvpnProxies}
	case pageVpnPptp:
		sections = []string{sectionVpnPptp}
	case pageVpnPptpPpp:
		sections = []string{sectionVpnPptpPpp}
	case pageVpnVpnc:
		sections = []string{sectionVpnVpnc}
	case pageVpnVpncAdvanced:
		sections = []string{sectionVpnVpncAdvanced}
	}
	return
}

// TODO rename
// listSections return all pages related sections
func listSections(data connectionData) (sections []string) {
	for _, page := range getAvailablePages(data) {
		sections = appendStrArrayUnion(sections, pageToSections(data, page)...)
	}
	return
}

func getSectionOfPageKey(data connectionData, page, key string) string {
	sections := pageToSections(data, page)
	for _, section := range sections {
		if generalIsKeyInSettingSection(section, key) {
			return section
		}
	}
	logger.Errorf("get corresponding filed of key in page failed, page=%s, key=%s", page, key)
	return ""
}

// get available keys for target page
func getAvailableKeys(data connectionData, page string) (keys []string) {
	sections := pageToSections(data, page)
	for _, section := range sections {
		keys = appendStrArrayUnion(keys, generalGetSettingAvailableKeys(data, section)...)
	}
	if len(keys) == 0 {
		logger.Warning("there is no available keys for page", page)
	}
	return
}
