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

import (
	"fmt"
	. "pkg.deepin.io/lib/gettext"
)

func initSettingSectionIpv4(data connectionData) {
	addSettingSection(data, sectionIpv4)
	setSettingIp4ConfigMethod(data, NM_SETTING_IP4_CONFIG_METHOD_AUTO)
}

// Initialize available values
var availableValuesIp4ConfigMethod = make(availableValues)

func initAvailableValuesIp4() {
	availableValuesIp4ConfigMethod[NM_SETTING_IP4_CONFIG_METHOD_AUTO] = kvalue{NM_SETTING_IP4_CONFIG_METHOD_AUTO, Tr("Auto")}
	availableValuesIp4ConfigMethod[NM_SETTING_IP4_CONFIG_METHOD_LINK_LOCAL] = kvalue{NM_SETTING_IP4_CONFIG_METHOD_LINK_LOCAL, Tr("Link-Local Only")}
	availableValuesIp4ConfigMethod[NM_SETTING_IP4_CONFIG_METHOD_MANUAL] = kvalue{NM_SETTING_IP4_CONFIG_METHOD_MANUAL, Tr("Manual")}
	availableValuesIp4ConfigMethod[NM_SETTING_IP4_CONFIG_METHOD_SHARED] = kvalue{NM_SETTING_IP4_CONFIG_METHOD_SHARED, Tr("Shared")}
	availableValuesIp4ConfigMethod[NM_SETTING_IP4_CONFIG_METHOD_DISABLED] = kvalue{NM_SETTING_IP4_CONFIG_METHOD_DISABLED, Tr("Disabled")}
}

// Get available keys
func getSettingIp4ConfigAvailableKeys(data connectionData) (keys []string) {
	keys = appendAvailableKeys(data, keys, sectionIpv4, NM_SETTING_IP4_CONFIG_METHOD)
	method := getSettingIp4ConfigMethod(data)
	switch method {
	default:
		logger.Error("ip4 config method is invalid:", method)
	case NM_SETTING_IP4_CONFIG_METHOD_AUTO:
		keys = appendAvailableKeys(data, keys, sectionIpv4, NM_SETTING_IP4_CONFIG_DNS)
	case NM_SETTING_IP4_CONFIG_METHOD_LINK_LOCAL: // ignore
	case NM_SETTING_IP4_CONFIG_METHOD_MANUAL:
		keys = appendAvailableKeys(data, keys, sectionIpv4, NM_SETTING_IP4_CONFIG_DNS)
		keys = appendAvailableKeys(data, keys, sectionIpv4, NM_SETTING_IP4_CONFIG_ADDRESSES)
	case NM_SETTING_IP4_CONFIG_METHOD_SHARED:
	case NM_SETTING_IP4_CONFIG_METHOD_DISABLED:
	}
	return
}

// Get available values
func getSettingIp4ConfigAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case NM_SETTING_IP4_CONFIG_METHOD:
		// TODO be careful, ipv4 method would be limited for different connection type
		// switch getCustomConnectionType(data) {
		// case typeWired:
		// case typeWireless:
		// case typePppoe:
		// }
		// values = []string{
		// 	NM_SETTING_IP4_CONFIG_METHOD_AUTO,
		// 	// NM_SETTING_IP4_CONFIG_METHOD_LINK_LOCAL, // ignore
		// 	NM_SETTING_IP4_CONFIG_METHOD_MANUAL,
		// 	// NM_SETTING_IP4_CONFIG_METHOD_SHARED,   // ignore
		// 	// NM_SETTING_IP4_CONFIG_METHOD_DISABLED, // ignore
		// }
		if getSettingConnectionType(data) == NM_SETTING_VPN_SETTING_NAME {
			values = []kvalue{
				availableValuesIp4ConfigMethod[NM_SETTING_IP4_CONFIG_METHOD_AUTO],
			}
		} else {
			values = []kvalue{
				availableValuesIp4ConfigMethod[NM_SETTING_IP4_CONFIG_METHOD_AUTO],
				availableValuesIp4ConfigMethod[NM_SETTING_IP4_CONFIG_METHOD_MANUAL],
			}
		}
	}
	return
}

// Check whether the values are correct
func checkSettingIp4ConfigValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)

	// check method
	ensureSettingIp4ConfigMethodNoEmpty(data, errs)
	switch getSettingIp4ConfigMethod(data) {
	default:
		rememberError(errs, sectionIpv4, NM_SETTING_IP4_CONFIG_METHOD, NM_KEY_ERROR_INVALID_VALUE)
		return
	case NM_SETTING_IP4_CONFIG_METHOD_AUTO:
	case NM_SETTING_IP4_CONFIG_METHOD_LINK_LOCAL: // ignore
		checkSettingIp4MethodConflict(data, errs)
	case NM_SETTING_IP4_CONFIG_METHOD_MANUAL:
		ensureSettingIp4ConfigAddressesNoEmpty(data, errs)
	case NM_SETTING_IP4_CONFIG_METHOD_SHARED:
		checkSettingIp4MethodConflict(data, errs)
	case NM_SETTING_IP4_CONFIG_METHOD_DISABLED: // ignore
		checkSettingIp4MethodConflict(data, errs)
	}

	// check value of dns
	checkSettingIp4ConfigDns(data, errs)

	// check value of address
	checkSettingIp4ConfigAddresses(data, errs)

	// TODO check value of route

	return
}
func checkSettingIp4MethodConflict(data connectionData, errs sectionErrors) {
	// check dns
	if isSettingIp4ConfigDnsExists(data) && len(getSettingIp4ConfigDns(data)) > 0 {
		rememberError(errs, sectionIpv4, NM_SETTING_IP4_CONFIG_DNS, fmt.Sprintf(NM_KEY_ERROR_IP4_METHOD_CONFLICT, NM_SETTING_IP4_CONFIG_DNS))
	}
	// check dns search
	if isSettingIp4ConfigDnsSearchExists(data) && len(getSettingIp4ConfigDnsSearch(data)) > 0 {
		rememberError(errs, sectionIpv4, NM_SETTING_IP4_CONFIG_DNS_SEARCH, fmt.Sprintf(NM_KEY_ERROR_IP4_METHOD_CONFLICT, NM_SETTING_IP4_CONFIG_DNS_SEARCH))
	}
	// check address
	if isSettingIp4ConfigAddressesExists(data) && len(getSettingIp4ConfigAddresses(data)) > 0 {
		rememberError(errs, sectionIpv4, NM_SETTING_IP4_CONFIG_ADDRESSES, fmt.Sprintf(NM_KEY_ERROR_IP4_METHOD_CONFLICT, NM_SETTING_IP4_CONFIG_ADDRESSES))
	}
	// check route
	if isSettingIp4ConfigRoutesExists(data) && len(getSettingIp4ConfigRoutes(data)) > 0 {
		rememberError(errs, sectionIpv4, NM_SETTING_IP4_CONFIG_ROUTES, fmt.Sprintf(NM_KEY_ERROR_IP4_METHOD_CONFLICT, NM_SETTING_IP4_CONFIG_ROUTES))
	}
}
func checkSettingIp4ConfigDns(data connectionData, errs sectionErrors) {
	if !isSettingIp4ConfigDnsExists(data) {
		return
	}
	dnses := getSettingIp4ConfigDns(data)
	for _, dns := range dnses {
		if dns == 0 {
			rememberError(errs, sectionIpv4, NM_SETTING_IP4_CONFIG_DNS, NM_KEY_ERROR_INVALID_VALUE)
			return
		}
	}
}
func checkSettingIp4ConfigAddresses(data connectionData, errs sectionErrors) {
	if !isSettingIp4ConfigAddressesExists(data) {
		return
	}
	addresses := getSettingIp4ConfigAddresses(data)
	for _, addr := range addresses {
		// check address struct
		if len(addr) != 3 {
			rememberError(errs, sectionIpv4, NM_SETTING_IP4_CONFIG_ADDRESSES, NM_KEY_ERROR_IP4_ADDRESSES_STRUCT)
		}
		// check address
		if addr[0] == 0 {
			rememberError(errs, sectionIpv4, NM_SETTING_VK_IP4_CONFIG_ADDRESSES_ADDRESS, NM_KEY_ERROR_INVALID_VALUE)
		}
		// check prefix
		if addr[1] < 1 || addr[1] > 32 {
			rememberError(errs, sectionIpv4, NM_SETTING_VK_IP4_CONFIG_ADDRESSES_MASK, NM_KEY_ERROR_INVALID_VALUE)
		}
	}
}

// Logic setter
func logicSetSettingIp4ConfigMethod(data connectionData, value string) (err error) {
	// just ignore error here and set value directly, error will be
	// check in checkSettingXXXValues()
	// TODO check logic for different connection types
	switch value {
	case NM_SETTING_IP4_CONFIG_METHOD_AUTO:
		removeSettingIp4ConfigAddresses(data)
	case NM_SETTING_IP4_CONFIG_METHOD_LINK_LOCAL: // ignore
		removeSettingIp4ConfigDns(data)
		removeSettingIp4ConfigDnsSearch(data)
		removeSettingIp4ConfigAddresses(data)
		removeSettingIp4ConfigRoutes(data)
	case NM_SETTING_IP4_CONFIG_METHOD_MANUAL:
	case NM_SETTING_IP4_CONFIG_METHOD_SHARED:
		removeSettingIp4ConfigDns(data)
		removeSettingIp4ConfigDnsSearch(data)
		removeSettingIp4ConfigAddresses(data)
		removeSettingIp4ConfigRoutes(data)
	case NM_SETTING_IP4_CONFIG_METHOD_DISABLED: // ignore
		removeSettingIp4ConfigDns(data)
		removeSettingIp4ConfigDnsSearch(data)
		removeSettingIp4ConfigAddresses(data)
		removeSettingIp4ConfigRoutes(data)
	}
	setSettingIp4ConfigMethod(data, value)
	return
}

// Virtual key utility
func isSettingIp4ConfigAddressesEmpty(data connectionData) bool {
	addresses := getSettingIp4ConfigAddresses(data)
	if len(addresses) == 0 {
		return true
	}
	if len(addresses[0]) != 3 {
		return true
	}
	return false
}
func getOrNewSettingIp4ConfigAddresses(data connectionData) (addresses [][]uint32) {
	if !isSettingIp4ConfigAddressesEmpty(data) {
		addresses = getSettingIp4ConfigAddresses(data)
	} else {
		addresses = make([][]uint32, 1)
		addresses[0] = make([]uint32, 3)
	}
	return
}

// Virtual key getter
func getSettingVkIp4ConfigDns(data connectionData) (value string) {
	dnses := getSettingIp4ConfigDns(data)
	if len(dnses) == 0 {
		return
	}
	value = convertIpv4AddressToString(dnses[0])
	return
}
func getSettingVkIp4ConfigAddressesAddress(data connectionData) (value string) {
	if isSettingIp4ConfigAddressesEmpty(data) {
		return
	}
	addresses := getSettingIp4ConfigAddresses(data)
	value = convertIpv4AddressToString(addresses[0][0])
	return
}
func getSettingVkIp4ConfigAddressesMask(data connectionData) (value string) {
	if isSettingIp4ConfigAddressesEmpty(data) {
		return
	}
	addresses := getSettingIp4ConfigAddresses(data)
	value = convertIpv4PrefixToNetMask(addresses[0][1])
	return
}
func getSettingVkIp4ConfigAddressesGateway(data connectionData) (value string) {
	if isSettingIp4ConfigAddressesEmpty(data) {
		return
	}
	addresses := getSettingIp4ConfigAddresses(data)
	value = convertIpv4AddressToStringNoZero(addresses[0][2])
	return
}
func getSettingVkIp4ConfigRoutesAddress(data connectionData) (value string) {
	// TODO
	// value := getSettingIp4ConfigRoutesAddress(data)
	return
}
func getSettingVkIp4ConfigRoutesMask(data connectionData) (value string) {
	// TODO
	// value := getSettingIp4ConfigRoutesMask(data)
	return
}
func getSettingVkIp4ConfigRoutesNexthop(data connectionData) (value string) {
	// TODO
	// value := getSettingIp4ConfigRoutesNexthop(data)
	return
}
func getSettingVkIp4ConfigRoutesMetric(data connectionData) (value string) {
	// TODO
	// value := getSettingIp4ConfigRoutesMetric(data)
	return
}

// Virtual key logic setter
func logicSetSettingVkIp4ConfigDns(data connectionData, value string) (err error) {
	if len(value) == 0 {
		removeSettingIp4ConfigDns(data)
		return
	}
	dnses := getSettingIp4ConfigDns(data)
	if len(dnses) == 0 {
		dnses = make([]uint32, 1)
	}
	tmpn, err := convertIpv4AddressToUint32Check(value)
	dnses[0] = tmpn
	if err != nil {
		err = fmt.Errorf(NM_KEY_ERROR_INVALID_VALUE)
		return
	}
	if dnses[0] != 0 {
		setSettingIp4ConfigDns(data, dnses)
	} else {
		removeSettingIp4ConfigDns(data)
	}
	return
}
func logicSetSettingVkIp4ConfigAddressesAddress(data connectionData, value string) (err error) {
	if len(value) == 0 {
		value = ipv4Zero
	}
	tmpn, err := convertIpv4AddressToUint32Check(value)
	if err != nil {
		err = fmt.Errorf(NM_KEY_ERROR_INVALID_VALUE)
		return
	}
	addresses := getOrNewSettingIp4ConfigAddresses(data)
	addr := addresses[0]
	addr[0] = tmpn
	if !isUint32ArrayEmpty(addr) {
		setSettingIp4ConfigAddresses(data, addresses)
	} else {
		removeSettingIp4ConfigAddresses(data)
	}
	return
}
func logicSetSettingVkIp4ConfigAddressesMask(data connectionData, value string) (err error) {
	if len(value) == 0 {
		value = ipv4Zero
	}
	tmpn, err := convertIpv4NetMaskToPrefixCheck(value)
	if err != nil {
		err = fmt.Errorf(NM_KEY_ERROR_INVALID_VALUE)
		return
	}
	addresses := getOrNewSettingIp4ConfigAddresses(data)
	addr := addresses[0]
	addr[1] = tmpn
	if !isUint32ArrayEmpty(addr) {
		setSettingIp4ConfigAddresses(data, addresses)
	} else {
		removeSettingIp4ConfigAddresses(data)
	}
	return
}
func logicSetSettingVkIp4ConfigAddressesGateway(data connectionData, value string) (err error) {
	if len(value) == 0 {
		value = ipv4Zero
	}
	tmpn, err := convertIpv4AddressToUint32Check(value)
	if err != nil {
		err = fmt.Errorf(NM_KEY_ERROR_INVALID_VALUE)
		return
	}
	addresses := getOrNewSettingIp4ConfigAddresses(data)
	addr := addresses[0]
	addr[2] = tmpn
	if !isUint32ArrayEmpty(addr) {
		setSettingIp4ConfigAddresses(data, addresses)
	} else {
		removeSettingIp4ConfigAddresses(data)
	}
	return
}
func logicSetSettingVkIp4ConfigRoutesAddress(data connectionData, value string) (err error) {
	// TODO
	// setSettingIp4ConfigRoutesAddressJSON(data)
	return
}
func logicSetSettingVkIp4ConfigRoutesMask(data connectionData, value string) (err error) {
	// TODO
	// setSettingIp4ConfigRoutesMaskJSON(data)
	return
}
func logicSetSettingVkIp4ConfigRoutesNexthop(data connectionData, value string) (err error) {
	// TODO
	// setSettingIp4ConfigRoutesNexthopJSON(data)
	return
}
func logicSetSettingVkIp4ConfigRoutesMetric(data connectionData, value string) (err error) {
	// TODO
	// setSettingIp4ConfigRoutesMetricJSON(data)
	return
}
