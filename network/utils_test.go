package main

import (
	"fmt"
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
	testConnectionIdJSON   = `"idname"`
	testConnectionUuidJSON = `"8e2f9aa2-42b8-47d5-b040-ae82c53fa1f2"`
	testConnectionTypeJSON = `"802-3-ethernet"`
)

const (
	testJSONKtypeString           = `"test string"`
	testJSONKtypeByte             = `97` // character 'a'
	testJSONKtypeInt32            = `-32`
	testJSONKtypeUint32           = `32`
	testJSONKtypeUint64           = `64`
	testJSONKtypeBoolean          = `true`
	testJSONKtypeArrayByte        = `"YXJyYXkgYnl0ZQ=="` // characters "array byte"
	testJSONKtypeArrayString      = `["str1","str2"]`
	testJSONKtypeArrayUint32      = `[32,32]`
	testJSONKtypeArrayArrayByte   = `["YXJyYXkgYnl0ZQ==","YXJyYXkgYnl0ZQ=="]`
	testJSONKtypeArrayArrayUint32 = `[[32,32],[32,32]]`
	testJSONKtypeDictStringString = `{"key1":"value1","key2":"value2"}`
	testJSONKtypeIpv6Addresses    = `[{"Address":"/oAAAAAAAAACImj//g9NCQ==","Prefix":32,"Gateway":"/oAAAAAAAAACImj//g9NCQ=="}]`
	testJSONKtypeIpv6Routes       = `[{"Address":"/oAAAAAAAAACImj//g9NCQ==","Prefix":32,"NextHop":"/oAAAAAAAAACImj//g9NCQ==","Metric":32}]` // TODO

	// key value wrapper
	testJSONKtypeWrapperString        = `"test wrapper string"`
	testJSONKtypeWrapperMacAddress    = `"00:12:34:56:78:AB"`
	testJSONKtypeWrapperIpv4Dns       = `["192.168.1.1","192.168.1.2"]`
	testJSONKtypeWrapperIpv4Addresses = `[{"Address":"192.168.1.100","Mask":"255.255.255.0","Gateway":"192.168.1.1"},{"Address":"192.168.1.150","Mask":"128.0.0.0","Gateway":"192.168.1.1"}]`
	testJSONKtypeWrapperIpv4Routes    = `[{"Address":"192.168.1.100","Mask":"255.255.192.0","NextHop":"192.168.1.1","Metric":100}]`
	testJSONKtypeWrapperIpv6Dns       = `["1111:2222:3333:4444:5555:6666;:aaaa:ffff"]`
	testJSONKtypeWrapperIpv6Addresses = `[{"Address":["1111:2222:3333:4444:5555:6666;:aaaa:ffff"],"Prefix":64,"Gateway":"1111:2222:3333:4444:5555:6666;:aaaa:1111"}]`
	testJSONKtypeWrapperIpv6Routes    = `[{"Address":["1111:2222:3333:4444:5555:6666;:aaaa:ffff"],"Prefix":64,"Gateway":"1111:2222:3333:4444:5555:6666;:aaaa:1111","Metric":32}]`
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

	// TODO more tests
	// set empty string to delete key
}

func (*Utils) TestConnectionDataDefaultValue(c *C) {
	// TODO
	// data := make(_ConnectionData)
}

func (*Utils) TestConvertMacAddressToString(c *C) {
	tests := []struct {
		test   []byte
		result string
	}{
		{[]byte{0, 0, 0, 0, 0, 0}, "00:00:00:00:00:00"},
		{[]byte{0, 18, 52, 86, 120, 171}, "00:12:34:56:78:AB"},
	}
	for _, t := range tests {
		c.Check(t.result, Equals, convertMacAddressToString(t.test))
	}
}

func (*Utils) TestConvertMacAddressToArrayByte(c *C) {
	tests := []struct {
		test   string
		result []byte
	}{
		{"00:00:00:00:00:00", []byte{0, 0, 0, 0, 0, 0}},
		{"00:12:34:56:78:AB", []byte{0, 18, 52, 86, 120, 171}},
	}
	for _, t := range tests {
		c.Check(t.result, DeepEquals, convertMacAddressToArrayByte(t.test))
	}
}

func (*Utils) TestConvertIpv4AddressToString(c *C) {
	tests := []struct {
		test   uint32
		result string
	}{
		{0, "0.0.0.0"},
		{0x0101a8c0, "192.168.1.1"},
	}
	for _, t := range tests {
		c.Check(t.result, Equals, convertIpv4AddressToString(t.test))
	}
}

func (*Utils) TestConvertIpv4AddressToUint32(c *C) {
	tests := []struct {
		test   string
		result uint32
	}{
		{"0.0.0.0", 0},
		{"192.168.1.1", 0x0101a8c0},
	}
	for _, t := range tests {
		c.Check(t.result, Equals, convertIpv4AddressToUint32(t.test))
	}
}

func (*Utils) TestConvertIpv4PrefixToNetMask(c *C) {
	tests := []struct {
		test   uint32
		result string
	}{
		{0, "0.0.0.0"},
		{1, "128.0.0.0"},
		{2, "192.0.0.0"},
		{3, "224.0.0.0"},
		{4, "240.0.0.0"},
		{5, "248.0.0.0"},
		{6, "252.0.0.0"},
		{7, "254.0.0.0"},
		{8, "255.0.0.0"},
		{9, "255.128.0.0"},
		{10, "255.192.0.0"},
		{11, "255.224.0.0"},
		{12, "255.240.0.0"},
		{13, "255.248.0.0"},
		{14, "255.252.0.0"},
		{15, "255.254.0.0"},
		{16, "255.255.0.0"},
		{17, "255.255.128.0"},
		{18, "255.255.192.0"},
		{19, "255.255.224.0"},
		{20, "255.255.240.0"},
		{21, "255.255.248.0"},
		{22, "255.255.252.0"},
		{23, "255.255.254.0"},
		{24, "255.255.255.0"},
		{25, "255.255.255.128"},
		{26, "255.255.255.192"},
		{27, "255.255.255.224"},
		{28, "255.255.255.240"},
		{29, "255.255.255.248"},
		{30, "255.255.255.252"},
		{31, "255.255.255.254"},
		{32, "255.255.255.255"},
	}
	for _, t := range tests {
		c.Check(t.result, Equals, convertIpv4PrefixToNetMask(t.test))
	}
}

func (*Utils) TestConvertIpv4NetMaskToPrefix(c *C) {
	tests := []struct {
		test   string
		result uint32
	}{
		{"0.0.0.0", 0},
		{"128.0.0.0", 1},
		{"192.0.0.0", 2},
		{"224.0.0.0", 3},
		{"240.0.0.0", 4},
		{"248.0.0.0", 5},
		{"252.0.0.0", 6},
		{"254.0.0.0", 7},
		{"255.0.0.0", 8},
		{"255.128.0.0", 9},
		{"255.192.0.0", 10},
		{"255.224.0.0", 11},
		{"255.240.0.0", 12},
		{"255.248.0.0", 13},
		{"255.252.0.0", 14},
		{"255.254.0.0", 15},
		{"255.255.0.0", 16},
		{"255.255.128.0", 17},
		{"255.255.192.0", 18},
		{"255.255.224.0", 19},
		{"255.255.240.0", 20},
		{"255.255.248.0", 21},
		{"255.255.252.0", 22},
		{"255.255.254.0", 23},
		{"255.255.255.0", 24},
		{"255.255.255.128", 25},
		{"255.255.255.192", 26},
		{"255.255.255.224", 27},
		{"255.255.255.240", 28},
		{"255.255.255.248", 29},
		{"255.255.255.252", 30},
		{"255.255.255.254", 31},
		{"255.255.255.255", 32},
	}
	for _, t := range tests {
		c.Check(t.result, Equals, convertIpv4NetMaskToPrefix(t.test))
	}
}

func (*Utils) TestReverseOrderUint32(c *C) {
	tests := []struct {
		test   uint32
		result uint32
	}{
		{0xaabbccdd, 0xddccbbaa},
		{0x12345678, 0x78563412},
	}
	for _, t := range tests {
		c.Check(t.result, Equals, reverseOrderUint32(t.test))
	}
}

// TODO
// func (*Utils) TestConvertIpv6AddressToString(c *C) {
// }

//func (*Utils) formatIpv6AddressToArrayByte(c *C) {}

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

	// key value wrapper
	v, _ = jsonToKeyValue(testJSONKtypeWrapperString, ktypeWrapperString)
	s, _ = keyValueToJSON(v, ktypeWrapperString)
	c.Check(s, Equals, testJSONKtypeWrapperString)

	v, _ = jsonToKeyValue(testJSONKtypeWrapperMacAddress, ktypeWrapperMacAddress)
	fmt.Println("mac address:", v)
	s, _ = keyValueToJSON(v, ktypeWrapperMacAddress)
	c.Check(s, Equals, testJSONKtypeWrapperMacAddress)

	v, _ = jsonToKeyValue(testJSONKtypeWrapperIpv4Dns, ktypeWrapperIpv4Dns)
	fmt.Printf("ipv4 dns: %x\n", v)
	s, _ = keyValueToJSON(v, ktypeWrapperIpv4Dns)
	c.Check(s, Equals, testJSONKtypeWrapperIpv4Dns)

	v, _ = jsonToKeyValue(testJSONKtypeWrapperIpv4Addresses, ktypeWrapperIpv4Addresses)
	fmt.Println("ipv4 addresses:", v)
	s, _ = keyValueToJSON(v, ktypeWrapperIpv4Addresses)
	c.Check(s, Equals, testJSONKtypeWrapperIpv4Addresses)

	v, _ = jsonToKeyValue(testJSONKtypeWrapperIpv4Routes, ktypeWrapperIpv4Routes)
	fmt.Println("ipv4 routes:", v)
	s, _ = keyValueToJSON(v, ktypeWrapperIpv4Routes)
	c.Check(s, Equals, testJSONKtypeWrapperIpv4Routes)

	// v, _ = jsonToKeyValue(testJSONKtypeWrapperIpv6Dns, ktypeWrapperIpv6Dns)
	// s, _ = keyValueToJSON(v, ktypeWrapperIpv6Dns)
	// c.Check(s, Equals, testJSONKtypeWrapperIpv6Dns)

	// v, _ = jsonToKeyValue(testJSONKtypeWrapperIpv6Addresses, ktypeWrapperIpv6Addresses)
	// s, _ = keyValueToJSON(v, ktypeWrapperIpv6Addresses)
	// c.Check(s, Equals, testJSONKtypeWrapperIpv6Addresses)

	// v, _ = jsonToKeyValue(testJSONKtypeWrapperIpv6Routes, ktypeWrapperIpv6Routes)
	// s, _ = keyValueToJSON(v, ktypeWrapperIpv6Routes)
	// c.Check(s, Equals, testJSONKtypeWrapperIpv6Routes)
}
