/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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
	"fmt"
	"strconv"
	"strings"
)

const (
	macAddrZero  = "00:00:00:00:00:00"
	ipv4Zero     = "0.0.0.0"
	ipv6AddrZero = "0000:0000:0000:0000:0000:0000:0000:0000"
)

// []byte{0,0,0,0,0,0} -> "00:00:00:00:00:00"
func convertMacAddressToString(v []byte) (macAddr string) {
	if len(v) != 6 {
		macAddr = macAddrZero
		logger.Error("machine address is invalid", v)
		return
	}
	macAddr = fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", v[0], v[1], v[2], v[3], v[4], v[5])
	return
}

// "00:00:00:00:00:00" -> []byte{0,0,0,0,0,0}
func convertMacAddressToArrayByte(v string) (macAddr []byte) {
	macAddr, err := convertMacAddressToArrayByteCheck(v)
	if err != nil {
		logger.Error(err)
	}
	return
}
func convertMacAddressToArrayByteCheck(v string) (macAddr []byte, err error) {
	macAddr = make([]byte, 6)
	a := strings.Split(v, ":")
	if len(a) != 6 {
		err = fmt.Errorf("machine address is invalid %s", v)
		return
	}
	for i := 0; i < 6; i++ {
		if len(a[i]) != 2 {
			err = fmt.Errorf("machine address is invalid %s", v)
			return
		}
		var n uint64
		n, err = strconv.ParseUint(a[i], 16, 8)
		if err != nil {
			err = fmt.Errorf("machine address is invalid %s", v)
			return
		}
		macAddr[i] = byte(n)
	}
	return
}

func convertIpv4AddressToString(v uint32) (ip4Addr string) {
	ip4Addr = fmt.Sprintf("%d.%d.%d.%d", byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
	return
}

// if address is 0, return empty string instead of "0.0.0.0"
func convertIpv4AddressToStringNoZero(v uint32) (ip4Addr string) {
	if v == 0 {
		return
	} else {
		ip4Addr = fmt.Sprintf("%d.%d.%d.%d", byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
	}
	return
}

func convertIpv4AddressToUint32(v string) (ip4Addr uint32) {
	ip4Addr, err := convertIpv4AddressToUint32Check(v)
	if err != nil {
		logger.Error(err)
	}
	return
}
func convertIpv4AddressToUint32Check(v string) (ip4Addr uint32, err error) {
	a := strings.Split(v, ".")
	if len(a) != 4 {
		ip4Addr = 0
		err = fmt.Errorf("ip address is invalid %s", v)
		return
	}
	for i := 3; i >= 0; i-- {
		var tmpn uint64
		tmpn, err = strconv.ParseUint(a[i], 10, 8)
		if err != nil {
			err = fmt.Errorf("ip address is invalid %s", v)
			return
		}
		ip4Addr = ip4Addr<<8 + uint32(tmpn)
	}
	return
}

// host order to network order, or network order to host order
func reverseOrderUint32(net uint32) (host uint32) {
	host = uint32(byte(net>>24)) << 0
	host |= uint32(byte(net>>16)) << 8
	host |= uint32(byte(net>>8)) << 16
	host |= uint32(byte(net>>0)) << 24
	return
}

// 24 -> "255.255.255.0"(string format)
func convertIpv4PrefixToNetMask(prefix uint32) (maskAddress string) {
	var mask uint32
	for i := uint32(0); i < prefix; i++ {
		mask = mask<<1 + 1
	}
	for i := prefix; i < 32; i++ {
		mask = mask<<1 + 0
	}
	mask = reverseOrderUint32(mask)
	maskAddress = convertIpv4AddressToString(mask)
	return
}

// "255.255.255.0"(string format) -> 24
func convertIpv4NetMaskToPrefix(maskAddress string) (prefix uint32) {
	prefix, err := convertIpv4NetMaskToPrefixCheck(maskAddress)
	if err != nil {
		logger.Error(err)
	}
	return
}
func convertIpv4NetMaskToPrefixCheck(maskAddress string) (prefix uint32, err error) {
	var mask uint32 // network order
	mask, err = convertIpv4AddressToUint32Check(maskAddress)
	if err != nil {
		return
	}
	mask = reverseOrderUint32(mask)
	foundZerorBit := false
	for i := uint32(0); i < 32; i++ {
		if mask>>(32-i-1)&0x01 == 1 {
			if !foundZerorBit {
				prefix++
			} else {
				err = fmt.Errorf("net mask address is invalid %s", maskAddress)
				return
			}
		} else {
			foundZerorBit = true
			continue
		}
	}
	return
}

// []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0} -> "0000:0000:0000:0000:0000:0000:0000:0000"
func convertIpv6AddressToString(v []byte) (ipv6Addr string) {
	ipv6Addr, err := convertIpv6AddressToStringCheck(v)
	if err != nil {
		logger.Error(err)
	}
	return
}
func convertIpv6AddressToStringCheck(v []byte) (ipv6Addr string, err error) {
	if len(v) != 16 {
		ipv6Addr = ipv6AddrZero
		err = fmt.Errorf("ipv6 address is invalid %s", v)
		return
	}
	for i := 0; i < 16; i += 2 {
		s := fmt.Sprintf("%02X%02X", v[i], v[i+1])
		if len(ipv6Addr) == 0 {
			ipv6Addr = s
		} else {
			ipv6Addr = ipv6Addr + ":" + s
		}
	}
	return
}

// if address is 0, return empty string instead of "0.0.0.0"
func convertIpv6AddressToStringNoZero(v []byte) (ipv6Addr string) {
	if isIpv6AddressZero(v) {
		return
	}
	return convertIpv6AddressToString(v)
}

// "0000:0000:0000:0000:0000:0000:0000:0000" -> []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}
func convertIpv6AddressToArrayByte(v string) (ipv6Addr []byte) {
	ipv6Addr, err := convertIpv6AddressToArrayByteCheck(v)
	if err != nil {
		logger.Error(err)
	}
	return
}
func convertIpv6AddressToArrayByteCheck(v string) (ipv6Addr []byte, err error) {
	v, err = expandIpv6Address(v)
	if err != nil {
		return
	}
	ipv6Addr = make([]byte, 16)
	a := strings.Split(v, ":")
	if len(a) != 8 {
		err = fmt.Errorf("ipv6 address is invalid %s", v)
		return
	}
	for i := 0; i < 8; i++ {
		s := a[i]
		if len(s) != 4 {
			err = fmt.Errorf("ipv6 address is invalid %s", v)
			return
		}

		var tmpn uint64
		tmpn, err = strconv.ParseUint(s[:2], 16, 8)
		ipv6Addr[i*2] = byte(tmpn)
		if err != nil {
			err = fmt.Errorf("ipv6 address is invalid %s", v)
			return
		}

		tmpn, err = strconv.ParseUint(s[2:], 16, 8)
		if err != nil {
			err = fmt.Errorf("ipv6 address is invalid %s", v)
			return
		}
		ipv6Addr[i*2+1] = byte(tmpn)
	}
	return
}

// expand ipv6 address to standard format, such as
// "0::0" -> "0000:0000:0000:0000:0000:0000:0000:0000"
// "2001:DB8:2de::e13" -> "2001:DB8:2de:0:0:0:0:e13"
// "2001::25de::cade" -> error
func expandIpv6Address(v string) (ipv6Addr string, err error) {
	a1 := strings.Split(v, ":")
	l1 := len(a1)
	if l1 > 8 {
		err = fmt.Errorf("invalid ipv6 address %s", v)
		return
	}

	a2 := strings.Split(v, "::")
	l2 := len(a2)
	if l2 > 2 {
		err = fmt.Errorf("invalid ipv6 address %s", v)
		return
	} else if l2 == 2 {
		// expand "::"
		abbrFields := ":"
		for i := 0; i <= 8-l1; i++ {
			abbrFields += "0000:"
		}
		v = strings.Replace(v, "::", abbrFields, -1)
	}

	// expand ":0:" to ":0000:"
	a1 = strings.Split(v, ":")
	for i, field := range a1 {
		l := len(field)
		if l > 4 {
			err = fmt.Errorf("invalid ipv6 address %s", v)
			return
		} else if l < 4 {
			field = strings.Repeat("0", 4-l) + field
		}
		if i == 0 {
			ipv6Addr = field
		} else {
			ipv6Addr += ":" + field
		}
	}
	return
}

func isIpv6AddressValid(v []byte) bool {
	if len(v) != 16 {
		return false
	}
	return true
}

func isIpv6AddressZero(v []byte) bool {
	// don't care if ipv6 address if is valid
	allAreZero := true
	for _, b := range v {
		if b != 0 {
			allAreZero = false
			break
		}
	}
	return allAreZero
}

func isIpv6AddressStructZero(addr ipv6Address) bool {
	if isIpv6AddressZero(addr.Address) && isIpv6AddressZero(addr.Gateway) && addr.Prefix == 0 {
		return true
	}
	return false
}

func isIpv6RouteStructZero(route ipv6Route) bool {
	return isIpv6AddressZero(route.Address) && isIpv6AddressZero(route.NextHop) &&
		(route.Prefix == 0) && (route.Metric == 0)
}
