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
	"crypto/rand"
	"dbus/org/freedesktop/notifications"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

func genUuid() string {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		panic("This can failed?")
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

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

func randString(n int) string {
	const alphanum = "0123456789abcdef"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func isInterfaceNil(v interface{}) bool {
	defer func() { recover() }()
	return v == nil || reflect.ValueOf(v).IsNil()
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

func unmarshalJSON(jsonStr string) (v interface{}, err error) {
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

func isUriPath(path string) bool {
	if strings.HasPrefix(path, "file://") {
		return true
	}
	return false
}

func isLocalPath(path string) bool {
	if isUriPath(path) {
		return false
	}
	return true
}

// "/the/path" -> "file:///the/path", "file:///the/path" -> "file:///the/path"
func toUriPath(path string) (uriPath string) {
	if strings.HasPrefix(path, "file://") {
		uriPath = path
	} else {
		uriPath = "file://" + path
	}
	return
}

// "/the/path" -> "/the/path", "file:///the/path" -> "/the/path"
func toLocalPath(path string) (localPath string) {
	localPath = strings.TrimPrefix(path, "file://")
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

func isFileExists(file string) bool {
	stat, err := os.Stat(file)
	if err == nil && !stat.IsDir() {
		return true
	}
	return false
}

func notify(icon, summary, body string) (err error) {
	var notify *notifications.Notifier
	notify, err = notifications.NewNotifier("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	if err != nil {
		logger.Error(err)
		return
	}
	appName := "Network"
	_, err = notify.Notify(appName, 0, icon, summary, body, nil, nil, 0)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}
