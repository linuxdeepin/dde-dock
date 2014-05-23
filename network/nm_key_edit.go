package main

import "dlib/dbus"
import "fmt"

// TODO remove
func pageGeneralGetId(con map[string]map[string]dbus.Variant) string {
	defer func() {
		if err := recover(); err != nil {
			logger.Warning("EditorGetID failed:", con, err)
		}
	}()
	return con[fieldConnection]["id"].Value().(string)
}

func getCustomConnectinoType(data connectionData) (connType string) {
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

// TODO
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

func getSettingKeyJSON(data connectionData, field, key string, t ktype) (valueJSON string) {
	// TODO
	value := getSettingKey(data, field, key)

	// var value interface{}
	// if isSettingKeyExists(data, field, key) {
	// 	value = getSettingKey(data, field, key)
	// } else {
	// 	// return default value if the key is not exists
	// 	value = generalGetSettingDefaultValue(field, key)
	// 	valueJSON, _ = marshalJSON(value)
	// 	return
	// }

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

func getSettingKey(data connectionData, field, key string) (value interface{}) {
	// special for vpn plugin keys
	if isSettingVpnPluginKey(field) {
		return getSettingVpnPluginKey(data, field, key)
	}

	value = generalGetSettingDefaultValue(field, key) // get default value firstly

	realField := getRealFieldName(field) // get real name of virtual fields
	fieldData, ok := data[realField]
	if !ok {
		logger.Errorf("invalid field: data[%s]", realField)
		return
	}

	variant, ok := fieldData[key]
	if !ok {
		// not exists, just return nil
		return
	}

	value = variant.Value()

	// logger.Debugf("getSettingKey: data[%s][%s]=%v", field, key, value) // TODO test
	if value == nil {
		logger.Warning("getSettingKey: data[%s][%s] is nil")
	}
	return
}

func setSettingKeyJSON(data connectionData, field, key, valueJSON string, t ktype) (kerr error) {
	if len(valueJSON) == 0 {
		logger.Error("setSettingKeyJSON: valueJSON is empty")
		kerr = fmt.Errorf(NM_KEY_ERROR_INVALID_VALUE)
		return
	}

	// remove connection data key if valueJSON is null or empty
	if isJSONValueMeansToDeleteKey(valueJSON, t) {
		logger.Debugf("removeSettingKey data[%s][%s], valueJSON=%s", field, key, valueJSON)
		removeSettingKey(data, field, key)
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
	logger.Debugf("setSettingKeyJSON data[%s][%s]=%#v, valueJSON=%s", field, key, value, valueJSON) // TODO test
	if isInterfaceNil(value) {
		removeSettingKey(data, field, key)
	} else {
		setSettingKey(data, field, key, value)
	}
	return
}

func setSettingKey(data connectionData, field, key string, value interface{}) {
	// special for vpn plugin keys
	if isSettingVpnPluginKey(field) {
		setSettingVpnPluginKey(data, field, key, value)
		return
	}

	realField := getRealFieldName(field) // get real name of virtual fields
	var fieldData map[string]dbus.Variant
	fieldData, ok := data[realField]
	if !ok {
		logger.Errorf(`set connection data failed, field "%s" is not exits yet`, realField)
		return
	}

	fieldData[key] = dbus.MakeVariant(value)

	// logger.Debugf("setSettingKey: data[%s][%s]=%s", field, key, value) // TODO test
	return
}

func removeSettingKey(data connectionData, field string, keys ...string) {
	// special for vpn plugin keys
	if isSettingVpnPluginKey(field) {
		removeSettingVpnPluginKey(data, field, keys...)
		return
	}

	realField := getRealFieldName(field) // get real name of virtual fields
	fieldData, ok := data[realField]
	if !ok {
		return
	}

	for _, k := range keys {
		delete(fieldData, k)
	}
}

func removeSettingKeyBut(data connectionData, field string, keys ...string) {
	// special for vpn plugin keys
	if isSettingVpnPluginKey(field) {
		removeSettingVpnPluginKeyBut(data, field, keys...)
		return
	}

	realField := getRealFieldName(field) // get real name of virtual fields
	fieldData, ok := data[realField]
	if !ok {
		return
	}

	for k := range fieldData {
		if !isStringInArray(k, keys) {
			delete(fieldData, k)
		}
	}
}

func isSettingKeyExists(data connectionData, field, key string) bool {
	// special for vpn plugin keys
	if isSettingVpnPluginKey(field) {
		return isSettingVpnPluginKeyExists(data, field, key)
	}

	realField := getRealFieldName(field) // get real name of virtual fields
	fieldData, ok := data[realField]
	if !ok {
		return false
	}

	_, ok = fieldData[key]
	if !ok {
		return false
	}

	return true
}

func addSettingField(data connectionData, field string) {
	realField := getRealFieldName(field) // get real name of virtual fields
	var fieldData map[string]dbus.Variant
	fieldData, ok := data[realField]
	if !ok {
		// add field if not exists
		fieldData = make(map[string]dbus.Variant)
		data[realField] = fieldData
	}
}

func removeSettingField(data connectionData, field string) {
	realField := getRealFieldName(field) // get real name of virtual fields
	_, ok := data[realField]
	if ok {
		// remove field if exists
		delete(data, realField)
	}
}

func isSettingFieldExists(data connectionData, field string) bool {
	realField := getRealFieldName(field) // get real name of virtual fields
	_, ok := data[realField]
	return ok
}
