package main

// Get key type
func getSettingWirelessKeyType(key string) (t ktype) {
	switch key {
	default:
		t = ktypeUnknown
	case NM_SETTING_WIRELESS_SSID:
		t = ktypeArrayByte
	case NM_SETTING_WIRELESS_MODE:
		t = ktypeString
	case NM_SETTING_WIRELESS_BAND:
		t = ktypeString
	case NM_SETTING_WIRELESS_CHANNEL:
		t = ktypeUint32
	case NM_SETTING_WIRELESS_BSSID:
		t = ktypeArrayByte
	case NM_SETTING_WIRELESS_RATE:
		t = ktypeUint32
	case NM_SETTING_WIRELESS_TX_POWER:
		t = ktypeUint32
	case NM_SETTING_WIRELESS_MAC_ADDRESS:
		t = ktypeArrayByte
	case NM_SETTING_WIRELESS_CLONED_MAC_ADDRESS:
		t = ktypeArrayByte
	case NM_SETTING_WIRELESS_MAC_ADDRESS_BLACKLIST:
		t = ktypeArrayString
	case NM_SETTING_WIRELESS_MTU:
		t = ktypeUint32
	case NM_SETTING_WIRELESS_SEEN_BSSIDS:
		t = ktypeArrayString
	case NM_SETTING_WIRELESS_SEC:
		t = ktypeString
	case NM_SETTING_WIRELESS_HIDDEN:
		t = ktypeBoolean
	}
	return
}

// TODO tmp
func getSettingWirelessSsid(data _ConnectionData) (value []byte) {
	value, _ = getConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_SSID).([]byte)
	return
}
func setSettingWirelessSsid(data _ConnectionData, value []byte) {
	setConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_SSID, value)
}
func setSettingWirelessSec(data _ConnectionData, value string) {
	setConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_SEC, value)
}

// Getter
func getSettingWirelessSsidJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_SSID, getSettingWirelessKeyType(NM_SETTING_WIRELESS_SSID))
	return
}
func getSettingWirelessModeJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_MODE, getSettingWirelessKeyType(NM_SETTING_WIRELESS_MODE))
	return
}
func getSettingWirelessBandJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_BAND, getSettingWirelessKeyType(NM_SETTING_WIRELESS_BAND))
	return
}
func getSettingWirelessChannelJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_CHANNEL, getSettingWirelessKeyType(NM_SETTING_WIRELESS_CHANNEL))
	return
}
func getSettingWirelessBssidJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_BSSID, getSettingWirelessKeyType(NM_SETTING_WIRELESS_BSSID))
	return
}
func getSettingWirelessRateJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_RATE, getSettingWirelessKeyType(NM_SETTING_WIRELESS_RATE))
	return
}
func getSettingWirelessTxPowerJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_TX_POWER, getSettingWirelessKeyType(NM_SETTING_WIRELESS_TX_POWER))
	return
}
func getSettingWirelessMacAddressJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_MAC_ADDRESS, getSettingWirelessKeyType(NM_SETTING_WIRELESS_MAC_ADDRESS))
	return
}
func getSettingWirelessClonedMacAddressJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_CLONED_MAC_ADDRESS, getSettingWirelessKeyType(NM_SETTING_WIRELESS_CLONED_MAC_ADDRESS))
	return
}
func getSettingWirelessMacAddressBlacklistJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_MAC_ADDRESS_BLACKLIST, getSettingWirelessKeyType(NM_SETTING_WIRELESS_MAC_ADDRESS_BLACKLIST))
	return
}
func getSettingWirelessMtuJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_MTU, getSettingWirelessKeyType(NM_SETTING_WIRELESS_MTU))
	return
}
func getSettingWirelessSeenBssidsJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_SEEN_BSSIDS, getSettingWirelessKeyType(NM_SETTING_WIRELESS_SEEN_BSSIDS))
	return
}
func getSettingWirelessSecJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_SEC, getSettingWirelessKeyType(NM_SETTING_WIRELESS_SEC))
	return
}
func getSettingWirelessHiddenJSON(data _ConnectionData) (value string) {
	value = getConnectionDataKeyJSON(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_HIDDEN, getSettingWirelessKeyType(NM_SETTING_WIRELESS_HIDDEN))
	return
}

// Setter
func setSettingWirelessSsidJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_SSID, value, getSettingWirelessKeyType(NM_SETTING_WIRELESS_SSID))
}
func setSettingWirelessModeJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_MODE, value, getSettingWirelessKeyType(NM_SETTING_WIRELESS_MODE))
}
func setSettingWirelessBandJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_BAND, value, getSettingWirelessKeyType(NM_SETTING_WIRELESS_BAND))
}
func setSettingWirelessChannelJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_CHANNEL, value, getSettingWirelessKeyType(NM_SETTING_WIRELESS_CHANNEL))
}
func setSettingWirelessBssidJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_BSSID, value, getSettingWirelessKeyType(NM_SETTING_WIRELESS_BSSID))
}
func setSettingWirelessRateJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_RATE, value, getSettingWirelessKeyType(NM_SETTING_WIRELESS_RATE))
}
func setSettingWirelessTxPowerJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_TX_POWER, value, getSettingWirelessKeyType(NM_SETTING_WIRELESS_TX_POWER))
}
func setSettingWirelessMacAddressJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_MAC_ADDRESS, value, getSettingWirelessKeyType(NM_SETTING_WIRELESS_MAC_ADDRESS))
}
func setSettingWirelessClonedMacAddressJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_CLONED_MAC_ADDRESS, value, getSettingWirelessKeyType(NM_SETTING_WIRELESS_CLONED_MAC_ADDRESS))
}
func setSettingWirelessMacAddressBlacklistJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_MAC_ADDRESS_BLACKLIST, value, getSettingWirelessKeyType(NM_SETTING_WIRELESS_MAC_ADDRESS_BLACKLIST))
}
func setSettingWirelessMtuJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_MTU, value, getSettingWirelessKeyType(NM_SETTING_WIRELESS_MTU))
}
func setSettingWirelessSeenBssidsJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_SEEN_BSSIDS, value, getSettingWirelessKeyType(NM_SETTING_WIRELESS_SEEN_BSSIDS))
}
func setSettingWirelessSecJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_SEC, value, getSettingWirelessKeyType(NM_SETTING_WIRELESS_SEC))
}
func setSettingWirelessHiddenJSON(data _ConnectionData, value string) {
	setConnectionDataKeyJSON(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_HIDDEN, value, getSettingWirelessKeyType(NM_SETTING_WIRELESS_HIDDEN))
}

// Remover
func removeSettingWirelessSsid(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_SSID)
}
func removeSettingWirelessMode(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_MODE)
}
func removeSettingWirelessBand(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_BAND)
}
func removeSettingWirelessChannel(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_CHANNEL)
}
func removeSettingWirelessBssid(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_BSSID)
}
func removeSettingWirelessRate(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_RATE)
}
func removeSettingWirelessTxPower(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_TX_POWER)
}
func removeSettingWirelessMacAddress(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_MAC_ADDRESS)
}
func removeSettingWirelessClonedMacAddress(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_CLONED_MAC_ADDRESS)
}
func removeSettingWirelessMacAddressBlacklist(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_MAC_ADDRESS_BLACKLIST)
}
func removeSettingWirelessMtu(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_MTU)
}
func removeSettingWirelessSeenBssids(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_SEEN_BSSIDS)
}
func removeSettingWirelessSec(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_SEC)
}
func removeSettingWirelessHidden(data _ConnectionData) {
	removeConnectionDataKey(data, NM_SETTING_WIRELESS_SETTING_NAME, NM_SETTING_WIRELESS_HIDDEN)
}
