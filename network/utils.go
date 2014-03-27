package main

import "dlib/dbus"
import "fmt"
import "io"
import "crypto/rand"

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

func getConnectionData(data _ConnectionData, field, key string, t ktype) (value string) {
	fieldData, ok := data[field]
	if !ok {
		LOGGER.Errorf("invalid field: data[%s]", field)
		return
	}

	valueVariant, ok := fieldData[key]
	if !ok {
		LOGGER.Errorf("invalid key: data[%s][%s]", field, key)
		return
	}

	value, err := unwrapVariant(valueVariant, t)
	if err != nil {
		LOGGER.Error("get connection data failed:", err)
		return
	}

	// LOGGER.Debugf("getConnectionData: data[%s][%s]=%s", field, key, value) // TODO
	return
}

func setConnectionData(data _ConnectionData, field, key, value string, t ktype) {
	var fieldData map[string]dbus.Variant
	fieldData, ok := data[field]
	if !ok {
		// create field if not exists yet
		fieldData = make(map[string]dbus.Variant)
		data[field] = fieldData
		return
	}

	valueVariant, err := wrapVariant(value, t)
	if err != nil {
		LOGGER.Error("set connection data failed:", err)
		return
	}
	fieldData[key] = valueVariant

	LOGGER.Debugf("setConnectionData: data[%s][%s]=%s", field, key, value) // TODO
	return
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
