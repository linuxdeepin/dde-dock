package main

import (
	"dlib"
	"fmt"
)

// TODO doc

const NM_SETTING_IP6_CONFIG_SETTING_NAME = "ipv6"

const (
	NM_SETTING_IP6_CONFIG_METHOD             = "method"
	NM_SETTING_IP6_CONFIG_DNS                = "dns"
	NM_SETTING_IP6_CONFIG_DNS_SEARCH         = "dns-search"
	NM_SETTING_IP6_CONFIG_ADDRESSES          = "addresses"
	NM_SETTING_IP6_CONFIG_ROUTES             = "routes"
	NM_SETTING_IP6_CONFIG_IGNORE_AUTO_ROUTES = "ignore-auto-routes"
	NM_SETTING_IP6_CONFIG_IGNORE_AUTO_DNS    = "ignore-auto-dns"
	NM_SETTING_IP6_CONFIG_NEVER_DEFAULT      = "never-default"
	NM_SETTING_IP6_CONFIG_MAY_FAIL           = "may-fail"
	NM_SETTING_IP6_CONFIG_IP6_PRIVACY        = "ip6-privacy"
	NM_SETTING_IP6_CONFIG_DHCP_HOSTNAME      = "dhcp-hostname"
)

const (
	NM_SETTING_IP6_CONFIG_METHOD_IGNORE     = "ignore"
	NM_SETTING_IP6_CONFIG_METHOD_AUTO       = "auto"
	NM_SETTING_IP6_CONFIG_METHOD_DHCP       = "dhcp"
	NM_SETTING_IP6_CONFIG_METHOD_LINK_LOCAL = "link-local"
	NM_SETTING_IP6_CONFIG_METHOD_MANUAL     = "manual"
	NM_SETTING_IP6_CONFIG_METHOD_SHARED     = "shared"
)

func initSettingFieldIpv6(data connectionData) {
	addSettingField(data, fieldIpv6)
	setSettingIp6ConfigMethod(data, NM_SETTING_IP6_CONFIG_METHOD_AUTO)
}

// Initialize available values
var availableValuesIp6ConfigMethod = make(availableValues)

func init() {
	availableValuesIp6ConfigMethod[NM_SETTING_IP6_CONFIG_METHOD_IGNORE] = kvalue{NM_SETTING_IP6_CONFIG_METHOD_IGNORE, dlib.Tr("Ignore")}
	availableValuesIp6ConfigMethod[NM_SETTING_IP6_CONFIG_METHOD_AUTO] = kvalue{NM_SETTING_IP6_CONFIG_METHOD_AUTO, dlib.Tr("Auto")}
	availableValuesIp6ConfigMethod[NM_SETTING_IP6_CONFIG_METHOD_DHCP] = kvalue{NM_SETTING_IP6_CONFIG_METHOD_DHCP, dlib.Tr("DHCP")}
	availableValuesIp6ConfigMethod[NM_SETTING_IP6_CONFIG_METHOD_LINK_LOCAL] = kvalue{NM_SETTING_IP6_CONFIG_METHOD_LINK_LOCAL, dlib.Tr("Link Local")}
	availableValuesIp6ConfigMethod[NM_SETTING_IP6_CONFIG_METHOD_MANUAL] = kvalue{NM_SETTING_IP6_CONFIG_METHOD_MANUAL, dlib.Tr("Manual")}
	availableValuesIp6ConfigMethod[NM_SETTING_IP6_CONFIG_METHOD_SHARED] = kvalue{NM_SETTING_IP6_CONFIG_METHOD_SHARED, dlib.Tr("Shared")}
}

// Get available keys
func getSettingIp6ConfigAvailableKeys(data connectionData) (keys []string) {
	method := getSettingIp6ConfigMethod(data)
	switch method {
	default:
		logger.Error("ip6 config method is invalid:", method)
	case NM_SETTING_IP6_CONFIG_METHOD_IGNORE:
		keys = appendAvailableKeys(data, keys, fieldIpv6, NM_SETTING_IP6_CONFIG_METHOD)
	case NM_SETTING_IP6_CONFIG_METHOD_AUTO:
		keys = appendAvailableKeys(data, keys, fieldIpv6, NM_SETTING_IP6_CONFIG_METHOD)
		keys = appendAvailableKeys(data, keys, fieldIpv6, NM_SETTING_IP6_CONFIG_DNS)
	case NM_SETTING_IP6_CONFIG_METHOD_DHCP: // ignore
		keys = appendAvailableKeys(data, keys, fieldIpv6, NM_SETTING_IP6_CONFIG_METHOD)
		keys = appendAvailableKeys(data, keys, fieldIpv6, NM_SETTING_IP6_CONFIG_DNS)
	case NM_SETTING_IP6_CONFIG_METHOD_LINK_LOCAL: // ignore
	case NM_SETTING_IP6_CONFIG_METHOD_MANUAL:
		keys = appendAvailableKeys(data, keys, fieldIpv6, NM_SETTING_IP6_CONFIG_METHOD)
		keys = appendAvailableKeys(data, keys, fieldIpv6, NM_SETTING_IP6_CONFIG_DNS)
		keys = appendAvailableKeys(data, keys, fieldIpv6, NM_SETTING_IP6_CONFIG_ADDRESSES)
	case NM_SETTING_IP6_CONFIG_METHOD_SHARED: // ignore
	}
	return
}

// Get available values
func getSettingIp6ConfigAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case NM_SETTING_IP6_CONFIG_METHOD:
		// values = []string{
		// 	// NM_SETTING_IP6_CONFIG_METHOD_IGNORE, // ignore
		// 	NM_SETTING_IP6_CONFIG_METHOD_AUTO,
		// 	// NM_SETTING_IP6_CONFIG_METHOD_DHCP, // ignore
		// 	// NM_SETTING_IP6_CONFIG_METHOD_LINK_LOCAL, // ignore
		// 	NM_SETTING_IP6_CONFIG_METHOD_MANUAL,
		// 	// NM_SETTING_IP6_CONFIG_METHOD_SHARED,// ignore
		// }
		if getSettingConnectionType(data) != NM_SETTING_VPN_SETTING_NAME {
			values = []kvalue{
				availableValuesIp6ConfigMethod[NM_SETTING_IP6_CONFIG_METHOD_AUTO],
				availableValuesIp6ConfigMethod[NM_SETTING_IP6_CONFIG_METHOD_MANUAL],
			}
		} else {
			values = []kvalue{
				availableValuesIp6ConfigMethod[NM_SETTING_IP6_CONFIG_METHOD_AUTO],
			}
		}
	}
	return
}

// Check whether the values are correct
func checkSettingIp6ConfigValues(data connectionData) (errs fieldErrors) {
	errs = make(map[string]string)

	// check method
	ensureSettingIp6ConfigMethodNoEmpty(data, errs)
	switch getSettingIp6ConfigMethod(data) {
	default:
		rememberError(errs, fieldIpv6, NM_SETTING_IP6_CONFIG_METHOD, NM_KEY_ERROR_INVALID_VALUE)
		return
	case NM_SETTING_IP6_CONFIG_METHOD_IGNORE: // ignore
		checkSettingIp6MethodConflict(data, errs)
	case NM_SETTING_IP6_CONFIG_METHOD_AUTO:
	case NM_SETTING_IP6_CONFIG_METHOD_DHCP: // ignore
	case NM_SETTING_IP6_CONFIG_METHOD_LINK_LOCAL: // ignore
		checkSettingIp6MethodConflict(data, errs)
	case NM_SETTING_IP6_CONFIG_METHOD_MANUAL:
		ensureSettingIp6ConfigAddressesNoEmpty(data, errs)
	case NM_SETTING_IP6_CONFIG_METHOD_SHARED: // ignore
		checkSettingIp6MethodConflict(data, errs)
	}

	// check value of dns
	checkSettingIp6ConfigDns(data, errs)

	// check value of address
	checkSettingIp6ConfigAddresses(data, errs)

	// TODO check value of route

	return
}

func checkSettingIp6MethodConflict(data connectionData, errs fieldErrors) {
	// check dns
	if isSettingIp6ConfigDnsExists(data) {
		rememberError(errs, fieldIpv6, NM_SETTING_IP6_CONFIG_DNS, fmt.Sprintf(NM_KEY_ERROR_IP6_METHOD_CONFLICT, NM_SETTING_IP6_CONFIG_DNS))
	}
	// check dns search
	if isSettingIp6ConfigDnsSearchExists(data) {
		rememberError(errs, fieldIpv6, NM_SETTING_IP6_CONFIG_DNS_SEARCH, fmt.Sprintf(NM_KEY_ERROR_IP6_METHOD_CONFLICT, NM_SETTING_IP6_CONFIG_DNS_SEARCH))
	}
	// check address
	if isSettingIp6ConfigAddressesExists(data) {
		rememberError(errs, fieldIpv6, NM_SETTING_IP6_CONFIG_ADDRESSES, fmt.Sprintf(NM_KEY_ERROR_IP6_METHOD_CONFLICT, NM_SETTING_IP6_CONFIG_ADDRESSES))
	}
	// check route
	if isSettingIp6ConfigRoutesExists(data) {
		rememberError(errs, fieldIpv6, NM_SETTING_IP6_CONFIG_ROUTES, fmt.Sprintf(NM_KEY_ERROR_IP6_METHOD_CONFLICT, NM_SETTING_IP6_CONFIG_ROUTES))
	}
}

func checkSettingIp6ConfigDns(data connectionData, errs fieldErrors) {
	if !isSettingIp6ConfigDnsExists(data) {
		return
	}
	dnses := getSettingIp6ConfigDns(data)
	for _, dns := range dnses {
		if !isIpv6AddressValid(dns) {
			rememberError(errs, fieldIpv6, NM_SETTING_VK_IP6_CONFIG_DNS, NM_KEY_ERROR_INVALID_VALUE)
			return
		}
		if isIpv6AddressZero(dns) {
			rememberError(errs, fieldIpv6, NM_SETTING_VK_IP6_CONFIG_DNS, NM_KEY_ERROR_INVALID_VALUE)
			return
		}
	}
}

func checkSettingIp6ConfigAddresses(data connectionData, errs fieldErrors) {
	if !isSettingIp6ConfigAddressesExists(data) {
		return
	}
	addresses := getSettingIp6ConfigAddresses(data)
	for _, addr := range addresses {
		// check address
		if !isIpv6AddressValid(addr.Address) {
			rememberError(errs, fieldIpv6, NM_SETTING_VK_IP6_CONFIG_ADDRESSES_ADDRESS, NM_KEY_ERROR_INVALID_VALUE)
			return
		}
		if isIpv6AddressZero(addr.Address) {
			rememberError(errs, fieldIpv6, NM_SETTING_VK_IP6_CONFIG_ADDRESSES_ADDRESS, NM_KEY_ERROR_INVALID_VALUE)
			return
		}
		// check prefix
		if addr.Prefix < 1 || addr.Prefix > 128 {
			rememberError(errs, fieldIpv6, NM_SETTING_VK_IP6_CONFIG_ADDRESSES_PREFIX, NM_KEY_ERROR_INVALID_VALUE)
			return
		}
		// check gateway
		if !isIpv6AddressValid(addr.Gateway) {
			rememberError(errs, fieldIpv6, NM_SETTING_VK_IP6_CONFIG_ADDRESSES_GATEWAY, NM_KEY_ERROR_INVALID_VALUE)
			return
		}
	}
}

// Logic setter
func logicSetSettingIp6ConfigMethod(data connectionData, value string) (err error) {
	switch value {
	case NM_SETTING_IP6_CONFIG_METHOD_IGNORE: // ignore
		removeSettingKeyBut(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_METHOD)
	case NM_SETTING_IP6_CONFIG_METHOD_AUTO:
		removeSettingIp6ConfigAddresses(data)
	case NM_SETTING_IP6_CONFIG_METHOD_DHCP: // ignore
	case NM_SETTING_IP6_CONFIG_METHOD_LINK_LOCAL: // ignore
		removeSettingIp6ConfigDns(data)
		removeSettingIp6ConfigDnsSearch(data)
		removeSettingIp6ConfigAddresses(data)
		removeSettingIp6ConfigRoutes(data)
	case NM_SETTING_IP6_CONFIG_METHOD_MANUAL:
	case NM_SETTING_IP6_CONFIG_METHOD_SHARED: // ignore
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
	}
	return
}

// Virtual key getter
func getSettingVkIp6ConfigDns(data connectionData) (value string) {
	dns := getSettingIp6ConfigDns(data)
	if len(dns) == 0 {
		return
	}
	value = convertIpv6AddressToString(dns[0])
	return
}
func getSettingVkIp6ConfigAddressesAddress(data connectionData) (value string) {
	if isSettingIp6ConfigAddressesEmpty(data) {
		return
	}
	addresses := getSettingIp6ConfigAddresses(data)
	value = convertIpv6AddressToString(addresses[0].Address)
	return
}
func getSettingVkIp6ConfigAddressesPrefix(data connectionData) (value uint32) {
	if isSettingIp6ConfigAddressesEmpty(data) {
		return
	}
	addresses := getSettingIp6ConfigAddresses(data)
	value = addresses[0].Prefix
	return
}
func getSettingVkIp6ConfigAddressesGateway(data connectionData) (value string) {
	if isSettingIp6ConfigAddressesEmpty(data) {
		return
	}
	addresses := getSettingIp6ConfigAddresses(data)
	value = convertIpv6AddressToString(addresses[0].Gateway)
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
	dns := getSettingIp6ConfigDns(data)
	if len(dns) == 0 {
		dns = make([][]byte, 1)
	}
	tmp, err := convertIpv6AddressToArrayByteCheck(value)
	dns[0] = tmp
	if err != nil {
		err = fmt.Errorf(NM_KEY_ERROR_INVALID_VALUE)
		return
	}
	if !isIpv6AddressZero(dns[0]) {
		setSettingIp6ConfigDns(data, dns)
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
