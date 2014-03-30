package main

import (
	"dlib/dbus"
	"encoding/json"
	"fmt"
	"strconv"
)

// dbus.Variant -> realdata -> string
func unwrapVariant(v dbus.Variant, t ktype) (s string, err error) {
	switch t {
	case ktypeString:
		s, err = unwrapVariantString(v)
	case ktypeByte:
		s, err = unwrapVariantByte(v)
	case ktypeInt32:
		s, err = unwrapVariantInt32(v)
	case ktypeUint32:
		s, err = unwrapVariantUint32(v)
	case ktypeUint64:
		s, err = unwrapVariantUint64(v)
	case ktypeBoolean:
		s, err = unwrapVariantBoolean(v)
	case ktypeArrayString:
		s, err = unwrapVariantArrayString(v)
	case ktypeArrayByte:
		s, err = unwrapVariantArrayByte(v)
	case ktypeArrayUint32:
		s, err = unwrapVariantArrayUint32(v)
	case ktypeArrayArrayByte:
		s, err = unwrapVariantArrayArrayByte(v)
	case ktypeArrayArrayUint32:
		s, err = unwrapVariantArrayArrayUint32(v)
	case ktypeDictStringString:
		s, err = unwrapVariantDictStringString(v)
	case ktypeIpv6Addresses:
		s, err = unwrapVariantIpv6Addresses(v)
	case ktypeIpv6Routes:
		s, err = unwrapVariantIpv6Routes(v)
	default:
		err = fmt.Errorf("invalid key type, %v", v)
	}
	return
}

func unwrapVariantString(v dbus.Variant) (s string, err error) {
	s, ok := v.Value().(string)
	if !ok {
		err = fmt.Errorf("unwrapVariantString() failed: %v", v)
		return
	}
	return
}

func unwrapVariantByte(v dbus.Variant) (s string, err error) {
	var d byte
	d, ok := v.Value().(byte)
	if !ok {
		err = fmt.Errorf("unwrapVariantByte() failed: %v", v)
		return
	}
	s = string(d)
	return
}

func unwrapVariantInt32(v dbus.Variant) (s string, err error) {
	var d int32
	d, ok := v.Value().(int32)
	if !ok {
		err = fmt.Errorf("unwrapVariantInt32() failed: %v", v)
		return
	}
	s = strconv.FormatInt(int64(d), 10)
	return
}

func unwrapVariantUint32(v dbus.Variant) (s string, err error) {
	var d uint32
	d, ok := v.Value().(uint32)
	if !ok {
		err = fmt.Errorf("unwrapVariantUint32() failed: %v", v)
		return
	}
	s = strconv.FormatUint(uint64(d), 10)
	return
}

func unwrapVariantUint64(v dbus.Variant) (s string, err error) {
	var d uint64
	d, ok := v.Value().(uint64)
	if !ok {
		err = fmt.Errorf("unwrapVariantUint64() failed: %v", v)
		return
	}
	s = strconv.FormatUint(d, 10)
	return
}

func unwrapVariantBoolean(v dbus.Variant) (s string, err error) {
	var d bool
	d, ok := v.Value().(bool)
	if !ok {
		err = fmt.Errorf("unwrapVariantBoolean() failed: %v", v)
		return
	}
	s = strconv.FormatBool(d)
	return
}

func unwrapVariantArrayByte(v dbus.Variant) (s string, err error) {
	var d []byte
	d, ok := v.Value().([]byte)
	if !ok {
		err = fmt.Errorf("unwrapVariantArrayByte() failed: %v", v)
		return
	}
	// TODO unwrap throuh json
	// b, err := json.Marshal(d)
	// s = string(b)
	s = string(d)
	return
}

func unwrapVariantArrayString(v dbus.Variant) (s string, err error) {
	var d []string
	d, ok := v.Value().([]string)
	if !ok {
		err = fmt.Errorf("unwrapVariantArrayString() failed: %v", v)
		return
	}
	b, err := json.Marshal(d)
	s = string(b)
	return
}

func unwrapVariantArrayUint32(v dbus.Variant) (s string, err error) {
	var d []uint32
	d, ok := v.Value().([]uint32)
	if !ok {
		err = fmt.Errorf("unwrapVariantArrayUint32() failed: %v", v)
		return
	}
	b, err := json.Marshal(d)
	s = string(b)
	return
}

func unwrapVariantArrayArrayByte(v dbus.Variant) (s string, err error) {
	var d [][]byte
	d, ok := v.Value().([][]byte)
	if !ok {
		err = fmt.Errorf("unwrapVariantArrayArrayByte() failed: %v", v)
		return
	}
	b, err := json.Marshal(d)
	s = string(b)
	return
}

func unwrapVariantArrayArrayUint32(v dbus.Variant) (s string, err error) {
	var d [][]uint32
	d, ok := v.Value().([][]uint32)
	if !ok {
		err = fmt.Errorf("unwrapVariantArrayArrayUint32() failed: %v", v)
		return
	}
	b, err := json.Marshal(d)
	s = string(b)
	return
}

func unwrapVariantDictStringString(v dbus.Variant) (s string, err error) {
	var d map[string]string
	d, ok := v.Value().(map[string]string)
	if !ok {
		err = fmt.Errorf("unwrapVariantDictStringString() failed: %v", v)
		return
	}
	b, err := json.Marshal(d)
	s = string(b)
	return
}

func unwrapVariantIpv6Addresses(v dbus.Variant) (s string, err error) {
	var d Ipv6Addresses
	d, ok := v.Value().(Ipv6Addresses)
	if !ok {
		err = fmt.Errorf("unwrapVariantIpv6Addresses() failed: %v", v)
		return
	}
	b, err := json.Marshal(d)
	s = string(b)
	return
}

func unwrapVariantIpv6Routes(v dbus.Variant) (s string, err error) {
	var d Ipv6Routes
	d, ok := v.Value().(Ipv6Routes)
	if !ok {
		err = fmt.Errorf("unwrapVariantIpv6Routes() failed: %v", v)
		return
	}
	b, err := json.Marshal(d)
	s = string(b)
	return
}
