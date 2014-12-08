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

package network

import (
	"encoding/json"
	"fmt"
	"pkg.linuxdeepin.com/lib/utils"
	"strings"
)

func isStringInArray(s string, list []string) bool {
	for _, i := range list {
		if i == s {
			return true
		}
	}
	return false
}

func stringArrayBut(list []string, ignoreList ...string) (newList []string) {
	for _, s := range list {
		if !isStringInArray(s, ignoreList) {
			newList = append(newList, s)
		}
	}
	return
}

func appendStrArrayUnique(a1 []string, a2 ...string) (a []string) {
	a = a1
	for _, s := range a2 {
		if !isStringInArray(s, a) {
			a = append(a, s)
		}
	}
	return
}

func isInterfaceNil(v interface{}) bool {
	return utils.IsInterfaceNil(v)
}

func isInterfaceEmpty(v interface{}) bool {
	if isInterfaceNil(v) {
		return true
	}
	switch v.(type) {
	case [][]interface{}: // ipv6Addresses
		if vd, ok := v.([][]interface{}); ok {
			if len(vd) == 0 {
				return true
			}
		}
	}
	return false
}

func marshalJSON(v interface{}) (jsonStr string, err error) {
	b, err := json.Marshal(v)
	if err != nil {
		logger.Error(err)
		return
	}
	jsonStr = string(b)
	return
}

func unmarshalJSON(jsonStr string, v interface{}) (err error) {
	err = json.Unmarshal([]byte(jsonStr), &v)
	if err != nil {
		logger.Error(err)
	}
	return
}

func isUint32ArrayEmpty(a []uint32) (empty bool) {
	empty = true
	for _, v := range a {
		if v != 0 {
			empty = false
			break
		}
	}
	return
}

// convert local path to uri, etc "/the/path" -> "file:///the/path"
func toUriPath(path string) (uriPath string) {
	return utils.EncodeURI(path, utils.SCHEME_FILE)
}

// convert uri to local path, etc "file:///the/path" -> "/the/path"
func toLocalPath(path string) (localPath string) {
	return utils.DecodeURI(path)
}

// convert local path to uri, etc "/the/path" -> "file:///the/path"
func toUriPathFor8021x(path string) (uriPath string) {
	// the uri for 8021x cert files is specially, we just need append
	// suffix "file://" for it
	if !utils.IsURI(path) {
		uriPath = "file://" + path
	} else {
		uriPath = path
	}
	return
}

// convert uri to local path, etc "file:///the/path" -> "/the/path"
func toLocalPathFor8021x(path string) (uriPath string) {
	// the uri for 8021x cert files is specially, we just need remove
	// suffix "file://" from it
	if utils.IsURI(path) {
		uriPath = strings.TrimPrefix(path, "file://")
	} else {
		uriPath = path
	}
	return
}

// byte array should end with null byte
func strToByteArrayPath(path string) (bytePath []byte) {
	bytePath = []byte(path)
	bytePath = append(bytePath, 0)
	return
}
func byteArrayToStrPath(bytePath []byte) (path string) {
	if len(bytePath) < 1 {
		return
	}
	path = string(bytePath[:len(bytePath)-1])
	return
}

// strToUuid convert any given string to md5, and then to uuid, for
// example, a device address string "00:12:34:56:ab:cd" will be
// converted to "1d417dad-8a98-fb90-e9df-016bd616d7dd"
func strToUuid(str string) (uuid string) {
	md5, _ := utils.SumStrMd5(str)
	return doStrToUuid(md5)
}
func doStrToUuid(str string) (uuid string) {
	str = strings.ToLower(str)
	for i := 0; i < len(str); i++ {
		if (str[i] >= '0' && str[i] <= '9') ||
			(str[i] >= 'a' && str[i] <= 'f') {
			uuid = uuid + string(str[i])
		}
	}
	if len(uuid) < 32 {
		misslen := 32 - len(uuid)
		uuid = strings.Repeat("0", misslen) + uuid
	}
	uuid = fmt.Sprintf("%s-%s-%s-%s-%s", uuid[0:8], uuid[8:12], uuid[12:16], uuid[16:20], uuid[20:32])
	return
}
