package main

import (
	"fmt"
)

// Convert dbus variant to other data type

func interfaceToString(v interface{}) (d string, err error) {
	d, ok := v.(string)
	if !ok {
		err = fmt.Errorf("variantToString() failed: %v", v)
		return
	}
	return
}

func interfaceToByte(v interface{}) (d byte, err error) {
	d, ok := v.(byte)
	if !ok {
		err = fmt.Errorf("variantToByte() failed: %v", v)
		return
	}
	return
}

func interfaceToInt32(v interface{}) (d int32, err error) {
	d, ok := v.(int32)
	if !ok {
		err = fmt.Errorf("variantToInt32() failed: %v", v)
		return
	}
	return
}

func interfaceToUint32(v interface{}) (d uint32, err error) {
	d, ok := v.(uint32)
	if !ok {
		err = fmt.Errorf("variantToUint32() failed: %v", v)
		return
	}
	return
}

func interfaceToUint64(v interface{}) (d uint64, err error) {
	d, ok := v.(uint64)
	if !ok {
		err = fmt.Errorf("variantToUint64() failed: %v", v)
		return
	}
	return
}

func interfaceToBoolean(v interface{}) (d bool, err error) {
	d, ok := v.(bool)
	if !ok {
		err = fmt.Errorf("variantToBoolean() failed: %v", v)
		return
	}
	return
}

func interfaceToArrayByte(v interface{}) (d []byte, err error) {
	d, ok := v.([]byte)
	if !ok {
		err = fmt.Errorf("variantToArrayByte() failed: %v", v)
		return
	}
	return
}

func interfaceToArrayString(v interface{}) (d []string, err error) {
	d, ok := v.([]string)
	if !ok {
		err = fmt.Errorf("variantToArrayString() failed: %v", v)
		return
	}
	return
}

func interfaceToArrayUint32(v interface{}) (d []uint32, err error) {
	d, ok := v.([]uint32)
	if !ok {
		err = fmt.Errorf("variantToArrayUint32() failed: %v", v)
		return
	}
	return
}

func interfaceToArrayArrayByte(v interface{}) (d [][]byte, err error) {
	d, ok := v.([][]byte)
	if !ok {
		err = fmt.Errorf("variantToArrayArrayByte() failed: %v", v)
		return
	}
	return
}

func interfaceToArrayArrayUint32(v interface{}) (d [][]uint32, err error) {
	d, ok := v.([][]uint32)
	if !ok {
		err = fmt.Errorf("variantToArrayArrayUint32() failed: %v", v)
		return
	}
	return
}

func interfaceToDictStringString(v interface{}) (d map[string]string, err error) {
	d, ok := v.(map[string]string)
	if !ok {
		err = fmt.Errorf("variantToDictStringString() failed: %v", v)
		return
	}
	return
}

func interfaceToIpv6Addresses(v interface{}) (d Ipv6Addresses, err error) {
	d, ok := v.(Ipv6Addresses)
	if !ok {
		err = fmt.Errorf("variantToIpv6Addresses() failed: %v", v)
		return
	}
	return
}

func interfaceToIpv6Routes(v interface{}) (d Ipv6Routes, err error) {
	d, ok := v.(Ipv6Routes)
	if !ok {
		err = fmt.Errorf("variantToIpv6Routes() failed: %v", v)
		return
	}
	return
}
