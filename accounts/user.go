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
        "strings"
)

const (
        POLKIT_CHANGED_OWN_DATA = "com.deepin.daemon.accounts.change-own-user-data"
        POLKIT_MANAGER_USER     = "com.deepin.daemon.accounts.user-administration"
        POLKIT_SET_LOGIN_OPTION = "com.deepin.daemon.accounts.set-login-option"
)

func (op *UserManager) SetUserName(username string) {
        if op.UserName != username {
                //authWithPolkit(POLKIT_CHANGED_OWN_DATA)
                op.setPropName("UserName", username)
                op.getPropName("UserName")
        }
}

func (op *UserManager) SetHomeDir(dir string) {
        if op.HomeDir != dir {
                //authWithPolkit(POLKIT_CHANGED_OWN_DATA)
                op.setPropName("HomeDir", dir)
                op.getPropName("HomeDir")
        }
}

func (op *UserManager) SetShell(shell string) {
        if op.Shell != shell {
                //authWithPolkit(POLKIT_CHANGED_OWN_DATA)
                op.setPropName("Shell", shell)
                op.getPropName("Shell")
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
        if op.AutomaticLogin != auto {
                //authWithPolkit(POLKIT_SET_LOGIN_OPTION)
                op.setPropName("AutomaticLogin", auto)
                op.getPropName("AutomaticLogin")
        }
}

func (op *UserManager) SetAccountType(t int32) {
        logObject.Warning("src type:%v", op.AccountType)
        logObject.Warning("dest type:%v", t)
        if op.AccountType != t {
                //authWithPolkit(POLKIT_MANAGER_USER)
                op.setPropName("AccountType", t)
                op.getPropName("AccountType")
        }
}

func (op *UserManager) SetLocked(locked bool) {
        if op.Locked != locked {
                //authWithPolkit(POLKIT_MANAGER_USER)
                op.setPropName("Locked", locked)
                op.getPropName("Locked")
        }
}

func (op *UserManager) SetIconFile(icon string) {
        if op.IconFile != icon {
                //authWithPolkit(POLKIT_CHANGED_OWN_DATA)
                op.setPropName("IconFile", icon)
                op.getPropName("IconFile")
        }
}

func (op *UserManager) SetBackgroundFile(bg string) {
        if op.BackgroundFile != bg {
                //authWithPolkit(POLKIT_CHANGED_OWN_DATA)
                op.setPropName("BackgroundFile", bg)
                op.getPropName("BackgroundFile")
        }
}

func newUserManager(uid string) *UserManager {
        m := &UserManager{}

        m.Uid = uid
        m.updateUserInfo()

        return m
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
