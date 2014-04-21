package main

import (
	"fmt"
	"os"
)

// "ktypeString" -> "string", "ktypeArrayByte" -> "[]byte"
func ToKeyTypeRealData(ktype string) (realData string) {
	switch ktype {
	default:
		fmt.Println("invalid ktype:", ktype)
		os.Exit(1)
	case "ktypeString":
		realData = "string"
	case "ktypeByte":
		realData = "byte"
	case "ktypeInt32":
		realData = "int32"
	case "ktypeUint32":
		realData = "uint32"
	case "ktypeUint64":
		realData = "uint64"
	case "ktypeBoolean":
		realData = "bool"
	case "ktypeArrayByte":
		realData = "[]byte"
	case "ktypeArrayString":
		realData = "[]string"
	case "ktypeArrayUint32":
		realData = "[]uint32"
	case "ktypeArrayArrayByte":
		realData = "[][]byte"
	case "ktypeArrayArrayUint32":
		realData = "[][]uint32"
	case "ktypeDictStringString":
		realData = "map[string]string"
	case "ktypeIpv6Addresses":
		realData = "Ipv6Addresses"
	case "ktypeIpv6Routes":
		realData = "Ipv6Routes"
	case "ktypeWrapperString":
		realData = "[]byte"
	case "ktypeWrapperMacAddress":
		realData = "[]byte"
	case "ktypeWrapperIpv4Dns":
		realData = "[]uint32"
	case "ktypeWrapperIpv4Addresses":
		realData = "[][]uint32"
	case "ktypeWrapperIpv4Routes":
		realData = "[][]uint32"
	case "ktypeWrapperIpv6Dns":
		realData = "[][]byte"
	case "ktypeWrapperIpv6Addresses":
		realData = "Ipv6Addresses"
	case "ktypeWrapperIpv6Routes":
		realData = "Ipv6Routes"
	}
	return
}

// "ktypeString" -> `""`, "ktypeBool" -> `false`
func ToKeyTypeDefaultValueJSON(ktype, customValue string) (valueJSON string) {
	if customValue == "<null>" {
		if ktype == "ktypeString" {
			return `""`
		} else {
			return "null"
		}
	} else if customValue != "<default>" {
		return customValue
	}
	switch ktype {
	default:
		fmt.Println("invalid ktype:", ktype)
		os.Exit(1)
	case "ktypeString":
		valueJSON = `""`
	case "ktypeByte":
		valueJSON = `0`
	case "ktypeInt32":
		valueJSON = `0`
	case "ktypeUint32":
		valueJSON = `0`
	case "ktypeUint64":
		valueJSON = `0`
	case "ktypeBoolean":
		valueJSON = `false`
	case "ktypeArrayByte":
		valueJSON = `""`
	case "ktypeArrayString":
		valueJSON = `null`
	case "ktypeArrayUint32":
		valueJSON = `null`
	case "ktypeArrayArrayByte":
		valueJSON = `null`
	case "ktypeArrayArrayUint32":
		valueJSON = `null`
	case "ktypeDictStringString":
		valueJSON = `null`
	case "ktypeIpv6Addresses":
		valueJSON = `null`
	case "ktypeIpv6Routes":
		valueJSON = `null`
	case "ktypeWrapperString":
		valueJSON = `""`
	case "ktypeWrapperMacAddress":
		valueJSON = `""`
	case "ktypeWrapperIpv4Dns":
		valueJSON = `null`
	case "ktypeWrapperIpv4Addresses":
		valueJSON = `null`
	case "ktypeWrapperIpv4Routes":
		valueJSON = `null`
	case "ktypeWrapperIpv6Dns":
		valueJSON = `null`
	case "ktypeWrapperIpv6Addresses":
		valueJSON = `null`
	case "ktypeWrapperIpv6Routes":
		valueJSON = `null`
	}
	return
}

// test if need check value length to ensure value not empty
func IfNeedCheckValueLength(ktype string) (need string) {
	switch ktype {
	default:
		fmt.Println("invalid ktype:", ktype)
		os.Exit(1)
	case "ktypeString":
		need = "t"
	case "ktypeByte":
		need = ""
	case "ktypeInt32":
		need = ""
	case "ktypeUint32":
		need = ""
	case "ktypeUint64":
		need = ""
	case "ktypeBoolean":
		need = ""
	case "ktypeArrayByte":
		need = "t"
	case "ktypeArrayString":
		need = "t"
	case "ktypeArrayUint32":
		need = "t"
	case "ktypeArrayArrayByte":
		need = "t"
	case "ktypeArrayArrayUint32":
		need = "t"
	case "ktypeDictStringString":
		need = "t"
	case "ktypeIpv6Addresses":
		need = "t"
	case "ktypeIpv6Routes":
		need = "t"
	case "ktypeWrapperString":
		need = "t"
	case "ktypeWrapperMacAddress":
		need = "t"
	case "ktypeWrapperIpv4Dns":
		need = "t"
	case "ktypeWrapperIpv4Addresses":
		need = "t"
	case "ktypeWrapperIpv4Routes":
		need = "t"
	case "ktypeWrapperIpv6Dns":
		need = "t"
	case "ktypeWrapperIpv6Addresses":
		need = "t"
	case "ktypeWrapperIpv6Routes":
		need = "t"
	}
	return
}
