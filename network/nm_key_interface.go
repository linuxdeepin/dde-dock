package main

import (
	"fmt"
)

// Convert dbus variant's value to other data type

func interfaceToString(v interface{}) (d string, err error) {
	d, ok := v.(string)
	if !ok {
		err = fmt.Errorf("interfaceToString() failed: %v", v)
		return
	}
	return
}

func interfaceToByte(v interface{}) (d byte, err error) {
	d, ok := v.(byte)
	if !ok {
		err = fmt.Errorf("interfaceToByte() failed: %v", v)
		return
	}
	return
}

func interfaceToInt32(v interface{}) (d int32, err error) {
	d, ok := v.(int32)
	if !ok {
		err = fmt.Errorf("interfaceToInt32() failed: %v", v)
		return
	}
	return
}

func interfaceToUint32(v interface{}) (d uint32, err error) {
	d, ok := v.(uint32)
	if !ok {
		err = fmt.Errorf("interfaceToUint32() failed: %v", v)
		return
	}
	return
}

func interfaceToUint64(v interface{}) (d uint64, err error) {
	d, ok := v.(uint64)
	if !ok {
		err = fmt.Errorf("interfaceToUint64() failed: %v", v)
		return
	}
	return
}

func interfaceToBoolean(v interface{}) (d bool, err error) {
	d, ok := v.(bool)
	if !ok {
		err = fmt.Errorf("interfaceToBoolean() failed: %v", v)
		return
	}
	return
}

func interfaceToArrayByte(v interface{}) (d []byte, err error) {
	d, ok := v.([]byte)
	if !ok {
		err = fmt.Errorf("interfaceToArrayByte() failed: %v", v)
		return
	}
	return
}

func interfaceToArrayString(v interface{}) (d []string, err error) {
	d, ok := v.([]string)
	if !ok {
		err = fmt.Errorf("interfaceToArrayString() failed: %v", v)
		return
	}
	return
}

func interfaceToArrayUint32(v interface{}) (d []uint32, err error) {
	d, ok := v.([]uint32)
	if !ok {
		err = fmt.Errorf("interfaceToArrayUint32() failed: %v", v)
		return
	}
	return
}

func interfaceToArrayArrayByte(v interface{}) (d [][]byte, err error) {
	d, ok := v.([][]byte)
	if !ok {
		err = fmt.Errorf("interfaceToArrayArrayByte() failed: %v", v)
		return
	}
	return
}

func interfaceToArrayArrayUint32(v interface{}) (d [][]uint32, err error) {
	d, ok := v.([][]uint32)
	if !ok {
		err = fmt.Errorf("interfaceToArrayArrayUint32() failed: %v", v)
		return
	}
	return
}

func interfaceToDictStringString(v interface{}) (d map[string]string, err error) {
	d, ok := v.(map[string]string)
	if !ok {
		err = fmt.Errorf("interfaceToDictStringString() failed: %v", v)
		return
	}
	return
}

func interfaceToIpv6Addresses(v interface{}) (d ipv6Addresses, err error) {
	d, ok := v.(ipv6Addresses)
	if !ok {
		err = fmt.Errorf("interfaceToIpv6Addresses() failed: %v", v)
		return
	}
	return
}

func interfaceToIpv6Routes(v interface{}) (d ipv6Routes, err error) {
	d, ok := v.(ipv6Routes)
	if !ok {
		err = fmt.Errorf("interfaceToIpv6Routes() failed: %v", v)
		return
	}
	return
}
