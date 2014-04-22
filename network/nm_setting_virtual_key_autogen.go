package main

// All virtual keys data
var virtualKeys = []VirtualKey{
	VirtualKey{NM_SETTING_VK_802_1X_EAP, ktypeString, NM_SETTING_802_1X_SETTING_NAME, NM_SETTING_802_1X_EAP, true, true},
	VirtualKey{NM_SETTING_VK_CONNECTION_PERMISSIONS, ktypeBoolean, NM_SETTING_CONNECTION_SETTING_NAME, NM_SETTING_CONNECTION_PERMISSIONS, true, true},
	VirtualKey{NM_SETTING_VK_IP4_CONFIG_DNS, ktypeString, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_DNS, true, true},
	VirtualKey{NM_SETTING_VK_IP4_CONFIG_ADDRESSES_ADDRESS, ktypeString, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_ADDRESSES, true, true},
	VirtualKey{NM_SETTING_VK_IP4_CONFIG_ADDRESSES_MASK, ktypeString, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_ADDRESSES, true, true},
	VirtualKey{NM_SETTING_VK_IP4_CONFIG_ADDRESSES_GATEWAY, ktypeString, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_ADDRESSES, true, false},
	VirtualKey{NM_SETTING_VK_IP4_CONFIG_ROUTES_ADDRESS, ktypeString, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_ROUTES, true, true},
	VirtualKey{NM_SETTING_VK_IP4_CONFIG_ROUTES_MASK, ktypeString, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_ROUTES, true, true},
	VirtualKey{NM_SETTING_VK_IP4_CONFIG_ROUTES_NEXTHOP, ktypeString, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_ROUTES, true, true},
	VirtualKey{NM_SETTING_VK_IP4_CONFIG_ROUTES_METRIC, ktypeString, NM_SETTING_IP4_CONFIG_SETTING_NAME, NM_SETTING_IP4_CONFIG_ROUTES, false, false},
	VirtualKey{NM_SETTING_VK_IP6_CONFIG_DNS, ktypeString, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_DNS, true, true},
	VirtualKey{NM_SETTING_VK_IP6_CONFIG_ADDRESSES_ADDRESS, ktypeString, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_ADDRESSES, true, true},
	VirtualKey{NM_SETTING_VK_IP6_CONFIG_ADDRESSES_PREFIX, ktypeString, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_ADDRESSES, true, true},
	VirtualKey{NM_SETTING_VK_IP6_CONFIG_ADDRESSES_GATEWAY, ktypeString, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_ADDRESSES, true, false},
	VirtualKey{NM_SETTING_VK_IP6_CONFIG_ROUTES_ADDRESS, ktypeString, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_ROUTES, true, true},
	VirtualKey{NM_SETTING_VK_IP6_CONFIG_ROUTES_PREFIX, ktypeString, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_ROUTES, true, true},
	VirtualKey{NM_SETTING_VK_IP6_CONFIG_ROUTES_NEXTHOP, ktypeString, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_ROUTES, true, true},
	VirtualKey{NM_SETTING_VK_IP6_CONFIG_ROUTES_METRIC, ktypeString, NM_SETTING_IP6_CONFIG_SETTING_NAME, NM_SETTING_IP6_CONFIG_ROUTES, false, false},
	VirtualKey{NM_SETTING_VK_WIRELESS_SECURITY_KEY_MGMT, ktypeString, NM_SETTING_WIRELESS_SECURITY_SETTING_NAME, NM_SETTING_WIRELESS_SECURITY_KEY_MGMT, true, true},
}

// Get JSON value generally
func generalGetVirtualKeyJSON(data _ConnectionData, field, key string) (valueJSON string) {
	switch field {
	case NM_SETTING_802_1X_SETTING_NAME:
	case NM_SETTING_CONNECTION_SETTING_NAME:
		switch key {
		case NM_SETTING_VK_CONNECTION_PERMISSIONS:
			valueJSON = getSettingVkConnectionPermissionsJSON(data)
		}
	case NM_SETTING_IP4_CONFIG_SETTING_NAME:
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
	case NM_SETTING_IP6_CONFIG_SETTING_NAME:
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
	case NM_SETTING_WIRED_SETTING_NAME:
	case NM_SETTING_WIRELESS_SETTING_NAME:
	case NM_SETTING_WIRELESS_SECURITY_SETTING_NAME:
		switch key {
		case NM_SETTING_VK_WIRELESS_SECURITY_KEY_MGMT:
			valueJSON = getSettingVkWirelessSecurityKeyMgmtJSON(data)
		}
	case NM_SETTING_PPPOE_SETTING_NAME:
	case NM_SETTING_PPP_SETTING_NAME:
	}
	return
}

// Set JSON value generally
func generalSetVirtualKeyJSON(data _ConnectionData, field, key string, valueJSON string) {
	Logger.Debug("generalSetVirtualKeyJSON:", field, key, valueJSON) // TODO test
	switch field {
	case NM_SETTING_802_1X_SETTING_NAME:
	case NM_SETTING_CONNECTION_SETTING_NAME:
		switch key {
		case NM_SETTING_VK_CONNECTION_PERMISSIONS:
			setSettingVkConnectionPermissionsJSON(data, valueJSON)
		}
	case NM_SETTING_IP4_CONFIG_SETTING_NAME:
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
	case NM_SETTING_IP6_CONFIG_SETTING_NAME:
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
	case NM_SETTING_WIRED_SETTING_NAME:
	case NM_SETTING_WIRELESS_SETTING_NAME:
	case NM_SETTING_WIRELESS_SECURITY_SETTING_NAME:
		switch key {
		case NM_SETTING_VK_WIRELESS_SECURITY_KEY_MGMT:
			logicSetSettingVkWirelessSecurityKeyMgmtJSON(data, valueJSON)
		}
	case NM_SETTING_PPPOE_SETTING_NAME:
	case NM_SETTING_PPP_SETTING_NAME:
	}
	return
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
