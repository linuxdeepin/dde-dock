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
	testConnectionId       = "idname"
	testConnectionUuid     = "8e2f9aa2-42b8-47d5-b040-ae82c53fa1f2"
	testConnectionType     = "802-3-ethernet"
	testConnectionIdJSON   = "idname"
	testConnectionUuidJSON = "8e2f9aa2-42b8-47d5-b040-ae82c53fa1f2"
	testConnectionTypeJSON = "802-3-ethernet"
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

func (*Utils) TestGetSetConnectionDataJSON(c *C) {
	data := make(_ConnectionData)
	addConnectionDataField(data, fieldConnection)
	setSettingConnectionIdJSON(data, testConnectionIdJSON)
	setSettingConnectionUuidJSON(data, testConnectionUuidJSON)
	setSettingConnectionTypeJSON(data, testConnectionTypeJSON)

	c.Check(getSettingConnectionIdJSON(data), Equals, testConnectionIdJSON)
	c.Check(getSettingConnectionUuidJSON(data), Equals, testConnectionUuidJSON)
	c.Check(getSettingConnectionTypeJSON(data), Equals, testConnectionTypeJSON)
}

func (*Utils) TestJSONWrapper(c *C) {
	var v interface{}
	var s string

	v, _ = jsonToKeyValue(testKtypeString, ktypeString)
	s, _ = keyValueToJSON(v, ktypeString)
	c.Check(s, Equals, testKtypeString)

	v, _ = jsonToKeyValue(testKtypeByte, ktypeByte)
	s, _ = keyValueToJSON(v, ktypeByte)
	c.Check(s, Equals, testKtypeByte)

	v, _ = jsonToKeyValue(testKtypeInt32, ktypeInt32)
	s, _ = keyValueToJSON(v, ktypeInt32)
	c.Check(s, Equals, testKtypeInt32)

	v, _ = jsonToKeyValue(testKtypeUint32, ktypeUint32)
	s, _ = keyValueToJSON(v, ktypeUint32)
	c.Check(s, Equals, testKtypeUint32)

	v, _ = jsonToKeyValue(testKtypeUint64, ktypeUint64)
	s, _ = keyValueToJSON(v, ktypeUint64)
	c.Check(s, Equals, testKtypeUint64)

	v, _ = jsonToKeyValue(testKtypeBoolean, ktypeBoolean)
	s, _ = keyValueToJSON(v, ktypeBoolean)
	c.Check(s, Equals, testKtypeBoolean)

	v, _ = jsonToKeyValue(testKtypeArrayByte, ktypeArrayByte)
	s, _ = keyValueToJSON(v, ktypeArrayByte)
	c.Check(s, Equals, testKtypeArrayByte)

	v, _ = jsonToKeyValue(testKtypeArrayString, ktypeArrayString)
	s, _ = keyValueToJSON(v, ktypeArrayString)
	c.Check(s, Equals, testKtypeArrayString)

	v, _ = jsonToKeyValue(testKtypeArrayUint32, ktypeArrayUint32)
	s, _ = keyValueToJSON(v, ktypeArrayUint32)
	c.Check(s, Equals, testKtypeArrayUint32)

	v, _ = jsonToKeyValue(testKtypeArrayArrayByte, ktypeArrayArrayByte)
	s, _ = keyValueToJSON(v, ktypeArrayArrayByte)
	c.Check(s, Equals, testKtypeArrayArrayByte)

	v, _ = jsonToKeyValue(testKtypeArrayArrayUint32, ktypeArrayArrayUint32)
	s, _ = keyValueToJSON(v, ktypeArrayArrayUint32)
	c.Check(s, Equals, testKtypeArrayArrayUint32)

	v, _ = jsonToKeyValue(testKtypeDictStringString, ktypeDictStringString)
	s, _ = keyValueToJSON(v, ktypeDictStringString)
	c.Check(s, Equals, testKtypeDictStringString)

	v, _ = jsonToKeyValue(testKtypeIpv6Addresses, ktypeIpv6Addresses)
	s, _ = keyValueToJSON(v, ktypeIpv6Addresses)
	c.Check(s, Equals, testKtypeIpv6Addresses)

	v, _ = jsonToKeyValue(testKtypeIpv6Routes, ktypeIpv6Routes)
	s, _ = keyValueToJSON(v, ktypeIpv6Routes)
	c.Check(s, Equals, testKtypeIpv6Routes)
}
