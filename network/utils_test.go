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
	testJSONKtypeWrapperIpv6Dns       = `["1234:2345:3456:4444:5555:6666:AAAA:FFFF"]`
	testJSONKtypeWrapperIpv6Addresses = `[{"Address":"1111:2222:3333:4444:5555:6666:AAAA:FFFF","Prefix":64,"Gateway":"1111:2222:3333:4444:5555:6666:AAAA:1111"}]`
	testJSONKtypeWrapperIpv6Routes    = `[{"Address":"1111:2222:3333:4444:5555:6666:AAAA:FFFF","Prefix":64,"NextHop":"1111:2222:3333:4444:5555:6666:AAAA:1111","Metric":32}]`
)

func (*Utils) TestGetSetConnectionData(c *C) {
	data := make(connectionData)

	addSettingField(data, fieldConnection)
	setSettingConnectionId(data, testConnectionId)
	setSettingConnectionUuid(data, testConnectionUuid)
	setSettingConnectionType(data, testConnectionType)

	c.Check(getSettingConnectionId(data), Equals, testConnectionId)
	c.Check(getSettingConnectionUuid(data), Equals, testConnectionUuid)
	c.Check(getSettingConnectionType(data), Equals, testConnectionType)
}

func (*Utils) TestGetSetConnectionDataJSON(c *C) {
	data := make(connectionData)
	addSettingField(data, fieldConnection)
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
	data := make(connectionData)
	var defaultValueJSON string
	var setValueJSON string
	addSettingField(data, fieldConnection)
	addSettingField(data, fieldWired)
	addSettingField(data, field8021x)
	addSettingField(data, fieldIpv4)
	addSettingField(data, fieldIpv6)

	// ktypeBoolean
	defaultValueJSON = `true`
	c.Check(getSettingConnectionAutoconnectJSON(data), Equals, defaultValueJSON)
	setSettingConnectionAutoconnectJSON(data, defaultValueJSON)
	c.Check(isSettingKeyExists(data, fieldConnection, NM_SETTING_CONNECTION_AUTOCONNECT), Equals, true)

	// ktypeArrayByte
	defaultValueJSON = `""`
	setValueJSON = `""`
	c.Check(getSetting8021xPasswordRawJSON(data), Equals, defaultValueJSON)
	setSetting8021xPasswordRawJSON(data, setValueJSON)
	c.Check(isSettingKeyExists(data, field8021x, NM_SETTING_802_1X_PASSWORD_RAW), Equals, false)

	// ktypeString
	defaultValueJSON = `""`
	setValueJSON = `""`
	c.Check(getSettingConnectionIdJSON(data), Equals, defaultValueJSON)
	setSettingConnectionIdJSON(data, setValueJSON)
	c.Check(isSettingKeyExists(data, fieldConnection, NM_SETTING_CONNECTION_ID), Equals, false)

	// ktypeWrapperMacAddress
	defaultValueJSON = `""`
	setValueJSON = `""`
	c.Check(getSettingWiredMacAddressJSON(data), Equals, defaultValueJSON)
	setSettingWiredMacAddressJSON(data, setValueJSON)
	c.Check(isSettingKeyExists(data, fieldWired, NM_SETTING_WIRED_MAC_ADDRESS), Equals, false)

	// ktypeWrapperString
	defaultValueJSON = `""`
	setValueJSON = `""`
	c.Check(getSetting8021xCaCertJSON(data), Equals, defaultValueJSON)
	setSetting8021xCaCertJSON(data, setValueJSON)
	c.Check(isSettingKeyExists(data, field8021x, NM_SETTING_802_1X_CA_CERT), Equals, false)

	// ktypeWrapperIpv4Dns
	defaultValueJSON = `null`
	setValueJSON = `[""]`
	c.Check(getSettingIp4ConfigDnsJSON(data), Equals, defaultValueJSON)
	setSettingIp4ConfigDnsJSON(data, setValueJSON)
	c.Check(isSettingKeyExists(data, fieldIpv4, NM_SETTING_IP4_CONFIG_DNS), Equals, false)

	// ktypeWrapperIpv4Addresses
	defaultValueJSON = `null`
	setValueJSON = `[{"Address":"","Mask":"","Gateway":""}]`
	c.Check(getSettingIp4ConfigAddressesJSON(data), Equals, defaultValueJSON)
	setSettingIp4ConfigAddressesJSON(data, setValueJSON)
	c.Check(isSettingKeyExists(data, fieldIpv4, NM_SETTING_IP4_CONFIG_ADDRESSES), Equals, false)

	// ktypeWrapperIpv4Routes
	defaultValueJSON = `null`
	setValueJSON = `[{"Address":"","Mask":"","NextHop":"","Metric":0}]`
	c.Check(getSettingIp4ConfigRoutesJSON(data), Equals, defaultValueJSON)
	setSettingIp4ConfigRoutesJSON(data, setValueJSON)
	c.Check(isSettingKeyExists(data, fieldIpv4, NM_SETTING_IP4_CONFIG_ROUTES), Equals, false)

	// ktypeWrapperIpv6Dns
	defaultValueJSON = `null`
	setValueJSON = `[""]`
	c.Check(getSettingIp6ConfigDnsJSON(data), Equals, defaultValueJSON)
	setSettingIp6ConfigDnsJSON(data, setValueJSON)
	c.Check(isSettingIp6ConfigDnsExists(data), Equals, false)

	// ktypeWrapperIpv6Addresses
	defaultValueJSON = `null`
	setValueJSON = `[{"Address":"","Prefix":0,"Gateway":""}]`
	c.Check(getSettingIp6ConfigAddressesJSON(data), Equals, defaultValueJSON)
	setSettingIp6ConfigAddressesJSON(data, setValueJSON)
	c.Check(isSettingIp6ConfigAddressesExists(data), Equals, false)

	// ktypeWrapperIpv6
	defaultValueJSON = `null`
	setValueJSON = `[{"Address":"","Prefix":0,"NextHop":"","Metric":0}]`
	c.Check(getSettingIp6ConfigRoutesJSON(data), Equals, defaultValueJSON)
	setSettingIp6ConfigRoutesJSON(data, setValueJSON)
	c.Check(isSettingIp6ConfigRoutesExists(data), Equals, false)
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

	// test error mask address
	_, err := convertIpv4NetMaskToPrefixCheck("255.255.255.250")
	c.Check(err, NotNil)
	_, err = convertIpv4NetMaskToPrefixCheck("255.255.100.2")
	c.Check(err, NotNil)
	_, err = convertIpv4NetMaskToPrefixCheck("255.100.0.0")
	c.Check(err, NotNil)
	_, err = convertIpv4NetMaskToPrefixCheck("191.0.0.0")
	c.Check(err, NotNil)
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

func (*Utils) TestConvertIpv6AddressToString(c *C) {
	tests := []struct {
		test   []byte
		result string
	}{
		{[]byte{0x12, 0x34, 0x23, 0x45, 0x34, 0x56, 0x44, 0x44, 0x55, 0x55, 0x66, 0x66, 0xaa, 0xaa, 0xff, 0xff}, "1234:2345:3456:4444:5555:6666:AAAA:FFFF"},
		{[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, "0000:0000:0000:0000:0000:0000:0000:0000"},
	}
	for _, t := range tests {
		c.Check(t.result, Equals, convertIpv6AddressToString(t.test)) // TODO
	}
}

func (*Utils) TestConvertIpv6AddressToArrayByte(c *C) {
	tests := []struct {
		test   string
		result []byte
	}{
		{"1234:2345:3456:4444:5555:6666:AAAA:FFFF", []byte{0x12, 0x34, 0x23, 0x45, 0x34, 0x56, 0x44, 0x44, 0x55, 0x55, 0x66, 0x66, 0xaa, 0xaa, 0xff, 0xff}},
		{"0000:0000:0000:0000:0000:0000:0000:0000", []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
	}
	for _, t := range tests {
		c.Check(t.result, DeepEquals, convertIpv6AddressToArrayByte(t.test))
	}

	// check error ipv6 format
	_, err := convertIpv6AddressToArrayByteCheck("-1234:2345:3456:4444:5555:6666:aaAA:ffFF")
	c.Check(err, NotNil)
	_, err = convertIpv6AddressToArrayByteCheck("1234:2345:3456:4444:5555:6666:aaAA:ffFh")
	c.Check(err, NotNil)
}

func (*Utils) TestExpandIpv6Address(c *C) {
	tests := []struct {
		test   string
		result string
	}{
		{"1234:2345:3456:4444:5555:6666:AAAA:FFFF", "1234:2345:3456:4444:5555:6666:AAAA:FFFF"},
		{"0000:0000:0000:0000:0000:0000:0000:0000", "0000:0000:0000:0000:0000:0000:0000:0000"},
		{"0::0", "0000:0000:0000:0000:0000:0000:0000:0000"},
		{"2001:DB8:2de::e13", "2001:0DB8:02de:0000:0000:0000:0000:0e13"},
		{"::ffff:874B:2B34", "0000:0000:0000:0000:0000:ffff:874B:2B34"},
	}
	for _, t := range tests {
		r, _ := expandIpv6Address(t.test)
		c.Check(t.result, Equals, r)
	}

	// check error ipv6 format
	_, err := expandIpv6Address("2001::25de::cade")
	c.Check(err, NotNil)
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

	v, _ = jsonToKeyValue(testJSONKtypeWrapperIpv6Dns, ktypeWrapperIpv6Dns)
	fmt.Println("ipv6 dns:", v)
	s, _ = keyValueToJSON(v, ktypeWrapperIpv6Dns)
	c.Check(s, Equals, testJSONKtypeWrapperIpv6Dns)

	v, _ = jsonToKeyValue(testJSONKtypeWrapperIpv6Addresses, ktypeWrapperIpv6Addresses)
	fmt.Println("ipv6 address:", v)
	s, _ = keyValueToJSON(v, ktypeWrapperIpv6Addresses)
	c.Check(s, Equals, testJSONKtypeWrapperIpv6Addresses)

	v, _ = jsonToKeyValue(testJSONKtypeWrapperIpv6Routes, ktypeWrapperIpv6Routes)
	fmt.Println("ipv6 route:", v)
	s, _ = keyValueToJSON(v, ktypeWrapperIpv6Routes)
	c.Check(s, Equals, testJSONKtypeWrapperIpv6Routes)
}

func (*Utils) TestGetterAndSetterForVirtualKey(c *C) {
	// TODO
}

func (*Utils) TestToUriPath(c *C) {
	tests := []struct {
		test   string
		result string
	}{
		{"/the/path", "file:///the/path"},
		{"file:///the/path", "file:///the/path"},
	}
	for _, t := range tests {
		c.Check(t.result, Equals, toUriPath(t.test))
	}
}
func (*Utils) TestToLocalPath(c *C) {
	tests := []struct {
		test   string
		result string
	}{
		{"/the/path", "/the/path"},
		{"file:///the/path", "/the/path"},
	}
	for _, t := range tests {
		c.Check(t.result, Equals, toLocalPath(t.test))
	}
}

func (*Utils) TestStrToByteArrayPath(c *C) {
	tests := []struct {
		test   string
		result []byte
	}{
		{"/the/path", []byte{0x2f, 0x74, 0x68, 0x65, 0x2f, 0x70, 0x61, 0x74, 0x68, 0x0}},
	}
	for _, t := range tests {
		c.Check(t.result, DeepEquals, strToByteArrayPath(t.test))
	}
}
func (*Utils) TestbyteArrayToStrPath(c *C) {
	tests := []struct {
		test   []byte
		result string
	}{
		{[]byte{0x2f, 0x74, 0x68, 0x65, 0x2f, 0x70, 0x61, 0x74, 0x68, 0x0}, "/the/path"},
		{[]byte{0x0}, ""},
		{[]byte{}, ""},
	}
	for _, t := range tests {
		c.Check(t.result, Equals, byteArrayToStrPath(t.test))
	}
}
