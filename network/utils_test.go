package main

import (
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

func (*Utils) TestGetSetConnectionDataJSON(c *C) {
	data := make(_ConnectionData)
	addConnectionDataField(data, fieldConnection)
	setSettingConnectionIdJSON(data, testConnectionId)
	setSettingConnectionUuidJSON(data, testConnectionUuid)
	setSettingConnectionTypeJSON(data, testConnectionType)

	c.Check(getSettingConnectionIdJSON(data), Equals, testConnectionId)
	c.Check(getSettingConnectionUuidJSON(data), Equals, testConnectionUuid)
	c.Check(getSettingConnectionTypeJSON(data), Equals, testConnectionType)
}

func (*Utils) TestJSONWrapper(c *C) {
	var v interface{}
	var s string

	v, _ = jsonToInterface(testKtypeString, ktypeString)
	s, _ = interfaceToJSON(v, ktypeString)
	c.Check(s, Equals, testKtypeString)

	v, _ = jsonToInterface(testKtypeByte, ktypeByte)
	s, _ = interfaceToJSON(v, ktypeByte)
	c.Check(s, Equals, testKtypeByte)

	v, _ = jsonToInterface(testKtypeInt32, ktypeInt32)
	s, _ = interfaceToJSON(v, ktypeInt32)
	c.Check(s, Equals, testKtypeInt32)

	v, _ = jsonToInterface(testKtypeUint32, ktypeUint32)
	s, _ = interfaceToJSON(v, ktypeUint32)
	c.Check(s, Equals, testKtypeUint32)

	v, _ = jsonToInterface(testKtypeUint64, ktypeUint64)
	s, _ = interfaceToJSON(v, ktypeUint64)
	c.Check(s, Equals, testKtypeUint64)

	v, _ = jsonToInterface(testKtypeBoolean, ktypeBoolean)
	s, _ = interfaceToJSON(v, ktypeBoolean)
	c.Check(s, Equals, testKtypeBoolean)

	v, _ = jsonToInterface(testKtypeArrayByte, ktypeArrayByte)
	s, _ = interfaceToJSON(v, ktypeArrayByte)
	c.Check(s, Equals, testKtypeArrayByte)

	v, _ = jsonToInterface(testKtypeArrayString, ktypeArrayString)
	s, _ = interfaceToJSON(v, ktypeArrayString)
	c.Check(s, Equals, testKtypeArrayString)

	v, _ = jsonToInterface(testKtypeArrayUint32, ktypeArrayUint32)
	s, _ = interfaceToJSON(v, ktypeArrayUint32)
	c.Check(s, Equals, testKtypeArrayUint32)

	v, _ = jsonToInterface(testKtypeArrayArrayByte, ktypeArrayArrayByte)
	s, _ = interfaceToJSON(v, ktypeArrayArrayByte)
	c.Check(s, Equals, testKtypeArrayArrayByte)

	v, _ = jsonToInterface(testKtypeArrayArrayUint32, ktypeArrayArrayUint32)
	s, _ = interfaceToJSON(v, ktypeArrayArrayUint32)
	c.Check(s, Equals, testKtypeArrayArrayUint32)

	v, _ = jsonToInterface(testKtypeDictStringString, ktypeDictStringString)
	s, _ = interfaceToJSON(v, ktypeDictStringString)
	c.Check(s, Equals, testKtypeDictStringString)

	v, _ = jsonToInterface(testKtypeIpv6Addresses, ktypeIpv6Addresses)
	s, _ = interfaceToJSON(v, ktypeIpv6Addresses)
	c.Check(s, Equals, testKtypeIpv6Addresses)

	v, _ = jsonToInterface(testKtypeIpv6Routes, ktypeIpv6Routes)
	s, _ = interfaceToJSON(v, ktypeIpv6Routes)
	c.Check(s, Equals, testKtypeIpv6Routes)
}
