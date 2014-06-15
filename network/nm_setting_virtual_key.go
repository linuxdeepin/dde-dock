/**
 * Copyright (c) 2014 Deepin, Inc.
 *               2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package network

import . "dlib/gettext"
import "fmt"

// virtual key types
const (
	vkTypeWrapper       = "wrapper"
	vkTypeEnableWrapper = "enable-wrapper"

	// control other sections or keys, no related key, and the related
	// section always is a virtual section, such as "vk-vpn-type", for
	// there is no real related section, the key's name must be unique
	vkTypeController = "controller"
)

type vkeyInfo struct {
	Value          string
	Type           ktype
	VkType         string // could be "wrapper", "enable-wrapper", "control"
	RelatedSection string
	RelatedKeys    []string
	Available      bool // check if is used by front-end
	Optional       bool // if key is optional(such as child key gateway of ip address), will ignore error for it
}

func getVkeyInfo(section, vkey string) (info vkeyInfo, ok bool) {
	for _, vk := range virtualKeys {
		if vk.RelatedSection == section && vk.Value == vkey {
			info = vk
			ok = true
			return
		}
	}
	logger.Errorf("invalid virtual key, section=%s, vkey=%s", section, vkey)
	ok = false
	return
}

func isVirtualKey(section, key string) bool {
	if isStringInArray(key, getVkeysOfSection(section)) {
		return true
	}
	return false
}

// get all virtual keys in target section
func getVkeysOfSection(section string) (vks []string) {
	for _, vk := range virtualKeys {
		if vk.RelatedSection == section {
			vks = append(vks, vk.Value)
		}
	}
	// logger.Debug("getVkeysOfSection: filed:", section, vks) // TODO test
	return
}

func getSettingVkeyType(section, key string) (t ktype) {
	t = ktypeUnknown
	for _, vk := range virtualKeys {
		if vk.RelatedSection == section && vk.Value == key {
			t = vk.Type
		}
	}
	if t == ktypeUnknown {
		logger.Errorf("get virtual key type failed, section=%s, key=%s", section, key)
	}
	return
}

func generalGetSettingVsectionAvailableKeys(data connectionData, vsection string) (keys []string) {
	switch vsection {
	case NM_SETTING_VS_SECURITY:
		if getCustomConnectionType(data) == connectionWired {
			keys = []string{NM_SETTING_VK_802_1X_ENABLE}
		}
	case NM_SETTING_VS_MOBILE:
		keys = []string{NM_SETTING_VK_MOBILE_SERVICE_TYPE}
	case NM_SETTING_VS_VPN:
		keys = []string{NM_SETTING_VK_VPN_TYPE}
		if !isStringInArray(getSettingVkVpnType(data), getLocalSupportedVpnTypes()) {
			keys = append(keys, NM_SETTING_VK_VPN_MISSING_PLUGIN)
		}
	}
	return
}

func generalGetSettingVkeyAvailableValues(data connectionData, section, key string) (values []kvalue) {
	switch section {
	case NM_SETTING_VS_MOBILE:
		switch key {
		case NM_SETTING_VK_MOBILE_SERVICE_TYPE:
			values = []kvalue{
				kvalue{connectionMobileGsm, Tr("GSM (GPRS, EDGE, UMTS, HSPA)")},
				kvalue{connectionMobileCdma, Tr("CDMA (1xRTT, EVDO)")},
			}
		}
	case NM_SETTING_VS_VPN:
		switch key {
		case NM_SETTING_VK_VPN_TYPE:
			values = []kvalue{
				kvalue{connectionVpnL2tp, Tr("L2TP")},
				kvalue{connectionVpnPptp, Tr("PPTP")},
				kvalue{connectionVpnOpenconnect, Tr("OpenConnect")},
				kvalue{connectionVpnOpenvpn, Tr("OpenVPN")},
				kvalue{connectionVpnVpnc, Tr("VPNC")},
			}
		}
	case section8021x:
		switch key {
		case NM_SETTING_VK_802_1X_EAP:
			values = getSetting8021xAvailableValues(data, NM_SETTING_802_1X_EAP)
		}
	case sectionWirelessSecurity:
		switch key {
		case NM_SETTING_VK_WIRELESS_SECURITY_KEY_MGMT:
			if getSettingWirelessMode(data) == NM_SETTING_WIRELESS_MODE_INFRA {
				values = []kvalue{
					kvalue{"none", Tr("None")},
					kvalue{"wep", Tr("WEP 40/128-bit Key")},
					kvalue{"wpa-psk", Tr("WPA/WPA2 Personal")},
					kvalue{"wpa-eap", Tr("WPA/WPA2 Enterprise")},
				}
			} else {
				values = []kvalue{
					kvalue{"none", Tr("None")},
					kvalue{"wep", Tr("WEP 40/128-bit Key")},
					kvalue{"wpa-psk", Tr("WPA/WPA2 Personal")},
				}
			}
		}
	case sectionVpnL2tpPpp:
		switch key {
		case NM_SETTING_VK_VPN_L2TP_MPPE_SECURITY:
			values = []kvalue{
				kvalue{"default", Tr("All Available (default)")},
				kvalue{"128-bit", Tr("128-bit (most secure)")},
				kvalue{"40-bit", Tr("40-bit (less secure)")},
			}
		}
	case sectionVpnPptpPpp:
		switch key {
		case NM_SETTING_VK_VPN_PPTP_MPPE_SECURITY:
			values = []kvalue{
				kvalue{"default", Tr("All Available (default)")},
				kvalue{"128-bit", Tr("128-bit (most secure)")},
				kvalue{"40-bit", Tr("40-bit (less secure)")},
			}
		}
	case sectionVpnVpncAdvanced:
		switch key {
		case NM_SETTING_VK_VPN_VPNC_KEY_ENCRYPTION_METHOD:
			values = []kvalue{
				kvalue{"secure", Tr("Secure (default)")},
				kvalue{"weak", Tr("Weak")},
				kvalue{"none", Tr("None")},
			}
		}
	}

	if len(values) == 0 {
		logger.Warningf("there is no available values for virtual key, %s->%s", section, key)
	}
	return
}

// general function to append available keys, will dispatch virtual keys specially
func appendAvailableKeys(data connectionData, keys []string, section, key string) (newKeys []string) {
	newKeys = appendStrArrayUnique(keys)
	relatedVks := getRelatedAvailableVkeys(section, key)
	if len(relatedVks) > 0 {
		for _, vk := range relatedVks {
			// if is enable wrapper virtual key, both virtual key and
			// real key will be appended
			if isEnableWrapperVkey(section, vk) {
				if isSettingKeyExists(data, section, key) {
					newKeys = appendStrArrayUnique(newKeys, key)
				}
			}
		}
		newKeys = appendStrArrayUnique(newKeys, relatedVks...)
	} else {
		newKeys = appendStrArrayUnique(newKeys, key)
	}
	return
}

func getRelatedAvailableVkeys(section, key string) (vks []string) {
	for _, vk := range virtualKeys {
		if vk.RelatedSection == section && isStringInArray(key, vk.RelatedKeys) && vk.Available {
			vks = append(vks, vk.Value)
		}
	}
	return
}

// get related virtual keys of target key
func getRelatedVkeys(section, key string) (vks []string) {
	for _, vk := range virtualKeys {
		if vk.RelatedSection == section && isStringInArray(key, vk.RelatedKeys) {
			vks = append(vks, vk.Value)
		}
	}
	return
}

// remove virtual key means to remove its related keys
func removeVirtualKey(data connectionData, section, vkey string) {
	if isControllerVkey(section, vkey) {
		// ignore controller virtual key for there is no related key for it
		return
	}
	vkeyInfo, ok := getVkeyInfo(section, vkey)
	if !ok {
		return
	}
	for _, key := range vkeyInfo.RelatedKeys {
		removeSettingKey(data, vkeyInfo.RelatedSection, key)
	}
}

func isWrapperVkey(section, vkey string) bool {
	vkInfo, ok := getVkeyInfo(section, vkey)
	if !ok {
		return false
	}
	if vkInfo.VkType == vkTypeWrapper {
		return true
	}
	return false
}

func isEnableWrapperVkey(section, vkey string) bool {
	vkInfo, ok := getVkeyInfo(section, vkey)
	if !ok {
		return false
	}
	if vkInfo.VkType == vkTypeEnableWrapper {
		return true
	}
	return false
}

func isControllerVkey(section, vkey string) bool {
	vkInfo, ok := getVkeyInfo(section, vkey)
	if !ok {
		return false
	}
	if vkInfo.VkType == vkTypeController {
		return true
	}
	return false
}

func isOptionalVkey(section, vkey string) (optional bool) {
	for _, vk := range virtualKeys {
		if vk.RelatedSection == section && vk.Value == vkey {
			optional = vk.Optional
		}
	}
	return
}

// Controller virtual key, which with no real related section
func logicSetSettingVk8021xEnable(data connectionData, value bool) (err error) {
	if value {
		addSettingSection(data, section8021x)
		err = logicSetSettingVk8021xEap(data, "tls")
	} else {
		removeSettingSection(data, section8021x)
	}
	return
}
func logicSetSettingVk8021xEap(data connectionData, value string) (err error) {
	return logicSetSetting8021xEap(data, []string{value})
}

func getSettingVkMobileServiceType(data connectionData) (serviceType string) {
	if isSettingSectionExists(data, sectionGsm) {
		serviceType = connectionMobileGsm
	} else if isSettingSectionExists(data, sectionCdma) {
		serviceType = connectionMobileCdma
	} else {
		logger.Error("get mobile service type failed, neither gsm section nor cdma section")
	}
	return
}
func logicSetSettingVkMobileServiceType(data connectionData, serviceType string) (err error) {
	switch serviceType {
	case connectionMobileGsm:
		removeSettingSection(data, sectionCdma)
		initSettingSectionGsm(data)
	case connectionMobileCdma:
		removeSettingSection(data, sectionGsm)
		initSettingSectionCdma(data)
	default:
		err = fmt.Errorf("invalid mobile service type", serviceType)
	}
	return
}

func getSettingVkVpnType(data connectionData) (vpnType string) {
	vpnType = getCustomConnectionType(data)
	return
}
func logicSetSettingVkVpnType(data connectionData, vpnType string) (err error) {
	removeSettingSection(data, sectionVpn)
	removeSettingSection(data, sectionIpv6)
	switch vpnType {
	case connectionVpnL2tp:
		initSettingSectionVpnL2tp(data)
	case connectionVpnPptp:
		initSettingSectionVpnPptp(data)
	case connectionVpnOpenconnect:
		initSettingSectionVpnOpenconnect(data)
		initSettingSectionIpv6(data)
	case connectionVpnOpenvpn:
		initSettingSectionVpnOpenvpn(data)
		initSettingSectionIpv6(data)
	case connectionVpnVpnc:
		initSettingSectionVpnVpnc(data)
	default:
		err = fmt.Errorf("invalid vpn type", vpnType)
	}
	return
}

func getSettingVkVpnMissingPlugin(data connectionData) (missingPlugin string) {
	vpnType := getCustomConnectionType(data)
	if !isStringInArray(vpnType, getLocalSupportedVpnTypes()) {
		switch vpnType {
		case connectionVpnL2tp:
			missingPlugin = "network-manager-l2tp"
		case connectionVpnPptp:
			missingPlugin = "network-manager-pptp"
		case connectionVpnOpenconnect:
			missingPlugin = "network-manager-openconnect"
		case connectionVpnOpenvpn:
			missingPlugin = "network-manager-openvpn"
		case connectionVpnVpnc:
			missingPlugin = "network-manager-vpnc"
		default:
			fmt.Errorf("invalid vpn type", vpnType)
		}
	}
	return
}
func logicSetSettingVkVpnMissingPlugin(data connectionData, vpnType string) (err error) {
	return
}
