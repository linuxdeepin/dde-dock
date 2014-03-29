package main

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

// TODO Adder, Remover

// Get key type
func getSettingIp4ConfigKeyType(key string) (t ktype) {
	switch key {
	default:
		t = ktypeUnknown
	case NM_SETTING_IP4_CONFIG_METHOD:
		t = ktypeString
	case NM_SETTING_IP4_CONFIG_DNS:
		t = ktypeArrayUint32
	case NM_SETTING_IP4_CONFIG_DNS_SEARCH:
		t = ktypeString
	case NM_SETTING_IP4_CONFIG_ADDRESSES:
		t = ktypeArrayArrayUint32
	case NM_SETTING_IP4_CONFIG_ROUTES:
		t = ktypeArrayArrayUint32
	case NM_SETTING_IP4_CONFIG_IGNORE_AUTO_ROUTES:
		t = ktypeBoolean
	case NM_SETTING_IP4_CONFIG_IGNORE_AUTO_DNS:
		t = ktypeBoolean
	case NM_SETTING_IP4_CONFIG_DHCP_CLIENT_ID:
		t = ktypeString
	case NM_SETTING_IP4_CONFIG_DHCP_SEND_HOSTNAME:
		t = ktypeBoolean
	case NM_SETTING_IP4_CONFIG_DHCP_HOSTNAME:
		t = ktypeString
	case NM_SETTING_IP4_CONFIG_NEVER_DEFAULT:
		t = ktypeBoolean
	case NM_SETTING_IP4_CONFIG_MAY_FAIL:
		t = ktypeBoolean
	}
	return
}

// Get and set key's value generally
func generalGetSettingIp4ConfigKey(data _ConnectionData, key string) (value string) {
	switch key {
	default:
		LOGGER.Error("generalGetSettingIp4ConfigKey: invalide key", key)
	case NM_SETTING_IP4_CONFIG_METHOD:
		value = getSettingIp4ConfigMethod(data)
	case NM_SETTING_IP4_CONFIG_DNS:
		value = getSettingIp4ConfigDns(data)
	case NM_SETTING_IP4_CONFIG_DNS_SEARCH:
		value = getSettingIp4ConfigDnsSearch(data)
	case NM_SETTING_IP4_CONFIG_ADDRESSES:
		value = getSettingIp4ConfigAddresses(data)
	case NM_SETTING_IP4_CONFIG_ROUTES:
		value = getSettingIp4ConfigRoutes(data)
	case NM_SETTING_IP4_CONFIG_IGNORE_AUTO_ROUTES:
		value = getSettingIp4ConfigIgnoreAutoRoutes(data)
	case NM_SETTING_IP4_CONFIG_IGNORE_AUTO_DNS:
		value = getSettingIp4ConfigIgnoreAutoDns(data)
	case NM_SETTING_IP4_CONFIG_DHCP_CLIENT_ID:
		value = getSettingIp4ConfigDhcpClientId(data)
	case NM_SETTING_IP4_CONFIG_DHCP_SEND_HOSTNAME:
		value = getSettingIp4ConfigDhcpSendHostname(data)
	case NM_SETTING_IP4_CONFIG_DHCP_HOSTNAME:
		value = getSettingIp4ConfigDhcpHostname(data)
	case NM_SETTING_IP4_CONFIG_NEVER_DEFAULT:
		value = getSettingIp4ConfigNeverDefault(data)
	case NM_SETTING_IP4_CONFIG_MAY_FAIL:
		value = getSettingIp4ConfigMayFail(data)
	}
	return
}

// TODO use logic setter
func generalSetSettingIp4ConfigKey(data _ConnectionData, key, value string) {
	switch key {
	default:
		LOGGER.Error("generalSetSettingIp4ConfigKey: invalide key", key)
	case NM_SETTING_IP4_CONFIG_METHOD:
		setSettingIp4ConfigMethod(data, value) // TODO
	case NM_SETTING_IP4_CONFIG_DNS:
		setSettingIp4ConfigDns(data, value)
	case NM_SETTING_IP4_CONFIG_DNS_SEARCH:
		setSettingIp4ConfigDnsSearch(data, value)
	case NM_SETTING_IP4_CONFIG_ADDRESSES:
		setSettingIp4ConfigAddresses(data, value)
	case NM_SETTING_IP4_CONFIG_ROUTES:
		setSettingIp4ConfigRoutes(data, value)
	case NM_SETTING_IP4_CONFIG_IGNORE_AUTO_ROUTES:
		setSettingIp4ConfigIgnoreAutoRoutes(data, value)
	case NM_SETTING_IP4_CONFIG_IGNORE_AUTO_DNS:
		setSettingIp4ConfigIgnoreAutoDns(data, value)
	case NM_SETTING_IP4_CONFIG_DHCP_CLIENT_ID:
		setSettingIp4ConfigDhcpClientId(data, value)
	case NM_SETTING_IP4_CONFIG_DHCP_SEND_HOSTNAME:
		setSettingIp4ConfigDhcpSendHostname(data, value)
	case NM_SETTING_IP4_CONFIG_DHCP_HOSTNAME:
		setSettingIp4ConfigDhcpHostname(data, value)
	case NM_SETTING_IP4_CONFIG_NEVER_DEFAULT:
		setSettingIp4ConfigNeverDefault(data, value)
	case NM_SETTING_IP4_CONFIG_MAY_FAIL:
		setSettingIp4ConfigMayFail(data, value)
	}
	return
}

// TODO Logic setter
func logicSetSettingIp4ConfigMethod(data _ConnectionData, value string) (err error) {
	setSettingIp4ConfigMethod(data, value)
	switch value {
	case NM_SETTING_IP4_CONFIG_METHOD_AUTO:

	case NM_SETTING_IP4_CONFIG_METHOD_LINK_LOCAL: // ignore
	case NM_SETTING_IP4_CONFIG_METHOD_MANUAL:
	case NM_SETTING_IP4_CONFIG_METHOD_SHARED: // ignore
	case NM_SETTING_IP4_CONFIG_METHOD_DISABLED: // ignore
	default:
		// TODO
		// err =
	}
	return
}
func logicSetSettingIp4ConfigDns(data _ConnectionData, value string) (err error) {
	setSettingIp4ConfigDns(data, value)
	return
}
func logicSetSettingIp4ConfigDnsSearch(data _ConnectionData, value string) (err error) {
	setSettingIp4ConfigDnsSearch(data, value)
	return
}
func logicSetSettingIp4ConfigAddresses(data _ConnectionData, value string) (err error) {
	setSettingIp4ConfigAddresses(data, value)
	return
}
func logicSetSettingIp4ConfigRoutes(data _ConnectionData, value string) (err error) {
	setSettingIp4ConfigRoutes(data, value)
	return
}
func logicSetSettingIp4ConfigIgnoreAutoRoutes(data _ConnectionData, value string) (err error) {
	setSettingIp4ConfigIgnoreAutoRoutes(data, value)
	return
}
func logicSetSettingIp4ConfigIgnoreAutoDns(data _ConnectionData, value string) (err error) {
	setSettingIp4ConfigIgnoreAutoDns(data, value)
	return
}
func logicSetSettingIp4ConfigDhcpClientId(data _ConnectionData, value string) (err error) {
	setSettingIp4ConfigDhcpClientId(data, value)
	return
}
func logicSetSettingIp4ConfigDhcpSendHostname(data _ConnectionData, value string) (err error) {
	setSettingIp4ConfigDhcpSendHostname(data, value)
	return
}
func logicSetSettingIp4ConfigDhcpHostname(data _ConnectionData, value string) (err error) {
	setSettingIp4ConfigDhcpHostname(data, value)
	return
}
func logicSetSettingIp4ConfigNeverDefault(data _ConnectionData, value string) (err error) {
	setSettingIp4ConfigNeverDefault(data, value)
	return
}
func logicSetSettingIp4ConfigMayFail(data _ConnectionData, value string) (err error) {
	setSettingIp4ConfigMayFail(data, value)
	return
}

// Getter
func getSettingIp4ConfigMethod(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_METHOD, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_METHOD))
	return
}
func getSettingIp4ConfigDns(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_DNS, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_DNS))
	return
}
func getSettingIp4ConfigDnsSearch(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_DNS_SEARCH, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_DNS_SEARCH))
	return
}
func getSettingIp4ConfigAddresses(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_ADDRESSES, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_ADDRESSES))
	return
}
func getSettingIp4ConfigRoutes(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_ROUTES, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_ROUTES))
	return
}
func getSettingIp4ConfigIgnoreAutoRoutes(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_IGNORE_AUTO_ROUTES, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_IGNORE_AUTO_ROUTES))
	return
}
func getSettingIp4ConfigIgnoreAutoDns(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_IGNORE_AUTO_DNS, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_IGNORE_AUTO_DNS))
	return
}
func getSettingIp4ConfigDhcpClientId(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_DHCP_CLIENT_ID, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_DHCP_CLIENT_ID))
	return
}
func getSettingIp4ConfigDhcpSendHostname(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_DHCP_SEND_HOSTNAME, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_DHCP_SEND_HOSTNAME))
	return
}
func getSettingIp4ConfigDhcpHostname(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_DHCP_HOSTNAME, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_DHCP_HOSTNAME))
	return
}
func getSettingIp4ConfigNeverDefault(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_NEVER_DEFAULT, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_NEVER_DEFAULT))
	return
}
func getSettingIp4ConfigMayFail(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_MAY_FAIL, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_MAY_FAIL))
	return
}

// Setter
func setSettingIp4ConfigMethod(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_METHOD, value, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_METHOD))
}
func setSettingIp4ConfigDns(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_DNS, value, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_DNS))
}
func setSettingIp4ConfigDnsSearch(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_DNS_SEARCH, value, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_DNS_SEARCH))
}
func setSettingIp4ConfigAddresses(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_ADDRESSES, value, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_ADDRESSES))
}
func setSettingIp4ConfigRoutes(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_ROUTES, value, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_ROUTES))
}
func setSettingIp4ConfigIgnoreAutoRoutes(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_IGNORE_AUTO_ROUTES, value, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_IGNORE_AUTO_ROUTES))
}
func setSettingIp4ConfigIgnoreAutoDns(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_IGNORE_AUTO_DNS, value, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_IGNORE_AUTO_DNS))
}
func setSettingIp4ConfigDhcpClientId(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_DHCP_CLIENT_ID, value, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_DHCP_CLIENT_ID))
}
func setSettingIp4ConfigDhcpSendHostname(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_DHCP_SEND_HOSTNAME, value, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_DHCP_SEND_HOSTNAME))
}
func setSettingIp4ConfigDhcpHostname(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_DHCP_HOSTNAME, value, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_DHCP_HOSTNAME))
}
func setSettingIp4ConfigNeverDefault(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_NEVER_DEFAULT, value, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_NEVER_DEFAULT))
}
func setSettingIp4ConfigMayFail(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_MAY_FAIL, value, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_MAY_FAIL))
}

// Remover
func removeSettingIp4ConfigMethod(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_METHOD)
}
func removeSettingIp4ConfigDns(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_DNS)
}
func removeSettingIp4ConfigDnsSearch(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_DNS_SEARCH)
}
func removeSettingIp4ConfigAddresses(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_ADDRESSES)
}
func removeSettingIp4ConfigRoutes(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_ROUTES)
}
func removeSettingIp4ConfigIgnoreAutoRoutes(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_IGNORE_AUTO_ROUTES)
}
func removeSettingIp4ConfigIgnoreAutoDns(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_IGNORE_AUTO_DNS)
}
func removeSettingIp4ConfigDhcpClientId(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_DHCP_CLIENT_ID)
}
func removeSettingIp4ConfigDhcpSendHostname(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_DHCP_SEND_HOSTNAME)
}
func removeSettingIp4ConfigDhcpHostname(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_DHCP_HOSTNAME)
}
func removeSettingIp4ConfigNeverDefault(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_NEVER_DEFAULT)
}
func removeSettingIp4ConfigMayFail(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_MAY_FAIL)
}
