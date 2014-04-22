package main

// Virtual key names.

// 802-1x
const (
	NM_SETTING_VK_802_1X_EAP = "vk-eap"
)

// connection
const (
	NM_SETTING_VK_CONNECTION_PERMISSIONS = "vk-permissions"
)

// ipv4
const (
	NM_SETTING_VK_IP4_CONFIG_DNS               = "vk-dns"
	NM_SETTING_VK_IP4_CONFIG_ADDRESSES_ADDRESS = "vk-addresses-address"
	NM_SETTING_VK_IP4_CONFIG_ADDRESSES_MASK    = "vk-addresses-mask"
	NM_SETTING_VK_IP4_CONFIG_ADDRESSES_GATEWAY = "vk-addresses-gateway"
	NM_SETTING_VK_IP4_CONFIG_ROUTES_ADDRESS    = "vk-routes-address"
	NM_SETTING_VK_IP4_CONFIG_ROUTES_MASK       = "vk-routes-mask"
	NM_SETTING_VK_IP4_CONFIG_ROUTES_NEXTHOP    = "vk-routes-nexthop"
	NM_SETTING_VK_IP4_CONFIG_ROUTES_METRIC     = "vk-routes-metric"
)

// ipv6
const (
	NM_SETTING_VK_IP6_CONFIG_DNS               = "vk-dns"
	NM_SETTING_VK_IP6_CONFIG_ADDRESSES_ADDRESS = "vk-addresses-address"
	NM_SETTING_VK_IP6_CONFIG_ADDRESSES_PREFIX  = "vk-addresses-prefix"
	NM_SETTING_VK_IP6_CONFIG_ADDRESSES_GATEWAY = "vk-addresses-gateway"
	NM_SETTING_VK_IP6_CONFIG_ROUTES_ADDRESS    = "vk-routes-address"
	NM_SETTING_VK_IP6_CONFIG_ROUTES_PREFIX     = "vk-routes-prefix"
	NM_SETTING_VK_IP6_CONFIG_ROUTES_NEXTHOP    = "vk-routes-nexthop"
	NM_SETTING_VK_IP6_CONFIG_ROUTES_METRIC     = "vk-routes-metric"
)

// wireless security
const (
	NM_SETTING_VK_WIRELESS_SECURITY_KEY_MGMT = "vk-key-mgmt"
)

// VirtualKey stores virtual key info for each fields.
type VirtualKey struct {
	Name         string
	Type         ktype
	RelatedField string
	RelatedKey   string
	Available    bool // check if is used by front-end
	Optional     bool // if key is optional, will ignore error for it
}

func isVirtualKey(field, key string) bool {
	if isStringInArray(key, getVirtualKeysOfField(field)) {
		return true
	}
	return false
}

func getVirtualKeysOfField(field string) (vks []string) {
	for _, vk := range virtualKeys {
		if vk.RelatedField == field {
			vks = append(vks, vk.Name)
		}
	}
	// Logger.Debug("getVirtualKeysOfField: filed:", field, vks) // TODO test
	return
}

func getSettingVkKeyType(field, key string) (t ktype) {
	t = ktypeUnknown
	for _, vk := range virtualKeys {
		if vk.RelatedField == field && vk.Name == key {
			t = vk.Type
		}
	}
	return
}

func generalGetSettingVkAvailableValues(field, key string) (values []string) {
	switch field {
	case field8021x:
		switch key {
		case NM_SETTING_VK_802_1X_EAP:
			values, _ = getSetting8021xAvailableValues(nil, NM_SETTING_802_1X_EAP)
		}
	case fieldConnection:
	case fieldIPv4:
	case fieldIPv6:
	case fieldWired:
	case fieldWireless:
	case fieldWirelessSecurity:
		switch key {
		case NM_SETTING_VK_WIRELESS_SECURITY_KEY_MGMT:
			values = []string{"none", "wep", "wpa-psk", "wpa-eap"}
		}
	case fieldPppoe:
	case fieldPpp:
	}
	return
}

func getRelatedAvailableVirtualKeys(field, key string) (vks []string) {
	for _, vk := range virtualKeys {
		if vk.RelatedField == field && vk.RelatedKey == key && vk.Available {
			vks = append(vks, vk.Name)
		}
	}
	return
}

// get related virtual key(s) for target key
func getRelatedVirtualKeys(field, key string) (vks []string) {
	for _, vk := range virtualKeys {
		if vk.RelatedField == field && vk.RelatedKey == key {
			vks = append(vks, vk.Name)
		}
	}
	return
}

func isOptionalChildVirtualKeys(field, vkey string) (optional bool) {
	for _, vk := range virtualKeys {
		if vk.RelatedField == field && vk.Name == vkey {
			optional = vk.Optional
		}
	}
	return
}

func doGetOrNewSettingIp4ConfigAddresses(data _ConnectionData) (addresses [][]uint32) {
	addresses = getSettingIp4ConfigAddresses(data)
	if len(addresses) == 0 {
		addresses = make([][]uint32, 1)
		addresses[0] = make([]uint32, 3)
	}
	if len(addresses[0]) != 3 {
		addresses[0] = make([]uint32, 3)
	}
	return
}

func doIsSettingIp4ConfigAddressesEmpty(data _ConnectionData) bool {
	addresses := getSettingIp4ConfigAddresses(data)
	if len(addresses) == 0 {
		return true
	}
	if len(addresses[0]) != 3 {
		return true
	}
	return false
}

// Getter
func getSettingVk8021xEap(data _ConnectionData) (eap string) {
	eaps := getSetting8021xEap(data)
	if len(eaps) == 0 {
		Logger.Error("eap value is empty")
		return
	}
	eap = eaps[0]
	return
}
func getSettingVkConnectionPermissions(data _ConnectionData) (value bool) {
	// TODO
	// value = getSettingConnectionPermissions(data)
	return
}
func getSettingVkIp4ConfigDns(data _ConnectionData) (value string) {
	dns := getSettingIp4ConfigDns(data)
	if len(dns) == 0 {
		return
	}
	value = convertIpv4AddressToString(dns[0])
	return
}
func getSettingVkIp4ConfigAddressesAddress(data _ConnectionData) (value string) {
	if doIsSettingIp4ConfigAddressesEmpty(data) {
		return
	}
	addresses := getSettingIp4ConfigAddresses(data)
	value = convertIpv4AddressToString(addresses[0][0])
	return
}
func getSettingVkIp4ConfigAddressesMask(data _ConnectionData) (value string) {
	if doIsSettingIp4ConfigAddressesEmpty(data) {
		return
	}
	addresses := getSettingIp4ConfigAddresses(data)
	value = convertIpv4PrefixToNetMask(addresses[0][1])
	return
}
func getSettingVkIp4ConfigAddressesGateway(data _ConnectionData) (value string) {
	if doIsSettingIp4ConfigAddressesEmpty(data) {
		return
	}
	addresses := getSettingIp4ConfigAddresses(data)
	value = convertIpv4AddressToStringNoZero(addresses[0][2])
	return
}
func getSettingVkIp4ConfigRoutesAddress(data _ConnectionData) (value string) {
	// TODO
	// value := getSettingIp4ConfigRoutesAddress(data)
	return
}
func getSettingVkIp4ConfigRoutesMask(data _ConnectionData) (value string) {
	// TODO
	// value := getSettingIp4ConfigRoutesMask(data)
	return
}
func getSettingVkIp4ConfigRoutesNexthop(data _ConnectionData) (value string) {
	// TODO
	// value := getSettingIp4ConfigRoutesNexthop(data)
	return
}
func getSettingVkIp4ConfigRoutesMetric(data _ConnectionData) (value string) {
	// TODO
	// value := getSettingIp4ConfigRoutesMetric(data)
	return
}
func getSettingVkIp6ConfigDns(data _ConnectionData) (value string) {
	// TODO
	// value := getSettingIp6ConfigDns(data)
	return
}
func getSettingVkIp6ConfigAddressesAddress(data _ConnectionData) (value string) {
	// TODO
	// value := getSettingIp6ConfigAddressesAddress(data)
	return
}
func getSettingVkIp6ConfigAddressesPrefix(data _ConnectionData) (value string) {
	// TODO
	// value := getSettingIp6ConfigAddressesPrefix(data)
	return
}
func getSettingVkIp6ConfigAddressesGateway(data _ConnectionData) (value string) {
	// TODO
	// value := getSettingIp6ConfigAddressesGateway(data)
	return
}
func getSettingVkIp6ConfigRoutesAddress(data _ConnectionData) (value string) {
	// TODO
	// value := getSettingIp6ConfigRoutesAddress(data)
	return
}
func getSettingVkIp6ConfigRoutesPrefix(data _ConnectionData) (value string) {
	// TODO
	// value := getSettingIp6ConfigRoutesPrefix(data)
	return
}
func getSettingVkIp6ConfigRoutesNexthop(data _ConnectionData) (value string) {
	// TODO
	// value := getSettingIp6ConfigRoutesNexthop(data)
	return
}
func getSettingVkIp6ConfigRoutesMetric(data _ConnectionData) (value string) {
	// TODO
	// value := getSettingIp6ConfigRoutesMetric(data)
	return
}
func getSettingVkWirelessSecurityKeyMgmt(data _ConnectionData) (value string) {
	if !isSettingFieldExists(data, fieldWirelessSecurity) {
		value = "none"
		return
	}
	keyMgmt := getSettingWirelessSecurityKeyMgmt(data)
	switch keyMgmt {
	case "none":
		value = "wep"
	case "wpa-psk":
		value = "wpa-psk"
	case "wpa-eap":
		value = "wpa-eap"
	}
	return
}

// Logic setter, all virtual keys has a logic setter
func logicSetSettingVk8021xEap(data _ConnectionData, value string) {
	logicSetSetting8021xEap(data, []string{value})
}
func logicSetSettingVkConnectionPermissions(data _ConnectionData, value bool) {
	// TODO
	// setSettingConnectionPermissionsJSON(data)
}
func logicSetSettingVkIp4ConfigDns(data _ConnectionData, value string) {
	dns := getSettingIp4ConfigDns(data)
	if len(dns) == 0 {
		dns = make([]uint32, 1)
	}
	dns[0] = convertIpv4AddressToUint32(value)
	if dns[0] != 0 {
		setSettingIp4ConfigDns(data, dns)
	} else {
		removeSettingIp4ConfigDns(data)
	}
}
func logicSetSettingVkIp4ConfigAddressesAddress(data _ConnectionData, value string) {
	addresses := doGetOrNewSettingIp4ConfigAddresses(data)
	addr := addresses[0]
	addr[0] = convertIpv4AddressToUint32(value)
	if !isUint32ArrayEmpty(addr) {
		setSettingIp4ConfigAddresses(data, addresses)
	} else {
		removeSettingIp4ConfigAddresses(data)
	}
}
func logicSetSettingVkIp4ConfigAddressesMask(data _ConnectionData, value string) {
	addresses := doGetOrNewSettingIp4ConfigAddresses(data)
	addr := addresses[0]
	addr[1] = convertIpv4NetMaskToPrefix(value)
	if !isUint32ArrayEmpty(addr) {
		setSettingIp4ConfigAddresses(data, addresses)
	} else {
		removeSettingIp4ConfigAddresses(data)
	}
}
func logicSetSettingVkIp4ConfigAddressesGateway(data _ConnectionData, value string) {
	if len(value) == 0 {
		value = ipv4Zero
	}
	addresses := doGetOrNewSettingIp4ConfigAddresses(data)
	addr := addresses[0]
	addr[2] = convertIpv4AddressToUint32(value)
	if !isUint32ArrayEmpty(addr) {
		setSettingIp4ConfigAddresses(data, addresses)
	} else {
		removeSettingIp4ConfigAddresses(data)
	}
}
func logicSetSettingVkIp4ConfigRoutesAddress(data _ConnectionData, value string) {
	// TODO
	// setSettingIp4ConfigRoutesAddressJSON(data)
}
func logicSetSettingVkIp4ConfigRoutesMask(data _ConnectionData, value string) {
	// TODO
	// setSettingIp4ConfigRoutesMaskJSON(data)
}
func logicSetSettingVkIp4ConfigRoutesNexthop(data _ConnectionData, value string) {
	// TODO
	// setSettingIp4ConfigRoutesNexthopJSON(data)
}
func logicSetSettingVkIp4ConfigRoutesMetric(data _ConnectionData, value string) {
	// TODO
	// setSettingIp4ConfigRoutesMetricJSON(data)
}
func logicSetSettingVkIp6ConfigDns(data _ConnectionData, value string) {
	// TODO
	// setSettingIp6ConfigDnsJSON(data)
}
func logicSetSettingVkIp6ConfigAddressesAddress(data _ConnectionData, value string) {
	// TODO
	// setSettingIp6ConfigAddressesAddressJSON(data)
}
func logicSetSettingVkIp6ConfigAddressesPrefix(data _ConnectionData, value string) {
	// TODO
	// setSettingIp6ConfigAddressesPrefixJSON(data)
}
func logicSetSettingVkIp6ConfigAddressesGateway(data _ConnectionData, value string) {
	// TODO
	// setSettingIp6ConfigAddressesGatewayJSON(data)
}
func logicSetSettingVkIp6ConfigRoutesAddress(data _ConnectionData, value string) {
	// TODO
	// setSettingIp6ConfigRoutesAddressJSON(data)
}
func logicSetSettingVkIp6ConfigRoutesPrefix(data _ConnectionData, value string) {
	// TODO
	// setSettingIp6ConfigRoutesPrefixJSON(data)
}
func logicSetSettingVkIp6ConfigRoutesNexthop(data _ConnectionData, value string) {
	// TODO
	// setSettingIp6ConfigRoutesNexthopJSON(data)
}
func logicSetSettingVkIp6ConfigRoutesMetric(data _ConnectionData, value string) {
	// TODO
	// setSettingIp6ConfigRoutesMetricJSON(data)
}
func logicSetSettingVkWirelessSecurityKeyMgmt(data _ConnectionData, value string) {
	switch value {
	default:
		Logger.Error("invalid value", value)
	case "none":
		removeSettingField(data, fieldWirelessSecurity)
		removeSettingField(data, field8021x)
	case "wep":
		addSettingField(data, fieldWirelessSecurity)
		removeSettingField(data, field8021x)
		setSettingWirelessSecurityKeyMgmt(data, "none")
		setSettingWirelessSecurityAuthAlg(data, "open")
		setSettingWirelessSecurityWepKeyType(data, 1)
	case "wpa-psk":
		addSettingField(data, fieldWirelessSecurity)
		removeSettingField(data, field8021x)
		setSettingWirelessSecurityKeyMgmt(data, "wpa-psk")
	case "wpa-eap":
		addSettingField(data, fieldWirelessSecurity)
		addSettingField(data, field8021x)
		setSettingWirelessSecurityKeyMgmt(data, "wpa-eap")
		logicSetSetting8021xEap(data, []string{"tls"})
		// TODO
	}
}
