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
        "os"
        "os/exec"
        "strings"
        "sync"
        "unsafe"
)

var (
        mutex   = new(sync.Mutex)
        genId   = func() func() uint32 {
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
