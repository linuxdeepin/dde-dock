/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
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

package accounts

// #cgo CFLAGS: -Wall -g
// #cgo LDFLAGS: -lcrypt
// #include <stdlib.h>
// #include "mkpasswd.h"
import "C"

import (
	libpolkit1 "dbus/org/freedesktop/policykit1"
	"dlib/dbus"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
	"unsafe"
)

const (
	POLKIT_DEST = "org.freedesktop.PolicyKit1"
	POLKIT_PATH = "/org/freedesktop/PolicyKit1/Authority"
	POLKIT_IFC  = "org.freedesktop.PolicyKit1.Authority"
)

func execCommand(cmd string, args []string) {
	if err := exec.Command(cmd, args...).Run(); err != nil {
		logger.Errorf("Exec '%s %v' failed: %v", cmd, args, err)
	}
}

func strIsInList(str string, list []string) bool {
	for _, v := range list {
		if str == v {
			return true
		}
	}

	return false
}

func isStrListEqual(list1, list2 []string) bool {
	l1 := len(list1)
	l2 := len(list2)

	if l1 != l2 {
		return false
	}

	for i := 0; i < l1; i++ {
		if list1[i] != list2[i] {
			return false
		}
	}

	return true
}

func deleteStrFromList(ele string, list []string) []string {
	tmp := []string{}
	for _, l := range list {
		if ele == l {
			continue
		}
		tmp = append(tmp, l)
	}

	return tmp
}

func compareStrList(list1, list2 []string) ([]string, int) {
	if isStrListEqual(list1, list2) {
		return []string{}, 0
	}

	l1 := len(list1)
	l2 := len(list2)

	tmp := []string{}
	if l1 < l2 {
		for i := 0; i < l2; i++ {
			j := 0
			for ; j < l1; j++ {
				if list1[j] == list2[i] {
					break
				}
			}
			if j == l1 {
				tmp = append(tmp, list2[i])
			}
		}

		return tmp, 1
	}

	if l1 > l2 {
		for i := 0; i < l1; i++ {
			j := 0
			for ; j < l2; j++ {
				if list1[i] == list2[j] {
					break
				}
			}
			if j == l2 {
				tmp = append(tmp, list1[i])
			}
		}

		return tmp, -1
	}

	return tmp, 0
}

func changeFileOwner(user, group, dir string) {
	args := []string{}
	args = append(args, "-R")
	args = append(args, user+":"+group)
	args = append(args, dir)

	go execCommand(CMD_CHOWN, args)
}

func encodePasswd(words string) string {
	str := C.CString(words)
	defer C.free(unsafe.Pointer(str))

	ret := C.mkpasswd(str)
	return C.GoString(ret)
}

func changePasswd(username, password string) {
	data, err := ioutil.ReadFile(ETC_SHADOW)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(data), "\n")
	index := 0
	line := ""
	okFlag := false
	for index, line = range lines {
		strs := strings.Split(line, ":")
		if strs[0] == username {
			if strs[1] == password {
				break
			}
			strs[1] = password
			l := len(strs)
			line = ""
			for i, s := range strs {
				if i == l-1 {
					line += s
					continue
				}
				line += s + ":"
			}
			okFlag = true
			break
		}
	}

	if okFlag {
		okFlag = false
		contents := ""
		l := len(lines)
		for i, tmp := range lines {
			if i == index {
				contents += line
			} else {
				contents += tmp
			}
			if i < l-1 {
				contents += "\n"
			}
		}

		f, err := os.Create(ETC_SHADOW_BAK)
		if err != nil {
			logger.Errorf("Create '%s' failed: %v\n",
				ETC_SHADOW_BAK, err)
			panic(err)
		}
		defer f.Close()

		var mutex sync.Mutex
		mutex.Lock()
		_, err = f.WriteString(contents)
		if err != nil {
			logger.Errorf("WriteString '%s' failed: %v\n",
				ETC_SHADOW_BAK, err)
			panic(err)
		}
		f.Sync()
		os.Rename(ETC_SHADOW_BAK, ETC_SHADOW)
		mutex.Unlock()
	}
}

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

func polkitAuthWithPid(actionId string, pid uint32) bool {
	objPolkit, err := libpolkit1.NewAuthority("org.freedesktop.PolicyKit1",
		"/org/freedesktop/PolicyKit1/Authority")
	if err != nil {
		logger.Error("New PolicyKit Object Failed: ", err)
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
		logger.Error("CheckAuthorization Failed: ", err1)
		return false
	}

	// Is Authority
	if !ret[0].(bool) {
		return false
	}

	return true
}
