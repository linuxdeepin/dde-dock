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
)

const (
        CMD_USERMOD = "/usr/sbin/usermod"
        CMD_GPASSWD = "/usr/bin/gpasswd"

        USER_ICON_DIR     = "/var/lib/AccountsService/icons/"
        USER_DEFAULT_ICON = USER_ICON_DIR + "001.jpg"
        USER_DEFAULT_BG   = "/usr/share/backgrounds/default_background.jpg"
        USER_CONFIG_FILE  = "/var/lib/AccountsService/users/"
)

type UserManager struct {
        Uid            string
        Gid            string
        UserName       string
        HomeDir        string
        Shell          string
        IconFile       string
        BackgroundFile string
        AutomaticLogin bool
        AccountType    int32
        Locked         bool
        LoginTime      uint64
        objectPath     string
}

func (op *UserManager) GetDBusInfo() dbus.DBusInfo {
        return dbus.DBusInfo{
                ACCOUNT_DEST,
                USER_MANAGER_PATH + op.Uid,
                USER_MANAGER_IFC,
        }
}

func (op *UserManager) setPropName(propName string, propValue interface{}) {
        switch propName {
        case "UserName":
                args := []string{}
                args = append(args, "-l")
                args = append(args, propValue.(string))
                args = append(args, op.UserName)
                execCommand(CMD_USERMOD, args)
        case "HomeDir":
                args := []string{}
                args = append(args, "-m")
                args = append(args, "-d")
                args = append(args, propValue.(string))
                args = append(args, op.UserName)
                execCommand(CMD_USERMOD, args)
        case "Shell":
                args := []string{}
                args = append(args, "-s")
                args = append(args, propValue.(string))
                args = append(args, op.UserName)
                execCommand(CMD_USERMOD, args)
        case "IconFile":
                file := USER_CONFIG_FILE + op.UserName
                writeKeyFileValue(file, "User", "Background",
                        KEY_TYPE_STRING, op.BackgroundFile)
        case "BackgroundFile":
                file := USER_CONFIG_FILE + op.UserName
                writeKeyFileValue(file, "User", "Background",
                        KEY_TYPE_STRING, op.BackgroundFile)
        case "AutomaticLogin":
                if op.AutomaticLogin {
                        setAutomaticLogin(op.UserName)
                } else {
                        setAutomaticLogin("")
                }
        case "AccountType":
                accountTyte := propValue.(int32)
                switch accountTyte {
                case ACCOUNT_TYPE_STANDARD:
                case ACCOUNT_TYPE_ADMINISTACTOR:
                        addUserToAdmList(op.UserName)
                }
        case "Locked":
                args := []string{}
                if propValue.(bool) {
                        args = append(args, "-L")
                        args = append(args, op.UserName)
                } else {
                        args = append(args, "-U")
                        args = append(args, op.UserName)
                }
                execCommand(CMD_USERMOD, args)
        }
}

func (op *UserManager) getPropName(propName string) {
        switch propName {
        case "UserName":
                info, ok := getInfoViaUid(op.Uid)
                if ok {
                        op.UserName = info.Name
                }
        case "HomeDir":
                info, ok := getInfoViaUid(op.Uid)
                if ok {
                        op.UserName = info.Home
                }
        case "Shell":
                info, ok := getInfoViaUid(op.Uid)
                if ok {
                        op.Shell = info.Shell
                }
        case "IconFile":
                file := USER_CONFIG_FILE + op.UserName
                if !fileIsExist(file) {
                        op.IconFile = USER_DEFAULT_ICON
                } else {
                        v, ok := readKeyFileValue(file, "User", "Background", KEY_TYPE_STRING)
                        if !ok {
                                op.IconFile = USER_DEFAULT_ICON
                        } else {
                                op.IconFile = v.(string)
                        }
                }
        case "BackgroundFile":
                file := USER_CONFIG_FILE + op.UserName
                if !fileIsExist(file) {
                        op.BackgroundFile = USER_DEFAULT_BG
                } else {
                        v, ok := readKeyFileValue(file, "User", "Background", KEY_TYPE_STRING)
                        if !ok {
                                op.BackgroundFile = USER_DEFAULT_BG
                        } else {
                                op.BackgroundFile = v.(string)
                        }
                }
        case "AutomaticLogin":
                ok := isAutoLogin(op.UserName)
                if ok {
                        op.AutomaticLogin = true
                } else {
                        op.AutomaticLogin = false
                }
        case "AccountType":
                admList := getAdministratorList()
                if isElementExist(op.UserName, admList) {
                        op.AccountType = ACCOUNT_TYPE_ADMINISTACTOR
                } else {
                        op.AccountType = ACCOUNT_TYPE_STANDARD
                }
        case "Locked":
                info, ok := getInfoViaUid(op.Uid)
                if ok {
                        op.Locked = info.Locked
                }
        }
        dbus.NotifyChange(op, propName)
}

func (op *UserManager) updateUserInfo() {
        info, ok := getInfoViaUid(op.Uid)
        if !ok {
                return
        }

        op.Gid = info.Gid
        op.UserName = info.Name
        op.HomeDir = info.Home
        op.Locked = info.Locked
        op.Shell = info.Shell
        op.getPropName("IconFile")
        op.getPropName("BackgroundFile")
        op.getPropName("AutomaticLogin")
        op.getPropName("AccountType")
        op.getPropName("LoginTime")
}

func addUserToAdmList(name string) {
        tmps := []string{}
        tmps = append(tmps, "-a")
        tmps = append(tmps, name)
        tmps = append(tmps, "sudo")
        go execCommand(CMD_GPASSWD, tmps)
}
