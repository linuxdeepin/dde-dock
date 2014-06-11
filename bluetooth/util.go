/**
 * Copyright (c) 2014 Deepin, Inc.
 *               2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package bluetooth

import (
	"encoding/json"
)

func isStringInArray(a string, list []string) bool {
	for _, b := range list {
		if b == a {
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
