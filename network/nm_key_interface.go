/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package network

// Convert dbus variant's value to other data type

func interfaceToString(v interface{}) (d string) {
	if isInterfaceNil(v) {
		return
	}
	d, ok := v.(string)
	if !ok {
		logger.Errorf("interfaceToString() failed: %#v", v)
		return
	}
	return
}

func interfaceToByte(v interface{}) (d byte) {
	if isInterfaceNil(v) {
		return
	}
	d, ok := v.(byte)
	if !ok {
		logger.Errorf("interfaceToByte() failed: %#v", v)
		return
	}
	return
}

func interfaceToInt32(v interface{}) (d int32) {
	if isInterfaceNil(v) {
		return
	}
	d, ok := v.(int32)
	if !ok {
		logger.Errorf("interfaceToInt32() failed: %#v", v)
		return
	}
	return
}

func interfaceToUint32(v interface{}) (d uint32) {
	if isInterfaceNil(v) {
		return
	}
	d, ok := v.(uint32)
	if !ok {
		logger.Errorf("interfaceToUint32() failed: %#v", v)
		return
	}
	return
}

func interfaceToInt64(v interface{}) (d int64) {
	if isInterfaceNil(v) {
		return
	}
	d, ok := v.(int64)
	if !ok {
		logger.Errorf("interfaceToInt64() failed: %#v", v)
		return
	}
	return
}

func interfaceToUint64(v interface{}) (d uint64) {
	if isInterfaceNil(v) {
		return
	}
	d, ok := v.(uint64)
	if !ok {
		logger.Errorf("interfaceToUint64() failed: %#v", v)
		return
	}
	return
}

func interfaceToBoolean(v interface{}) (d bool) {
	if isInterfaceNil(v) {
		return
	}
	d, ok := v.(bool)
	if !ok {
		logger.Errorf("interfaceToBoolean() failed: %#v", v)
		return
	}
	return
}

func interfaceToArrayByte(v interface{}) (d []byte) {
	if isInterfaceNil(v) {
		return
	}
	d, ok := v.([]byte)
	if !ok {
		logger.Errorf("interfaceToArrayByte() failed: %#v", v)
		return
	}
	return
}

func interfaceToArrayString(v interface{}) (d []string) {
	if isInterfaceNil(v) {
		return
	}
	d, ok := v.([]string)
	if !ok {
		logger.Errorf("interfaceToArrayString() failed: %#v", v)
		return
	}
	return
}

func interfaceToArrayUint32(v interface{}) (d []uint32) {
	if isInterfaceNil(v) {
		return
	}
	d, ok := v.([]uint32)
	if !ok {
		logger.Errorf("interfaceToArrayUint32() failed: %#v", v)
		return
	}
	return
}

func interfaceToArrayArrayByte(v interface{}) (d [][]byte) {
	if isInterfaceNil(v) {
		return
	}
	d, ok := v.([][]byte)
	if !ok {
		logger.Errorf("interfaceToArrayArrayByte() failed: %#v", v)
		return
	}
	return
}

func interfaceToArrayArrayUint32(v interface{}) (d [][]uint32) {
	if isInterfaceNil(v) {
		return
	}
	d, ok := v.([][]uint32)
	if !ok {
		logger.Errorf("interfaceToArrayArrayUint32() failed: %#v", v)
		return
	}
	return
}

func interfaceToDictStringString(v interface{}) (d map[string]string) {
	if isInterfaceNil(v) {
		return
	}
	d, ok := v.(map[string]string)
	if !ok {
		logger.Errorf("interfaceToDictStringString() failed: %#v", v)
		return
	}
	return
}

func interfaceToIpv6Addresses(v interface{}) (d ipv6Addresses) {
	if isInterfaceNil(v) {
		return
	}

	// try convert interface to [][]interface{} and ipv6Addresses
	tmpData, ok := v.([][]interface{})
	if !ok {
		d, ok = v.(ipv6Addresses)
		if !ok {
			logger.Errorf("interfaceToIpv6Addresses() failed: %#v", v)
		}
		return
	}
	d = make(ipv6Addresses, len(tmpData))
	for i, _ := range tmpData {
		if len(tmpData[i]) >= 3 {
			var ok0, ok1, ok2 bool
			d[i].Address, ok0 = tmpData[i][0].([]byte)
			d[i].Prefix, ok1 = tmpData[i][1].(uint32)
			d[i].Gateway, ok2 = tmpData[i][2].([]byte)
			if !(ok0 && ok1 && ok2) {
				logger.Errorf("interfaceToIpv6Addresses() failed: %#v", v)
				return
			}
		}
	}
	return
}

func interfaceToIpv6Routes(v interface{}) (d ipv6Routes) {
	if isInterfaceNil(v) {
		return
	}

	// try convert interface to [][]interface{} and ipv6Routes
	tmpData, ok := v.([][]interface{})
	if !ok {
		d, ok = v.(ipv6Routes)
		if !ok {
			logger.Errorf("interfaceToIpv6Routes() failed: %#v", v)
		}
		return
	}
	d = make(ipv6Routes, len(tmpData))
	for i, _ := range tmpData {
		if len(tmpData) >= 4 {
			var ok0, ok1, ok2, ok3 bool
			d[i].Address, ok0 = tmpData[i][0].([]byte)
			d[i].Prefix, ok1 = tmpData[i][1].(uint32)
			d[i].NextHop, ok2 = tmpData[i][2].([]byte)
			d[i].Metric, ok3 = tmpData[i][3].(uint32)
			if !(ok0 && ok1 && ok2 && ok3) {
				logger.Errorf("interfaceToIpv6Routes() failed: %#v", v)
				return
			}
		}
	}
	return
}
