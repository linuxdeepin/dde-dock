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
        "io/ioutil"
        "os"
        "strings"
)

const (
        POLKIT_CHANGED_OWN_DATA = "com.deepin.daemon.accounts.change-own-user-data"
        POLKIT_MANAGER_USER     = "com.deepin.daemon.accounts.user-administration"
        POLKIT_SET_LOGIN_OPTION = "com.deepin.daemon.accounts.set-login-option"
        ICON__SYSTEM_DIR        = "/var/lib/AccountsService/icons/"
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

func (op *UserManager) SetPassword(passwd string) {
        defer func() {
                if err := recover(); err != nil {
                        logObject.Warning("Recover Error In SetPassword:%v",
                                err)
                }
        }()
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
        if op.IconFile != icon {
                op.applyPropertiesChanged("IconFile", icon)
                op.setPropName("IconFile")
        }
}

func (op *UserManager) SetBackgroundFile(bg string) {
        //authWithPolkit(POLKIT_CHANGED_OWN_DATA)
        if op.BackgroundFile != bg {
                op.applyPropertiesChanged("BackgroundFile", bg)
                op.setPropName("BackgroundFile")
        }
}

func (op *UserManager) GetIconList() []string {
        list := []string{}

        sysList := getSystemIconList()
        list = append(list, sysList...)

        return list
}

func (op *UserManager) DeleteHistoryIcon(icon string) {
        file := USER_CONFIG_FILE + op.UserName
        op.HistoryIcons = deleteHistoryIcon(file, icon)
}

func newUserManager(uid string) *UserManager {
        m := &UserManager{}

        m.Uid = uid
        m.updateUserInfo()
        m.listenUserInfoChanged(ETC_GROUP)
        m.listenUserInfoChanged(ETC_SHADOW)

        return m
}

func getSystemIconList() []string {
        iconfd, err := os.Open(ICON__SYSTEM_DIR)
        if err != nil {
                logObject.Warning("Open '%s' failed: %v\n",
                        ICON__SYSTEM_DIR, err)
                return []string{}
        }

        names, _ := iconfd.Readdirnames(0)
        list := []string{}
        for _, v := range names {
                if strings.Contains(v, "guest") {
                        continue
                } else if strings.Contains(v, "jpg") ||
                        strings.Contains(v, "JPG") ||
                        strings.Contains(v, "png") ||
                        strings.Contains(v, "PNG") {
                        list = append(list, ICON__SYSTEM_DIR+v)
                }
        }

        return list
}

func getAdministratorList() []string {
        contents, err := ioutil.ReadFile(ETC_GROUP)
        if err != nil {
                logObject.Warning("ReadFile '%s' failed: %s\n", ETC_PASSWD, err)
                panic(err)
        }

        list := ""
        lines := strings.Split(string(contents), "\n")
        for _, line := range lines {
                strs := strings.Split(line, ":")
                if len(strs) != GROUP_SPLIT_LEN {
                        continue
                }

                if strs[0] == "sudo" {
                        list = strs[3]
                        break
                }
        }

        return strings.Split(list, ",")
}

func setAutomaticLogin(name string) {
        dsp := getDefaultDisplayManager()
        switch dsp {
        case "lightdm":
                if fileIsExist(ETC_LIGHTDM_CONFIG) {
                        writeKeyFileValue(ETC_LIGHTDM_CONFIG,
                                LIGHTDM_AUTOLOGIN_GROUP,
                                LIGHTDM_AUTOLOGIN_USER,
                                KEY_TYPE_STRING, name)
                }
        case "gdm":
                if fileIsExist(ETC_GDM_CONFIG) {
                        writeKeyFileValue(ETC_GDM_CONFIG,
                                GDM_AUTOLOGIN_GROUP,
                                GDM_AUTOLOGIN_USER,
                                KEY_TYPE_STRING, name)
                }
        case "kdm":
                if fileIsExist(ETC_KDM_CONFIG) {
                        writeKeyFileValue(ETC_KDM_CONFIG,
                                KDM_AUTOLOGIN_GROUP,
                                KDM_AUTOLOGIN_ENABLE,
                                KEY_TYPE_BOOL, true)
                        writeKeyFileValue(ETC_KDM_CONFIG,
                                KDM_AUTOLOGIN_GROUP,
                                KDM_AUTOLOGIN_USER,
                                KEY_TYPE_STRING, name)
                } else if fileIsExist(USER_KDM_CONFIG) {
                        writeKeyFileValue(ETC_KDM_CONFIG,
                                KDM_AUTOLOGIN_GROUP,
                                KDM_AUTOLOGIN_ENABLE,
                                KEY_TYPE_BOOL, true)
                        writeKeyFileValue(USER_KDM_CONFIG,
                                KDM_AUTOLOGIN_GROUP,
                                KDM_AUTOLOGIN_USER,
                                KEY_TYPE_STRING, name)
                }
        default:
                logObject.Warning("No support display manager")
        }
}

func isAutoLogin(username string) bool {
        dsp := getDefaultDisplayManager()

        switch dsp {
        case "lightdm":
                if fileIsExist(ETC_LIGHTDM_CONFIG) {
                        v, ok := readKeyFileValue(ETC_LIGHTDM_CONFIG,
                                LIGHTDM_AUTOLOGIN_GROUP,
                                LIGHTDM_AUTOLOGIN_USER,
                                KEY_TYPE_STRING)
                        if ok && v.(string) == username {
                                return true
                        }
                }
        case "gdm":
                if fileIsExist(ETC_GDM_CONFIG) {
                        v, ok := readKeyFileValue(ETC_GDM_CONFIG,
                                GDM_AUTOLOGIN_GROUP,
                                GDM_AUTOLOGIN_USER,
                                KEY_TYPE_STRING)
                        if ok && v.(string) == username {
                                return true
                        }
                }
        case "kdm":
                if fileIsExist(ETC_KDM_CONFIG) {
                        v, ok := readKeyFileValue(ETC_KDM_CONFIG,
                                KDM_AUTOLOGIN_GROUP,
                                KDM_AUTOLOGIN_USER,
                                KEY_TYPE_STRING)
                        if ok && v.(string) == username {
                                return true
                        }
                } else if fileIsExist(USER_KDM_CONFIG) {
                        v, ok := readKeyFileValue(USER_KDM_CONFIG,
                                KDM_AUTOLOGIN_GROUP,
                                KDM_AUTOLOGIN_USER,
                                KEY_TYPE_STRING)
                        if ok && v.(string) == username {
                                return true
                        }
                }
        }

        return false
}

func getDefaultDisplayManager() string {
        contents, err := ioutil.ReadFile(ETC_DISPLAY_MANAGER)
        if err != nil {
                logObject.Warning("ReadFile '%s' failed: %s\n",
                        ETC_DISPLAY_MANAGER, err)
                panic(err)
        }

        tmp := ""
        for _, b := range contents {
                if b == '\n' {
                        tmp += ""
                        continue
                }
                tmp += string(b)
        }

        return getBaseName(tmp)
}
