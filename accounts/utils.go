/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
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

package main

// #cgo CFLAGS: -Wall -g
// #cgo LDFLAGS: -lcrypt
// #include <stdlib.h>
// #include "mkpasswd.h"
import "C"

import (
	libpolkit1 "dbus/org/freedesktop/policykit1"
	"dlib/dbus"
	"os"
	"os/exec"
	"strings"
	"sync"
	"unsafe"
)

type polkitSubject struct {
	/*
	 * The following kinds of subjects are known:
	 * Unix Process: should be set to unix-process with keys
	 *                  pid (of type uint32) and
	 *                  start-time (of type uint64)
	 * Unix Session: should be set to unix-session with the key
	 *                  session-id (of type string)
	 * System Bus Name: should be set to system-bus-name with the key
	 *                  name (of type string)
	 */
	SubjectKind    string
	SubjectDetails map[string]dbus.Variant
}

var (
	mutex = new(sync.Mutex)
	genId = func() func() uint32 {
		id := uint32(0)
		return func() uint32 {
			mutex.Lock()
			tmp := id
			id += 1
			mutex.Unlock()
			return tmp
		}
	}()
)

const (
	POLKIT_DEST = "org.freedesktop.PolicyKit1"
	POLKIT_PATH = "/org/freedesktop/PolicyKit1/Authority"
	POLKIT_IFC  = "org.freedesktop.PolicyKit1.Authority"
)

func execCommand(cmdline string, args []string) {
	err := exec.Command(cmdline, args...).Run()
	if err != nil {
		logObject.Warningf("Exec '%v %v' failed:%v",
			cmdline, args, err)
		panic(err)
	}
}

func encodePasswd(words string) string {
	str := C.CString(words)
	defer C.free(unsafe.Pointer(str))

	ret := C.mkpasswd(str)
	return C.GoString(ret)
}

func getBaseName(path string) string {
	strs := strings.Split(path, "/")
	return strs[len(strs)-1]
}

func fileIsExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil || os.IsExist(err)
}

func isElementExist(element string, list []string) bool {
	for _, v := range list {
		if v == element {
			return true
		}
	}

	return false
}

/*
 * To determine whether the character is [A-Za-z0-9]
 */
func charIsAlNum(ch byte) bool {
	if (ch >= '0' && ch <= '9') ||
		(ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') {
		return true
	}

	return false
}

func deleteElementFromList(ele string, list []string) []string {
	tmp := []string{}
	for _, l := range list {
		if ele == l {
			continue
		}
		tmp = append(tmp, l)
	}

	return tmp
}

func polkitAuthWithPid(actionId string, pid uint32) bool {
	objPolkit, err := libpolkit1.NewAuthority("org.freedesktop.PolicyKit1",
		"/org/freedesktop/PolicyKit1/Authority")
	if err != nil {
		logObject.Warning("New PolicyKit Object Failed: ", err)
		return false
	}

	subject := polkitSubject{}
	subject.SubjectKind = "unix-process"
	subject.SubjectDetails = make(map[string]dbus.Variant)
	subject.SubjectDetails["pid"] = dbus.MakeVariant(uint32(pid))
	subject.SubjectDetails["start-time"] = dbus.MakeVariant(uint64(0))
	details := make(map[string]string)
	details[""] = ""
	flags := uint32(1)
	cancelId := ""

	ret, err1 := objPolkit.CheckAuthorization(subject, actionId, details,
		flags, cancelId)
	if err1 != nil {
		logObject.Warning("CheckAuthorization Failed: ", err1)
		return false
	}

	// Is Authority
	if !ret[0].(bool) {
		return false
	}

	return true
}
