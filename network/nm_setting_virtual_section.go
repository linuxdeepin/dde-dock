/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package network

import (
	"fmt"
)

// Sections, correspondence to "NM_SETTING_XXX" in network manager.
const (
	// TODO: refactor code, add mappings to virtual sections and merge
	// sectionVpnL2tp with NM_SETTING_VS_VPN_L2TP
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
	sectionVpnL2tpIpsec       = NM_SETTING_ALIAS_VPN_L2TP_IPSEC_SETTING_NAME
	sectionVpnL2tpPpp         = NM_SETTING_ALIAS_VPN_L2TP_PPP_SETTING_NAME
	sectionVpnStrongswan      = NM_SETTING_ALIAS_VPN_STRONGSWAN_SETTING_NAME
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
	NM_SETTING_ALIAS_VPN_STRONGSWAN_SETTING_NAME       = "alias-vpn-strongswan"
	NM_SETTING_ALIAS_VPN_VPNC_SETTING_NAME             = "alias-vpn-vpnc"
	NM_SETTING_ALIAS_VPN_VPNC_ADVANCED_SETTING_NAME    = "alias-vpn-vpnc-advanced"
)

// Cache section, used for connection session when editing some
// special virtual keys, e.g. NM_SETTING_VK_MOBILE_COUNTRYREGION
const sectionCache = "cache"

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
	case NM_SETTING_ALIAS_VPN_STRONGSWAN_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_ALIAS_VPN_VPNC_SETTING_NAME:
		realName = sectionVpn
	case NM_SETTING_ALIAS_VPN_VPNC_ADVANCED_SETTING_NAME:
		realName = sectionVpn
	}
	return
}

// VsectionInfo defines virtual section, used by front-end to build connection edit page run time.
type VsectionInfo struct {
	relatedSections []string // all related sections

	// VirtualSection is the id of the virtual section, such as
	// "vs-general".
	VirtualSection string

	// Name is the display name for front-end expand widget.
	Name string

	// Expanded tells if the front-end expand widget should be
	// expanded default.
	Expanded bool

	// Keys contains all the keys information under current virtual
	// section.
	Keys []*GeneralKeyInfo
}

// GeneralKeyInfo defines settings key that used by front-end.
type GeneralKeyInfo struct {
	// Section most time is a real section name such as "connection",
	// but will a virtual section for vkTypeController virtual key,
	// such as "vk-mobile-service-type" for "vs-mobile".
	Section string

	// Key will be a real or virtual key name, such as "id" and "vk-eap".
	Key string

	// Name is the display name for front-end widget.
	Name string

	// WidgetType is the showing widget type to create for the key run
	// time.
	WidgetType string

	// Readonly tells if the front-end widget should edit the key value.
	Readonly bool

	// AlwaysUpdate tells if front-end widget should re-get value once
	// other keys changed.
	AlwaysUpdate bool

	// UseValueRange tells the custom value range will be used for
	// integer keys.
	UseValueRange bool
	MinValue      int
	MaxValue      int
}

const (
	vsectionGeneral            = NM_SETTING_VS_GENERAL              // -> sectionConnection
	vsectionEthernet           = NM_SETTING_VS_ETHERNET             // -> sectionWired
	vsectionMobile             = NM_SETTING_VS_MOBILE               // -> sectionGsm, sectionCdma
	vsectionMobileGsm          = NM_SETTING_VS_MOBILE_GSM           // -> sectionGsm // TODO: remove
	vsectionMobileCdma         = NM_SETTING_VS_MOBILE_CDMA          // -> sectionCdma
	vsectionWifi               = NM_SETTING_VS_WIFI                 // -> sectionWireless
	vsectionIpv4               = NM_SETTING_VS_IPV4                 // -> sectionIpv4
	vsectionIpv6               = NM_SETTING_VS_IPV6                 // -> sectionIpv6
	vsectionSecurity           = NM_SETTING_VS_SECURITY             // -> section8021x, sectionWirelessSecurity
	vsectionPppoe              = NM_SETTING_VS_PPPOE                // -> sectionPppoe
	vsectionPpp                = NM_SETTING_VS_PPP                  // -> sectionPpp
	vsectionVpn                = NM_SETTING_VS_VPN                  // -> sectionVpnL2tp, sectionVpnOpenconnect, sectionVpnOpenvpn, sectionVpnPptp, sectionVpnVpnc
	vsectionVpnL2tp            = NM_SETTING_VS_VPN_L2TP             // -> sectionVpnL2tp
	vsectionVpnL2tpPpp         = NM_SETTING_VS_VPN_L2TP_PPP         // -> sectionVpnL2tpPpp
	vsectionVpnL2tpIpsec       = NM_SETTING_VS_VPN_L2TP_IPSEC       // -> sectionVpnL2tpIpsec
	vsectionVpnOpenconnect     = NM_SETTING_VS_VPN_OPENCONNECT      // -> sectionVpnOpenconnect
	vsectionVpnOpenvpn         = NM_SETTING_VS_VPN_OPENVPN          // -> sectionVpnOpenvpn
	vsectionVpnOpenvpnAdvanced = NM_SETTING_VS_VPN_OPENVPN_ADVANCED // -> sectionVpnOpenVpnAdvanced
	vsectionVpnOpenvpnSecurity = NM_SETTING_VS_VPN_OPENVPN_SECURITY // -> sectionVpnOpenVpnSecurity
	vsectionVpnOpenvpnTlsauth  = NM_SETTING_VS_VPN_OPENVPN_TLSAUTH  // -> sectionVpnOpenVpnTlsauth
	vsectionVpnOpenvpnProxies  = NM_SETTING_VS_VPN_OPENVPN_PROXIES  // -> sectionVpnOpenVpnProxies
	vsectionVpnStrongswan      = NM_SETTING_VS_VPN_STRONGSWAN       // -> sectionVpnStrongswan
	vsectionVpnPptp            = NM_SETTING_VS_VPN_PPTP             // -> sectionVpnPptp
	vsectionVpnPptpPpp         = NM_SETTING_VS_VPN_PPTP_PPP         // -> sectionVpnPptpPpp
	vsectionVpnVpnc            = NM_SETTING_VS_VPN_VPNC             // -> sectionVpnVpnc
	vsectionVpnVpncAdvanced    = NM_SETTING_VS_VPN_VPNC_ADVANCED    // -> sectionVpnVpncAdvanced
)

func isVirtualSection(section string) bool {
	for _, vs := range virtualSections {
		if vs.VirtualSection == section {
			return true
		}
	}
	return false
}

// get avaiable virtual sections for target connection type run time
func getAvailableVsections(data connectionData) (vsections []string) {
	return doGetRelatedVsections(data, false)
}

// get all related virtual sections for target connection type
func getAllVsections(data connectionData) (vsections []string) {
	return doGetRelatedVsections(data, true)
}

// get related virtual sections for target connection type. If
// keepAll is true, will return all may used sections instead of
// filter some of theme through the context.
func doGetRelatedVsections(data connectionData, keepAll bool) (vsections []string) {
	connectionType := getCustomConnectionType(data)
	switch connectionType {
	case connectionWired:
		vsections = []string{
			vsectionGeneral,
			vsectionEthernet,
			vsectionSecurity,
			vsectionIpv4,
			vsectionIpv6,
		}
	case connectionWireless:
		vsections = []string{
			vsectionGeneral,
			vsectionWifi,
			vsectionSecurity,
			vsectionIpv4,
			vsectionIpv6,
		}
	case connectionWirelessAdhoc:
		vsections = []string{
			vsectionGeneral,
			vsectionWifi,
			vsectionSecurity,
			vsectionIpv4,
			vsectionIpv6,
		}
	case connectionWirelessHotspot:
		vsections = []string{
			vsectionGeneral,
			vsectionWifi,
			vsectionSecurity,
			vsectionIpv4,
			vsectionIpv6,
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
		if keepAll || getSettingVpnOpenvpnKeyConnectionType(data) != NM_OPENVPN_CONTYPE_STATIC_KEY {
			vsections = append(vsections, vsectionVpnOpenvpnTlsauth)
		}
	case connectionVpnPptp:
		vsections = []string{
			vsectionGeneral,
			vsectionVpn,
			vsectionVpnPptpPpp,
			vsectionIpv4,
		}
	case connectionVpnStrongswan:
		vsections = []string{
			vsectionGeneral,
			vsectionVpn,
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

// get available sections of virtual section run time
func getAvailableSectionsOfVsection(data connectionData, vsection string) (sections []string) {
	return doGetRelatedSectionsOfVsection(data, vsection, false)
}

// get all related sections of virtual section
func getAllSectionsOfVsection(data connectionData, vsection string) (sections []string) {
	return doGetRelatedSectionsOfVsection(data, vsection, true)
}

// get related sections of virtual section run time. The returned
// sections may contains virtual sections like vsectionMobile, for
// that some vkTypeController keys will be contained in them. If
// keepAll is true, will return all may used sections instead of
// filter some of them through the context.
func doGetRelatedSectionsOfVsection(data connectionData, vsection string, keepAll bool) (sections []string) {
	connectionType := getCustomConnectionType(data)
	switch vsection {
	default:
		logger.Error("getRelatedSectionsOfVsection: invalid vsection name", vsection)
	case vsectionGeneral:
		sections = []string{sectionConnection}
	case vsectionMobile:
		sections = []string{vsectionMobile}
		switch connectionType {
		case connectionMobileGsm:
			sections = append(sections, sectionGsm)
		case connectionMobileCdma:
			sections = append(sections, sectionCdma)
		}
	case vsectionMobileGsm:
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
			sections = []string{vsectionSecurity}
			if keepAll || isSettingSectionExists(data, section8021x) {
				sections = append(sections, section8021x)
			}
		case connectionWireless, connectionWirelessAdhoc, connectionWirelessHotspot:
			sections = []string{sectionWirelessSecurity}
			if keepAll || isSettingSectionExists(data, section8021x) {
				sections = append(sections, section8021x)
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
		case connectionVpnStrongswan:
			sections = []string{sectionVpnStrongswan}
		case connectionVpnVpnc:
			sections = []string{sectionVpnVpnc}
		}
		sections = append(sections, vsectionVpn)
	case vsectionVpnL2tp:
		sections = []string{sectionVpnL2tp}
	case vsectionVpnL2tpPpp:
		sections = []string{sectionVpnL2tpPpp}
	case vsectionVpnL2tpIpsec:
		sections = []string{sectionVpnL2tpIpsec}
	case vsectionVpnOpenconnect:
		sections = []string{sectionVpnOpenconnect}
	case vsectionVpnOpenvpn:
		sections = []string{sectionVpnOpenvpn}
	case vsectionVpnOpenvpnAdvanced:
		sections = []string{sectionVpnOpenvpnAdvanced}
	case vsectionVpnOpenvpnSecurity:
		sections = []string{sectionVpnOpenvpnSecurity}
	case vsectionVpnOpenvpnTlsauth:
		sections = []string{sectionVpnOpenvpnTlsauth}
	case vsectionVpnOpenvpnProxies:
		sections = []string{sectionVpnOpenvpnProxies}
	case vsectionVpnPptp:
		sections = []string{sectionVpnPptp}
	case vsectionVpnPptpPpp:
		sections = []string{sectionVpnPptpPpp}
	case vsectionVpnVpnc:
		sections = []string{sectionVpnVpnc}
	case vsectionVpnVpncAdvanced:
		sections = []string{sectionVpnVpncAdvanced}
	}
	return
}

// getAvailableSections return all related available sections run time
func getAvailableSections(data connectionData) (sections []string) {
	return doGetRelatedSections(data, false)
}

// getAllSections return all related sections
func getAllSections(data connectionData) (sections []string) {
	return doGetRelatedSections(data, true)
}

func doGetRelatedSections(data connectionData, keepAll bool) (sections []string) {
	for _, vsection := range doGetRelatedVsections(data, keepAll) {
		sections = appendStrArrayUnique(sections, doGetRelatedSectionsOfVsection(data, vsection, keepAll)...)
	}
	return
}

// get real section name of target key in virtual section
func getSectionOfKeyInVsection(data connectionData, vsection, key string) (section string) {
	sections := doGetRelatedSectionsOfVsection(data, vsection, false)
	for _, section := range sections {
		if generalIsKeyInSettingSection(section, key) {
			return section
		}
	}
	logger.Errorf("get corresponding section of key in virtual section failed, vsection=%s, key=%s", vsection, key)
	return ""
}

func isKeyAvailable(data connectionData, section, key string) bool {
	if isStringInArray(section, getAvailableSections(data)) {
		if isStringInArray(key, generalGetSettingAvailableKeys(data, section)) {
			return true
		}
	}
	return false
}

// get all available keys of virtual section run time
func getAvailableKeysOfVsection(data connectionData, vsection string) (keys []string) {
	return doGetRelatedKeysOfVsection(data, vsection, false)
}

// get all related keys of virtual section run time
func getAllKeysOfVsection(data connectionData, vsection string) (keys []string) {
	return doGetRelatedKeysOfVsection(data, vsection, true)
}

func doGetRelatedKeysOfVsection(data connectionData, vsection string, keepAll bool) (keys []string) {
	sections := doGetRelatedSectionsOfVsection(data, vsection, keepAll)
	for _, section := range sections {
		keys = appendStrArrayUnique(keys, generalGetSettingAvailableKeys(data, section)...)
	}
	if len(keys) == 0 {
		logger.Warning("there is no available keys for virtual section", vsection)
	}
	return
}

func getRelatedKeyName(data connectionData, vsection, key string) (name string, err error) {
	keyInfo, err := getGeneralKeyInfo(data, vsection, key)
	if keyInfo != nil {
		name = keyInfo.Name
	}
	return
}

func getGeneralKeyInfo(data connectionData, vsection, key string) (keyInfo *GeneralKeyInfo, err error) {
	for _, sectionInfo := range virtualSections {
		for _, tmpKeyInfo := range sectionInfo.Keys {
			if tmpKeyInfo.Section == vsection && tmpKeyInfo.Key == key {
				keyInfo = tmpKeyInfo
				return
			}
		}
	}
	err = fmt.Errorf("get key information failed for vsection=%s, key=%s", vsection, key)
	logger.Error(err)
	return
}

func (vs *VsectionInfo) fixExpanded(data connectionData) {
	vs.Expanded = isVsectionExpandedDefault(data, vs.VirtualSection)
}
func isVsectionExpandedDefault(data connectionData, vsection string) (expanded bool) {
	if vsection == vsectionGeneral {
		return true
	}

	connectionType := getCustomConnectionType(data)
	switch connectionType {
	case connectionWired:
		switch vsection {
		case vsectionIpv4:
			expanded = true
		}
	case connectionPppoe:
		switch vsection {
		case vsectionPppoe:
			expanded = true
		}
	case connectionWireless, connectionWirelessAdhoc, connectionWirelessHotspot:
		switch vsection {
		case vsectionIpv4, vsectionSecurity:
			expanded = true
		}
	case connectionMobileGsm, connectionMobileCdma:
		switch vsection {
		case vsectionMobile:
			expanded = true
		}
	case connectionVpnL2tp:
		switch vsection {
		case vsectionVpn:
			expanded = true
		}
	case connectionVpnOpenconnect:
		switch vsection {
		case vsectionVpn:
			expanded = true
		}
	case connectionVpnOpenvpn:
		switch vsection {
		case vsectionVpn:
			expanded = true
		}
	case connectionVpnPptp:
		switch vsection {
		case vsectionVpn:
			expanded = true
		}
	case connectionVpnStrongswan:
		switch vsection {
		case vsectionVpn:
			expanded = true
		}
	case connectionVpnVpnc:
		switch vsection {
		case vsectionVpn:
			expanded = true
		}
	default:
		logger.Error("unknown custom connection type", connectionType)
	}
	return
}

func (k *GeneralKeyInfo) fixReadonly(data connectionData) {
	k.Readonly = isGeneralKeyReadonly(data, k.Section, k.Key)
}
func isGeneralKeyReadonly(data connectionData, section, key string) (readonly bool) {
	connectionType := getCustomConnectionType(data)
	switch connectionType {
	case connectionWired, connectionMobileGsm, connectionMobileCdma, connectionWireless, connectionWirelessAdhoc, connectionWirelessHotspot:
		if section == sectionConnection && key == NM_SETTING_CONNECTION_ID {
			return true
		}
	}
	return false
}
