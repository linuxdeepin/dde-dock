package main

// available values structure
type availableValues map[string]kvalue
type kvalue struct {
	Value interface{}
	Text  string // used for internationalization
}

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
type ipv4AddressesWrapper []ipv4AddressWrapper
type ipv4AddressWrapper struct {
	Address string
	Mask    string
	Gateway string
}

// Ipv4RoutesWrapper
type ipv4RoutesWrapper []ipv4RouteWrapper
type ipv4RouteWrapper struct {
	Address string
	Mask    string
	NextHop string
	Metric  uint32
}

// Ipv6AddressesWrapper
type ipv6AddressesWrapper []ipv6AddressWrapper
type ipv6AddressWrapper struct {
	Address string
	Prefix  uint32
	Gateway string
}

// Ipv6Addresses is an array of (byte array, uint32, byte array)
type ipv6Addresses []ipv6Address
type ipv6Address struct {
	Address []byte
	Prefix  uint32
	Gateway []byte
}

// ipv6RoutesWrapper
type ipv6RoutesWrapper []ipv6RouteWrapper
type ipv6RouteWrapper struct {
	Address string
	Prefix  uint32
	NextHop string
	Metric  uint32
}

// ipv6Routes is an array of (byte array, uint32, byte array, uint32)
type ipv6Route struct {
	Address []byte
	Prefix  uint32
	NextHop []byte
	Metric  uint32
}
type ipv6Routes []ipv6Route

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
		desc = "ipv6Routes, array of (byte array, uint32, byte array, uint32), encode by json"
	}
	return
}
