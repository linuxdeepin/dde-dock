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

// TODO Check whether the values are correct
func checkSettingIp6ConfigValues(data _ConnectionData) (errs map[string]string) {
	errs = make(map[string]string)
	return
}

// Set JSON value generally
// TODO use logic setter
func generalSetSettingIp6ConfigKeyJSON(data _ConnectionData, key, value string) {
	switch key {
	default:
		LOGGER.Error("generalSetSettingIp6ConfigKey: invalide key", key)
	case NM_SETTING_IP6_CONFIG_METHOD:
		setSettingIp6ConfigMethodJSON(data, value)
	case NM_SETTING_IP6_CONFIG_DNS:
		setSettingIp6ConfigDnsJSON(data, value)
	case NM_SETTING_IP6_CONFIG_DNS_SEARCH:
		setSettingIp6ConfigDnsSearchJSON(data, value)
	case NM_SETTING_IP6_CONFIG_ADDRESSES:
		setSettingIp6ConfigAddressesJSON(data, value)
	case NM_SETTING_IP6_CONFIG_ROUTES:
		setSettingIp6ConfigRoutesJSON(data, value)
	case NM_SETTING_IP6_CONFIG_IGNORE_AUTO_ROUTES:
		setSettingIp6ConfigIgnoreAutoRoutesJSON(data, value)
	case NM_SETTING_IP6_CONFIG_IGNORE_AUTO_DNS:
		setSettingIp6ConfigIgnoreAutoDnsJSON(data, value)
	case NM_SETTING_IP6_CONFIG_NEVER_DEFAULT:
		setSettingIp6ConfigNeverDefaultJSON(data, value)
	case NM_SETTING_IP6_CONFIG_MAY_FAIL:
		setSettingIp6ConfigMayFailJSON(data, value)
	case NM_SETTING_IP6_CONFIG_IP6_PRIVACY:
		setSettingIp6ConfigIp6PrivacyJSON(data, value)
	case NM_SETTING_IP6_CONFIG_DHCP_HOSTNAME:
		setSettingIp6ConfigDhcpHostnameJSON(data, value)
	}
	return
}
