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
)

func initSettingSectionIpv6(data connectionData) {
	addSettingSection(data, sectionIpv6)
	setSettingIp6ConfigMethod(data, NM_SETTING_IP6_CONFIG_METHOD_AUTO)
}

// Initialize available values
var availableValuesIp6ConfigMethod = make(availableValues)

func initAvailableValuesIp6() {
	availableValuesIp6ConfigMethod[NM_SETTING_IP6_CONFIG_METHOD_IGNORE] = kvalue{NM_SETTING_IP6_CONFIG_METHOD_IGNORE, Tr("Ignore")}
	availableValuesIp6ConfigMethod[NM_SETTING_IP6_CONFIG_METHOD_AUTO] = kvalue{NM_SETTING_IP6_CONFIG_METHOD_AUTO, Tr("Auto")}
	availableValuesIp6ConfigMethod[NM_SETTING_IP6_CONFIG_METHOD_DHCP] = kvalue{NM_SETTING_IP6_CONFIG_METHOD_DHCP, Tr("DHCP")}
	availableValuesIp6ConfigMethod[NM_SETTING_IP6_CONFIG_METHOD_LINK_LOCAL] = kvalue{NM_SETTING_IP6_CONFIG_METHOD_LINK_LOCAL, Tr("Link-Local Only")}
	availableValuesIp6ConfigMethod[NM_SETTING_IP6_CONFIG_METHOD_MANUAL] = kvalue{NM_SETTING_IP6_CONFIG_METHOD_MANUAL, Tr("Manual")}
	availableValuesIp6ConfigMethod[NM_SETTING_IP6_CONFIG_METHOD_SHARED] = kvalue{NM_SETTING_IP6_CONFIG_METHOD_SHARED, Tr("Shared")}
}

// Get available keys
func getSettingIp6ConfigAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, sectionIpv6, NM_SETTING_IP_CONFIG_METHOD)
	method := getSettingIp6ConfigMethod(data)
	switch method {
	default:
		logger.Error("ip6 config method is invalid:", method)
	case NM_SETTING_IP6_CONFIG_METHOD_IGNORE:
	case NM_SETTING_IP6_CONFIG_METHOD_AUTO:
		keys = appendAvailableKeys(data, keys, sectionIpv6, NM_SETTING_IP_CONFIG_DNS)
	case NM_SETTING_IP6_CONFIG_METHOD_DHCP: // ignore
		keys = appendAvailableKeys(data, keys, sectionIpv6, NM_SETTING_IP_CONFIG_DNS)
	case NM_SETTING_IP6_CONFIG_METHOD_LINK_LOCAL: // ignore
	case NM_SETTING_IP6_CONFIG_METHOD_MANUAL:
		keys = appendAvailableKeys(data, keys, sectionIpv6, NM_SETTING_IP_CONFIG_DNS)
		keys = appendAvailableKeys(data, keys, sectionIpv6, NM_SETTING_IP_CONFIG_ADDRESSES)
	case NM_SETTING_IP6_CONFIG_METHOD_SHARED:
	}
	return
}

// Get available values
func getSettingIp6ConfigAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case NM_SETTING_IP_CONFIG_METHOD:
		// values = []string{
		// 	// NM_SETTING_IP6_CONFIG_METHOD_IGNORE, // ignore
		// 	NM_SETTING_IP6_CONFIG_METHOD_AUTO,
		// 	// NM_SETTING_IP6_CONFIG_METHOD_DHCP, // ignore
		// 	// NM_SETTING_IP6_CONFIG_METHOD_LINK_LOCAL, // ignore
		// 	NM_SETTING_IP6_CONFIG_METHOD_MANUAL,
		// 	// NM_SETTING_IP6_CONFIG_METHOD_SHARED,// ignore
		// }
		if getSettingConnectionType(data) == NM_SETTING_VPN_SETTING_NAME {
			values = []kvalue{
				availableValuesIp6ConfigMethod[NM_SETTING_IP6_CONFIG_METHOD_AUTO],
				availableValuesIp6ConfigMethod[NM_SETTING_IP6_CONFIG_METHOD_IGNORE],
			}
		} else {
			values = []kvalue{
				availableValuesIp6ConfigMethod[NM_SETTING_IP6_CONFIG_METHOD_AUTO],
				availableValuesIp6ConfigMethod[NM_SETTING_IP6_CONFIG_METHOD_MANUAL],
				availableValuesIp6ConfigMethod[NM_SETTING_IP6_CONFIG_METHOD_IGNORE],
			}
		}
	}
	return
}

// Check whether the values are correct
func checkSettingIp6ConfigValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)

	// check method
	ensureSettingIp6ConfigMethodNoEmpty(data, errs)
	switch getSettingIp6ConfigMethod(data) {
	default:
		rememberError(errs, sectionIpv6, NM_SETTING_IP_CONFIG_METHOD, NM_KEY_ERROR_INVALID_VALUE)
		return
	case NM_SETTING_IP6_CONFIG_METHOD_IGNORE:
		checkSettingIp6MethodConflict(data, errs)
	case NM_SETTING_IP6_CONFIG_METHOD_AUTO:
	case NM_SETTING_IP6_CONFIG_METHOD_DHCP: // ignore
	case NM_SETTING_IP6_CONFIG_METHOD_LINK_LOCAL: // ignore
		checkSettingIp6MethodConflict(data, errs)
	case NM_SETTING_IP6_CONFIG_METHOD_MANUAL:
		ensureSettingIp6ConfigAddressesNoEmpty(data, errs)
	case NM_SETTING_IP6_CONFIG_METHOD_SHARED:
		checkSettingIp6MethodConflict(data, errs)
	}

	// check value of dns
	checkSettingIp6ConfigDns(data, errs)

	// check value of address
	checkSettingIp6ConfigAddresses(data, errs)

	// TODO check value of route

	return
}

func checkSettingIp6MethodConflict(data connectionData, errs sectionErrors) {
	// check dns
	if isSettingIp6ConfigDnsExists(data) && len(getSettingIp6ConfigDns(data)) > 0 {
		rememberError(errs, sectionIpv6, NM_SETTING_IP_CONFIG_DNS, fmt.Sprintf(NM_KEY_ERROR_IP6_METHOD_CONFLICT, NM_SETTING_IP_CONFIG_DNS))
	}
	// check dns search
	if isSettingIp6ConfigDnsSearchExists(data) && len(getSettingIp6ConfigDnsSearch(data)) > 0 {
		rememberError(errs, sectionIpv6, NM_SETTING_IP_CONFIG_DNS_SEARCH, fmt.Sprintf(NM_KEY_ERROR_IP6_METHOD_CONFLICT, NM_SETTING_IP_CONFIG_DNS_SEARCH))
	}
	// check address
	if isSettingIp6ConfigAddressesExists(data) && len(getSettingIp6ConfigAddresses(data)) > 0 {
		rememberError(errs, sectionIpv6, NM_SETTING_IP_CONFIG_ADDRESSES, fmt.Sprintf(NM_KEY_ERROR_IP6_METHOD_CONFLICT, NM_SETTING_IP_CONFIG_ADDRESSES))
	}
	// check route
	if isSettingIp6ConfigRoutesExists(data) && len(getSettingIp6ConfigRoutes(data)) > 0 {
		rememberError(errs, sectionIpv6, NM_SETTING_IP_CONFIG_ROUTES, fmt.Sprintf(NM_KEY_ERROR_IP6_METHOD_CONFLICT, NM_SETTING_IP_CONFIG_ROUTES))
	}
}

func checkSettingIp6ConfigDns(data connectionData, errs sectionErrors) {
	if !isSettingIp6ConfigDnsExists(data) {
		return
	}
	dnses := getSettingIp6ConfigDns(data)
	for _, dns := range dnses {
		if !isIpv6AddressValid(dns) {
			rememberError(errs, sectionIpv6, NM_SETTING_VK_IP6_CONFIG_DNS, NM_KEY_ERROR_INVALID_VALUE)
			return
		}
		if isIpv6AddressZero(dns) {
			rememberError(errs, sectionIpv6, NM_SETTING_VK_IP6_CONFIG_DNS, NM_KEY_ERROR_INVALID_VALUE)
			return
		}
	}
}

func checkSettingIp6ConfigAddresses(data connectionData, errs sectionErrors) {
	if !isSettingIp6ConfigAddressesExists(data) {
		return
	}
	addresses := getSettingIp6ConfigAddresses(data)
	for _, addr := range addresses {
		// check address
		if !isIpv6AddressValid(addr.Address) {
			rememberError(errs, sectionIpv6, NM_SETTING_VK_IP6_CONFIG_ADDRESSES_ADDRESS, NM_KEY_ERROR_INVALID_VALUE)
			// TODO test
			logger.Warning(NM_KEY_ERROR_INVALID_VALUE, addr.Address)
		}
		if isIpv6AddressZero(addr.Address) {
			rememberError(errs, sectionIpv6, NM_SETTING_VK_IP6_CONFIG_ADDRESSES_ADDRESS, NM_KEY_ERROR_INVALID_VALUE)
		}
		// check prefix
		if addr.Prefix < 1 || addr.Prefix > 128 {
			rememberError(errs, sectionIpv6, NM_SETTING_VK_IP6_CONFIG_ADDRESSES_PREFIX, NM_KEY_ERROR_INVALID_VALUE)
		}
		// check gateway
		if !isIpv6AddressValid(addr.Gateway) {
			rememberError(errs, sectionIpv6, NM_SETTING_VK_IP6_CONFIG_ADDRESSES_GATEWAY, NM_KEY_ERROR_INVALID_VALUE)
		}
	}
}

// Logic setter
func logicSetSettingIp6ConfigMethod(data connectionData, value string) (err error) {
	switch value {
	case NM_SETTING_IP6_CONFIG_METHOD_IGNORE:
		removeSettingKeyBut(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP_CONFIG_METHOD)
	case NM_SETTING_IP6_CONFIG_METHOD_AUTO:
		removeSettingIp6ConfigAddresses(data)
	case NM_SETTING_IP6_CONFIG_METHOD_DHCP: // ignore
	case NM_SETTING_IP6_CONFIG_METHOD_LINK_LOCAL: // ignore
		removeSettingIp6ConfigDns(data)
		removeSettingIp6ConfigDnsSearch(data)
		removeSettingIp6ConfigAddresses(data)
		removeSettingIp6ConfigRoutes(data)
	case NM_SETTING_IP6_CONFIG_METHOD_MANUAL:
	case NM_SETTING_IP6_CONFIG_METHOD_SHARED:
		removeSettingIp6ConfigDns(data)
		removeSettingIp6ConfigDnsSearch(data)
		removeSettingIp6ConfigAddresses(data)
		removeSettingIp6ConfigRoutes(data)
	}
	setSettingIp6ConfigMethod(data, value)
	return
}

// Virtual key utility
func isSettingIp6ConfigAddressesEmpty(data connectionData) bool {
	addresses := getSettingIp6ConfigAddresses(data)
	if len(addresses) == 0 {
		return true
	}
	return false
}
func getOrNewSettingIp6ConfigAddresses(data connectionData) (addresses ipv6Addresses) {
	if !isSettingIp6ConfigAddressesEmpty(data) {
		addresses = getSettingIp6ConfigAddresses(data)
	} else {
		addresses = make(ipv6Addresses, 1)
		addresses[0].Gateway = make([]byte, 16)
	}
	return
}

// Virtual key getter
func getSettingVkIp6ConfigDns(data connectionData) (value string) {
	dnses := getSettingIp6ConfigDns(data)
	if len(dnses) == 0 {
		return
	}
	value = convertIpv6AddressToString(dnses[0])
	return
}
func getSettingVkIp6ConfigAddressesAddress(data connectionData) (value string) {
	if isSettingIp6ConfigAddressesEmpty(data) {
		return
	}
	addresses := getSettingIp6ConfigAddresses(data)
	if isIpv6AddressValid(addresses[0].Address) {
		value = convertIpv6AddressToString(addresses[0].Address)
	}
	return
}
func getSettingVkIp6ConfigAddressesPrefix(data connectionData) (value uint32) {
	if isSettingIp6ConfigAddressesEmpty(data) {
		return
	}
	addresses := getSettingIp6ConfigAddresses(data)
	value = addresses[0].Prefix
	logger.Info(addresses) // TODO test
	return
}
func getSettingVkIp6ConfigAddressesGateway(data connectionData) (value string) {
	if isSettingIp6ConfigAddressesEmpty(data) {
		return
	}
	addresses := getSettingIp6ConfigAddresses(data)
	value = convertIpv6AddressToStringNoZero(addresses[0].Gateway)
	return
}
func getSettingVkIp6ConfigRoutesAddress(data connectionData) (value string) {
	// TODO
	// value := getSettingIp6ConfigRoutesAddress(data)
	return
}
func getSettingVkIp6ConfigRoutesPrefix(data connectionData) (value string) {
	// TODO
	// value := getSettingIp6ConfigRoutesPrefix(data)
	return
}
func getSettingVkIp6ConfigRoutesNexthop(data connectionData) (value string) {
	// TODO
	// value := getSettingIp6ConfigRoutesNexthop(data)
	return
}
func getSettingVkIp6ConfigRoutesMetric(data connectionData) (value string) {
	// TODO
	// value := getSettingIp6ConfigRoutesMetric(data)
	return
}

// Virtual key logic setter
func logicSetSettingVkIp6ConfigDns(data connectionData, value string) (err error) {
	if len(value) == 0 {
		removeSettingIp6ConfigDns(data)
		return
	}
	dnses := getSettingIp6ConfigDns(data)
	if len(dnses) == 0 {
		dnses = make([][]byte, 1)
	}
	tmp, err := convertIpv6AddressToArrayByteCheck(value)
	dnses[0] = tmp
	if err != nil {
		err = fmt.Errorf(NM_KEY_ERROR_INVALID_VALUE)
		return
	}
	if !isIpv6AddressZero(dnses[0]) {
		setSettingIp6ConfigDns(data, dnses)
	} else {
		removeSettingIp6ConfigDns(data)
	}
	return
}
func logicSetSettingVkIp6ConfigAddressesAddress(data connectionData, value string) (err error) {
	if len(value) == 0 {
		value = ipv6AddrZero
	}
	tmp, err := convertIpv6AddressToArrayByteCheck(value)
	if err != nil {
		err = fmt.Errorf(NM_KEY_ERROR_INVALID_VALUE)
		return
	}
	addresses := getOrNewSettingIp6ConfigAddresses(data)
	addr := addresses[0]
	addr.Address = tmp
	addresses[0] = addr
	if !isIpv6AddressStructZero(addr) {
		setSettingIp6ConfigAddresses(data, addresses)
	} else {
		removeSettingIp6ConfigAddresses(data)
	}
	return
}
func logicSetSettingVkIp6ConfigAddressesPrefix(data connectionData, value uint32) (err error) {
	addresses := getOrNewSettingIp6ConfigAddresses(data)
	addr := addresses[0]
	addr.Prefix = value
	addresses[0] = addr
	if !isIpv6AddressStructZero(addr) {
		setSettingIp6ConfigAddresses(data, addresses)
	} else {
		removeSettingIp6ConfigAddresses(data)
	}
	return
}
func logicSetSettingVkIp6ConfigAddressesGateway(data connectionData, value string) (err error) {
	if len(value) == 0 {
		value = ipv6AddrZero
	}
	tmp, err := convertIpv6AddressToArrayByteCheck(value)
	if err != nil {
		err = fmt.Errorf(NM_KEY_ERROR_INVALID_VALUE)
		return
	}
	addresses := getOrNewSettingIp6ConfigAddresses(data)
	addr := addresses[0]
	addr.Gateway = tmp
	addresses[0] = addr
	if !isIpv6AddressStructZero(addr) {
		setSettingIp6ConfigAddresses(data, addresses)
	} else {
		removeSettingIp6ConfigAddresses(data)
	}
	return
}
func logicSetSettingVkIp6ConfigRoutesAddress(data connectionData, value string) (err error) {
	// TODO
	// setSettingIp6ConfigRoutesAddressJSON(data)
	return
}
func logicSetSettingVkIp6ConfigRoutesPrefix(data connectionData, value uint32) (err error) {
	// TODO
	// setSettingIp6ConfigRoutesPrefixJSON(data)
	return
}
func logicSetSettingVkIp6ConfigRoutesNexthop(data connectionData, value string) (err error) {
	// TODO
	// setSettingIp6ConfigRoutesNexthopJSON(data)
	return
}
func logicSetSettingVkIp6ConfigRoutesMetric(data connectionData, value uint32) (err error) {
	// TODO
	// setSettingIp6ConfigRoutesMetricJSON(data)
	return
}
