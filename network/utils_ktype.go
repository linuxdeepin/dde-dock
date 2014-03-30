package main

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

func getKtypeDescription(t ktype) (desc string) {
	switch t {
	case ktypeUnknown:
		desc = "Unknown"
	case ktypeString:
		desc = "String"
	case ktypeByte:
		desc = "Byte"
	case ktypeInt32:
		desc = "Int32"
	case ktypeUint32:
		desc = "Uint32"
	case ktypeUint64:
		desc = "Uint64"
	case ktypeBoolean:
		desc = "Boolean"
	case ktypeArrayByte:
		desc = "ArrayByte"
	case ktypeArrayString:
		desc = "ArrayString, encode by json"
	case ktypeArrayUint32:
		desc = "ArrayUint32, encode by json"
	case ktypeArrayArrayByte:
		desc = "ArrayArrayByte, array of array of byte, encode by json"
	case ktypeArrayArrayUint32:
		desc = "ArrayArrayUint32, array of array of uint32, encode by json"
	case ktypeDictStringString:
		desc = "DictStringString, dict of (string::string), encode by json"
	case ktypeIpv6Addresses:
		desc = "Ipv6Addresses, array of (byte array, uint32, byte array), encode by json"
	case ktypeIpv6Routes:
		desc = "Ipv6Routes, array of (byte array, uint32, byte array, uint32), encode by json"
	}
	return
}

// Ipv6Addresses is an array of (byte array, uint32, byte array)
type Ipv6Address struct {
	Address []byte
	Prefix  uint32
	Gateway []byte
}
type Ipv6Addresses []Ipv6Address

// Ipv6Routes is an array of (byte array, uint32, byte array, uint32)
type Ipv6Route struct {
	Address []byte
	Prefix  uint32
	NextHop []byte
	Metric  uint32
}
type Ipv6Routes []Ipv6Route
