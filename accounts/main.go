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
        "dlib"
        "dlib/dbus"
        "dlib/logger"
        dutils "dlib/utils"
        "os"
)

var (
        idUserManagerMap = make(map[string]*UserManager)
        logObject        = logger.NewLogger("daemon/Accounts")
        opUtils          *dutils.Manager
)

func main() {
        defer func() {
                if err := recover(); err != nil {
                        logObject.Fatal("Recover Error:", err)
                }
        }()

        if !dlib.UniqueOnSystem(ACCOUNT_DEST) {
                logObject.Warning("There already has an Accounts daemon running.")
                return
        }

        // Configure Logger
        logObject.SetRestartCommand("/usr/lib/deepin-daemon/Accounts")
        var err error
        opUtils = dutils.NewUtils()

        opAccount := newAccountManager()
        err = dbus.InstallOnSystem(opAccount)
        if err != nil {
                logObject.Warningf("Install Account Object On System Failed:%v", err)
                logObject.Fatalf("%v", err)
        }

        updateUserList()

        dbus.DealWithUnhandledMessage()

        //select {}
        if err = dbus.Wait(); err != nil {
                logObject.Warningf("lost dbus session:%v", err)
                os.Exit(1)
        } else {
                os.Exit(0)
        }
}

func updateUserList() {
        destroyAllUserObject()

        infos := getUserInfoList()
        for _, info := range infos {
                opUser := newUserManager(info.Uid)
                err := dbus.InstallOnSystem(opUser)
                if err != nil {
                        logObject.Debugf("Install User:%s Object On System Failed:%s",
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
