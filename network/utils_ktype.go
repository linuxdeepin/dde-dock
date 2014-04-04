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
	ktypeArrayString
	ktypeArrayUint32
	ktypeArrayArrayByte   // [array of array of byte]
	ktypeArrayArrayUint32 // [array of array of uint32]
	ktypeDictStringString // [dict of (string::string)]
	ktypeIpv6Addresses    // [array of (byte array, uint32, byte array)]
	ktypeIpv6Routes       // [array of (byte array, uint32, byte array, uint32)]

	// wrapper for special key type, used by json setter and getter
	ktypeWrapperString        // wrap ktypeArrayByte to [string]
	ktypeWrapperMacAddress    // wrap ktypeArrayByte to [string]
	ktypeWrapperIpv4Dns       // wrap ktypeArrayUint32 to [array of string]
	ktypeWrapperIpv4Addresses // wrap ktypeArrayArrayUint32 to [array of (string, string, string)]
	ktypeWrapperIpv4Routes    // wrap ktypeArrayArrayUint32 to [array of (string, string, string, uint32)]
	ktypeWrapperIpv6Dns       // wrap ktypeArrayArrayByte to [array of string]
	ktypeWrapperIpv6Addresses // wrap ktypeIpv6Addresses to [array of (string, uint32, string)]
	ktypeWrapperIpv6Routes    // wrap ktypeIpv6Routes to [array of (string, uint32, string, uint32)]
)

// Ipv4AddressesWrapper
type Ipv4AddressesWrapper []Ipv4AddressWrapper
type Ipv4AddressWrapper struct {
	Address string
	Prefix  string
	Gateway string
}

// Ipv4RoutesWrapper
type Ipv4RoutesWrapper []Ipv4RouteWrapper
type Ipv4RouteWrapper struct {
	Address string
	Prefix  string
	NextHop string
	Metric  uint32
}

// Ipv6AddressesWrapper
type Ipv6AddressesWrapper []Ipv6AddressWrapper
type Ipv6AddressWrapper struct {
	Address string
	Prefix  uint32
	Gateway string
}

// Ipv6Addresses is an array of (byte array, uint32, byte array)
type Ipv6Addresses []Ipv6Address
type Ipv6Address struct {
	Address []byte
	Prefix  uint32
	Gateway []byte
}

// Ipv6RoutesWrapper
type Ipv6RoutesWrapper []Ipv6RouteWrapper
type Ipv6RouteWrapper struct {
	Address string
	Prefix  uint32
	NextHop string
	Metric  uint32
}

// Ipv6Routes is an array of (byte array, uint32, byte array, uint32)
type Ipv6Route struct {
	Address []byte
	Prefix  uint32
	NextHop []byte
	Metric  uint32
}
type Ipv6Routes []Ipv6Route

// TODO remove
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
