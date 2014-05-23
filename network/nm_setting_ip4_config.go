package main

import (
	"dlib"
	"fmt"
)

const NM_SETTING_IP4_CONFIG_SETTING_NAME = "ipv4"

const (
	// IPv4 configuration method. If 'auto' is specified then the
	// appropriate automatic method (DHCP, PPP, etc) is used for the
	// interface and most other properties can be left unset. If
	// 'link-local' is specified, then a link-local address in the
	// 169.254/16 range will be assigned to the interface. If 'manual'
	// is specified, static IP addressing is used and at least one IP
	// address must be given in the 'addresses' property. If 'shared'
	// is specified (indicating that this connection will provide
	// network access to other computers) then the interface is
	// assigned an address in the 10.42.x.1/24 range and a DHCP and
	// forwarding DNS server are started, and the interface is NAT-ed
	// to the current default network connection. 'disabled' means
	// IPv4 will not be used on this connection. This property must be
	// set.
	// Default value: NULL
	NM_SETTING_IP4_CONFIG_METHOD = "method"

	// List of DNS servers (network byte order). For the 'auto'
	// method, these DNS servers are appended to those (if any)
	// returned by automatic configuration. DNS servers cannot be used
	// with the 'shared', 'link-local', or 'disabled' methods as there
	// is no upstream network. In all other methods, these DNS servers
	// are used as the only DNS servers for this connection.
	NM_SETTING_IP4_CONFIG_DNS = "dns"

	// List of DNS search domains. For the 'auto' method, these search
	// domains are appended to those returned by automatic
	// configuration. Search domains cannot be used with the 'shared',
	// 'link-local', or 'disabled' methods as there is no upstream
	// network. In all other methods, these search domains are used as
	// the only search domains for this connection.
	NM_SETTING_IP4_CONFIG_DNS_SEARCH = "dns-search"

	// Array of IPv4 address structures. Each IPv4 address structure
	// is composed of 3 32-bit values; the first being the IPv4
	// address (network byte order), the second the prefix (1 - 32),
	// and last the IPv4 gateway (network byte order). The gateway may
	// be left as 0 if no gateway exists for that subnet. For the
	// 'auto' method, given IP addresses are appended to those
	// returned by automatic configuration. Addresses cannot be used
	// with the 'shared', 'link-local', or 'disabled' methods as
	// addressing is either automatic or disabled with these methods.
	NM_SETTING_IP4_CONFIG_ADDRESSES = "addresses"

	// Array of IPv4 route structures. Each IPv4 route structure is
	// composed of 4 32-bit values; the first being the destination
	// IPv4 network or address (network byte order), the second the
	// destination network or address prefix (1 - 32), the third being
	// the next-hop (network byte order) if any, and the fourth being
	// the route metric. For the 'auto' method, given IP routes are
	// appended to those returned by automatic configuration. Routes
	// cannot be used with the 'shared', 'link-local', or 'disabled'
	// methods because there is no upstream network.
	NM_SETTING_IP4_CONFIG_ROUTES = "routes"

	// When the method is set to 'auto' and this property to TRUE,
	// automatically configured routes are ignored and only routes
	// specified in "routes", if any, are used.
	// Default value: FALSE
	NM_SETTING_IP4_CONFIG_IGNORE_AUTO_ROUTES = "ignore-auto-routes"

	// When the method is set to 'auto' and this property to TRUE,
	// automatically configured nameservers and search domains are
	// ignored and only nameservers and search domains specified in
	// "dns" and "dns-search", if any, are used.
	// Default value: FALSE
	NM_SETTING_IP4_CONFIG_IGNORE_AUTO_DNS = "ignore-auto-dns"

	// A string sent to the DHCP server to identify the local machine
	// which the DHCP server may use to customize the DHCP lease and
	// options.
	// Default value: NULL
	NM_SETTING_IP4_CONFIG_DHCP_CLIENT_ID = "dhcp-client-id"

	// If TRUE, a hostname is sent to the DHCP server when acquiring a
	// lease. Some DHCP servers use this hostname to update DNS
	// databases, essentially providing a static hostname for the
	// computer. If "dhcp-hostname" is empty and this property is
	// TRUE, the current persistent hostname of the computer is sent.
	// Default value: TRUE
	NM_SETTING_IP4_CONFIG_DHCP_SEND_HOSTNAME = "dhcp-send-hostname"

	// If the "dhcp-send-hostname" property is TRUE, then the
	// specified name will be sent to the DHCP server when acquiring a
	// lease.
	// Default value: NULL
	NM_SETTING_IP4_CONFIG_DHCP_HOSTNAME = "dhcp-hostname"

	// If TRUE, this connection will never be the default IPv4
	// connection, meaning it will never be assigned the default route
	// by NetworkManager.
	// Default value: FALSE
	NM_SETTING_IP4_CONFIG_NEVER_DEFAULT = "never-default"

	// If TRUE, allow overall network configuration to proceed even if
	// IPv4 configuration times out. Note that at least one IP
	// configuration must succeed or overall network configuration
	// will still fail. For example, in IPv6-only networks, setting
	// this property to TRUE allows the overall network configuration
	// to succeed if IPv4 configuration fails but IPv6 configuration
	// completes successfully.
	// Default value: TRUE
	NM_SETTING_IP4_CONFIG_MAY_FAIL = "may-fail"
)

const (
	// IPv4 configuration should be automatically determined via a
	// method appropriate for the hardware interface, ie DHCP or PPP
	// or some other device-specific manner.
	NM_SETTING_IP4_CONFIG_METHOD_AUTO = "auto"

	// IPv4 configuration should be automatically configured for
	// link-local-only operation.
	NM_SETTING_IP4_CONFIG_METHOD_LINK_LOCAL = "link-local"

	// All necessary IPv4 configuration (addresses, prefix, DNS, etc)
	// is specified in the setting's properties.
	NM_SETTING_IP4_CONFIG_METHOD_MANUAL = "manual"

	// This connection specifies configuration that allows other
	// computers to connect through it to the default network (usually
	// the Internet). The connection's interface will be assigned a
	// private address, and a DHCP server, caching DNS server, and
	// Network Address Translation (NAT) functionality will be started
	// on this connection's interface to allow other devices to
	// connect through that interface to the default network.
	NM_SETTING_IP4_CONFIG_METHOD_SHARED = "shared"

	// This connection does not use or require IPv4 address and it
	// should be disabled.
	NM_SETTING_IP4_CONFIG_METHOD_DISABLED = "disabled"
)

func initSettingFieldIpv4(data connectionData) {
	addSettingField(data, fieldIpv4)
	setSettingIp4ConfigMethod(data, NM_SETTING_IP4_CONFIG_METHOD_AUTO)
}

// Initialize available values
var availableValuesIp4ConfigMethod = make(availableValues)

func init() {
	availableValuesIp4ConfigMethod[NM_SETTING_IP4_CONFIG_METHOD_AUTO] = kvalue{NM_SETTING_IP4_CONFIG_METHOD_AUTO, dlib.Tr("Auto")}
	availableValuesIp4ConfigMethod[NM_SETTING_IP4_CONFIG_METHOD_LINK_LOCAL] = kvalue{NM_SETTING_IP4_CONFIG_METHOD_LINK_LOCAL, dlib.Tr("Link Local")}
	availableValuesIp4ConfigMethod[NM_SETTING_IP4_CONFIG_METHOD_MANUAL] = kvalue{NM_SETTING_IP4_CONFIG_METHOD_MANUAL, dlib.Tr("Manual")}
	availableValuesIp4ConfigMethod[NM_SETTING_IP4_CONFIG_METHOD_SHARED] = kvalue{NM_SETTING_IP4_CONFIG_METHOD_SHARED, dlib.Tr("Shared")}
	availableValuesIp4ConfigMethod[NM_SETTING_IP4_CONFIG_METHOD_DISABLED] = kvalue{NM_SETTING_IP4_CONFIG_METHOD_DISABLED, dlib.Tr("Disabled")}
}

// Get available keys
func getSettingIp4ConfigAvailableKeys(data connectionData) (keys []string) {
	method := getSettingIp4ConfigMethod(data)
	switch method {
	default:
		logger.Error("ip4 config method is invalid:", method)
	case NM_SETTING_IP4_CONFIG_METHOD_AUTO:
		keys = appendAvailableKeys(data, keys, fieldIpv4, NM_SETTING_IP4_CONFIG_METHOD)
		keys = appendAvailableKeys(data, keys, fieldIpv4, NM_SETTING_IP4_CONFIG_DNS)
	case NM_SETTING_IP4_CONFIG_METHOD_LINK_LOCAL: // ignore
	case NM_SETTING_IP4_CONFIG_METHOD_MANUAL:
		keys = appendAvailableKeys(data, keys, fieldIpv4, NM_SETTING_IP4_CONFIG_METHOD)
		keys = appendAvailableKeys(data, keys, fieldIpv4, NM_SETTING_IP4_CONFIG_DNS)
		keys = appendAvailableKeys(data, keys, fieldIpv4, NM_SETTING_IP4_CONFIG_ADDRESSES)
	case NM_SETTING_IP4_CONFIG_METHOD_SHARED: // ignore
	case NM_SETTING_IP4_CONFIG_METHOD_DISABLED:
		keys = appendAvailableKeys(data, keys, fieldIpv4, NM_SETTING_IP4_CONFIG_METHOD)
	}
	return
}

// Get available values
func getSettingIp4ConfigAvailableValues(data connectionData, key string) (values []kvalue) {
	switch key {
	case NM_SETTING_IP4_CONFIG_METHOD:
		// TODO be careful, ipv4 method would be limited for different connection type
		// switch getCustomConnectinoType(data) {
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
		if getSettingConnectionType(data) != NM_SETTING_VPN_SETTING_NAME {
			values = []kvalue{
				availableValuesIp4ConfigMethod[NM_SETTING_IP4_CONFIG_METHOD_AUTO],
				availableValuesIp4ConfigMethod[NM_SETTING_IP4_CONFIG_METHOD_MANUAL],
			}
		} else {
			values = []kvalue{
				availableValuesIp4ConfigMethod[NM_SETTING_IP4_CONFIG_METHOD_AUTO],
			}
		}
	}
	return
}

// Check whether the values are correct
func checkSettingIp4ConfigValues(data connectionData) (errs fieldErrors) {
	errs = make(map[string]string)

	// check method
	ensureSettingIp4ConfigMethodNoEmpty(data, errs)
	switch getSettingIp4ConfigMethod(data) {
	default:
		rememberError(errs, fieldIpv4, NM_SETTING_IP4_CONFIG_METHOD, NM_KEY_ERROR_INVALID_VALUE)
		return
	case NM_SETTING_IP4_CONFIG_METHOD_AUTO:
	case NM_SETTING_IP4_CONFIG_METHOD_LINK_LOCAL: // ignore
		checkSettingIp4MethodConflict(data, errs)
	case NM_SETTING_IP4_CONFIG_METHOD_MANUAL:
		ensureSettingIp4ConfigAddressesNoEmpty(data, errs)
	case NM_SETTING_IP4_CONFIG_METHOD_SHARED: // ignore
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
func checkSettingIp4MethodConflict(data connectionData, errs fieldErrors) {
	// check dns
	if isSettingIp4ConfigDnsExists(data) {
		rememberError(errs, fieldIpv4, NM_SETTING_IP4_CONFIG_DNS, fmt.Sprintf(NM_KEY_ERROR_IP4_METHOD_CONFLICT, NM_SETTING_IP4_CONFIG_DNS))
	}
	// check dns search
	if isSettingIp4ConfigDnsSearchExists(data) {
		rememberError(errs, fieldIpv4, NM_SETTING_IP4_CONFIG_DNS_SEARCH, fmt.Sprintf(NM_KEY_ERROR_IP4_METHOD_CONFLICT, NM_SETTING_IP4_CONFIG_DNS_SEARCH))
	}
	// check address
	if isSettingIp4ConfigAddressesExists(data) {
		rememberError(errs, fieldIpv4, NM_SETTING_IP4_CONFIG_ADDRESSES, fmt.Sprintf(NM_KEY_ERROR_IP4_METHOD_CONFLICT, NM_SETTING_IP4_CONFIG_ADDRESSES))
	}
	// check route
	if isSettingIp4ConfigRoutesExists(data) {
		rememberError(errs, fieldIpv4, NM_SETTING_IP4_CONFIG_ROUTES, fmt.Sprintf(NM_KEY_ERROR_IP4_METHOD_CONFLICT, NM_SETTING_IP4_CONFIG_ROUTES))
	}
}
func checkSettingIp4ConfigDns(data connectionData, errs fieldErrors) {
	if !isSettingIp4ConfigDnsExists(data) {
		return
	}
	dnses := getSettingIp4ConfigDns(data)
	for _, dns := range dnses {
		if dns == 0 {
			rememberError(errs, fieldIpv4, NM_SETTING_IP4_CONFIG_DNS, NM_KEY_ERROR_INVALID_VALUE)
			return
		}
	}
}
func checkSettingIp4ConfigAddresses(data connectionData, errs fieldErrors) {
	if !isSettingIp4ConfigAddressesExists(data) {
		return
	}
	addresses := getSettingIp4ConfigAddresses(data)
	for _, addr := range addresses {
		// check address struct
		if len(addr) != 3 {
			rememberError(errs, fieldIpv4, NM_SETTING_IP4_CONFIG_ADDRESSES, NM_KEY_ERROR_IP4_ADDRESSES_STRUCT)
		}
		// check address
		if addr[0] == 0 {
			rememberError(errs, fieldIpv4, NM_SETTING_VK_IP4_CONFIG_ADDRESSES_ADDRESS, NM_KEY_ERROR_INVALID_VALUE)
		}
		// check prefix
		if addr[1] < 1 || addr[1] > 32 {
			rememberError(errs, fieldIpv4, NM_SETTING_VK_IP4_CONFIG_ADDRESSES_MASK, NM_KEY_ERROR_INVALID_VALUE)
		}
	}
}

// Logic setter
func logicSetSettingIp4ConfigMethod(data connectionData, value string) (err error) {
	// just ignore error here and set value directly, error will be
	// check in checkSettingXXXValues()
	switch value {
	case NM_SETTING_IP4_CONFIG_METHOD_AUTO:
		removeSettingIp4ConfigAddresses(data)
	case NM_SETTING_IP4_CONFIG_METHOD_LINK_LOCAL: // ignore
		removeSettingIp4ConfigDns(data)
		removeSettingIp4ConfigDnsSearch(data)
		removeSettingIp4ConfigAddresses(data)
		removeSettingIp4ConfigRoutes(data)
	case NM_SETTING_IP4_CONFIG_METHOD_MANUAL:
	case NM_SETTING_IP4_CONFIG_METHOD_SHARED: // ignore
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
	dns := getSettingIp4ConfigDns(data)
	if len(dns) == 0 {
		return
	}
	value = convertIpv4AddressToString(dns[0])
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
	dns := getSettingIp4ConfigDns(data)
	if len(dns) == 0 {
		dns = make([]uint32, 1)
	}
	tmpn, err := convertIpv4AddressToUint32Check(value)
	dns[0] = tmpn
	if err != nil {
		err = fmt.Errorf(NM_KEY_ERROR_INVALID_VALUE)
		return
	}
	if dns[0] != 0 {
		setSettingIp4ConfigDns(data, dns)
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
