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

func generalIsKeyInSettingField(field, key string) bool {
	if isVirtualKey(field, key) {
		return true
	}
	switch field {
	default:
		LOGGER.Warning("invalid field name", field)
	case field8021x:
		return isKeyInSetting8021x(key)
	case fieldConnection:
		return isKeyInSettingConnection(key)
	case fieldIPv4:
		return isKeyInSettingIp4Config(key)
	case fieldIPv6:
		return isKeyInSettingIp6Config(key)
	case fieldWired:
		return isKeyInSettingWired(key)
	case fieldWireless:
		return isKeyInSettingWireless(key)
	case fieldWirelessSecurity:
		return isKeyInSettingWirelessSecurity(key)
	}
	return false
}

func generalGetSettingKeyType(field, key string) (t ktype) {
	if isVirtualKey(field, key) {
		t = getSettingVkKeyType(field, key)
		return
	}
	switch field {
	default:
		LOGGER.Warning("invalid field name", field)
	case field8021x:
		t = getSetting8021xKeyType(key)
	case fieldConnection:
		t = getSettingConnectionKeyType(key)
	case fieldIPv4:
		t = getSettingIp4ConfigKeyType(key)
	case fieldIPv6:
		t = getSettingIp6ConfigKeyType(key)
	case fieldWired:
		t = getSettingWiredKeyType(key)
	case fieldWireless:
		t = getSettingWirelessKeyType(key)
	case fieldWirelessSecurity:
		t = getSettingWirelessSecurityKeyType(key)
	}
	return
}

func generalGetSettingAvailableKeys(data _ConnectionData, field string) (keys []string) {
	switch field {
	case field8021x:
		keys = getSetting8021xAvailableKeys(data)
	case fieldConnection:
		keys = getSettingConnectionAvailableKeys(data)
	case fieldIPv4:
		keys = getSettingIp4ConfigAvailableKeys(data)
	case fieldIPv6:
		keys = getSettingIp6ConfigAvailableKeys(data)
	case fieldWired:
		keys = getSettingWiredAvailableKeys(data)
	case fieldWireless:
		keys = getSettingWirelessAvailableKeys(data)
	case fieldWirelessSecurity:
		keys = getSettingWirelessSecurityAvailableKeys(data)
	}
	return
}

func generalGetSettingAvailableValues(field, key string) (values []string, customizable bool) {
	if isVirtualKey(field, key) {
		values = generalGetSettingVkAvailableValues(field, key)
		return
	}
	switch field {
	case field8021x:
		values, customizable = getSetting8021xAvailableValues(key)
	case fieldConnection:
		values, customizable = getSettingConnectionAvailableValues(key)
	case fieldIPv4:
		values, customizable = getSettingIp4ConfigAvailableValues(key)
	case fieldIPv6:
		values, customizable = getSettingIp6ConfigAvailableValues(key)
	case fieldWired:
		// values,customizable = getSettingWiredAvailableValues(key)
	case fieldWireless:
		values, customizable = getSettingWirelessAvailableValues(key)
	case fieldWirelessSecurity:
		values, customizable = getSettingWirelessSecurityAvailableValues(key)
	}
	return
}

func generalCheckSettingValues(data _ConnectionData, field string) (errs map[string]string) {
	switch field {
	default:
		LOGGER.Error("updatePropErrors: invalid field name", field)
	case field8021x:
		errs = checkSetting8021xValues(data)
	case fieldConnection:
		errs = checkSettingConnectionValues(data)
	case fieldIPv4:
		errs = checkSettingIp4ConfigValues(data)
	case fieldIPv6:
		errs = checkSettingIp6ConfigValues(data)
	case fieldWired:
		errs = checkSettingWiredValues(data)
	case fieldWireless:
		errs = checkSettingWirelessValues(data)
	case fieldWirelessSecurity:
		errs = checkSettingWirelessSecurityValues(data)
	}
	return
}

func generalGetSettingKeyJSON(data _ConnectionData, field, key string) (valueJSON string) {
	if isVirtualKey(field, key) {
		valueJSON = generalGetVirtualKeyJSON(data, field, key)
		return
	}
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

func generalSetSettingKeyJSON(data _ConnectionData, field, key, valueJSON string) {
	if isVirtualKey(field, key) {
		generalSetVirtualKeyJSON(data, field, key, valueJSON)
		return
	}
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

func getSettingKeyDefaultValueJSON(field, key string) (valueJSON string) {
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

func getSettingKeyJSON(data _ConnectionData, field, key string, t ktype) (valueJSON string) {
	var value interface{}
	if isSettingKeyExists(data, field, key) {
		value = getSettingKey(data, field, key)
	} else {
		// return default value if the key is not exists
		valueJSON = getSettingKeyDefaultValueJSON(field, key)
		return
	}

	valueJSON, err := keyValueToJSON(value, t)
	if err != nil {
		LOGGER.Error("get connection data failed:", err)
		return
	}

	if len(valueJSON) == 0 {
		LOGGER.Error("getSettingKeyJSON: valueJSON is empty")
	}

	return
}

func setSettingKeyJSON(data _ConnectionData, field, key, valueJSON string, t ktype) {
	if len(valueJSON) == 0 {
		LOGGER.Error("setSettingKeyJSON: valueJSON is empty")
		return
	}

	// remove connection data key if valueJSON is null or empty
	if isJSONKeyValueMeansToDeleteKey(valueJSON, t) {
		removeSettingKey(data, field, key)
		return
	}

	value, err := jsonToKeyValue(valueJSON, t)
	if err != nil {
		LOGGER.Errorf("set connection data failed, valueJSON=%s, ktype=%s, error message:%v",
			valueJSON, getKtypeDescription(t), err)
		return
	}
	// LOGGER.Debugf("setSettingKeyJSON data[%s][%s]=%#v, valueJSON=%s", field, key, value, valueJSON) // TODO test
	if isInterfaceNil(value) {
		removeSettingKey(data, field, key)
	} else {
		setSettingKey(data, field, key, value)
	}
	return
}

func getSettingKey(data _ConnectionData, field, key string) (value interface{}) {
	fieldData, ok := data[field]
	if !ok {
		LOGGER.Errorf("invalid field: data[%s]", field)
		return
	}

	variant, ok := fieldData[key]
	if !ok {
		return
	}

	value = variant.Value()

	// LOGGER.Debugf("getSettingKey: data[%s][%s]=%v", field, key, value) // TODO test
	return
}

func setSettingKey(data _ConnectionData, field, key string, value interface{}) {
	var fieldData map[string]dbus.Variant
	fieldData, ok := data[field]
	if !ok {
		LOGGER.Error(fmt.Errorf(`set connection data failed, field "%s" is not exits yet`, field))
		return
	}

	fieldData[key] = dbus.MakeVariant(value)

	// LOGGER.Debugf("setSettingKey: data[%s][%s]=%s", field, key, value) // TODO test
	return
}

func removeSettingKey(data _ConnectionData, field, key string) {
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
	var fieldData map[string]dbus.Variant
	fieldData, ok := data[field]
	if !ok {
		// add field if not exists
		fieldData = make(map[string]dbus.Variant)
		data[field] = fieldData
	}
}

func removeSettingField(data _ConnectionData, field string) {
	_, ok := data[field]
	if !ok {
		// remove field if exists
		delete(data, field)
	}
}

func isSettingFieldExists(data _ConnectionData, field string) bool {
	_, ok := data[field]
	return ok
}
