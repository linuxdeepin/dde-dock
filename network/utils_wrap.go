package main

import (
	"dlib/dbus"
	"encoding/json"
	"fmt"
	"strconv"
)

// define key type
type ktype uint32

const (
	ktypeUnknown ktype = iota
	ktypeString
	ktypeByte // for byte and gchar type
	ktypeInt32
	ktypeUint32
	ktypeUint64
	ktypeBoolean
	ktypeArrayByte
	ktypeArrayString      // json
	ktypeArrayUint32      // json
	ktypeArrayArrayByte   // json, array of array of byte
	ktypeArrayArrayUint32 // json, array of array of uint32
	ktypeDictStringString // json, dict of (string::string)
	ktypeIpv6Addresses    // json, array of (byte array, uint32, byte array)
	ktypeIpv6Routes       // json, array of (byte array, uint32, byte array, uint32)
)

type ktypeDescription struct {
	t    uint32
	desc string
}

var ktypeDescriptions = []ktypeDescription{
	// TODO
	{uint32(ktypeUnknown), "Unknown"},
	{uint32(ktypeString), "String"},
	{uint32(ktypeByte), "Byte"},
	{uint32(ktypeInt32), "Int32"},
	{uint32(ktypeUint32), "Uint32"},
	{uint32(ktypeUint64), "Uint64"},
	{uint32(ktypeBoolean), "Boolean"},
	{uint32(ktypeArrayByte), "ArrayByte"},
	{uint32(ktypeArrayString), "ArrayString, encode by json"},
	{uint32(ktypeArrayUint32), "ArrayUint32, encode by json"},
	{uint32(ktypeArrayArrayByte), "ArrayArrayByte, array of array of byte, encode by json"},
	{uint32(ktypeArrayArrayUint32), "ArrayArrayUint32, array of array of uint32, encode by json"},
	{uint32(ktypeDictStringString), "DictStringString, dict of (string::string), encode by json"},
	{uint32(ktypeIpv6Addresses), "Ipv6Addresses, array of (byte array, uint32, byte array), encode by json"},
	{uint32(ktypeIpv6Routes), "Ipv6Routes, array of (byte array, uint32, byte array, uint32), encode by json"},
}

// Ipv6Addresses is an array of (byte array, uint32, byte array)
type Ipv6Addresses []struct {
	Address []byte
	Prefix  uint32
	Gateway []byte
}

// Ipv6Routes is an array of (byte array, uint32, byte array, uint32)
type Ipv6Routes []struct {
	Address []byte
	Prefix  uint32
	NextHop []byte
	Metric  uint32
}

// TODO
// string -> realdata -> dbus.Variant
func wrapVariant(s string, t ktype) (v dbus.Variant) {
	var err error
	switch t {
	case ktypeString:
		v, err = wrapVariantString(s)
	case ktypeByte:
		v, err = wrapVariantByte(s)
	case ktypeInt32:
		v, err = wrapVariantInt32(s)
	case ktypeUint32:
		v, err = wrapVariantUint32(s)
	case ktypeUint64:
		v, err = wrapVariantUint64(s)
	case ktypeBoolean:
		v, err = wrapVariantBoolean(s)
	case ktypeArrayString:
		v, err = wrapVariantArrayString(s)
	case ktypeArrayByte:
		v, err = wrapVariantArrayByte(s)
	case ktypeArrayUint32:
		v, err = wrapVariantArrayUint32(s)
	case ktypeArrayArrayByte:
		v, err = wrapVariantArrayArrayByte(s)
	case ktypeArrayArrayUint32:
		v, err = wrapVariantArrayArrayUint32(s)
	case ktypeDictStringString:
		v, err = wrapVariantDictStringString(s)
	case ktypeIpv6Addresses:
		v, err = wrapVariantIpv6Addresses(s)
	case ktypeIpv6Routes:
		v, err = wrapVariantIpv6Routes(s)
	default:
		err = fmt.Errorf("invalid variant type, %s", s)
	}

	if err != nil {
		LOGGER.Error("wrapVariant() failed:", err)
	}

	return
}

func wrapVariantString(s string) (v dbus.Variant, err error) {
	v = dbus.MakeVariant(s)
	return
}

func wrapVariantByte(s string) (v dbus.Variant, err error) {
	if len(s) == 0 {
		err = fmt.Errorf("string is empty")
		return
	}
	d := s[0]
	v = dbus.MakeVariant(d)
	return
}

func wrapVariantInt32(s string) (v dbus.Variant, err error) {
	var d int32
	tmpd, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		return
	}
	d = int32(tmpd)
	v = dbus.MakeVariant(d)
	return
}

func wrapVariantUint32(s string) (v dbus.Variant, err error) {
	var d uint32
	tmpd, err := strconv.ParseUint(s, 10, 0)
	if err != nil {
		return
	}
	d = uint32(tmpd)
	v = dbus.MakeVariant(d)
	return
}

func wrapVariantUint64(s string) (v dbus.Variant, err error) {
	var d uint64
	d, err = strconv.ParseUint(s, 10, 0)
	if err != nil {
		return
	}
	v = dbus.MakeVariant(d)
	return
}

func wrapVariantBoolean(s string) (v dbus.Variant, err error) {
	var d bool
	d, err = strconv.ParseBool(s)
	if err != nil {
		return
	}
	v = dbus.MakeVariant(d)
	return
}

func wrapVariantArrayByte(s string) (v dbus.Variant, err error) {
	var d []byte
	d = []byte(s)
	v = dbus.MakeVariant(d)
	return
}

// TODO
func wrapVariantArrayString(s string) (v dbus.Variant, err error) {
	var d = make([]string, 0)
	err = json.Unmarshal([]byte(s), d)
	v = dbus.MakeVariant(d)
	return
}

// TODO
func wrapVariantArrayUint32(s string) (v dbus.Variant, err error) {
	var d = make([]uint32, 0)
	err = json.Unmarshal([]byte(s), d)
	v = dbus.MakeVariant(d)
	return
}

// TODO
func wrapVariantArrayArrayByte(s string) (v dbus.Variant, err error) {
	var d = make([][]byte, 0)
	err = json.Unmarshal([]byte(s), d)
	v = dbus.MakeVariant(d)
	return
}

// TODO
func wrapVariantArrayArrayUint32(s string) (v dbus.Variant, err error) {
	var d = make([][]uint32, 0)
	err = json.Unmarshal([]byte(s), d)
	v = dbus.MakeVariant(d)
	return
}

// TODO
func wrapVariantDictStringString(s string) (v dbus.Variant, err error) {
	var d = make(map[string]string)
	err = json.Unmarshal([]byte(s), d)
	v = dbus.MakeVariant(d)
	return
}

// TODO
func wrapVariantIpv6Addresses(s string) (v dbus.Variant, err error) {
	var d = make(Ipv6Addresses, 0)
	err = json.Unmarshal([]byte(s), d)
	v = dbus.MakeVariant(d)
	return
}

// TODO
func wrapVariantIpv6Routes(s string) (v dbus.Variant, err error) {
	var d = make(Ipv6Routes, 0)
	err = json.Unmarshal([]byte(s), d)
	v = dbus.MakeVariant(d)
	return
}
