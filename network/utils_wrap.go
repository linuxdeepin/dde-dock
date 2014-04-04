package main

import (
	"fmt"
	"strconv"
	"strings"
)

// []byte{0,0,0,0,0,0} -> "00:00:00:00:00:00"
func formatMacAddressToString(v []byte) (macAddr string) {
	if len(v) != 6 {
		macAddr = "00:00:00:00:00:00"
		LOGGER.Error("formatMacAddressToString, machine address is invalid", v)
		return
	}
	macAddr = fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", v[0], v[1], v[2], v[3], v[4], v[5])
	return
}

// "00:00:00:00:00:00" -> []byte{0,0,0,0,0,0}
func formatMacAddressToArrayByte(v string) (macAddr []byte) {
	a := strings.Split(v, ":")
	if len(a) != 6 {
		macAddr = []byte{0, 0, 0, 0, 0, 0}
		LOGGER.Error("formatMacAddressToArrayByte, machine address is invalid", v)
		return
	}
	macAddr = make([]byte, 6)
	for i := 0; i < 6; i++ {
		tmpn, err := strconv.ParseUint(a[i], 16, 8)
		if err != nil {
			LOGGER.Error("formatMacAddressToArrayByte, machine address is invalid", v)
		}
		macAddr[i] = byte(tmpn)
	}
	return
}

func formatIpv4AddressToString(v uint32) (ip4Addr string) {
	ip4Addr = fmt.Sprintf("%d.%d.%d.%d", byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
	return
}
func formatIpv4AddressToUint32(v string) (ip4Addr uint32) {
	a := strings.Split(v, ".")
	if len(a) != 4 {
		ip4Addr = 0
		LOGGER.Error("formatIpv4AddressToUint32, ip address is invalid", v)
		return
	}
	for i := 3; i >= 0; i-- {
		tmpn, err := strconv.ParseUint(a[i], 10, 8)
		if err != nil {
			LOGGER.Error("formatIpv4AddressToUint32, ip address is invalid", v)
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
	maskAddress = formatIpv4AddressToString(mask)
	return
}

// "255.255.255.0"(string format) -> 24
func convertIpv4NetMaskToPrefix(maskAddress string) (prefix uint32) {
	var mask uint32 // network order
	mask = formatIpv4AddressToUint32(maskAddress)
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
func formatIpv6AddressToString(v []byte) string {
	return ""
}
func formatIpv6AddressToArrayByte(v string) []byte {
	return nil
}

func wrapIpv4Dns(data []uint32) (wrapData []string) {
	for _, a := range data {
		wrapData = append(wrapData, formatIpv4AddressToString(a))
	}
	return
}
func unwrapIpv4Dns(wrapData []string) (data []uint32) {
	for _, a := range wrapData {
		data = append(data, formatIpv4AddressToUint32(a))
	}
	return
}

func wrapIpv4Addresses(data [][]uint32) (wrapData Ipv4AddressesWrapper) {
	for _, d := range data {
		if len(d) != 3 {
			LOGGER.Error("wrapIpv4Addresses, ipv4 address invalid", d)
			continue
		}
		ipv4Addr := Ipv4AddressWrapper{}
		ipv4Addr.Address = formatIpv4AddressToString(d[0])
		ipv4Addr.Prefix = d[1]
		ipv4Addr.Gateway = formatIpv4AddressToString(d[2])
		wrapData = append(wrapData, ipv4Addr)
	}
	return
}
func unwrapIpv4Addresses(wrapData Ipv4AddressesWrapper) (data [][]uint32) {
	for _, d := range wrapData {
		ipv4Addr := make([]uint32, 3)
		ipv4Addr[0] = formatIpv4AddressToUint32(d.Address)
		ipv4Addr[1] = d.Prefix
		ipv4Addr[2] = formatIpv4AddressToUint32(d.Gateway)
		data = append(data, ipv4Addr)
	}
	return
}

func wrapIpv4Routes(data [][]uint32) (wrapData Ipv4RoutesWrapper) {
	for _, d := range data {
		if len(d) != 4 {
			LOGGER.Error("wrapIpv4Routes: invalid ipv2 route", d)
			continue
		}
		ipv4Route := Ipv4RouteWrapper{}
		ipv4Route.Address = formatIpv4AddressToString(d[0])
		ipv4Route.Prefix = d[1]
		ipv4Route.NextHop = formatIpv4AddressToString(d[2])
		ipv4Route.Metric = d[3]
		wrapData = append(wrapData, ipv4Route)
	}
	return
}
func unwrapIpv4Routes(wrapData Ipv4RoutesWrapper) (data [][]uint32) {
	for _, d := range wrapData {
		ipv4Route := make([]uint32, 4)
		ipv4Route[0] = formatIpv4AddressToUint32(d.Address)
		ipv4Route[1] = d.Prefix
		ipv4Route[2] = formatIpv4AddressToUint32(d.NextHop)
		ipv4Route[3] = d.Metric
		data = append(data, ipv4Route)
	}
	return
}

func wrapIpv6Dns(data [][]byte) (wrapData []string) {
	for _, a := range data {
		wrapData = append(wrapData, formatIpv6AddressToString(a))
	}
	return
}
func unwrapIpv6Dns(wrapData []string) (data [][]byte) {
	for _, a := range wrapData {
		data = append(data, formatIpv6AddressToArrayByte(a))
	}
	return
}

func wrapIpv6Addresses(data Ipv6Addresses) (wrapData Ipv6AddressesWrapper) {
	// TODO
	return
}
func unwrapIpv6Addresses(wrapData Ipv6AddressesWrapper) (data Ipv6Addresses) {
	// TODO
	return
}

func wrapIpv6Routes(data Ipv6Routes) (wrapData Ipv6RoutesWrapper) {
	// TODO
	return
}
func unwrapIpv6Routes(wrapData Ipv6RoutesWrapper) (data Ipv6Routes) {
	// TODO
	return
}
