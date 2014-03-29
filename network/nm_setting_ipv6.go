package main

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

// TODO Get available keys
func getSettingIp6ConfigAvailableKeys(data _ConnectionData) (keys []string) {
	method := getSettingIp6ConfigMethod(data)
	switch method {
	default:
		LOGGER.Error("ip6 config method is invalid:", method)
	case NM_SETTING_IP6_CONFIG_METHOD_IGNORE:
		keys = []string{
			NM_SETTING_IP6_CONFIG_METHOD,
		}
	case NM_SETTING_IP6_CONFIG_METHOD_AUTO:
		keys = []string{
			NM_SETTING_IP6_CONFIG_METHOD,
			NM_SETTING_IP6_CONFIG_DNS,
		}
	case NM_SETTING_IP6_CONFIG_METHOD_DHCP:
		keys = []string{
			NM_SETTING_IP6_CONFIG_METHOD,
			NM_SETTING_IP6_CONFIG_DNS,
		}
	case NM_SETTING_IP6_CONFIG_METHOD_LINK_LOCAL: // ignore
	case NM_SETTING_IP6_CONFIG_METHOD_MANUAL: // ignore
		keys = []string{
			NM_SETTING_IP6_CONFIG_METHOD,
			NM_SETTING_IP6_CONFIG_DNS,
			NM_SETTING_IP6_CONFIG_ADDRESSES,
		}
	case NM_SETTING_IP6_CONFIG_METHOD_SHARED: // ignore
	}
	return
}

// Get key type
func getSettingIp6ConfigKeyType(key string) (t ktype) {
	switch key {
	default:
		t = ktypeUnknown
	case NM_SETTING_IP6_CONFIG_METHOD:
		t = ktypeString
	case NM_SETTING_IP6_CONFIG_DHCP_HOSTNAME:
		t = ktypeString
	case NM_SETTING_IP6_CONFIG_DNS:
		t = ktypeArrayByte
	case NM_SETTING_IP6_CONFIG_DNS_SEARCH:
		t = ktypeArrayString
	case NM_SETTING_IP6_CONFIG_ADDRESSES:
		t = ktypeIpv6Addresses
	case NM_SETTING_IP6_CONFIG_ROUTES:
		t = ktypeIpv6Routes
	case NM_SETTING_IP6_CONFIG_IGNORE_AUTO_ROUTES:
		t = ktypeBoolean
	case NM_SETTING_IP6_CONFIG_IGNORE_AUTO_DNS:
		t = ktypeBoolean
	case NM_SETTING_IP6_CONFIG_NEVER_DEFAULT:
		t = ktypeBoolean
	case NM_SETTING_IP6_CONFIG_MAY_FAIL:
		t = ktypeBoolean
	case NM_SETTING_IP6_CONFIG_IP6_PRIVACY:
		t = ktypeInt32
	}
	return
}

// Get and set key's value generally
func generalGetSettingIp6ConfigKey(data _ConnectionData, key string) (value string) {
	switch key {
	default:
		LOGGER.Error("generalGetSettingIp6ConfigKey: invalide key", key)
	case NM_SETTING_IP6_CONFIG_METHOD:
		value = getSettingIp6ConfigMethod(data)
	case NM_SETTING_IP6_CONFIG_DNS:
		value = getSettingIp6ConfigDns(data)
	case NM_SETTING_IP6_CONFIG_DNS_SEARCH:
		value = getSettingIp6ConfigDnsSearch(data)
	case NM_SETTING_IP6_CONFIG_ADDRESSES:
		value = getSettingIp6ConfigAddresses(data)
	case NM_SETTING_IP6_CONFIG_ROUTES:
		value = getSettingIp6ConfigRoutes(data)
	case NM_SETTING_IP6_CONFIG_IGNORE_AUTO_ROUTES:
		value = getSettingIp6ConfigIgnoreAutoRoutes(data)
	case NM_SETTING_IP6_CONFIG_IGNORE_AUTO_DNS:
		value = getSettingIp6ConfigIgnoreAutoDns(data)
	case NM_SETTING_IP6_CONFIG_NEVER_DEFAULT:
		value = getSettingIp6ConfigNeverDefault(data)
	case NM_SETTING_IP6_CONFIG_MAY_FAIL:
		value = getSettingIp6ConfigMayFail(data)
	case NM_SETTING_IP6_CONFIG_IP6_PRIVACY:
		value = getSettingIp6ConfigIp6Privacy(data)
	case NM_SETTING_IP6_CONFIG_DHCP_HOSTNAME:
		value = getSettingIp6ConfigDhcpHostname(data)
	}
	return
}

// TODO use logic setter
func generalSetSettingIp6ConfigKey(data _ConnectionData, key, value string) {
	switch key {
	default:
		LOGGER.Error("generalSetSettingIp6ConfigKey: invalide key", key)
	case NM_SETTING_IP6_CONFIG_METHOD:
		setSettingIp6ConfigMethod(data, value)
	case NM_SETTING_IP6_CONFIG_DNS:
		setSettingIp6ConfigDns(data, value)
	case NM_SETTING_IP6_CONFIG_DNS_SEARCH:
		setSettingIp6ConfigDnsSearch(data, value)
	case NM_SETTING_IP6_CONFIG_ADDRESSES:
		setSettingIp6ConfigAddresses(data, value)
	case NM_SETTING_IP6_CONFIG_ROUTES:
		setSettingIp6ConfigRoutes(data, value)
	case NM_SETTING_IP6_CONFIG_IGNORE_AUTO_ROUTES:
		setSettingIp6ConfigIgnoreAutoRoutes(data, value)
	case NM_SETTING_IP6_CONFIG_IGNORE_AUTO_DNS:
		setSettingIp6ConfigIgnoreAutoDns(data, value)
	case NM_SETTING_IP6_CONFIG_NEVER_DEFAULT:
		setSettingIp6ConfigNeverDefault(data, value)
	case NM_SETTING_IP6_CONFIG_MAY_FAIL:
		setSettingIp6ConfigMayFail(data, value)
	case NM_SETTING_IP6_CONFIG_IP6_PRIVACY:
		setSettingIp6ConfigIp6Privacy(data, value)
	case NM_SETTING_IP6_CONFIG_DHCP_HOSTNAME:
		setSettingIp6ConfigDhcpHostname(data, value)
	}
	return
}

// Getter
func getSettingIp6ConfigMethod(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_METHOD, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_METHOD))
	return
}
func getSettingIp6ConfigDns(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_DNS, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_DNS))
	return
}
func getSettingIp6ConfigDnsSearch(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_DNS_SEARCH, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_DNS_SEARCH))
	return
}
func getSettingIp6ConfigAddresses(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_ADDRESSES, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_ADDRESSES))
	return
}
func getSettingIp6ConfigRoutes(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_ROUTES, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_ROUTES))
	return
}
func getSettingIp6ConfigIgnoreAutoRoutes(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_IGNORE_AUTO_ROUTES, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_IGNORE_AUTO_ROUTES))
	return
}
func getSettingIp6ConfigIgnoreAutoDns(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_IGNORE_AUTO_DNS, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_IGNORE_AUTO_DNS))
	return
}
func getSettingIp6ConfigNeverDefault(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_NEVER_DEFAULT, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_NEVER_DEFAULT))
	return
}
func getSettingIp6ConfigMayFail(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_MAY_FAIL, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_MAY_FAIL))
	return
}
func getSettingIp6ConfigIp6Privacy(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_IP6_PRIVACY, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_IP6_PRIVACY))
	return
}
func getSettingIp6ConfigDhcpHostname(data _ConnectionData) (value string) {
	value = getConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_DHCP_HOSTNAME, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_DHCP_HOSTNAME))
	return
}

// Setter
func setSettingIp6ConfigMethod(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_METHOD, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_METHOD))
}
func setSettingIp6ConfigDns(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_DNS, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_DNS))
}
func setSettingIp6ConfigDnsSearch(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_DNS_SEARCH, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_DNS_SEARCH))
}
func setSettingIp6ConfigAddresses(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_ADDRESSES, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_ADDRESSES))
}
func setSettingIp6ConfigRoutes(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_ROUTES, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_ROUTES))
}
func setSettingIp6ConfigIgnoreAutoRoutes(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_IGNORE_AUTO_ROUTES, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_IGNORE_AUTO_ROUTES))
}
func setSettingIp6ConfigIgnoreAutoDns(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_IGNORE_AUTO_DNS, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_IGNORE_AUTO_DNS))
}
func setSettingIp6ConfigNeverDefault(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_NEVER_DEFAULT, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_NEVER_DEFAULT))
}
func setSettingIp6ConfigMayFail(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_MAY_FAIL, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_MAY_FAIL))
}
func setSettingIp6ConfigIp6Privacy(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_IP6_PRIVACY, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_IP6_PRIVACY))
}
func setSettingIp6ConfigDhcpHostname(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_DHCP_HOSTNAME, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_DHCP_HOSTNAME))
}

// Remover
func removeSettingIp6ConfigMethod(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_METHOD)
}
func removeSettingIp6ConfigDns(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_DNS)
}
func removeSettingIp6ConfigDnsSearch(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_DNS_SEARCH)
}
func removeSettingIp6ConfigAddresses(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_ADDRESSES)
}
func removeSettingIp6ConfigRoutes(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_ROUTES)
}
func removeSettingIp6ConfigIgnoreAutoRoutes(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_IGNORE_AUTO_ROUTES)
}
func removeSettingIp6ConfigIgnoreAutoDns(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_IGNORE_AUTO_DNS)
}
func removeSettingIp6ConfigNeverDefault(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_NEVER_DEFAULT)
}
func removeSettingIp6ConfigMayFail(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_MAY_FAIL)
}
func removeSettingIp6ConfigIp6Privacy(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_IP6_PRIVACY)
}
func removeSettingIp6ConfigDhcpHostname(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_DHCP_HOSTNAME)
}
