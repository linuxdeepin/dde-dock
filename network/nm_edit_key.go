package main

import "dlib/dbus"
import "fmt"

// TODO remove
func pageGeneralGetId(con map[string]map[string]dbus.Variant) string {
	defer func() {
		if err := recover(); err != nil {
			LOGGER.Warning("EditorGetID failed:", con, err)
		}
	}()
	return con[fieldConnection]["id"].Value().(string)
}

func addConnectionDataField(data _ConnectionData, field string) {
	var fieldData map[string]dbus.Variant
	fieldData, ok := data[field]
	if !ok {
		// add field if not exists
		fieldData = make(map[string]dbus.Variant)
		data[field] = fieldData
	}
}

func removeConnectionDataField(data _ConnectionData, field string) {
	_, ok := data[field]
	if !ok {
		// remove field if exists
		delete(data, field)
	}
}

func isConnectionDataFieldExists(data _ConnectionData, field string) bool {
	_, ok := data[field]
	return ok
}

func isConnectionDataKeyExists(data _ConnectionData, field, key string) bool {
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

func generalGetKeyJSON(data _ConnectionData, field, key string) (valueJSON string) {
	switch field {
	default:
		LOGGER.Warning("invalid field name", field)
	case field8021x:
		valueJSON = generalGetSetting8021xKeyJSON(data, key)
	case fieldConnection:
		valueJSON = generalGetSettingConnectionKeyJSON(data, key)
	case fieldIPv4:
		valueJSON = generalGetSettingIp4ConfigKeyJSON(data, key)
	case fieldIPv6:
		valueJSON = generalGetSettingIp6ConfigKeyJSON(data, key)
	case fieldWired:
		valueJSON = generalGetSettingWiredKeyJSON(data, key)
	case fieldWireless:
		valueJSON = generalGetSettingWirelessKeyJSON(data, key)
	case fieldWirelessSecurity:
		valueJSON = generalGetSettingWirelessSecurityKeyJSON(data, key)
	}
	return
}

func generalSetKeyJSON(data _ConnectionData, field, key, valueJSON string) {
	switch field {
	default:
		LOGGER.Warning("invalid field name", field)
	case field8021x:
		generalSetSetting8021xKeyJSON(data, key, valueJSON)
	case fieldConnection:
		generalSetSettingConnectionKeyJSON(data, key, valueJSON)
	case fieldIPv4:
		generalSetSettingIp4ConfigKeyJSON(data, key, valueJSON)
	case fieldIPv6:
		generalSetSettingIp6ConfigKeyJSON(data, key, valueJSON)
	case fieldWired:
		generalSetSettingWiredKeyJSON(data, key, valueJSON)
	case fieldWireless:
		generalSetSettingWirelessKeyJSON(data, key, valueJSON)
	case fieldWirelessSecurity:
		generalSetSettingWirelessSecurityKeyJSON(data, key, valueJSON)
	}
}

func getConnectionDataKeyDefaultValueJSON(field, key string) (valueJSON string) {
	switch field {
	default:
		LOGGER.Warning("invalid field name", field)
	case field8021x:
		valueJSON = getSetting8021xKeyDefaultValueJSON(key)
	case fieldConnection:
		valueJSON = getSettingConnectionKeyDefaultValueJSON(key)
	case fieldIPv4:
		valueJSON = getSettingIp4ConfigKeyDefaultValueJSON(key)
	case fieldIPv6:
		valueJSON = getSettingIp6ConfigKeyDefaultValueJSON(key)
	case fieldWired:
		valueJSON = getSettingWiredKeyDefaultValueJSON(key)
	case fieldWireless:
		valueJSON = getSettingWirelessKeyDefaultValueJSON(key)
	case fieldWirelessSecurity:
		valueJSON = getSettingWirelessSecurityKeyDefaultValueJSON(key)
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

func getConnectionDataKeyJSON(data _ConnectionData, field, key string, t ktype) (valueJSON string) {
	var value interface{}
	if isConnectionDataKeyExists(data, field, key) {
		value = getConnectionDataKey(data, field, key)
	} else {
		// return default value if the key is not exists
		valueJSON = getConnectionDataKeyDefaultValueJSON(field, key)
		return
	}

	valueJSON, err := keyValueToJSON(value, t)
	if err != nil {
		LOGGER.Error("get connection data failed:", err)
		return
	}

	if len(valueJSON) == 0 {
		LOGGER.Error("getConnectionDataKeyJSON: valueJSON is empty")
	}

	return
}

func setConnectionDataKeyJSON(data _ConnectionData, field, key, valueJSON string, t ktype) {
	if len(valueJSON) == 0 {
		LOGGER.Error("setConnectionDataKeyJSON: valueJSON is empty")
		return
	}

	// remove connection data key if valueJSON is null or empty
	if isJSONKeyValueMeansToDeleteKey(valueJSON, t) {
		removeConnectionDataKey(data, field, key)
		return
	}

	value, err := jsonToKeyValue(valueJSON, t)
	if err != nil {
		LOGGER.Errorf("set connection data failed, valueJSON=%s, ktype=%s, error message:%v",
			valueJSON, getKtypeDescription(t), err)
		return
	}
	// LOGGER.Debugf("setConnectionDataKeyJSON data[%s][%s]=%#v, valueJSON=%s", field, key, value, valueJSON) // TODO test
	if isInterfaceNil(value) {
		removeConnectionDataKey(data, field, key)
	} else {
		setConnectionDataKey(data, field, key, value)
	}
	return
}

func getConnectionDataKey(data _ConnectionData, field, key string) (value interface{}) {
	fieldData, ok := data[field]
	if !ok {
		LOGGER.Errorf("invalid field: data[%s]", field)
		return
	}

	variant, ok := fieldData[key]
	if !ok {
		LOGGER.Errorf("invalid key: data[%s][%s]", field, key)
		return
	}

	value = variant.Value()

	// LOGGER.Debugf("getConnectionDataKey: data[%s][%s]=%v", field, key, value) // TODO test
	return
}

func setConnectionDataKey(data _ConnectionData, field, key string, value interface{}) {
	var fieldData map[string]dbus.Variant
	fieldData, ok := data[field]
	if !ok {
		LOGGER.Error(fmt.Errorf(`set connection data failed, field "%s" is not exits yet`, field))
		return
	}

	fieldData[key] = dbus.MakeVariant(value)

	// LOGGER.Debugf("setConnectionDataKey: data[%s][%s]=%s", field, key, value) // TODO test
	return
}

func removeConnectionDataKey(data _ConnectionData, field, key string) {
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

func removeConnectionDataKeyBut(data _ConnectionData, field string, keys ...string) {
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
