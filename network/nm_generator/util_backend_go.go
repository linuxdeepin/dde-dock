/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package main

import (
	"fmt"
	"os"
	"strings"
)

// "ktypeString" -> "String", "ktypeBoolean" -> "Boolean"
func ToKeyTypeShortName(ktype string) string {
	return strings.TrimPrefix(ktype, "ktype")
}

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
	case "ktypeInt64":
		realData = "int64"
	case "ktypeUint64":
		realData = "uint64"
	case "ktypeBoolean":
		realData = "bool"
	case "ktypeArrayByte", "ktypeWrapperString", "ktypeWrapperMacAddress":
		realData = "[]byte"
	case "ktypeArrayString":
		realData = "[]string"
	case "ktypeArrayUint32", "ktypeWrapperIpv4Dns":
		realData = "[]uint32"
	case "ktypeArrayArrayByte", "ktypeWrapperIpv6Dns":
		realData = "[][]byte"
	case "ktypeArrayArrayUint32", "ktypeWrapperIpv4Addresses", "ktypeWrapperIpv4Routes":
		realData = "[][]uint32"
	case "ktypeDictStringString":
		realData = "map[string]string"
	case "ktypeIpv6Addresses", "ktypeWrapperIpv6Addresses":
		realData = "ipv6Addresses"
	case "ktypeIpv6Routes", "ktypeWrapperIpv6Routes":
		realData = "ipv6Routes"
	}
	return
}

// "ktypeString" -> `""`, "ktypeBool" -> `false`
func ToKeyDefaultValue(keyName string) (value string) {
	keyInfo := getKeyInfo(keyName)
	ktype := keyInfo.Type
	customValue := keyInfo.Default
	if customValue == "<null>" {
		if ktype == "ktypeString" {
			return `""`
		} else {
			return "nil"
		}
	} else if customValue != "<default>" {
		return customValue
	}
	switch ktype {
	default:
		fmt.Println("invalid ktype:", ktype)
		os.Exit(1)
	case "ktypeString":
		value = `""`
	case "ktypeByte":
		value = `byte(0)`
	case "ktypeInt32":
		value = `int32(0)`
	case "ktypeUint32":
		value = `uint32(0)`
	case "ktypeInt64":
		value = `int64(0)`
	case "ktypeUint64":
		value = `uint64(0)`
	case "ktypeBoolean":
		value = `false`
	case "ktypeArrayByte", "ktypeWrapperString", "ktypeWrapperMacAddress":
		value = `make([]byte, 0)`
		// value = `nil`
	case "ktypeArrayString":
		value = `make([]string, 0)`
		// value = `nil`
	case "ktypeArrayUint32", "ktypeWrapperIpv4Dns":
		value = `make([]uint32, 0)`
		// value = `nil`
	case "ktypeArrayArrayByte", "ktypeWrapperIpv6Dns":
		value = `make([][]byte, 0)`
		// value = `nil`
	case "ktypeArrayArrayUint32", "ktypeWrapperIpv4Addresses", "ktypeWrapperIpv4Routes":
		value = `make([][]uint32, 0)`
		// value = `nil`
	case "ktypeDictStringString":
		value = `make(map[string]string)`
		// value = `nil`
	case "ktypeIpv6Addresses", "ktypeWrapperIpv6Addresses":
		value = `make(ipv6Addresses, 0)`
		// value = `nil`
	case "ktypeIpv6Routes", "ktypeWrapperIpv6Routes":
		value = `make(ipv6Routes, 0)`
		// value = `nil`
	}
	return
}

// "ktypeString" -> interfaceToString, "ktypeBool" -> interfaceToBoolean
func ToKeyTypeInterfaceConverter(ktype string) (converter string) {
	switch ktype {
	default:
		fmt.Println("invalid ktype:", ktype)
		os.Exit(1)
	case "ktypeString":
		converter = "interfaceToString"
	case "ktypeByte":
		converter = "interfaceToByte"
	case "ktypeInt32":
		converter = "interfaceToInt32"
	case "ktypeUint32":
		converter = "interfaceToUint32"
	case "ktypeInt64":
		converter = "interfaceToInt64"
	case "ktypeUint64":
		converter = "interfaceToUint64"
	case "ktypeBoolean":
		converter = "interfaceToBoolean"
	case "ktypeArrayByte", "ktypeWrapperString", "ktypeWrapperMacAddress":
		converter = "interfaceToArrayByte"
	case "ktypeArrayString":
		converter = "interfaceToArrayString"
	case "ktypeArrayUint32", "ktypeWrapperIpv4Dns":
		converter = "interfaceToArrayUint32"
	case "ktypeArrayArrayByte", "ktypeWrapperIpv6Dns":
		converter = "interfaceToArrayArrayByte"
	case "ktypeArrayArrayUint32", "ktypeWrapperIpv4Addresses", "ktypeWrapperIpv4Routes":
		converter = "interfaceToArrayArrayUint32"
	case "ktypeDictStringString":
		converter = "interfaceToDictStringString"
	case "ktypeIpv6Addresses", "ktypeWrapperIpv6Addresses":
		converter = "interfaceToIpv6Addresses"
	case "ktypeIpv6Routes", "ktypeWrapperIpv6Routes":
		converter = "interfaceToIpv6Routes"
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
	case "ktypeInt64":
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

// get all related sections of virtual keys
func GetAllVkeysRelatedSections(nmVkeys []NMVkeyStruct) (sections []string) {
	for _, vk := range nmVkeys {
		sections = appendStrArrayUnique(sections, vk.RelatedSection)
	}
	return
}

// get all virtual keys in target section
func GetVkeysOfSection(nmVkeys []NMVkeyStruct, section string) (keys []string) {
	for _, vk := range nmVkeys {
		if vk.RelatedSection == section {
			keys = append(keys, vk.Name)
		}
	}
	return
}
