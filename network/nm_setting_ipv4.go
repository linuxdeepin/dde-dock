package main

import (
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

// TODO Get available keys
func getSettingIp4ConfigAvailableKeys(data _ConnectionData) (keys []string) {
	method := getSettingIp4ConfigMethod(data)
	switch method {
	default:
		LOGGER.Error("ip4 config method is invalid:", method)
	case NM_SETTING_IP4_CONFIG_METHOD_AUTO:
		keys = []string{
			NM_SETTING_IP4_CONFIG_METHOD,
			NM_SETTING_IP4_CONFIG_DNS,
		}
	case NM_SETTING_IP4_CONFIG_METHOD_LINK_LOCAL: // ignore
	case NM_SETTING_IP4_CONFIG_METHOD_MANUAL:
		keys = []string{
			NM_SETTING_IP4_CONFIG_METHOD,
			NM_SETTING_IP4_CONFIG_DNS,
			NM_SETTING_IP4_CONFIG_ADDRESSES,
		}
	case NM_SETTING_IP4_CONFIG_METHOD_SHARED: // ignore
	case NM_SETTING_IP4_CONFIG_METHOD_DISABLED:
		keys = []string{
			NM_SETTING_IP4_CONFIG_METHOD,
		}
	}
	return
}

// TODO Get available values
func getSettingIp4ConfigAvailableValues(key string) (values []string, customizable bool) {
	customizable = true
	switch key {
	case NM_SETTING_IP4_CONFIG_METHOD:
		values = []string{
			NM_SETTING_IP4_CONFIG_METHOD_AUTO,
			NM_SETTING_IP4_CONFIG_METHOD_LINK_LOCAL,
			NM_SETTING_IP4_CONFIG_METHOD_MANUAL,
			NM_SETTING_IP4_CONFIG_METHOD_SHARED,
			NM_SETTING_IP4_CONFIG_METHOD_DISABLED,
		}
		customizable = false
	}
	return
}

// TODO Check whether the values are correct
func checkSettingIp4ConfigValues(data _ConnectionData) (errs map[string]string) {
	errs = make(map[string]string)
	fieldData, ok := data[fieldIPv4]
	if !ok {
		LOGGER.Warning("field ipv4 does not exist")
		return
	}
	if len(fieldData) == 0 {
		LOGGER.Warning("field ipv4 is empty")
		return
	}

	for key, _ := range fieldData {
		availableValues, customizable := getSettingIp4ConfigAvailableValues(key)
		wrappedValue := generalGetSettingIp4ConfigKeyJSON(data, key) // TODO
		if !customizable && !isStringInArray(wrappedValue, availableValues) {
			errs[key] = "invalid value: " + wrappedValue
			continue
		}

		method := getSettingIp4ConfigMethod(data)
		methodSharedLinkloacalDisable := []string{NM_SETTING_IP4_CONFIG_METHOD_SHARED,
			NM_SETTING_IP4_CONFIG_METHOD_LINK_LOCAL,
			NM_SETTING_IP4_CONFIG_METHOD_DISABLED}
		switch key {
		case NM_SETTING_IP4_CONFIG_METHOD:
		case NM_SETTING_IP4_CONFIG_DNS:
			if isStringInArray(method, methodSharedLinkloacalDisable) {
				errs[key] = fmt.Sprintf(`key "%s" cannot be used with the 'shared', 'link-local', or 'disabled' methods`, key)
				continue
			}
			// TODO
		case NM_SETTING_IP4_CONFIG_DNS_SEARCH: // TODO
			if isStringInArray(method, methodSharedLinkloacalDisable) {
				errs[key] = fmt.Sprintf(`key "%s" cannot be used with the 'shared', 'link-local', or 'disabled' methods`, key)
				continue
			}
		case NM_SETTING_IP4_CONFIG_ADDRESSES:
			if isStringInArray(method, methodSharedLinkloacalDisable) {
				errs[key] = fmt.Sprintf(`key "%s" cannot be used with the 'shared', 'link-local', or 'disabled' methods`, key)
				continue
			}

			// if method is "manual", addresses cannot be empty
			if method == NM_SETTING_IP4_CONFIG_METHOD_MANUAL && len(wrappedValue) == 0 {
				errs[key] = "ip address cannot be empty"
				continue
			}

			// TODO check address format
		case NM_SETTING_IP4_CONFIG_ROUTES:
			if isStringInArray(method, methodSharedLinkloacalDisable) {
				errs[key] = fmt.Sprintf(`key "%s" cannot be used with the 'shared', 'link-local', or 'disabled' methods`, key)
				continue
			}
			// TODO
		case NM_SETTING_IP4_CONFIG_IGNORE_AUTO_ROUTES: // ignore
		case NM_SETTING_IP4_CONFIG_IGNORE_AUTO_DNS: // ignore
		case NM_SETTING_IP4_CONFIG_DHCP_CLIENT_ID: // ignore
		case NM_SETTING_IP4_CONFIG_DHCP_SEND_HOSTNAME: // ignore
		case NM_SETTING_IP4_CONFIG_DHCP_HOSTNAME: // ignore
		case NM_SETTING_IP4_CONFIG_NEVER_DEFAULT: // ignore
		case NM_SETTING_IP4_CONFIG_MAY_FAIL: // ignore
		}
	}

	// method := getSettingIp4ConfigMethodJSON(data)
	// switch method {
	// case NM_SETTING_IP4_CONFIG_METHOD_AUTO:
	// 	// removeSettingIp4ConfigAddresses(data)
	// case NM_SETTING_IP4_CONFIG_METHOD_LINK_LOCAL: // ignore
	// case NM_SETTING_IP4_CONFIG_METHOD_MANUAL:
	// case NM_SETTING_IP4_CONFIG_METHOD_SHARED: // ignore
	// case NM_SETTING_IP4_CONFIG_METHOD_DISABLED:
	// 	// removeSettingIp4ConfigDns(data)
	// 	// removeSettingIp4ConfigAddresses(data)
	// }
	return
}

// TODO Adder

// TODO use logic setter
func generalSetSettingIp4ConfigKeyJSON(data _ConnectionData, key, value string) {
	switch key {
	default:
		LOGGER.Error("generalSetSettingIp4ConfigKey: invalide key", key)
	case NM_SETTING_IP4_CONFIG_METHOD:
		logicSetSettingIp4ConfigMethod(data, value) // TODO
	case NM_SETTING_IP4_CONFIG_DNS:
		setSettingIp4ConfigDnsJSON(data, value)
	case NM_SETTING_IP4_CONFIG_DNS_SEARCH:
		setSettingIp4ConfigDnsSearchJSON(data, value)
	case NM_SETTING_IP4_CONFIG_ADDRESSES:
		setSettingIp4ConfigAddressesJSON(data, value)
	case NM_SETTING_IP4_CONFIG_ROUTES:
		setSettingIp4ConfigRoutesJSON(data, value)
	case NM_SETTING_IP4_CONFIG_IGNORE_AUTO_ROUTES:
		setSettingIp4ConfigIgnoreAutoRoutesJSON(data, value)
	case NM_SETTING_IP4_CONFIG_IGNORE_AUTO_DNS:
		setSettingIp4ConfigIgnoreAutoDnsJSON(data, value)
	case NM_SETTING_IP4_CONFIG_DHCP_CLIENT_ID:
		setSettingIp4ConfigDhcpClientIdJSON(data, value)
	case NM_SETTING_IP4_CONFIG_DHCP_SEND_HOSTNAME:
		setSettingIp4ConfigDhcpSendHostnameJSON(data, value)
	case NM_SETTING_IP4_CONFIG_DHCP_HOSTNAME:
		setSettingIp4ConfigDhcpHostnameJSON(data, value)
	case NM_SETTING_IP4_CONFIG_NEVER_DEFAULT:
		setSettingIp4ConfigNeverDefaultJSON(data, value)
	case NM_SETTING_IP4_CONFIG_MAY_FAIL:
		setSettingIp4ConfigMayFailJSON(data, value)
	}
	return
}

// TODO Logic setterJSON
func logicSetSettingIp4ConfigMethod(data _ConnectionData, value string) {
	switch value {
	case NM_SETTING_IP4_CONFIG_METHOD_AUTO:
		removeSettingIp4ConfigAddresses(data)
	case NM_SETTING_IP4_CONFIG_METHOD_LINK_LOCAL: // ignore
	case NM_SETTING_IP4_CONFIG_METHOD_MANUAL:
	case NM_SETTING_IP4_CONFIG_METHOD_SHARED: // ignore
	case NM_SETTING_IP4_CONFIG_METHOD_DISABLED:
		removeSettingIp4ConfigDns(data)
		removeSettingIp4ConfigAddresses(data)
	}
	setSettingIp4ConfigMethodJSON(data, value)
	return
}
