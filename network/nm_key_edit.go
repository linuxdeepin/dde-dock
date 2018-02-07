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
	"pkg.deepin.io/lib/dbus"
)

func isJSONValueMeansToDeleteKey(valueJSON string, t ktype) (doDelete bool) {
	if valueJSON == jsonNull || valueJSON == jsonEmptyString || valueJSON == jsonEmptyArray {
		return true
	}
	switch t {
	case ktypeIpv6Addresses:
	case ktypeIpv6Routes:
	case ktypeWrapperIpv4Dns:
		if valueJSON == `[""]` {
			doDelete = true
		}
	case ktypeWrapperIpv4Addresses:
		if valueJSON == `[{"Address":"","Mask":"","Gateway":""}]` {
			doDelete = true
		}
	case ktypeWrapperIpv4Routes:
		if valueJSON == `[{"Address":"","Mask":"","NextHop":"","Metric":0}]` {
			doDelete = true
		}
	case ktypeWrapperIpv6Dns:
		if valueJSON == `[""]` {
			doDelete = true
		}
	case ktypeWrapperIpv6Addresses:
		if valueJSON == `[{"Address":"","Prefix":0,"Gateway":""}]` {
			doDelete = true
		}
	case ktypeWrapperIpv6Routes:
		if valueJSON == `[{"Address":"","Prefix":0,"NextHop":"","Metric":0}]` {
			doDelete = true
		}
	}
	return
}

func getSettingKeyJSON(data connectionData, section, key string, t ktype) (valueJSON string) {
	if isWrapperKeyType(t) && !isSettingKeyExists(data, section, key) {
		// if the key is not exists and is a wrapper key, get its
		// default value and marshaled to json directly instead of use
		// keyValueToJSON(), which will dispatch wrapper keys
		// specially
		valueJSON, _ = marshalJSON(generalGetSettingDefaultValue(section, key))
		return
	}

	value := getSettingKey(data, section, key)
	valueJSON, err := keyValueToJSON(value, t)
	if err != nil {
		logger.Error("get connection data failed:", err)
		return
	}

	if len(valueJSON) == 0 {
		logger.Error("getSettingKeyJSON: valueJSON is empty")
	}

	return
}

func getSettingKey(data connectionData, section, key string) (value interface{}) {
	// special for vpn plugin keys
	if isSettingVpnPluginKey(section) {
		return getSettingVpnPluginKey(data, section, key)
	}

	if !isSettingKeyExists(data, section, key) {
		// if key not exists, return the default value
		return generalGetSettingDefaultValue(section, key)
	}

	realSetting := getAliasSettingRealName(section)
	return doGetSettingKey(data, realSetting, key)
}
func doGetSettingKey(data connectionData, section, key string) (value interface{}) {
	sectionData, ok := data[section]
	if !ok {
		logger.Errorf("invalid section: data[%s]", section)
		return
	}
	variant, ok := sectionData[key]
	if !ok {
		// not exists, just return nil
		return
	}

	value = variant.Value()
	// only debug for develop
	// logger.Debugf("getSettingKey: data[%s][%s]=%v", section, key, value)
	if isInterfaceNil(value) {
		// variant exists, but the value is nil, so we give an error
		// message
		logger.Errorf("getSettingKey: data[%s][%s] is nil", section, key)
	}

	return
}

func setSettingKeyJSON(data connectionData, section, key, valueJSON string, t ktype) (kerr error) {
	if len(valueJSON) == 0 {
		logger.Error("setSettingKeyJSON: valueJSON is empty")
		kerr = fmt.Errorf(nmKeyErrorInvalidValue)
		return
	}

	// remove connection data key if valueJSON is null or empty
	if isJSONValueMeansToDeleteKey(valueJSON, t) {
		logger.Debugf("json value means to remove key, data[%s][%s]=%#v", section, key, valueJSON)
		removeSettingKey(data, section, key)
		return
	}

	value, err := jsonToKeyValue(valueJSON, t)
	if err != nil {
		logger.Debugf("set connection data failed, valueJSON=%s, ktype=%s, error message:%v",
			valueJSON, getKtypeDesc(t), err)
		kerr = fmt.Errorf(nmKeyErrorInvalidValue)
		return
	}
	logger.Debugf("setSettingKeyJSON data[%s][%s]=%#v, valueJSON=%s", section, key, value, valueJSON)
	if isInterfaceNil(value) {
		removeSettingKey(data, section, key)
	} else {
		setSettingKey(data, section, key, value)
	}
	return
}

func setSettingKey(data connectionData, section, key string, value interface{}) {
	// special for vpn plugin keys
	if isSettingVpnPluginKey(section) {
		setSettingVpnPluginKey(data, section, key, value)
		return
	}
	realSetting := getAliasSettingRealName(section)
	doSetSettingKey(data, realSetting, key, value)
}
func doSetSettingKey(data connectionData, section, key string, value interface{}) {
	var sectionData map[string]dbus.Variant
	sectionData, ok := data[section]
	if !ok {
		logger.Errorf(`set connection data failed, section "%s" is not exits yet`, section)
		return
	}
	sectionData[key] = dbus.MakeVariant(value)
	logger.Debugf("setSettingKey: data[%s][%s]=%#v", section, key, value)
}

func removeSettingKey(data connectionData, section string, keys ...string) {
	logger.Debugf("removeSettingKey data[%s], %s", section, keys)

	// special for vpn plugin keys
	if isSettingVpnPluginKey(section) {
		removeSettingVpnPluginKey(data, section, keys...)
		return
	}

	realSetting := getAliasSettingRealName(section)
	sectionData, ok := data[realSetting]
	if !ok {
		return
	}

	for _, k := range keys {
		delete(sectionData, k)
	}
}

func removeSettingKeyBut(data connectionData, section string, keys ...string) {
	// special for vpn plugin keys
	if isSettingVpnPluginKey(section) {
		removeSettingVpnPluginKeyBut(data, section, keys...)
		return
	}

	realSetting := getAliasSettingRealName(section)
	sectionData, ok := data[realSetting]
	if !ok {
		return
	}

	for k := range sectionData {
		if !isStringInArray(k, keys) {
			delete(sectionData, k)
		}
	}
}

func isSettingKeyExists(data connectionData, section, key string) bool {
	// special for vpn plugin keys
	if isSettingVpnPluginKey(section) {
		return isSettingVpnPluginKeyExists(data, section, key)
	}

	realSetting := getAliasSettingRealName(section)
	sectionData, ok := data[realSetting]
	if !ok {
		return false
	}

	_, ok = sectionData[key]
	if !ok {
		return false
	}

	return true
}

func addSetting(data connectionData, setting string) {
	realSetting := getAliasSettingRealName(setting)
	var settingData map[string]dbus.Variant
	settingData, ok := data[realSetting]
	if !ok {
		// add setting if not exists
		settingData = make(map[string]dbus.Variant)
		data[realSetting] = settingData
	}
}

func removeSetting(data connectionData, setting string) {
	realSetting := getAliasSettingRealName(setting)
	_, ok := data[realSetting]
	if ok {
		// remove setting if exists
		delete(data, realSetting)
	}
}

func isSettingExists(data connectionData, setting string) bool {
	realSetting := getAliasSettingRealName(setting)
	_, ok := data[realSetting]
	return ok
}

func generalSetSettingAutoconnect(data connectionData, autoConnect bool) {
	switch getSettingConnectionType(data) {
	case nm.NM_SETTING_VPN_SETTING_NAME:
		uuid := getSettingConnectionUuid(data)
		manager.config.setVpnConnectionAutoConnect(uuid, autoConnect)
	default:
		setSettingConnectionAutoconnect(data, autoConnect)
	}
}

// operator for cache section
func getSettingCacheKey(data connectionData, section, key string) (value interface{}) {
	return doGetSettingKey(data, sectionCache, section+"/"+key)
}
func getSettingCacheKeyString(data connectionData, section, key string) (value string) {
	return interfaceToString(getSettingCacheKey(data, section, key))
}
func setSettingCacheKey(data connectionData, section, key string, value interface{}) {
	doSetSettingKey(data, sectionCache, section+"/"+key, value)
}
func fillSectionCache(data connectionData) {
	addSetting(data, sectionCache)

	// ip4
	if isSettingExists(data, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME) {
		dnses := getSettingIP4ConfigDns(data)
		switch len(dnses) {
		case 0:
			logicSetSettingVkIp4ConfigDns(data, "")
			logicSetSettingVkIp4ConfigDns2(data, "")
		case 1:
			logicSetSettingVkIp4ConfigDns(data, convertIpv4AddressToString(dnses[0]))
			logicSetSettingVkIp4ConfigDns2(data, "")
		default:
			logicSetSettingVkIp4ConfigDns(data, convertIpv4AddressToString(dnses[0]))
			logicSetSettingVkIp4ConfigDns2(data, convertIpv4AddressToString(dnses[1]))
		}
	}

	// ip6
	if isSettingExists(data, nm.NM_SETTING_IP6_CONFIG_SETTING_NAME) {
		dnses := getSettingIP6ConfigDns(data)
		switch len(dnses) {
		case 0:
			logicSetSettingVkIp6ConfigDns(data, "")
			logicSetSettingVkIp6ConfigDns2(data, "")
		case 1:
			logicSetSettingVkIp6ConfigDns(data, convertIpv6AddressToString(dnses[0]))
			logicSetSettingVkIp6ConfigDns2(data, "")
		default:
			logicSetSettingVkIp6ConfigDns(data, convertIpv6AddressToString(dnses[0]))
			logicSetSettingVkIp6ConfigDns2(data, convertIpv6AddressToString(dnses[1]))
		}
	}

	// mobile
	switch getCustomConnectionType(data) {
	case connectionMobileGsm, connectionMobileCdma:
		uuid := getSettingConnectionUuid(data)
		manager.config.ensureMobileConfigExists(uuid)
		doLogicSetSettingVkMobileCountry(data, manager.config.getMobileConnectionCountry(uuid))
		doLogicSetSettingVkMobileProvider(data, manager.config.getMobileConnectionProvider(uuid))
		doLogicSetSettingVkMobilePlan(data, manager.config.getMobileConnectionPlan(uuid))
	}
}
func refileSectionCache(data connectionData) {
	// ip4
	if isSettingExists(data, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME) {
		dnses := make([]uint32, 0)
		dns1Str := getSettingVkIp4ConfigDns(data)
		dns2Str := getSettingVkIp4ConfigDns2(data)
		if dns, err := convertIpv4AddressToUint32Check(dns1Str); err == nil {
			dnses = append(dnses, dns)
		}
		if dns, err := convertIpv4AddressToUint32Check(dns2Str); err == nil {
			dnses = append(dnses, dns)
		}
		if len(dnses) == 0 {
			removeSettingIP4ConfigDns(data)
		} else {
			setSettingIP4ConfigDns(data, dnses)
		}
	}

	// ip6
	if isSettingExists(data, nm.NM_SETTING_IP6_CONFIG_SETTING_NAME) {
		dnses := make([][]byte, 0)
		dns1Str := getSettingVkIp6ConfigDns(data)
		dns2Str := getSettingVkIp6ConfigDns2(data)
		if dns, err := convertIpv6AddressToArrayByteCheck(dns1Str); err == nil {
			dnses = append(dnses, dns)
		}
		if dns, err := convertIpv6AddressToArrayByteCheck(dns2Str); err == nil {
			dnses = append(dnses, dns)
		}
		if len(dnses) == 0 {
			removeSettingIP6ConfigDns(data)
		} else {
			setSettingIP6ConfigDns(data, dnses)
		}
	}

	// mobile
	switch getCustomConnectionType(data) {
	case connectionMobileGsm, connectionMobileCdma:
		uuid := getSettingConnectionUuid(data)
		manager.config.ensureMobileConfigExists(uuid)
		manager.config.setMobileConnectionCountry(uuid, getSettingVkMobileCountry(data))
		manager.config.setMobileConnectionProvider(uuid, getSettingVkMobileProvider(data))
		manager.config.setMobileConnectionPlan(uuid, getSettingVkMobilePlan(data))
	}
	removeSetting(data, sectionCache)
}
