package main

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

// TODO tmp
func getSettingIp4ConfigMethod(data _ConnectionData) (value string) {
	value, _ = getConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_METHOD).(string)
	return
}
func setSettingIp4ConfigMethod(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_METHOD, value)
}

// Getter
func getSettingIp4ConfigMethodJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_METHOD, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_METHOD))
	return
}
func getSettingIp4ConfigDnsJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_DNS, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_DNS))
	return
}
func getSettingIp4ConfigDnsSearchJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_DNS_SEARCH, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_DNS_SEARCH))
	return
}
func getSettingIp4ConfigAddressesJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_ADDRESSES, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_ADDRESSES))
	return
}
func getSettingIp4ConfigRoutesJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_ROUTES, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_ROUTES))
	return
}
func getSettingIp4ConfigIgnoreAutoRoutesJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_IGNORE_AUTO_ROUTES, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_IGNORE_AUTO_ROUTES))
	return
}
func getSettingIp4ConfigIgnoreAutoDnsJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_IGNORE_AUTO_DNS, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_IGNORE_AUTO_DNS))
	return
}
func getSettingIp4ConfigDhcpClientIdJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_DHCP_CLIENT_ID, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_DHCP_CLIENT_ID))
	return
}
func getSettingIp4ConfigDhcpSendHostnameJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_DHCP_SEND_HOSTNAME, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_DHCP_SEND_HOSTNAME))
	return
}
func getSettingIp4ConfigDhcpHostnameJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_DHCP_HOSTNAME, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_DHCP_HOSTNAME))
	return
}
func getSettingIp4ConfigNeverDefaultJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_NEVER_DEFAULT, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_NEVER_DEFAULT))
	return
}
func getSettingIp4ConfigMayFailJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_MAY_FAIL, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_MAY_FAIL))
	return
}

// Setter
func setSettingIp4ConfigMethodJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_METHOD, value, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_METHOD))
}
func setSettingIp4ConfigDnsJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_DNS, value, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_DNS))
}
func setSettingIp4ConfigDnsSearchJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_DNS_SEARCH, value, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_DNS_SEARCH))
}
func setSettingIp4ConfigAddressesJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_ADDRESSES, value, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_ADDRESSES))
}
func setSettingIp4ConfigRoutesJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_ROUTES, value, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_ROUTES))
}
func setSettingIp4ConfigIgnoreAutoRoutesJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_IGNORE_AUTO_ROUTES, value, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_IGNORE_AUTO_ROUTES))
}
func setSettingIp4ConfigIgnoreAutoDnsJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_IGNORE_AUTO_DNS, value, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_IGNORE_AUTO_DNS))
}
func setSettingIp4ConfigDhcpClientIdJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_DHCP_CLIENT_ID, value, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_DHCP_CLIENT_ID))
}
func setSettingIp4ConfigDhcpSendHostnameJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_DHCP_SEND_HOSTNAME, value, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_DHCP_SEND_HOSTNAME))
}
func setSettingIp4ConfigDhcpHostnameJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_DHCP_HOSTNAME, value, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_DHCP_HOSTNAME))
}
func setSettingIp4ConfigNeverDefaultJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_NEVER_DEFAULT, value, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_NEVER_DEFAULT))
}
func setSettingIp4ConfigMayFailJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_MAY_FAIL, value, getSettingIp4ConfigKeyType(NM_SETTING_IP4_CONFIG_MAY_FAIL))
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
