package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// TODO
// arrayByteToVariant, variantToArrayByte
// wrapArrayByte, unwrapArrayByte, wrapArrayByteByJSON
// ktypeStringToVariant

// string[json] -> realdata -> dbus.Variant.Value()
func jsonToKeyValue(s string, t ktype) (v interface{}, err error) {
	switch t {
	case ktypeString:
		v, err = jsonToKeyValueString(s)
	case ktypeByte:
		v, err = jsonToKeyValueByte(s)
	case ktypeInt32:
		v, err = jsonToKeyValueInt32(s)
	case ktypeUint32:
		v, err = jsonToKeyValueUint32(s)
	case ktypeUint64:
		v, err = jsonToKeyValueUint64(s)
	case ktypeBoolean:
		v, err = jsonToKeyValueBoolean(s)
	case ktypeArrayString:
		v, err = jsonToKeyValueArrayString(s)
	case ktypeArrayByte:
		v, err = jsonToKeyValueArrayByte(s)
	case ktypeArrayUint32:
		v, err = jsonToKeyValueArrayUint32(s)
	case ktypeArrayArrayByte:
		v, err = jsonToKeyValueArrayArrayByte(s)
	case ktypeArrayArrayUint32:
		v, err = jsonToKeyValueArrayArrayUint32(s)
	case ktypeDictStringString:
		v, err = jsonToKeyValueDictStringString(s)
	case ktypeIpv6Addresses:
		v, err = jsonToKeyValueIpv6Addresses(s)
	case ktypeIpv6Routes:
		v, err = jsonToKeyValueIpv6Routes(s)
	default:
		err = fmt.Errorf("invalid variant type, %s", s)
	}
	return
}

// dbus.Variant.Value() -> realdata -> string[json]
func keyValueToJSON(v interface{}, t ktype) (s string, err error) {
	switch t {
	case ktypeString:
		s, err = keyValueToJSONString(v)
	case ktypeByte:
		s, err = keyValueToJSONByte(v)
	case ktypeInt32:
		s, err = keyValueToJSONInt32(v)
	case ktypeUint32:
		s, err = keyValueToJSONUint32(v)
	case ktypeUint64:
		s, err = keyValueToJSONUint64(v)
	case ktypeBoolean:
		s, err = keyValueToJSONBoolean(v)
	case ktypeArrayString:
		s, err = keyValueToJSONArrayString(v)
	case ktypeArrayByte:
		s, err = keyValueToJSONArrayByte(v)
	case ktypeArrayUint32:
		s, err = keyValueToJSONArrayUint32(v)
	case ktypeArrayArrayByte:
		s, err = keyValueToJSONArrayArrayByte(v)
	case ktypeArrayArrayUint32:
		s, err = keyValueToJSONArrayArrayUint32(v)
	case ktypeDictStringString:
		s, err = keyValueToJSONDictStringString(v)
	case ktypeIpv6Addresses:
		s, err = keyValueToJSONIpv6Addresses(v)
	case ktypeIpv6Routes:
		s, err = keyValueToJSONIpv6Routes(v)
	default:
		err = fmt.Errorf("invalid key type, %v", v)
	}
	return
}

// Convert sepcial key type which wrapped by json to dbus variant's value
// TODO
func jsonToKeyValueString(s string) (v interface{}, err error) {
	// var d string
	// json.Unmarshal([]byte(s), &d)
	// v = d
	v = s
	return
}
func jsonToKeyValueByte(s string) (v interface{}, err error) {
	if len(s) == 0 {
		err = fmt.Errorf("string is empty")
		return
	}
	d := byte(s[0])
	v = d
	return
}
func jsonToKeyValueInt32(s string) (v interface{}, err error) {
	var d int32
	tmpd, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		return
	}
	d = int32(tmpd)
	v = d
	return
}
func jsonToKeyValueUint32(s string) (v interface{}, err error) {
	var d uint32
	tmpd, err := strconv.ParseUint(s, 10, 0)
	if err != nil {
		return
	}
	d = uint32(tmpd)
	v = d
	return
}
func jsonToKeyValueUint64(s string) (v interface{}, err error) {
	var d uint64
	d, err = strconv.ParseUint(s, 10, 0)
	if err != nil {
		return
	}
	v = d
	return
}
func jsonToKeyValueBoolean(s string) (v interface{}, err error) {
	var d bool
	d, err = strconv.ParseBool(s)
	if err != nil {
		return
	}
	v = d
	return
}
func jsonToKeyValueArrayByte(s string) (v interface{}, err error) {
	var d []byte
	// TODO wrap throuh json
	// err = json.Unmarshal([]byte(s), &d)
	d = []byte(s)
	v = d
	return
}
func jsonToKeyValueArrayString(s string) (v interface{}, err error) {
	var d []string
	err = json.Unmarshal([]byte(s), &d)
	v = d
	return
}
func jsonToKeyValueArrayUint32(s string) (v interface{}, err error) {
	var d []uint32
	err = json.Unmarshal([]byte(s), &d)
	v = d
	return
}
func jsonToKeyValueArrayArrayByte(s string) (v interface{}, err error) {
	var d [][]byte
	err = json.Unmarshal([]byte(s), &d)
	v = d
	return
}
func jsonToKeyValueArrayArrayUint32(s string) (v interface{}, err error) {
	var d [][]uint32
	err = json.Unmarshal([]byte(s), &d)
	v = d
	return
}
func jsonToKeyValueDictStringString(s string) (v interface{}, err error) {
	var d map[string]string
	err = json.Unmarshal([]byte(s), &d)
	v = d
	return
}
func jsonToKeyValueIpv6Addresses(s string) (v interface{}, err error) {
	var d Ipv6Addresses
	err = json.Unmarshal([]byte(s), &d)
	v = d
	return
}
func jsonToKeyValueIpv6Routes(s string) (v interface{}, err error) {
	var d Ipv6Routes
	err = json.Unmarshal([]byte(s), &d)
	v = d
	return
}

// Convert dbus variant's value to special key type and wrap to json
// TODO
func keyValueToJSONString(v interface{}) (s string, err error) {
	// d, err := interfaceToString(v)
	// if err != nil {
	// 	return
	// }
	// b, err := json.Marshal(d)
	// if err != nil {
	// 	return
	// }
	// s = string(b)
	// return
	return interfaceToString(v)
}
func keyValueToJSONByte(v interface{}) (s string, err error) {
	d, err := interfaceToByte(v)
	if err != nil {
		return
	}
	s = string(d)
	return
}
func keyValueToJSONInt32(v interface{}) (s string, err error) {
	d, err := interfaceToInt32(v)
	if err != nil {
		return
	}
	s = strconv.FormatInt(int64(d), 10)
	return
}
func keyValueToJSONUint32(v interface{}) (s string, err error) {
	d, err := interfaceToUint32(v)
	if err != nil {
		return
	}
	s = strconv.FormatUint(uint64(d), 10)
	return
}
func keyValueToJSONUint64(v interface{}) (s string, err error) {
	d, err := interfaceToUint64(v)
	if err != nil {
		return
	}
	s = strconv.FormatUint(d, 10)
	return
}
func keyValueToJSONBoolean(v interface{}) (s string, err error) {
	d, err := interfaceToBoolean(v)
	if err != nil {
		return
	}
	s = strconv.FormatBool(d)
	return
}
func keyValueToJSONArrayByte(v interface{}) (s string, err error) {
	d, err := interfaceToArrayByte(v)
	if err != nil {
		return
	}
	// TODO unwrap throuh json
	// b, err := json.Marshal(d)
	// s = string(b)
	s = string(d)
	return
}
func keyValueToJSONArrayString(v interface{}) (s string, err error) {
	d, err := interfaceToArrayString(v)
	if err != nil {
		return
	}
	b, err := json.Marshal(d)
	s = string(b)
	return
}
func keyValueToJSONArrayUint32(v interface{}) (s string, err error) {
	d, err := interfaceToArrayUint32(v)
	if err != nil {
		return
	}
	b, err := json.Marshal(d)
	s = string(b)
	return
}
func keyValueToJSONArrayArrayByte(v interface{}) (s string, err error) {
	d, err := interfaceToArrayArrayByte(v)
	if err != nil {
		return
	}
	b, err := json.Marshal(d)
	s = string(b)
	return
}
func keyValueToJSONArrayArrayUint32(v interface{}) (s string, err error) {
	d, err := interfaceToArrayArrayUint32(v)
	if err != nil {
		return
	}
	b, err := json.Marshal(d)
	s = string(b)
	return
}
func keyValueToJSONDictStringString(v interface{}) (s string, err error) {
	d, err := interfaceToDictStringString(v)
	if err != nil {
		return
	}
	b, err := json.Marshal(d)
	s = string(b)
	return
}
func keyValueToJSONIpv6Addresses(v interface{}) (s string, err error) {
	d, err := interfaceToIpv6Addresses(v)
	if err != nil {
		return
	}
	b, err := json.Marshal(d)
	s = string(b)
	return
}
func keyValueToJSONIpv6Routes(v interface{}) (s string, err error) {
	d, err := interfaceToIpv6Routes(v)
	if err != nil {
		return
	}
	b, err := json.Marshal(d)
	s = string(b)
	return
}
