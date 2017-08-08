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
	"pkg.deepin.io/dde/daemon/network/nm"
	. "pkg.deepin.io/lib/gettext"
)

func initSettingSectionIpv4(data connectionData) {
	addSetting(data, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME)
	setSettingIP4ConfigMethod(data, nm.NM_SETTING_IP4_CONFIG_METHOD_AUTO)
	setSettingIP4ConfigNeverDefault(data, false)
}

// Initialize available values
var availableValuesIp4ConfigMethod = make(availableValues)

func initAvailableValuesIp4() {
	availableValuesIp4ConfigMethod[nm.NM_SETTING_IP4_CONFIG_METHOD_AUTO] = kvalue{nm.NM_SETTING_IP4_CONFIG_METHOD_AUTO, Tr("Auto")}
	availableValuesIp4ConfigMethod[nm.NM_SETTING_IP4_CONFIG_METHOD_LINK_LOCAL] = kvalue{nm.NM_SETTING_IP4_CONFIG_METHOD_LINK_LOCAL, Tr("Link-Local Only")}
	availableValuesIp4ConfigMethod[nm.NM_SETTING_IP4_CONFIG_METHOD_MANUAL] = kvalue{nm.NM_SETTING_IP4_CONFIG_METHOD_MANUAL, Tr("Manual")}
	availableValuesIp4ConfigMethod[nm.NM_SETTING_IP4_CONFIG_METHOD_SHARED] = kvalue{nm.NM_SETTING_IP4_CONFIG_METHOD_SHARED, Tr("Shared")}
	availableValuesIp4ConfigMethod[nm.NM_SETTING_IP4_CONFIG_METHOD_DISABLED] = kvalue{nm.NM_SETTING_IP4_CONFIG_METHOD_DISABLED, Tr("Disabled")}
}

// Get available keys
func getSettingIP4ConfigAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME, nm.NM_SETTING_IP_CONFIG_METHOD)
	method := getSettingIP4ConfigMethod(data)
	switch method {
	default:
		logger.Error("ip4 config method is invalid:", method)
	case nm.NM_SETTING_IP4_CONFIG_METHOD_AUTO:
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME, nm.NM_SETTING_IP_CONFIG_DNS)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME, nm.NM_SETTING_IP_CONFIG_NEVER_DEFAULT)
	case nm.NM_SETTING_IP4_CONFIG_METHOD_LINK_LOCAL: // ignore
	case nm.NM_SETTING_IP4_CONFIG_METHOD_MANUAL:
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME, nm.NM_SETTING_IP_CONFIG_DNS)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME, nm.NM_SETTING_IP_CONFIG_ADDRESSES)
		keys = appendAvailableKeys(data, keys, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME, nm.NM_SETTING_IP_CONFIG_NEVER_DEFAULT)
	case nm.NM_SETTING_IP4_CONFIG_METHOD_SHARED:
	case nm.NM_SETTING_IP4_CONFIG_METHOD_DISABLED:
	}
	return
}

// Get available values
func getSettingIP4ConfigAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case nm.NM_SETTING_IP_CONFIG_METHOD:
		if getSettingConnectionType(data) == nm.NM_SETTING_VPN_SETTING_NAME {
			values = []kvalue{
				availableValuesIp4ConfigMethod[nm.NM_SETTING_IP4_CONFIG_METHOD_AUTO],
			}
		} else {
			values = []kvalue{
				availableValuesIp4ConfigMethod[nm.NM_SETTING_IP4_CONFIG_METHOD_AUTO],
				availableValuesIp4ConfigMethod[nm.NM_SETTING_IP4_CONFIG_METHOD_MANUAL],
			}
		}
	}
	return
}

// Check whether the values are correct
func checkSettingIP4ConfigValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)

	// check method
	ensureSettingIP4ConfigMethodNoEmpty(data, errs)
	switch getSettingIP4ConfigMethod(data) {
	default:
		rememberError(errs, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME, nm.NM_SETTING_IP_CONFIG_METHOD, nmKeyErrorInvalidValue)
		return
	case nm.NM_SETTING_IP4_CONFIG_METHOD_AUTO:
	case nm.NM_SETTING_IP4_CONFIG_METHOD_LINK_LOCAL: // ignore
		checkSettingIP4MethodConflict(data, errs)
	case nm.NM_SETTING_IP4_CONFIG_METHOD_MANUAL:
		ensureSettingIP4ConfigAddressesNoEmpty(data, errs)
	case nm.NM_SETTING_IP4_CONFIG_METHOD_SHARED:
		checkSettingIP4MethodConflict(data, errs)
	case nm.NM_SETTING_IP4_CONFIG_METHOD_DISABLED: // ignore
		checkSettingIP4MethodConflict(data, errs)
	}

	// check value of dns
	checkSettingIP4ConfigDns(data, errs)

	// check value of address
	checkSettingIP4ConfigAddresses(data, errs)

	return
}
func checkSettingIP4MethodConflict(data connectionData, errs sectionErrors) {
	// check dns
	if isSettingIP4ConfigDnsExists(data) && len(getSettingIP4ConfigDns(data)) > 0 {
		rememberError(errs, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME, nm.NM_SETTING_IP_CONFIG_DNS, fmt.Sprintf(nmKeyErrorIp4MethodConflict, nm.NM_SETTING_IP_CONFIG_DNS))
	}
	// check dns search
	if isSettingIP4ConfigDnsSearchExists(data) && len(getSettingIP4ConfigDnsSearch(data)) > 0 {
		rememberError(errs, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME, nm.NM_SETTING_IP_CONFIG_DNS_SEARCH, fmt.Sprintf(nmKeyErrorIp4MethodConflict, nm.NM_SETTING_IP_CONFIG_DNS_SEARCH))
	}
	// check address
	if isSettingIP4ConfigAddressesExists(data) && len(getSettingIP4ConfigAddresses(data)) > 0 {
		rememberError(errs, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME, nm.NM_SETTING_IP_CONFIG_ADDRESSES, fmt.Sprintf(nmKeyErrorIp4MethodConflict, nm.NM_SETTING_IP_CONFIG_ADDRESSES))
	}
	// check route
	if isSettingIP4ConfigRoutesExists(data) && len(getSettingIP4ConfigRoutes(data)) > 0 {
		rememberError(errs, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME, nm.NM_SETTING_IP_CONFIG_ROUTES, fmt.Sprintf(nmKeyErrorIp4MethodConflict, nm.NM_SETTING_IP_CONFIG_ROUTES))
	}
}
func checkSettingIP4ConfigDns(data connectionData, errs sectionErrors) {
	if !isSettingIP4ConfigDnsExists(data) {
		return
	}
	dnses := getSettingIP4ConfigDns(data)
	for _, dns := range dnses {
		if dns == 0 {
			rememberError(errs, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME, nm.NM_SETTING_IP_CONFIG_DNS, nmKeyErrorInvalidValue)
			return
		}
	}
}
func checkSettingIP4ConfigAddresses(data connectionData, errs sectionErrors) {
	if !isSettingIP4ConfigAddressesExists(data) {
		return
	}
	addresses := getSettingIP4ConfigAddresses(data)
	for _, addr := range addresses {
		// check address struct
		if len(addr) != 3 {
			rememberError(errs, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME, nm.NM_SETTING_IP_CONFIG_ADDRESSES, nmKeyErrorIp4AddressesStruct)
		}
		// check address
		if addr[0] == 0 {
			rememberError(errs, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME, nm.NM_SETTING_VK_IP4_CONFIG_ADDRESSES_ADDRESS, nmKeyErrorInvalidValue)
		}
		// check prefix
		if addr[1] < 1 || addr[1] > 32 {
			rememberError(errs, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME, nm.NM_SETTING_VK_IP4_CONFIG_ADDRESSES_MASK, nmKeyErrorInvalidValue)
		}
	}
}

// Logic setter
func logicSetSettingIP4ConfigMethod(data connectionData, value string) (err error) {
	// just ignore error here and set value directly, error will be
	// check in checkSettingXXXValues()
	switch value {
	case nm.NM_SETTING_IP4_CONFIG_METHOD_AUTO:
		removeSettingIP4ConfigAddresses(data)
	case nm.NM_SETTING_IP4_CONFIG_METHOD_LINK_LOCAL: // ignore
		removeSettingIP4ConfigDns(data)
		removeSettingIP4ConfigDnsSearch(data)
		removeSettingIP4ConfigAddresses(data)
		removeSettingIP4ConfigRoutes(data)
	case nm.NM_SETTING_IP4_CONFIG_METHOD_MANUAL:
	case nm.NM_SETTING_IP4_CONFIG_METHOD_SHARED:
		removeSettingIP4ConfigDns(data)
		removeSettingIP4ConfigDnsSearch(data)
		removeSettingIP4ConfigAddresses(data)
		removeSettingIP4ConfigRoutes(data)
	case nm.NM_SETTING_IP4_CONFIG_METHOD_DISABLED: // ignore
		removeSettingIP4ConfigDns(data)
		removeSettingIP4ConfigDnsSearch(data)
		removeSettingIP4ConfigAddresses(data)
		removeSettingIP4ConfigRoutes(data)
	}
	setSettingIP4ConfigMethod(data, value)
	return
}

// Virtual key utility
func isSettingIP4ConfigAddressesEmpty(data connectionData) bool {
	addresses := getSettingIP4ConfigAddresses(data)
	if len(addresses) == 0 {
		return true
	}
	if len(addresses[0]) != 3 {
		return true
	}
	return false
}
func getOrNewSettingIP4ConfigAddresses(data connectionData) (addresses [][]uint32) {
	if !isSettingIP4ConfigAddressesEmpty(data) {
		addresses = getSettingIP4ConfigAddresses(data)
	} else {
		addresses = make([][]uint32, 1)
		addresses[0] = make([]uint32, 3)
	}
	return
}

// Virtual key getter
func getSettingVkIp4ConfigDns(data connectionData) (value string) {
	return getSettingCacheKeyString(data, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME, nm.NM_SETTING_VK_IP4_CONFIG_DNS)
}
func getSettingVkIp4ConfigDns2(data connectionData) (value string) {
	return getSettingCacheKeyString(data, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME, nm.NM_SETTING_VK_IP4_CONFIG_DNS2)
}
func getSettingVkIp4ConfigAddressesAddress(data connectionData) (value string) {
	if isSettingIP4ConfigAddressesEmpty(data) {
		return
	}
	addresses := getSettingIP4ConfigAddresses(data)
	value = convertIpv4AddressToString(addresses[0][0])
	return
}
func getSettingVkIp4ConfigAddressesMask(data connectionData) (value string) {
	if isSettingIP4ConfigAddressesEmpty(data) {
		return
	}
	addresses := getSettingIP4ConfigAddresses(data)
	value = convertIpv4PrefixToNetMask(addresses[0][1])
	return
}
func getSettingVkIp4ConfigAddressesGateway(data connectionData) (value string) {
	if isSettingIP4ConfigAddressesEmpty(data) {
		return
	}
	addresses := getSettingIP4ConfigAddresses(data)
	value = convertIpv4AddressToStringNoZero(addresses[0][2])
	return
}
func getSettingVkIp4ConfigRoutesAddress(data connectionData) (value string) {
	// TODO value := getSettingIP4ConfigRoutesAddress(data)
	return
}
func getSettingVkIp4ConfigRoutesMask(data connectionData) (value string) {
	// TODO value := getSettingIP4ConfigRoutesMask(data)
	return
}
func getSettingVkIp4ConfigRoutesNexthop(data connectionData) (value string) {
	// TODO value := getSettingIP4ConfigRoutesNexthop(data)
	return
}
func getSettingVkIp4ConfigRoutesMetric(data connectionData) (value string) {
	// TODO value := getSettingIP4ConfigRoutesMetric(data)
	return
}

// Virtual key logic setter
func logicSetSettingVkIp4ConfigDns(data connectionData, value string) (err error) {
	setSettingCacheKey(data, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME, nm.NM_SETTING_VK_IP4_CONFIG_DNS, value)
	if len(value) > 0 {
		if _, errWrap := convertIpv4AddressToUint32Check(value); errWrap != nil {
			err = fmt.Errorf(nmKeyErrorInvalidValue)
		}
	}
	return
}
func logicSetSettingVkIp4ConfigDns2(data connectionData, value string) (err error) {
	setSettingCacheKey(data, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME, nm.NM_SETTING_VK_IP4_CONFIG_DNS2, value)
	if len(value) > 0 {
		if _, errWrap := convertIpv4AddressToUint32Check(value); errWrap != nil {
			err = fmt.Errorf(nmKeyErrorInvalidValue)
		}
	}
	return
}
func logicSetSettingVkIp4ConfigAddressesAddress(data connectionData, value string) (err error) {
	if len(value) == 0 {
		value = ipv4Zero
	}
	tmpn, err := convertIpv4AddressToUint32Check(value)
	if err != nil {
		err = fmt.Errorf(nmKeyErrorInvalidValue)
		return
	}
	addresses := getOrNewSettingIP4ConfigAddresses(data)
	addr := addresses[0]
	addr[0] = tmpn
	if !isUint32ArrayEmpty(addr) {
		setSettingIP4ConfigAddresses(data, addresses)
	} else {
		removeSettingIP4ConfigAddresses(data)
	}
	return
}
func logicSetSettingVkIp4ConfigAddressesMask(data connectionData, value string) (err error) {
	if len(value) == 0 {
		value = ipv4Zero
	}
	tmpn, err := convertIpv4NetMaskToPrefixCheck(value)
	if err != nil {
		err = fmt.Errorf(nmKeyErrorInvalidValue)
		return
	}
	addresses := getOrNewSettingIP4ConfigAddresses(data)
	addr := addresses[0]
	addr[1] = tmpn
	if !isUint32ArrayEmpty(addr) {
		setSettingIP4ConfigAddresses(data, addresses)
	} else {
		removeSettingIP4ConfigAddresses(data)
	}
	return
}
func logicSetSettingVkIp4ConfigAddressesGateway(data connectionData, value string) (err error) {
	if len(value) == 0 {
		value = ipv4Zero
	}
	tmpn, err := convertIpv4AddressToUint32Check(value)
	if err != nil {
		err = fmt.Errorf(nmKeyErrorInvalidValue)
		return
	}
	addresses := getOrNewSettingIP4ConfigAddresses(data)
	addr := addresses[0]
	addr[2] = tmpn
	if !isUint32ArrayEmpty(addr) {
		setSettingIP4ConfigAddresses(data, addresses)
	} else {
		removeSettingIP4ConfigAddresses(data)
	}
	return
}
func logicSetSettingVkIp4ConfigRoutesAddress(data connectionData, value string) (err error) {
	// TODO setSettingIP4ConfigRoutesAddressJSON(data)
	return
}
func logicSetSettingVkIp4ConfigRoutesMask(data connectionData, value string) (err error) {
	// TODO setSettingIP4ConfigRoutesMaskJSON(data)
	return
}
func logicSetSettingVkIp4ConfigRoutesNexthop(data connectionData, value string) (err error) {
	// TODO setSettingIP4ConfigRoutesNexthopJSON(data)
	return
}
func logicSetSettingVkIp4ConfigRoutesMetric(data connectionData, value string) (err error) {
	// TODO setSettingIP4ConfigRoutesMetricJSON(data)
	return
}
