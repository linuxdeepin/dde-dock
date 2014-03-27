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
        "os"
)

const (
        POLKIT_CHANGED_OWN_DATA = "com.deepin.daemon.accounts.change-own-user-data"
        POLKIT_MANAGER_USER     = "com.deepin.daemon.accounts.user-administration"
        POLKIT_SET_LOGIN_OPTION = "com.deepin.daemon.accounts.set-login-option"

        ICON_SYSTEM_DIR = "/var/lib/AccountsService/icons"
        ICON_LOCAL_DIR  = "/var/lib/AccountsService/icons/local"
)

func (op *UserManager) SetUserName(username string) {
        //authWithPolkit(POLKIT_CHANGED_OWN_DATA)
        if op.UserName != username {
                op.applyPropertiesChanged("UserName", username)
                op.setPropName("UserName")
        }
}

func (op *UserManager) SetHomeDir(dir string) {
        //authWithPolkit(POLKIT_CHANGED_OWN_DATA)
        if op.HomeDir != dir {
                op.applyPropertiesChanged("HomeDir", dir)
                op.setPropName("HomeDir")
        }
}

func (op *UserManager) SetShell(shell string) {
        //authWithPolkit(POLKIT_CHANGED_OWN_DATA)
        if op.Shell != shell {
                op.applyPropertiesChanged("Shell", shell)
                op.setPropName("Shell")
        }
}

func (op *UserManager) SetPassword(words string) {
        defer func() {
                if err := recover(); err != nil {
                        logObject.Warningf("Recover Error In SetPassword:%v",
                                err)
                }
        }()
        passwd := encodePasswd(words)
        args := []string{}

        args = append(args, "-p")
        args = append(args, passwd)
        args = append(args, op.UserName)
        execCommand(CMD_USERMOD, args)

        op.SetLocked(false)
}

func (op *UserManager) SetAutomaticLogin(auto bool) {
        //authWithPolkit(POLKIT_SET_LOGIN_OPTION)
        if op.AutomaticLogin != auto {
                op.applyPropertiesChanged("AutomaticLogin", auto)
                op.setPropName("AutomaticLogin")
        }
}

func (op *UserManager) SetAccountType(t int32) {
        //authWithPolkit(POLKIT_MANAGER_USER)
        if op.AccountType != t {
                op.applyPropertiesChanged("AccountType", t)
                op.setPropName("AccountType")
        }
}

func (op *UserManager) SetLocked(locked bool) {
        //authWithPolkit(POLKIT_MANAGER_USER)
        if op.Locked != locked {
                op.applyPropertiesChanged("Locked", locked)
                op.setPropName("Locked")
        }
}

func (op *UserManager) SetIconFile(icon string) {
        //authWithPolkit(POLKIT_CHANGED_OWN_DATA)
        if ok := opUtils.IsFileExist(icon); !ok || op.IconFile == icon {
                return
        }

        if !opUtils.IsElementExist(icon, op.IconList) {
                if ok := opUtils.IsFileExist(ICON_LOCAL_DIR); !ok {
                        if err := os.MkdirAll(ICON_LOCAL_DIR, 0755); err != nil {
                                return
                        }
                }
                name, _ := opUtils.GetBaseName(icon)
                dest := ICON_LOCAL_DIR + "/" + op.UserName + "-" + name
                if ok := opUtils.CopyFile(icon, dest); !ok {
                        return
                }
                icon = dest
        }
        op.applyPropertiesChanged("IconFile", icon)
        op.setPropName("IconFile")
}

func (op *UserManager) SetBackgroundFile(bg string) {
        //authWithPolkit(POLKIT_CHANGED_OWN_DATA)
        if op.BackgroundFile != bg {
                op.applyPropertiesChanged("BackgroundFile", bg)
                op.setPropName("BackgroundFile")
        }
}

func (op *UserManager) DeleteHistoryIcon(icon string) {
        file := USER_CONFIG_FILE + op.UserName
        op.HistoryIcons = deleteHistoryIcon(file, icon)
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
