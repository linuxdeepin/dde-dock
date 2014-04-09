package main

// Virtual key names.

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

// VirtualKey store virtual key info for each fields.
type VirtualKey struct {
	name       string
	field      string
	keyType    ktype
	relatedKey string
	available  bool
	required   bool // check if child virtual key is optional
}

var virtualKeys = []VirtualKey{
	VirtualKey{NM_SETTING_VK_CONNECTION_PERMISSIONS, fieldConnection, ktypeBoolean, NM_SETTING_CONNECTION_PERMISSIONS, true, true},
	VirtualKey{NM_SETTING_VK_IP4_CONFIG_DNS, fieldIPv4, ktypeString, NM_SETTING_IP4_CONFIG_DNS, true, true},
	VirtualKey{NM_SETTING_VK_IP4_CONFIG_ADDRESSES_ADDRESS, fieldIPv4, ktypeString, NM_SETTING_IP4_CONFIG_ADDRESSES, true, true},
	VirtualKey{NM_SETTING_VK_IP4_CONFIG_ADDRESSES_MASK, fieldIPv4, ktypeString, NM_SETTING_IP4_CONFIG_ADDRESSES, true, true},
	VirtualKey{NM_SETTING_VK_IP4_CONFIG_ADDRESSES_GATEWAY, fieldIPv4, ktypeString, NM_SETTING_IP4_CONFIG_ADDRESSES, true, false},
	VirtualKey{NM_SETTING_VK_IP4_CONFIG_ROUTES_ADDRESS, fieldIPv4, ktypeString, NM_SETTING_IP4_CONFIG_ROUTES, true, true},
	VirtualKey{NM_SETTING_VK_IP4_CONFIG_ROUTES_MASK, fieldIPv4, ktypeString, NM_SETTING_IP4_CONFIG_ROUTES, true, true},
	VirtualKey{NM_SETTING_VK_IP4_CONFIG_ROUTES_NEXTHOP, fieldIPv4, ktypeString, NM_SETTING_IP4_CONFIG_ROUTES, true, true},
	VirtualKey{NM_SETTING_VK_IP4_CONFIG_ROUTES_METRIC, fieldIPv4, ktypeString, NM_SETTING_IP4_CONFIG_ROUTES, false, false},
	VirtualKey{NM_SETTING_VK_IP6_CONFIG_DNS, fieldIPv6, ktypeString, NM_SETTING_IP6_CONFIG_DNS, true, true},
	VirtualKey{NM_SETTING_VK_IP6_CONFIG_ADDRESSES_ADDRESS, fieldIPv6, ktypeString, NM_SETTING_IP6_CONFIG_ADDRESSES, true, true},
	VirtualKey{NM_SETTING_VK_IP6_CONFIG_ADDRESSES_PREFIX, fieldIPv6, ktypeString, NM_SETTING_IP6_CONFIG_ADDRESSES, true, true},
	VirtualKey{NM_SETTING_VK_IP6_CONFIG_ADDRESSES_GATEWAY, fieldIPv6, ktypeString, NM_SETTING_IP6_CONFIG_ADDRESSES, true, false},
	VirtualKey{NM_SETTING_VK_IP6_CONFIG_ROUTES_ADDRESS, fieldIPv6, ktypeString, NM_SETTING_IP6_CONFIG_ROUTES, true, true},
	VirtualKey{NM_SETTING_VK_IP6_CONFIG_ROUTES_PREFIX, fieldIPv6, ktypeString, NM_SETTING_IP6_CONFIG_ROUTES, true, true},
	VirtualKey{NM_SETTING_VK_IP6_CONFIG_ROUTES_NEXTHOP, fieldIPv6, ktypeString, NM_SETTING_IP6_CONFIG_ROUTES, true, true},
	VirtualKey{NM_SETTING_VK_IP6_CONFIG_ROUTES_METRIC, fieldIPv6, ktypeString, NM_SETTING_IP6_CONFIG_ROUTES, false, false},
	VirtualKey{NM_SETTING_VK_WIRELESS_SECURITY_KEY_MGMT, fieldWirelessSecurity, ktypeString, NM_SETTING_WIRELESS_SECURITY_KEY_MGMT, true, true},
}

func isVirtualKey(field, key string) bool {
	if isStringInArray(key, getVirtualKeysOfField(field)) {
		return true
	}
	return false
}

func getVirtualKeysOfField(field string) (vks []string) {
	for _, vk := range virtualKeys {
		if vk.field == field {
			vks = append(vks, vk.name)
		}
	}
	return
}

func getSettingVkKeyType(field, key string) (t ktype) {
	t = ktypeUnknown
	for _, vk := range virtualKeys {
		if vk.field == field && vk.name == key {
			t = vk.keyType
		}
	}
	return
}

func generalGetSettingVkAvailableValues(field, key string) (values []string) {
	switch field {
	case field8021x:
	case fieldConnection:
	case fieldIPv4:
	case fieldIPv6:
	case fieldWired:
	case fieldWireless:
	case fieldWirelessSecurity:
		switch key {
		case NM_SETTING_VK_WIRELESS_SECURITY_KEY_MGMT:
			// TODO
			values = []string{"none", "wep-128", "leap", "wpa", "wpa-eap"}
		}
	}
	return
}

func getRelatedAvailableVirtualKeys(field, key string) (vks []string) {
	for _, vk := range virtualKeys {
		if vk.field == field && vk.relatedKey == key && vk.available {
			vks = append(vks, vk.name)
		}
	}
	return
}

// get related virtual key(s) for target key
func getRelatedVirtualKeys(field, key string) (vks []string) {
	for _, vk := range virtualKeys {
		if vk.field == field && vk.relatedKey == key {
			vks = append(vks, vk.name)
		}
	}
	return
}

// check if is required child virtual key
func isRequiredChildVirtualKeys(field, vkey string) (required bool) {
	for _, vk := range virtualKeys {
		if vk.field == field && vk.name == vkey {
			required = vk.required
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
			logicSetSettingVkWirelessSecurityKeyMgmtJSON(data, valueJSON)
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
	value = convertIpv4AddressToStringNoZeor(addresses[0][2])
	return
}
func getSettingVkIp4ConfigRoutesAddress(data _ConnectionData) (value string) {
	// TODO
	// valueJSON = jsonEmptyString
	// value := getSettingIp4ConfigRoutesAddress(data)
	return
}
func getSettingVkIp4ConfigRoutesMask(data _ConnectionData) (value string) {
	// TODO
	// valueJSON = jsonEmptyString
	// value := getSettingIp4ConfigRoutesMask(data)
	return
}
func getSettingVkIp4ConfigRoutesNexthop(data _ConnectionData) (value string) {
	// TODO
	// valueJSON = jsonEmptyString
	// value := getSettingIp4ConfigRoutesNexthop(data)
	return
}
func getSettingVkIp4ConfigRoutesMetric(data _ConnectionData) (value string) {
	// TODO
	// valueJSON = jsonEmptyString
	// value := getSettingIp4ConfigRoutesMetric(data)
	return
}
func getSettingVkIp6ConfigDns(data _ConnectionData) (value string) {
	// TODO
	// valueJSON = jsonEmptyString
	// value := getSettingIp6ConfigDns(data)
	return
}
func getSettingVkIp6ConfigAddressesAddress(data _ConnectionData) (value string) {
	// TODO
	// valueJSON = jsonEmptyString
	// value := getSettingIp6ConfigAddressesAddress(data)
	return
}
func getSettingVkIp6ConfigAddressesPrefix(data _ConnectionData) (value string) {
	// TODO
	// valueJSON = jsonEmptyString
	// value := getSettingIp6ConfigAddressesPrefix(data)
	return
}
func getSettingVkIp6ConfigAddressesGateway(data _ConnectionData) (value string) {
	// TODO
	// valueJSON = jsonEmptyString
	// value := getSettingIp6ConfigAddressesGateway(data)
	return
}
func getSettingVkIp6ConfigRoutesAddress(data _ConnectionData) (value string) {
	// TODO
	// valueJSON = jsonEmptyString
	// value := getSettingIp6ConfigRoutesAddress(data)
	return
}
func getSettingVkIp6ConfigRoutesPrefix(data _ConnectionData) (value string) {
	// TODO
	// valueJSON = jsonEmptyString
	// value := getSettingIp6ConfigRoutesPrefix(data)
	return
}
func getSettingVkIp6ConfigRoutesNexthop(data _ConnectionData) (value string) {
	// TODO
	// valueJSON = jsonEmptyString
	// value := getSettingIp6ConfigRoutesNexthop(data)
	return
}
func getSettingVkIp6ConfigRoutesMetric(data _ConnectionData) (value string) {
	// TODO
	// valueJSON = jsonEmptyString
	// value := getSettingIp6ConfigRoutesMetric(data)
	return
}
func getSettingVkWirelessSecurityKeyMgmt(data _ConnectionData) (value string) {
	// TODO
	// valueJSON = jsonEmptyString
	// value := getSettingWirelessSecurityKeyMgmtJSON(data)
	return
}

// Setter
func setSettingVkConnectionPermissions(data _ConnectionData, value bool) {
	// TODO
	// setSettingConnectionPermissionsJSON(data)
}
func setSettingVkIp4ConfigDns(data _ConnectionData, value string) {
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
func setSettingVkIp4ConfigAddressesAddress(data _ConnectionData, value string) {
	addresses := doGetOrNewSettingIp4ConfigAddresses(data)
	addr := addresses[0]
	addr[0] = convertIpv4AddressToUint32(value)
	if !isUint32ArrayEmpty(addr) {
		setSettingIp4ConfigAddresses(data, addresses)
	} else {
		removeSettingIp4ConfigAddresses(data)
	}
}
func setSettingVkIp4ConfigAddressesMask(data _ConnectionData, value string) {
	addresses := doGetOrNewSettingIp4ConfigAddresses(data)
	addr := addresses[0]
	addr[1] = convertIpv4NetMaskToPrefix(value)
	if !isUint32ArrayEmpty(addr) {
		setSettingIp4ConfigAddresses(data, addresses)
	} else {
		removeSettingIp4ConfigAddresses(data)
	}
}
func setSettingVkIp4ConfigAddressesGateway(data _ConnectionData, value string) {
	addresses := doGetOrNewSettingIp4ConfigAddresses(data)
	addr := addresses[0]
	addr[2] = convertIpv4AddressToUint32(value)
	if !isUint32ArrayEmpty(addr) {
		setSettingIp4ConfigAddresses(data, addresses)
	} else {
		removeSettingIp4ConfigAddresses(data)
	}
}
func setSettingVkIp4ConfigRoutesAddress(data _ConnectionData, value string) {
	// TODO
	// setSettingIp4ConfigRoutesAddressJSON(data)
}
func setSettingVkIp4ConfigRoutesMask(data _ConnectionData, value string) {
	// TODO
	// setSettingIp4ConfigRoutesMaskJSON(data)
}
func setSettingVkIp4ConfigRoutesNexthop(data _ConnectionData, value string) {
	// TODO
	// setSettingIp4ConfigRoutesNexthopJSON(data)
}
func setSettingVkIp4ConfigRoutesMetric(data _ConnectionData, value string) {
	// TODO
	// setSettingIp4ConfigRoutesMetricJSON(data)
}
func setSettingVkIp6ConfigDns(data _ConnectionData, value string) {
	// TODO
	// setSettingIp6ConfigDnsJSON(data)
}
func setSettingVkIp6ConfigAddressesAddress(data _ConnectionData, value string) {
	// TODO
	// setSettingIp6ConfigAddressesAddressJSON(data)
}
func setSettingVkIp6ConfigAddressesPrefix(data _ConnectionData, value string) {
	// TODO
	// setSettingIp6ConfigAddressesPrefixJSON(data)
}
func setSettingVkIp6ConfigAddressesGateway(data _ConnectionData, value string) {
	// TODO
	// setSettingIp6ConfigAddressesGatewayJSON(data)
}
func setSettingVkIp6ConfigRoutesAddress(data _ConnectionData, value string) {
	// TODO
	// setSettingIp6ConfigRoutesAddressJSON(data)
}
func setSettingVkIp6ConfigRoutesPrefix(data _ConnectionData, value string) {
	// TODO
	// setSettingIp6ConfigRoutesPrefixJSON(data)
}
func setSettingVkIp6ConfigRoutesNexthop(data _ConnectionData, value string) {
	// TODO
	// setSettingIp6ConfigRoutesNexthopJSON(data)
}
func setSettingVkIp6ConfigRoutesMetric(data _ConnectionData, value string) {
	// TODO
	// setSettingIp6ConfigRoutesMetricJSON(data)
}
func setSettingVkWirelessSecurityKeyMgmt(data _ConnectionData, value string) {
	// TODO
	// setSettingWirelessSecurityKeyMgmtJSON(data)
}

// JSON getter
func getSettingVkConnectionPermissionsJSON(data _ConnectionData) (valueJSON string) {
	valueJSON, _ = marshalJSON(getSettingVkConnectionPermissions(data))
	return
}
func getSettingVkIp4ConfigDnsJSON(data _ConnectionData) (valueJSON string) {
	valueJSON, _ = marshalJSON(getSettingVkIp4ConfigDns(data))
	return
}
func getSettingVkIp4ConfigAddressesAddressJSON(data _ConnectionData) (valueJSON string) {
	valueJSON, _ = marshalJSON(getSettingVkIp4ConfigAddressesAddress(data))
	return
}
func getSettingVkIp4ConfigAddressesMaskJSON(data _ConnectionData) (valueJSON string) {
	valueJSON, _ = marshalJSON(getSettingVkIp4ConfigAddressesMask(data))
	return
}
func getSettingVkIp4ConfigAddressesGatewayJSON(data _ConnectionData) (valueJSON string) {
	valueJSON, _ = marshalJSON(getSettingVkIp4ConfigAddressesGateway(data))
	return
}
func getSettingVkIp4ConfigRoutesAddressJSON(data _ConnectionData) (valueJSON string) {
	valueJSON, _ = marshalJSON(getSettingVkIp4ConfigRoutesAddress(data))
	return
}
func getSettingVkIp4ConfigRoutesMaskJSON(data _ConnectionData) (valueJSON string) {
	valueJSON, _ = marshalJSON(getSettingVkIp4ConfigRoutesMask(data))
	return
}
func getSettingVkIp4ConfigRoutesNexthopJSON(data _ConnectionData) (valueJSON string) {
	valueJSON, _ = marshalJSON(getSettingVkIp4ConfigRoutesNexthop(data))
	return
}
func getSettingVkIp4ConfigRoutesMetricJSON(data _ConnectionData) (valueJSON string) {
	valueJSON, _ = marshalJSON(getSettingVkIp4ConfigRoutesMetric(data))
	return
}
func getSettingVkIp6ConfigDnsJSON(data _ConnectionData) (valueJSON string) {
	valueJSON, _ = marshalJSON(getSettingVkIp6ConfigDns(data))
	return
}
func getSettingVkIp6ConfigAddressesAddressJSON(data _ConnectionData) (valueJSON string) {
	valueJSON, _ = marshalJSON(getSettingVkIp6ConfigAddressesAddress(data))
	return
}
func getSettingVkIp6ConfigAddressesPrefixJSON(data _ConnectionData) (valueJSON string) {
	valueJSON, _ = marshalJSON(getSettingVkIp6ConfigAddressesPrefix(data))
	return
}
func getSettingVkIp6ConfigAddressesGatewayJSON(data _ConnectionData) (valueJSON string) {
	valueJSON, _ = marshalJSON(getSettingVkIp6ConfigAddressesGateway(data))
	return
}
func getSettingVkIp6ConfigRoutesAddressJSON(data _ConnectionData) (valueJSON string) {
	valueJSON, _ = marshalJSON(getSettingVkIp6ConfigRoutesAddress(data))
	return
}
func getSettingVkIp6ConfigRoutesPrefixJSON(data _ConnectionData) (valueJSON string) {
	valueJSON, _ = marshalJSON(getSettingVkIp6ConfigRoutesPrefix(data))
	return
}
func getSettingVkIp6ConfigRoutesNexthopJSON(data _ConnectionData) (valueJSON string) {
	valueJSON, _ = marshalJSON(getSettingVkIp6ConfigRoutesNexthop(data))
	return
}
func getSettingVkIp6ConfigRoutesMetricJSON(data _ConnectionData) (valueJSON string) {
	valueJSON, _ = marshalJSON(getSettingVkIp6ConfigRoutesMetric(data))
	return
}
func getSettingVkWirelessSecurityKeyMgmtJSON(data _ConnectionData) (valueJSON string) {
	valueJSON, _ = marshalJSON(getSettingVkWirelessSecurityKeyMgmt(data))
	return
}

// JSON setter
func setSettingVkConnectionPermissionsJSON(data _ConnectionData, valueJSON string) {
	value, _ := jsonToKeyValueBoolean(valueJSON)
	setSettingVkConnectionPermissions(data, value)
}
func setSettingVkIp4ConfigDnsJSON(data _ConnectionData, valueJSON string) {
	value, _ := jsonToKeyValueString(valueJSON)
	setSettingVkIp4ConfigDns(data, value)
}
func setSettingVkIp4ConfigAddressesAddressJSON(data _ConnectionData, valueJSON string) {
	value, _ := jsonToKeyValueString(valueJSON)
	setSettingVkIp4ConfigAddressesAddress(data, value)
}
func setSettingVkIp4ConfigAddressesMaskJSON(data _ConnectionData, valueJSON string) {
	value, _ := jsonToKeyValueString(valueJSON)
	setSettingVkIp4ConfigAddressesMask(data, value)
}
func setSettingVkIp4ConfigAddressesGatewayJSON(data _ConnectionData, valueJSON string) {
	value, _ := jsonToKeyValueString(valueJSON)
	setSettingVkIp4ConfigAddressesGateway(data, value)
}
func setSettingVkIp4ConfigRoutesAddressJSON(data _ConnectionData, valueJSON string) {
	value, _ := jsonToKeyValueString(valueJSON)
	setSettingVkIp4ConfigRoutesAddress(data, value)
}
func setSettingVkIp4ConfigRoutesMaskJSON(data _ConnectionData, valueJSON string) {
	value, _ := jsonToKeyValueString(valueJSON)
	setSettingVkIp4ConfigRoutesMask(data, value)
}
func setSettingVkIp4ConfigRoutesNexthopJSON(data _ConnectionData, valueJSON string) {
	value, _ := jsonToKeyValueString(valueJSON)
	setSettingVkIp4ConfigRoutesNexthop(data, value)
}
func setSettingVkIp4ConfigRoutesMetricJSON(data _ConnectionData, valueJSON string) {
	value, _ := jsonToKeyValueString(valueJSON)
	setSettingVkIp4ConfigRoutesMetric(data, value)
}
func setSettingVkIp6ConfigDnsJSON(data _ConnectionData, valueJSON string) {
	value, _ := jsonToKeyValueString(valueJSON)
	setSettingVkIp6ConfigDns(data, value)
}
func setSettingVkIp6ConfigAddressesAddressJSON(data _ConnectionData, valueJSON string) {
	value, _ := jsonToKeyValueString(valueJSON)
	setSettingVkIp6ConfigAddressesAddress(data, value)
}
func setSettingVkIp6ConfigAddressesPrefixJSON(data _ConnectionData, valueJSON string) {
	value, _ := jsonToKeyValueString(valueJSON)
	setSettingVkIp6ConfigAddressesPrefix(data, value)
}
func setSettingVkIp6ConfigAddressesGatewayJSON(data _ConnectionData, valueJSON string) {
	value, _ := jsonToKeyValueString(valueJSON)
	setSettingVkIp6ConfigAddressesGateway(data, value)
}
func setSettingVkIp6ConfigRoutesAddressJSON(data _ConnectionData, valueJSON string) {
	value, _ := jsonToKeyValueString(valueJSON)
	setSettingVkIp6ConfigRoutesAddress(data, value)
}
func setSettingVkIp6ConfigRoutesPrefixJSON(data _ConnectionData, valueJSON string) {
	value, _ := jsonToKeyValueString(valueJSON)
	setSettingVkIp6ConfigRoutesPrefix(data, value)
}
func setSettingVkIp6ConfigRoutesNexthopJSON(data _ConnectionData, valueJSON string) {
	value, _ := jsonToKeyValueString(valueJSON)
	setSettingVkIp6ConfigRoutesNexthop(data, value)
}
func setSettingVkIp6ConfigRoutesMetricJSON(data _ConnectionData, valueJSON string) {
	value, _ := jsonToKeyValueString(valueJSON)
	setSettingVkIp6ConfigRoutesMetric(data, value)
}
func setSettingVkWirelessSecurityKeyMgmtJSON(data _ConnectionData, valueJSON string) {
	value, _ := jsonToKeyValueString(valueJSON)
	setSettingVkWirelessSecurityKeyMgmt(data, value)
}

// Logic setter
func logicSetSettingVkWirelessSecurityKeyMgmtJSON(data _ConnectionData, valueJSON string) {
	// TODO
}
