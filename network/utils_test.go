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
	testJSONKtypeString  = `"test string"`
	testJSONKtypeByte    = "a"
	testJSONKtypeInt32   = `-32`
	testJSONKtypeUint32  = `32`
	testJSONKtypeUint64  = `64`
	testJSONKtypeBoolean = `true`
	// testJSONKtypeArrayByte        = `"YXJyYXkgYnl0ZQ=="` // json, "array byte"
	testJSONKtypeArrayByte        = `array byte`
	testJSONKtypeArrayString      = `["str1","str2"]`
	testJSONKtypeArrayUint32      = `[32,32]`
	testJSONKtypeArrayArrayByte   = `["YXJyYXkgYnl0ZQ==","YXJyYXkgYnl0ZQ=="]`
	testJSONKtypeArrayArrayUint32 = `[[32,32],[32,32]]`
	testJSONKtypeDictStringString = `{"key1":"value1","key2":"value2"}`
	testJSONKtypeIpv6Addresses    = `[{"Address":"/oAAAAAAAAACImj//g9NCQ==","Prefix":32,"Gateway":"/oAAAAAAAAACImj//g9NCQ=="}]`
	testJSONKtypeIpv6Routes       = `[{"Address":"/oAAAAAAAAACImj//g9NCQ==","Prefix":32,"NextHop":"/oAAAAAAAAACImj//g9NCQ==","Metric":32}]` // TODO
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

	v, _ = jsonToKeyValue(testJSONKtypeString, ktypeString)
	s, _ = keyValueToJSON(v, ktypeString)
	c.Check(s, Equals, testJSONKtypeString)

	v, _ = jsonToKeyValue(testJSONKtypeByte, ktypeByte)
	s, _ = keyValueToJSON(v, ktypeByte)
	c.Check(s, Equals, testJSONKtypeByte)

	v, _ = jsonToKeyValue(testJSONKtypeInt32, ktypeInt32)
	s, _ = keyValueToJSON(v, ktypeInt32)
	c.Check(s, Equals, testJSONKtypeInt32)

	v, _ = jsonToKeyValue(testJSONKtypeUint32, ktypeUint32)
	s, _ = keyValueToJSON(v, ktypeUint32)
	c.Check(s, Equals, testJSONKtypeUint32)

	v, _ = jsonToKeyValue(testJSONKtypeUint64, ktypeUint64)
	s, _ = keyValueToJSON(v, ktypeUint64)
	c.Check(s, Equals, testJSONKtypeUint64)

	v, _ = jsonToKeyValue(testJSONKtypeBoolean, ktypeBoolean)
	s, _ = keyValueToJSON(v, ktypeBoolean)
	c.Check(s, Equals, testJSONKtypeBoolean)

	v, _ = jsonToKeyValue(testJSONKtypeArrayByte, ktypeArrayByte)
	s, _ = keyValueToJSON(v, ktypeArrayByte)
	c.Check(s, Equals, testJSONKtypeArrayByte)

	v, _ = jsonToKeyValue(testJSONKtypeArrayString, ktypeArrayString)
	s, _ = keyValueToJSON(v, ktypeArrayString)
	c.Check(s, Equals, testJSONKtypeArrayString)

	v, _ = jsonToKeyValue(testJSONKtypeArrayUint32, ktypeArrayUint32)
	s, _ = keyValueToJSON(v, ktypeArrayUint32)
	c.Check(s, Equals, testJSONKtypeArrayUint32)

	v, _ = jsonToKeyValue(testJSONKtypeArrayArrayByte, ktypeArrayArrayByte)
	s, _ = keyValueToJSON(v, ktypeArrayArrayByte)
	c.Check(s, Equals, testJSONKtypeArrayArrayByte)

	v, _ = jsonToKeyValue(testJSONKtypeArrayArrayUint32, ktypeArrayArrayUint32)
	s, _ = keyValueToJSON(v, ktypeArrayArrayUint32)
	c.Check(s, Equals, testJSONKtypeArrayArrayUint32)

	v, _ = jsonToKeyValue(testJSONKtypeDictStringString, ktypeDictStringString)
	s, _ = keyValueToJSON(v, ktypeDictStringString)
	c.Check(s, Equals, testJSONKtypeDictStringString)

	v, _ = jsonToKeyValue(testJSONKtypeIpv6Addresses, ktypeIpv6Addresses)
	s, _ = keyValueToJSON(v, ktypeIpv6Addresses)
	c.Check(s, Equals, testJSONKtypeIpv6Addresses)

	v, _ = jsonToKeyValue(testJSONKtypeIpv6Routes, ktypeIpv6Routes)
	s, _ = keyValueToJSON(v, ktypeIpv6Routes)
	c.Check(s, Equals, testJSONKtypeIpv6Routes)
}
