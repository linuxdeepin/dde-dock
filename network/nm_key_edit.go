package main

import "dlib/dbus"
import "fmt"

// TODO remove
func pageGeneralGetId(con map[string]map[string]dbus.Variant) string {
	defer func() {
		if err := recover(); err != nil {
			Logger.Warning("EditorGetID failed:", con, err)
		}
	}()
	return con[fieldConnection]["id"].Value().(string)
}

func generalGetConnectionType(data _ConnectionData) (connType string) {
	switch getSettingConnectionType(data) {
	case typeWired:
		connType = typeWired
	case typeWireless:
		connType = typeWireless
	case typePppoe:
		connType = typePppoe
	case typeVpn:
		switch getSettingVpnServiceType(data) {
		case NM_DBUS_SERVICE_L2TP:
			connType = typeVpnL2tp
		case NM_DBUS_SERVICE_OPENCONNECT:
			connType = typeVpnOpenconnect
		case NM_DBUS_SERVICE_OPENVPN:
			connType = typeVpnOpenvpn
		case NM_DBUS_SERVICE_PPTP:
			connType = typeVpnPptp
		case NM_DBUS_SERVICE_VPNC:
			connType = typeVpnVpnc
		}
	}
	if len(connType) == 0 {
		Logger.Error("get connection type failed")
	}
	return
}

func isJSONKeyValueMeansToDeleteKey(valueJSON string, t ktype) (doDelete bool) {
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

func getSettingKeyJSON(data _ConnectionData, field, key string, t ktype) (valueJSON string) {
	_filed := field
	field = getRealFieldName(field) // get real field name
	var value interface{}
	if isSettingKeyExists(data, field, key) {
		value = getSettingKey(data, field, key)
	} else {
		// return default value if the key is not exists
		valueJSON = getSettingKeyDefaultValueJSON(_filed, key)
		return
	}

	valueJSON, err := keyValueToJSON(value, t)
	if err != nil {
		Logger.Error("get connection data failed:", err)
		return
	}

	if len(valueJSON) == 0 {
		Logger.Error("getSettingKeyJSON: valueJSON is empty")
	}

	return
}

func setSettingKeyJSON(data _ConnectionData, field, key, valueJSON string, t ktype) {
	field = getRealFieldName(field) // get real field name
	if len(valueJSON) == 0 {
		Logger.Error("setSettingKeyJSON: valueJSON is empty")
		return
	}

	// remove connection data key if valueJSON is null or empty
	if isJSONKeyValueMeansToDeleteKey(valueJSON, t) {
		removeSettingKey(data, field, key)
		return
	}

	value, err := jsonToKeyValue(valueJSON, t)
	if err != nil {
		Logger.Errorf("set connection data failed, valueJSON=%s, ktype=%s, error message:%v",
			valueJSON, getKtypeDescription(t), err)
		return
	}
	// Logger.Debugf("setSettingKeyJSON data[%s][%s]=%#v, valueJSON=%s", field, key, value, valueJSON) // TODO test
	if isInterfaceNil(value) {
		removeSettingKey(data, field, key)
	} else {
		setSettingKey(data, field, key, value)
	}
	return
}

func getSettingKey(data _ConnectionData, field, key string) (value interface{}) {
	field = getRealFieldName(field) // get real field name
	fieldData, ok := data[field]
	if !ok {
		Logger.Errorf("invalid field: data[%s]", field)
		return
	}

	variant, ok := fieldData[key]
	if !ok {
		return
	}

	value = variant.Value()

	// Logger.Debugf("getSettingKey: data[%s][%s]=%v", field, key, value) // TODO test
	return
}

func setSettingKey(data _ConnectionData, field, key string, value interface{}) {
	field = getRealFieldName(field) // get real field name
	var fieldData map[string]dbus.Variant
	fieldData, ok := data[field]
	if !ok {
		Logger.Error(fmt.Errorf(`set connection data failed, field "%s" is not exits yet`, field))
		return
	}

	fieldData[key] = dbus.MakeVariant(value)

	Logger.Debugf("setSettingKey: data[%s][%s]=%s", field, key, value) // TODO test
	return
}

func removeSettingKey(data _ConnectionData, field, key string) {
	field = getRealFieldName(field) // get real field name
	fieldData, ok := data[field]
	if !ok {
		return
	}

	_, ok = fieldData[key]
	if !ok {
		return
	}

	delete(fieldData, key)
}

func removeSettingKeyBut(data _ConnectionData, field string, keys ...string) {
	field = getRealFieldName(field) // get real field name
	fieldData, ok := data[field]
	if !ok {
		return
	}

	for k := range fieldData {
		if isStringInArray(k, keys) {
			delete(fieldData, k)
		}
	}
}

func isSettingKeyExists(data _ConnectionData, field, key string) bool {
	field = getRealFieldName(field) // get real field name
	fieldData, ok := data[field]
	if !ok {
		return false
	}

	_, ok = fieldData[key]
	if !ok {
		return false
	}

	return true
}

func addSettingField(data _ConnectionData, field string) {
	field = getRealFieldName(field) // get real field name
	var fieldData map[string]dbus.Variant
	fieldData, ok := data[field]
	if !ok {
		// add field if not exists
		fieldData = make(map[string]dbus.Variant)
		data[field] = fieldData
	}
}

func removeSettingField(data _ConnectionData, field string) {
	field = getRealFieldName(field) // get real field name
	_, ok := data[field]
	if ok {
		// remove field if exists
		delete(data, field)
	}
}

func isSettingFieldExists(data _ConnectionData, field string) bool {
	field = getRealFieldName(field) // get real field name
	_, ok := data[field]
	return ok
}
