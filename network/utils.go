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
	"dbus/org/freedesktop/notifications"
	"encoding/json"
	"pkg.linuxdeepin.com/lib/utils"
)

const (
	dbusNotifyDest = "org.freedesktop.Notifications"
	dbusNotifyPath = "/org/freedesktop/Notifications"
)

var notifier, _ = notifications.NewNotifier(dbusNotifyDest, dbusNotifyPath)

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

// "/the/path" -> "file:///the/path", "file:///the/path" -> "file:///the/path"
func toUriPath(path string) (uriPath string) {
	return utils.EncodeURI(path, utils.SCHEME_FILE)
}

// "/the/path" -> "/the/path", "file:///the/path" -> "/the/path"
func toLocalPath(path string) (localPath string) {
	return utils.DecodeURI(path)
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

func notify(icon, summary, body string) (err error) {
	if notifier == nil {
		logger.Error("connect to org.freedesktop.Notifications failed")
		return
	}
	_, err = notifier.Notify("Network", 0, icon, summary, body, nil, nil, 0)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}
