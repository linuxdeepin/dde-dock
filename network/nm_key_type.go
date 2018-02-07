/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package network

import (
	"sort"
)

// available values structure
type availableValues map[string]kvalue
type kvalue struct {
	Value interface{}
	Text  string // used for internationalization
}

type kvalues []kvalue

func (ks kvalues) Len() int {
	return len(ks)
}
func (ks kvalues) Swap(i, j int) {
	ks[i], ks[j] = ks[j], ks[i]
}
func (ks kvalues) Less(i, j int) bool {
	return ks[i].Text < ks[j].Text
}

func sortKvalues(ks []kvalue) {
	sort.Sort(kvalues(ks))
}

// define key type
type ktype uint32

const (
	ktypeUnknown ktype = iota
	ktypeString
	ktypeByte // for byte and gchar type
	ktypeInt32
	ktypeUint32
	ktypeInt64
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

	// wrapper for special key type, used by json getter and setter,
	// in other words, only works for front-end
	ktypeWrapperString        // wrap ktypeArrayByte to [string]
	ktypeWrapperMacAddress    // wrap ktypeArrayByte to [string]
	ktypeWrapperIpv4Dns       // wrap ktypeArrayUint32 to [array of string]
	ktypeWrapperIpv4Addresses // wrap ktypeArrayArrayUint32 to [array of (string, string, string)]
	ktypeWrapperIpv4Routes    // wrap ktypeArrayArrayUint32 to [array of (string, string, string, uint32)]
	ktypeWrapperIpv6Dns       // wrap ktypeArrayArrayByte to [array of string]
	ktypeWrapperIpv6Addresses // wrap ktypeIpv6Addresses to [array of (string, uint32, string)]
	ktypeWrapperIpv6Routes    // wrap ktypeIpv6Routes to [array of (string, uint32, string, uint32)]
)

func isWrapperKeyType(t ktype) bool {
	switch t {
	case ktypeWrapperString:
		return true
	case ktypeWrapperMacAddress:
		return true
	case ktypeWrapperIpv4Dns:
		return true
	case ktypeWrapperIpv4Addresses:
		return true
	case ktypeWrapperIpv4Routes:
		return true
	case ktypeWrapperIpv6Dns:
		return true
	case ktypeWrapperIpv6Addresses:
		return true
	case ktypeWrapperIpv6Routes:
		return true
	}
	return false
}

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

func getKtypeDesc(t ktype) (desc string) {
	switch t {
	default:
		logger.Error("Unknown type", t)
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
	case ktypeWrapperString:
		desc = "wrap ktypeArrayByte to [string]"
	case ktypeWrapperMacAddress:
		desc = "wrap ktypeArrayByte to [string]"
	case ktypeWrapperIpv4Dns:
		desc = "wrap ktypeArrayUint32 to [array of string]"
	case ktypeWrapperIpv4Addresses:
		desc = "wrap ktypeArrayArrayUint32 to [array of (string, string, string)]"
	case ktypeWrapperIpv4Routes:
		desc = "wrap ktypeArrayArrayUint32 to [array of (string, string, string, uint32)]"
	case ktypeWrapperIpv6Dns:
		desc = "wrap ktypeArrayArrayByte to [array of string]"
	case ktypeWrapperIpv6Addresses:
		desc = "wrap ktypeIpv6Addresses to [array of (string, uint32, string)]"
	case ktypeWrapperIpv6Routes:
		desc = "wrap ktypeIpv6Routes to [array of (string, uint32, string, uint32)]"
	}
	return
}
