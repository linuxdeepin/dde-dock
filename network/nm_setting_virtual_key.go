package main

// Virtual keys for each fields.

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

func isVirtualKey(field, key string) bool {
	if isStringInArray(key, getFieldVirtualKeys(field)) {
		return true
	}
	return false
}

func getFieldVirtualKeys(field string) (vks []string) {
	switch field {
	case field8021x:
	case fieldConnection:
		vks = []string{NM_SETTING_VK_CONNECTION_PERMISSIONS}
	case fieldIPv4:
		vks = []string{
			NM_SETTING_VK_IP4_CONFIG_DNS,
			NM_SETTING_VK_IP4_CONFIG_ADDRESSES_ADDRESS,
			NM_SETTING_VK_IP4_CONFIG_ADDRESSES_MASK,
			NM_SETTING_VK_IP4_CONFIG_ADDRESSES_GATEWAY,
			NM_SETTING_VK_IP4_CONFIG_ROUTES_ADDRESS,
			NM_SETTING_VK_IP4_CONFIG_ROUTES_MASK,
			NM_SETTING_VK_IP4_CONFIG_ROUTES_NEXTHOP,
			NM_SETTING_VK_IP4_CONFIG_ROUTES_METRIC,
		}
	case fieldIPv6:
		vks = []string{
			NM_SETTING_VK_IP6_CONFIG_DNS,
			NM_SETTING_VK_IP6_CONFIG_ADDRESSES_ADDRESS,
			NM_SETTING_VK_IP6_CONFIG_ADDRESSES_PREFIX,
			NM_SETTING_VK_IP6_CONFIG_ADDRESSES_GATEWAY,
			NM_SETTING_VK_IP6_CONFIG_ROUTES_ADDRESS,
			NM_SETTING_VK_IP6_CONFIG_ROUTES_PREFIX,
			NM_SETTING_VK_IP6_CONFIG_ROUTES_NEXTHOP,
			NM_SETTING_VK_IP6_CONFIG_ROUTES_METRIC,
		}
	case fieldWired:
	case fieldWireless:
	case fieldWirelessSecurity:
		vks = []string{NM_SETTING_VK_WIRELESS_SECURITY_KEY_MGMT}
	}
	return
}

func getSettingVkKeyType(field, key string) (t ktype) {
	t = ktypeUnknown
	switch field {
	case field8021x:
	case fieldConnection:
		switch key {
		case NM_SETTING_VK_CONNECTION_PERMISSIONS:
			t = ktypeBoolean
		}
	case fieldIPv4:
		switch key {
		case NM_SETTING_VK_IP4_CONFIG_DNS:
			t = ktypeString
		case NM_SETTING_VK_IP4_CONFIG_ADDRESSES_ADDRESS:
			t = ktypeString
		case NM_SETTING_VK_IP4_CONFIG_ADDRESSES_MASK:
			t = ktypeString
		case NM_SETTING_VK_IP4_CONFIG_ADDRESSES_GATEWAY:
			t = ktypeString
		case NM_SETTING_VK_IP4_CONFIG_ROUTES_ADDRESS:
			t = ktypeString
		case NM_SETTING_VK_IP4_CONFIG_ROUTES_MASK:
			t = ktypeString
		case NM_SETTING_VK_IP4_CONFIG_ROUTES_NEXTHOP:
			t = ktypeString
		case NM_SETTING_VK_IP4_CONFIG_ROUTES_METRIC:
			t = ktypeString
		}
	case fieldIPv6:
		switch key {
		case NM_SETTING_VK_IP6_CONFIG_DNS:
			t = ktypeString
		case NM_SETTING_VK_IP6_CONFIG_ADDRESSES_ADDRESS:
			t = ktypeString
		case NM_SETTING_VK_IP6_CONFIG_ADDRESSES_PREFIX:
			t = ktypeString
		case NM_SETTING_VK_IP6_CONFIG_ADDRESSES_GATEWAY:
			t = ktypeString
		case NM_SETTING_VK_IP6_CONFIG_ROUTES_ADDRESS:
			t = ktypeString
		case NM_SETTING_VK_IP6_CONFIG_ROUTES_PREFIX:
			t = ktypeString
		case NM_SETTING_VK_IP6_CONFIG_ROUTES_NEXTHOP:
			t = ktypeString
		case NM_SETTING_VK_IP6_CONFIG_ROUTES_METRIC:
			t = ktypeString
		}
	case fieldWired:
	case fieldWireless:
	case fieldWirelessSecurity:
		switch key {
		case NM_SETTING_VK_WIRELESS_SECURITY_KEY_MGMT:
			t = ktypeString
		}
	}
	return
}

func getRelatedVirtualKeys(field, key string) (vks []string) {
	switch field {
	case field8021x:
	case fieldConnection:
		switch key {
		case NM_SETTING_CONNECTION_PERMISSIONS:
			vks = []string{NM_SETTING_VK_CONNECTION_PERMISSIONS}
		}
	case fieldIPv4:
		switch key {
		case NM_SETTING_IP4_CONFIG_DNS:
			vks = []string{NM_SETTING_VK_IP4_CONFIG_DNS}
		case NM_SETTING_IP4_CONFIG_ADDRESSES:
			vks = []string{
				NM_SETTING_VK_IP4_CONFIG_ADDRESSES_ADDRESS,
				NM_SETTING_VK_IP4_CONFIG_ADDRESSES_MASK,
				NM_SETTING_VK_IP4_CONFIG_ADDRESSES_GATEWAY,
			}
		case NM_SETTING_IP4_CONFIG_ROUTES:
			vks = []string{
				NM_SETTING_VK_IP4_CONFIG_ROUTES_ADDRESS,
				NM_SETTING_VK_IP4_CONFIG_ROUTES_MASK,
				NM_SETTING_VK_IP4_CONFIG_ROUTES_NEXTHOP,
				NM_SETTING_VK_IP4_CONFIG_ROUTES_METRIC,
			}
		}
	case fieldIPv6:
		switch key {
		case NM_SETTING_IP6_CONFIG_DNS:
			vks = []string{NM_SETTING_VK_IP6_CONFIG_DNS}
		case NM_SETTING_IP6_CONFIG_ADDRESSES:
			vks = []string{
				NM_SETTING_VK_IP6_CONFIG_ADDRESSES_ADDRESS,
				NM_SETTING_VK_IP6_CONFIG_ADDRESSES_PREFIX,
				NM_SETTING_VK_IP6_CONFIG_ADDRESSES_GATEWAY,
			}
		case NM_SETTING_IP6_CONFIG_ROUTES:
			vks = []string{
				NM_SETTING_VK_IP6_CONFIG_ROUTES_ADDRESS,
				NM_SETTING_VK_IP6_CONFIG_ROUTES_PREFIX,
				NM_SETTING_VK_IP6_CONFIG_ROUTES_NEXTHOP,
				NM_SETTING_VK_IP6_CONFIG_ROUTES_METRIC,
			}
		}
	case fieldWired:
	case fieldWireless:
	case fieldWirelessSecurity:
		switch key {
		case NM_SETTING_WIRELESS_SECURITY_KEY_MGMT:
			vks = []string{NM_SETTING_VK_WIRELESS_SECURITY_KEY_MGMT}
		}
	}
	return
}

func generalGetVirtualKeyJSON(data _ConnectionData, field, key string) (valueJSON string) {
	switch field {
	case field8021x:
	case fieldConnection:
		switch key {
		case NM_SETTING_VK_CONNECTION_PERMISSIONS:
			valueJSON = getSettingVkConnectionPermissionsJSON(data)
		}
	case fieldIPv4:
		switch key {
		case NM_SETTING_VK_IP4_CONFIG_DNS:
			valueJSON = getSettingVkIp4ConfigDnsJSON(data)
		case NM_SETTING_VK_IP4_CONFIG_ADDRESSES_ADDRESS:
			valueJSON = getSettingVkIp4ConfigAddressesAddressJSON(data)
		case NM_SETTING_VK_IP4_CONFIG_ADDRESSES_MASK:
			valueJSON = getSettingVkIp4ConfigAddressesMaskJSON(data)
		case NM_SETTING_VK_IP4_CONFIG_ADDRESSES_GATEWAY:
			valueJSON = getSettingVkIp4ConfigAddressesGatewayJSON(data)
		case NM_SETTING_VK_IP4_CONFIG_ROUTES_ADDRESS:
			valueJSON = getSettingVkIp4ConfigRoutesAddressJSON(data)
		case NM_SETTING_VK_IP4_CONFIG_ROUTES_MASK:
			valueJSON = getSettingVkIp4ConfigRoutesMaskJSON(data)
		case NM_SETTING_VK_IP4_CONFIG_ROUTES_NEXTHOP:
			valueJSON = getSettingVkIp4ConfigRoutesNexthopJSON(data)
		case NM_SETTING_VK_IP4_CONFIG_ROUTES_METRIC:
			valueJSON = getSettingVkIp4ConfigRoutesMetricJSON(data)
		}
	case fieldIPv6:
		switch key {
		case NM_SETTING_VK_IP6_CONFIG_DNS:
			valueJSON = getSettingVkIp6ConfigDnsJSON(data)
		case NM_SETTING_VK_IP6_CONFIG_ADDRESSES_ADDRESS:
			valueJSON = getSettingVkIp6ConfigAddressesAddressJSON(data)
		case NM_SETTING_VK_IP6_CONFIG_ADDRESSES_PREFIX:
			valueJSON = getSettingVkIp6ConfigAddressesPrefixJSON(data)
		case NM_SETTING_VK_IP6_CONFIG_ADDRESSES_GATEWAY:
			valueJSON = getSettingVkIp6ConfigAddressesGatewayJSON(data)
		case NM_SETTING_VK_IP6_CONFIG_ROUTES_ADDRESS:
			valueJSON = getSettingVkIp6ConfigRoutesAddressJSON(data)
		case NM_SETTING_VK_IP6_CONFIG_ROUTES_PREFIX:
			valueJSON = getSettingVkIp6ConfigRoutesPrefixJSON(data)
		case NM_SETTING_VK_IP6_CONFIG_ROUTES_NEXTHOP:
			valueJSON = getSettingVkIp6ConfigRoutesNexthopJSON(data)
		case NM_SETTING_VK_IP6_CONFIG_ROUTES_METRIC:
			valueJSON = getSettingVkIp6ConfigRoutesMetricJSON(data)
		}
	case fieldWired:
	case fieldWireless:
	case fieldWirelessSecurity:
		switch key {
		case NM_SETTING_VK_WIRELESS_SECURITY_KEY_MGMT:
			valueJSON = getSettingVkWirelessSecurityKeyMgmtJSON(data)
		}
	}
	return
}

func generalSetVirtualKeyJSON(data _ConnectionData, field, key string, valueJSON string) {
	switch field {
	case field8021x:
	case fieldConnection:
		switch key {
		case NM_SETTING_VK_CONNECTION_PERMISSIONS:
			setSettingVkConnectionPermissionsJSON(data, valueJSON)
		}
	case fieldIPv4:
		switch key {
		case NM_SETTING_VK_IP4_CONFIG_DNS:
			setSettingVkIp4ConfigDnsJSON(data, valueJSON)
		case NM_SETTING_VK_IP4_CONFIG_ADDRESSES_ADDRESS:
			setSettingVkIp4ConfigAddressesAddressJSON(data, valueJSON)
		case NM_SETTING_VK_IP4_CONFIG_ADDRESSES_MASK:
			setSettingVkIp4ConfigAddressesMaskJSON(data, valueJSON)
		case NM_SETTING_VK_IP4_CONFIG_ADDRESSES_GATEWAY:
			setSettingVkIp4ConfigAddressesGatewayJSON(data, valueJSON)
		case NM_SETTING_VK_IP4_CONFIG_ROUTES_ADDRESS:
			setSettingVkIp4ConfigRoutesAddressJSON(data, valueJSON)
		case NM_SETTING_VK_IP4_CONFIG_ROUTES_MASK:
			setSettingVkIp4ConfigRoutesMaskJSON(data, valueJSON)
		case NM_SETTING_VK_IP4_CONFIG_ROUTES_NEXTHOP:
			setSettingVkIp4ConfigRoutesNexthopJSON(data, valueJSON)
		case NM_SETTING_VK_IP4_CONFIG_ROUTES_METRIC:
			setSettingVkIp4ConfigRoutesMetricJSON(data, valueJSON)
		}
	case fieldIPv6:
		switch key {
		case NM_SETTING_VK_IP6_CONFIG_DNS:
			setSettingVkIp6ConfigDnsJSON(data, valueJSON)
		case NM_SETTING_VK_IP6_CONFIG_ADDRESSES_ADDRESS:
			setSettingVkIp6ConfigAddressesAddressJSON(data, valueJSON)
		case NM_SETTING_VK_IP6_CONFIG_ADDRESSES_PREFIX:
			setSettingVkIp6ConfigAddressesPrefixJSON(data, valueJSON)
		case NM_SETTING_VK_IP6_CONFIG_ADDRESSES_GATEWAY:
			setSettingVkIp6ConfigAddressesGatewayJSON(data, valueJSON)
		case NM_SETTING_VK_IP6_CONFIG_ROUTES_ADDRESS:
			setSettingVkIp6ConfigRoutesAddressJSON(data, valueJSON)
		case NM_SETTING_VK_IP6_CONFIG_ROUTES_PREFIX:
			setSettingVkIp6ConfigRoutesPrefixJSON(data, valueJSON)
		case NM_SETTING_VK_IP6_CONFIG_ROUTES_NEXTHOP:
			setSettingVkIp6ConfigRoutesNexthopJSON(data, valueJSON)
		case NM_SETTING_VK_IP6_CONFIG_ROUTES_METRIC:
			setSettingVkIp6ConfigRoutesMetricJSON(data, valueJSON)
		}
	case fieldWired:
	case fieldWireless:
	case fieldWirelessSecurity:
		switch key {
		case NM_SETTING_VK_WIRELESS_SECURITY_KEY_MGMT:
			setSettingVkWirelessSecurityKeyMgmtJSON(data, valueJSON)
		}
	}
	return
}

// JSON getter for virtual keys
func getSettingVkConnectionPermissionsJSON(data _ConnectionData) (valueJSON string) {
	// TODO
	// value = getSettingConnectionPermissions(data)
	return
}
func getSettingVkIp4ConfigDnsJSON(data _ConnectionData) (valueJSON string) {
	valueJSON = jsonEmptyString
	value := getSettingIp4ConfigDns(data)
	if len(value) == 0 {
		return
	}
	valueJSON, _ = marshalJSON(convertIpv4AddressToString(value[0]))
	return
}
func getSettingVkIp4ConfigAddressesAddressJSON(data _ConnectionData) (valueJSON string) {
	// TODO
	valueJSON = jsonEmptyString
	value := getSettingIp4ConfigAddresses(data)
	if len(value) == 0 {
		return
	}
	addr := value[0]
	if len(addr) != 3 {
		return
	}
	valueJSON, _ = marshalJSON(convertIpv4AddressToString(addr[0]))
	return
}
func getSettingVkIp4ConfigAddressesMaskJSON(data _ConnectionData) (valueJSON string) {
	// TODO
	valueJSON = jsonEmptyString
	value := getSettingIp4ConfigAddresses(data)
	if len(value) == 0 {
		return
	}
	addr := value[0]
	if len(addr) != 3 {
		return
	}
	valueJSON, _ = marshalJSON(convertIpv4PrefixToNetMask(addr[1]))
	return
}
func getSettingVkIp4ConfigAddressesGatewayJSON(data _ConnectionData) (valueJSON string) {
	// TODO
	valueJSON = jsonEmptyString
	value := getSettingIp4ConfigAddresses(data)
	if len(value) == 0 {
		return
	}
	addr := value[0]
	if len(addr) != 3 {
		return
	}
	valueJSON, _ = marshalJSON(convertIpv4AddressToString(addr[2]))
	return
}
func getSettingVkIp4ConfigRoutesAddressJSON(data _ConnectionData) (valueJSON string) {
	// TODO
	// valueJSON = jsonEmptyString
	// value := getSettingIp4ConfigRoutesAddress(data)
	return
}
func getSettingVkIp4ConfigRoutesMaskJSON(data _ConnectionData) (valueJSON string) {
	// TODO
	// valueJSON = jsonEmptyString
	// value := getSettingIp4ConfigRoutesMask(data)
	return
}
func getSettingVkIp4ConfigRoutesNexthopJSON(data _ConnectionData) (valueJSON string) {
	// TODO
	// valueJSON = jsonEmptyString
	// value := getSettingIp4ConfigRoutesNexthop(data)
	return
}
func getSettingVkIp4ConfigRoutesMetricJSON(data _ConnectionData) (valueJSON string) {
	// TODO
	// valueJSON = jsonEmptyString
	// value := getSettingIp4ConfigRoutesMetric(data)
	return
}
func getSettingVkIp6ConfigDnsJSON(data _ConnectionData) (valueJSON string) {
	// TODO
	// valueJSON = jsonEmptyString
	// value := getSettingIp6ConfigDns(data)
	return
}
func getSettingVkIp6ConfigAddressesAddressJSON(data _ConnectionData) (valueJSON string) {
	// TODO
	// valueJSON = jsonEmptyString
	// value := getSettingIp6ConfigAddressesAddress(data)
	return
}
func getSettingVkIp6ConfigAddressesPrefixJSON(data _ConnectionData) (valueJSON string) {
	// TODO
	// valueJSON = jsonEmptyString
	// value := getSettingIp6ConfigAddressesPrefix(data)
	return
}
func getSettingVkIp6ConfigAddressesGatewayJSON(data _ConnectionData) (valueJSON string) {
	// TODO
	// valueJSON = jsonEmptyString
	// value := getSettingIp6ConfigAddressesGateway(data)
	return
}
func getSettingVkIp6ConfigRoutesAddressJSON(data _ConnectionData) (valueJSON string) {
	// TODO
	// valueJSON = jsonEmptyString
	// value := getSettingIp6ConfigRoutesAddress(data)
	return
}
func getSettingVkIp6ConfigRoutesPrefixJSON(data _ConnectionData) (valueJSON string) {
	// TODO
	// valueJSON = jsonEmptyString
	// value := getSettingIp6ConfigRoutesPrefix(data)
	return
}
func getSettingVkIp6ConfigRoutesNexthopJSON(data _ConnectionData) (valueJSON string) {
	// TODO
	// valueJSON = jsonEmptyString
	// value := getSettingIp6ConfigRoutesNexthop(data)
	return
}
func getSettingVkIp6ConfigRoutesMetricJSON(data _ConnectionData) (valueJSON string) {
	// TODO
	// valueJSON = jsonEmptyString
	// value := getSettingIp6ConfigRoutesMetric(data)
	return
}
func getSettingVkWirelessSecurityKeyMgmtJSON(data _ConnectionData) (valueJSON string) {
	// TODO
	// valueJSON = jsonEmptyString
	// value := getSettingWirelessSecurityKeyMgmtJSON(data)
	return
}

// JSON setter for virtual keys
func setSettingVkConnectionPermissionsJSON(data _ConnectionData, valueJSON string) {
	// TODO
	// setSettingConnectionPermissionsJSON(data)
}
func setSettingVkIp4ConfigDnsJSON(data _ConnectionData, valueJSON string) {
	// TODO
	strAddr, _ := jsonToKeyValueString(valueJSON)
	value := []uint32{convertIpv4AddressToUint32(strAddr)}
	setSettingIp4ConfigDns(data, value)
}
func setSettingVkIp4ConfigAddressesAddressJSON(data _ConnectionData, valueJSON string) {
	// TODO
	strAddr, _ := jsonToKeyValueString(valueJSON)
	value := getSettingIp4ConfigAddresses(data)
	if len(value) == 0 {
		value = make([][]uint32, 1)
		value[0] = make([]uint32, 3)
	}
	addr := value[0]
	if len(addr) != 3 {
		value[0] = make([]uint32, 3)
	}
	addr[0] = convertIpv4AddressToUint32(strAddr)
	setSettingIp4ConfigAddresses(data, value)
}
func setSettingVkIp4ConfigAddressesMaskJSON(data _ConnectionData, valueJSON string) {
	// TODO
	strAddr, _ := jsonToKeyValueString(valueJSON)
	value := getSettingIp4ConfigAddresses(data)
	if len(value) == 0 {
		value = make([][]uint32, 1)
		value[0] = make([]uint32, 3)
	}
	addr := value[0]
	if len(addr) != 3 {
		value[0] = make([]uint32, 3)
	}
	addr[1] = convertIpv4NetMaskToPrefix(strAddr)
	setSettingIp4ConfigAddresses(data, value)
}
func setSettingVkIp4ConfigAddressesGatewayJSON(data _ConnectionData, valueJSON string) {
	// TODO
	strAddr, _ := jsonToKeyValueString(valueJSON)
	value := getSettingIp4ConfigAddresses(data)
	if len(value) == 0 {
		value = make([][]uint32, 1)
		value[0] = make([]uint32, 3)
	}
	addr := value[0]
	if len(addr) != 3 {
		value[0] = make([]uint32, 3)
	}
	addr[2] = convertIpv4AddressToUint32(strAddr)
	setSettingIp4ConfigAddresses(data, value)
}
func setSettingVkIp4ConfigRoutesAddressJSON(data _ConnectionData, valueJSON string) {
	// TODO
	// setSettingIp4ConfigRoutesAddressJSON(data)
}
func setSettingVkIp4ConfigRoutesMaskJSON(data _ConnectionData, valueJSON string) {
	// TODO
	// setSettingIp4ConfigRoutesMaskJSON(data)
}
func setSettingVkIp4ConfigRoutesNexthopJSON(data _ConnectionData, valueJSON string) {
	// TODO
	// setSettingIp4ConfigRoutesNexthopJSON(data)
}
func setSettingVkIp4ConfigRoutesMetricJSON(data _ConnectionData, valueJSON string) {
	// TODO
	// setSettingIp4ConfigRoutesMetricJSON(data)
}
func setSettingVkIp6ConfigDnsJSON(data _ConnectionData, valueJSON string) {
	// TODO
	// setSettingIp6ConfigDnsJSON(data)
}
func setSettingVkIp6ConfigAddressesAddressJSON(data _ConnectionData, valueJSON string) {
	// TODO
	// setSettingIp6ConfigAddressesAddressJSON(data)
}
func setSettingVkIp6ConfigAddressesPrefixJSON(data _ConnectionData, valueJSON string) {
	// TODO
	// setSettingIp6ConfigAddressesPrefixJSON(data)
}
func setSettingVkIp6ConfigAddressesGatewayJSON(data _ConnectionData, valueJSON string) {
	// TODO
	// setSettingIp6ConfigAddressesGatewayJSON(data)
}
func setSettingVkIp6ConfigRoutesAddressJSON(data _ConnectionData, valueJSON string) {
	// TODO
	// setSettingIp6ConfigRoutesAddressJSON(data)
}
func setSettingVkIp6ConfigRoutesPrefixJSON(data _ConnectionData, valueJSON string) {
	// TODO
	// setSettingIp6ConfigRoutesPrefixJSON(data)
}
func setSettingVkIp6ConfigRoutesNexthopJSON(data _ConnectionData, valueJSON string) {
	// TODO
	// setSettingIp6ConfigRoutesNexthopJSON(data)
}
func setSettingVkIp6ConfigRoutesMetricJSON(data _ConnectionData, valueJSON string) {
	// TODO
	// setSettingIp6ConfigRoutesMetricJSON(data)
}
func setSettingVkWirelessSecurityKeyMgmtJSON(data _ConnectionData, valueJSON string) {
	// TODO
	// setSettingWirelessSecurityKeyMgmtJSON(data)
}
