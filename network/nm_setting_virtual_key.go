/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/iso"
	"pkg.deepin.io/lib/mobileprovider"
)

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
	value          string
	ktype          ktype
	vkType         string // could be "wrapper", "enable-wrapper", "controller"
	relatedSection string
	relatedKeys    []string

	// child keys that split from one key, such as IP address, mask
	// and gateway for IP addresses
	childKey bool

	// if key is optional(such as the child key gateway for ip
	// addresses), then will ignore missing error for it
	optional bool
}

func getVkeyInfo(section, vkey string) (info vkeyInfo, ok bool) {
	for _, vk := range virtualKeys {
		if vk.relatedSection == section && vk.value == vkey {
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
		if vk.relatedSection == section {
			vks = append(vks, vk.value)
		}
	}
	return
}

func getSettingVkeyType(section, key string) (t ktype) {
	t = ktypeUnknown
	for _, vk := range virtualKeys {
		if vk.relatedSection == section && vk.value == key {
			t = vk.ktype
		}
	}
	if t == ktypeUnknown {
		logger.Errorf("get virtual key type failed, section=%s, key=%s", section, key)
	}
	return
}

func generalGetSettingVsectionAvailableKeys(data connectionData, vsection string) (keys []string) {
	switch vsection {
	case nm.NM_SETTING_VS_SECURITY:
		if getCustomConnectionType(data) == connectionWired {
			keys = []string{nm.NM_SETTING_VK_802_1X_ENABLE}
		}
	case nm.NM_SETTING_VS_MOBILE:
		keys = []string{
			nm.NM_SETTING_VK_MOBILE_COUNTRY,
			nm.NM_SETTING_VK_MOBILE_PROVIDER,
		}
		if getSettingVkMobileProvider(data) == mobileProviderValueCustom {
			keys = append(keys, nm.NM_SETTING_VK_MOBILE_SERVICE_TYPE)
		} else {
			keys = append(keys, nm.NM_SETTING_VK_MOBILE_PLAN)
			// TODO: is "apn-readonly" widget necessary?
			// if getSettingVkMobileServiceType(data) == connectionMobileGsm {
			// keys = append(keys, nm.NM_SETTING_VK_MOBILE_APN_READONLY)
			// }
		}
	case nm.NM_SETTING_VS_VPN:
		keys = []string{}
		if !isStringInArray(getSettingVkVpnType(data), getLocalSupportedVpnTypes()) {
			keys = append(keys, nm.NM_SETTING_VK_VPN_MISSING_PLUGIN)
		}
	}
	return
}

func generalGetSettingVkeyAvailableValues(data connectionData, section, key string) (values []kvalue) {
	switch section {
	case nm.NM_SETTING_VS_MOBILE:
		switch key {
		case nm.NM_SETTING_VK_MOBILE_COUNTRY:
			codeList, _ := mobileprovider.GetAllCountryCode()
			for _, code := range codeList {
				if name, err := iso.GetCountryNameForCode(code); err == nil {
					values = append(values, kvalue{code, name})
				} else {
					logger.Error(err, code)
				}
			}
			// sort country list
			sortKvalues(values)
		case nm.NM_SETTING_VK_MOBILE_PROVIDER:
			countryCode := getSettingVkMobileCountry(data)
			names, _ := mobileprovider.GetProviderNames(countryCode)
			for _, name := range names {
				values = append(values, kvalue{name, name})
			}
			values = append(values, kvalue{mobileProviderValueCustom, Tr("Custom")})
		case nm.NM_SETTING_VK_MOBILE_PLAN:
			countryCode := getSettingVkMobileCountry(data)
			providerName := getSettingVkMobileProvider(data)
			plans, _ := mobileprovider.GetPlans(countryCode, providerName)
			for _, p := range plans {
				if len(p.Name) > 0 {
					values = append(values, kvalue{mobileprovider.MarshalPlan(p), p.Name})
				} else {
					values = append(values, kvalue{mobileprovider.MarshalPlan(p), Tr("Default")})
				}
			}
		case nm.NM_SETTING_VK_MOBILE_SERVICE_TYPE:
			values = []kvalue{
				kvalue{connectionMobileGsm, Tr("GSM (GPRS, UMTS)")},
				kvalue{connectionMobileCdma, Tr("CDMA (1xRTT, EVDO)")},
			}
		}
	case nm.NM_SETTING_VS_VPN:
		switch key {
		case nm.NM_SETTING_VK_VPN_TYPE:
			values = []kvalue{
				kvalue{connectionVpnL2tp, Tr("L2TP")},
				kvalue{connectionVpnPptp, Tr("PPTP")},
				kvalue{connectionVpnOpenconnect, Tr("OpenConnect")},
				kvalue{connectionVpnOpenvpn, Tr("OpenVPN")},
				kvalue{connectionVpnVpnc, Tr("VPNC")},
			}
		}
	case nm.NM_SETTING_802_1X_SETTING_NAME:
		switch key {
		case nm.NM_SETTING_VK_802_1X_EAP:
			values = getSetting8021xAvailableValues(data, nm.NM_SETTING_802_1X_EAP)
		}
	case nm.NM_SETTING_WIRELESS_SECURITY_SETTING_NAME:
		switch key {
		case nm.NM_SETTING_VK_WIRELESS_SECURITY_KEY_MGMT:
			if getSettingWirelessMode(data) == nm.NM_SETTING_WIRELESS_MODE_INFRA {
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
	case nm.NM_SETTING_ALIAS_VPN_L2TP_PPP_SETTING_NAME:
		switch key {
		case nm.NM_SETTING_VK_VPN_L2TP_MPPE_SECURITY:
			values = []kvalue{
				kvalue{"default", Tr("All Available (default)")},
				kvalue{"128-bit", Tr("128-bit (most secure)")},
				kvalue{"40-bit", Tr("40-bit (less secure)")},
			}
		}
	case nm.NM_SETTING_ALIAS_VPN_PPTP_PPP_SETTING_NAME:
		switch key {
		case nm.NM_SETTING_VK_VPN_PPTP_MPPE_SECURITY:
			values = []kvalue{
				kvalue{"default", Tr("All Available (default)")},
				kvalue{"128-bit", Tr("128-bit (most secure)")},
				kvalue{"40-bit", Tr("40-bit (less secure)")},
			}
		}
	case nm.NM_SETTING_ALIAS_VPN_VPNC_ADVANCED_SETTING_NAME:
		switch key {
		case nm.NM_SETTING_VK_VPN_VPNC_KEY_ENCRYPTION_METHOD:
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

func removeAvailableKeys(keys []string, section, key string) (newKeys []string) {
	newKeys = appendStrArrayUnique(keys)
	relatedVks := getRelatedAvailableVkeys(section, key)
	if len(relatedVks) == 0 {
		newKeys = removeStrArray(newKeys, key)
		return
	}

	newKeys = removeStrArray(newKeys, relatedVks...)
	return
}

func getRelatedAvailableVkeys(section, key string) (vks []string) {
	for _, vk := range virtualKeys {
		if vk.relatedSection == section && isStringInArray(key, vk.relatedKeys) {
			vks = append(vks, vk.value)
		}
	}
	return
}

// get related virtual keys of target key
func getRelatedVkeys(section, key string) (vks []string) {
	for _, vk := range virtualKeys {
		if vk.relatedSection == section && isStringInArray(key, vk.relatedKeys) {
			vks = append(vks, vk.value)
		}
	}
	return
}

// remove virtual key means to remove its related keys
func removeVirtualKey(data connectionData, section, vkey string) {
	if isControllerVkey(section, vkey) {
		// ignore controller virtual key for that there is no related key for it
		return
	}
	if isChildVkey(section, vkey) {
		// ignore child virtual key or its brother keys will be affected
		return
	}
	vkeyInfo, ok := getVkeyInfo(section, vkey)
	if !ok {
		return
	}
	for _, key := range vkeyInfo.relatedKeys {
		removeSettingKey(data, vkeyInfo.relatedSection, key)
	}
}

func isWrapperVkey(section, vkey string) bool {
	vkInfo, ok := getVkeyInfo(section, vkey)
	if !ok {
		return false
	}
	if vkInfo.vkType == vkTypeWrapper {
		return true
	}
	return false
}

func isEnableWrapperVkey(section, vkey string) bool {
	vkInfo, ok := getVkeyInfo(section, vkey)
	if !ok {
		return false
	}
	if vkInfo.vkType == vkTypeEnableWrapper {
		return true
	}
	return false
}

func isControllerVkey(section, vkey string) bool {
	vkInfo, ok := getVkeyInfo(section, vkey)
	if !ok {
		return false
	}
	if vkInfo.vkType == vkTypeController {
		return true
	}
	return false
}

func isChildVkey(section, vkey string) (childKey bool) {
	for _, vk := range virtualKeys {
		if vk.relatedSection == section && vk.value == vkey {
			childKey = vk.childKey
		}
	}
	return
}

func isOptionalVkey(section, vkey string) (optional bool) {
	for _, vk := range virtualKeys {
		if vk.relatedSection == section && vk.value == vkey {
			optional = vk.optional
		}
	}
	return
}

// Controller virtual keys, that without really related section
func logicSetSettingVk8021xEnable(data connectionData, value bool) (err error) {
	if value {
		addSetting(data, nm.NM_SETTING_802_1X_SETTING_NAME)
		err = logicSetSettingVk8021xEap(data, "tls")
	} else {
		removeSetting(data, nm.NM_SETTING_802_1X_SETTING_NAME)
	}
	return
}
func logicSetSettingVk8021xEap(data connectionData, value string) (err error) {
	return logicSetSetting8021xEap(data, []string{value})
}

func getSettingVkMobileCountry(data connectionData) (countryCode string) {
	return getSettingCacheKeyString(data, nm.NM_SETTING_VS_MOBILE, nm.NM_SETTING_VK_MOBILE_COUNTRY)
}
func doLogicSetSettingVkMobileCountry(data connectionData, countryCode string) {
	setSettingCacheKey(data, nm.NM_SETTING_VS_MOBILE, nm.NM_SETTING_VK_MOBILE_COUNTRY, countryCode)
}
func logicSetSettingVkMobileCountry(data connectionData, countryCode string) (err error) {
	logger.Info("set", nm.NM_SETTING_VK_MOBILE_COUNTRY, countryCode)
	doLogicSetSettingVkMobileCountry(data, countryCode)
	defaultProvider, err := mobileprovider.GetDefaultProvider(countryCode)
	if err == nil {
		err = logicSetSettingVkMobileProvider(data, defaultProvider)
	}
	return
}

func getSettingVkMobileProvider(data connectionData) (provider string) {
	return getSettingCacheKeyString(data, nm.NM_SETTING_VS_MOBILE, nm.NM_SETTING_VK_MOBILE_PROVIDER)
}
func doLogicSetSettingVkMobileProvider(data connectionData, provider string) {
	setSettingCacheKey(data, nm.NM_SETTING_VS_MOBILE, nm.NM_SETTING_VK_MOBILE_PROVIDER, provider)
}
func logicSetSettingVkMobileProvider(data connectionData, provider string) (err error) {
	logger.Info("set", nm.NM_SETTING_VK_MOBILE_PROVIDER, provider)
	doLogicSetSettingVkMobileProvider(data, provider)
	if provider != mobileProviderValueCustom {
		defaultPlan, err := mobileprovider.GetDefaultPlan(getSettingVkMobileCountry(data), provider)
		if err == nil {
			err = logicSetSettingVkMobilePlan(data, mobileprovider.MarshalPlan(defaultPlan))
		}
	} else {
		syncMoibleConnectionId(data)
	}
	return
}

func getSettingVkMobilePlan(data connectionData) (planValue string) {
	return getSettingCacheKeyString(data, nm.NM_SETTING_VS_MOBILE, nm.NM_SETTING_VK_MOBILE_PLAN)
}
func doLogicSetSettingVkMobilePlan(data connectionData, planValue string) {
	setSettingCacheKey(data, nm.NM_SETTING_VS_MOBILE, nm.NM_SETTING_VK_MOBILE_PLAN, planValue)
}
func logicSetSettingVkMobilePlan(data connectionData, planValue string) (err error) {
	logger.Info("set", nm.NM_SETTING_VK_MOBILE_PLAN, planValue)
	doLogicSetSettingVkMobilePlan(data, planValue)
	p, err := mobileprovider.UnmarshalPlan(planValue)
	if err != nil {
		logger.Error(err)
		return
	}

	countryCode := getSettingVkMobileCountry(data)
	providerName := getSettingVkMobileProvider(data)
	if p.IsGSM {
		logicSetSettingVkMobileServiceType(data, connectionMobileGsm)
		apn, err := mobileprovider.GetAPN(countryCode, providerName, p.APNValue, p.APNUsageType)
		if err == nil {
			setSettingGsmApn(data, apn.Value)
			if len(apn.Username) > 0 {
				setSettingGsmUsername(data, apn.Username)
			}
			if len(apn.Password) > 0 {
				setSettingGsmPassword(data, apn.Password)
			}
		}
	} else {
		logicSetSettingVkMobileServiceType(data, connectionMobileCdma)
		cdma, err := mobileprovider.GetCDMA(countryCode, providerName)
		if err == nil {
			if len(cdma.Username) > 0 {
				setSettingCdmaUsername(data, cdma.Username)
			}
			if len(cdma.Password) > 0 {
				setSettingCdmaPassword(data, cdma.Password)
			}
		}
	}
	return
}

func getSettingVkMobileServiceType(data connectionData) (serviceType string) {
	if isSettingExists(data, nm.NM_SETTING_GSM_SETTING_NAME) {
		serviceType = connectionMobileGsm
	} else if isSettingExists(data, nm.NM_SETTING_CDMA_SETTING_NAME) {
		serviceType = connectionMobileCdma
	} else {
		logger.Error("get mobile service type failed, neither gsm nor cdma")
	}
	return
}
func logicSetSettingVkMobileServiceType(data connectionData, serviceType string) (err error) {
	// always reset mobile settings
	removeSetting(data, nm.NM_SETTING_GSM_SETTING_NAME)
	removeSetting(data, nm.NM_SETTING_CDMA_SETTING_NAME)
	switch serviceType {
	case connectionMobileGsm:
		initSettingSectionGsm(data)
	case connectionMobileCdma:
		initSettingSectionCdma(data)
	default:
		err = fmt.Errorf("invalid mobile service type %s", serviceType)
	}
	syncMoibleConnectionId(data)
	return
}

func getSettingVkMobileApnReadonly(data connectionData) (value string) {
	return getSettingGsmApn(data)
}
func logicSetSettingVkMobileApnReadonly(data connectionData, value string) (err error) {
	return
}

func getSettingVkVpnType(data connectionData) (vpnType string) {
	vpnType = getCustomConnectionType(data)
	return
}
func logicSetSettingVkVpnType(data connectionData, vpnType string) (err error) {
	removeSetting(data, nm.NM_SETTING_VPN_SETTING_NAME)
	removeSetting(data, nm.NM_SETTING_IP6_CONFIG_SETTING_NAME)
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
	case connectionVpnStrongswan:
		initSettingSectionVpnStrongswan(data)
	case connectionVpnVpnc:
		initSettingSectionVpnVpnc(data)
	default:
		err = fmt.Errorf("invalid vpn type %s", vpnType)
	}
	return
}

func getSettingVkVpnMissingPlugin(data connectionData) (missingPlugin string) {
	// FIXME: the package names only works for debian and ubuntu
	vpnType := getCustomConnectionType(data)
	if !isStringInArray(vpnType, getLocalSupportedVpnTypes()) {
		switch vpnType {
		case connectionVpnL2tp:
			missingPlugin = "network-manager-l2tp-gnome"
		case connectionVpnOpenconnect:
			missingPlugin = "network-manager-openconnect-gnome"
		case connectionVpnOpenvpn:
			missingPlugin = "network-manager-openvpn-gnome"
		case connectionVpnPptp:
			missingPlugin = "network-manager-pptp-gnome"
		case connectionVpnStrongswan:
			missingPlugin = "network-manager-strongswan"
		case connectionVpnVpnc:
			missingPlugin = "network-manager-vpnc-gnome"
		default:
			err := fmt.Errorf("invalid vpn type %s", vpnType)
			logger.Error(err)
		}
	}
	return
}
func logicSetSettingVkVpnMissingPlugin(data connectionData, vpnType string) (err error) {
	return
}
