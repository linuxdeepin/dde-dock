/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package bluetooth

import (
	apidevice "dbus/com/deepin/api/device"
	"encoding/json"
)

func isStringInArray(str string, list []string) bool {
	for _, tmp := range list {
		if tmp == str {
			return true
		}
	}
	return false
}

func marshalJSON(v interface{}) (strJSON string) {
	byteJSON, err := json.Marshal(v)
	if err != nil {
		logger.Error(err)
		return
	}
	strJSON = string(byteJSON)
	return
}

func unmarshalJSON(strJSON string, v interface{}) {
	err := json.Unmarshal([]byte(strJSON), v)
	if err != nil {
		logger.Error(err)
	}
	return
}

func isDBusObjectKeyExists(data dbusObjectData, key string) (ok bool) {
	_, ok = data[key]
	return
}

func getDBusObjectValueString(data dbusObjectData, key string) (r string) {
	v, ok := data[key]
	if ok {
		r = interfaceToString(v.Value())
	}
	return
}

func getDBusObjectValueInt16(data dbusObjectData, key string) (r int16) {
	v, ok := data[key]
	if ok {
		r = interfaceToInt16(v.Value())
	}
	return
}

func interfaceToString(v interface{}) (r string) {
	r, _ = v.(string)
	return
}

func interfaceToInt16(v interface{}) (r int16) {
	r, _ = v.(int16)
	return
}

func requestUnblockBluetoothDevice() {
	d, err := apidevice.NewDevice("com.deepin.api.Device", "/com/deepin/api/Device")
	if err != nil {
		logger.Error(err)
		return
	}
	err = d.UnblockDevice("bluetooth")
	if err != nil {
		logger.Error(err)
	}
}
