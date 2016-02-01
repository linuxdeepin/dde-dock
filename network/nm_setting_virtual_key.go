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
	available      bool // check if is used by front-end
	childKey       bool // such as ip address, mask and gateway
	optional       bool // if key is optional(such as child key gateway of ip address), will ignore error for it
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
	case NM_SETTING_VS_SECURITY:
		if getCustomConnectionType(data) == connectionWired {
			keys = []string{NM_SETTING_VK_802_1X_ENABLE}
		}
	case NM_SETTING_VS_MOBILE:
		keys = []string{
			NM_SETTING_VK_MOBILE_COUNTRY,
			NM_SETTING_VK_MOBILE_PROVIDER,
		}
		if getSettingVkMobileProvider(data) == mobileProviderValueCustom {
			keys = append(keys, NM_SETTING_VK_MOBILE_SERVICE_TYPE)
		} else {
			keys = append(keys, NM_SETTING_VK_MOBILE_PLAN)
			// TODO: is apn-readonly widget necessary?
			// if getSettingVkMobileServiceType(data) == connectionMobileGsm {
			// keys = append(keys, NM_SETTING_VK_MOBILE_APN_READONLY)
			// }
		}
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
		case NM_SETTING_VK_MOBILE_COUNTRY:
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
		case NM_SETTING_VK_MOBILE_PROVIDER:
			countryCode := getSettingVkMobileCountry(data)
			names, _ := mobileprovider.GetProviderNames(countryCode)
			for _, name := range names {
				values = append(values, kvalue{name, name})
			}
			values = append(values, kvalue{mobileProviderValueCustom, Tr("Custom")})
		case NM_SETTING_VK_MOBILE_PLAN:
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
		case NM_SETTING_VK_MOBILE_SERVICE_TYPE:
			values = []kvalue{
				kvalue{connectionMobileGsm, Tr("GSM (GPRS, UMTS)")},
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
		if vk.relatedSection == section && isStringInArray(key, vk.relatedKeys) && vk.available {
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

func getSettingVkMobileCountry(data connectionData) (countryCode string) {
	ivalue := getSettingCacheKey(data, NM_SETTING_VK_MOBILE_COUNTRY)
	countryCode = interfaceToString(ivalue)
	return
}
func logicSetSettingVkMobileCountry(data connectionData, countryCode string) (err error) {
	logger.Info("set", NM_SETTING_VK_MOBILE_COUNTRY, countryCode)
	setSettingCacheKey(data, NM_SETTING_VK_MOBILE_COUNTRY, countryCode)
	defaultProvider, err := mobileprovider.GetDefaultProvider(countryCode)
	if err == nil {
		err = logicSetSettingVkMobileProvider(data, defaultProvider)
	}
	return
}

func getSettingVkMobileProvider(data connectionData) (provider string) {
	ivalue := getSettingCacheKey(data, NM_SETTING_VK_MOBILE_PROVIDER)
	provider = interfaceToString(ivalue)
	return
}
func logicSetSettingVkMobileProvider(data connectionData, provider string) (err error) {
	logger.Info("set", NM_SETTING_VK_MOBILE_PROVIDER, provider)
	setSettingCacheKey(data, NM_SETTING_VK_MOBILE_PROVIDER, provider)
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
	ivalue := getSettingCacheKey(data, NM_SETTING_VK_MOBILE_PLAN)
	planValue = interfaceToString(ivalue)
	return
}
func logicSetSettingVkMobilePlan(data connectionData, planValue string) (err error) {
	logger.Info("set", NM_SETTING_VK_MOBILE_PLAN, planValue)
	setSettingCacheKey(data, NM_SETTING_VK_MOBILE_PLAN, planValue)
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
	if isSettingSectionExists(data, sectionGsm) {
		serviceType = connectionMobileGsm
	} else if isSettingSectionExists(data, sectionCdma) {
		serviceType = connectionMobileCdma
	} else {
		logger.Error("get mobile service type failed, neither gsm nor cdma")
	}
	return
}
func logicSetSettingVkMobileServiceType(data connectionData, serviceType string) (err error) {
	// always reset mobile settings
	removeSettingSection(data, sectionGsm)
	removeSettingSection(data, sectionCdma)
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
			// TODO:
			// fmt.Errorf("invalid vpn type %s", vpnType)
		}
	}
	return
}
func logicSetSettingVkVpnMissingPlugin(data connectionData, vpnType string) (err error) {
	return
}
