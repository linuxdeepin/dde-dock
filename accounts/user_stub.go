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
        dutils "dbus/com/deepin/api/utils"
        "dlib/dbus"
)

const (
        CMD_USERMOD = "/usr/sbin/usermod"
        CMD_GPASSWD = "/usr/bin/gpasswd"

        USER_ICON_DIR     = "/var/lib/AccountsService/icons/"
        USER_DEFAULT_ICON = USER_ICON_DIR + "1.jpg"
        USER_DEFAULT_BG   = "file:///usr/share/backgrounds/default_background.jpg"
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
        HistoryIcons   []string
        objectPath     string
}

func (op *UserManager) GetDBusInfo() dbus.DBusInfo {
        return dbus.DBusInfo{
                ACCOUNT_DEST,
                USER_MANAGER_PATH + op.Uid,
                USER_MANAGER_IFC,
        }
}

func (op *UserManager) applyPropertiesChanged(propName string, value interface{}) {
        switch propName {
        case "UserName":
                if v, ok := value.(string); ok && v != op.UserName {
                        args := []string{}
                        args = append(args, "-l")
                        args = append(args, v)
                        args = append(args, op.UserName)
                        execCommand(CMD_USERMOD, args)
                }
        case "HomeDir":
                if v, ok := value.(string); ok && v != op.HomeDir {
                        args := []string{}
                        args = append(args, "-m")
                        args = append(args, "-d")
                        args = append(args, v)
                        args = append(args, op.UserName)
                        execCommand(CMD_USERMOD, args)
                }
        case "Shell":
                if v, ok := value.(string); ok && v != op.Shell {
                        args := []string{}
                        args = append(args, "-s")
                        args = append(args, v)
                        args = append(args, op.UserName)
                        execCommand(CMD_USERMOD, args)
                }
        case "IconFile":
                if v, ok := value.(string); ok && v != op.IconFile {
                        file := USER_CONFIG_FILE + op.UserName
                        writeKeyFileValue(file, "User", "Icon",
                                KEY_TYPE_STRING, v)
                        addHistoryIcon(file, v)
                        op.setPropName("HistoryIcons")
                }
        case "BackgroundFile":
                if v, ok := value.(string); ok && v != op.BackgroundFile {
                        file := USER_CONFIG_FILE + op.UserName
                        writeKeyFileValue(file, "User", "Background",
                                KEY_TYPE_STRING, v)
                }
        case "AutomaticLogin":
                if v, ok := value.(bool); ok && v != op.AutomaticLogin {
                        if v {
                                setAutomaticLogin(op.UserName)
                        } else {
                                setAutomaticLogin("")
                        }
                }
        case "AccountType":
                if v, ok := value.(int32); ok && v != op.AccountType {
                        switch v {
                        case ACCOUNT_TYPE_STANDARD:
                                admList := getAdministratorList()
                                if isElementExist(op.UserName, admList) {
                                        deleteUserFromAdmList(op.UserName)
                                }
                        case ACCOUNT_TYPE_ADMINISTACTOR:
                                admList := getAdministratorList()
                                if !isElementExist(op.UserName, admList) {
                                        addUserToAdmList(op.UserName)
                                }
                        }
                }
        case "Locked":
                if v, ok := value.(bool); ok && v != op.Locked {
                        args := []string{}
                        if v {
                                args = append(args, "-L")
                                args = append(args, op.UserName)
                        } else {
                                args = append(args, "-U")
                                args = append(args, op.UserName)
                        }
                        execCommand(CMD_USERMOD, args)
                }
        }
}

func (op *UserManager) setPropName(propName string) {
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
                        v, ok := readKeyFileValue(file, "User", "Icon", KEY_TYPE_STRING)
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
                                tmp := v.(string)
                                opUtils, _ := dutils.NewUtils("com.deepin.api.Utils", "/com/deepin/api/Utils")
                                uri, _, _ := opUtils.PathToFileURI(tmp)
                                op.BackgroundFile = uri
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
        case "HistoryIcons":
                file := USER_CONFIG_FILE + op.UserName
                if !fileIsExist(file) {
                        op.HistoryIcons = append(op.HistoryIcons,
                                USER_DEFAULT_ICON)
                } else {
                        v, ok := readKeyFileValue(file, "User",
                                "HistoryIcons", KEY_TYPE_STRING_LIST)
                        if !ok {
                                op.HistoryIcons = append(op.HistoryIcons,
                                        USER_DEFAULT_ICON)
                        } else {
                                op.HistoryIcons = v.([]string)
                        }
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
        op.setPropName("IconFile")
        op.setPropName("BackgroundFile")
        op.setPropName("AutomaticLogin")
        op.setPropName("AccountType")
        op.setPropName("LoginTime")
        op.setPropName("HistoryIcons")
}

func addHistoryIcon(filename, iconPath string) []string {
        list, _ := readKeyFileValue(filename, "User",
                "HistoryIcons", KEY_TYPE_STRING_LIST)

        ret := []string{}
        ret = append(ret, iconPath)
        if list != nil {
                tmp := deleteElementFromList(iconPath, list.([]string))
                ret = append(ret, tmp...)
        }
        writeKeyFileValue(filename, "User",
                "HistoryIcons", KEY_TYPE_STRING_LIST, ret)

        return ret
}

func deleteHistoryIcon(filename, iconPath string) []string {
        list, ok := readKeyFileValue(filename, "User",
                "HistoryIcons", KEY_TYPE_STRING_LIST)
        if !ok || len(list.([]string)) <= 0 {
                return []string{}
        }

        tmp := deleteElementFromList(iconPath, list.([]string))
        writeKeyFileValue(filename, "User",
                "HistoryIcons", KEY_TYPE_STRING_LIST, tmp)

        return tmp
}

func addUserToAdmList(name string) {
        tmps := []string{}
        tmps = append(tmps, "-a")
        tmps = append(tmps, name)
        tmps = append(tmps, "sudo")
        go execCommand(CMD_GPASSWD, tmps)
}

func deleteUserFromAdmList(name string) {
        tmps := []string{}
        tmps = append(tmps, "-d")
        tmps = append(tmps, name)
        tmps = append(tmps, "sudo")
        go execCommand(CMD_GPASSWD, tmps)
}
