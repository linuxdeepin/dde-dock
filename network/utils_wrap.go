package main

import (
	"fmt"
	"strconv"
	"strings"
)

// []byte{0,0,0,0,0,0} -> "00:00:00:00:00:00"
func formatMacAddressToString(v []byte) (macAddr string, err error) {
	if len(v) != 6 {
		macAddr = "00:00:00:00:00:00"
		err = fmt.Errorf("formatMacAddressToString error, machine address is invalid: %v", v)
		return
	}
	macAddr = fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", v[0], v[1], v[2], v[3], v[4], v[5])
	return
}

// "00:00:00:00:00:00" -> []byte{0,0,0,0,0,0}
func formatMacAddressToArrayByte(v string) (macAddr []byte, err error) {
	a := strings.Split(v, ":")
	if len(a) != 6 {
		macAddr = []byte{0, 0, 0, 0, 0, 0}
		err = fmt.Errorf("formatMacAddressToArrayByte error, machine address is invalid: %v", v)
		return
	}
	macAddr = make([]byte, 6)
	var tmpn uint64
	for i := 0; i < 6; i++ {
		tmpn, err = strconv.ParseUint(a[i], 16, 8) // TODO
		macAddr[i] = byte(tmpn)
	}
	return
}

func formatIpv4AddressToString(v uint32) (ip4Addr string) {
	ip4Addr = fmt.Sprintf("%d.%d.%d.%d", byte(v), byte(v>>8), byte(v>>16), byte(v>>24))
	return
}
func formatIpv4AddressToUint32(v string) (ip4Addr uint32, err error) {
	a := strings.Split(v, ".")
	if len(a) != 4 {
		ip4Addr = 0
		err = fmt.Errorf("formatIpv4AddressToUint32 error, ip address is invalid: %v", v)
		return
	}
	var tmpn uint64
	for i := 3; i >= 0; i-- {
		tmpn, err = strconv.ParseUint(a[i], 10, 8)
		ip4Addr = ip4Addr<<8 + uint32(tmpn)
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

func wrapIpv4Dns(v []uint32) (addrs []string) {
	for _, a := range v {
		addrs = append(addrs, formatIpv4AddressToString(a))
	}
	return
}
func unwrapIpv4Dns(v []string) (addrs []uint32) {
	for _, a := range v {
		addr, _ := formatIpv4AddressToUint32(a)
		addrs = append(addrs, addr)
	}
	return
}

func wrapIpv4Addresses(v [][]uint32) Ipv4AddressesWrapper {
	// TODO
	return nil
}
func unwrapIpv4Addresses(v Ipv4AddressesWrapper) [][]uint32 {
	// TODO
	return nil
}

func wrapIpv4Routes(v [][]uint32) Ipv4RoutesWrapper {
	// TODO
	return nil
}
func unwrapIpv4Routes(v Ipv4RoutesWrapper) [][]uint32 {
	// TODO
	return nil
}

func wrapIpv6Dns(v [][]byte) (addrs []string) {
	for _, a := range v {
		addrs = append(addrs, formatIpv6AddressToString(a))
	}
	return
}
func unwrapIpv6Dns(v []string) (addrs [][]byte) {
	for _, a := range v {
		addrs = append(addrs, formatIpv6AddressToArrayByte(a))
	}
	return
}

func wrapIpv6Addresses(v Ipv6Addresses) Ipv6AddressesWrapper {
	// TODO
	return nil
}
func unwrapIpv6Addresses(v Ipv6AddressesWrapper) Ipv6Addresses {
	// TODO
	return nil
}

func wrapIpv6Routes(v Ipv6Routes) Ipv6RoutesWrapper {
	// TODO
	return nil
}
func unwrapIpv6Routes(v Ipv6RoutesWrapper) Ipv6Routes {
	// TODO
	return nil
}
