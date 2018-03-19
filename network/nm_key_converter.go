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
	"encoding/json"
	"fmt"
)

const (
	jsonNull        = `null`
	jsonEmptyString = `""`
	jsonEmptyArray  = `[]`
)

// dbus.Variant.Value() -> realdata -> wrapped data(if need) -> json string
func keyValueToJSON(v interface{}, t ktype) (jsonStr string, err error) {
	// dispatch wrapper keys
	switch t {
	case ktypeWrapperString:
		tmpv := interfaceToArrayByte(v)
		v = string(tmpv)
	case ktypeWrapperMacAddress:
		tmpv := interfaceToArrayByte(v)
		v = convertMacAddressToString(tmpv)
	case ktypeWrapperIpv4Dns:
		tmpv := interfaceToArrayUint32(v)
		v = wrapIpv4Dns(tmpv)
	case ktypeWrapperIpv4Addresses:
		tmpv := interfaceToArrayArrayUint32(v)
		v = wrapIpv4Addresses(tmpv)
	case ktypeWrapperIpv4Routes:
		tmpv := interfaceToArrayArrayUint32(v)
		v = wrapIpv4Routes(tmpv)
	case ktypeWrapperIpv6Dns:
		tmpv := interfaceToArrayArrayByte(v)
		v = wrapIpv6Dns(tmpv)
	case ktypeWrapperIpv6Addresses:
		tmpv := interfaceToIpv6Addresses(v)
		v = wrapIpv6Addresses(tmpv)
	case ktypeWrapperIpv6Routes:
		tmpv := interfaceToIpv6Routes(v)
		v = wrapIpv6Routes(tmpv)
	}

	jsonStr, err = marshalJSON(v)
	return
}

// json string -> wrapped data(if need) -> realdata -> dbus.Variant.Value()
func jsonToKeyValue(jsonStr string, t ktype) (v interface{}, err error) {
	switch t {
	default:
		err = fmt.Errorf("invalid variant type, %s", jsonStr)
	case ktypeString:
		v, err = jsonToKeyValueString(jsonStr)
	case ktypeByte:
		v, err = jsonToKeyValueByte(jsonStr)
	case ktypeInt32:
		v, err = jsonToKeyValueInt32(jsonStr)
	case ktypeUint32:
		v, err = jsonToKeyValueUint32(jsonStr)
	case ktypeUint64:
		v, err = jsonToKeyValueUint64(jsonStr)
	case ktypeBoolean:
		v, err = jsonToKeyValueBoolean(jsonStr)
	case ktypeArrayString:
		v, err = jsonToKeyValueArrayString(jsonStr)
	case ktypeArrayByte:
		v, err = jsonToKeyValueArrayByte(jsonStr)
	case ktypeArrayUint32:
		v, err = jsonToKeyValueArrayUint32(jsonStr)
	case ktypeArrayArrayByte:
		v, err = jsonToKeyValueArrayArrayByte(jsonStr)
	case ktypeArrayArrayUint32:
		v, err = jsonToKeyValueArrayArrayUint32(jsonStr)
	case ktypeDictStringString:
		v, err = jsonToKeyValueDictStringString(jsonStr)
	case ktypeIpv6Addresses:
		v, err = jsonToKeyValueIpv6Addresses(jsonStr)
	case ktypeIpv6Routes:
		v, err = jsonToKeyValueIpv6Routes(jsonStr)
	case ktypeWrapperString:
		v, err = jsonToKeyValueWrapperString(jsonStr)
	case ktypeWrapperMacAddress:
		v, err = jsonToKeyValueWrapperMacAddress(jsonStr)
	case ktypeWrapperIpv4Dns:
		v, err = jsonToKeyValueWrapperIpv4Dns(jsonStr)
	case ktypeWrapperIpv4Addresses:
		v, err = jsonToKeyValueWrapperIpv4Addresses(jsonStr)
	case ktypeWrapperIpv4Routes:
		v, err = jsonToKeyValueWrapperIpv4Routes(jsonStr)
	case ktypeWrapperIpv6Dns:
		v, err = jsonToKeyValueWrapperIpv6Dns(jsonStr)
	case ktypeWrapperIpv6Addresses:
		v, err = jsonToKeyValueWrapperIpv6Addresses(jsonStr)
	case ktypeWrapperIpv6Routes:
		v, err = jsonToKeyValueWrapperIpv6Routes(jsonStr)
	}
	return
}

// Convert sepcial key type which wrapped by json to dbus variant value
func jsonToKeyValueString(jsonStr string) (v string, err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	return
}
func jsonToKeyValueByte(jsonStr string) (v byte, err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	return
}
func jsonToKeyValueInt32(jsonStr string) (v int32, err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	return
}
func jsonToKeyValueUint32(jsonStr string) (v uint32, err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	return
}
func jsonToKeyValueUint64(jsonStr string) (v uint64, err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	return
}
func jsonToKeyValueBoolean(jsonStr string) (v bool, err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	return
}
func jsonToKeyValueArrayByte(jsonStr string) (v []byte, err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	return
}
func jsonToKeyValueArrayString(jsonStr string) (v []string, err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	return
}
func jsonToKeyValueArrayUint32(jsonStr string) (v []uint32, err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	return
}
func jsonToKeyValueArrayArrayByte(jsonStr string) (v [][]byte, err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	return
}
func jsonToKeyValueArrayArrayUint32(jsonStr string) (v [][]uint32, err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	return
}
func jsonToKeyValueDictStringString(jsonStr string) (v map[string]string, err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	return
}
func jsonToKeyValueIpv6Addresses(jsonStr string) (v ipv6Addresses, err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	return
}
func jsonToKeyValueIpv6Routes(jsonStr string) (v ipv6Routes, err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	return
}

// key type wrapper
func jsonToKeyValueWrapperString(jsonStr string) (v []byte, err error) {
	// wrap ktypeArrayByte to [string]
	var wrapData string
	err = json.Unmarshal([]byte(jsonStr), &wrapData)
	if err != nil {
		return
	}
	v = []byte(wrapData)
	return
}
func jsonToKeyValueWrapperMacAddress(jsonStr string) (v []byte, err error) {
	// wrap ktypeArrayByte to [string]
	var wrapData string
	err = json.Unmarshal([]byte(jsonStr), &wrapData)
	if err != nil {
		return
	}
	v, err = convertMacAddressToArrayByteCheck(wrapData)
	return
}
func jsonToKeyValueWrapperIpv4Dns(jsonStr string) (v []uint32, err error) {
	// wrap ktypeArrayUint32 to [array of string]
	var wrapData []string
	err = json.Unmarshal([]byte(jsonStr), &wrapData)
	if err != nil {
		return
	}
	v = unwrapIpv4Dns(wrapData)
	return
}
func jsonToKeyValueWrapperIpv4Addresses(jsonStr string) (v [][]uint32, err error) {
	// wrap ktypeArrayArrayUint32 to [array of (string, uint32, string)]
	var wrapData ipv4AddressesWrapper
	err = json.Unmarshal([]byte(jsonStr), &wrapData)
	if err != nil {
		return
	}
	v = unwrapIpv4Addresses(wrapData)
	return
}
func jsonToKeyValueWrapperIpv4Routes(jsonStr string) (v [][]uint32, err error) {
	// wrap ktypeArrayArrayUint32 to [array of (string, uint32, string, uint32)]
	var wrapData ipv4RoutesWrapper
	err = json.Unmarshal([]byte(jsonStr), &wrapData)
	if err != nil {
		return
	}
	v = unwrapIpv4Routes(wrapData)
	return
}
func jsonToKeyValueWrapperIpv6Dns(jsonStr string) (v [][]byte, err error) {
	// wrap ktypeArrayArrayByte to [array of string]
	var wrapData []string
	err = json.Unmarshal([]byte(jsonStr), &wrapData)
	if err != nil {
		return
	}
	v = unwrapIpv6Dns(wrapData)
	return
}
func jsonToKeyValueWrapperIpv6Addresses(jsonStr string) (v ipv6Addresses, err error) {
	// wrap ktypeIpv6Addresses to [array of (string, uint32, string)]
	var wrapData ipv6AddressesWrapper
	err = json.Unmarshal([]byte(jsonStr), &wrapData)
	if err != nil {
		return
	}
	v = unwrapIpv6Addresses(wrapData)
	return
}
func jsonToKeyValueWrapperIpv6Routes(jsonStr string) (v ipv6Routes, err error) {
	// wrap ktypeIpv6Routes to [array of (string, uint32, string, uint32)]
	var wrapData ipv6RoutesWrapper
	err = json.Unmarshal([]byte(jsonStr), &wrapData)
	if err != nil {
		return
	}
	v = unwrapIpv6Routes(wrapData)
	return
}

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

// Wrappers

func wrapIpv4Dns(data []uint32) (wrapData []string) {
	for _, a := range data {
		wrapData = append(wrapData, convertIpv4AddressToString(a))
	}
	return
}
func unwrapIpv4Dns(wrapData []string) (data []uint32) {
	for _, a := range wrapData {
		data = append(data, convertIpv4AddressToUint32(a))
	}
	return
}

func wrapIpv4Addresses(data [][]uint32) (wrapData ipv4AddressesWrapper) {
	for _, d := range data {
		if len(d) != 3 {
			logger.Error("ipv4 address invalid", d)
			continue
		}
		ipv4Addr := ipv4AddressWrapper{}
		ipv4Addr.Address = convertIpv4AddressToString(d[0])
		ipv4Addr.Mask = convertIpv4PrefixToNetMask(d[1])
		ipv4Addr.Gateway = convertIpv4AddressToString(d[2])
		wrapData = append(wrapData, ipv4Addr)
	}
	return
}
func unwrapIpv4Addresses(wrapData ipv4AddressesWrapper) (data [][]uint32) {
	for _, d := range wrapData {
		ipv4Addr := make([]uint32, 3)
		ipv4Addr[0] = convertIpv4AddressToUint32(d.Address)
		ipv4Addr[1] = convertIpv4NetMaskToPrefix(d.Mask)
		ipv4Addr[2] = convertIpv4AddressToUint32(d.Gateway)
		data = append(data, ipv4Addr)
	}
	return
}

func wrapIpv4Routes(data [][]uint32) (wrapData ipv4RoutesWrapper) {
	for _, d := range data {
		if len(d) != 4 {
			logger.Error("invalid ipv4 route", d)
			continue
		}
		ipv4Route := ipv4RouteWrapper{}
		ipv4Route.Address = convertIpv4AddressToString(d[0])
		ipv4Route.Mask = convertIpv4PrefixToNetMask(d[1])
		ipv4Route.NextHop = convertIpv4AddressToString(d[2])
		ipv4Route.Metric = d[3]
		wrapData = append(wrapData, ipv4Route)
	}
	return
}
func unwrapIpv4Routes(wrapData ipv4RoutesWrapper) (data [][]uint32) {
	for _, d := range wrapData {
		ipv4Route := make([]uint32, 4)
		ipv4Route[0] = convertIpv4AddressToUint32(d.Address)
		ipv4Route[1] = convertIpv4NetMaskToPrefix(d.Mask)
		ipv4Route[2] = convertIpv4AddressToUint32(d.NextHop)
		ipv4Route[3] = d.Metric
		data = append(data, ipv4Route)
	}
	return
}

func wrapIpv6Dns(data [][]byte) (wrapData []string) {
	for _, a := range data {
		wrapData = append(wrapData, convertIpv6AddressToString(a))
	}
	return
}
func unwrapIpv6Dns(wrapData []string) (data [][]byte) {
	for _, a := range wrapData {
		data = append(data, convertIpv6AddressToArrayByte(a))
	}
	return
}

func wrapIpv6Addresses(data ipv6Addresses) (wrapData ipv6AddressesWrapper) {
	for _, d := range data {
		ipv6Addr := ipv6AddressWrapper{}
		ipv6Addr.Address = convertIpv6AddressToString(d.Address)
		ipv6Addr.Prefix = d.Prefix
		ipv6Addr.Gateway = convertIpv6AddressToString(d.Gateway)
		wrapData = append(wrapData, ipv6Addr)
	}
	return
}

func unwrapIpv6Addresses(wrapData ipv6AddressesWrapper) (data ipv6Addresses) {
	for _, d := range wrapData {
		ipv6Addr := ipv6Address{}
		ipv6Addr.Address = convertIpv6AddressToArrayByte(d.Address)
		ipv6Addr.Prefix = d.Prefix
		ipv6Addr.Gateway = convertIpv6AddressToArrayByte(d.Gateway)
		data = append(data, ipv6Addr)
	}
	return
}

func wrapIpv6Routes(data ipv6Routes) (wrapData ipv6RoutesWrapper) {
	for _, d := range data {
		ipv6Route := ipv6RouteWrapper{}
		ipv6Route.Address = convertIpv6AddressToString(d.Address)
		ipv6Route.Prefix = d.Prefix
		ipv6Route.NextHop = convertIpv6AddressToString(d.NextHop)
		ipv6Route.Metric = d.Metric
		wrapData = append(wrapData, ipv6Route)
	}
	return
}
func unwrapIpv6Routes(wrapData ipv6RoutesWrapper) (data ipv6Routes) {
	for _, d := range wrapData {
		ipv6Route := ipv6Route{}
		ipv6Route.Address = convertIpv6AddressToArrayByte(d.Address)
		ipv6Route.Prefix = d.Prefix
		ipv6Route.NextHop = convertIpv6AddressToArrayByte(d.NextHop)
		ipv6Route.Metric = d.Metric
		data = append(data, ipv6Route)
	}
	return
}
