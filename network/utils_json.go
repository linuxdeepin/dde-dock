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

// string[json] -> realdata -> dbus.Variant
func jsonToInterface(s string, t ktype) (v interface{}, err error) {
	switch t {
	case ktypeString:
		v, err = jsonToInterfaceString(s)
	case ktypeByte:
		v, err = jsonToInterfaceByte(s)
	case ktypeInt32:
		v, err = jsonToInterfaceInt32(s)
	case ktypeUint32:
		v, err = jsonToInterfaceUint32(s)
	case ktypeUint64:
		v, err = jsonToInterfaceUint64(s)
	case ktypeBoolean:
		v, err = jsonToInterfaceBoolean(s)
	case ktypeArrayString:
		v, err = jsonToInterfaceArrayString(s)
	case ktypeArrayByte:
		v, err = jsonToInterfaceArrayByte(s)
	case ktypeArrayUint32:
		v, err = jsonToInterfaceArrayUint32(s)
	case ktypeArrayArrayByte:
		v, err = jsonToInterfaceArrayArrayByte(s)
	case ktypeArrayArrayUint32:
		v, err = jsonToInterfaceArrayArrayUint32(s)
	case ktypeDictStringString:
		v, err = jsonToInterfaceDictStringString(s)
	case ktypeIpv6Addresses:
		v, err = jsonToInterfaceIpv6Addresses(s)
	case ktypeIpv6Routes:
		v, err = jsonToInterfaceIpv6Routes(s)
	default:
		err = fmt.Errorf("invalid variant type, %s", s)
	}
	return
}

// dbus.Variant -> realdata -> string[json]
func interfaceToJSON(v interface{}, t ktype) (s string, err error) {
	switch t {
	case ktypeString:
		s, err = interfaceToJSONString(v)
	case ktypeByte:
		s, err = interfaceToJSONByte(v)
	case ktypeInt32:
		s, err = interfaceToJSONInt32(v)
	case ktypeUint32:
		s, err = interfaceToJSONUint32(v)
	case ktypeUint64:
		s, err = interfaceToJSONUint64(v)
	case ktypeBoolean:
		s, err = interfaceToJSONBoolean(v)
	case ktypeArrayString:
		s, err = interfaceToJSONArrayString(v)
	case ktypeArrayByte:
		s, err = interfaceToJSONArrayByte(v)
	case ktypeArrayUint32:
		s, err = interfaceToJSONArrayUint32(v)
	case ktypeArrayArrayByte:
		s, err = interfaceToJSONArrayArrayByte(v)
	case ktypeArrayArrayUint32:
		s, err = interfaceToJSONArrayArrayUint32(v)
	case ktypeDictStringString:
		s, err = interfaceToJSONDictStringString(v)
	case ktypeIpv6Addresses:
		s, err = interfaceToJSONIpv6Addresses(v)
	case ktypeIpv6Routes:
		s, err = interfaceToJSONIpv6Routes(v)
	default:
		err = fmt.Errorf("invalid key type, %v", v)
	}
	return
}

// Convert sepcial key type which wrapped by json to dbus variant's value
// TODO
func jsonToInterfaceString(s string) (v interface{}, err error) {
	// var d string
	// json.Unmarshal([]byte(s), &d)
	// v = d
	v = s
	return
}
func jsonToInterfaceByte(s string) (v interface{}, err error) {
	if len(s) == 0 {
		err = fmt.Errorf("string is empty")
		return
	}
	d := byte(s[0])
	v = d
	return
}
func jsonToInterfaceInt32(s string) (v interface{}, err error) {
	var d int32
	tmpd, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		return
	}
	d = int32(tmpd)
	v = d
	return
}
func jsonToInterfaceUint32(s string) (v interface{}, err error) {
	var d uint32
	tmpd, err := strconv.ParseUint(s, 10, 0)
	if err != nil {
		return
	}
	d = uint32(tmpd)
	v = d
	return
}
func jsonToInterfaceUint64(s string) (v interface{}, err error) {
	var d uint64
	d, err = strconv.ParseUint(s, 10, 0)
	if err != nil {
		return
	}
	v = d
	return
}
func jsonToInterfaceBoolean(s string) (v interface{}, err error) {
	var d bool
	d, err = strconv.ParseBool(s)
	if err != nil {
		return
	}
	v = d
	return
}
func jsonToInterfaceArrayByte(s string) (v interface{}, err error) {
	var d []byte
	// TODO wrap throuh json
	// err = json.Unmarshal([]byte(s), &d)
	d = []byte(s)
	v = d
	return
}
func jsonToInterfaceArrayString(s string) (v interface{}, err error) {
	var d []string
	err = json.Unmarshal([]byte(s), &d)
	v = d
	return
}
func jsonToInterfaceArrayUint32(s string) (v interface{}, err error) {
	var d []uint32
	err = json.Unmarshal([]byte(s), &d)
	v = d
	return
}
func jsonToInterfaceArrayArrayByte(s string) (v interface{}, err error) {
	var d [][]byte
	err = json.Unmarshal([]byte(s), &d)
	v = d
	return
}
func jsonToInterfaceArrayArrayUint32(s string) (v interface{}, err error) {
	var d [][]uint32
	err = json.Unmarshal([]byte(s), &d)
	v = d
	return
}
func jsonToInterfaceDictStringString(s string) (v interface{}, err error) {
	var d map[string]string
	err = json.Unmarshal([]byte(s), &d)
	v = d
	return
}
func jsonToInterfaceIpv6Addresses(s string) (v interface{}, err error) {
	var d Ipv6Addresses
	err = json.Unmarshal([]byte(s), &d)
	v = d
	return
}
func jsonToInterfaceIpv6Routes(s string) (v interface{}, err error) {
	var d Ipv6Routes
	err = json.Unmarshal([]byte(s), &d)
	v = d
	return
}

// Convert dbus variant's value to special key type and wrap to json
// TODO
func interfaceToJSONString(v interface{}) (s string, err error) {
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
func interfaceToJSONByte(v interface{}) (s string, err error) {
	d, err := interfaceToByte(v)
	if err != nil {
		return
	}
	s = string(d)
	return
}
func interfaceToJSONInt32(v interface{}) (s string, err error) {
	d, err := interfaceToInt32(v)
	if err != nil {
		return
	}
	s = strconv.FormatInt(int64(d), 10)
	return
}
func interfaceToJSONUint32(v interface{}) (s string, err error) {
	d, err := interfaceToUint32(v)
	if err != nil {
		return
	}
	s = strconv.FormatUint(uint64(d), 10)
	return
}
func interfaceToJSONUint64(v interface{}) (s string, err error) {
	d, err := interfaceToUint64(v)
	if err != nil {
		return
	}
	s = strconv.FormatUint(d, 10)
	return
}
func interfaceToJSONBoolean(v interface{}) (s string, err error) {
	d, err := interfaceToBoolean(v)
	if err != nil {
		return
	}
	s = strconv.FormatBool(d)
	return
}
func interfaceToJSONArrayByte(v interface{}) (s string, err error) {
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
func interfaceToJSONArrayString(v interface{}) (s string, err error) {
	d, err := interfaceToArrayString(v)
	if err != nil {
		return
	}
	b, err := json.Marshal(d)
	s = string(b)
	return
}
func interfaceToJSONArrayUint32(v interface{}) (s string, err error) {
	d, err := interfaceToArrayUint32(v)
	if err != nil {
		return
	}
	b, err := json.Marshal(d)
	s = string(b)
	return
}
func interfaceToJSONArrayArrayByte(v interface{}) (s string, err error) {
	d, err := interfaceToArrayArrayByte(v)
	if err != nil {
		return
	}
	b, err := json.Marshal(d)
	s = string(b)
	return
}
func interfaceToJSONArrayArrayUint32(v interface{}) (s string, err error) {
	d, err := interfaceToArrayArrayUint32(v)
	if err != nil {
		return
	}
	b, err := json.Marshal(d)
	s = string(b)
	return
}
func interfaceToJSONDictStringString(v interface{}) (s string, err error) {
	d, err := interfaceToDictStringString(v)
	if err != nil {
		return
	}
	b, err := json.Marshal(d)
	s = string(b)
	return
}
func interfaceToJSONIpv6Addresses(v interface{}) (s string, err error) {
	d, err := interfaceToIpv6Addresses(v)
	if err != nil {
		return
	}
	b, err := json.Marshal(d)
	s = string(b)
	return
}
func interfaceToJSONIpv6Routes(v interface{}) (s string, err error) {
	d, err := interfaceToIpv6Routes(v)
	if err != nil {
		return
	}
	b, err := json.Marshal(d)
	s = string(b)
	return
}
