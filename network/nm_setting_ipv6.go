package main

// TODO

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

// Getter
func getSettingIp6ConfigMethod(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_METHOD, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_METHOD))
	return
}
func getSettingIp6ConfigDns(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_DNS, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_DNS))
	return
}
func getSettingIp6ConfigDnsSEARCH(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_DNS_SEARCH, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_DNS_SEARCH))
	return
}
func getSettingIp6ConfigAddresses(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_ADDRESSES, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_ADDRESSES))
	return
}
func getSettingIp6ConfigRoutes(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_ROUTES, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_ROUTES))
	return
}
func getSettingIp6ConfigIgnoreAUTOROUTES(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_IGNORE_AUTO_ROUTES, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_IGNORE_AUTO_ROUTES))
	return
}
func getSettingIp6ConfigIgnoreAUTODNS(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_IGNORE_AUTO_DNS, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_IGNORE_AUTO_DNS))
	return
}
func getSettingIp6ConfigNeverDEFAULT(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_NEVER_DEFAULT, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_NEVER_DEFAULT))
	return
}
func getSettingIp6ConfigMayFAIL(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_MAY_FAIL, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_MAY_FAIL))
	return
}
func getSettingIp6ConfigIp6PRIVACY(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_IP6_PRIVACY, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_IP6_PRIVACY))
	return
}
func getSettingIp6ConfigDhcpHOSTNAME(data _ConnectionData) (value string, err error) {
	value, err = getConnectionData(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_DHCP_HOSTNAME, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_DHCP_HOSTNAME))
	return
}

// Setter
func setSettingIp6ConfigMethod(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_METHOD, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_METHOD))
	return
}
func setSettingIp6ConfigDns(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_DNS, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_DNS))
	return
}
func setSettingIp6ConfigDnsSEARCH(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_DNS_SEARCH, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_DNS_SEARCH))
	return
}
func setSettingIp6ConfigAddresses(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_ADDRESSES, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_ADDRESSES))
	return
}
func setSettingIp6ConfigRoutes(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_ROUTES, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_ROUTES))
	return
}
func setSettingIp6ConfigIgnoreAUTOROUTES(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_IGNORE_AUTO_ROUTES, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_IGNORE_AUTO_ROUTES))
	return
}
func setSettingIp6ConfigIgnoreAUTODNS(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_IGNORE_AUTO_DNS, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_IGNORE_AUTO_DNS))
	return
}
func setSettingIp6ConfigNeverDEFAULT(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_NEVER_DEFAULT, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_NEVER_DEFAULT))
	return
}
func setSettingIp6ConfigMayFAIL(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_MAY_FAIL, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_MAY_FAIL))
	return
}
func setSettingIp6ConfigIp6PRIVACY(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_IP6_PRIVACY, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_IP6_PRIVACY))
	return
}
func setSettingIp6ConfigDhcpHOSTNAME(data _ConnectionData, value string) (err error) {
	err = setConnectionData(data, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_DHCP_HOSTNAME, value, getSettingIp6ConfigKeyType(NM_SETTING_IP6_CONFIG_DHCP_HOSTNAME))
	return
}

// TODO remove

// SettingIp6ConfigMethod NM_SETTING_IP6_CONFIG_METHOD
// SettingIp6ConfigDns NM_SETTING_IP6_CONFIG_DNS
// SettingIp6ConfigDnsSEARCH NM_SETTING_IP6_CONFIG_DNS_SEARCH
// SettingIp6ConfigAddresses NM_SETTING_IP6_CONFIG_ADDRESSES
// SettingIp6ConfigRoutes NM_SETTING_IP6_CONFIG_ROUTES
// SettingIp6ConfigIgnoreAUTOROUTES NM_SETTING_IP6_CONFIG_IGNORE_AUTO_ROUTES
// SettingIp6ConfigIgnoreAUTODNS NM_SETTING_IP6_CONFIG_IGNORE_AUTO_DNS
// SettingIp6ConfigNeverDEFAULT NM_SETTING_IP6_CONFIG_NEVER_DEFAULT
// SettingIp6ConfigMayFAIL NM_SETTING_IP6_CONFIG_MAY_FAIL
// SettingIp6ConfigIp6PRIVACY NM_SETTING_IP6_CONFIG_IP6_PRIVACY
// SettingIp6ConfigDhcpHOSTNAME NM_SETTING_IP6_CONFIG_DHCP_HOSTNAME
