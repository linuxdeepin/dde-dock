/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package network

import (
	C "gopkg.in/check.v1"
	"os"
	"pkg.deepin.io/dde/daemon/network/nm"
	"testing"
)

func Test(t *testing.T) { C.TestingT(t) }

type testWrapper struct{}

func init() {
	C.Suite(&testWrapper{})
}

func (*testWrapper) TestGetSetConnectionData(c *C.C) {
	testConnectionId := "idname"
	testConnectionUuid := "8e2f9aa2-42b8-47d5-b040-ae82c53fa1f2"
	testConnectionType := "802-3-ethernet"

	data := make(connectionData)

	addSetting(data, nm.NM_SETTING_CONNECTION_SETTING_NAME)
	setSettingConnectionId(data, testConnectionId)
	setSettingConnectionUuid(data, testConnectionUuid)
	setSettingConnectionType(data, testConnectionType)

	c.Check(getSettingConnectionId(data), C.Equals, testConnectionId)
	c.Check(getSettingConnectionUuid(data), C.Equals, testConnectionUuid)
	c.Check(getSettingConnectionType(data), C.Equals, testConnectionType)
}

func (*testWrapper) TestGetSetConnectionDataJSON(c *C.C) {
	testConnectionIdJSON := `"idname"`
	testConnectionUuidJSON := `"8e2f9aa2-42b8-47d5-b040-ae82c53fa1f2"`
	testConnectionTypeJSON := `"802-3-ethernet"`

	data := make(connectionData)
	addSetting(data, nm.NM_SETTING_CONNECTION_SETTING_NAME)
	setSettingConnectionIdJSON(data, testConnectionIdJSON)
	setSettingConnectionUuidJSON(data, testConnectionUuidJSON)
	setSettingConnectionTypeJSON(data, testConnectionTypeJSON)

	c.Check(getSettingConnectionIdJSON(data), C.Equals, testConnectionIdJSON)
	c.Check(getSettingConnectionUuidJSON(data), C.Equals, testConnectionUuidJSON)
	c.Check(getSettingConnectionTypeJSON(data), C.Equals, testConnectionTypeJSON)
}

func (*testWrapper) TestConnectionDataDefaultValue(c *C.C) {
	data := make(connectionData)
	var defaultValueJSON string
	var setValueJSON string
	addSetting(data, nm.NM_SETTING_CONNECTION_SETTING_NAME)
	addSetting(data, nm.NM_SETTING_WIRED_SETTING_NAME)
	addSetting(data, nm.NM_SETTING_802_1X_SETTING_NAME)
	addSetting(data, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME)
	addSetting(data, nm.NM_SETTING_IP6_CONFIG_SETTING_NAME)

	// ktypeBoolean
	defaultValueJSON = `true`
	c.Check(getSettingConnectionAutoconnectJSON(data), C.Equals, defaultValueJSON)
	setSettingConnectionAutoconnectJSON(data, defaultValueJSON)
	c.Check(isSettingKeyExists(data, nm.NM_SETTING_CONNECTION_SETTING_NAME, nm.NM_SETTING_CONNECTION_AUTOCONNECT), C.Equals, true)

	// ktypeArrayByte
	defaultValueJSON = `""`
	setValueJSON = `""`
	c.Check(getSetting8021xPasswordRawJSON(data), C.Equals, defaultValueJSON)
	setSetting8021xPasswordRawJSON(data, setValueJSON)
	c.Check(isSettingKeyExists(data, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_PASSWORD_RAW), C.Equals, false)

	// ktypeString
	defaultValueJSON = `""`
	setValueJSON = `""`
	c.Check(getSettingConnectionIdJSON(data), C.Equals, defaultValueJSON)
	setSettingConnectionIdJSON(data, setValueJSON)
	c.Check(isSettingKeyExists(data, nm.NM_SETTING_CONNECTION_SETTING_NAME, nm.NM_SETTING_CONNECTION_ID), C.Equals, false)

	// ktypeWrapperMacAddress
	defaultValueJSON = `""`
	setValueJSON = `""`
	c.Check(getSettingWiredMacAddressJSON(data), C.Equals, defaultValueJSON)
	setSettingWiredMacAddressJSON(data, setValueJSON)
	c.Check(isSettingKeyExists(data, nm.NM_SETTING_WIRED_SETTING_NAME, nm.NM_SETTING_WIRED_MAC_ADDRESS), C.Equals, false)

	// ktypeWrapperString
	defaultValueJSON = `""`
	setValueJSON = `""`
	c.Check(getSetting8021xCaCertJSON(data), C.Equals, defaultValueJSON)
	setSetting8021xCaCertJSON(data, setValueJSON)
	c.Check(isSettingKeyExists(data, nm.NM_SETTING_802_1X_SETTING_NAME, nm.NM_SETTING_802_1X_CA_CERT), C.Equals, false)

	// ktypeWrapperIpv4Dns
	defaultValueJSON = `[]`
	setValueJSON = `[""]`
	c.Check(getSettingIP4ConfigDnsJSON(data), C.Equals, defaultValueJSON)
	setSettingIP4ConfigDnsJSON(data, setValueJSON)
	c.Check(isSettingKeyExists(data, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME, nm.NM_SETTING_IP_CONFIG_DNS), C.Equals, false)

	// ktypeWrapperIpv4Addresses
	defaultValueJSON = `[[0,24,0]]`
	setValueJSON = `[{"Address":"","Mask":"","Gateway":""}]`
	c.Check(getSettingIP4ConfigAddressesJSON(data), C.Equals, defaultValueJSON)
	setSettingIP4ConfigAddressesJSON(data, setValueJSON)
	c.Check(isSettingKeyExists(data, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME, nm.NM_SETTING_IP_CONFIG_ADDRESSES), C.Equals, false)

	// ktypeWrapperIpv4Routes
	defaultValueJSON = `[]`
	setValueJSON = `[{"Address":"","Mask":"","NextHop":"","Metric":0}]`
	c.Check(getSettingIP4ConfigRoutesJSON(data), C.Equals, defaultValueJSON)
	setSettingIP4ConfigRoutesJSON(data, setValueJSON)
	c.Check(isSettingKeyExists(data, nm.NM_SETTING_IP4_CONFIG_SETTING_NAME, nm.NM_SETTING_IP_CONFIG_ROUTES), C.Equals, false)

	// ktypeWrapperIpv6Dns
	defaultValueJSON = `[]`
	setValueJSON = `[""]`
	c.Check(getSettingIP6ConfigDnsJSON(data), C.Equals, defaultValueJSON)
	setSettingIP6ConfigDnsJSON(data, setValueJSON)
	c.Check(isSettingIP6ConfigDnsExists(data), C.Equals, false)

	// ktypeWrapperIpv6Addresses
	defaultValueJSON = `[]`
	setValueJSON = `[{"Address":"","Prefix":0,"Gateway":""}]`
	c.Check(getSettingIP6ConfigAddressesJSON(data), C.Equals, defaultValueJSON)
	setSettingIP6ConfigAddressesJSON(data, setValueJSON)
	c.Check(isSettingIP6ConfigAddressesExists(data), C.Equals, false)

	// ktypeWrapperIpv6
	defaultValueJSON = `[]`
	setValueJSON = `[{"Address":"","Prefix":0,"NextHop":"","Metric":0}]`
	c.Check(getSettingIP6ConfigRoutesJSON(data), C.Equals, defaultValueJSON)
	setSettingIP6ConfigRoutesJSON(data, setValueJSON)
	c.Check(isSettingIP6ConfigRoutesExists(data), C.Equals, false)
}

func (*testWrapper) TestKeyError(c *C.C) {
	cwd, err := os.Getwd()
	if err != nil {
		c.Skip(err.Error())
	}

	var errs sectionErrors
	data := make(connectionData)
	addSetting(data, nm.NM_SETTING_802_1X_SETTING_NAME)

	// check 8021x cert file
	errs = make(sectionErrors)
	logicSetSettingVk8021xCaCert(data, cwd+"/testdata/ca.crt")
	checkSetting8021xCaCert(data, errs)
	c.Check(len(errs), C.Equals, 0)

	errs = make(sectionErrors)
	logicSetSettingVk8021xCaCert(data, cwd+"/testdata/ca.crt.notexists")
	checkSetting8021xCaCert(data, errs)
	c.Check(len(errs), C.Equals, 1)

	errs = make(sectionErrors)
	logicSetSettingVk8021xCaCert(data, "abc"+cwd+"/testdata/ca.crt")
	checkSetting8021xCaCert(data, errs)
	c.Check(len(errs), C.Equals, 1)
}

func (*testWrapper) TestConvertMacAddressToString(c *C.C) {
	tests := []struct {
		test   []byte
		result string
	}{
		{[]byte{0, 0, 0, 0, 0, 0}, "00:00:00:00:00:00"},
		{[]byte{0, 18, 52, 86, 120, 171}, "00:12:34:56:78:AB"},
	}
	for _, t := range tests {
		c.Check(t.result, C.Equals, convertMacAddressToString(t.test))
	}
}

func (*testWrapper) TestConvertMacAddressToArrayByte(c *C.C) {
	tests := []struct {
		test   string
		result []byte
	}{
		{"00:00:00:00:00:00", []byte{0, 0, 0, 0, 0, 0}},
		{"00:12:34:56:78:AB", []byte{0, 18, 52, 86, 120, 171}},
	}
	for _, t := range tests {
		c.Check(t.result, C.DeepEquals, convertMacAddressToArrayByte(t.test))
	}
}

func (*testWrapper) TestConvertIpv4AddressToString(c *C.C) {
	tests := []struct {
		test   uint32
		result string
	}{
		{0, "0.0.0.0"},
		{0x0101a8c0, "192.168.1.1"},
	}
	for _, t := range tests {
		c.Check(t.result, C.Equals, convertIpv4AddressToString(t.test))
	}
}

func (*testWrapper) TestConvertIpv4AddressToUint32(c *C.C) {
	tests := []struct {
		test   string
		result uint32
	}{
		{"0.0.0.0", 0},
		{"192.168.1.1", 0x0101a8c0},
	}
	for _, t := range tests {
		c.Check(t.result, C.Equals, convertIpv4AddressToUint32(t.test))
	}
}

func (*testWrapper) TestConvertIpv4PrefixToNetMask(c *C.C) {
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
		c.Check(t.result, C.Equals, convertIpv4PrefixToNetMask(t.test))
	}
}

func (*testWrapper) TestConvertIpv4NetMaskToPrefix(c *C.C) {
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
		c.Check(t.result, C.Equals, convertIpv4NetMaskToPrefix(t.test))
	}

	// test error mask address
	_, err := convertIpv4NetMaskToPrefixCheck("255.255.255.250")
	c.Check(err, C.NotNil)
	_, err = convertIpv4NetMaskToPrefixCheck("255.255.100.2")
	c.Check(err, C.NotNil)
	_, err = convertIpv4NetMaskToPrefixCheck("255.100.0.0")
	c.Check(err, C.NotNil)
	_, err = convertIpv4NetMaskToPrefixCheck("191.0.0.0")
	c.Check(err, C.NotNil)
}

func (*testWrapper) TestReverseOrderUint32(c *C.C) {
	tests := []struct {
		test   uint32
		result uint32
	}{
		{0xaabbccdd, 0xddccbbaa},
		{0x12345678, 0x78563412},
	}
	for _, t := range tests {
		c.Check(t.result, C.Equals, reverseOrderUint32(t.test))
	}
}

func (*testWrapper) TestConvertIpv6AddressToString(c *C.C) {
	tests := []struct {
		test   []byte
		result string
	}{
		{[]byte{0x12, 0x34, 0x23, 0x45, 0x34, 0x56, 0x44, 0x44, 0x55, 0x55, 0x66, 0x66, 0xaa, 0xaa, 0xff, 0xff}, "1234:2345:3456:4444:5555:6666:AAAA:FFFF"},
		{[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, "0000:0000:0000:0000:0000:0000:0000:0000"},
	}
	for _, t := range tests {
		c.Check(t.result, C.Equals, convertIpv6AddressToString(t.test))
	}
}

func (*testWrapper) TestConvertIpv6AddressToArrayByte(c *C.C) {
	tests := []struct {
		test   string
		result []byte
	}{
		{"1234:2345:3456:4444:5555:6666:AAAA:FFFF", []byte{0x12, 0x34, 0x23, 0x45, 0x34, 0x56, 0x44, 0x44, 0x55, 0x55, 0x66, 0x66, 0xaa, 0xaa, 0xff, 0xff}},
		{"0000:0000:0000:0000:0000:0000:0000:0000", []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
	}
	for _, t := range tests {
		c.Check(t.result, C.DeepEquals, convertIpv6AddressToArrayByte(t.test))
	}

	// check error ipv6 format
	_, err := convertIpv6AddressToArrayByteCheck("-1234:2345:3456:4444:5555:6666:aaAA:ffFF")
	c.Check(err, C.NotNil)
	_, err = convertIpv6AddressToArrayByteCheck("1234:2345:3456:4444:5555:6666:aaAA:ffFh")
	c.Check(err, C.NotNil)
}

func (*testWrapper) TestExpandIpv6Address(c *C.C) {
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
		c.Check(t.result, C.Equals, r)
	}

	// check error ipv6 format
	_, err := expandIpv6Address("2001::25de::cade")
	c.Check(err, C.NotNil)
}

func (*testWrapper) TestJSONWrapper(c *C.C) {
	// test json values
	testJSONKtypeString := `"test string"`
	testJSONKtypeByte := `97` // character 'a'
	testJSONKtypeInt32 := `-32`
	testJSONKtypeUint32 := `32`
	testJSONKtypeUint64 := `64`
	testJSONKtypeBoolean := `true`
	testJSONKtypeArrayByte := `"YXJyYXkgYnl0ZQ=="` // characters "array byte"
	testJSONKtypeArrayString := `["str1","str2"]`
	testJSONKtypeArrayUint32 := `[32,32]`
	testJSONKtypeArrayArrayByte := `["YXJyYXkgYnl0ZQ==","YXJyYXkgYnl0ZQ=="]`
	testJSONKtypeArrayArrayUint32 := `[[32,32],[32,32]]`
	testJSONKtypeDictStringString := `{"key1":"value1","key2":"value2"}`
	testJSONKtypeIpv6Addresses := `[{"Address":"/oAAAAAAAAACImj//g9NCQ==","Prefix":32,"Gateway":"/oAAAAAAAAACImj//g9NCQ=="}]`
	testJSONKtypeIpv6Routes := `[{"Address":"/oAAAAAAAAACImj//g9NCQ==","Prefix":32,"NextHop":"/oAAAAAAAAACImj//g9NCQ==","Metric":32}]`

	// test json values for wrapper
	testJSONKtypeWrapperString := `"123四五六"`
	testJSONKtypeWrapperStringRuneValues := []byte{0x31, 0x32, 0x33, 0xe5, 0x9b, 0x9b, 0xe4, 0xba, 0x94, 0xe5, 0x85, 0xad}
	testJSONKtypeWrapperMacAddress := `"00:12:34:56:78:AB"`
	testJSONKtypeWrapperIpv4Dns := `["192.168.1.1","192.168.1.2"]`
	testJSONKtypeWrapperIpv4Addresses := `[{"Address":"192.168.1.100","Mask":"255.255.255.0","Gateway":"192.168.1.1"},{"Address":"192.168.1.150","Mask":"128.0.0.0","Gateway":"192.168.1.1"}]`
	testJSONKtypeWrapperIpv4Routes := `[{"Address":"192.168.1.100","Mask":"255.255.192.0","NextHop":"192.168.1.1","Metric":100}]`
	testJSONKtypeWrapperIpv6Dns := `["1234:2345:3456:4444:5555:6666:AAAA:FFFF"]`
	testJSONKtypeWrapperIpv6Addresses := `[{"Address":"1111:2222:3333:4444:5555:6666:AAAA:FFFF","Prefix":64,"Gateway":"1111:2222:3333:4444:5555:6666:AAAA:1111"}]`
	testJSONKtypeWrapperIpv6Routes := `[{"Address":"1111:2222:3333:4444:5555:6666:AAAA:FFFF","Prefix":64,"NextHop":"1111:2222:3333:4444:5555:6666:AAAA:1111","Metric":32}]`

	var v interface{}
	var s string

	v, _ = jsonToKeyValue(testJSONKtypeString, ktypeString)
	s, _ = keyValueToJSON(v, ktypeString)
	c.Check(s, C.Equals, testJSONKtypeString)

	v, _ = jsonToKeyValue(testJSONKtypeByte, ktypeByte)
	s, _ = keyValueToJSON(v, ktypeByte)
	c.Check(s, C.Equals, testJSONKtypeByte)

	v, _ = jsonToKeyValue(testJSONKtypeInt32, ktypeInt32)
	s, _ = keyValueToJSON(v, ktypeInt32)
	c.Check(s, C.Equals, testJSONKtypeInt32)

	v, _ = jsonToKeyValue(testJSONKtypeUint32, ktypeUint32)
	s, _ = keyValueToJSON(v, ktypeUint32)
	c.Check(s, C.Equals, testJSONKtypeUint32)

	v, _ = jsonToKeyValue(testJSONKtypeUint64, ktypeUint64)
	s, _ = keyValueToJSON(v, ktypeUint64)
	c.Check(s, C.Equals, testJSONKtypeUint64)

	v, _ = jsonToKeyValue(testJSONKtypeBoolean, ktypeBoolean)
	s, _ = keyValueToJSON(v, ktypeBoolean)
	c.Check(s, C.Equals, testJSONKtypeBoolean)

	v, _ = jsonToKeyValue(testJSONKtypeArrayByte, ktypeArrayByte)
	s, _ = keyValueToJSON(v, ktypeArrayByte)
	c.Check(s, C.Equals, testJSONKtypeArrayByte)

	v, _ = jsonToKeyValue(testJSONKtypeArrayString, ktypeArrayString)
	s, _ = keyValueToJSON(v, ktypeArrayString)
	c.Check(s, C.Equals, testJSONKtypeArrayString)

	v, _ = jsonToKeyValue(testJSONKtypeArrayUint32, ktypeArrayUint32)
	s, _ = keyValueToJSON(v, ktypeArrayUint32)
	c.Check(s, C.Equals, testJSONKtypeArrayUint32)

	v, _ = jsonToKeyValue(testJSONKtypeArrayArrayByte, ktypeArrayArrayByte)
	s, _ = keyValueToJSON(v, ktypeArrayArrayByte)
	c.Check(s, C.Equals, testJSONKtypeArrayArrayByte)

	v, _ = jsonToKeyValue(testJSONKtypeArrayArrayUint32, ktypeArrayArrayUint32)
	s, _ = keyValueToJSON(v, ktypeArrayArrayUint32)
	c.Check(s, C.Equals, testJSONKtypeArrayArrayUint32)

	v, _ = jsonToKeyValue(testJSONKtypeDictStringString, ktypeDictStringString)
	s, _ = keyValueToJSON(v, ktypeDictStringString)
	c.Check(s, C.Equals, testJSONKtypeDictStringString)

	v, _ = jsonToKeyValue(testJSONKtypeIpv6Addresses, ktypeIpv6Addresses)
	s, _ = keyValueToJSON(v, ktypeIpv6Addresses)
	c.Check(s, C.Equals, testJSONKtypeIpv6Addresses)

	v, _ = jsonToKeyValue(testJSONKtypeIpv6Routes, ktypeIpv6Routes)
	s, _ = keyValueToJSON(v, ktypeIpv6Routes)
	c.Check(s, C.Equals, testJSONKtypeIpv6Routes)

	// key value wrapper
	v, _ = jsonToKeyValue(testJSONKtypeWrapperString, ktypeWrapperString)
	c.Check(v, C.DeepEquals, testJSONKtypeWrapperStringRuneValues)
	s, _ = keyValueToJSON(v, ktypeWrapperString)
	c.Check(s, C.Equals, testJSONKtypeWrapperString)

	v, _ = jsonToKeyValue(testJSONKtypeWrapperMacAddress, ktypeWrapperMacAddress)
	s, _ = keyValueToJSON(v, ktypeWrapperMacAddress)
	c.Check(s, C.Equals, testJSONKtypeWrapperMacAddress)

	v, _ = jsonToKeyValue(testJSONKtypeWrapperIpv4Dns, ktypeWrapperIpv4Dns)
	s, _ = keyValueToJSON(v, ktypeWrapperIpv4Dns)
	c.Check(s, C.Equals, testJSONKtypeWrapperIpv4Dns)

	v, _ = jsonToKeyValue(testJSONKtypeWrapperIpv4Addresses, ktypeWrapperIpv4Addresses)
	s, _ = keyValueToJSON(v, ktypeWrapperIpv4Addresses)
	c.Check(s, C.Equals, testJSONKtypeWrapperIpv4Addresses)

	v, _ = jsonToKeyValue(testJSONKtypeWrapperIpv4Routes, ktypeWrapperIpv4Routes)
	s, _ = keyValueToJSON(v, ktypeWrapperIpv4Routes)
	c.Check(s, C.Equals, testJSONKtypeWrapperIpv4Routes)

	v, _ = jsonToKeyValue(testJSONKtypeWrapperIpv6Dns, ktypeWrapperIpv6Dns)
	s, _ = keyValueToJSON(v, ktypeWrapperIpv6Dns)
	c.Check(s, C.Equals, testJSONKtypeWrapperIpv6Dns)

	v, _ = jsonToKeyValue(testJSONKtypeWrapperIpv6Addresses, ktypeWrapperIpv6Addresses)
	s, _ = keyValueToJSON(v, ktypeWrapperIpv6Addresses)
	c.Check(s, C.Equals, testJSONKtypeWrapperIpv6Addresses)

	v, _ = jsonToKeyValue(testJSONKtypeWrapperIpv6Routes, ktypeWrapperIpv6Routes)
	s, _ = keyValueToJSON(v, ktypeWrapperIpv6Routes)
	c.Check(s, C.Equals, testJSONKtypeWrapperIpv6Routes)
}

func (*testWrapper) TestGetterAndSetterForVirtualKey(c *C.C) {
	data := newWirelessConnectionData("", "", nil, apSecNone)

	logicSetSettingVkWirelessSecurityKeyMgmt(data, "none")
	c.Check("none", C.Equals, getSettingVkWirelessSecurityKeyMgmt(data))

	logicSetSettingVkWirelessSecurityKeyMgmt(data, "wep")
	c.Check("wep", C.Equals, getSettingVkWirelessSecurityKeyMgmt(data))

	logicSetSettingVkWirelessSecurityKeyMgmt(data, "wpa-psk")
	c.Check("wpa-psk", C.Equals, getSettingVkWirelessSecurityKeyMgmt(data))

	logicSetSettingVkWirelessSecurityKeyMgmt(data, "wpa-eap")
	c.Check("wpa-eap", C.Equals, getSettingVkWirelessSecurityKeyMgmt(data))
}

func (*testWrapper) Test8021xCertPath(c *C.C) {
	data := newWirelessConnectionData("", "", nil, apSecEap)
	logicSetSettingVk8021xEap(data, "tls")

	tests := []struct {
		test, result string
	}{
		{"/the/path", "/the/path"},
		{"file:///the/path", "/the/path"},
		{"/the/path/中文", "/the/path/中文"},
		{"file:///the/path/中文", "/the/path/中文"},
		{"/the/path/%E4%B8%AD%E6%96%87", "/the/path/%E4%B8%AD%E6%96%87"},
		{"file:///the/path/%E4%B8%AD%E6%96%87", "/the/path/%E4%B8%AD%E6%96%87"},
	}
	for _, t := range tests {
		logicSetSettingVk8021xCaCert(data, t.test)
		c.Check(getSettingVk8021xCaCert(data), C.Equals, t.result)

		logicSetSettingVk8021xClientCert(data, t.test)
		c.Check(getSettingVk8021xClientCert(data), C.Equals, t.result)

		logicSetSettingVk8021xPrivateKey(data, t.test)
		c.Check(getSettingVk8021xPrivateKey(data), C.Equals, t.result)
	}
}

func (*testWrapper) TestToUriPathFor8021x(c *C.C) {
	tests := []struct {
		test   string
		result string
	}{
		{"/the/path", "file:///the/path"},
		{"file:///the/path", "file:///the/path"},
		{"/the/path/中文", "file:///the/path/中文"},
		{"file:///the/path/中文", "file:///the/path/中文"},
		{"/the/path/%E4%B8%AD%E6%96%87", "file:///the/path/%E4%B8%AD%E6%96%87"},
	}
	for _, t := range tests {
		c.Check(t.result, C.Equals, toUriPathFor8021x(t.test))
	}
}
func (*testWrapper) TestToLocalPathFor8021x(c *C.C) {
	tests := []struct {
		test   string
		result string
	}{
		{"/the/path", "/the/path"},
		{"file:///the/path", "/the/path"},
		{"file:///the/path/%E4%B8%AD%E6%96%87", "/the/path/%E4%B8%AD%E6%96%87"},
		{"/the/path/中文", "/the/path/中文"},
		{"file:///the/path/中文", "/the/path/中文"},
		{"/the/path/%E4%B8%AD%E6%96%87", "/the/path/%E4%B8%AD%E6%96%87"},
	}
	for _, t := range tests {
		c.Check(t.result, C.Equals, toLocalPathFor8021x(t.test))
	}
}

func (*testWrapper) TestToUriPath(c *C.C) {
	tests := []struct {
		test   string
		result string
	}{
		{"/the/path", "file:///the/path"},
		{"file:///the/path", "file:///the/path"},
		{"/the/path/中文", "file:///the/path/%E4%B8%AD%E6%96%87"},
		{"file:///the/path/中文", "file:///the/path/%E4%B8%AD%E6%96%87"},
	}
	for _, t := range tests {
		c.Check(t.result, C.Equals, toUriPath(t.test))
	}
}
func (*testWrapper) TestToLocalPath(c *C.C) {
	tests := []struct {
		test   string
		result string
	}{
		{"/the/path", "/the/path"},
		{"file:///the/path", "/the/path"},
		{"file:///the/path/%E4%B8%AD%E6%96%87", "/the/path/中文"},
		{"/the/path/中文", "/the/path/中文"},
		{"file:///the/path/中文", "/the/path/中文"},
		{"/the/path/%E4%B8%AD%E6%96%87", "/the/path/%E4%B8%AD%E6%96%87"},
	}
	for _, t := range tests {
		c.Check(t.result, C.Equals, toLocalPath(t.test))
	}
}

func (*testWrapper) TestStrToByteArrayPath(c *C.C) {
	tests := []struct {
		test   string
		result []byte
	}{
		{"/the/path", []byte{0x2f, 0x74, 0x68, 0x65, 0x2f, 0x70, 0x61, 0x74, 0x68, 0x0}},
		{"/the/path/中文", []byte{0x2f, 0x74, 0x68, 0x65, 0x2f, 0x70, 0x61, 0x74, 0x68, 0x2f, 0xe4, 0xb8, 0xad, 0xe6, 0x96, 0x87, 0x0}},
	}
	for _, t := range tests {
		c.Check(t.result, C.DeepEquals, strToByteArrayPath(t.test))
	}
}
func (*testWrapper) TestByteArrayToStrPath(c *C.C) {
	tests := []struct {
		test   []byte
		result string
	}{
		{[]byte{}, ""},
		{[]byte{0x0}, ""},
		{[]byte{0x2f, 0x74, 0x68, 0x65, 0x2f, 0x70, 0x61, 0x74, 0x68, 0x0}, "/the/path"},
		{[]byte{0x2f, 0x74, 0x68, 0x65, 0x2f, 0x70, 0x61, 0x74, 0x68, 0x2f, 0xe4, 0xb8, 0xad, 0xe6, 0x96, 0x87, 0x0}, "/the/path/中文"},
	}
	for _, t := range tests {
		c.Check(t.result, C.Equals, byteArrayToStrPath(t.test))
	}
}

func (*testWrapper) TestMarshalVpnPluginKey(c *C.C) {
	tests := []struct {
		t      ktype
		test   interface{}
		result string
	}{
		{ktypeString, "string", "string"},
		{ktypeBoolean, true, "yes"},
		{ktypeBoolean, false, "no"},
		{ktypeUint32, 0, "0"},
		{ktypeUint32, 1000, "1000"},
	}
	for _, t := range tests {
		c.Check(t.result, C.Equals, marshalVpnPluginKey(t.test, t.t))
	}
}
func (*testWrapper) TestUnmarshalVpnPluginKey(c *C.C) {
	tests := []struct {
		t      ktype
		test   string
		result interface{}
	}{
		{ktypeString, "string", "string"},
		{ktypeBoolean, "yes", true},
		{ktypeBoolean, "no", false},
		{ktypeUint32, "0", uint32(0)},
		{ktypeUint32, "1000", uint32(1000)},
	}
	for _, t := range tests {
		c.Check(t.result, C.Equals, unmarshalVpnPluginKey(t.test, t.t))
	}
}

func (*testWrapper) TestParseVpnServiceFile(c *C.C) {
	fileContent := `[VPN Connection]
name=l2tp
service=org.freedesktop.NetworkManager.l2tp
program=/usr/lib/NetworkManager/nm-l2tp-service

[GNOME]
auth-dialog=/usr/lib/NetworkManager/nm-l2tp-auth-dialog
properties=/usr/lib/libnm-l2tp-properties
supports-external-ui-mode=true
`
	service, program, authdialog, properties := doParseVpnNameFile(fileContent)
	c.Check(service, C.Equals, "org.freedesktop.NetworkManager.l2tp")
	c.Check(program, C.Equals, "/usr/lib/NetworkManager/nm-l2tp-service")
	c.Check(authdialog, C.Equals, "/usr/lib/NetworkManager/nm-l2tp-auth-dialog")
	c.Check(properties, C.Equals, "/usr/lib/libnm-l2tp-properties")
}

func (*testWrapper) TestStrToUuid(c *C.C) {
	data := []struct {
		addr, uuid string
	}{
		{"", "d41d8cd9-8f00-b204-e980-0998ecf8427e"},
		{"你好", "7eca689f-0d33-89d9-dea6-6ae112e5cfd7"},
		{"12:34:56:ab:cd:ef", "fdeaa9e5-b0a9-d05a-4c5a-624d6375bc0b"},
		{"fe:dc:ba:65:43:21", "9d9bc082-cc1b-ddbb-c502-46d7499954d8"},
		{"12:34:56:AB:CD:EF", "e2667717-e697-702d-7167-4bb2c5b9f58a"},
		{"123456abcdef", "6f3b8ded-65bd-7a4d-b116-25ac84e579bb"},
		{"12:34:56:ab:cd:xy", "c3701a18-6af4-aa02-7c54-53c09ea75e62"},
		{":34:56:ab:cd:ef", "2f2aab1d-d983-2df8-fe91-8598e79fc009"},
		{"123456abcdef1234abcd123456abcdef", "2fc8f109-cc40-de78-b0c4-1744b9ea62f0"},
		{"123456abcdef1234abcd123456abcdef1234", "18a1eaac-9a1e-3828-8191-511317dc2921"},
	}
	for _, d := range data {
		c.Check(d.uuid, C.Equals, strToUuid(d.addr))
	}
}

func (*testWrapper) TestDoStrToUuid(c *C.C) {
	data := []struct {
		addr, uuid string
	}{
		{"", "00000000-0000-0000-0000-000000000000"},
		{"你好", "00000000-0000-0000-0000-000000000000"},
		{"12:34:56:ab:cd:ef", "00000000-0000-0000-0000-123456abcdef"},
		{"fe:dc:ba:65:43:21", "00000000-0000-0000-0000-fedcba654321"},
		{"12:34:56:AB:CD:EF", "00000000-0000-0000-0000-123456abcdef"},
		{"123456abcdef", "00000000-0000-0000-0000-123456abcdef"},
		{"12:34:56:ab:cd:xy", "00000000-0000-0000-0000-00123456abcd"},
		{":34:56:ab:cd:ef", "00000000-0000-0000-0000-003456abcdef"},
		{"123456abcdef1234abcd123456abcdef", "123456ab-cdef-1234-abcd-123456abcdef"},
		{"123456abcdef1234abcd123456abcdef1234", "123456ab-cdef-1234-abcd-123456abcdef"},
	}
	for _, d := range data {
		c.Check(d.uuid, C.Equals, doStrToUuid(d.addr))
	}
}

func (*testWrapper) TestMarshalMobilePlanKey(c *C.C) {
	wantJSON := "\"{\\\"IsGSM\\\":true,\\\"Name\\\":\\\"LaptopConnect (data cards)\\\",\\\"ProviderName\\\":\\\"AT\\\\u0026T\\\",\\\"APNValue\\\":\\\"Broadband\\\",\\\"APNUsageType\\\":\\\"internet\\\"}\""
	wantValue := "{\"IsGSM\":true,\"Name\":\"LaptopConnect (data cards)\",\"ProviderName\":\"AT\\u0026T\",\"APNValue\":\"Broadband\",\"APNUsageType\":\"internet\"}"

	jsonStr, _ := keyValueToJSON(wantValue, ktypeString)
	c.Check(jsonStr, C.Equals, wantJSON)

	value, _ := jsonToKeyValueString(wantJSON)
	c.Check(value, C.Equals, wantValue)
}

func (*testWrapper) TestFixupDeviceDesc(c *C.C) {
	data := []struct {
		desc, fixedDesc string
	}{
		{"Intel Corporation 82567LM Gigabit Network Connection", "Intel 82567LM Gigabit"},
		{"Intel Corporation PRO/Wireless 5100 AGN [Shiloh] Network Connection", "Intel PRO/Wireless 5100 AGN [Shiloh]"},
		{"Ralink Technology, Corp. RT5370 Wireless Adapter", "Ralink RT5370"},
		{"Realtek RTL8111/8168/8411 PCI Express Gigabit Ethernet Controller (Motherboard)", "Realtek RTL8111/8168/8411 Gigabit"},
		{"Kontron (Industrial Computer Source / ICS Advent) DM9601 Fast Ethernet Adapter", "Kontron DM9601"},
	}
	for _, d := range data {
		c.Check(fixupDeviceDesc(d.desc), C.Equals, d.fixedDesc)
	}
}

// no font or locale in jenkins
// func (*testWrapper) TestDecodeSsid(c *C.C) {
// 	var ssidList = []struct {
// 		ssid []byte
// 		name string
// 	}{
// 		{
// 			// gbk
// 			ssid: []byte{178, 226, 202, 212, 45, 119, 105, 102, 105},
// 			name: "测试-wifi",
// 		},
// 		{
// 			// utf8
// 			ssid: []byte{231, 159, 179, 231, 154, 132, 119, 105, 102, 105},
// 			name: "石的wifi",
// 		},
// 	}

// 	for _, info := range ssidList {
// 		c.Check(decodeSsid(info.ssid), C.Equals, info.name)
// 	}
// }
