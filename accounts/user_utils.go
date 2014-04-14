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
        "math/rand"
        "os"
        "strings"
)

func getRandUserIcon() string {
        list := getIconList(ICON_SYSTEM_DIR)
        l := len(list)
        if l <= 0 {
                return ""
        }

        index := rand.Int31n(int32(l))
        //logObject.Infof("Rand Icon Index: %d", index)
        //logObject.Infof("Rand Icon Path: %s", list[index])
        return list[index]
}

func getIconList(dir string) []string {
        iconfd, err := os.Open(dir)
        if err != nil {
                logObject.Warningf("Open '%s' failed: %v",
                        dir, err)
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
                        strings.Contains(v, "PNG") ||
                        strings.Contains(v, "jpeg") ||
                        strings.Contains(v, "JPEG") ||
                        strings.Contains(v, "gif") ||
                        strings.Contains(v, "GIF") {
                        list = append(list, dir+"/"+v)
                }
        }

        return list
}

func getAdministratorList() []string {
        contents, err := ioutil.ReadFile(ETC_GROUP)
        if err != nil {
                logObject.Warningf("ReadFile '%s' failed: %s", ETC_PASSWD, err)
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
                        opUtils.WriteKeyToKeyFile(ETC_LIGHTDM_CONFIG,
                                LIGHTDM_AUTOLOGIN_GROUP,
                                LIGHTDM_AUTOLOGIN_USER,
                                name)
                }
        case "gdm":
                if fileIsExist(ETC_GDM_CONFIG) {
                        opUtils.WriteKeyToKeyFile(ETC_GDM_CONFIG,
                                GDM_AUTOLOGIN_GROUP,
                                GDM_AUTOLOGIN_USER,
                                name)
                }
        case "kdm":
                if fileIsExist(ETC_KDM_CONFIG) {
                        opUtils.WriteKeyToKeyFile(ETC_KDM_CONFIG,
                                KDM_AUTOLOGIN_GROUP,
                                KDM_AUTOLOGIN_ENABLE,
                                true)
                        opUtils.WriteKeyToKeyFile(ETC_KDM_CONFIG,
                                KDM_AUTOLOGIN_GROUP,
                                KDM_AUTOLOGIN_USER,
                                name)
                } else if fileIsExist(USER_KDM_CONFIG) {
                        opUtils.WriteKeyToKeyFile(ETC_KDM_CONFIG,
                                KDM_AUTOLOGIN_GROUP,
                                KDM_AUTOLOGIN_ENABLE,
                                true)
                        opUtils.WriteKeyToKeyFile(USER_KDM_CONFIG,
                                KDM_AUTOLOGIN_GROUP,
                                KDM_AUTOLOGIN_USER,
                                name)
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
                        v, ok := opUtils.ReadKeyFromKeyFile(ETC_LIGHTDM_CONFIG,
                                LIGHTDM_AUTOLOGIN_GROUP,
                                LIGHTDM_AUTOLOGIN_USER,
                                "")
                        if ok && v.(string) == username {
                                return true
                        }
                }
        case "gdm":
                if fileIsExist(ETC_GDM_CONFIG) {
                        v, ok := opUtils.ReadKeyFromKeyFile(ETC_GDM_CONFIG,
                                GDM_AUTOLOGIN_GROUP,
                                GDM_AUTOLOGIN_USER,
                                "")
                        if ok && v.(string) == username {
                                return true
                        }
                }
        case "kdm":
                if fileIsExist(ETC_KDM_CONFIG) {
                        v, ok := opUtils.ReadKeyFromKeyFile(ETC_KDM_CONFIG,
                                KDM_AUTOLOGIN_GROUP,
                                KDM_AUTOLOGIN_USER,
                                "")
                        if ok && v.(string) == username {
                                return true
                        }
                } else if fileIsExist(USER_KDM_CONFIG) {
                        v, ok := opUtils.ReadKeyFromKeyFile(USER_KDM_CONFIG,
                                KDM_AUTOLOGIN_GROUP,
                                KDM_AUTOLOGIN_USER,
                                "")
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
                logObject.Warningf("ReadFile '%s' failed: %s",
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
