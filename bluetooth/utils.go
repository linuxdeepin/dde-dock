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
