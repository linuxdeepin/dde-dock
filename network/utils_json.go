package main

import (
	"encoding/json"
	"fmt"
	// "strconv"
)

// TODO
// arrayByteToVariant, variantToArrayByte
// wrapArrayByte, unwrapArrayByte, wrapArrayByteByJSON
// ktypeStringToVariant

// string[json] -> realdata -> dbus.Variant.Value()
func jsonToKeyValue(jsonStr string, t ktype) (v interface{}, err error) {
	switch t {
	default:
		err = fmt.Errorf("invalid variant type, %jsonStr", jsonStr)
	case ktypeString:
		v, err = jsonToKeyValueString(jsonStr)
	case ktypeByte:
		v, err = jsonToKeyValueByte(jsonStr)
	case ktypeInt32:
		v, err = jsonToKeyValueInt32(jsonStr)
	case ktypeUint32:
		v, err = jsonToKeyValueUint32(jsonStr)
	case ktypeUint64:
		v, err = jsonToKeyValueUint64(jsonStr)
	case ktypeBoolean:
		v, err = jsonToKeyValueBoolean(jsonStr)
	case ktypeArrayString:
		v, err = jsonToKeyValueArrayString(jsonStr)
	case ktypeArrayByte:
		v, err = jsonToKeyValueArrayByte(jsonStr)
	case ktypeArrayUint32:
		v, err = jsonToKeyValueArrayUint32(jsonStr)
	case ktypeArrayArrayByte:
		v, err = jsonToKeyValueArrayArrayByte(jsonStr)
	case ktypeArrayArrayUint32:
		v, err = jsonToKeyValueArrayArrayUint32(jsonStr)
	case ktypeDictStringString:
		v, err = jsonToKeyValueDictStringString(jsonStr)
	case ktypeIpv6Addresses:
		v, err = jsonToKeyValueIpv6Addresses(jsonStr)
	case ktypeIpv6Routes:
		v, err = jsonToKeyValueIpv6Routes(jsonStr)
	}
	return
}

// dbus.Variant.Value() -> realdata -> string[json]
func keyValueToJSON(v interface{}, t ktype) (jsonStr string, err error) {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return
	}
	jsonStr = string(jsonBytes)
	return
	// switch t {
	// default:
	// 	err = fmt.Errorf("invalid key type, %v", v)
	// case ktypeString:
	// 	jsonStr, err = keyValueToJSONString(v)
	// case ktypeByte:
	// 	jsonStr, err = keyValueToJSONByte(v)
	// case ktypeInt32:
	// 	jsonStr, err = keyValueToJSONInt32(v)
	// case ktypeUint32:
	// 	jsonStr, err = keyValueToJSONUint32(v)
	// case ktypeUint64:
	// 	jsonStr, err = keyValueToJSONUint64(v)
	// case ktypeBoolean:
	// 	jsonStr, err = keyValueToJSONBoolean(v)
	// case ktypeArrayString:
	// 	jsonStr, err = keyValueToJSONArrayString(v)
	// case ktypeArrayByte:
	// 	jsonStr, err = keyValueToJSONArrayByte(v)
	// case ktypeArrayUint32:
	// 	jsonStr, err = keyValueToJSONArrayUint32(v)
	// case ktypeArrayArrayByte:
	// 	jsonStr, err = keyValueToJSONArrayArrayByte(v)
	// case ktypeArrayArrayUint32:
	// 	jsonStr, err = keyValueToJSONArrayArrayUint32(v)
	// case ktypeDictStringString:
	// 	jsonStr, err = keyValueToJSONDictStringString(v)
	// case ktypeIpv6Addresses:
	// 	jsonStr, err = keyValueToJSONIpv6Addresses(v)
	// case ktypeIpv6Routes:
	// 	jsonStr, err = keyValueToJSONIpv6Routes(v)
	// }

	// d, err := interfaceToString(v)
	// if err != nil {
	// 	return
	// }
}

// Convert sepcial key type which wrapped by json to dbus variant'jsonStr value
// TODO
func jsonToKeyValueString(jsonStr string) (v string, err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	return
}
func jsonToKeyValueByte(jsonStr string) (v byte, err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	return
}
func jsonToKeyValueInt32(jsonStr string) (v int32, err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	return
}
func jsonToKeyValueUint32(jsonStr string) (v uint32, err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	return
}
func jsonToKeyValueUint64(jsonStr string) (v uint64, err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	return
}
func jsonToKeyValueBoolean(jsonStr string) (v bool, err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	return
}
func jsonToKeyValueArrayByte(jsonStr string) (v []byte, err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	return
}
func jsonToKeyValueArrayString(jsonStr string) (v []string, err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	return
}
func jsonToKeyValueArrayUint32(jsonStr string) (v []uint32, err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	return
}
func jsonToKeyValueArrayArrayByte(jsonStr string) (v [][]byte, err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	return
}
func jsonToKeyValueArrayArrayUint32(jsonStr string) (v [][]uint32, err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	return
}
func jsonToKeyValueDictStringString(jsonStr string) (v map[string]string, err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	return
}
func jsonToKeyValueIpv6Addresses(jsonStr string) (v Ipv6Addresses, err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	return
}
func jsonToKeyValueIpv6Routes(jsonStr string) (v Ipv6Routes, err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	return
}

// Convert dbus variant'jsonStr value to special key type and wrap to json
// TODO
// func keyValueToJSONString(v interface{}) (jsonStr string, err error) {
// 	d, err := interfaceToString(v)
// 	if err != nil {
// 		return
// 	}
// 	b, err := json.Marshal(d)
// 	if err != nil {
// 		return
// 	}
// 	jsonStr = string(b)
// 	return
// 	// return interfaceToString(v)
// }
// func keyValueToJSONByte(v interface{}) (jsonStr string, err error) {
// 	d, err := interfaceToByte(v)
// 	if err != nil {
// 		return
// 	}
// 	jsonStr = string(d)
// 	return
// }
// func keyValueToJSONInt32(v interface{}) (jsonStr string, err error) {
// 	d, err := interfaceToInt32(v)
// 	if err != nil {
// 		return
// 	}
// 	jsonStr = strconv.FormatInt(int64(d), 10)
// 	return
// }
// func keyValueToJSONUint32(v interface{}) (jsonStr string, err error) {
// 	d, err := interfaceToUint32(v)
// 	if err != nil {
// 		return
// 	}
// 	jsonStr = strconv.FormatUint(uint64(d), 10)
// 	return
// }
// func keyValueToJSONUint64(v interface{}) (jsonStr string, err error) {
// 	d, err := interfaceToUint64(v)
// 	if err != nil {
// 		return
// 	}
// 	jsonStr = strconv.FormatUint(d, 10)
// 	return
// }
// func keyValueToJSONBoolean(v interface{}) (jsonStr string, err error) {
// 	d, err := interfaceToBoolean(v)
// 	if err != nil {
// 		return
// 	}
// 	jsonStr = strconv.FormatBool(d)
// 	return
// }
// func keyValueToJSONArrayByte(v interface{}) (jsonStr string, err error) {
// 	d, err := interfaceToArrayByte(v)
// 	if err != nil {
// 		return
// 	}
// 	// TODO unwrap throuh json
// 	// b, err := json.Marshal(d)
// 	// jsonStr = string(b)
// 	jsonStr = string(d)
// 	return
// }
// func keyValueToJSONArrayString(v interface{}) (jsonStr string, err error) {
// 	d, err := interfaceToArrayString(v)
// 	if err != nil {
// 		return
// 	}
// 	b, err := json.Marshal(d)
// 	jsonStr = string(b)
// 	return
// }
// func keyValueToJSONArrayUint32(v interface{}) (jsonStr string, err error) {
// 	d, err := interfaceToArrayUint32(v)
// 	if err != nil {
// 		return
// 	}
// 	b, err := json.Marshal(d)
// 	jsonStr = string(b)
// 	return
// }
// func keyValueToJSONArrayArrayByte(v interface{}) (jsonStr string, err error) {
// 	d, err := interfaceToArrayArrayByte(v)
// 	if err != nil {
// 		return
// 	}
// 	b, err := json.Marshal(d)
// 	jsonStr = string(b)
// 	return
// }
// func keyValueToJSONArrayArrayUint32(v interface{}) (jsonStr string, err error) {
// 	d, err := interfaceToArrayArrayUint32(v)
// 	if err != nil {
// 		return
// 	}
// 	b, err := json.Marshal(d)
// 	jsonStr = string(b)
// 	return
// }
// func keyValueToJSONDictStringString(v interface{}) (jsonStr string, err error) {
// 	d, err := interfaceToDictStringString(v)
// 	if err != nil {
// 		return
// 	}
// 	b, err := json.Marshal(d)
// 	jsonStr = string(b)
// 	return
// }
// func keyValueToJSONIpv6Addresses(v interface{}) (jsonStr string, err error) {
// 	d, err := interfaceToIpv6Addresses(v)
// 	if err != nil {
// 		return
// 	}
// 	b, err := json.Marshal(d)
// 	jsonStr = string(b)
// 	return
// }
// func keyValueToJSONIpv6Routes(v interface{}) (jsonStr string, err error) {
// 	d, err := interfaceToIpv6Routes(v)
// 	if err != nil {
// 		return
// 	}
// 	b, err := json.Marshal(d)
// 	jsonStr = string(b)
// 	return
// }
