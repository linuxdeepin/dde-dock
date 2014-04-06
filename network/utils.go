package main

import "dlib/dbus"
import "fmt"
import "io"
import "crypto/rand"
import "reflect"

func newUUID() string {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		panic("This can failed?")
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

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

// TODO key: add()
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

func getConnectionDataKeyDefaultValue(field, key string) (value interface{}) {
	switch field {
	default:
		LOGGER.Warning("invalid field name", field)
	case field8021x:
		value = getSetting8021xKeyDefaultValue(key)
	case fieldConnection:
		value = getSettingConnectionKeyDefaultValue(key)
	case fieldIPv4:
		value = getSettingIp4ConfigKeyDefaultValue(key)
	case fieldIPv6:
		value = getSettingIp6ConfigKeyDefaultValue(key)
	case fieldWired:
		value = getSettingWiredKeyDefaultValue(key)
	case fieldWireless:
		value = getSettingWirelessKeyDefaultValue(key)
	case fieldWirelessSecurity:
		value = getSettingWirelessSecurityKeyDefaultValue(key)
	}
	return
}

func getConnectionDataKeyJSON(data _ConnectionData, field, key string, t ktype) (valueJSON string) {
	var value interface{}
	if isConnectionDataKeyExists(data, field, key) {
		value = getConnectionDataKey(data, field, key)
	} else {
		// return default value if the key is not exists
		value = getConnectionDataKeyDefaultValue(field, key)
		// LOGGER.Debugf("default value data[%s][%s]=%#v", field, key, value) // TODO test
		if isInterfaceNil(value) {
			if t == ktypeString || t == ktypeWrapperString || t == ktypeWrapperMacAddress {
				return `""`
			}
		}
	}

	// TODO
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

	// TODO default value for ktypeWrapper
	if valueJSON == `""` {
		// if valueJSON is empty string, just means to remove current key
		if t == ktypeString || t == ktypeWrapperString || t == ktypeWrapperMacAddress {
			removeConnectionDataKey(data, field, key)
			return
		}
	}

	value, err := jsonToKeyValue(valueJSON, t)
	if err != nil {
		LOGGER.Errorf("set connection data failed, valueJSON=%s, ktype=%s, error message:%v",
			valueJSON, getKtypeDescription(t), err)
		return
	}
	// LOGGER.Debugf("setConnectionDataKeyJSON data[%s][%s]=%#v, valueJSON=%s", field, key, value, valueJSON) // TODO test
	if isInterfaceNil(value) {
		// if value == getConnectionDataKeyDefaultValue(field, key) {
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

func isStringInArray(s string, list []string) bool {
	for _, i := range list {
		if i == s {
			return true
		}
	}
	return false
}

func randString(n int) string {
	const alphanum = "0123456789abcdef"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func isInterfaceNil(v interface{}) bool {
	defer func() { recover() }()
	return v == nil || reflect.ValueOf(v).IsNil()
}
