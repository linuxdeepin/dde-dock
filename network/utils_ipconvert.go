package main

import (
	"fmt"
	"strconv"
	"strings"
)

// []byte{0,0,0,0,0,0} -> "00:00:00:00:00:00"
func convertMacAddressToString(v []byte) (macAddr string) {
	if len(v) != 6 {
		macAddr = "00:00:00:00:00:00"
		LOGGER.Error("convertMacAddressToString, machine address is invalid", v)
		return
	}
	macAddr = fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", v[0], v[1], v[2], v[3], v[4], v[5])
	return
}

// "00:00:00:00:00:00" -> []byte{0,0,0,0,0,0}
func convertMacAddressToArrayByte(v string) (macAddr []byte) {
	a := strings.Split(v, ":")
	if len(a) != 6 {
		macAddr = []byte{0, 0, 0, 0, 0, 0}
		LOGGER.Error("convertMacAddressToArrayByte, machine address is invalid", v)
		return
	}
	macAddr = make([]byte, 6)
	for i := 0; i < 6; i++ {
		tmpn, err := strconv.ParseUint(a[i], 16, 8)
		if err != nil {
			LOGGER.Error("convertMacAddressToArrayByte, machine address is invalid", v)
		}
		macAddr[i] = byte(tmpn)
	}
	return
}

func convertIpv4AddressToString(v uint32) (ip4Addr string) {
	ip4Addr = fmt.Sprintf("%d.%d.%d.%d", byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
	return
}
func convertIpv4AddressToUint32(v string) (ip4Addr uint32) {
	a := strings.Split(v, ".")
	if len(a) != 4 {
		ip4Addr = 0
		LOGGER.Error("convertIpv4AddressToUint32, ip address is invalid", v)
		return
	}
	for i := 3; i >= 0; i-- {
		tmpn, err := strconv.ParseUint(a[i], 10, 8)
		if err != nil {
			LOGGER.Error("convertIpv4AddressToUint32, ip address is invalid", v)
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

// TODO
func convertIpv6AddressToString(v []byte) string {
	return ""
}
func convertIpv6AddressToArrayByte(v string) []byte {
	return nil
}
