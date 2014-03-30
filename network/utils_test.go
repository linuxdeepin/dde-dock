package main

import (
	"dlib/dbus"
	. "launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type Utils struct{}

func init() {
	utils := &Utils{}
	Suite(utils)
}

const (
	testConnectionId   = "idname"
	testConnectionUuid = "8e2f9aa2-42b8-47d5-b040-ae82c53fa1f2"
	testConnectionType = typeWired
)

const (
	testKtypeString  = "test string"
	testKtypeByte    = "a"
	testKtypeInt32   = "-32"
	testKtypeUint32  = "32"
	testKtypeUint64  = "64"
	testKtypeBoolean = "true"
	// testKtypeArrayByte        = `"YXJyYXkgYnl0ZQ=="` // json, "array byte"
	testKtypeArrayByte        = `array byte`
	testKtypeArrayString      = `["str1","str2"]`
	testKtypeArrayUint32      = `[32,32]`
	testKtypeArrayArrayByte   = `["YXJyYXkgYnl0ZQ==","YXJyYXkgYnl0ZQ=="]`
	testKtypeArrayArrayUint32 = `[[32,32],[32,32]]`
	testKtypeDictStringString = `{"key1":"value1","key2":"value2"}`
	testKtypeIpv6Addresses    = `[{"Address":"/oAAAAAAAAACImj//g9NCQ==","Prefix":32,"Gateway":"/oAAAAAAAAACImj//g9NCQ=="}]`
	testKtypeIpv6Routes       = `[{"Address":"/oAAAAAAAAACImj//g9NCQ==","Prefix":32,"NextHop":"/oAAAAAAAAACImj//g9NCQ==","Metric":32}]` // TODO
	// 'addresses': [([254, 128, 0, 0, 0, 0, 0, 0, 2, 34, 104, 255, 254, 15, 77, 9], 64L, [254, 128, 0, 0, 0, 0, 0, 0, 2, 34, 104, 255, 254, 15, 77, 9])]
	// 'routes': [([254, 128, 0, 0, 0, 0, 0, 0, 2, 34, 104, 255, 254, 15, 77, 9], 64L, [254, 128, 0, 0, 0, 0, 0, 0, 2, 34, 104, 255, 254, 15, 77, 9], 12L)]
)

func (*Utils) TestGetSetConnectionData(c *C) {
	data := make(_ConnectionData)
	addConnectionDataField(data, fieldConnection)
	setSettingConnectionId(data, testConnectionId)
	setSettingConnectionUuid(data, testConnectionUuid)
	setSettingConnectionType(data, testConnectionType)

	c.Check(getSettingConnectionId(data), Equals, testConnectionId)
	c.Check(getSettingConnectionUuid(data), Equals, testConnectionUuid)
	c.Check(getSettingConnectionType(data), Equals, testConnectionType)
}

func (*Utils) TestVariantWrapper(c *C) {
	var v dbus.Variant
	var s string

	v, _ = wrapVariant(testKtypeString, ktypeString)
	s, _ = unwrapVariant(v, ktypeString)
	c.Check(s, Equals, testKtypeString)

	v, _ = wrapVariant(testKtypeByte, ktypeByte)
	s, _ = unwrapVariant(v, ktypeByte)
	c.Check(s, Equals, testKtypeByte)

	v, _ = wrapVariant(testKtypeInt32, ktypeInt32)
	s, _ = unwrapVariant(v, ktypeInt32)
	c.Check(s, Equals, testKtypeInt32)

	v, _ = wrapVariant(testKtypeUint32, ktypeUint32)
	s, _ = unwrapVariant(v, ktypeUint32)
	c.Check(s, Equals, testKtypeUint32)

	v, _ = wrapVariant(testKtypeUint64, ktypeUint64)
	s, _ = unwrapVariant(v, ktypeUint64)
	c.Check(s, Equals, testKtypeUint64)

	v, _ = wrapVariant(testKtypeBoolean, ktypeBoolean)
	s, _ = unwrapVariant(v, ktypeBoolean)
	c.Check(s, Equals, testKtypeBoolean)

	v, _ = wrapVariant(testKtypeArrayByte, ktypeArrayByte)
	s, _ = unwrapVariant(v, ktypeArrayByte)
	c.Check(s, Equals, testKtypeArrayByte)

	v, _ = wrapVariant(testKtypeArrayString, ktypeArrayString)
	s, _ = unwrapVariant(v, ktypeArrayString)
	c.Check(s, Equals, testKtypeArrayString)

	v, _ = wrapVariant(testKtypeArrayUint32, ktypeArrayUint32)
	s, _ = unwrapVariant(v, ktypeArrayUint32)
	c.Check(s, Equals, testKtypeArrayUint32)

	v, _ = wrapVariant(testKtypeArrayArrayByte, ktypeArrayArrayByte)
	s, _ = unwrapVariant(v, ktypeArrayArrayByte)
	c.Check(s, Equals, testKtypeArrayArrayByte)

	v, _ = wrapVariant(testKtypeArrayArrayUint32, ktypeArrayArrayUint32)
	s, _ = unwrapVariant(v, ktypeArrayArrayUint32)
	c.Check(s, Equals, testKtypeArrayArrayUint32)

	v, _ = wrapVariant(testKtypeDictStringString, ktypeDictStringString)
	s, _ = unwrapVariant(v, ktypeDictStringString)
	c.Check(s, Equals, testKtypeDictStringString)

	v, _ = wrapVariant(testKtypeIpv6Addresses, ktypeIpv6Addresses)
	s, _ = unwrapVariant(v, ktypeIpv6Addresses)
	c.Check(s, Equals, testKtypeIpv6Addresses)

	v, _ = wrapVariant(testKtypeIpv6Routes, ktypeIpv6Routes)
	s, _ = unwrapVariant(v, ktypeIpv6Routes)
	c.Check(s, Equals, testKtypeIpv6Routes)
}
