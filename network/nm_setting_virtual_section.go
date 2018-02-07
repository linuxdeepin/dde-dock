/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package network

import (
	"fmt"
	"pkg.deepin.io/dde/daemon/network/nm"
)

// Virtual cache section, used for connection session when editing
// some special virtual keys,
// e.g. nm.NM_SETTING_VK_MOBILE_COUNTRYREGION
const sectionCache = "cache"

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
			nm.NM_SETTING_VS_GENERAL,
			nm.NM_SETTING_VS_SECURITY,
			nm.NM_SETTING_VS_IPV4,
			nm.NM_SETTING_VS_IPV6,
			nm.NM_SETTING_VS_ETHERNET,
		}
	case connectionWireless:
		vsections = []string{
			nm.NM_SETTING_VS_GENERAL,
			nm.NM_SETTING_VS_SECURITY,
			nm.NM_SETTING_VS_IPV4,
			nm.NM_SETTING_VS_IPV6,
			nm.NM_SETTING_VS_WIFI,
		}
	case connectionWirelessAdhoc:
		vsections = []string{
			nm.NM_SETTING_VS_GENERAL,
			nm.NM_SETTING_VS_SECURITY,
			nm.NM_SETTING_VS_IPV4,
			nm.NM_SETTING_VS_IPV6,
			nm.NM_SETTING_VS_WIFI,
		}
	case connectionWirelessHotspot:
		vsections = []string{
			nm.NM_SETTING_VS_GENERAL,
			nm.NM_SETTING_VS_SECURITY,
			// nm.NM_SETTING_VS_IPV4,
			// nm.NM_SETTING_VS_IPV6,
			nm.NM_SETTING_VS_WIFI,
		}
	case connectionPppoe:
		vsections = []string{
			nm.NM_SETTING_VS_GENERAL,
			nm.NM_SETTING_VS_PPPOE,
			nm.NM_SETTING_VS_IPV4,
			nm.NM_SETTING_VS_ETHERNET,
			nm.NM_SETTING_VS_PPP,
		}
	case connectionVpnL2tp:
		vsections = []string{
			nm.NM_SETTING_VS_GENERAL,
			nm.NM_SETTING_VS_VPN,
			nm.NM_SETTING_VS_VPN_L2TP_PPP,
			nm.NM_SETTING_VS_VPN_L2TP_IPSEC,
			nm.NM_SETTING_VS_IPV4,
		}
	case connectionVpnOpenconnect:
		vsections = []string{
			nm.NM_SETTING_VS_GENERAL,
			nm.NM_SETTING_VS_VPN,
			nm.NM_SETTING_VS_IPV4,
			nm.NM_SETTING_VS_IPV6,
		}
	case connectionVpnOpenvpn:
		vsections = []string{
			nm.NM_SETTING_VS_GENERAL,
			nm.NM_SETTING_VS_VPN,
			nm.NM_SETTING_VS_VPN_OPENVPN_ADVANCED,
			nm.NM_SETTING_VS_VPN_OPENVPN_SECURITY,
			nm.NM_SETTING_VS_VPN_OPENVPN_PROXIES,
			nm.NM_SETTING_VS_IPV4,
			nm.NM_SETTING_VS_IPV6,
		}
		// when connection connection is static key, nm.NM_SETTING_VS_VPN_OPENVPN_TLSAUTH is not available
		if keepAll || getSettingVpnOpenvpnKeyConnectionType(data) != nm.NM_OPENVPN_CONTYPE_STATIC_KEY {
			vsections = append(vsections, nm.NM_SETTING_VS_VPN_OPENVPN_TLSAUTH)
		}
	case connectionVpnPptp:
		vsections = []string{
			nm.NM_SETTING_VS_GENERAL,
			nm.NM_SETTING_VS_VPN,
			nm.NM_SETTING_VS_VPN_PPTP_PPP,
			nm.NM_SETTING_VS_IPV4,
		}
	case connectionVpnStrongswan:
		vsections = []string{
			nm.NM_SETTING_VS_GENERAL,
			nm.NM_SETTING_VS_VPN,
			nm.NM_SETTING_VS_IPV4,
		}
	case connectionVpnVpnc:
		vsections = []string{
			nm.NM_SETTING_VS_GENERAL,
			nm.NM_SETTING_VS_VPN,
			nm.NM_SETTING_VS_VPN_VPNC_ADVANCED,
			nm.NM_SETTING_VS_IPV4,
		}
	case connectionMobileGsm:
		vsections = []string{
			nm.NM_SETTING_VS_GENERAL,
			nm.NM_SETTING_VS_MOBILE,
			nm.NM_SETTING_VS_PPP,
			nm.NM_SETTING_VS_IPV4,
		}
	case connectionMobileCdma:
		vsections = []string{
			nm.NM_SETTING_VS_GENERAL,
			nm.NM_SETTING_VS_MOBILE,
			nm.NM_SETTING_VS_PPP,
			nm.NM_SETTING_VS_IPV4,
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
// sections may contains virtual sections like nm.NM_SETTING_VS_MOBILE, for
// that some vkTypeController keys will be contained in them. If
// keepAll is true, will return all may used sections instead of
// filter some of them through the context.
func doGetRelatedSectionsOfVsection(data connectionData, vsection string, keepAll bool) (sections []string) {
	connectionType := getCustomConnectionType(data)
	switch vsection {
	default:
		logger.Error("getRelatedSectionsOfVsection: invalid vsection name", vsection)
	case nm.NM_SETTING_VS_GENERAL:
		sections = []string{nm.NM_SETTING_CONNECTION_SETTING_NAME}
	case nm.NM_SETTING_VS_MOBILE:
		sections = []string{nm.NM_SETTING_VS_MOBILE}
		switch connectionType {
		case connectionMobileGsm:
			sections = append(sections, nm.NM_SETTING_GSM_SETTING_NAME)
		case connectionMobileCdma:
			sections = append(sections, nm.NM_SETTING_CDMA_SETTING_NAME)
		}
	case nm.NM_SETTING_VS_MOBILE_GSM:
		sections = []string{nm.NM_SETTING_GSM_SETTING_NAME}
	case nm.NM_SETTING_VS_MOBILE_CDMA:
		sections = []string{nm.NM_SETTING_CDMA_SETTING_NAME}
	case nm.NM_SETTING_VS_ETHERNET:
		sections = []string{nm.NM_SETTING_WIRED_SETTING_NAME}
	case nm.NM_SETTING_VS_WIFI:
		sections = []string{nm.NM_SETTING_WIRELESS_SETTING_NAME}
	case nm.NM_SETTING_VS_IPV4:
		sections = []string{nm.NM_SETTING_IP4_CONFIG_SETTING_NAME}
	case nm.NM_SETTING_VS_IPV6:
		sections = []string{nm.NM_SETTING_IP6_CONFIG_SETTING_NAME}
	case nm.NM_SETTING_VS_SECURITY:
		switch connectionType {
		case connectionWired:
			sections = []string{nm.NM_SETTING_VS_SECURITY}
			if keepAll || isSettingExists(data, nm.NM_SETTING_802_1X_SETTING_NAME) {
				sections = append(sections, nm.NM_SETTING_802_1X_SETTING_NAME)
			}
		case connectionWireless, connectionWirelessAdhoc, connectionWirelessHotspot:
			sections = []string{nm.NM_SETTING_WIRELESS_SECURITY_SETTING_NAME}
			if keepAll || isSettingExists(data, nm.NM_SETTING_802_1X_SETTING_NAME) {
				sections = append(sections, nm.NM_SETTING_802_1X_SETTING_NAME)
			}
		}
	case nm.NM_SETTING_VS_PPPOE:
		sections = []string{nm.NM_SETTING_PPPOE_SETTING_NAME}
	case nm.NM_SETTING_VS_PPP:
		sections = []string{nm.NM_SETTING_PPP_SETTING_NAME}
	case nm.NM_SETTING_VS_VPN:
		switch connectionType {
		case connectionVpnL2tp:
			sections = []string{nm.NM_SETTING_ALIAS_VPN_L2TP_SETTING_NAME}
		case connectionVpnOpenconnect:
			sections = []string{nm.NM_SETTING_ALIAS_VPN_OPENCONNECT_SETTING_NAME}
		case connectionVpnOpenvpn:
			sections = []string{nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME}
		case connectionVpnPptp:
			sections = []string{nm.NM_SETTING_ALIAS_VPN_PPTP_SETTING_NAME}
		case connectionVpnStrongswan:
			sections = []string{nm.NM_SETTING_ALIAS_VPN_STRONGSWAN_SETTING_NAME}
		case connectionVpnVpnc:
			sections = []string{nm.NM_SETTING_ALIAS_VPN_VPNC_SETTING_NAME}
		}
		sections = append(sections, nm.NM_SETTING_VS_VPN)
	case nm.NM_SETTING_VS_VPN_L2TP:
		sections = []string{nm.NM_SETTING_ALIAS_VPN_L2TP_SETTING_NAME}
	case nm.NM_SETTING_VS_VPN_L2TP_PPP:
		sections = []string{nm.NM_SETTING_ALIAS_VPN_L2TP_PPP_SETTING_NAME}
	case nm.NM_SETTING_VS_VPN_L2TP_IPSEC:
		sections = []string{nm.NM_SETTING_ALIAS_VPN_L2TP_IPSEC_SETTING_NAME}
	case nm.NM_SETTING_VS_VPN_OPENCONNECT:
		sections = []string{nm.NM_SETTING_ALIAS_VPN_OPENCONNECT_SETTING_NAME}
	case nm.NM_SETTING_VS_VPN_OPENVPN:
		sections = []string{nm.NM_SETTING_ALIAS_VPN_OPENVPN_SETTING_NAME}
	case nm.NM_SETTING_VS_VPN_OPENVPN_ADVANCED:
		sections = []string{nm.NM_SETTING_ALIAS_VPN_OPENVPN_ADVANCED_SETTING_NAME}
	case nm.NM_SETTING_VS_VPN_OPENVPN_SECURITY:
		sections = []string{nm.NM_SETTING_ALIAS_VPN_OPENVPN_SECURITY_SETTING_NAME}
	case nm.NM_SETTING_VS_VPN_OPENVPN_TLSAUTH:
		sections = []string{nm.NM_SETTING_ALIAS_VPN_OPENVPN_TLSAUTH_SETTING_NAME}
	case nm.NM_SETTING_VS_VPN_OPENVPN_PROXIES:
		sections = []string{nm.NM_SETTING_ALIAS_VPN_OPENVPN_PROXIES_SETTING_NAME}
	case nm.NM_SETTING_VS_VPN_PPTP:
		sections = []string{nm.NM_SETTING_ALIAS_VPN_PPTP_SETTING_NAME}
	case nm.NM_SETTING_VS_VPN_PPTP_PPP:
		sections = []string{nm.NM_SETTING_ALIAS_VPN_PPTP_PPP_SETTING_NAME}
	case nm.NM_SETTING_VS_VPN_VPNC:
		sections = []string{nm.NM_SETTING_ALIAS_VPN_VPNC_SETTING_NAME}
	case nm.NM_SETTING_VS_VPN_VPNC_ADVANCED:
		sections = []string{nm.NM_SETTING_ALIAS_VPN_VPNC_ADVANCED_SETTING_NAME}
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
		if generalIsKeyShouldInSettingSection(section, key) {
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
	if vsection == nm.NM_SETTING_VS_GENERAL {
		return true
	}

	connectionType := getCustomConnectionType(data)
	switch connectionType {
	case connectionWired:
		switch vsection {
		case nm.NM_SETTING_VS_IPV4:
			expanded = true
		}
	case connectionPppoe:
		switch vsection {
		case nm.NM_SETTING_VS_PPPOE:
			expanded = true
		}
	case connectionWireless, connectionWirelessAdhoc, connectionWirelessHotspot:
		switch vsection {
		case nm.NM_SETTING_VS_IPV4, nm.NM_SETTING_VS_SECURITY:
			expanded = true
		}
	case connectionMobileGsm, connectionMobileCdma:
		switch vsection {
		case nm.NM_SETTING_VS_MOBILE:
			expanded = true
		}
	case connectionVpnL2tp:
		switch vsection {
		case nm.NM_SETTING_VS_VPN:
			expanded = true
		}
	case connectionVpnOpenconnect:
		switch vsection {
		case nm.NM_SETTING_VS_VPN:
			expanded = true
		}
	case connectionVpnOpenvpn:
		switch vsection {
		case nm.NM_SETTING_VS_VPN:
			expanded = true
		}
	case connectionVpnPptp:
		switch vsection {
		case nm.NM_SETTING_VS_VPN:
			expanded = true
		}
	case connectionVpnStrongswan:
		switch vsection {
		case nm.NM_SETTING_VS_VPN:
			expanded = true
		}
	case connectionVpnVpnc:
		switch vsection {
		case nm.NM_SETTING_VS_VPN:
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
		if section == nm.NM_SETTING_CONNECTION_SETTING_NAME && key == nm.NM_SETTING_CONNECTION_ID {
			return true
		}
	}
	return false
}
