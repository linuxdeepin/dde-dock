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
        //freedbus "dbus/org/freedesktop/dbus"
        polkit "dbus/org/freedesktop/policykit1"
        "dlib/dbus"
        "dlib/glib-2.0"
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

func readKeyFileValue(filename, group, key string, t int32) (interface{}, bool) {
        if !fileIsExist(filename) {
                return nil, false
        }

        keyFile := glib.NewKeyFile()
        defer keyFile.Free()
        ok, _ := keyFile.LoadFromFile(filename, glib.KeyFileFlagsKeepComments)
        if !ok {
                logObject.Warningf("LoadKeyFile '%s' failed", filename)
                return nil, false
        }

        switch t {
        case KEY_TYPE_BOOL:
                v, err := keyFile.GetBoolean(group, key)
                if err != nil {
                        //logObject.Warningf("Get '%s' from '%s' failed: %s",
                        //key, filename, err)
                        break
                }
                return v, true
        case KEY_TYPE_INT:
                v, err := keyFile.GetInteger(group, key)
                if err != nil {
                        //logObject.Warningf("Get '%s' from '%s' failed: %s",
                        //key, filename, err)
                        break
                }
                return v, true
        case KEY_TYPE_STRING:
                v, err := keyFile.GetString(group, key)
                if err != nil {
                        //logObject.Warningf("Get '%s' from '%s' failed: %s",
                        //key, filename, err)
                        break
                }
                return v, true
        case KEY_TYPE_STRING_LIST:
                _, v, err := keyFile.GetStringList(group, key)
                if err != nil {
                        break
                }
                return v, true
        }

        return nil, false
}

func writeKeyFileValue(filename, group, key string, t int32, value interface{}) {
        if !fileIsExist(filename) {
                f, err := os.Create(filename)
                if err != nil {
                        logObject.Infof("Create '%s' failed: %v",
                                filename, err)
                        return
                }
                f.Close()
                writeKeyFileValue(filename, "User", "Icon",
                        KEY_TYPE_STRING, USER_DEFAULT_ICON)
                writeKeyFileValue(filename, "User", "Background",
                        KEY_TYPE_STRING, USER_DEFAULT_BG)
        }

        mutex.Lock()
        defer mutex.Unlock()
        keyFile := glib.NewKeyFile()
        defer keyFile.Free()
        ok, _ := keyFile.LoadFromFile(filename, glib.KeyFileFlagsKeepComments)
        if !ok {
                logObject.Warningf("LoadKeyFile '%s' failed", filename)
                return
        }

        switch t {
        case KEY_TYPE_BOOL:
                keyFile.SetBoolean(group, key, value.(bool))
        case KEY_TYPE_INT:
                keyFile.SetInteger(group, key, value.(int))
        case KEY_TYPE_STRING:
                keyFile.SetString(group, key, value.(string))
        case KEY_TYPE_STRING_LIST:
                keyFile.SetStringList(group, key, value.([]string))
        }

        _, contents, err := keyFile.ToData()
        if err != nil {
                logObject.Warningf("KeyFile '%s' ToData failed: %s", filename, err)
                panic(err)
        }

        writeKeyFile(contents, filename)
}

func writeKeyFile(contents, file string) {
        if len(file) <= 0 {
                return
        }

        //logObject.Warning(contents)
        //return
        f, err := os.Create(file + "~")
        if err != nil {
                logObject.Warningf("OpenFile '%s' failed: %s",
                        file+"~", err)
                panic(err)
        }
        defer f.Close()

        if _, err = f.WriteString(contents); err != nil {
                logObject.Warningf("Write in '%s' failed: %s", file, err)
                panic(err)
        }
        f.Sync()
        os.Rename(file+"~", file)
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

func authWithPolkit(actionId string, pid uint32) bool {
        var (
                objPolkit *polkit.Authority
                err       error
        )

        objPolkit, err = polkit.NewAuthority(POLKIT_DEST, POLKIT_PATH)
        if err != nil {
                logObject.Warningf("New Authority Failed:%v", err)
                //panic(err)
                return false
        }

        subject := polkitSubject{}
        subject.SubjectKind = "unix-process"
        subject.SubjectDetails = make(map[string]dbus.Variant)
        subject.SubjectDetails["pid"] = dbus.MakeVariant(uint32(pid))
        subject.SubjectDetails["start-time"] = dbus.MakeVariant(uint64(0))
        details := make(map[string]string)
        //details[""] = ""
        flags := uint32(1)
        cancelId := ""

        rets, _err := objPolkit.CheckAuthorization(subject, actionId, details, flags, cancelId)
        //println("Is Authority: ", rets[0].(bool))
        //println("Is Challonge: ", rets[1].(bool))
        if _err != nil {
                logObject.Warningf("CheckAuthorization Failed:%v", _err)
                return false
        }

        // Is Authority
        if !rets[0].(bool) {
                return false
        }

        return true
}
