package main

// Get key type
func getSettingWiredKeyType(key string) (t ktype) {
	switch key {
	default:
		t = ktypeUnknown
	case NM_SETTING_WIRED_PORT:
		t = ktypeString
	case NM_SETTING_WIRED_SPEED:
		t = ktypeUint32
	case NM_SETTING_WIRED_DUPLEX:
		t = ktypeString
	case NM_SETTING_WIRED_AUTO_NEGOTIATE:
		t = ktypeBoolean
	case NM_SETTING_WIRED_MAC_ADDRESS:
		t = ktypeArrayByte
	case NM_SETTING_WIRED_CLONED_MAC_ADDRESS:
		t = ktypeArrayByte
	case NM_SETTING_WIRED_MAC_ADDRESS_BLACKLIST:
		t = ktypeArrayString
	case NM_SETTING_WIRED_MTU:
		t = ktypeUint32
	case NM_SETTING_WIRED_S390_SUBCHANNELS:
		t = ktypeArrayString
	case NM_SETTING_WIRED_S390_NETTYPE:
		t = ktypeString
	case NM_SETTING_WIRED_S390_OPTIONS:
		t = ktypeDictStringString
	}
	return
}

// TODO tmp
func setSettingWiredDuplex(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_DUPLEX, value)
}

// Getter
func getSettingWiredPortJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_PORT, getSettingWiredKeyType(NM_SETTING_WIRED_PORT))
	return
}
func getSettingWiredSpeedJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_SPEED, getSettingWiredKeyType(NM_SETTING_WIRED_SPEED))
	return
}
func getSettingWiredDuplexJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_DUPLEX, getSettingWiredKeyType(NM_SETTING_WIRED_DUPLEX))
	return
}
func getSettingWiredAutoNegotiateJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_AUTO_NEGOTIATE, getSettingWiredKeyType(NM_SETTING_WIRED_AUTO_NEGOTIATE))
	return
}
func getSettingWiredMacAddressJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_MAC_ADDRESS, getSettingWiredKeyType(NM_SETTING_WIRED_MAC_ADDRESS))
	return
}
func getSettingWiredClonedMacAddressJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_CLONED_MAC_ADDRESS, getSettingWiredKeyType(NM_SETTING_WIRED_CLONED_MAC_ADDRESS))
	return
}
func getSettingWiredMacAddressBlacklistJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_MAC_ADDRESS_BLACKLIST, getSettingWiredKeyType(NM_SETTING_WIRED_MAC_ADDRESS_BLACKLIST))
	return
}
func getSettingWiredMtuJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_MTU, getSettingWiredKeyType(NM_SETTING_WIRED_MTU))
	return
}
func getSettingWiredS390SubchannelsJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_S390_SUBCHANNELS, getSettingWiredKeyType(NM_SETTING_WIRED_S390_SUBCHANNELS))
	return
}
func getSettingWiredS390NettypeJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_S390_NETTYPE, getSettingWiredKeyType(NM_SETTING_WIRED_S390_NETTYPE))
	return
}
func getSettingWiredS390OptionsJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_S390_OPTIONS, getSettingWiredKeyType(NM_SETTING_WIRED_S390_OPTIONS))
	return
}

// Setter
func setSettingWiredPortJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_PORT, value, getSettingWiredKeyType(NM_SETTING_WIRED_PORT))
}
func setSettingWiredSpeedJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_SPEED, value, getSettingWiredKeyType(NM_SETTING_WIRED_SPEED))
}
func setSettingWiredDuplexJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_DUPLEX, value, getSettingWiredKeyType(NM_SETTING_WIRED_DUPLEX))
}
func setSettingWiredAutoNegotiateJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_AUTO_NEGOTIATE, value, getSettingWiredKeyType(NM_SETTING_WIRED_AUTO_NEGOTIATE))
}
func setSettingWiredMacAddressJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_MAC_ADDRESS, value, getSettingWiredKeyType(NM_SETTING_WIRED_MAC_ADDRESS))
}
func setSettingWiredClonedMacAddressJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_CLONED_MAC_ADDRESS, value, getSettingWiredKeyType(NM_SETTING_WIRED_CLONED_MAC_ADDRESS))
}
func setSettingWiredMacAddressBlacklistJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_MAC_ADDRESS_BLACKLIST, value, getSettingWiredKeyType(NM_SETTING_WIRED_MAC_ADDRESS_BLACKLIST))
}
func setSettingWiredMtuJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_MTU, value, getSettingWiredKeyType(NM_SETTING_WIRED_MTU))
}
func setSettingWiredS390SubchannelsJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_S390_SUBCHANNELS, value, getSettingWiredKeyType(NM_SETTING_WIRED_S390_SUBCHANNELS))
}
func setSettingWiredS390NettypeJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_S390_NETTYPE, value, getSettingWiredKeyType(NM_SETTING_WIRED_S390_NETTYPE))
}
func setSettingWiredS390OptionsJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_S390_OPTIONS, value, getSettingWiredKeyType(NM_SETTING_WIRED_S390_OPTIONS))
}

// Remover
func removeSettingWiredPort(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_PORT)
}
func removeSettingWiredSpeed(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_SPEED)
}
func removeSettingWiredDuplex(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_DUPLEX)
}
func removeSettingWiredAutoNegotiate(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_AUTO_NEGOTIATE)
}
func removeSettingWiredMacAddress(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_MAC_ADDRESS)
}
func removeSettingWiredClonedMacAddress(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_CLONED_MAC_ADDRESS)
}
func removeSettingWiredMacAddressBlacklist(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_MAC_ADDRESS_BLACKLIST)
}
func removeSettingWiredMtu(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_MTU)
}
func removeSettingWiredS390Subchannels(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_S390_SUBCHANNELS)
}
func removeSettingWiredS390Nettype(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_S390_NETTYPE)
}
func removeSettingWiredS390Options(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRED_SETTING_NAME, NM_SETTING_WIRED_S390_OPTIONS)
}
