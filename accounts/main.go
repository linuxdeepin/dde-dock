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

import (
        "dlib/dbus"
        "fmt"
)

var idUserManagerMap map[string]*UserManager = make(map[string]*UserManager)

func main() {
        defer func() {
                if err := recover(); err != nil {
                        fmt.Println("Recover Error:", err)
                }
        }()

        idUserManagerMap = make(map[string]*UserManager)

        opAccount := newAccountManager()
        err := dbus.InstallOnSystem(opAccount)
        if err != nil {
                fmt.Println("Install Account Object On System Failed:", err)
                panic(err)
        }

        updateUserList()

        listenFileChanged(ETC_PASSWD)
        listenFileChanged(ETC_GROUP)
        listenFileChanged(ETC_SHADOW)

        dbus.DealWithUnhandledMessage()

        select {}
}

func updateUserList() {
        destroyAllUserObject()

        infos := getUserInfoList()
        for _, info := range infos {
                opUser := newUserManager(info.Uid)
                err := dbus.InstallOnSystem(opUser)
                if err != nil {
                        fmt.Printf("Install User:%s Object On System Failed:%s\n",
                                info.Name, err)
                        panic(err)
                }

                idUserManagerMap[info.Uid] = opUser
        }
}

func destroyAllUserObject() {
        for k, v := range idUserManagerMap {
                dbus.UnInstallObject(v)
                delete(idUserManagerMap, k)
        }
}

func printUserInfo(info UserInfo) {
        fmt.Println("Name:", info.Name)
        fmt.Println("Uid:", info.Uid)
        fmt.Println("Gid:", info.Gid)
        fmt.Println("Home:", info.Home)
        fmt.Println("Shell:", info.Shell)
        fmt.Printf("\n")
}
