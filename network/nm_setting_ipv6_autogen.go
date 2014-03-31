package main

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

// TODO tmp
func getSettingIp6ConfigMethod(data _ConnectionData) (value string) {
	value, _ = getConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_METHOD).(string)
	return
}
func setSettingIp6ConfigMethod(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_METHOD, value)
}

// Getter
func getSettingIp6ConfigMethodJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_METHOD, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_METHOD))
	return
}
func getSettingIp6ConfigDnsJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_DNS, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_DNS))
	return
}
func getSettingIp6ConfigDnsSearchJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_DNS_SEARCH, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_DNS_SEARCH))
	return
}
func getSettingIp6ConfigAddressesJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_ADDRESSES, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_ADDRESSES))
	return
}
func getSettingIp6ConfigRoutesJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_ROUTES, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_ROUTES))
	return
}
func getSettingIp6ConfigIgnoreAutoRoutesJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_IGNORE_AUTO_ROUTES, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_IGNORE_AUTO_ROUTES))
	return
}
func getSettingIp6ConfigIgnoreAutoDnsJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_IGNORE_AUTO_DNS, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_IGNORE_AUTO_DNS))
	return
}
func getSettingIp6ConfigNeverDefaultJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_NEVER_DEFAULT, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_NEVER_DEFAULT))
	return
}
func getSettingIp6ConfigMayFailJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_MAY_FAIL, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_MAY_FAIL))
	return
}
func getSettingIp6ConfigIp6PrivacyJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_IP6_PRIVACY, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_IP6_PRIVACY))
	return
}
func getSettingIp6ConfigDhcpHostnameJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_DHCP_HOSTNAME, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_DHCP_HOSTNAME))
	return
}

// Setter
func setSettingIp6ConfigMethodJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_METHOD, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_METHOD))
}
func setSettingIp6ConfigDnsJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_DNS, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_DNS))
}
func setSettingIp6ConfigDnsSearchJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_DNS_SEARCH, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_DNS_SEARCH))
}
func setSettingIp6ConfigAddressesJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_ADDRESSES, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_ADDRESSES))
}
func setSettingIp6ConfigRoutesJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_ROUTES, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_ROUTES))
}
func setSettingIp6ConfigIgnoreAutoRoutesJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_IGNORE_AUTO_ROUTES, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_IGNORE_AUTO_ROUTES))
}
func setSettingIp6ConfigIgnoreAutoDnsJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_IGNORE_AUTO_DNS, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_IGNORE_AUTO_DNS))
}
func setSettingIp6ConfigNeverDefaultJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_NEVER_DEFAULT, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_NEVER_DEFAULT))
}
func setSettingIp6ConfigMayFailJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_MAY_FAIL, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_MAY_FAIL))
}
func setSettingIp6ConfigIp6PrivacyJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_IP6_PRIVACY, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_IP6_PRIVACY))
}
func setSettingIp6ConfigDhcpHostnameJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_DHCP_HOSTNAME, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_DHCP_HOSTNAME))
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
