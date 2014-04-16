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
        "io/ioutil"
        "os"
        "strings"
)

const (
        POLKIT_CHANGED_OWN_DATA = "com.deepin.daemon.accounts.change-own-user-data"
        POLKIT_MANAGER_USER     = "com.deepin.daemon.accounts.user-administration"
        POLKIT_SET_LOGIN_OPTION = "com.deepin.daemon.accounts.set-login-option"

        ICON_SYSTEM_DIR = "/var/lib/AccountsService/icons"
        ICON_LOCAL_DIR  = "/var/lib/AccountsService/icons/local"
        SHADOW_FILE     = "/etc/shadow"
        SHADOW_BAK_FILE = "/etc/shadow.bak"
)

func (op *UserManager) SetUserName(dbusMsg dbus.DMessage, username string) bool {
        defer func() {
                if err := recover(); err != nil {
                        logObject.Warningf("Recover Error In SetUserName:%v",
                                err)
                }
        }()

        if ok := opUtils.PolkitAuthWithPid(POLKIT_MANAGER_USER,
                dbusMsg.GetSenderPID()); !ok {
                return false
        }

        if op.UserName != username {
                op.applyPropertiesChanged("UserName", username)
                op.setPropName("UserName")
        }

        return true
}

func (op *UserManager) SetHomeDir(dbusMsg dbus.DMessage, dir string) bool {
        defer func() {
                if err := recover(); err != nil {
                        logObject.Warningf("Recover Error In SetHomeDir:%v",
                                err)
                }
        }()

        if ok := opUtils.PolkitAuthWithPid(POLKIT_MANAGER_USER,
                dbusMsg.GetSenderPID()); !ok {
                return false
        }

        if op.HomeDir != dir {
                op.applyPropertiesChanged("HomeDir", dir)
                op.setPropName("HomeDir")
        }

        return true
}

func (op *UserManager) SetShell(dbusMsg dbus.DMessage, shell string) bool {
        defer func() {
                if err := recover(); err != nil {
                        logObject.Warningf("Recover Error In SetShell:%v",
                                err)
                }
        }()

        if ok := opUtils.PolkitAuthWithPid(POLKIT_MANAGER_USER,
                dbusMsg.GetSenderPID()); !ok {
                return false
        }

        if op.Shell != shell {
                op.applyPropertiesChanged("Shell", shell)
                op.setPropName("Shell")
        }

        return true
}

func (op *UserManager) SetPassword(dbusMsg dbus.DMessage, words string) bool {
        defer func() {
                if err := recover(); err != nil {
                        logObject.Warningf("Recover Error In SetPassword:%v",
                                err)
                }
        }()

        if ok := opUtils.PolkitAuthWithPid(POLKIT_MANAGER_USER,
                dbusMsg.GetSenderPID()); !ok {
                return false
        }

        passwd := encodePasswd(words)
        changePasswd(op.UserName, passwd)

        //args := []string{}
        //args = append(args, "-p")
        //args = append(args, passwd)
        //args = append(args, op.UserName)
        //execCommand(CMD_USERMOD, args)

        op.applyPropertiesChanged("Locked", false)
        op.setPropName("Locked")

        return true
}

func (op *UserManager) SetAutomaticLogin(dbusMsg dbus.DMessage, auto bool) bool {
        defer func() {
                if err := recover(); err != nil {
                        logObject.Warningf("Recover Error In SetAutomaticLogin:%v",
                                err)
                }
        }()

        if ok := opUtils.PolkitAuthWithPid(POLKIT_MANAGER_USER,
                dbusMsg.GetSenderPID()); !ok {
                return false
        }

        if op.AutomaticLogin != auto {
                op.applyPropertiesChanged("AutomaticLogin", auto)
                op.setPropName("AutomaticLogin")
        }

        return true
}

func (op *UserManager) SetAccountType(dbusMsg dbus.DMessage, t int32) bool {
        defer func() {
                if err := recover(); err != nil {
                        logObject.Warningf("Recover Error In SetAccountType:%v",
                                err)
                }
        }()

        if ok := opUtils.PolkitAuthWithPid(POLKIT_MANAGER_USER,
                dbusMsg.GetSenderPID()); !ok {
                return false
        }

        if op.AccountType != t {
                op.applyPropertiesChanged("AccountType", t)
                op.setPropName("AccountType")
        }

        return true
}

func (op *UserManager) SetLocked(dbusMsg dbus.DMessage, locked bool) bool {
        defer func() {
                if err := recover(); err != nil {
                        logObject.Warningf("Recover Error In SetLocked:%v",
                                err)
                }
        }()

        if ok := opUtils.PolkitAuthWithPid(POLKIT_MANAGER_USER,
                dbusMsg.GetSenderPID()); !ok {
                return false
        }

        if op.Locked != locked {
                op.applyPropertiesChanged("Locked", locked)
                op.setPropName("Locked")
        }

        return true
}

func (op *UserManager) SetIconFile(dbusMsg dbus.DMessage, icon string) bool {
        defer func() {
                if err := recover(); err != nil {
                        logObject.Warningf("Recover Error In SetIconFile:%v",
                                err)
                }
        }()

        if ok := opUtils.PolkitAuthWithPid(POLKIT_MANAGER_USER,
                dbusMsg.GetSenderPID()); !ok {
                return false
        }

        if ok := opUtils.IsFileExist(icon); !ok || op.IconFile == icon {
                return false
        }

        if !opUtils.IsElementExist(icon, op.IconList) {
                if ok := opUtils.IsFileExist(ICON_LOCAL_DIR); !ok {
                        if err := os.MkdirAll(ICON_LOCAL_DIR, 0755); err != nil {
                                return false
                        }
                }
                name, _ := opUtils.GetBaseName(icon)
                dest := ICON_LOCAL_DIR + "/" + op.UserName + "-" + name
                if ok := opUtils.CopyFile(icon, dest); !ok {
                        return false
                }
                icon = dest
        }
        op.applyPropertiesChanged("IconFile", icon)
        op.setPropName("IconFile")

        return true
}

func (op *UserManager) SetBackgroundFile(dbusMsg dbus.DMessage, bg string) bool {
        defer func() {
                if err := recover(); err != nil {
                        logObject.Warningf("Recover Error In SetBackgroundFile:%v",
                                err)
                }
        }()

        if ok := opUtils.PolkitAuthWithPid(POLKIT_MANAGER_USER,
                dbusMsg.GetSenderPID()); !ok {
                return false
        }

        if op.BackgroundFile != bg {
                op.applyPropertiesChanged("BackgroundFile", bg)
                op.setPropName("BackgroundFile")
        }

        return true
}

func (op *UserManager) DeleteHistoryIcon(dbusMsg dbus.DMessage, icon string) bool {
        defer func() {
                if err := recover(); err != nil {
                        logObject.Warningf("Recover Error In DeleteHistoryIcon:%v",
                                err)
                }
        }()

        if ok := opUtils.PolkitAuthWithPid(POLKIT_MANAGER_USER,
                dbusMsg.GetSenderPID()); !ok {
                return false
        }

        file := USER_CONFIG_FILE + op.UserName
        deleteHistoryIcon(file, icon)
        op.setPropName("HistoryIcons")

        return true
}

func changePasswd(username, password string) {
        mutex.Lock()
        defer mutex.Unlock()

        data, err := ioutil.ReadFile(SHADOW_FILE)
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

                f, err := os.Create(SHADOW_BAK_FILE)
                if err != nil {
                        logObject.Warningf("Create '%s' failed: %v\n",
                                SHADOW_BAK_FILE, err)
                        panic(err)
                }
                defer f.Close()

                _, err = f.WriteString(contents)
                if err != nil {
                        logObject.Warningf("WriteString '%s' failed: %v\n",
                                SHADOW_BAK_FILE, err)
                        panic(err)
                }
                f.Sync()
                os.Rename(SHADOW_BAK_FILE, SHADOW_FILE)
        }

}

func newUserManager(uid string) *UserManager {
        defer func() {
                if err := recover(); err != nil {
                        logObject.Infof("Recover Error: %v", err)
                }
        }()

        m := &UserManager{}

        m.Uid = uid
        m.updateUserInfo()
        m.setPropName("IconList")
        m.listenUserInfoChanged(ETC_GROUP)
        m.listenUserInfoChanged(ETC_SHADOW)
        m.listenIconListChanged(ICON_SYSTEM_DIR)
        m.listenIconListChanged(ICON_LOCAL_DIR)

        return m
}
