package main

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
		Logger.Error("machine address is invalid", v)
		return
	}
	macAddr = fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", v[0], v[1], v[2], v[3], v[4], v[5])
	return
}

// "00:00:00:00:00:00" -> []byte{0,0,0,0,0,0}
func convertMacAddressToArrayByte(v string) (macAddr []byte) {
	macAddr = make([]byte, 6)
	a := strings.Split(v, ":")
	if len(a) != 6 {
		Logger.Error("machine address is invalid", v)
		return
	}
	for i := 0; i < 6; i++ {
		n, err := strconv.ParseUint(a[i], 16, 8)
		if err != nil {
			Logger.Error("machine address is invalid", v)
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
	if len(v) == 0 {
		v = ipv4Zero // convert empty string to "0.0.0.0"
	}
	a := strings.Split(v, ".")
	if len(a) != 4 {
		ip4Addr = 0
		Logger.Error("ip address is invalid", v)
		return
	}
	for i := 3; i >= 0; i-- {
		tmpn, err := strconv.ParseUint(a[i], 10, 8)
		if err != nil {
			Logger.Error("ip address is invalid", v)
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
	var mask uint32 // network order
	mask = convertIpv4AddressToUint32(maskAddress)
	mask = reverseOrderUint32(mask)
	for i := uint32(0); i < 32; i++ {
		if mask>>(32-i-1)&0x01 == 1 {
			prefix++
		} else {
			break
		}
	}
	return
}

// []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0} -> "0000:0000:0000:0000:0000:0000:0000:0000"
func convertIpv6AddressToString(v []byte) (ipv6Addr string) {
	if len(v) != 16 {
		ipv6Addr = ipv6AddrZero
		Logger.Error("ipv6 address is invalid", v)
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

// "0000:0000:0000:0000:0000:0000:0000:0000" -> []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}
func convertIpv6AddressToArrayByte(v string) (ipv6Addr []byte) {
	ipv6Addr = make([]byte, 16)
	a := strings.Split(v, ":")
	if len(a) != 8 {
		Logger.Error("ipv6 address is invalid", v)
		return // TODO
	}
	for i := 0; i < 8; i++ {
		s := a[i]
		if len(s) != 4 {
			Logger.Error("ipv6 address is invalid", v)
			return
		}

		n, err := strconv.ParseUint(s[:2], 16, 8)
		ipv6Addr[i*2] = byte(n)
		if err != nil {
			Logger.Error("ipv6 address is invalid", v)
			return
		}

		n, err = strconv.ParseUint(s[2:], 16, 8)
		if err != nil {
			Logger.Error("ipv6 address is invalid", v)
			return
		}
		ipv6Addr[i*2+1] = byte(n)
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
	allAreZero := true
	for _, b := range v {
		if b != 0 {
			allAreZero = false
			break
		}
	}
	return allAreZero
}
