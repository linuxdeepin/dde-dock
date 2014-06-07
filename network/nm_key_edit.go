package network

import "dlib/dbus"
import "fmt"

// return custom connection type, and the wrapper type will be ignored
func getCustomConnectionType(data connectionData) (connType string) {
	t := getSettingConnectionType(data)
	switch t {
	case NM_SETTING_WIRED_SETTING_NAME:
		connType = connectionWired
	case NM_SETTING_WIRELESS_SETTING_NAME:
		if isSettingWirelessModeExists(data) {
			switch getSettingWirelessMode(data) {
			case NM_SETTING_WIRELESS_MODE_INFRA:
				connType = connectionWireless
			case NM_SETTING_WIRELESS_MODE_ADHOC:
				connType = connectionWirelessAdhoc
			case NM_SETTING_WIRELESS_MODE_AP:
				connType = connectionWirelessHotspot
			}
		} else {
			connType = connectionWireless
		}
	case NM_SETTING_PPPOE_SETTING_NAME:
		connType = connectionPppoe
	case NM_SETTING_GSM_SETTING_NAME:
		connType = connectionMobileGsm
	case NM_SETTING_CDMA_SETTING_NAME:
		connType = connectionMobileCdma
	case NM_SETTING_VPN_SETTING_NAME:
		switch getSettingVpnServiceType(data) {
		case NM_DBUS_SERVICE_L2TP:
			connType = connectionVpnL2tp
		case NM_DBUS_SERVICE_OPENCONNECT:
			connType = connectionVpnOpenconnect
		case NM_DBUS_SERVICE_OPENVPN:
			connType = connectionVpnOpenvpn
		case NM_DBUS_SERVICE_PPTP:
			connType = connectionVpnPptp
		case NM_DBUS_SERVICE_VPNC:
			connType = connectionVpnVpnc
		}
	}
	if len(connType) == 0 {
		connType = connectionUnknown
	}
	return
}

func isJSONValueMeansToDeleteKey(valueJSON string, t ktype) (doDelete bool) {
	if valueJSON == jsonNull || valueJSON == jsonEmptyString || valueJSON == jsonEmptyArray {
		return true
	}
	switch t {
	case ktypeIpv6Addresses:
	case ktypeIpv6Routes:
	case ktypeWrapperIpv4Dns:
		if valueJSON == `[""]` {
			doDelete = true
		}
	case ktypeWrapperIpv4Addresses:
		if valueJSON == `[{"Address":"","Mask":"","Gateway":""}]` {
			doDelete = true
		}
	case ktypeWrapperIpv4Routes:
		if valueJSON == `[{"Address":"","Mask":"","NextHop":"","Metric":0}]` {
			doDelete = true
		}
	case ktypeWrapperIpv6Dns:
		if valueJSON == `[""]` {
			doDelete = true
		}
	case ktypeWrapperIpv6Addresses:
		if valueJSON == `[{"Address":"","Prefix":0,"Gateway":""}]` {
			doDelete = true
		}
	case ktypeWrapperIpv6Routes:
		if valueJSON == `[{"Address":"","Prefix":0,"NextHop":"","Metric":0}]` {
			doDelete = true
		}
	}
	return
}

func getSettingKeyJSON(data connectionData, section, key string, t ktype) (valueJSON string) {
	if isWrapperKeyType(t) && !isSettingKeyExists(data, section, key) {
		// if the key is not exists and is a wrapper key, get its
		// default value and marshaled to json directly instead of use
		// keyValueToJSON(), which will dispatch wrapper keys
		// specially
		valueJSON, _ = marshalJSON(generalGetSettingDefaultValue(section, key))
		return
	}

	value := getSettingKey(data, section, key)
	valueJSON, err := keyValueToJSON(value, t)
	if err != nil {
		logger.Error("get connection data failed:", err)
		return
	}

	if len(valueJSON) == 0 {
		logger.Error("getSettingKeyJSON: valueJSON is empty")
	}

	return
}

func getSettingKey(data connectionData, section, key string) (value interface{}) {
	// special for vpn plugin keys
	if isSettingVpnPluginKey(section) {
		return getSettingVpnPluginKey(data, section, key)
	}

	value = generalGetSettingDefaultValue(section, key) // get default value firstly

	realSection := getRealSectionName(section) // get real name of virtual sections
	sectionData, ok := data[realSection]
	if !ok {
		logger.Errorf("invalid section: data[%s]", realSection)
		return
	}

	variant, ok := sectionData[key]
	if !ok {
		// not exists, just return nil
		return
	}

	value = variant.Value()

	// logger.Debugf("getSettingKey: data[%s][%s]=%v", section, key, value) // TODO test
	if isInterfaceNil(value) {
		logger.Errorf("getSettingKey: data[%s][%s] is nil", section, key)
	}

	return
}

func setSettingKeyJSON(data connectionData, section, key, valueJSON string, t ktype) (kerr error) {
	if len(valueJSON) == 0 {
		logger.Error("setSettingKeyJSON: valueJSON is empty")
		kerr = fmt.Errorf(NM_KEY_ERROR_INVALID_VALUE)
		return
	}

	// remove connection data key if valueJSON is null or empty
	if isJSONValueMeansToDeleteKey(valueJSON, t) {
		logger.Debugf("json value means to remove key, data[%s][%s]=%#v", section, key, valueJSON)
		removeSettingKey(data, section, key)
		return
	}

	value, err := jsonToKeyValue(valueJSON, t)
	if err != nil {
		// TODO test
		// logger.Errorf("set connection data failed, valueJSON=%s, ktype=%s, error message:%v",
		// valueJSON, getKtypeDescription(t), err)
		kerr = fmt.Errorf(NM_KEY_ERROR_INVALID_VALUE)
		return
	}
	// logger.Debugf("setSettingKeyJSON data[%s][%s]=%#v, valueJSON=%s", section, key, value, valueJSON) // TODO test
	if isInterfaceNil(value) {
		removeSettingKey(data, section, key)
	} else {
		setSettingKey(data, section, key, value)
	}
	return
}

func setSettingKey(data connectionData, section, key string, value interface{}) {
	// special for vpn plugin keys
	if isSettingVpnPluginKey(section) {
		setSettingVpnPluginKey(data, section, key, value)
		return
	}

	realSection := getRealSectionName(section) // get real name of virtual sections
	var sectionData map[string]dbus.Variant
	sectionData, ok := data[realSection]
	if !ok {
		logger.Errorf(`set connection data failed, section "%s" is not exits yet`, realSection)
		return
	}

	sectionData[key] = dbus.MakeVariant(value)

	logger.Debugf("setSettingKey: data[%s][%s]=%#v", section, key, value) // TODO test
	return
}

func removeSettingKey(data connectionData, section string, keys ...string) {
	logger.Debugf("removeSettingKey data[%s], %s", section, keys)

	// special for vpn plugin keys
	if isSettingVpnPluginKey(section) {
		removeSettingVpnPluginKey(data, section, keys...)
		return
	}

	realSection := getRealSectionName(section) // get real name of virtual sections
	sectionData, ok := data[realSection]
	if !ok {
		return
	}

	for _, k := range keys {
		delete(sectionData, k)
	}
}

func removeSettingKeyBut(data connectionData, section string, keys ...string) {
	// special for vpn plugin keys
	if isSettingVpnPluginKey(section) {
		removeSettingVpnPluginKeyBut(data, section, keys...)
		return
	}

	realSection := getRealSectionName(section) // get real name of virtual sections
	sectionData, ok := data[realSection]
	if !ok {
		return
	}

	for k := range sectionData {
		if !isStringInArray(k, keys) {
			delete(sectionData, k)
		}
	}
}

func isSettingKeyExists(data connectionData, section, key string) bool {
	// special for vpn plugin keys
	if isSettingVpnPluginKey(section) {
		return isSettingVpnPluginKeyExists(data, section, key)
	}

	realSection := getRealSectionName(section) // get real name of virtual sections
	sectionData, ok := data[realSection]
	if !ok {
		return false
	}

	_, ok = sectionData[key]
	if !ok {
		return false
	}

	return true
}

func addSettingSection(data connectionData, section string) {
	realSection := getRealSectionName(section) // get real name of virtual sections
	var sectionData map[string]dbus.Variant
	sectionData, ok := data[realSection]
	if !ok {
		// add section if not exists
		sectionData = make(map[string]dbus.Variant)
		data[realSection] = sectionData
	}
}

func removeSettingSection(data connectionData, section string) {
	realSection := getRealSectionName(section) // get real name of virtual sections
	_, ok := data[realSection]
	if ok {
		// remove section if exists
		delete(data, realSection)
	}
}

func isSettingSectionExists(data connectionData, section string) bool {
	realSection := getRealSectionName(section) // get real name of virtual sections
	_, ok := data[realSection]
	return ok
}
